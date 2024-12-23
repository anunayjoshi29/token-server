# Token Route Optimization Server
## Overview

The **Token Route Optimization Server** is a Go-based application developed to efficiently handle token routing requests. It calculates and returns all possible routes between a specified `fromToken` and `toToken` for a given input amount, sorted by the highest expected output amount while adhering to specific constraints.

## Features

- **Comprehensive Route Discovery:** Identifies all possible token routes up to a maximum path length of 5.
- **Accurate Calculations:** Uses the constant product formula to compute expected output amounts for each route.
- **Performance Optimization:** Implements in-memory caching to enhance response times for repeated requests.
- **MongoDB Integration:** Efficiently retrieves pool data from a MongoDB database.

## Route Calculation

1. **Graph Representation:**
    - Tokens are nodes in a graph.
    - Pools between tokens are bidirectional edges with their respective reserves.

2. **Graph Traversal:**
    - Utilizes **Depth-First Search (DFS)** to explore all possible routes from `fromToken` to `toToken`.
    - Limits path length to 5 tokens. If no routes are found within this limit, it searches for the best available route exceeding the limit.

3. **Amount Out Calculation:**
    - For each route, calculates `expectedAmountOut` using the constant product formula:
      ```
      reserveIn * reserveOut = newReserveIn * newReserveOut
      ```
    - Ensures the product of reserves remains constant post-swap to determine the output amount.

4. **Route Sorting:**
    - Sorts all valid routes in descending order based on `expectedAmountOut`, presenting the most profitable routes first.

## File Logics

The important files of the project with their logics is as follows:

- `internal/routecalc/cache.go`: This file implements an in-memory cache to store and retrieve route calculations. It helps improve response times for repeated requests.

- `internal/routecalc/calculations.go`: This file maintains the Constant Product pool.

- `internal/routecalc/finder.go`: This file has the logic to find the route and calculate out price.

- `internal/routecalc/graph.go`: This file has the logic to build graph representation of the given pools using adjacency.

- `sample_pools_data.json`: This file is used in this project to verify and validate the project.

## Getting Started

### Prerequisites

- **Go:** Version 1.18 or higher. [Download Go](https://golang.org/dl/)
- **MongoDB:** Installed and running. [MongoDB Installation Guide](https://www.mongodb.com/docs/manual/installation/)

### Installation

1. **Clone the Repository:**
     ```bash
     git clone https://github.com/anunayjoshi29/token-server.git
     cd token-server
     ```

2. **Initialize Go Modules:**
     ```bash
     go mod tidy
     ```

### Importing Pool Data

1. **Prepare `pools_data.json`:**
    - Ensure it contains pool configurations with `token0`, `token1`, `reserve0`, and `reserve1`.

2. **Import into MongoDB:**
     ```bash
     mongoimport --uri="mongodb://localhost:27017" --db yourdb --collection pools --file /path/to/pools_data.json --jsonArray
     ```

3. **Verify Import:**
     ```bash
     mongo
     use yourdb
     db.pools.find().limit(5).pretty()
     ```

### Running the Server

1. **Build and Run:**
     ```bash
     go run ./cmd/server
     ```
     - The server will start on [http://localhost:8080](http://localhost:8080).

## API Usage

### Endpoint

- **URL:** `POST /routes`
- **Content-Type:** `application/json`

### Request Example

```json
{
    "fromToken": "TokenA",
    "toToken": "TokenB",
    "amountIn": 100.0
}
```

### Response Example

```json
{
    "routes": [
        {
            "path": [
                "TokenA",
                "TokenB"
            ],
            "expectedAmountOut": 199.80019980020006
        },
        {
            "path": [
                "TokenA",
                "TokenC",
                "TokenB"
            ],
            "expectedAmountOut": 99.67605282830482
        }
    ]
}
```

### Testing with curl

```bash
curl -X POST -H "Content-Type: application/json" \
         -d '{"fromToken":"TokenA","toToken":"TokenB","amountIn":100.0}' \
         http://localhost:8080/routes
```

