package ordersService

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/models"
	"github.com/ivkhr/orderstream/internal/services/ordersService/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
)

type OrdersServiceUpsertSuite struct {
	suite.Suite

	ctx context.Context

	storage *mocks.MockOrdersStorage
	pub     *mocks.MockOrdersEventsPublisher
	ackReg  *mocks.MockAckWaitRegistry
}

func (s *OrdersServiceUpsertSuite) SetupTest() {
	s.ctx = context.Background()
	s.storage = mocks.NewMockOrdersStorage(s.T())
	s.pub = mocks.NewMockOrdersEventsPublisher(s.T())
	s.ackReg = mocks.NewMockAckWaitRegistry(s.T())
}

func (s *OrdersServiceUpsertSuite) newSvc(ack AckWaitRegistry, timeout time.Duration) *OrdersService {
	return NewOrdersService(s.storage, s.pub, ack, timeout)
}

func (s *OrdersServiceUpsertSuite) TestCreate_NoAck_Success() {
	svc := s.newSvc(nil, 0)

	userID := uuid.New()
	payload := json.RawMessage(`{"id":"ext-1","items":[1]}`)

	s.storage.EXPECT().
		Create(s.ctx, mock.Anything).
		Run(func(_ context.Context, ord *models.Order) {
			assert.Equal(s.T(), userID, ord.UserID)
			assert.Equal(s.T(), "new", ord.Status)
			assert.DeepEqual(s.T(), payload, ord.Payload)
			assert.Check(s.T(), ord.OrderID != uuid.Nil)
			assert.Check(s.T(), ord.Bucket >= 0 && ord.Bucket < 4)
		}).
		Return(nil)

	s.pub.EXPECT().
		Publish(s.ctx, mock.Anything).
		Return(nil)

	gotID, gotStatus, err := svc.UpsertOrder(s.ctx, "0", userID, "", payload)
	assert.NilError(s.T(), err)
	assert.Check(s.T(), gotID != "")
	assert.Equal(s.T(), "created", gotStatus)
}

func (s *OrdersServiceUpsertSuite) TestCreate_StorageError() {
	svc := s.newSvc(nil, 0)

	s.storage.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(errors.New("db error"))

	_, _, err := svc.UpsertOrder(s.ctx, "0", uuid.New(), "", json.RawMessage(`{}`))
	assert.Check(s.T(), err != nil)
}

func (s *OrdersServiceUpsertSuite) TestCreate_PublishError() {
	svc := s.newSvc(nil, 0)

	s.storage.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(nil)

	wantErr := errors.New("kafka error")
	s.pub.EXPECT().
		Publish(s.ctx, mock.Anything).
		Return(wantErr)

	_, _, err := svc.UpsertOrder(s.ctx, "0", uuid.New(), "", json.RawMessage(`{}`))
	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *OrdersServiceUpsertSuite) TestCreate_WaitAck_Success() {
	ch := make(chan struct{}, 1)
	cleanupCalled := false
	cleanup := func() { cleanupCalled = true }

	s.ackReg.EXPECT().
		Register(mock.Anything).
		Return((<-chan struct{})(ch), cleanup)

	s.storage.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(nil)

	s.pub.EXPECT().
		Publish(s.ctx, mock.Anything).
		Run(func(_ context.Context, _ []byte) { ch <- struct{}{} }).
		Return(nil)

	svc := s.newSvc(s.ackReg, 2*time.Second)
	_, _, err := svc.UpsertOrder(s.ctx, "0", uuid.New(), "", json.RawMessage(`{"id":"ext"}`))
	assert.NilError(s.T(), err)
	assert.Check(s.T(), cleanupCalled)
}

func (s *OrdersServiceUpsertSuite) TestCreate_WaitAck_Timeout() {
	ch := make(chan struct{}, 1)
	cleanup := func() {}

	s.ackReg.EXPECT().
		Register(mock.Anything).
		Return((<-chan struct{})(ch), cleanup)

	s.storage.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(nil)

	s.pub.EXPECT().
		Publish(s.ctx, mock.Anything).
		Return(nil)

	svc := s.newSvc(s.ackReg, 5*time.Millisecond)
	_, _, err := svc.UpsertOrder(s.ctx, "0", uuid.New(), "", json.RawMessage(`{}`))
	assert.ErrorIs(s.T(), err, context.DeadlineExceeded)
}

func (s *OrdersServiceUpsertSuite) TestCreate_WaitAck_CtxCanceled() {
	ch := make(chan struct{}, 1)
	cleanup := func() {}

	ctx, cancel := context.WithCancel(context.Background())

	s.ackReg.EXPECT().
		Register(mock.Anything).
		Return((<-chan struct{})(ch), cleanup)

	s.storage.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(nil)

	s.pub.EXPECT().
		Publish(mock.Anything, mock.Anything).
		Run(func(_ context.Context, _ []byte) { cancel() }).
		Return(nil)

	svc := s.newSvc(s.ackReg, 2*time.Second)
	_, _, err := svc.UpsertOrder(ctx, "0", uuid.New(), "", json.RawMessage(`{}`))
	assert.ErrorIs(s.T(), err, context.Canceled)
}

func (s *OrdersServiceUpsertSuite) TestUpdate_InvalidOrderID() {
	svc := s.newSvc(nil, 0)
	_, _, err := svc.UpsertOrder(s.ctx, "bad-uuid", uuid.New(), "", json.RawMessage(`{}`))
	assert.Check(s.T(), err != nil)
}

