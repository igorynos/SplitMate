package debts

import "testing"

func TestSettlement(t *testing.T) {
	s := NewService()
	if err := s.Add("trip", Expense{Description: "Dinner", Amount: 3000, PaidBy: "Igor", Participants: []string{"Igor", "Anna", "Max"}}); err != nil {
		t.Fatal(err)
	}
	got := s.Settle("trip")
	if len(got) != 2 || got[0].Amount != 1000 || got[1].Amount != 1000 {
		t.Fatalf("unexpected settlement: %#v", got)
	}
}
