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

type GetOrderByIDSuite struct {
	suite.Suite
	ctx     context.Context
	storage *mocks.MockOrdersStorage
	svc     *OrdersService
}

func (s *GetOrderByIDSuite) SetupTest() {
	s.ctx = context.Background()
	s.storage = mocks.NewMockOrdersStorage(s.T())
	s.svc = NewOrdersService(s.storage, nil, nil, 0)
}

func (s *GetOrderByIDSuite) TestOk_DelegatesToStorage() {
	id := uuid.New()
	want := &models.Order{OrderID: id, UserID: uuid.New(), Status: "new"}

	s.storage.EXPECT().
		GetByID(s.ctx, id).
		Return(want, nil)

	got, err := s.svc.GetOrderByID(s.ctx, id)
	assert.NilError(s.T(), err)
	assert.Equal(s.T(), want, got)
}

func (s *GetOrderByIDSuite) TestErr_Propagates() {
	id := uuid.New()
	wantErr := errors.New("db error")

	s.storage.EXPECT().
		GetByID(s.ctx, id).
		Return((*models.Order)(nil), wantErr)

	got, err := s.svc.GetOrderByID(s.ctx, id)
	assert.ErrorIs(s.T(), err, wantErr)
	assert.Equal(s.T(), (*models.Order)(nil), got)
}

func TestGetOrderByIDSuite(t *testing.T) {
	suite.Run(t, new(GetOrderByIDSuite))
}


