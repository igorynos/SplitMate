package httpapi

import (
	"encoding/json"
	"github.com/igorynos/SplitMate/internal/model"
	"github.com/igorynos/SplitMate/internal/repository"
	"net/http"
	"strconv"
)

type Server struct{ Repo *repository.Postgres }

func (s *Server) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) { write(w, 200, map[string]string{"status": "ok"}) })
	m.HandleFunc("POST /users", s.users)
	m.HandleFunc("POST /ledgers", s.ledgers)
	m.HandleFunc("POST /ledgers/{id}/members", s.members)
	m.HandleFunc("POST /ledgers/{id}/expenses", s.expenses)
	m.HandleFunc("POST /ledgers/{id}/payments", s.payments)
	m.HandleFunc("GET /ledgers/{id}/settlement", s.settlement)
	m.HandleFunc("DELETE /documents/{kind}/{id}", s.documents)
	return m
}
func (s *Server) users(w http.ResponseWriter, r *http.Request) {
	var v model.User
	if !body(w, r, &v) {
		return
	}
	out, e := s.Repo.CreateUser(r.Context(), v)
	result(w, out, e, 201)
}
func (s *Server) ledgers(w http.ResponseWriter, r *http.Request) {
	var in struct {
		model.Ledger
		Members []int64 `json:"members"`
	}
	if !body(w, r, &in) {
		return
	}
	out, e := s.Repo.CreateLedger(r.Context(), in.Ledger, in.Members)
	result(w, out, e, 201)
}
func (s *Server) members(w http.ResponseWriter, r *http.Request) {
	var in struct {
		UserID int64 `json:"user_id"`
	}
	if !body(w, r, &in) {
		return
	}
	e := s.Repo.Join(r.Context(), id(r, "id"), in.UserID)
	result(w, map[string]bool{"joined": true}, e, 201)
}
func (s *Server) expenses(w http.ResponseWriter, r *http.Request) {
	var v model.Expense
	if !body(w, r, &v) {
		return
	}
	v.LedgerID = id(r, "id")
	out, e := s.Repo.AddExpense(r.Context(), v)
	result(w, out, e, 201)
}
func (s *Server) payments(w http.ResponseWriter, r *http.Request) {
	var v model.Payment
	if !body(w, r, &v) {
		return
	}
	v.LedgerID = id(r, "id")
	out, e := s.Repo.AddPayment(r.Context(), v)
	result(w, out, e, 201)
}
func (s *Server) settlement(w http.ResponseWriter, r *http.Request) {
	out, e := s.Repo.Settlement(r.Context(), id(r, "id"))
	result(w, out, e, 200)
}
func (s *Server) documents(w http.ResponseWriter, r *http.Request) {
	user, _ := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	e := s.Repo.DeleteDocument(r.Context(), r.PathValue("kind"), id(r, "id"), user)
	result(w, map[string]bool{"deleted": true}, e, 200)
}
func body(w http.ResponseWriter, r *http.Request, v any) bool {
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(v) != nil {
		write(w, 400, map[string]string{"error": "invalid JSON"})
		return false
	}
	return true
}
func id(r *http.Request, name string) int64 {
	v, _ := strconv.ParseInt(r.PathValue(name), 10, 64)
	return v
}
func result(w http.ResponseWriter, v any, e error, code int) {
	if e != nil {
		write(w, 500, map[string]string{"error": e.Error()})
		return
	}
	write(w, code, v)
}
func write(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
