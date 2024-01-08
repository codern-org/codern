package controller

import (
	"time"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/platform/server/middleware"
	"github.com/codern-org/codern/platform/server/payload"
	"github.com/codern-org/codern/platform/server/response"
	"github.com/gofiber/fiber/v2"
)

type AssignmentController struct {
	validator domain.PayloadValidator

	assignmentUsecase domain.AssignmentUsecase
}

func NewAssignmentController(
	validator domain.PayloadValidator,
	assignmentUsecase domain.AssignmentUsecase,
) *AssignmentController {
	return &AssignmentController{
		validator:         validator,
		assignmentUsecase: assignmentUsecase,
	}
}

func (c *AssignmentController) Create(ctx *fiber.Ctx) error {
	var pl payload.CreateAssignmentPayload
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}
	if err := payload.ValidateTestcaseFiles(pl.TestcaseInputFiles, pl.TestcaseOutputFiles); err != nil {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)
	testcaseFiles := domain.CreateTestcaseFiles(pl.TestcaseInputFiles, pl.TestcaseOutputFiles)

	if err := c.assignmentUsecase.Create(
		user.Id,
		pl.WorkspaceId,
		&domain.CreateAssignment{
			Name:          pl.Name,
			Description:   pl.Description,
			MemoryLimit:   pl.MemoryLimit,
			TimeLimit:     pl.TimeLimit,
			Level:         pl.Level,
			DetailFile:    pl.DetailFile,
			TestcaseFiles: testcaseFiles,
		},
	); err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"created_at": time.Now(),
	})
}

// CreateSubmission godoc
//
// @Summary 		Create a new submission
// @Description	Submit a submission of the assignment
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				workspaceId					path	int				true	"Workspace ID"
// @Param				assignmentId				path	int				true	"Assignment ID"
// @Param				payload							body	payload.CreateSubmissionPayload true "Payload"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/workspaces/{workspaceId}/assignments/{assignmentId}/submissions [post]
func (c *AssignmentController) CreateSubmission(ctx *fiber.Ctx) error {
	var pl payload.CreateSubmissionPayload
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	if err := c.assignmentUsecase.CreateSubmission(
		user.Id,
		pl.AssignmentId,
		pl.WorkspaceId,
		pl.Language,
		pl.SourceCode,
	); err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"submitted_at": time.Now(),
	})
}

// List godoc
//
// @Summary 		List assignment
// @Description	Get all assignment from a workspace id on path parameter
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				workspaceId					path	int				true	"Workspace ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/workspaces/{workspaceId}/assignments [get]
func (c *AssignmentController) List(ctx *fiber.Ctx) error {
	var pl payload.WorkspacePath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	assignments, err := c.assignmentUsecase.List(user.Id, pl.WorkspaceId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignments)
}

// ListSubmission godoc
//
// @Summary 		List submission
// @Description	Get all submission from a workspace id on path parameter
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				workspaceId					path	int				true	"Workspace ID"
// @Param				assignmentId				path	int				true	"Assignment ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/workspaces/{workspaceId}/assignments/{assignmentId}/submissions [get]
func (c *AssignmentController) ListSubmission(ctx *fiber.Ctx) error {
	var pl payload.AssignmentPath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	submissions, err := c.assignmentUsecase.ListSubmission(user.Id, pl.AssignmentId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, submissions)
}

// Get godoc
//
// @Summary 		Get an assignment
// @Description	Get an assignment from workspace id on path parameter
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				workspaceId					path	int				true	"Workspace ID"
// @Param				assignmentId				path	int				true	"Assignment ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/workspaces/{workspaceId}/assignments/{assignmentId} [get]
func (c *AssignmentController) Get(ctx *fiber.Ctx) error {
	var pl payload.AssignmentPath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	assignment, err := c.assignmentUsecase.GetWithStatus(pl.AssignmentId, user.Id)
	if err != nil {
		return err
	} else if assignment == nil {
		return errs.New(errs.ErrAssignmentNotFound, "assignment id %d not found", pl.AssignmentId)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignment)
}

func (c *AssignmentController) Update(ctx *fiber.Ctx) error {
	var pl payload.UpdateAssignment
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}
	if err := payload.ValidateTestcaseFiles(pl.TestcaseInputFiles, pl.TestcaseOutputFiles); err != nil {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)
	testcaseFiles := domain.CreateTestcaseFiles(pl.TestcaseInputFiles, pl.TestcaseOutputFiles)

	if err := c.assignmentUsecase.Update(
		user.Id,
		pl.AssignmentId,
		&domain.UpdateAssignment{
			Name:          pl.Name,
			Description:   pl.Description,
			MemoryLimit:   pl.MemoryLimit,
			TimeLimit:     pl.TimeLimit,
			Level:         pl.Level,
			DetailFile:    pl.DetailFile,
			TestcaseFiles: &testcaseFiles,
		},
	); err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"updated_at": time.Now(),
	})
}

func (c *AssignmentController) Delete(ctx *fiber.Ctx) error {
	var pl payload.DeleteAssignment
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	if err := c.assignmentUsecase.Delete(user.Id, pl.AssignmentId); err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"deleted_at": time.Now(),
	})
}
