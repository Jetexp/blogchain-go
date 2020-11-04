package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zikwall/blogchain/src/models/user"
	"github.com/zikwall/blogchain/src/utils"
)

func UseBlogchainUserIdentity(ctx *fiber.Ctx) error {
	instance := &user.User{}

	claims := ctx.Locals("claims")

	if token, ok := claims.(*utils.TokenClaims); ok {
		u := user.NewUserModel()

		if i, err := u.FindById(token.UUID); err == nil {
			instance = i
		}
	}

	ctx.Locals("user", instance)

	return ctx.Next()
}