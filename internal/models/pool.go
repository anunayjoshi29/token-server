package models

type Pool struct {
	ID       string  `bson:"_id,omitempty"`
	Token0   string  `bson:"token0"`
	Token1   string  `bson:"token1"`
	Reserve0 float64 `bson:"reserve0"`
	Reserve1 float64 `bson:"reserve1"`
}
