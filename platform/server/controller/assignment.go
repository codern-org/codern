package controller

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
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

type TestcaseFileHeader struct {
	Input  *multipart.FileHeader
	Output *multipart.FileHeader
}

// CreateAssignment godoc
//
// @Summary 		Create a new assignment
// @Description	Create a new assignment on workspace id on path parameter
// @Tags 				workspace
// @Accept 			multipart/form-data
// @Produce 		json
// @Param				workspaceId					path	int				true	"Workspace ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Param 			name formData string true "Assignment name"
// @Param 			description formData string true "Assignment description"
// @Param 			memory_limit formData int true "Assignment memory limit"
// @Param 			time_limit formData int true "Assignment time limit"
// @Param 			level formData int true "Assignment level"
// @Param 			detail formData file true "Assignment detail"
// @Param 			in formData file true "Assignment testcase input"
// @Param 			out formData file true "Assignment testcase output"
// @Router 			/api/workspaces/{workspaceId}/assignments [post]
func (c *AssignmentController) CreateAssignment(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	testcaseHeaderByIndex := make(map[int]*TestcaseFileHeader)

	for formFieldName, fileHeaders := range form.File {
		if formFieldName != "in" && formFieldName != "out" {
			continue
		}

		for _, fileHeader := range fileHeaders {
			sequenceStr := strings.Split(fileHeader.Filename, ".")[0]
			sequence, err := strconv.Atoi(sequenceStr)
			if err != nil {
				return response.NewSuccessResponse(ctx, fiber.StatusBadRequest, fiber.Map{
					"message": fmt.Sprintf("testcase file name `%s` is not valid", fileHeader.Filename),
				})
			}

			testcaseFile, found := testcaseHeaderByIndex[sequence]
			if !found {
				testcaseFile = &TestcaseFileHeader{}
				testcaseHeaderByIndex[sequence] = testcaseFile
			}

			switch formFieldName {
			case "in":
				testcaseFile.Input = fileHeader
			case "out":
				testcaseFile.Output = fileHeader
			}
		}
	}

	for _, testcaseFileHeader := range testcaseHeaderByIndex {
		if testcaseFileHeader.Input == nil || testcaseFileHeader.Output == nil {
			return response.NewSuccessResponse(ctx, fiber.StatusBadRequest, fiber.Map{
				"message": "testcase file is not complete",
			})
		}
	}

	for i := 1; i <= len(testcaseHeaderByIndex); i++ {
		_, found := testcaseHeaderByIndex[i]
		if !found {
			return response.NewSuccessResponse(ctx, fiber.StatusBadRequest, fiber.Map{
				"message": "testcase file is not sequential",
			})
		}
	}

	testcaseFiles := make([]domain.TestcaseFile, 0)
	for _, testcaseFileHeader := range testcaseHeaderByIndex {
		input, err := testcaseFileHeader.Input.Open()
		if err != nil {
			return err
		}

		output, err := testcaseFileHeader.Output.Open()
		if err != nil {
			return err
		}

		testcaseFiles = append(testcaseFiles, domain.TestcaseFile{
			Input:  input,
			Output: output,
		})
	}

	var payload payload.CreateAssignmentPayload
	if ok, err := c.validator.Validate(&payload, ctx); !ok {
		return err
	}

	workspaceId, _ := ctx.ParamsInt("workspaceId")

	assignment, err := c.assignmentUsecase.CreateAssigment(workspaceId, payload.Name, payload.Description, payload.MemoryLimit, payload.TimeLimit, payload.Level, payload.DetailFile)
	if err != nil {
		return err
	}

	err = c.assignmentUsecase.CreateTestcase(assignment.Id, testcaseFiles)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"submitted_at": time.Now(),
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
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/workspaces/{workspaceId}/assignments/{assignmentId}/submissions [post]
func (c *AssignmentController) CreateSubmission(ctx *fiber.Ctx) error {
	var pl payload.CreateSubmissionPayload
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)

	err := c.assignmentUsecase.CreateSubmission(user.Id, pl.AssignmentId, workspaceId, pl.Language, pl.SourceCode)
	if err != nil {
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
	user := middleware.GetUserFromCtx(ctx)
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)

	assignments, err := c.assignmentUsecase.List(user.Id, workspaceId)
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
	user := middleware.GetUserFromCtx(ctx)
	assignmentId := middleware.GetAssignmentIdFromCtx(ctx)

	submissions, err := c.assignmentUsecase.ListSubmission(user.Id, assignmentId)
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
	user := middleware.GetUserFromCtx(ctx)
	assignmentId := middleware.GetAssignmentIdFromCtx(ctx)

	assignment, err := c.assignmentUsecase.Get(assignmentId, user.Id)
	if err != nil {
		return err
	} else if assignment == nil {
		return errs.New(errs.ErrAssignmentNotFound, "assignment id %d not found", assignmentId)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignment)
}
