package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
)

type assignmentRepository struct {
	db *sqlx.DB
}

func NewAssignmentRepository(db *sqlx.DB) domain.AssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) CreateAssigment(assignment *domain.Assignment) error {
	_, err := r.db.NamedExec("INSERT INTO assignment (id, workspace_id, name, description, detail_url, memory_limit, time_limit, level) VALUES (:id, :workspace_id, :name, :description, :detail_url, :memory_limit, :time_limit, :level)", assignment)
	if err != nil {
		return fmt.Errorf("cannot query to insert assignment: %w", err)
	}

	return nil
}

func (r *assignmentRepository) CreateTestcases(testcases []domain.Testcase) error {
	query := "INSERT INTO testcase (id, assignment_id, input_file_url, output_file_url) VALUES "
	for i := range testcases {
		query += fmt.Sprintf("('%d', '%d', '%s', '%s'),", testcases[i].Id, testcases[i].AssignmentId, testcases[i].InputFileUrl, testcases[i].OutputFileUrl)
	}
	query = query[:len(query)-1]

	if _, err := r.db.Exec(query); err != nil {
		return fmt.Errorf("cannot query to create testcase: %w", err)
	}

	return nil
}

func (r *assignmentRepository) CreateSubmission(
	submission *domain.Submission,
	testcases []domain.Testcase,
) (retErr error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("cannot begin transaction to create submission: %w", err)
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

	_, err = tx.NamedExec("INSERT INTO submission (id, assignment_id, user_id, language, file_url) VALUES (:id, :assignment_id, :user_id, :language, :file_url)", submission)
	if err != nil {
		panic(fmt.Errorf("cannot query to insert submission: %w", err))
	}

	query := "INSERT INTO submission_result (submission_id, testcase_id, status) VALUES "
	for i := range testcases {
		query += fmt.Sprintf("('%d', '%d', '%s'),", submission.Id, testcases[i].Id, "GRADING")
	}
	query = query[:len(query)-1]

	if _, err := tx.Exec(query); err != nil {
		panic(fmt.Errorf("cannot execute transaction to create submission: %w", err))
	}

	if err = tx.Commit(); err != nil {
		panic(fmt.Errorf("cannot commit transaction to create submission: %w", err))
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
			t2.last_submitted_at,
			IFNULL(t2.status, 'TODO') AS status
		FROM (
			SELECT
				a.*,
				MAX(t1.submitted_at) AS last_submitted_at,
				CASE
					WHEN SUM(CASE WHEN t1.status = "GRADING" THEN 1 ELSE 0 END) > 0 THEN "GRADING"
					WHEN SUM(CASE WHEN t1.status = "DONE" THEN 1 ELSE 0 END) > 0 THEN 'DONE'
					ELSE "ERROR"
				END as status
			FROM (
				SELECT
					s.id,
					s.assignment_id,
					s.submitted_at,
					CASE
						WHEN SUM(CASE WHEN sr.status = "GRADING" THEN 1 ELSE 0 END) > 0 THEN "GRADING"
						WHEN COUNT(DISTINCT sr.status) = 1 AND MIN(sr.status) = 'DONE' THEN 'DONE'
						ELSE "ERROR"
					END as status
				FROM submission s
				INNER JOIN submission_result sr ON sr.submission_id = s.id
				WHERE s.user_id = ? AND s.assignment_id %[1]s
				GROUP BY sr.submission_id
			) t1
			INNER JOIN assignment a ON a.id = t1.assignment_id
			GROUP BY a.id
		) t2
		RIGHT JOIN assignment a ON a.id = t2.id
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

func (r *assignmentRepository) ListSubmission(filter *domain.SubmissionFilter) ([]domain.Submission, error) {
	submissions := make([]domain.Submission, 0)

	query := "SELECT * FROM submission "
	args := make([]interface{}, 0)

	conditions := make([]string, 0)

	if filter.AssignmentId != nil {
		conditions = append(conditions, "assignment_id = ?")
		args = append(args, *filter.AssignmentId)
	}
	if filter.UserId != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, *filter.UserId)
	}

	where := strings.Join(conditions, " AND ")
	if where != "" {
		query += "WHERE " + where
	}

	err := r.db.Select(&submissions, query, args...)
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

func (r *assignmentRepository) UpdateSubmissionResults(
	submissionId int,
	compilationLog string,
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

	_, err = tx.Exec("UPDATE submission SET compilation_log = ? WHERE id = ?", compilationLog, submissionId)
	if err != nil {
		panic(fmt.Errorf("cannot query to update submission from submission result: %w", err))
	}

	// TODO: optimization
	for i := range results {
		_, err := tx.NamedExec(`
			UPDATE submission_result
			SET
				status = :status,
				status_detail = :status_detail,
				memory_usage = :memory_usage,
				time_usage = :time_usage
			WHERE submission_id = :submission_id AND testcase_id = :testcase_id;
		`, results[i])
		if err != nil {
			panic(fmt.Errorf("cannot query to update submission result: %w", err))
		}
	}

	if err = tx.Commit(); err != nil {
		panic(fmt.Errorf("cannot commit transaction to update submission result: %w", err))
	}

	return
}
