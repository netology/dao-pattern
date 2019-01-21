package postgresql

import (
	"github.com/netology/godesignpatterns/dao/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
	"testing"
)

func TestOrderItem_GetByOrderID(t *testing.T) {
	expectedOrderID := models.OrderID(1020)

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectPrepare("SELECT order_item_id, order_id, product_id, quantity, price, currency FROM order_items").ExpectQuery().
			WillReturnRows(sqlmock.NewRows([]string{"order_item_id", "order_id", "product_id", "quantity", "price", "currency"}).
				AddRow(1, expectedOrderID, models.ProductID(2), 1, 3.0, "usd"))

		orderRepository := NewOrderItemRepository(db)
		orderItems, err := orderRepository.GetByOrderID(expectedOrderID)
		require.NoError(t, err)
		require.NotEmpty(t, orderItems)
		require.Equal(t, expectedOrderID, orderItems[0].OrderID)
	})

	t.Run("errors", func(t *testing.T) {
		dummyError := errors.New("dummy-error")

		t.Run("prepare and query return an error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectPrepare("SELECT order_item_id, order_id, product_id, quantity, price, currency FROM order_items").ExpectQuery().WillReturnError(dummyError)

			orderRepository := NewOrderItemRepository(db)
			orderItems, err := orderRepository.GetByOrderID(expectedOrderID)
			require.Empty(t, orderItems)
			require.Error(t, err)
			require.Equal(t, errors.Cause(err), dummyError)
		})

	})

}

func TestOrderItem_Save(t *testing.T) {
	expectedOrderID := models.OrderID(1020)
	expectedInput := &models.OrderItem{1, expectedOrderID, models.ProductID(2), 1, models.Money{3, models.USD}}

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectPrepare(`INSERT INTO order_items \(order_id, product_id, quantity, price, currency\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING order_item_id`).
			ExpectQuery().
			WithArgs(expectedOrderID, models.ProductID(2), 1, 3.0, "usd").
			WillReturnRows(sqlmock.NewRows([]string{"order_item_id"}).AddRow(1))

		orderRepository := NewOrderItemRepository(db)
		err = orderRepository.Save(expectedInput)
		require.NoError(t, err)
	})

	t.Run("errors", func(t *testing.T) {
		dummyError := errors.New("dummy-error")

		t.Run("prepare and query return an error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectPrepare(`INSERT INTO order_items \(order_id, product_id, quantity, price, currency\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING order_item_id`).
				ExpectQuery().
				WithArgs(expectedOrderID, models.ProductID(2), 1, 3.0, "usd").WillReturnError(dummyError)

			orderRepository := NewOrderItemRepository(db)
			err = orderRepository.Save(expectedInput)
			require.Error(t, err)
			require.Equal(t, errors.Cause(err), dummyError)
		})

	})
}

func TestOrderItem_SaveWithTransaction(t *testing.T) {
	expectedOrderID := models.OrderID(1020)

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectPrepare(`INSERT INTO order_items \(order_id, product_id, quantity, price, currency\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING order_item_id`).
			ExpectQuery().
			WithArgs(expectedOrderID, models.ProductID(2), 1, 3.0, "usd").
			WillReturnRows(sqlmock.NewRows([]string{"order_item_id"}).AddRow(1))

		orderRepository := NewOrderItemRepository(db)
		tx, err := db.Begin()
		require.NoError(t, err)

		err = orderRepository.SaveWithTransaction(tx, &models.OrderItem{1, expectedOrderID, models.ProductID(2), 1, models.Money{3, models.USD}})
		require.NoError(t, err)
	})
}
