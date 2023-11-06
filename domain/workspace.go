package domain

import (
	"time"
)

type Workspace struct {
	Id         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	ProfileUrl string    `json:"profileUrl" db:"profile_url"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`

	// Always aggregation
	OwnerName           string    `json:"ownerName" db:"owner_name"`
	OwnerProfileUrl     string    `json:"ownerProfileUrl" db:"owner_profile_url"`
	ParticipantCount    int       `json:"participantCount" db:"participant_count"`
	TotalAssignment     int       `json:"totalAssignment" db:"total_assignment"`
	CompletedAssignment int       `json:"completedAssignment" db:"completed_assignment"`
	Role                string    `json:"role" db:"role"`
	Favorite            bool      `json:"favorite" db:"favorite"`
	JoinedAt            time.Time `json:"joinedAt" db:"joined_at"`
	RecentlyVisitedAt   time.Time `json:"recentlyVisitedAt" db:"recently_visited_at"`

	// Optional aggregation
	Participants []WorkspaceParticipant `json:"participants,omitempty"`
}

type WorkspaceRole string

const (
	MemberRole WorkspaceRole = "MEMBER"
	AdminRole  WorkspaceRole = "ADMIN"
	OwnerRole  WorkspaceRole = "OWNER"
)

type WorkspaceParticipant struct {
	WorkspaceId       int           `json:"-" db:"workspace_id"`
	UserId            string        `json:"userId" db:"user_id"`
	Role              WorkspaceRole `json:"role" db:"role"`
	Favorite          bool          `json:"-" db:"favorite"`
	JoinedAt          time.Time     `json:"joinedAt" db:"joined_at"`
	RecentlyVisitedAt time.Time     `json:"-" db:"recently_visited_at"`

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
	Get(id int, userId string, selector *WorkspaceSelector) (*Workspace, error)
	GetRole(userId string, workspaceId int) (*WorkspaceRole, error)
	List(userId string, selector *WorkspaceSelector) ([]Workspace, error)
	UpdateRecent(userId string, workspaceId int) error
	UpdateRole(userId string, workspaceId int, role WorkspaceRole) error
}

type WorkspaceUsecase interface {
	HasUser(userId string, workspaceId int) (bool, error)
	HasAssignment(assignmentId int, workspaceId int) (bool, error)
	Get(id int, selector *WorkspaceSelector, userId string) (*Workspace, error)
	List(userId string, selector *WorkspaceSelector) ([]Workspace, error)
	UpdateRole(updaterUserId string, targetUserId string, workspaceId int, role WorkspaceRole) error
}
