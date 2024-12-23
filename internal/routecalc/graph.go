package routecalc

import (
	"context"

	"github.com/anunayjoshi29/token-server/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PoolEdge struct {
	TokenA   string
	TokenB   string
	ReserveA float64
	ReserveB float64
}

type Graph struct {
	adjacency map[string][]PoolEdge
}

func NewGraph() *Graph {
	return &Graph{
		adjacency: make(map[string][]PoolEdge),
	}
}

// BuildGraph loads all pools from MongoDB and constructs the graph.
func (g *Graph) BuildGraph(ctx context.Context, poolsColl *mongo.Collection) error {
	cursor, err := poolsColl.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p models.Pool
		if err := cursor.Decode(&p); err != nil {
			return err
		}

		// Add bidirectional edges
		g.adjacency[p.Token0] = append(g.adjacency[p.Token0], PoolEdge{
			TokenA: p.Token0, TokenB: p.Token1,
			ReserveA: p.Reserve0, ReserveB: p.Reserve1,
		})
		g.adjacency[p.Token1] = append(g.adjacency[p.Token1], PoolEdge{
			TokenA: p.Token1, TokenB: p.Token0,
			ReserveA: p.Reserve1, ReserveB: p.Reserve0,
		})
	}
	return nil
}
