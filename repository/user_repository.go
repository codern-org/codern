package repository

import (
	"database/sql"

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
		"INSERT INTO user (id, email, password, display_name, profile_url, provider, created_at)"+
			"VALUES (:id, :email, :password, :display_name, :profile_url, :provider, :created_at)",
		user,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Get(id string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.Get(&user, "SELECT * FROM user WHERE id = ? LIMIT 1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetBySessionId(id string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.Get(
		&user,
		"SELECT user.* FROM user JOIN session ON user.id = session.user_id WHERE session.id = ? LIMIT 1",
		id,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetSelfProviderUser(email string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.Get(
		&user,
		"SELECT * FROM user WHERE email = ? AND provider = \"SELF\" LIMIT 1",
		email,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
