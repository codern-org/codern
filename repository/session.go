package repository

import (
	"database/sql"
	"fmt"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
)

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) domain.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *domain.Session) error {
	_, err := r.db.NamedExec(
		"INSERT INTO session VALUES (:id, :user_id, :ip_address, :user_agent, :created_at, :expired_at)",
		session,
	)
	if err != nil {
		return fmt.Errorf("cannot query to create session: %w", err)
	}
	return nil
}

func (r *sessionRepository) Get(id string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.Get(&session, "SELECT * FROM session WHERE id = ? LIMIT 1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get session: %w", err)
	}
	return &session, nil
}

func (r *sessionRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM session WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("cannot query to delete session: %w", err)
	}
	return nil
}

func (r *sessionRepository) DeleteDuplicates(userId string, ipAddress string, userAgent string) error {
	_, err := r.db.Exec(
		"DELETE FROM session WHERE user_id = ? AND user_agent = ? AND ip_address = ?",
		userId, userAgent, ipAddress,
	)
	if err != nil {
		return fmt.Errorf("cannot query to delete duplicated session: %w", err)
	}
	return nil
}
