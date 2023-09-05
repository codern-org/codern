package repository

import (
	"database/sql"
	"fmt"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
)

type workspaceRepository struct {
	db *sqlx.DB
}

func NewWorkspaceRepository(db *sqlx.DB) domain.WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) CreateSubmission(submission *domain.Submission) error {
	var testcases []domain.Testcase
	err := r.db.Select(&testcases, "SELECT * FROM testcase WHERE assignment_id = ?", submission.AssignmentId)
	if err != nil {
		return err
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec("INSERT INTO submission (id, assignment_id, user_id, language, file_url) VALUES (:id, :assignment_id, :user_id, :language, :file_url)", submission)
	if err != nil {
		return err
	}

	query := "INSERT INTO submission_result (submission_id, testcase_id, status) VALUES "
	for i := range testcases {
		query += fmt.Sprintf("('%d', '%d', '%s'),", submission.Id, testcases[i].Id, "GRADING")
	}
	query = query[:len(query)-1]

	if _, err := tx.Exec(query); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
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
		return false, err
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
		return false, err
	}
	return true, nil
}

func (r *workspaceRepository) Get(id int, selector *domain.WorkspaceSelector) (*domain.Workspace, error) {
	workspaces, err := r.list([]int{id}, selector)
	if err != nil {
		return nil, err
	} else if len(*workspaces) != 1 {
		return nil, nil
	}
	return &(*workspaces)[0], nil
}

func (r *workspaceRepository) GetAssignment(
	id int,
	userId string,
	workspaceId int,
) (*domain.Assignment, error) {
	assignments, err := r.listAssignment(userId, workspaceId, &id)
	if err != nil {
		return nil, err
	}
	return &(*assignments)[0], nil
}

func (r *workspaceRepository) GetSubmission(id int) (*domain.Submission, error) {
	var submission domain.Submission
	err := r.db.Get(&submission, "SELECT * FROM submission WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var results []domain.SubmissionResult
	err = r.db.Select(&results, "SELECT * FROM submission_result WHERE submission_id = ?", id)
	if err != nil {
		return nil, err
	}
	submission.Results = results

	return &submission, nil
}

func (r *workspaceRepository) List(
	userId string,
	selector *domain.WorkspaceSelector,
) (*[]domain.Workspace, error) {
	var workspaceIds []int
	err := r.db.Select(&workspaceIds, "SELECT workspace_id FROM workspace_participant WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	return r.list(workspaceIds, selector)
}

func (r *workspaceRepository) list(ids []int, selector *domain.WorkspaceSelector) (*[]domain.Workspace, error) {
	workspaces := make([]domain.Workspace, 0)
	if len(ids) == 0 {
		return &workspaces, nil
	}

	query, args, err := sqlx.In(`
		SELECT
			w.*,
			user.display_name AS owner_name,
			count(wp.user_id) AS participant_count,
			count(a.id) AS total_assignment
		FROM workspace w
		INNER JOIN user ON user.id = w.owner_id
		INNER JOIN workspace_participant wp ON wp.workspace_id = w.id
		LEFT JOIN assignment a ON a.workspace_id = w.id
		WHERE w.id IN (?)
		GROUP BY w.id
		LIMIT ?
		`, ids, len(ids),
	)
	if err != nil {
		return nil, err
	}
	if err := r.db.Select(&workspaces, query, args...); err != nil {
		return nil, err
	}

	if selector.Participants {
		participants := make([]domain.WorkspaceParticipant, 0)
		query, args, err := sqlx.In("SELECT * FROM workspace_participant WHERE workspace_id IN (?)", ids)
		if err != nil {
			return nil, err
		}
		if err := r.db.Select(&participants, query, args...); err != nil {
			return nil, err
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

	return &workspaces, nil
}

func (r *workspaceRepository) ListAssignment(userId string, workspaceId int) (*[]domain.Assignment, error) {
	return r.listAssignment(userId, workspaceId, nil)
}

func (r *workspaceRepository) listAssignment(
	userId string,
	workspaceId int,
	assignmentId *int,
) (*[]domain.Assignment, error) {
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
		return nil, err
	}

	if len(assignments) > 0 {
		var assignmentIds []int
		for i := range assignments {
			assignmentIds = append(assignmentIds, assignments[i].Id)
		}

		var testcases []domain.Testcase
		query, args, err := sqlx.In("SELECT * FROM testcase WHERE assignment_id IN (?)", assignmentIds)
		if err != nil {
			return nil, err
		}
		if err = r.db.Select(&testcases, query, args...); err != nil {
			return nil, err
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

	return &assignments, nil
}

func (r *workspaceRepository) ListSubmission(userId string, assignmentId int) (*[]domain.Submission, error) {
	submissions := make([]domain.Submission, 0)
	err := r.db.Select(&submissions, "SELECT * FROM submission WHERE assignment_id = ?", assignmentId)
	if err != nil {
		return nil, err
	}

	if len(submissions) > 0 {
		var submissionIds []int
		for i := range submissions {
			submissionIds = append(submissionIds, submissions[i].Id)
		}

		var results []domain.SubmissionResult
		query, args, err := sqlx.In("SELECT * FROM submission_result WHERE submission_id IN (?)", submissionIds)
		if err != nil {
			return nil, err
		}
		if err = r.db.Select(&results, query, args...); err != nil {
			return nil, err
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

	return &submissions, nil
}
