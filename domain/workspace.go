package domain

import "time"

type Workspace struct {
	Id         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	ProfileUrl string    `json:"profileUrl" db:"profile_url"`
	OwnerId    string    `json:"ownerId" db:"owner_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`

	// For aggregation
	TotalAssignment  int                    `json:"totalAssignment"`
	UserProgression  int                    `json:"progression"`
	OwnerName        string                 `json:"ownerName,omitempty"`
	Participants     []WorkspaceParticipant `json:"participants,omitempty"`
	ParticipantCount int                    `json:"participantCount,omitempty" db:"participant_count"`
}

type WorkspaceParticipant struct {
	WorkspaceId int       `json:"-" db:"workspace_id"`
	UserId      string    `json:"userId" db:"user_id"`
	JoinedAt    time.Time `json:"joinedAt" db:"joined_at"`
}

type WorkspaceSelector struct {
	Progression  bool
	OwnerName    bool
	Participants bool
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
