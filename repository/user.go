package repository

import (
	"database/sql"
	"fmt"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	_, err := r.db.NamedExec(
		"INSERT INTO user (id, email, password, display_name, profile_url, account_type, provider, created_at)"+
			"VALUES (:id, :email, :password, :display_name, :profile_url, :account_type, :provider, :created_at)",
		user,
	)
	if err != nil {
		return fmt.Errorf("cannot query to create user: %w", err)
	}
	return nil
}

func (r *userRepository) Get(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Get(&user, "SELECT * FROM user WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetBySessionId(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Get(
		&user,
		"SELECT user.* FROM user JOIN session ON user.id = session.user_id WHERE session.id = ?",
		id,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get user by sesion id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(
	email string,
	provider domain.AuthProvider,
) (*domain.User, error) {
	var user domain.User
	err := r.db.Get(
		&user,
		"SELECT * FROM user WHERE email = ? AND provider = ?",
		email, provider,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	_, err := r.db.NamedExec(`
		UPDATE user
		SET
			email = :email,
			password = :password,
			display_name = :display_name,
			profile_url = :profile_url,
			account_type = :account_type,
			provider = :provider
		WHERE id = :id
	`, user)
	if err != nil {
		return fmt.Errorf("cannot query to update user: %w", err)
	}
	return nil
}
