package domain

import (
	"time"
)

type Workspace struct {
	Id                  int       `json:"id" db:"id"`
	Name                string    `json:"name" db:"name"`
	ProfileUrl          string    `json:"profileUrl" db:"profile_url"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	OwnerName           string    `json:"ownerName" db:"owner_name"`
	OwnerProfileUrl     string    `json:"ownerProfileUrl" db:"owner_profile_url"`
	ParticipantCount    int       `json:"participantCount" db:"participant_count"`
	TotalAssignment     int       `json:"totalAssignment" db:"total_assignment"`
	CompletedAssignment int       `json:"completedAssignment" db:"completed_assignment"`
	Role                string    `json:"role" db:"role"`
	Favorite            bool      `json:"favorite" db:"favorite"`
	JoinedAt            time.Time `json:"joinedAt" db:"joined_at"`
	RecentlyVisitedAt   time.Time `json:"recentlyVisitedAt" db:"recently_visited_at"`
}

type UpdateWorkspace struct {
	Name     *string
	Favorite *bool
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
	Name              string        `json:"name" db:"name"`
	Role              WorkspaceRole `json:"role" db:"role"`
	ProfileUrl        string        `json:"profileUrl" db:"profile_url"`
	Favorite          bool          `json:"-" db:"favorite"`
	JoinedAt          time.Time     `json:"joinedAt" db:"joined_at"`
	RecentlyVisitedAt time.Time     `json:"-" db:"recently_visited_at"`
}

type WorkspaceInvitation struct {
	Id          string    `json:"id" db:"id"`
	WorkspaceId int       `json:"workspaceId" db:"workspace_id"`
	InviterId   string    `json:"inviterId" db:"inviter_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	ValidAt     time.Time `json:"validAt" db:"valid_at"`
	ValidUntil  time.Time `json:"validUntil" db:"valid_until"`
}

type WorkspaceRank struct {
	UserId              string `json:"userId" db:"id"`
	DisplayName         string `json:"displayName" db:"display_name"`
	ProfileUrl          string `json:"profileUrl" db:"profile_url"`
	Score               int    `json:"score" db:"score"`
	CompletedAssignment int    `json:"completedAssignment" db:"completed_assignment"`
}

type WorkspaceRepository interface {
	CreateInvitation(invitation *WorkspaceInvitation) error
	GetInvitation(id string) (*WorkspaceInvitation, error)
	GetInvitations(workspaceId int) ([]WorkspaceInvitation, error)
	DeleteInvitation(invitationId string) error
	CreateParticipant(participant *WorkspaceParticipant) error
	HasUser(userId string, workspaceId int) (bool, error)
	HasAssignment(assignmentId int, workspaceId int) (bool, error)
	Get(id int, userId string) (*Workspace, error)
	GetRole(userId string, workspaceId int) (*WorkspaceRole, error)
	GetScoreboard(workspaceId int) ([]WorkspaceRank, error)
	List(userId string) ([]Workspace, error)
	ListParticipant(workspaceId int) ([]WorkspaceParticipant, error)
	Update(workspace *Workspace) error
	UpdateRecent(userId string, workspaceId int) error
	UpdateRole(userId string, workspaceId int, role WorkspaceRole) error
}

type WorkspaceUsecase interface {
	CreateInvitation(workspaceId int, inviterId string, validAt time.Time, validUntil time.Time) (string, error)
	GetInvitation(id string) (*WorkspaceInvitation, error)
	GetInvitations(workspaceId int) ([]WorkspaceInvitation, error)
	DeleteInvitation(invitationId string, userId string) error
	CreateParticipant(workspaceId int, userId string, role WorkspaceRole) error
	JoinByInvitation(userId string, invitationCode string) error
	HasUser(userId string, workspaceId int) (bool, error)
	HasAssignment(assignmentId int, workspaceId int) (bool, error)
	Get(id int, userId string) (*Workspace, error)
	GetRole(userId string, workspaceId int) (*WorkspaceRole, error)
	GetScoreboard(workspaceId int) ([]WorkspaceRank, error)
	CheckPerm(userId string, workspaceId int) (bool, error)
	List(userId string) ([]Workspace, error)
	ListParticipant(workspaceId int) ([]WorkspaceParticipant, error)
	Update(userId string, workspaceId int, workspace *UpdateWorkspace) error
	UpdateRole(updaterUserId string, targetUserId string, workspaceId int, role WorkspaceRole) error
}
