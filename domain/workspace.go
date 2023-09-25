package domain

import (
	"io"
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

type AssignmentLevel string

const (
	AssignmentEasyLevel   AssignmentLevel = "EASY"
	AssignmentMediumLevel AssignmentLevel = "MEDIUM"
	AssignmentHardLevel   AssignmentLevel = "HARD"
)

type AssignmentStatus string

const (
	AssignmentStatusTodo    AssignmentStatus = "TODO"
	AssignmentStatusGrading AssignmentStatus = "GRADING"
	AssignmentStatusError   AssignmentStatus = "ERROR"
	AssignmentStatusDone    AssignmentStatus = "DONE"
)

type Assignment struct {
	Id          int             `json:"id" db:"id"`
	WorkspaceId int             `json:"-" db:"workspace_id"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description" db:"description"`
	DetailUrl   string          `json:"detailUrl" db:"detail_url"`
	MemoryLimit int             `json:"memoryLimit" db:"memory_limit"`
	TimeLimit   int             `json:"timeLimit" db:"time_limit"`
	Level       AssignmentLevel `json:"level" db:"level"`
	CreatedAt   time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" db:"updated_at"`

	// Always aggregation
	Testcases       []Testcase       `json:"-"`
	LastSubmittedAt *time.Time       `json:"lastSubmittedAt" db:"last_submitted_at"`
	Status          AssignmentStatus `json:"status" db:"status"`
}

type Submission struct {
	Id             int       `json:"id" db:"id"`
	AssignmentId   int       `json:"-" db:"assignment_id"`
	UserId         string    `json:"-" db:"user_id"`
	Language       string    `json:"language" db:"language"`
	FileUrl        string    `json:"-" db:"file_url"`
	SubmittedAt    time.Time `json:"submittedAt" db:"submitted_at"`
	CompilationLog *string   `json:"compilationLog" db:"compilation_log"`

	// Always aggregation
	Results []SubmissionResult `json:"results"`
}

type SubmissionResultStatus string

const (
	SubmissionResultGrading SubmissionResultStatus = "GRADING"
	SubmissionResultError   SubmissionResultStatus = "ERROR"
	SubmissionResultDone    SubmissionResultStatus = "DONE"
)

type SubmissionResult struct {
	SubmissionId int                    `json:"-" db:"submission_id"`
	TestcaseId   int                    `json:"-" db:"testcase_id"`
	Status       SubmissionResultStatus `json:"status" db:"status"`

	// Can be null if status is `GRADING`
	StatusDetail *string `json:"statusDetail" db:"status_detail"`
	MemoryUsage  *int    `json:"memoryUsage" db:"memory_usage"`
	TimeUsage    *int    `json:"timeUsage" db:"time_usage"`
}

type Testcase struct {
	Id            int    `json:"id" db:"id"`
	AssignmentId  int    `json:"assignmentId" db:"assignment_id"`
	InputFileUrl  string `json:"inputFileUrl" db:"input_file_url"`
	OutputFileUrl string `json:"outputFileUrl" db:"output_file_url"`
}

type TestcaseFile struct {
	Input  io.Reader
	Output io.Reader
}

type WorkspaceRepository interface {
	CreateWorkspace(workspace *Workspace) error
	CreateTestcases(testcases []Testcase) error
	CreateAssigment(assignment *Assignment) error
	CreateSubmission(submission *Submission, testcases []Testcase) error
	IsUserIn(userId string, workspaceId int) (bool, error)
	IsAssignmentIn(assignmentId int, workspaceId int) (bool, error)
	Get(id int, selector *WorkspaceSelector) (*Workspace, error)
	GetAssignment(id int, userId string) (*Assignment, error)
	GetSubmission(id int) (*Submission, error)
	List(userId string, selector *WorkspaceSelector) ([]Workspace, error)
	ListRecent(userId string) ([]Workspace, error)
	ListAssignment(userId string, workspaceId int) ([]Assignment, error)
	ListSubmission(userId string, assignmentId int) ([]Submission, error)
	UpdateRecent(userId string, workspaceId int) error
	UpdateSubmissionResults(submissionId int, compilationLog string, results []SubmissionResult) error
}

type WorkspaceUsecase interface {
	CreateWorkspace(userId string, name string, file io.Reader) error
	CreateTestcase(assignmentId int, testcaseFiles []TestcaseFile) error
	CreateAssigment(workspaceId int, name string, description string, memoryLimit int, timeLimit int, level AssignmentLevel, file io.Reader) (*Assignment, error)
	CreateSubmission(userId string, assignmentId int, workspaceId int, language string, file io.Reader) error
	IsUserIn(userId string, workspaceId int) (bool, error)
	IsAssignmentIn(assignmentId int, workspaceId int) (bool, error)
	Get(id int, selector *WorkspaceSelector, userId string) (*Workspace, error)
	GetAssignment(id int, userId string) (*Assignment, error)
	GetSubmission(id int) (*Submission, error)
	List(userId string, selector *WorkspaceSelector) ([]Workspace, error)
	ListRecent(userId string) ([]Workspace, error)
	ListAssignment(userId string, workspaceId int) ([]Assignment, error)
	ListSubmission(userId string, assignmentId int) ([]Submission, error)
	UpdateSubmissionResults(submissionId int, compilationLog string, results []SubmissionResult) error
}
