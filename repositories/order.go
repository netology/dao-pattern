//go:generate mockgen -source=order.go -package repositories -destination order_mock.go

package repositories

import (
	"github.com/netology/dao-pattern/models"
)

// OrderRepository is a repository
type OrderRepository interface {
	GetByID(orderID models.OrderID) (*models.Order, error)
	Save(order *models.Order) error
}
