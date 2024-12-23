package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/anunayjoshi29/token-server/internal/routecalc"
)

type RoutesRequest struct {
	FromToken string  `json:"fromToken"`
	ToToken   string  `json:"toToken"`
	AmountIn  float64 `json:"amountIn"`
}

type RoutesResponse struct {
	Routes []struct {
		Path              []string `json:"path"`
		ExpectedAmountOut float64  `json:"expectedAmountOut"`
	} `json:"routes"`
}

func RoutesHandler(finder *routecalc.Finder, cache *routecalc.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RoutesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cacheKey := req.FromToken + "-" + req.ToToken + "-" + fmt.Sprintf("%f", req.AmountIn)
		if val, found := cache.Get(cacheKey); found {
			// Serve from cache and update cache time when a value is served in response
			cache.Set(cacheKey, val, 30*time.Second)
			writeRoutesResponse(w, val)
			return
		}

		routes := finder.FindAllRoutes(req.FromToken, req.ToToken, req.AmountIn)
		cache.Set(cacheKey, routes, 30*time.Second)
		writeRoutesResponse(w, routes)
	}
}

func writeRoutesResponse(w http.ResponseWriter, routes []routecalc.RouteResult) {
	w.Header().Set("Content-Type", "application/json")
	response := RoutesResponse{}
	for _, rt := range routes {
		response.Routes = append(response.Routes, struct {
			Path              []string `json:"path"`
			ExpectedAmountOut float64  `json:"expectedAmountOut"`
		}{
			Path: rt.Path, ExpectedAmountOut: rt.ExpectedAmountOut,
		})
	}
	json.NewEncoder(w).Encode(response)
}
