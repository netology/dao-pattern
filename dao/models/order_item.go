package models

// OrderItem is an entity
type OrderItem struct {
	ID        int64
	OrderID   OrderID
	ProductID ProductID
	Quantity  int
	Price     Money
}