func (s *OrdersServiceUpsertSuite) TestUpdate_NotFound() {
	svc := s.newSvc(nil, 0)
	oid := uuid.New()

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return((*models.Order)(nil), errors.New("db: not found"))

	_, _, err := svc.UpsertOrder(s.ctx, oid.String(), uuid.New(), "", json.RawMessage(`{}`))
	assert.ErrorIs(s.T(), err, ErrNotFound)
}

func (s *OrdersServiceUpsertSuite) TestUpdate_DeletedConflict() {
	svc := s.newSvc(nil, 0)
	oid := uuid.New()
	existing := &models.Order{
		OrderID: oid,
		UserID:  uuid.New(),
		Status:  "deleted",
		Payload: json.RawMessage(`{"id":"ext"}`),
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	gotID, gotStatus, err := svc.UpsertOrder(s.ctx, oid.String(), uuid.New(), "", json.RawMessage(`{}`))
	assert.ErrorIs(s.T(), err, ErrDeletedConflict)
	assert.Equal(s.T(), oid.String(), gotID)
	assert.Equal(s.T(), "deleted", gotStatus)
}

func (s *OrdersServiceUpsertSuite) TestDelete_Success() {
	svc := s.newSvc(nil, 0)
	oid := uuid.New()
	existing := &models.Order{
		OrderID: oid,
		UserID:  uuid.New(),
		Status:  "new",
		Payload: json.RawMessage(`{"id":"ext"}`),
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	s.storage.EXPECT().
		DeleteOrder(s.ctx, oid.String()).
		Return(nil)

	s.pub.EXPECT().
		Publish(s.ctx, mock.Anything).
		Return(nil)

	gotID, gotStatus, err := svc.UpsertOrder(s.ctx, oid.String(), uuid.New(), "deleted", json.RawMessage(`{}`))
	assert.NilError(s.T(), err)
	assert.Equal(s.T(), oid.String(), gotID)
	assert.Equal(s.T(), "deleted", gotStatus)
}

func (s *OrdersServiceUpsertSuite) TestDelete_StorageError() {
	svc := s.newSvc(nil, 0)
	oid := uuid.New()
	existing := &models.Order{
		OrderID: oid,
		UserID:  uuid.New(),
		Status:  "new",
		Payload: json.RawMessage(`{"id":"ext"}`),
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	wantErr := errors.New("db error")
	s.storage.EXPECT().
		DeleteOrder(s.ctx, oid.String()).
		Return(wantErr)

	_, _, err := svc.UpsertOrder(s.ctx, oid.String(), uuid.New(), "deleted", json.RawMessage(`{}`))
	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *OrdersServiceUpsertSuite) TestDelete_PublishError() {
	svc := s.newSvc(nil, 0)
	oid := uuid.New()
	existing := &models.Order{
		OrderID: oid,
		UserID:  uuid.New(),
		Status:  "new",
		Payload: json.RawMessage(`{"id":"ext"}`),
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	s.storage.EXPECT().
		DeleteOrder(s.ctx, oid.String()).
		Return(nil)

	wantErr := errors.New("kafka error")
	s.pub.EXPECT().
		Publish(s.ctx, mock.Anything).
		Return(wantErr)

	_, _, err := svc.UpsertOrder(s.ctx, oid.String(), uuid.New(), "deleted", json.RawMessage(`{}`))
	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *OrdersServiceUpsertSuite) TestUpdate_Success() {
	svc := s.newSvc(nil, 0)
	oid := uuid.New()
	existing := &models.Order{
		OrderID: oid,
		UserID:  uuid.New(),
		Status:  "new",
		Payload: json.RawMessage(`{"id":"ext"}`),
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	s.storage.EXPECT().
		Update(s.ctx, mock.Anything).
		Return(nil)

	s.pub.EXPECT().
		Publish(s.ctx, mock.Anything).
		Return(nil)

	gotID, gotStatus, err := svc.UpsertOrder(s.ctx, oid.String(), uuid.New(), "", json.RawMessage(`{"id":"ext","v":1}`))
	assert.NilError(s.T(), err)
	assert.Equal(s.T(), oid.String(), gotID)
	assert.Equal(s.T(), "updated", gotStatus)
}

func (s *OrdersServiceUpsertSuite) TestUpdate_UpdateError() {
	svc := s.newSvc(nil, 0)
	oid := uuid.New()
	existing := &models.Order{
		OrderID: oid,
		UserID:  uuid.New(),
		Status:  "new",
		Payload: json.RawMessage(`{"id":"ext"}`),
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	wantErr := errors.New("db error")
	s.storage.EXPECT().
		Update(s.ctx, mock.Anything).
		Return(wantErr)

	_, _, err := svc.UpsertOrder(s.ctx, oid.String(), uuid.New(), "", json.RawMessage(`{"id":"ext","v":1}`))
	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *OrdersServiceUpsertSuite) TestUpdate_PublishError() {
	svc := s.newSvc(nil, 0)
	oid := uuid.New()
	existing := &models.Order{
		OrderID: oid,
		UserID:  uuid.New(),
		Status:  "new",
		Payload: json.RawMessage(`{"id":"ext"}`),
	}

	s.storage.EXPECT().
		GetByID(s.ctx, oid).
		Return(existing, nil)

	s.storage.EXPECT().
		Update(s.ctx, mock.Anything).
		Return(nil)

	wantErr := errors.New("kafka error")
	s.pub.EXPECT().
		Publish(s.ctx, mock.Anything).
		Return(wantErr)

	_, _, err := svc.UpsertOrder(s.ctx, oid.String(), uuid.New(), "", json.RawMessage(`{"id":"ext","v":1}`))
	assert.ErrorIs(s.T(), err, wantErr)
}

func TestOrdersServiceUpsertSuite(t *testing.T) {
	suite.Run(t, new(OrdersServiceUpsertSuite))
}


