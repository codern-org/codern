package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform"
	"github.com/jmoiron/sqlx"
)

type assignmentRepository struct {
	db *platform.MySql
}

func NewAssignmentRepository(db *platform.MySql) domain.AssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) Create(assignment *domain.Assignment) error {
	_, err := r.db.NamedExec(`
		INSERT INTO assignment
			(id, workspace_id, name, description, detail_url, memory_limit, time_limit, level, publish_date, due_date)
		VALUES
			(:id, :workspace_id, :name, :description, :detail_url, :memory_limit, :time_limit, :level, :publish_date, :due_date)
		`, assignment)
	if err != nil {
		return fmt.Errorf("cannot query to insert assignment: %w", err)
	}

	return nil
}

func (r *assignmentRepository) Update(assignment *domain.Assignment) error {
	_, err := r.db.NamedExec(`
		UPDATE assignment SET
			name = :name,
			description = :description,
			detail_url = :detail_url,
			memory_limit = :memory_limit,
			time_limit = :time_limit,
			level = :level,
			publish_date = :publish_date,
			due_date = :due_date
		WHERE id = :id
	`, assignment)

	if err != nil {
		return fmt.Errorf("cannot query to update assignment: %w", err)
	}

	return nil
}

func (r *assignmentRepository) Delete(id int) error {
	_, err := r.db.Exec("UPDATE assignment SET is_deleted = TRUE WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("cannot query to soft delete assignment: %w", err)
	}
	return nil
}

func (r *assignmentRepository) CreateTestcases(testcases []domain.Testcase) error {
	var revision int
	err := r.db.Get(
		&revision,
		"SELECT MAX(revision) AS revision FROM testcase WHERE assignment_id = ? GROUP BY assignment_id",
		testcases[0].AssignmentId,
	)
	if err == sql.ErrNoRows {
		revision = 1
	} else if err != nil {
		return fmt.Errorf("cannot query revision to create testcase: %w", err)
	} else {
		revision += 1
	}

	query := "INSERT INTO testcase (id, assignment_id, revision, input_file_url, output_file_url) VALUES "
	args := make([]interface{}, 0, len(testcases)*4)
	for _, testcase := range testcases {
		query += "(?, ?, ?, ?, ?),"
		args = append(args, testcase.Id, testcase.AssignmentId, revision, testcase.InputFileUrl, testcase.OutputFileUrl)
	}

	query = query[:len(query)-1]

	if _, err := r.db.Exec(query, args...); err != nil {
		return fmt.Errorf("cannot query to create testcase: %w", err)
	}

	return nil
}

func (r *assignmentRepository) DeleteTestcases(assignmentId int) error {
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
	score float64,
	results []domain.SubmissionResult,
) error {
	return r.db.ExecuteTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			"UPDATE submission SET compilation_log = ?, status = ?, score = ? WHERE id = ?",
			compilationLog, status, score, submissionId,
		)
		if err != nil {
			return fmt.Errorf("cannot query to update submission from submission result: %w", err)
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
			return fmt.Errorf("cannot query to create submission result: %w", err)
		}

		return nil
	})
}

func (r *assignmentRepository) GetWithStatus(id int, userId string) (*domain.AssignmentWithStatus, error) {
	assignments, err := r.list(userId, nil, &id)
	if err != nil {
		return nil, fmt.Errorf("cannot query to get assignment: %w", err)
	} else if len(assignments) == 0 {
		return nil, nil
	}
	return &assignments[0], nil
}

