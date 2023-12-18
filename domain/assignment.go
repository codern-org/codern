package domain

import (
	"io"
	"mime/multipart"
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
	AssignmentStatusTodo        AssignmentStatus = "TODO"
	AssignmentStatusGrading     AssignmentStatus = "GRADING"
	AssignmentStatusIncompleted AssignmentStatus = "INCOMPLETED"
	AssignmentStatusComplete    AssignmentStatus = "COMPLETED"
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
	DueDate     *time.Time      `json:"dueDate" db:"due_date"`

	// Always aggregation
	Testcases []Testcase `json:"-"`
}

type CreateAssignment struct {
	Name          string          `validate:"required"`
	Description   string          `validate:"required"`
	MemoryLimit   int             `validate:"required"`
	TimeLimit     int             `validate:"required"`
	Level         AssignmentLevel `validate:"required"`
	DetailFile    io.Reader       `validate:"required"`
	TestcaseFiles []TestcaseFile  `validate:"required"`
}

type UpdateAssignment struct {
	Name          *string
	Description   *string
	MemoryLimit   *int
	TimeLimit     *int
	Level         *AssignmentLevel
	DetailFile    io.Reader
	TestcaseFiles *[]TestcaseFile
}

type AssignmentWithStatus struct {
	Assignment

	Status          AssignmentStatus `json:"status" db:"status"`
	LastSubmittedAt *time.Time       `json:"lastSubmittedAt" db:"last_submitted_at"`
}

type Submission struct {
	Id             int              `json:"id" db:"id"`
	AssignmentId   int              `json:"-" db:"assignment_id"`
	UserId         string           `json:"-" db:"user_id"`
	Language       string           `json:"language" db:"language"`
	Status         AssignmentStatus `json:"status" db:"status"`
	Score          int              `json:"score" db:"score"`
	FileUrl        string           `json:"fileUrl" db:"file_url"`
	SubmittedAt    time.Time        `json:"submittedAt" db:"submitted_at"`
	CompilationLog *string          `json:"compilationLog,omitempty" db:"compilation_log"`
	IsLate         bool             `json:"isLate" db:"is_late"`

	// Always aggregation
	Results []SubmissionResult `json:"results,omitempty"`
}

type SubmissionResult struct {
	SubmissionId int    `json:"-" db:"submission_id"`
	TestcaseId   int    `json:"-" db:"testcase_id"`
	IsPassed     bool   `json:"isPassed" db:"is_passed"`
	Status       string `json:"status" db:"status"`
	MemoryUsage  *int   `json:"memoryUsage" db:"memory_usage"`
	TimeUsage    *int   `json:"timeUsage" db:"time_usage"`
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

func CreateTestcaseFiles(inputs []multipart.File, outputs []multipart.File) []TestcaseFile {
	files := make([]TestcaseFile, len(inputs))
	for i, input := range inputs {
		files[i] = TestcaseFile{
			Input:  input,
			Output: outputs[i],
		}
	}
	return files
}

type AssignmentRepository interface {
	Create(assignment *Assignment) error
	Update(assignment *Assignment) error
	CreateTestcases(testcases []Testcase) error
	DeleteTestcases(assignmentId int) error
	CreateSubmission(submission *Submission, testcases []Testcase) error
	CreateSubmissionResults(submissionId int, compilationLog string, status AssignmentStatus, score int, results []SubmissionResult) error
	Get(id int) (*Assignment, error)
	GetWithStatus(id int, userId string) (*AssignmentWithStatus, error)
	GetSubmission(id int) (*Submission, error)
	List(userId string, workspaceId int) ([]AssignmentWithStatus, error)
	ListSubmission(userId string, assignmentId int) ([]Submission, error)
}

type AssignmentUsecase interface {
	Create(userId string, workspaceId int, assignment *CreateAssignment) error
	Update(userId string, assignmentId int, assignment *UpdateAssignment) error
	CreateTestcases(assignmentId int, files []TestcaseFile) error
	UpdateTestcases(assignmentId int, files []TestcaseFile) error
	CreateSubmission(userId string, assignmentId int, workspaceId int, language string, file io.Reader) error
	CreateSubmissionResults(submissionId int, compilationLog string, results []SubmissionResult) error
	Get(id int) (*Assignment, error)
	GetWithStatus(id int, userId string) (*AssignmentWithStatus, error)
	CheckPerm(userId string, workspaceId int) (bool, error)
	GetSubmission(id int) (*Submission, error)
	List(userId string, workspaceId int) ([]AssignmentWithStatus, error)
	ListSubmission(userId string, assignmentId int) ([]Submission, error)
}
