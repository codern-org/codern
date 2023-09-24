package domain

import (
	"time"
)

type Workspace struct {
	Id         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	ProfileUrl string    `json:"profileUrl" db:"profile_url"`
	OwnerId    string    `json:"ownerId" db:"owner_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`

	// Always aggregation
	OwnerName        string `json:"ownerName" db:"owner_name"`
	OwnerProfileUrl  string `json:"ownerProfileUrl" db:"owner_profile_url"`
	ParticipantCount int    `json:"participantCount" db:"participant_count"`
	TotalAssignment  int    `json:"totalAssignment" db:"total_assignment"`

	// Optional aggregation
	Participants []WorkspaceParticipant `json:"participants,omitempty"`
}

type WorkspaceParticipant struct {
	WorkspaceId       int       `json:"-" db:"workspace_id"`
	UserId            string    `json:"userId" db:"user_id"`
	Role              string    `json:"role" db:"role"`
	JoinedAt          time.Time `json:"joinedAt" db:"joined_at"`
	RecentlyVisitedAt time.Time `json:"-" db:"recently_visited_at"`

	// Always aggregation
	Name       string `json:"name" db:"name"`
	ProfileUrl string `json:"profileUrl" db:"profile_url"`
}

type WorkspaceSelector struct {
	Participants bool
}

type WorkspaceRepository interface {
	HasUser(userId string, workspaceId int) (bool, error)
	HasAssignment(assignmentId int, workspaceId int) (bool, error)
	Get(id int, selector *WorkspaceSelector) (*Workspace, error)
	List(userId string, selector *WorkspaceSelector) ([]Workspace, error)
	ListRecent(userId string) ([]Workspace, error)
	UpdateRecent(userId string, workspaceId int) error
}

type WorkspaceUsecase interface {
	HasUser(userId string, workspaceId int) (bool, error)
	HasAssignment(assignmentId int, workspaceId int) (bool, error)
	Get(id int, selector *WorkspaceSelector, userId string) (*Workspace, error)
	List(userId string, selector *WorkspaceSelector) ([]Workspace, error)
	ListRecent(userId string) ([]Workspace, error)
}
