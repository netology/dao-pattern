package models

// OrderID is a value object
type OrderID int

// Order is an entity
type Order struct {
	ID         OrderID
	CustomerID CustomerID
	Amount     Money
	Items      []*OrderItem
}
