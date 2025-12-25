package worker

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/services/worker/internal/models"
	"github.com/ivkhr/orderstream/services/worker/internal/services/worker/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WorkerHandleOrderEventSuite struct {
	suite.Suite

	ctx context.Context

	storage *mocks.MockStorage
	svc     *Service
}

func (s *WorkerHandleOrderEventSuite) SetupTest() {
	s.ctx = context.Background()
	s.storage = mocks.NewMockStorage(s.T())
	s.svc = New(s.storage)
}

func (s *WorkerHandleOrderEventSuite) TestNilEvent_Error() {
	_, err := s.svc.HandleOrderEvent(s.ctx, nil)
	assert.Error(s.T(), err)
}

func (s *WorkerHandleOrderEventSuite) TestBadUserID_Error() {
	_, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID: "0",
		UserID:  "bad",
		Payload: json.RawMessage(`{"id":"x"}`),
	})
	assert.Error(s.T(), err)
}

func (s *WorkerHandleOrderEventSuite) TestCreate_BackCompat_OrderIDZero_Success() {
	userID := uuid.New()
	payload := json.RawMessage(`{"id":"ext-1","amount":12.5}`)

	s.storage.EXPECT().
		Create(s.ctx, mock.Anything).
		Run(func(_ context.Context, ord *models.Order) {
			assert.Equal(s.T(), userID, ord.UserID)
			assert.Equal(s.T(), "done", ord.Status)
			assert.Equal(s.T(), payload, ord.Payload)
			assert.NotEqual(s.T(), uuid.Nil, ord.OrderID)
			assert.True(s.T(), ord.Bucket >= 0 && ord.Bucket < 4)
			assert.Equal(s.T(), 12.5, ord.Amount)
		}).
		Return(nil)

	res, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID: "0",
		UserID:  userID.String(),
		Payload: payload,
	})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), res)
	assert.Equal(s.T(), "created", res.Status)
	assert.Equal(s.T(), "done", res.OrderStatus)
	assert.NotEmpty(s.T(), res.OrderID)
}

func (s *WorkerHandleOrderEventSuite) TestCreate_BackCompat_OrderIDZero_CreateError() {
	userID := uuid.New()
	wantErr := errors.New("db create error")

	s.storage.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(wantErr)

	_, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID: "0",
		UserID:  userID.String(),
		Payload: json.RawMessage(`{"id":"ext-err","amount":1}`),
	})
	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *WorkerHandleOrderEventSuite) TestBadOrderIDParse_Error() {
	_, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID: "bad-uuid",
		UserID:  uuid.New().String(),
		Payload: json.RawMessage(`{"id":"x"}`),
	})
	assert.Error(s.T(), err)
}

func (s *WorkerHandleOrderEventSuite) TestGetByID_Error_ThenUpsertCreated() {
	oid := uuid.New()
	userID := uuid.New()
	payload := json.RawMessage(`{"id":"ext-2","amount":7}`)

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return((*models.Order)(nil), errors.New("not found"))

	s.storage.EXPECT().
		Upsert(s.ctx, mock.Anything).
		Run(func(_ context.Context, ord *models.Order) {
			assert.Equal(s.T(), oid, ord.OrderID)
			assert.Equal(s.T(), userID, ord.UserID)
			assert.Equal(s.T(), "done", ord.Status)
			assert.Equal(s.T(), payload, ord.Payload)
			assert.Equal(s.T(), 7.0, ord.Amount)
		}).
		Return(true, nil)

	// To cover branch where event.ExternalID is empty but payload has id (closure in slog.Info).
	res, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID:    oid.String(),
		UserID:     userID.String(),
		ExternalID: "",
		Payload:    payload,
		Status:     "",
	})
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "created", res.Status)
	assert.Equal(s.T(), "done", res.OrderStatus)
	assert.Equal(s.T(), oid.String(), res.OrderID)
}

func (s *WorkerHandleOrderEventSuite) TestGetByID_Error_ThenUpsertError() {
	oid := uuid.New()
	userID := uuid.New()

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return((*models.Order)(nil), errors.New("not found"))

	wantErr := errors.New("db upsert error")
	s.storage.EXPECT().
		Upsert(s.ctx, mock.Anything).
		Return(false, wantErr)

	_, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID: oid.String(),
		UserID:  userID.String(),
		Payload: json.RawMessage(`{"id":"x"}`),
	})
	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *WorkerHandleOrderEventSuite) TestExistingDeleted_RejectsUpdate() {
	oid := uuid.New()
	userID := uuid.New()
	existing := &models.Order{
		OrderID:   oid,
		UserID:    userID,
		Amount:    1,
		Status:    "deleted",
		Payload:   json.RawMessage(`{"id":"ext"}`),
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Minute),
		Bucket:    1,
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	res, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID: oid.String(),
		UserID:  userID.String(),
		Status:  "processing",
		Payload: json.RawMessage(`{"id":"ext","amount":999}`),
	})
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "deleted", res.Status)
	assert.Equal(s.T(), "deleted", res.OrderStatus)
}

