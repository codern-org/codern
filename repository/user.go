package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
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
	err := r.db.Get(&user, "SELECT * FROM user WHERE id = ? LIMIT 1", id)
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
		"SELECT user.* FROM user JOIN session ON user.id = session.user_id WHERE session.id = ? LIMIT 1",
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
		"SELECT * FROM user WHERE email = ? AND provider = ? LIMIT 1",
		email, provider,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get user by email: %w", err)
	}
	return &user, nil
}

// TODO: what if the user field we want to update is number 🤔
func (r *userRepository) Update(user *domain.User) (bool, error) {
	setQueries := make([]string, 0)
	args := make([]interface{}, 0)

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return false, fmt.Errorf("cannot hash password while update user %s", err)
		}

		setQueries = append(setQueries, "password = ?")
		args = append(args, hashedPassword)
	}

	if user.DisplayName != "" {
		setQueries = append(setQueries, "display_name = ?")
		args = append(args, user.DisplayName)
	}

	if user.Email != "" {
		setQueries = append(setQueries, "email = ?")
		args = append(args, user.Email)
	}

	if len(setQueries) == 0 {
		return false, nil
	}

	updateUserQuery := "UPDATE user SET " + strings.Join(setQueries, ", ") + " WHERE id = ?"

	_, err := r.db.Exec(updateUserQuery, append(args, user.Id)...)
	if err != nil {
		return false, fmt.Errorf("cannot query to update user: %w", err)
	}

	return true, nil
}
