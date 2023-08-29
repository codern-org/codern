package domain

import "time"

type Workspace struct {
	Id         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	ProfileUrl string    `json:"profileUrl" db:"profile_url"`
	OwnerId    string    `json:"ownerId" db:"owner_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`

	// Always aggregation
	OwnerName        string `json:"ownerName" db:"owner_name"`
	ParticipantCount int    `json:"participantCount" db:"participant_count"`
	TotalAssignment  int    `json:"totalAssignment" db:"total_assignment"`

	// Optional aggregation
	Participants []WorkspaceParticipant `json:"participants,omitempty"`
}

type WorkspaceParticipant struct {
	WorkspaceId int       `json:"-" db:"workspace_id"`
	UserId      string    `json:"userId" db:"user_id"`
	JoinedAt    time.Time `json:"joinedAt" db:"joined_at"`
}

type WorkspaceSelector struct {
	Participants bool
}

type AssignmentLevel string

const (
	AssignmentEasyLevel   AssignmentLevel = "EASY"
	AssignmentNormalLevel AssignmentLevel = "NORMAL"
	AssignmentHardLevel   AssignmentLevel = "HARD"
)

type Assignment struct {
	Id          int             `json:"id" db:"id"`
	WorkspaceId int             `json:"-" db:"workspace_id"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description" db:"description"`
	DetailUrl   string          `json:"detailUrl" db:"detail_url"`
	MemoryLimit string          `json:"memoryLimit" db:"memory_limit"`
	TimeLimit   string          `json:"timeLimit" db:"time_limit"`
	Level       AssignmentLevel `json:"level" db:"level"`
	CreatedAt   time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" db:"updated_at"`
}

type WorkspaceRepository interface {
	IsUserIn(userId string, workspaceId int) (bool, error)
	Get(id int, selector *WorkspaceSelector) (*Workspace, error)
	List(ids []int, selector *WorkspaceSelector) (*[]Workspace, error)
	ListFromUserId(userId string, selector *WorkspaceSelector) (*[]Workspace, error)
}

type WorkspaceUsecase interface {
	CanUserView(userId string, workspaceIds []int) (bool, error)
	IsUserIn(userId string, workspaceId int) (bool, error)
	Get(id int, selector *WorkspaceSelector) (*Workspace, error)
	ListFromUserId(userId string, selector *WorkspaceSelector) (*[]Workspace, error)
}
