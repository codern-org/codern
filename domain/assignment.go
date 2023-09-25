package domain

import (
	"io"
	"time"
)

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

type AssignmentRepository interface {
	CreateAssigment(assignment *Assignment) error
	CreateTestcases(testcases []Testcase) error
	CreateSubmission(submission *Submission, testcases []Testcase) error
	Get(id int, userId string) (*Assignment, error)
	GetSubmission(id int) (*Submission, error)
	List(userId string, workspaceId int) ([]Assignment, error)
	ListSubmission(userId string, assignmentId int) ([]Submission, error)
	UpdateSubmissionResults(submissionId int, compilationLog string, results []SubmissionResult) error
}

type AssignmentUsecase interface {
	CreateAssigment(workspaceId int, name string, description string, memoryLimit int, timeLimit int, level AssignmentLevel, file io.Reader) (*Assignment, error)
	CreateTestcase(assignmentId int, testcaseFiles []TestcaseFile) error
	CreateSubmission(userId string, assignmentId int, workspaceId int, language string, file io.Reader) error
	Get(id int, userId string) (*Assignment, error)
	GetSubmission(id int) (*Submission, error)
	List(userId string, workspaceId int) ([]Assignment, error)
	ListSubmission(userId string, assignmentId int) ([]Submission, error)
	UpdateSubmissionResults(submissionId int, compilationLog string, results []SubmissionResult) error
}
