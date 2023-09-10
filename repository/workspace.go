package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
)

type workspaceRepository struct {
	db *sqlx.DB
}

func NewWorkspaceRepository(db *sqlx.DB) domain.WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) CreateSubmission(
	submission *domain.Submission,
	testcases []domain.Testcase,
) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("cannot begin transaction to create submission: %w", err)
	}

	_, err = tx.NamedExec("INSERT INTO submission (id, assignment_id, user_id, language, file_url) VALUES (:id, :assignment_id, :user_id, :language, :file_url)", submission)
	if err != nil {
		return fmt.Errorf("cannot query to insert submission: %w", err)
	}

	query := "INSERT INTO submission_result (submission_id, testcase_id, status) VALUES "
	for i := range testcases {
		query += fmt.Sprintf("('%d', '%d', '%s'),", submission.Id, testcases[i].Id, "GRADING")
	}
	query = query[:len(query)-1]

	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("cannot execute transaction to create submission: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("cannot commit transaction to create submission: %w", err)
	}

	return nil
}

func (r *workspaceRepository) IsUserIn(userId string, workspaceId int) (bool, error) {
	var result domain.WorkspaceParticipant
	err := r.db.Get(
		&result,
		"SELECT workspace_id FROM workspace_participant WHERE workspace_id = ? AND user_id = ? LIMIT 1",
		workspaceId, userId,
	)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("cannot query to check user in workspace participant: %w", err)
	}
	return true, nil
}

func (r *workspaceRepository) IsAssignmentIn(assignmentId int, workspaceId int) (bool, error) {
	var result domain.Assignment
	err := r.db.Get(
		&result,
		"SELECT id FROM assignment WHERE id = ? AND workspace_id = ? LIMIT 1",
		assignmentId, workspaceId,
	)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("cannot query to check assignment in workspace: %w", err)
	}
	return true, nil
}

func (r *workspaceRepository) Get(id int, selector *domain.WorkspaceSelector) (*domain.Workspace, error) {
	workspaces, err := r.list([]int{id}, selector)
	if err != nil {
		return nil, fmt.Errorf("cannot query to get workspace: %w", err)
	}
	return &workspaces[0], nil
}

func (r *workspaceRepository) GetAssignment(
	id int,
	userId string,
	workspaceId int,
) (*domain.Assignment, error) {
	assignments, err := r.listAssignment(userId, workspaceId, &id)
	if err != nil {
		return nil, fmt.Errorf("cannot query to get assignment: %w", err)
	}
	return &assignments[0], nil
}

func (r *workspaceRepository) List(
	userId string,
	selector *domain.WorkspaceSelector,
) ([]domain.Workspace, error) {
	var workspaceIds []int
	err := r.db.Select(&workspaceIds, "SELECT workspace_id FROM workspace_participant WHERE user_id = ?", userId)
	if err != nil {
		return nil, fmt.Errorf("cannot query to list workspace id: %w", err)
	}
	return r.list(workspaceIds, selector)
}

func (r *workspaceRepository) list(ids []int, selector *domain.WorkspaceSelector) ([]domain.Workspace, error) {
	workspaces := make([]domain.Workspace, 0)
	if len(ids) == 0 {
		return workspaces, nil
	}

	query, args, err := sqlx.In(`
		SELECT
			w.*,
			user.display_name AS owner_name,
			(SELECT COUNT(*) FROM workspace_participant wp WHERE wp.workspace_id = w.id) AS participant_count,
			(SELECT COUNT(*) FROM assignment a WHERE a.workspace_id = w.id) AS total_assignment
		FROM workspace w
		INNER JOIN user ON user.id = w.owner_id
		WHERE w.id IN (?)
		LIMIT ?
	`, ids, len(ids))
	if err != nil {
		return nil, fmt.Errorf("cannot query to create query to list workspace: %w", err)
	}
	if err := r.db.Select(&workspaces, query, args...); err != nil {
		return nil, fmt.Errorf("cannot query to list workspace: %w", err)
	}

	if selector.Participants {
		participants := make([]domain.WorkspaceParticipant, 0)
		query, args, err := sqlx.In("SELECT * FROM workspace_participant WHERE workspace_id IN (?)", ids)
		if err != nil {
			return nil, fmt.Errorf("cannot query to create query to list workspace participant: %w", err)
		}
		if err := r.db.Select(&participants, query, args...); err != nil {
			return nil, fmt.Errorf("cannot query to list workspace participant: %w", err)
		}

		for i := range workspaces {
			workspace := &workspaces[i]
			for j := range participants {
				if workspace.Id == participants[j].WorkspaceId {
					workspace.Participants = append(workspace.Participants, participants[j])
				}
			}
		}
	}

	return workspaces, nil
}

func (r *workspaceRepository) ListRecent(userId string) ([]domain.Workspace, error) {
	workspaces := make([]domain.Workspace, 0)
	query, args, err := sqlx.In(`
		SELECT
			w.*,
			user.display_name AS owner_name,
			(SELECT COUNT(*) FROM workspace_participant wp WHERE wp.workspace_id = w.id) AS participant_count,
			(SELECT COUNT(*) FROM assignment a WHERE a.workspace_id = w.id) AS total_assignment
		FROM workspace w
		INNER JOIN user ON user.id = w.owner_id
		INNER JOIN workspace_participant wp ON wp.workspace_id = w.id
		WHERE wp.user_id = ? AND w.id IN (SELECT workspace_id FROM workspace_participant WHERE user_id = ?)
		ORDER BY wp.recently_visited_at DESC
		LIMIT 4
	`, userId, userId)
	if err != nil {
		return nil, fmt.Errorf("cannot query to create query to list recent workspace: %w", err)
	}
	if err := r.db.Select(&workspaces, query, args...); err != nil {
		return nil, fmt.Errorf("cannot query to list recent workspace: %w", err)
	}
	return workspaces, nil
}

func (r *workspaceRepository) ListAssignment(userId string, workspaceId int) ([]domain.Assignment, error) {
	return r.listAssignment(userId, workspaceId, nil)
}

func (r *workspaceRepository) listAssignment(
	userId string,
	workspaceId int,
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
		param = *assignmentId
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

func (r *workspaceRepository) listTestcase(assignmentIds []int) ([]domain.Testcase, error) {
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

func (r *workspaceRepository) ListSubmission(userId string, assignmentId int) ([]domain.Submission, error) {
	submissions := make([]domain.Submission, 0)
	err := r.db.Select(&submissions, "SELECT * FROM submission WHERE assignment_id = ?", assignmentId)
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

func (r *workspaceRepository) UpdateRecent(userId string, workspaceId int) error {
	_, err := r.db.Exec(`
		UPDATE workspace_participant SET recently_visited_at = ? WHERE user_id = ? AND workspace_id = ?
	`, time.Now(), userId, workspaceId)
	if err != nil {
		return fmt.Errorf("cannot query to update recent workspace: %w", err)
	}
	return nil
}

func (r *workspaceRepository) UpdateSubmissionResult(result *domain.SubmissionResult) error {
	_, err := r.db.NamedExec(`
		UPDATE submission_result
		SET
			status = :status,
			status_detail = :status_detail,
			memory_usage = :memory_usage,
			time_usage = :time_usage,
			compilation_log = :compilation_log
		WHERE submission_id = :submission_id AND testcase_id = :testcase_id
	`, result)
	if err != nil {
		return fmt.Errorf("cannot query to update submission result: %w", err)
	}
	return nil
}
