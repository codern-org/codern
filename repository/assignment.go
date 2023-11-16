package repository

import (
	"database/sql"
	"fmt"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
)

type assignmentRepository struct {
	db *sqlx.DB
}

func NewAssignmentRepository(db *sqlx.DB) domain.AssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) CreateSubmission(
	submission *domain.Submission,
	testcases []domain.Testcase,
) error {
	_, err := r.db.NamedExec(`
		INSERT INTO submission (id, assignment_id, user_id, language, status, score, file_url)
		VALUES (:id, :assignment_id, :user_id, :language, 'GRADING', 0, :file_url)
	`, submission)
	if err != nil {
		return fmt.Errorf("cannot query to create submission: %w", err)
	}
	return nil
}

func (r *assignmentRepository) CreateSubmissionResults(
	submissionId int,
	compilationLog string,
	status domain.AssignmentStatus,
	score int,
	results []domain.SubmissionResult,
) (retErr error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("cannot begin transaction to update submission result: %w", err)
	}

	defer func() {
		if err := recover(); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				retErr = fmt.Errorf("cannot rollback transaction: %w", err.(error))
			} else {
				retErr = err.(error)
			}
		}
	}()

	_, err = tx.Exec(
		"UPDATE submission SET compilation_log = ?, status = ?, score = ? WHERE id = ?",
		compilationLog, status, score, submissionId,
	)
	if err != nil {
		panic(fmt.Errorf("cannot query to update submission from submission result: %w", err))
	}

	query := "INSERT INTO submission_result (submission_id, testcase_id, is_passed, status, memory_usage, time_usage) VALUES "
	for _, result := range results {
		query += fmt.Sprintf(
			"('%d', '%d', %t, '%s', '%d', '%d'),",
			result.SubmissionId, result.TestcaseId, result.IsPassed, result.Status,
			*result.MemoryUsage, *result.TimeUsage,
		)
	}
	query = query[:len(query)-1] // Remove trailing comma
	if _, err := tx.Exec(query); err != nil {
		panic(fmt.Errorf("cannot query to create submission result: %w", err))
	}

	if err = tx.Commit(); err != nil {
		panic(fmt.Errorf("cannot commit transaction to update submission result: %w", err))
	}

	return
}

func (r *assignmentRepository) Get(id int, userId string) (*domain.Assignment, error) {
	assignments, err := r.list(userId, nil, &id)
	if err != nil {
		return nil, fmt.Errorf("cannot query to get assignment: %w", err)
	} else if len(assignments) == 0 {
		return nil, nil
	}
	return &assignments[0], nil
}

func (r *assignmentRepository) GetSubmission(id int) (*domain.Submission, error) {
	var submission domain.Submission
	err := r.db.Get(&submission, "SELECT * FROM submission WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get submission: %w", err)
	}

	var results []domain.SubmissionResult
	query, args, err := sqlx.In("SELECT * FROM submission_result WHERE submission_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("cannot query to create query to list submission result: %w", err)
	}
	if err = r.db.Select(&results, query, args...); err != nil {
		return nil, fmt.Errorf("cannot query to list submission result: %w", err)
	}
	submission.Results = results
	return &submission, nil
}

func (r *assignmentRepository) List(userId string, workspaceId int) ([]domain.Assignment, error) {
	return r.list(userId, &workspaceId, nil)
}

