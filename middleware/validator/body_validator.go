package validator

import (
	"booking/internal/constant"
	"booking/internal/vo"

	"github.com/gofiber/fiber/v2"
)

func BodyValidator[T any]() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := new(T)
		c.BodyParser(&body)
		customValidator := GetCustomValidatorInstance()
		if errs := customValidator.Validate(body); errs != nil {
			res := vo.Response{}
			res.SetErrorMessage(constant.ERROR_MESSAGE_INVALID_PARAMETERS)
			for _, err := range errs {
				res.AppendErrors(err)
			}
			return c.Status(fiber.StatusBadRequest).JSON(res)
		}
		return c.Next()
	}
}
