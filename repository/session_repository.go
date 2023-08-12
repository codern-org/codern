package repository

import (
	"database/sql"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
)

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) domain.SessionRepository {
	return &sessionRepository{db: db}
}

func (repository *sessionRepository) Create(session *domain.Session) error {
	_, err := repository.db.NamedExec(
		"INSERT INTO session VALUES (:id, :user_id, :ip_address, :user_agent, :created_at, :expired_at)",
		session,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repository *sessionRepository) Get(id string) (*domain.Session, error) {
	session := domain.Session{}
	err := repository.db.Get(&session, "SELECT * FROM session WHERE id = ? LIMIT 1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &session, nil
}

func (repository *sessionRepository) Delete(id string) error {
	_, err := repository.db.Exec("DELETE FROM session WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *sessionRepository) DeleteDuplicates(userId string, ipAddress string, userAgent string) error {
	_, err := repository.db.Exec(
		"DELETE FROM session WHERE user_id = ? AND user_agent = ? AND ip_address = ?",
		userId, userAgent, ipAddress,
	)
	if err != nil {
		return err
	}
	return nil
}
