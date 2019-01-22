// +build integration

package postgresql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/netology/dao-pattern/models"
	"github.com/netology/dao-pattern/repositories/postgresql"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntegration(t *testing.T) {
	db := postgresql.NewConnection()
	defer db.Close()

	t.Run("success commint", func(t *testing.T) {
		orderItemRepository := postgresql.NewOrderItemRepository(db)
		orderRepository := postgresql.NewOrderRepository(db, orderItemRepository)
		orderEntity := &models.Order{
			CustomerID: 1,
			Amount:     models.Money{20, models.USD},
			Items: []*models.OrderItem{
				{
					ProductID: 1,
					Quantity:  1,
					Price:     models.Money{8, models.USD},
				},
				{
					ProductID: 2,
					Quantity:  26,
					Price:     models.Money{12, models.USD},
				},
			},
		}
		err := orderRepository.Save(orderEntity)
		require.NoError(t, err)

		order, err := orderRepository.GetByID(orderEntity.ID)
		require.NoError(t, err)
		require.NotNil(t, order)
		require.Len(t, order.Items, 2)
	})

	t.Run("failure and rollback", func(t *testing.T) {
		orderItemRepository := postgresql.NewOrderItemRepository(db)
		orderRepository := postgresql.NewOrderRepository(db, orderItemRepository)
		orderEntity := &models.Order{
			CustomerID: 1,
			Amount:     models.Money{20, models.USD},
			Items: []*models.OrderItem{
				{
					ProductID: 1,
					Quantity:  1,
					Price:     models.Money{-8, models.USD},
				},
				{
					ProductID: 2,
					Quantity:  26,
					Price:     models.Money{12, models.USD},
				},
			},
		}
		err := orderRepository.Save(orderEntity)
		require.Error(t, err)

		order, err := orderRepository.GetByID(orderEntity.ID)
		require.Error(t, err)
		require.Equal(t, errors.Cause(err), sql.ErrNoRows)
		require.Nil(t, order)
	})

}
