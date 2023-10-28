package middleware

import (
	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/gofiber/fiber/v2"
)

func NewWorkspaceMiddleware(
	workspaceUsecase domain.WorkspaceUsecase,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Params("workspaceId") == "" {
			return ctx.Next()
		}

		workspaceId, err := ctx.ParamsInt("workspaceId")
		if err != nil {
			return errs.New(errs.ErrPayloadValidator, "param workspaceId is invalid")
		}

		user := GetUserFromCtx(ctx)
		ok, err := workspaceUsecase.HasUser(user.Id, workspaceId)
		if !ok {
			return errs.New(errs.ErrWorkspaceNoPerm, "cannot access workspace id %d", workspaceId)
		} else if err != nil {
			return err
		}

		var assignmentId int
		if ctx.Params("assignmentId") != "" {
			assignmentId, err = ctx.ParamsInt("assignmentId")
			if err != nil {
				return errs.New(errs.ErrPayloadValidator, "param assignmentId is invalid", err)
			}

			ok, err := workspaceUsecase.HasAssignment(assignmentId, workspaceId)
			if !ok {
				return errs.New(errs.ErrWorkspaceNoPerm, "cannot access assignment id %d", assignmentId)
			} else if err != nil {
				return err
			}
		}

		ctx.Locals("workspaceId", workspaceId)
		ctx.Locals("assignmentId", assignmentId)

		return ctx.Next()
	}
}

func NewWorkspaceRoleMiddleware(workspaceUsecase domain.WorkspaceUsecase) func(expectedRoles ...domain.WorkspaceRole) func(*fiber.Ctx) error {
	return func(expectedRoles ...domain.WorkspaceRole) fiber.Handler {
		return func(ctx *fiber.Ctx) error {
			user := GetUserFromCtx(ctx)
			workspaceId := GetWorkspaceIdFromCtx(ctx)

			role, err := workspaceUsecase.GetRole(user.Id, workspaceId)
			if err != nil {
				return errs.New(errs.ErrGetWorkspaceRole, "cannot get workspace role to list submission", err)
			}
			if role == nil {
				return errs.New(errs.ErrWorkspaceNoPerm, "user %s does not have role and permission to access workspace id %d", user.Id, workspaceId)
			}

			for _, expectedRole := range expectedRoles {
				if *role == expectedRole {
					return ctx.Next()
				}
			}

			return errs.New(errs.ErrWorkspaceNoPerm, "user %s does not have permission to access workspace id %d", user.Id, workspaceId)
		}
	}
}

func GetWorkspaceIdFromCtx(ctx *fiber.Ctx) int {
	return ctx.Locals("workspaceId").(int)
}

func GetAssignmentIdFromCtx(ctx *fiber.Ctx) int {
	return ctx.Locals("assignmentId").(int)
}
