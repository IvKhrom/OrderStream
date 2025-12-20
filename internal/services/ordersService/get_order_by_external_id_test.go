package ordersService

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/services/ordersService/mocks"
	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"

	"github.com/ivkhr/orderstream/internal/models"
)

type GetOrderByExternalIDSuite struct {
	suite.Suite
	ctx     context.Context
	storage *mocks.MockOrdersStorage
	svc     *OrdersService
}

func (s *GetOrderByExternalIDSuite) SetupTest() {
	s.ctx = context.Background()
	s.storage = mocks.NewMockOrdersStorage(s.T())
	s.svc = NewOrdersService(s.storage, nil, nil, 0)
}

func (s *GetOrderByExternalIDSuite) TestOk_DelegatesToStorage() {
	userID := uuid.New()
	externalID := "ext-1"
	want := &models.Order{OrderID: uuid.New(), UserID: userID, Status: "new"}

	s.storage.EXPECT().
		GetByExternalID(s.ctx, externalID, userID).
		Return(want, nil)

	got, err := s.svc.GetOrderByExternalID(s.ctx, externalID, userID)
	assert.NilError(s.T(), err)
	assert.Equal(s.T(), want, got)
}

func (s *GetOrderByExternalIDSuite) TestErr_Propagates() {
	userID := uuid.New()
	externalID := "ext-404"
	wantErr := errors.New("db error")

	s.storage.EXPECT().
		GetByExternalID(s.ctx, externalID, userID).
		Return((*models.Order)(nil), wantErr)

	got, err := s.svc.GetOrderByExternalID(s.ctx, externalID, userID)
	assert.ErrorIs(s.T(), err, wantErr)
	assert.Equal(s.T(), (*models.Order)(nil), got)
}

func TestGetOrderByExternalIDSuite(t *testing.T) {
	suite.Run(t, new(GetOrderByExternalIDSuite))
}
