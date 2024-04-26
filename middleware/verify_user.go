package middleware

import (
	"booking/internal/cache"
	"booking/internal/constant"
	"booking/internal/util"
	"booking/internal/vo"

	"github.com/gofiber/fiber/v2"
)

func VerifyUser(cache cache.ICache, util util.IUtility) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		accessTokenString, err := util.ExtractJWTTokenStringFromHeaders(c.GetReqHeaders())
		res := vo.Response{}
		if err != nil {
			res.SetErrorMessage(constant.ERROR_MESSAGE_UNAUTHORIZED)
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		accessToken, err := util.JWTParseAccessToken(accessTokenString)

		if err != nil {
			res.SetErrorMessage(constant.ERROR_MESSAGE_UNAUTHORIZED)
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		accessTokenStringFromCache, err := cache.GetAccessTokenString(accessToken.Username)
		if err != nil || accessTokenStringFromCache != accessTokenString {
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		c.Locals(constant.LOCAL_KEY_ACCESS_TOKEN, accessToken)
		return c.Next()
	}

}
