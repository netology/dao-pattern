package models

// ProductID is value object
type ProductID int

// Product is an entity
type Product struct {
	ID    ProductID
	Price Money
}