func (r *assignmentRepository) Get(id int) (*domain.Assignment, error) {
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

func (r *assignmentRepository) List(userId string, workspaceId int) ([]domain.AssignmentWithStatus, error) {
	return r.list(userId, &workspaceId, nil)
}

func (r *assignmentRepository) list(
	userId string,
	workspaceId *int,
	assignmentId *int,
) ([]domain.AssignmentWithStatus, error) {
	assignments := make([]domain.AssignmentWithStatus, 0)

	query := `
		SELECT
			a.*,
			t1.last_submitted_at,
			IFNULL(t1.status, 'TODO') AS status,
			t1.score
		FROM (
			SELECT
				s.assignment_id,
				MAX(s.submitted_at) AS last_submitted_at,
				CASE
					WHEN SUM(CASE WHEN s.status = 'GRADING' THEN 1 ELSE 0 END) > 0 THEN 'GRADING'
					WHEN SUM(CASE WHEN s.status = 'COMPLETED' THEN 1 ELSE 0 END) > 0 THEN 'COMPLETED'
					ELSE 'INCOMPLETED'
				END AS status,
				MAX(s.score) AS score
			FROM submission s
			WHERE s.user_id = ? AND s.assignment_id %[1]s
			GROUP BY s.assignment_id
		) t1
		RIGHT JOIN assignment a ON a.id = t1.assignment_id
		WHERE a.id %[1]s
	`

	whereAssignmentId := "IN (SELECT id FROM assignment WHERE workspace_id = ? AND is_deleted = FALSE)"
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
		params := make([]*domain.Assignment, 0, len(assignments))
		for i := range assignments {
			params = append(params, &assignments[i].Assignment)
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
) ([]domain.Assignment, error) {
	rawAssignments := make([]domain.Assignment, 0)

	if workspaceId != nil {
		query := `SELECT * FROM assignment WHERE workspace_id = ? AND is_deleted = FALSE`
		if err := r.db.Select(&rawAssignments, query, workspaceId); err != nil {
			return nil, fmt.Errorf("cannot query to list raw assignment with workspace id: %w", err)
		}
	}

	if assignmentId != nil {
		query := `SELECT * FROM assignment WHERE id = ? AND is_deleted = FALSE`
		if err := r.db.Select(&rawAssignments, query, assignmentId); err != nil {
			return nil, fmt.Errorf("cannot query to list raw assignment with id: %w", err)
		}
	}

	if len(rawAssignments) > 0 {
		params := make([]*domain.Assignment, 0, len(rawAssignments))
		for i := range rawAssignments {
			params = append(params, &rawAssignments[i])
		}
		if err := r.mutateTestcases(params); err != nil {
			return nil, err
		}
	}
	return rawAssignments, nil
}

func (r *assignmentRepository) mutateTestcases(assignments []*domain.Assignment) error {
	var assignmentIds []int
	for i := range assignments {
		assignmentIds = append(assignmentIds, assignments[i].Id)
	}

	testcases, err := r.listTestcase(assignmentIds)
	if err != nil {
		return fmt.Errorf("cannot query to list testcase for assignment: %w", err)
	}

	assignmentById := make(map[int]*domain.Assignment)
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
	query, args, err := sqlx.In(`
		WITH assignment_latest_revision AS (
			SELECT assignment_id, MAX(revision) AS lastet_revision
			FROM testcase
			WHERE assignment_id IN (?)
			GROUP BY assignment_id
		)
		SELECT testcase.*
		FROM assignment_latest_revision t1
		INNER JOIN testcase ON testcase.assignment_id = t1.assignment_id AND revision = t1.lastet_revision
	`, assignmentIds)
	if err != nil {
		return nil, fmt.Errorf("cannot query to create query to list testcase: %w", err)
	}
	if err = r.db.Select(&testcases, query, args...); err != nil {
		return nil, fmt.Errorf("cannot query to list testcase: %w", err)
	}
	return testcases, nil
}

func (r *assignmentRepository) ListSubmission(
	userId *string,
	assignmentId *int,
) ([]domain.Submission, error) {
	submissions := make([]domain.Submission, 0)

	queryArgs := make([]interface{}, 0)
	whereQueries := make([]string, 0)

	if userId != nil {
		queryArgs = append(queryArgs, userId)
		whereQueries = append(whereQueries, "s.user_id = ?")
	}

	if assignmentId != nil {
		queryArgs = append(queryArgs, assignmentId)
		whereQueries = append(whereQueries, "s.assignment_id = ?")
	}

	whereQueryString := fmt.Sprintf("WHERE %s", strings.Join(whereQueries, " AND "))
	query := fmt.Sprintf(`
		SELECT
			s.*,
			u.display_name AS user_display_name,
			u.profile_url AS user_profile_url,
			CASE
				WHEN s.submitted_at > a.due_date THEN TRUE
				WHEN s.submitted_at < a.due_date THEN FALSE
				WHEN s.submitted_at = a.due_date THEN FALSE
				ELSE FALSE
			END AS is_late
		FROM submission s
		INNER JOIN user u ON u.id = s.user_id
		INNER JOIN assignment a ON a.id = s.assignment_id
		%s
	`, whereQueryString)

	err := r.db.Select(&submissions, query, queryArgs...)
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
