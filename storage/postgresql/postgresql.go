package postgresql

import (
	"database/sql"

	"github.com/jonfriesen/subscriber-tracker-api/model"

	_ "github.com/lib/pq"
)

// PostgreSQL is our DB object that satifies our interface
type PostgreSQL struct {
	db *sql.DB
}

// NewConnection creates a DB connection
func NewConnection(databaseURL string) (*PostgreSQL, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreSQL{
		db: db,
	}, nil
}

// ListSubscribers subs
func (p *PostgreSQL) ListSubscribers() (*model.Subscriber, error) {
	return nil, nil
}

// AddSubscriber adds
func (p *PostgreSQL) AddSubscriber(newSub *model.Subscriber) (*model.Subscriber, error) {
	return nil, nil
}

func (p *PostgreSQL) Close() error {
	return p.db.Close()
}
