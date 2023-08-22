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
	err := r.db.Get(
		&workspace,
		"SELECT *, (SELECT COUNT(*) FROM workspace_participant WHERE workspace_id = ?) AS participant_count FROM workspace WHERE id = ? LIMIT 1",
		id, id,
	)
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

	if selector.OwnerName {
		err = r.db.Get(&workspace.OwnerName, "SELECT display_name FROM user WHERE id = ?", workspace.OwnerId)
		if err == sql.ErrNoRows {
			return nil, nil
		} else if err != nil {
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

	query, args, err := sqlx.In(
		"SELECT *, (SELECT COUNT(*) FROM workspace_participant WHERE workspace_id = id) AS participant_count FROM workspace WHERE id IN (?)",
		ids,
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

	if selector.OwnerName {
		var ownerIds []string
		for i := range workspaces {
			ownerIds = append(ownerIds, workspaces[i].OwnerId)
		}

		query, args, err := sqlx.In("SELECT id, display_name FROM user WHERE id IN (?)", ownerIds)
		if err != nil {
			return nil, err
		}
		var results []struct {
			Id          string `db:"id"`
			DisplayName string `db:"display_name"`
		}
		err = r.db.Select(&results, query, args...)
		if err != nil {
			return nil, err
		}

		for i := range workspaces {
			workspace := &workspaces[i]
			for j := range results {
				if workspace.OwnerId == results[j].Id {
					workspace.OwnerName = results[j].DisplayName
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
