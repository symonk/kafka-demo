package models

// An encapsulation of a simple data type, generated
// and dispatched int he kafka message payload
// also subsequently unmarshalled on the consumer side.
type Superhero struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Power int    `json:"power"`
	Melee bool   `json:"melee"`
}
