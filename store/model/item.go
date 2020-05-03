package model

//ite model
type Item struct {
	ID       string  `json:"_id" bson:"_id"`
	Name     string  `json:"name" bson:"name"`
	Price    float64 `json:"price" bson:"price"`
	Quantity int64   `json:"quantity" bson:"quantity"`
}
