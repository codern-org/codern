package controller

import (
	"fmt"
	"net/url"

	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/platform/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type FileController struct {
	filerUrl string
}

func NewFileController(cfg *config.Config) *FileController {
	return &FileController{filerUrl: cfg.Client.SeaweedFs.FilerUrls[0]}
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
		return err
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
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)

	path := fmt.Sprintf("/workspaces/%d/profile", workspaceId)
	url, err := url.JoinPath(c.filerUrl, path)
	if err != nil {
		return err
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
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router			/file/workspaces/{workspaceId}/assignments/{assignmentId}/detail [get]
func (c *FileController) GetAssignmentDetail(ctx *fiber.Ctx) error {
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)
	assignmentId := middleware.GetAssignmentIdFromCtx(ctx)

	path := fmt.Sprintf("/workspaces/%d/assignments/%d/problem.md", workspaceId, assignmentId)
	url, err := url.JoinPath(c.filerUrl, path)
	if err != nil {
		return err
	}
	return proxy.Forward(url)(ctx)
}
