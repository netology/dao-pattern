package postgresql

import (
	"github.com/golang/mock/gomock"
	"github.com/netology/godesignpatterns/dao/models"
	"github.com/netology/godesignpatterns/dao/repositories"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
	"testing"
)

func TestOrder_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedID := models.OrderID(123)

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectPrepare(`INSERT INTO orders \(customer_id, amount, currency\) VALUES \(\$1, \$2, \$3\) RETURNING order_id`).
			ExpectQuery().
			WithArgs(1, float64(1), "usd").
			WillReturnRows(sqlmock.NewRows([]string{"order_id"}).AddRow(expectedID))
		mock.ExpectCommit()

		ctrl := gomock.NewController(t)
		mockOrderItemRepository := repositories.NewMockOrderItemRepository(ctrl)
		mockOrderItemRepository.EXPECT().SaveWithTransaction(gomock.Any(), gomock.Any()).Return(nil)

		orderRepository := NewOrderRepository(db, mockOrderItemRepository)

		orderEntity := &models.Order{
			CustomerID: models.CustomerID(1),
			Amount:     models.Money{1, models.USD},
			Items: []*models.OrderItem{
				{
					ProductID: 1,
					Quantity:  1,
					Price:     models.Money{8, models.USD},
				},
			},
		}
		err = orderRepository.Save(orderEntity)
		require.NoError(t, err)
		require.Equal(t, expectedID, orderEntity.ID)

		ctrl.Finish()
	})

	t.Run("errors", func(t *testing.T) {
		dummyError := errors.New("dummy-error")

		t.Run("begin transaction return an error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectBegin().WillReturnError(dummyError)

			orderRepository := NewOrderRepository(db, nil)
			err = orderRepository.Save(&models.Order{})
			require.Error(t, err)
			require.Equal(t, errors.Cause(err), dummyError)
		})

		t.Run("prepare and query return an error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectBegin()
			mock.ExpectPrepare(`INSERT INTO orders \(customer_id, amount, currency\) VALUES \(\$1, \$2, \$3\) RETURNING order_id`).ExpectQuery().
				WithArgs(1, float64(1), "usd").WillReturnError(dummyError)

			orderRepository := NewOrderRepository(db, nil)
			err = orderRepository.Save(&models.Order{
				CustomerID: models.CustomerID(1),
				Amount:     models.Money{1, models.USD},
			})
			require.Error(t, err)
			require.Equal(t, errors.Cause(err), dummyError)
		})

		t.Run("orderitem repository return an error", func(t *testing.T) {
			expectedID := models.OrderID(123)

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectBegin()
			mock.ExpectPrepare(`INSERT INTO orders \(customer_id, amount, currency\) VALUES \(\$1, \$2, \$3\) RETURNING order_id`).
				ExpectQuery().
				WithArgs(1, float64(1), "usd").
				WillReturnRows(sqlmock.NewRows([]string{"order_id"}).AddRow(expectedID))
			mock.ExpectRollback()

			ctrl := gomock.NewController(t)
			mockOrderItemRepository := repositories.NewMockOrderItemRepository(ctrl)
			mockOrderItemRepository.EXPECT().SaveWithTransaction(gomock.Any(), gomock.Any()).Return(dummyError)

			orderRepository := NewOrderRepository(db, mockOrderItemRepository)

			orderEntity := &models.Order{
				CustomerID: models.CustomerID(1),
				Amount:     models.Money{1, models.USD},
				Items: []*models.OrderItem{
					{
						ProductID: 1,
						Quantity:  1,
						Price:     models.Money{8, models.USD},
					},
				},
			}
			err = orderRepository.Save(orderEntity)

			ctrl.Finish()
			require.Error(t, err)
			require.Equal(t, errors.Cause(err), dummyError)
		})
	})

}

func TestOrder_GetByID(t *testing.T) {
	expectedOrderID := models.OrderID(1)

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectPrepare("SELECT order_id, customer_id, amount, currency FROM orders").ExpectQuery().
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "customer_id", "amount", "currency"}).AddRow(expectedOrderID, 2, 3.0, "usd"))

		ctrl := gomock.NewController(t)
		mockOrderItemRepository := repositories.NewMockOrderItemRepository(ctrl)

		mockOrderItemRepository.EXPECT().GetByOrderID(expectedOrderID).Return([]*models.OrderItem{}, nil)

		orderRepository := NewOrderRepository(db, mockOrderItemRepository)
		order, err := orderRepository.GetByID(expectedOrderID)
		require.NoError(t, err)
		require.NotNil(t, order)

		ctrl.Finish()
	})

	t.Run("errors", func(t *testing.T) {
		dummyError := errors.New("dummy-error")

		t.Run("prepare and query return an error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectPrepare("SELECT order_id, customer_id, amount, currency FROM orders").ExpectQuery().WillReturnError(dummyError)

			orderRepository := NewOrderRepository(db, nil)
			order, err := orderRepository.GetByID(expectedOrderID)
			require.Nil(t, order)
			require.Error(t, err)
			require.Equal(t, errors.Cause(err), dummyError)
		})

	})
}
