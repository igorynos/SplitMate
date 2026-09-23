package repository

import (
	"context"
	"github.com/igorynos/SplitMate/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct{ db *pgxpool.Pool }

func NewPostgres(db *pgxpool.Pool) *Postgres          { return &Postgres{db: db} }
func (p *Postgres) Migrate(ctx context.Context) error { _, err := p.db.Exec(ctx, schema); return err }
func (p *Postgres) CreateUser(ctx context.Context, v model.User) (model.User, error) {
	err := p.db.QueryRow(ctx, `INSERT INTO users(id,name) VALUES($1,$2) ON CONFLICT(id) DO UPDATE SET name=EXCLUDED.name RETURNING id`, v.ID, v.Name).Scan(&v.ID)
	return v, err
}
func (p *Postgres) CreateLedger(ctx context.Context, v model.Ledger, members []int64) (model.Ledger, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return v, err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `INSERT INTO ledgers(name,owner_id) VALUES($1,$2) RETURNING id,created_at`, v.Name, v.OwnerID).Scan(&v.ID, &v.CreatedAt)
	if err != nil {
		return v, err
	}
	for _, id := range append([]int64{v.OwnerID}, members...) {
		var wallet int64
		err = tx.QueryRow(ctx, `INSERT INTO wallets(ledger_id,name) VALUES($1,(SELECT name FROM users WHERE id=$2)) RETURNING id`, v.ID, id).Scan(&wallet)
		if err != nil {
			return v, err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO members(ledger_id,user_id,wallet_id) VALUES($1,$2,$3) ON CONFLICT DO NOTHING`, v.ID, id, wallet); err != nil {
			return v, err
		}
	}
	return v, tx.Commit(ctx)
}
func (p *Postgres) Join(ctx context.Context, ledger, user int64) error {
	var wallet int64
	err := p.db.QueryRow(ctx, `INSERT INTO wallets(ledger_id,name) VALUES($1,(SELECT name FROM users WHERE id=$2)) RETURNING id`, ledger, user).Scan(&wallet)
	if err != nil {
		return err
	}
	_, err = p.db.Exec(ctx, `INSERT INTO members(ledger_id,user_id,wallet_id) VALUES($1,$2,$3)`, ledger, user, wallet)
	return err
}
func (p *Postgres) AddExpense(ctx context.Context, v model.Expense) (model.Expense, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return v, err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `INSERT INTO expenses(ledger_id,paid_by,description,amount) VALUES($1,$2,$3,$4) RETURNING id,created_at`, v.LedgerID, v.PaidBy, v.Description, v.Amount).Scan(&v.ID, &v.CreatedAt)
	if err != nil {
		return v, err
	}
	for _, id := range v.Participants {
		if _, err = tx.Exec(ctx, `INSERT INTO expense_participants(expense_id,user_id) VALUES($1,$2)`, v.ID, id); err != nil {
			return v, err
		}
	}
	return v, tx.Commit(ctx)
}
func (p *Postgres) AddPayment(ctx context.Context, v model.Payment) (model.Payment, error) {
	err := p.db.QueryRow(ctx, `INSERT INTO payments(ledger_id,from_user,to_user,amount,comment) VALUES($1,$2,$3,$4,$5) RETURNING id,created_at`, v.LedgerID, v.FromUser, v.ToUser, v.Amount, v.Comment).Scan(&v.ID, &v.CreatedAt)
	return v, err
}
func (p *Postgres) DeleteDocument(ctx context.Context, kind string, id, user int64) error {
	table, owner := "expenses", "paid_by"
	if kind == "payment" {
		table, owner = "payments", "from_user"
	}
	_, err := p.db.Exec(ctx, `DELETE FROM `+table+` WHERE id=$1 AND `+owner+`=$2`, id, user)
	return err
}
func (p *Postgres) Settlement(ctx context.Context, ledger int64) ([]model.Transfer, error) {
	rows, err := p.db.Query(ctx, settlementSQL, ledger)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	balances := map[int64]int64{}
	for rows.Next() {
		var id, amount int64
		if err = rows.Scan(&id, &amount); err != nil {
			return nil, err
		}
		balances[id] = amount
	}
	type b struct{ id, amount int64 }
	var debtors, creditors []b
	for id, v := range balances {
		if v < 0 {
			debtors = append(debtors, b{id, -v})
		} else if v > 0 {
			creditors = append(creditors, b{id, v})
		}
	}
	var out []model.Transfer
	for i, j := 0, 0; i < len(debtors) && j < len(creditors); {
		n := debtors[i].amount
		if creditors[j].amount < n {
			n = creditors[j].amount
		}
		out = append(out, model.Transfer{FromUser: debtors[i].id, ToUser: creditors[j].id, Amount: n})
		debtors[i].amount -= n
		creditors[j].amount -= n
		if debtors[i].amount == 0 {
			i++
		}
		if creditors[j].amount == 0 {
			j++
		}
	}
	return out, nil
}

const settlementSQL = `WITH expense_shares AS (SELECT ep.user_id,e.amount/COUNT(*) OVER(PARTITION BY e.id) owed,e.paid_by,e.amount,e.id FROM expenses e JOIN expense_participants ep ON ep.expense_id=e.id WHERE e.ledger_id=$1), movements AS (SELECT user_id,-owed amount FROM expense_shares UNION ALL SELECT paid_by,MAX(amount) FROM expense_shares GROUP BY id,paid_by UNION ALL SELECT from_user,-amount FROM payments WHERE ledger_id=$1 UNION ALL SELECT to_user,amount FROM payments WHERE ledger_id=$1) SELECT user_id,SUM(amount)::bigint FROM movements GROUP BY user_id ORDER BY user_id`
const schema = `CREATE TABLE IF NOT EXISTS users(id BIGINT PRIMARY KEY,name VARCHAR(64) NOT NULL);CREATE TABLE IF NOT EXISTS ledgers(id BIGSERIAL PRIMARY KEY,name VARCHAR(128) NOT NULL,owner_id BIGINT NOT NULL REFERENCES users(id),created_at TIMESTAMPTZ NOT NULL DEFAULT now());CREATE TABLE IF NOT EXISTS wallets(id BIGSERIAL PRIMARY KEY,ledger_id BIGINT NOT NULL REFERENCES ledgers(id) ON DELETE CASCADE,name VARCHAR(128) NOT NULL);CREATE TABLE IF NOT EXISTS members(ledger_id BIGINT REFERENCES ledgers(id) ON DELETE CASCADE,user_id BIGINT REFERENCES users(id),wallet_id BIGINT REFERENCES wallets(id),joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),PRIMARY KEY(ledger_id,user_id));CREATE TABLE IF NOT EXISTS expenses(id BIGSERIAL PRIMARY KEY,ledger_id BIGINT REFERENCES ledgers(id) ON DELETE CASCADE,paid_by BIGINT REFERENCES users(id),description VARCHAR(255) NOT NULL,amount BIGINT NOT NULL CHECK(amount>0),created_at TIMESTAMPTZ NOT NULL DEFAULT now());CREATE TABLE IF NOT EXISTS expense_participants(expense_id BIGINT REFERENCES expenses(id) ON DELETE CASCADE,user_id BIGINT REFERENCES users(id),PRIMARY KEY(expense_id,user_id));CREATE TABLE IF NOT EXISTS payments(id BIGSERIAL PRIMARY KEY,ledger_id BIGINT REFERENCES ledgers(id) ON DELETE CASCADE,from_user BIGINT REFERENCES users(id),to_user BIGINT REFERENCES users(id),amount BIGINT NOT NULL CHECK(amount>0),comment VARCHAR(255),created_at TIMESTAMPTZ NOT NULL DEFAULT now());`
