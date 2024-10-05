package repository

import (
	"database/sql"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform"
)

type miscRepository struct {
	db *platform.MySql
}

func NewMiscRepsitory(db *platform.MySql) domain.MiscRepository {
	return &miscRepository{db: db}
}

func (r *miscRepository) GetFeatureFlag(feature string) (bool, error) {
	var enabled bool
	err := r.db.Get(&enabled, "SELECT enabled FROM feature_flag WHERE feature = ?", feature)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return enabled, nil
}
