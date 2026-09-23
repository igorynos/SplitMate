package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/igorynos/SplitMate/internal/debts"
)

func main() {
	svc := debts.NewService()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	mux.HandleFunc("/groups/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 3 || parts[0] != "groups" {
			http.NotFound(w, r)
			return
		}
		group, action := parts[1], parts[2]
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && action == "expenses":
			var e debts.Expense
			if json.NewDecoder(r.Body).Decode(&e) != nil || svc.Add(group, e) != nil {
				http.Error(w, "invalid expense", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodGet && action == "settlement":
			json.NewEncoder(w).Encode(svc.Settle(group))
		default:
			http.NotFound(w, r)
		}
	})
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("SplitMate listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
