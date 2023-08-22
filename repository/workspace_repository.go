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

func (r *workspaceRepository) Get(id int, hasParticipant bool) (*domain.Workspace, error) {
	var workspace domain.Workspace
	err := r.db.Get(&workspace, "SELECT * FROM workspace WHERE id = ? LIMIT 1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if hasParticipant {
		var participants []domain.WorkspaceParticipant
		err = r.db.Select(&participants, "SELECT * FROM workspace_participant WHERE workspace_id = ?", id)
		if err != nil {
			return nil, err
		}
		workspace.Participants = participants
		return &workspace, nil
	}

	var participantCount int
	err = r.db.Get(&participantCount, "SELECT COUNT(*) FROM workspace_participant WHERE workspace_id = ?", id)
	if err != nil {
		return nil, err
	}
	workspace.ParticipantCount = participantCount
	return &workspace, nil
}

func (r *workspaceRepository) List(ids []int, hasParticipant bool) (*[]domain.Workspace, error) {
	var workspaces []domain.Workspace
	query, args, err := sqlx.In("SELECT * FROM workspace WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&workspaces, query, args...)
	if err != nil {
		return nil, err
	}

	if hasParticipant {
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
			for _, participant := range participants {
				if workspace.Id == participant.WorkspaceId {
					workspace.Participants = append(workspace.Participants, participant)
				}
			}
		}
		return &workspaces, nil
	}

	var participantCounts []struct {
		WorkspaceId int `db:"workspace_id"`
		Count       int `db:"count"`
	}
	query, args, err = sqlx.In(
		"SELECT workspace_id, COUNT(*) AS count FROM workspace_participant WHERE workspace_id IN (?) GROUP BY workspace_id",
		ids,
	)
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&participantCounts, query, args...)
	if err != nil {
		return nil, err
	}
	for i := range workspaces {
		workspace := &workspaces[i]
		for _, participantCount := range participantCounts {
			if workspace.Id == participantCount.WorkspaceId {
				workspace.ParticipantCount = participantCount.Count
			}
		}
	}
	return &workspaces, nil
}

func (r *workspaceRepository) ListFromUserId(userId string, hasParticipant bool) (*[]domain.Workspace, error) {
	var workspaceId []int
	err := r.db.Select(&workspaceId, "SELECT workspace_id FROM workspace_participant WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	return r.List(workspaceId, hasParticipant)
}