func (r *assignmentRepository) list(
	userId string,
	workspaceId *int,
	assignmentId *int,
) ([]domain.Assignment, error) {
	assignments := make([]domain.Assignment, 0)

	query := `
		SELECT
			a.*,
			t1.last_submitted_at,
			IFNULL(t1.status, 'TODO') AS status
		FROM (
			SELECT
				s.assignment_id,
				MAX(s.submitted_at) AS last_submitted_at,
				CASE
					WHEN SUM(CASE WHEN s.status = 'GRADING' THEN 1 ELSE 0 END) > 0 THEN 'GRADING'
					WHEN COUNT(DISTINCT s.status) = 1 AND MIN(s.status) = 'COMPLETED' THEN 'COMPLETED'
					ELSE 'INCOMPLETED'
				END AS status
			FROM submission s
			WHERE s.user_id = ? AND s.assignment_id %[1]s
			GROUP BY s.assignment_id
		) t1
		RIGHT JOIN assignment a ON a.id = t1.assignment_id
		WHERE a.id %[1]s
	`

	whereAssignmentId := "IN (SELECT id FROM assignment WHERE workspace_id = ?)"
	param := workspaceId
	if assignmentId != nil {
		whereAssignmentId = "= ?"
		param = assignmentId
	}
	query = fmt.Sprintf(query, whereAssignmentId)

	if err := r.db.Select(&assignments, query, userId, param, param); err != nil {
		return nil, fmt.Errorf("cannot query to list assignment: %w", err)
	}

	if len(assignments) > 0 {
		var assignmentIds []int
		for i := range assignments {
			assignmentIds = append(assignmentIds, assignments[i].Id)
		}

		testcases, err := r.listTestcase(assignmentIds)
		if err != nil {
			return nil, fmt.Errorf("cannot query to list testcase for assignment: %w", err)
		}

		assignmentById := make(map[int]*domain.Assignment)
		for i := range assignments {
			assignmentById[assignments[i].Id] = &assignments[i]
		}
		for i := range testcases {
			assignment := assignmentById[testcases[i].AssignmentId]
			assignment.Testcases = append(assignment.Testcases, testcases[i])
		}
	}

	return assignments, nil
}

func (r *assignmentRepository) listTestcase(assignmentIds []int) ([]domain.Testcase, error) {
	var testcases []domain.Testcase
	query, args, err := sqlx.In("SELECT * FROM testcase WHERE assignment_id IN (?)", assignmentIds)
	if err != nil {
		return nil, fmt.Errorf("cannot query to create query to list testcase: %w", err)
	}
	if err = r.db.Select(&testcases, query, args...); err != nil {
		return nil, fmt.Errorf("cannot query to list testcase: %w", err)
	}
	return testcases, nil
}

func (r *assignmentRepository) ListSubmission(userId string, assignmentId int) ([]domain.Submission, error) {
	submissions := make([]domain.Submission, 0)
	err := r.db.Select(&submissions, "SELECT * FROM submission WHERE assignment_id = ? AND user_id = ?", assignmentId, userId)
	if err != nil {
		return nil, fmt.Errorf("cannot query to list submission: %w", err)
	}

	if len(submissions) > 0 {
		var submissionIds []int
		for i := range submissions {
			submissionIds = append(submissionIds, submissions[i].Id)
		}
		var results []domain.SubmissionResult
		query, args, err := sqlx.In("SELECT * FROM submission_result WHERE submission_id IN (?)", submissionIds)
		if err != nil {
			return nil, fmt.Errorf("cannot query to create query to list submission result: %w", err)
		}
		if err = r.db.Select(&results, query, args...); err != nil {
			return nil, fmt.Errorf("cannot query to list submission result: %w", err)
		}

		submissionById := make(map[int]*domain.Submission)
		assignment, err := r.Get(assignmentId, userId)
		if assignment == nil {
			return nil, fmt.Errorf("cannot get assignment due date to determine late flag: assignment not found")
		}
		if err != nil {
			return nil, fmt.Errorf("cannot get assignment due date to determine late flag: %w", err)
		}
		for i := range submissions {
			submissionById[submissions[i].Id] = &submissions[i]
			submission := submissionById[submissions[i].Id]
			if assignment.DueDate == nil {
				submission.IsLate = false
			} else if submission.SubmittedAt.After(*assignment.DueDate) {
				submission.IsLate = true
			} else {
				submission.IsLate = false
			}
		}

		for i := range results {
			submission := submissionById[results[i].SubmissionId]
			submission.Results = append(submission.Results, results[i])
		}
	}

	return submissions, nil
}
