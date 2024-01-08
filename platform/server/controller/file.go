package controller

import (
	"fmt"
	"net/url"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/platform/server/middleware"
	"github.com/codern-org/codern/platform/server/payload"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type FileController struct {
	validator        domain.PayloadValidator
	filerUrl         string
	WorkspaceUsecase domain.WorkspaceUsecase
}

func NewFileController(
	cfg *config.Config,
	validator domain.PayloadValidator,
	WorkspaceUsecase domain.WorkspaceUsecase,
) *FileController {
	return &FileController{
		validator:        validator,
		filerUrl:         cfg.Client.SeaweedFs.FilerUrls.Internal,
		WorkspaceUsecase: WorkspaceUsecase,
	}
}

// GetUserProfile godoc
//
// @Summary 		Get an user profile image
// @Description Get an user profile image from internal file system by proxy
// @Tags				file
// @Produce			png,jpeg,gif
// @Param				userId			path string true "User ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router			/file/user/{userId}/profile [get]
func (c *FileController) GetUserProfile(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	path := fmt.Sprintf("/user/%s/profile", userId)
	url, err := url.JoinPath(c.filerUrl, path)
	if err != nil {
		return errs.New(errs.ErrCreateUrlPath, "invalid url", err)
	}
	return proxy.Forward(url)(ctx)
}

// GetWorkspaceProfile godoc
//
// @Summary 		Get a workspace profile image
// @Description Get a workspace profile image from internal file system by proxy
// @Tags				file
// @Produce			png,jpeg,gif
// @Param				workspaceId			path number true "Workspace ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router			/file/workspaces/{workspaceId}/profile [get]
func (c *FileController) GetWorkspaceProfile(ctx *fiber.Ctx) error {
	var pl payload.WorkspacePath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	path := fmt.Sprintf("/workspaces/%d/profile", pl.WorkspaceId)
	url, err := url.JoinPath(c.filerUrl, path)
	if err != nil {
		return errs.New(errs.ErrCreateUrlPath, "invalid url", err)
	}
	return proxy.Forward(url)(ctx)
}

// GetAssignmentDetail godoc
//
// @Summary 		Get a workspace detail markdown
// @Description Get a workspace detail markdown from internal file system by proxy
// @Tags				file
// @Produce			png,jpeg,gif
// @Param				workspaceId			path number true "Workspace ID"
// @Param				assignmentId		path number true "Assignment ID"
// @Param				subPath					path string false "Sub path"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router			/file/workspaces/{workspaceId}/assignments/{assignmentId}/detail/{subPath} [get]
func (c *FileController) GetAssignmentDetail(ctx *fiber.Ctx) error {
	var pl payload.AssignmentPath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	subPath := ctx.Params("*")
	if subPath == "" {
		return ctx.Status(fiber.StatusNotFound).Send(nil)
	}

	path := fmt.Sprintf(
		"/workspaces/%d/assignments/%d/detail/%s",
		pl.WorkspaceId, pl.AssignmentId, subPath,
	)
	url, err := url.JoinPath(c.filerUrl, path)
	if err != nil {
		return errs.New(errs.ErrCreateUrlPath, "invalid url", err)
	}
	return proxy.Forward(url)(ctx)
}

func (c *FileController) GetAssignmentTestcase(ctx *fiber.Ctx) error {
	var pl payload.AssignmentPath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	testcaseFile := ctx.Params("testcaseFile")
	path := fmt.Sprintf(
		"/workspaces/%d/assignments/%d/testcase/%s",
		pl.WorkspaceId, pl.AssignmentId, testcaseFile,
	)
	url, err := url.JoinPath(c.filerUrl, path)
	if err != nil {
		return errs.New(errs.ErrCreateUrlPath, "invalid url", err)
	}
	return proxy.Forward(url)(ctx)
}

func (c *FileController) GetSubmission(ctx *fiber.Ctx) error {
	var pl payload.SubmissionPath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)
	submittedUserId := ctx.Params("userId")

	path := fmt.Sprintf(
		"/workspaces/%d/assignments/%d/submissions/%s/%d",
		pl.WorkspaceId, pl.AssignmentId, submittedUserId, pl.SubmissionId,
	)

	userRole, err := c.WorkspaceUsecase.GetRole(user.Id, pl.WorkspaceId)
	if err != nil {
		return err
	} else if *userRole == domain.MemberRole && user.Id != submittedUserId {
		return errs.New(errs.ErrFilePerm, "no permission to get file not own", err)
	}

	url, err := url.JoinPath(c.filerUrl, path)
	if err != nil {
		return errs.New(errs.ErrCreateUrlPath, "invalid url", err)
	}
	return proxy.Forward(url)(ctx)
}
