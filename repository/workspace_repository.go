package repository

import (
	"database/sql"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
)

type workspaceRepository struct {
	db *sqlx.DB
}

func NewWorkspaceRepository(db *sqlx.DB) domain.WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) IsUserIn(userId string, workspaceId int) (bool, error) {
	var result domain.WorkspaceParticipant
	err := r.db.Get(
		&result,
		"SELECT workspace_id FROM workspace_participant WHERE workspace_id = ? AND user_id = ? LIMIT 1",
		workspaceId, userId,
	)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (r *workspaceRepository) Get(id int, selector *domain.WorkspaceSelector) (*domain.Workspace, error) {
	var workspace domain.Workspace
	err := r.db.Get(&workspace, `
		SELECT
			w.*,
			user.display_name AS owner_name,
			count(wp.user_id) AS participant_count,
			count(a.id) AS total_assignment
		FROM workspace w
		INNER JOIN user ON user.id = w.owner_id
		INNER JOIN workspace_participant wp ON wp.workspace_id = w.id
		LEFT JOIN assignment a ON a.workspace_id = w.id
		WHERE w.id = ?
		GROUP BY w.id
		LIMIT 1
	`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if selector.Participants {
		err = r.db.Select(&workspace.Participants, "SELECT * FROM workspace_participant WHERE workspace_id = ?", id)
		if err != nil {
			return nil, err
		}
	}

	return &workspace, nil
}

func (r *workspaceRepository) List(ids []int, selector *domain.WorkspaceSelector) (*[]domain.Workspace, error) {
	workspaces := make([]domain.Workspace, 0)
	if len(ids) == 0 {
		return &workspaces, nil
	}

	query, args, err := sqlx.In(`
		SELECT
			w.*,
			user.display_name AS owner_name,
			count(wp.user_id) AS participant_count,
			count(a.id) AS total_assignment
		FROM workspace w
		INNER JOIN user ON user.id = w.owner_id
		INNER JOIN workspace_participant wp ON wp.workspace_id = w.id
		LEFT JOIN assignment a ON a.workspace_id = w.id
		WHERE w.id IN (?)
		GROUP BY w.id
		`, ids,
	)
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&workspaces, query, args...)
	if err != nil {
		return nil, err
	}

	if selector.Participants {
		var participants []domain.WorkspaceParticipant
		query, args, err := sqlx.In("SELECT * FROM workspace_participant WHERE workspace_id IN (?)", ids)
		if err != nil {
			return nil, err
		}
		err = r.db.Select(&participants, query, args...)
		if err != nil {
			return nil, err
		}

		for i := range workspaces {
			workspace := &workspaces[i]
			for j := range participants {
				if workspace.Id == participants[j].WorkspaceId {
					workspace.Participants = append(workspace.Participants, participants[j])
				}
			}
		}
	}

	return &workspaces, nil
}

func (r *workspaceRepository) ListFromUserId(
	userId string,
	selector *domain.WorkspaceSelector,
) (*[]domain.Workspace, error) {
	var workspaceIds []int
	err := r.db.Select(&workspaceIds, "SELECT workspace_id FROM workspace_participant WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	return r.List(workspaceIds, selector)
}
