package models

// CustomerID is value object
type CustomerID int

// Customer is an entity
type Customer struct {
	ID      CustomerID
	Balance Money
}
