package domain

import "time"

type Workspace struct {
	Id         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	ProfileUrl string    `json:"profileUrl" db:"profile_url"`
	OwnerId    string    `json:"ownerId" db:"owner_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`

	// For aggregation
	Participants     []WorkspaceParticipant `json:"participants,omitempty"`
	ParticipantCount int                    `json:"participantCount,omitempty"`
}

type WorkspaceParticipant struct {
	WorkspaceId int       `json:"-" db:"workspace_id"`
	UserId      string    `json:"userId" db:"user_id"`
	JoinedAt    time.Time `json:"joinedAt" db:"joined_at"`
}

type WorkspaceRepository interface {
	Get(id int, hasParticipant bool) (*Workspace, error)
	GetAllFromUserId(userId string, hasParticipant bool) (*[]Workspace, error)
}

type WorkspaceUsecase interface {
	Get(id int, hasParticipant bool) (*Workspace, error)
	GetAllFromUserId(userId string, hasParticipant bool) (*[]Workspace, error)
}
