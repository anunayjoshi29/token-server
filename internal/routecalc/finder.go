package routecalc

import (
	"errors"
	"fmt"
	"sort"
)

// represent single route
type RouteResult struct {
	Path              []string
	ExpectedAmountOut float64
}

// finder for finding route
type Finder struct {
	g *Graph
}

// new finder instance
func NewFinder(g *Graph) *Finder {
	return &Finder{g: g}
}

// find all routes possible
// It first attempt to find routes with max depth of 5, if no routes are found, it will try to find routes with no depth limit
func (f *Finder) FindAllRoutes(fromToken, toToken string, amountIn float64) []RouteResult {
	var routes []RouteResult

	visited := make(map[string]bool)

	// Try to find routes with max depth of 5
	f.findRoutes(fromToken, toToken, visited, []string{fromToken}, &routes, 5, amountIn)

	if len(routes) == 0 {
		// If no routes are found, try to find routes with no depth limit
		f.findRoutes(fromToken, toToken, visited, []string{fromToken}, &routes, 0, amountIn)
	}

	// Sort routes by expected amount out
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].ExpectedAmountOut > routes[j].ExpectedAmountOut
	})

	return routes
}

// find routes between two tokens with max depth and amount
func (f *Finder) findRoutes(currentToken, targetToken string, visited map[string]bool, path []string, routes *[]RouteResult, maxDepth int, amountIn float64) {
	if maxDepth > 0 && len(path) > maxDepth {
		return
	}

	if currentToken == targetToken && len(path) > 1 {
		expected, err := f.calculateRouteAmount(path, amountIn)
		if err == nil {
			*routes = append(*routes, RouteResult{
				Path:              append([]string(nil), path...),
				ExpectedAmountOut: expected,
			})
		}
		return
	}

	visited[currentToken] = true
	defer func() { visited[currentToken] = false }()

	for _, edge := range f.g.adjacency[currentToken] {
		nextToken := edge.TokenB
		if !visited[nextToken] {
			newPath := append(path, nextToken)
			f.findRoutes(nextToken, targetToken, visited, newPath, routes, maxDepth, amountIn)
		}
	}
}

// calculate expected amount out for a route
func (f *Finder) calculateRouteAmount(path []string, amountIn float64) (float64, error) {
	if len(path) < 2 {
		return 0, errors.New("path must contain at least two tokens")
	}

	output := amountIn
	for i := 0; i < len(path)-1; i++ {
		from := path[i]
		to := path[i+1]

		// Find pool edge for the route
		var selectedEdge *PoolEdge
		for _, edge := range f.g.adjacency[from] {
			if edge.TokenB == to {
				selectedEdge = &edge
				break
			}
		}

		if selectedEdge == nil {
			return 0, fmt.Errorf("no pool found for %s -> %s", from, to)
		}

		// Calculate expected amount out
		var reserveIn, reserveOut float64
		if from == selectedEdge.TokenA && to == selectedEdge.TokenB {
			reserveIn = selectedEdge.ReserveA
			reserveOut = selectedEdge.ReserveB
		} else {
			reserveIn = selectedEdge.ReserveB
			reserveOut = selectedEdge.ReserveA
		}

		amtOut, err := CalculateAmountOut(output, reserveIn, reserveOut)
		if err != nil {
			return 0, err
		}

		output = amtOut
	}

	return output, nil
}
