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

func (r *assignmentRepository) CreateAssignment(assignment *domain.RawAssignment) error {
	_, err := r.db.NamedExec(`
		INSERT INTO assignment
			(id, workspace_id, name, description, detail_url, memory_limit, time_limit, level)
		VALUES
			(:id, :workspace_id, :name, :description, :detail_url, :memory_limit, :time_limit, :level)
		`, assignment)
	if err != nil {
		return fmt.Errorf("cannot query to insert assignment: %w", err)
	}

	return nil
}

func (r *assignmentRepository) UpdateAssignment(assignment *domain.RawAssignment) error {
	_, err := r.db.Exec(
		`UPDATE assignment SET name = ?, description = ?, detail_url = ?, memory_limit = ?, time_limit = ?, level = ? WHERE id = ?`,
		assignment.Name, assignment.Description, assignment.DetailUrl, assignment.MemoryLimit, assignment.TimeLimit, assignment.Level, assignment.Id,
	)

	if err != nil {
		return fmt.Errorf("cannot query to update assignment: %w", err)
	}

	return nil
}

func (r *assignmentRepository) CreateTestcases(testcases []domain.Testcase) error {
	query := "INSERT INTO testcase (id, assignment_id, input_file_url, output_file_url) VALUES "

	args := make([]interface{}, 0, len(testcases)*4)
	for i := range testcases {
		query += "(?, ?, ?, ?),"
		args = append(args, testcases[i].Id, testcases[i].AssignmentId, testcases[i].InputFileUrl, testcases[i].OutputFileUrl)
	}

	query = query[:len(query)-1]

	if _, err := r.db.Exec(query, args...); err != nil {
		return fmt.Errorf("cannot query to create testcase: %w", err)
	}

	return nil
}

func (r *assignmentRepository) DeleteTestcasesByAssignmentId(assignmentId int) error {
	_, err := r.db.Exec("DELETE FROM testcase WHERE assignment_id = ?", assignmentId)
	if err != nil {
		return fmt.Errorf("cannot query to delete testcase: %w", err)
	}
	return nil
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

func (r *assignmentRepository) GetRaw(id int) (*domain.RawAssignment, error) {
	assignments, err := r.listRaw(nil, &id)
	if err != nil {
		return nil, fmt.Errorf("cannot query to get raw assignment: %w", err)
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
					WHEN SUM(CASE WHEN s.status = 'COMPLETED' THEN 1 ELSE 0 END) > 0 THEN 'COMPLETED'
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
		params := make([]*domain.RawAssignment, 0, len(assignments))
		for i := range assignments {
			params = append(params, &assignments[i].RawAssignment)
		}
		if err := r.mutateTestcases(params); err != nil {
			return nil, err
		}
	}

	return assignments, nil
}

func (r *assignmentRepository) listRaw(
	workspaceId *int,
	assignmentId *int,
) ([]domain.RawAssignment, error) {
	rawAssignments := make([]domain.RawAssignment, 0)

	if workspaceId != nil {
		query := `SELECT * FROM assignment WHERE workspace_id = ?`
		if err := r.db.Select(&rawAssignments, query, workspaceId); err != nil {
			return nil, fmt.Errorf("cannot query to list raw assignment with workspace id: %w", err)
		}
	}

	if assignmentId != nil {
		query := `SELECT * FROM assignment WHERE id = ?`
		if err := r.db.Select(&rawAssignments, query, assignmentId); err != nil {
			return nil, fmt.Errorf("cannot query to list raw assignment with id: %w", err)
		}
	}

	if len(rawAssignments) > 0 {
		params := make([]*domain.RawAssignment, 0, len(rawAssignments))
		for i := range rawAssignments {
			params = append(params, &rawAssignments[i])
		}
		if err := r.mutateTestcases(params); err != nil {
			return nil, err
		}
	}
	return rawAssignments, nil
}

func (r *assignmentRepository) mutateTestcases(assignments []*domain.RawAssignment) error {
	var assignmentIds []int
	for i := range assignments {
		assignmentIds = append(assignmentIds, assignments[i].Id)
	}

	testcases, err := r.listTestcase(assignmentIds)
	if err != nil {
		return fmt.Errorf("cannot query to list testcase for assignment: %w", err)
	}

	assignmentById := make(map[int]*domain.RawAssignment)
	for i := range assignments {
		assignmentById[assignments[i].Id] = assignments[i]
	}
	for i := range testcases {
		assignment := assignmentById[testcases[i].AssignmentId]
		assignment.Testcases = append(assignment.Testcases, testcases[i])
	}

	return nil
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
	query := `
		SELECT
			s.*,
			CASE
				WHEN s.submitted_at > a.due_date THEN TRUE
				WHEN s.submitted_at < a.due_date THEN FALSE
				WHEN s.submitted_at = a.due_date THEN FALSE
				ELSE FALSE
			END AS is_late
		FROM submission s
		INNER JOIN assignment a on a.id = s.assignment_id
		WHERE s.assignment_id = ? AND s.user_id = ?
	`
	err := r.db.Select(&submissions, query, assignmentId, userId)
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
		for i := range submissions {
			submissionById[submissions[i].Id] = &submissions[i]
		}
		for i := range results {
			submission := submissionById[results[i].SubmissionId]
			submission.Results = append(submission.Results, results[i])
		}
	}

	return submissions, nil
}
