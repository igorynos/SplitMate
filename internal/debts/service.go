package debts

import (
	"errors"
	"sort"
	"sync"
)

var ErrInvalidExpense = errors.New("payer, participants and positive amount are required")

type Expense struct {
	Description  string   `json:"description"`
	Amount       int64    `json:"amount"` // smallest currency unit
	PaidBy       string   `json:"paid_by"`
	Participants []string `json:"participants"`
}

type Transfer struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

type Service struct {
	mu       sync.RWMutex
	expenses map[string][]Expense
}

func NewService() *Service { return &Service{expenses: make(map[string][]Expense)} }

func (s *Service) Add(group string, e Expense) error {
	if group == "" || e.PaidBy == "" || e.Amount <= 0 || len(e.Participants) == 0 {
		return ErrInvalidExpense
	}
	seen := map[string]bool{}
	for _, p := range e.Participants {
		if p == "" || seen[p] {
			return ErrInvalidExpense
		}
		seen[p] = true
	}
	s.mu.Lock()
	s.expenses[group] = append(s.expenses[group], e)
	s.mu.Unlock()
	return nil
}

func (s *Service) Settle(group string) []Transfer {
	s.mu.RLock()
	items := append([]Expense(nil), s.expenses[group]...)
	s.mu.RUnlock()
	balances := map[string]int64{}
	for _, e := range items {
		share, remainder := e.Amount/int64(len(e.Participants)), e.Amount%int64(len(e.Participants))
		balances[e.PaidBy] += e.Amount
		for i, p := range e.Participants {
			owed := share
			if int64(i) < remainder {
				owed++
			}
			balances[p] -= owed
		}
	}
	type balance struct {
		name   string
		amount int64
	}
	var debtors, creditors []balance
	for name, amount := range balances {
		if amount < 0 {
			debtors = append(debtors, balance{name, -amount})
		}
		if amount > 0 {
			creditors = append(creditors, balance{name, amount})
		}
	}
	sort.Slice(debtors, func(i, j int) bool { return debtors[i].name < debtors[j].name })
	sort.Slice(creditors, func(i, j int) bool { return creditors[i].name < creditors[j].name })
	var result []Transfer
	for i, j := 0, 0; i < len(debtors) && j < len(creditors); {
		amount := min(debtors[i].amount, creditors[j].amount)
		result = append(result, Transfer{From: debtors[i].name, To: creditors[j].name, Amount: amount})
		debtors[i].amount -= amount
		creditors[j].amount -= amount
		if debtors[i].amount == 0 {
			i++
		}
		if creditors[j].amount == 0 {
			j++
		}
	}
	return result
}