func (s *WorkerHandleOrderEventSuite) TestDelete_Success() {
	oid := uuid.New()
	userID := uuid.New()
	existing := &models.Order{
		OrderID:   oid,
		UserID:    userID,
		Status:    "done",
		Payload:   json.RawMessage(`{"id":"ext"}`),
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Minute),
		Bucket:    2,
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)
	s.storage.EXPECT().
		DeleteOrder(s.ctx, oid.String()).
		Return(nil)

	// Cover mismatch warning branch (event.ExternalID != payload id).
	res, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID:    oid.String(),
		UserID:     userID.String(),
		Status:     "deleted",
		ExternalID: "from-event",
		Payload:    json.RawMessage(`{"id":"from-payload"}`),
	})
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "deleted", res.Status)
	assert.Equal(s.T(), "deleted", res.OrderStatus)
}

func (s *WorkerHandleOrderEventSuite) TestDelete_Error() {
	oid := uuid.New()
	userID := uuid.New()
	existing := &models.Order{OrderID: oid, UserID: userID, Status: "done"}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	wantErr := errors.New("delete error")
	s.storage.EXPECT().
		DeleteOrder(s.ctx, oid.String()).
		Return(wantErr)

	_, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID: oid.String(),
		UserID:  userID.String(),
		Status:  "deleted",
		Payload: json.RawMessage(`{"id":"x"}`),
	})
	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *WorkerHandleOrderEventSuite) TestUpdate_WithPayload_Success() {
	oid := uuid.New()
	userID := uuid.New()
	oldPayload := json.RawMessage(`{"id":"old","amount":1}`)
	existing := &models.Order{
		OrderID:   oid,
		UserID:    userID,
		Amount:    1,
		Status:    "new",
		Payload:   oldPayload,
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Minute),
		Bucket:    3,
	}
	newPayload := json.RawMessage(`{"id":"old","amount":10}`)

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	s.storage.EXPECT().
		Update(s.ctx, mock.Anything).
		Run(func(_ context.Context, ord *models.Order) {
			assert.Equal(s.T(), newPayload, ord.Payload)
			assert.Equal(s.T(), 10.0, ord.Amount)
			assert.False(s.T(), ord.UpdatedAt.IsZero())
		}).
		Return(nil)

	s.storage.EXPECT().
		UpdateStatus(s.ctx, oid.String(), "done").
		Return(nil)

	res, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID: oid.String(),
		UserID:  userID.String(),
		Status:  "processing",
		Payload: newPayload,
	})
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "updated", res.Status)
	assert.Equal(s.T(), "done", res.OrderStatus)
	assert.Equal(s.T(), 10.0, res.Amount)
}

func (s *WorkerHandleOrderEventSuite) TestUpdate_WithoutPayload_KeepsExisting() {
	oid := uuid.New()
	userID := uuid.New()
	oldPayload := json.RawMessage(`{"id":"old","amount":2}`)
	existing := &models.Order{
		OrderID:   oid,
		UserID:    userID,
		Amount:    2,
		Status:    "new",
		Payload:   oldPayload,
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Minute),
		Bucket:    0,
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	s.storage.EXPECT().
		Update(s.ctx, mock.Anything).
		Run(func(_ context.Context, ord *models.Order) {
			assert.Equal(s.T(), oldPayload, ord.Payload)
			assert.Equal(s.T(), 2.0, ord.Amount)
		}).
		Return(nil)

	// Even if UpdateStatus fails, handler ignores it; still cover the call.
	s.storage.EXPECT().
		UpdateStatus(s.ctx, oid.String(), "done").
		Return(errors.New("ignored"))

	res, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID: oid.String(),
		UserID:  userID.String(),
		Status:  "",
		Payload: nil,
	})
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "updated", res.Status)
	assert.Equal(s.T(), "done", res.OrderStatus)
	assert.Equal(s.T(), oldPayload, res.Payload)
}

func (s *WorkerHandleOrderEventSuite) TestUpdate_Error() {
	oid := uuid.New()
	userID := uuid.New()
	existing := &models.Order{
		OrderID: oid,
		UserID:  userID,
		Status:  "new",
		Payload: json.RawMessage(`{"id":"x"}`),
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	wantErr := errors.New("update error")
	s.storage.EXPECT().
		Update(s.ctx, mock.Anything).
		Return(wantErr)

	_, err := s.svc.HandleOrderEvent(s.ctx, &models.OrderEvent{
		OrderID: oid.String(),
		UserID:  userID.String(),
		Payload: json.RawMessage(`{"id":"x","amount":5}`),
	})
	assert.ErrorIs(s.T(), err, wantErr)
}

func TestWorkerHandleOrderEventSuite(t *testing.T) {
	suite.Run(t, new(WorkerHandleOrderEventSuite))
}
