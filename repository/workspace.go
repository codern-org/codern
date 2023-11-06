package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
)

type workspaceRepository struct {
	db *sqlx.DB
}

func NewWorkspaceRepository(db *sqlx.DB) domain.WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) HasUser(userId string, workspaceId int) (bool, error) {
	var result domain.WorkspaceParticipant
	err := r.db.Get(
		&result,
		"SELECT workspace_id FROM workspace_participant WHERE workspace_id = ? AND user_id = ?",
		workspaceId, userId,
	)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("cannot query to check user in workspace participant: %w", err)
	}
	return true, nil
}

func (r *workspaceRepository) HasAssignment(assignmentId int, workspaceId int) (bool, error) {
	var result domain.Assignment
	err := r.db.Get(
		&result,
		"SELECT id FROM assignment WHERE id = ? AND workspace_id = ?",
		assignmentId, workspaceId,
	)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("cannot query to check assignment in workspace: %w", err)
	}
	return true, nil
}

func (r *workspaceRepository) Get(
	id int,
	userId string,
	selector *domain.WorkspaceSelector,
) (*domain.Workspace, error) {
	workspaces, err := r.list([]int{id}, userId, selector)
	if err != nil {
		return nil, fmt.Errorf("cannot query to get workspace: %w", err)
	} else if len(workspaces) == 0 {
		return nil, nil
	}
	return &workspaces[0], nil
}

func (r *workspaceRepository) GetRole(userId string, workspaceId int) (*domain.WorkspaceRole, error) {
	var role domain.WorkspaceRole
	err := r.db.Get(
		&role,
		"SELECT role FROM workspace_participant WHERE user_id = ? AND workspace_id = ?",
		userId, workspaceId,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get workspace role: %w", err)
	}
	return &role, nil
}

func (r *workspaceRepository) List(
	userId string,
	selector *domain.WorkspaceSelector,
) ([]domain.Workspace, error) {
	var workspaceIds []int
	err := r.db.Select(
		&workspaceIds,
		"SELECT workspace_id FROM workspace_participant WHERE user_id = ?",
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot query to list workspace id: %w", err)
	}
	return r.list(workspaceIds, userId, selector)
}

func (r *workspaceRepository) list(
	ids []int,
	userId string,
	selector *domain.WorkspaceSelector,
) ([]domain.Workspace, error) {
	workspaces := make([]domain.Workspace, 0)
	if len(ids) == 0 {
		return workspaces, nil
	}

	query, args, err := sqlx.In(`
		SELECT
			w.*,
			user.display_name AS owner_name,
			user.profile_url AS owner_profile_url,
			(SELECT COUNT(*) FROM workspace_participant wp WHERE wp.workspace_id = w.id) AS participant_count,
			(SELECT COUNT(*) FROM assignment a WHERE a.workspace_id = w.id) AS total_assignment,
			wp.role, wp.favorite, wp.joined_at, wp.recently_visited_at
		FROM workspace w
		INNER JOIN user ON user.id = (SELECT user_id FROM workspace_participant WHERE workspace_id = w.id AND role = 'OWNER')
		INNER JOIN workspace_participant wp ON wp.workspace_id = w.id AND wp.user_id = ?
		WHERE w.id IN (?)
	`, userId, ids)
	if err != nil {
		return nil, fmt.Errorf("cannot query to create query to list workspace: %w", err)
	}
	if err := r.db.Select(&workspaces, query, args...); err != nil {
		return nil, fmt.Errorf("cannot query to list workspace: %w", err)
	}

	if selector.Participants {
		participants := make([]domain.WorkspaceParticipant, 0)
		query, args, err := sqlx.In(`
			SELECT
				user.display_name as name,
				user.profile_url,
				wp.*
			FROM workspace_participant wp
			INNER JOIN user ON user.id = wp.user_id
			WHERE workspace_id IN (?)
		`, ids)
		if err != nil {
			return nil, fmt.Errorf("cannot query to create query to list workspace participant: %w", err)
		}
		if err := r.db.Select(&participants, query, args...); err != nil {
			return nil, fmt.Errorf("cannot query to list workspace participant: %w", err)
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

	return workspaces, nil
}

func (r *workspaceRepository) UpdateRecent(userId string, workspaceId int) error {
	_, err := r.db.Exec(`
		UPDATE workspace_participant SET recently_visited_at = ? WHERE user_id = ? AND workspace_id = ?
	`, time.Now(), userId, workspaceId)
	if err != nil {
		return fmt.Errorf("cannot query to update recent workspace: %w", err)
	}
	return nil
}

func (r *workspaceRepository) UpdateRole(
	userId string,
	workspaceId int,
	role domain.WorkspaceRole,
) error {
	_, err := r.db.Exec(
		"UPDATE workspace_participant SET role = ? WHERE user_id = ? AND workspace_id = ?",
		role, userId, workspaceId,
	)
	if err != nil {
		return fmt.Errorf("cannot query to update role: %w", err)
	}
	return nil
}
