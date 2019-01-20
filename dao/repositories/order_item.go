//go:generate mockgen -source=order_item.go -package repositories -destination order_item_mock.go

package repositories

import (
	"database/sql"
	"github.com/netology/godesignpatterns/dao/models"
)

// OrderItemRepository is a repository
type OrderItemRepository interface {
	GetByOrderID(orderID models.OrderID) ([]*models.OrderItem, error)
	SaveWithTransaction(tx *sql.Tx, orderItem *models.OrderItem) error
	Save(orderItem *models.OrderItem) error
}
