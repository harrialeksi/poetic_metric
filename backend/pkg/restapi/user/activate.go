package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/poeticmetric/poeticmetric/backend/pkg/depot"
	dm "github.com/poeticmetric/poeticmetric/backend/pkg/restapi/middleware/depot"
	"github.com/poeticmetric/poeticmetric/backend/pkg/service/useraccesstoken"
	"github.com/poeticmetric/poeticmetric/backend/pkg/service/userself"
)

func activate(c *fiber.Ctx) error {
	dp := dm.Get(c)

	payload := &userself.ActivatePayload{}

	err := c.BodyParser(payload)
	if err != nil {
		return err
	}

	userAccessToken := &useraccesstoken.UserAccessToken{}
	userSelf := &userself.UserSelf{}

	err = dp.WithPostgresTransaction(func(dp2 *depot.Depot) error {
		var err2 error

		userSelf, err2 = userself.Activate(dp2, payload)
		if err2 != nil {
			return err2
		}

		userAccessToken, err2 = useraccesstoken.Create(dp2, userSelf.Id)
		if err2 != nil {
			return err2
		}

		return nil
	})
	if err != nil {
		return err
	}

	return c.JSON(map[string]any{
		"user":            userSelf,
		"userAccessToken": userAccessToken,
	})
}