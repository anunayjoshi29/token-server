package main

import (
	"context"
	"log"

	"github.com/anunayjoshi29/token-server/internal/db"
	"github.com/anunayjoshi29/token-server/internal/routecalc"
	"github.com/anunayjoshi29/token-server/internal/server"
)

func main() {
	// Connect to MongoDB
	mongoDB, err := db.Connect("mongodb://localhost:27017", "euclid-db")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Build graph
	graph := routecalc.NewGraph()
	if err := graph.BuildGraph(context.Background(), mongoDB.Pools); err != nil {
		log.Fatal("Failed to build graph:", err)
	}

	finder := routecalc.NewFinder(graph)
	cache := routecalc.NewCache()

	srv := server.NewServer(finder, cache)

	log.Println("Server listening on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
