package validator_test

import (
	"booking/middleware/validator"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

func TestQueryValidator_ResponseError(t *testing.T) {
	type QueryParams struct {
		Text    string `validate:"required" query:"text"`
		Number  int    `validate:"required" query:"number"`
		Boolean bool   `validate:"required" query:"boolean"`
	}

	app := fiber.New()

	app.Use(fiberi18n.New(&fiberi18n.Config{
		RootPath:        "../../localize",
		AcceptLanguages: []language.Tag{language.Thai, language.English},
		DefaultLanguage: language.English,
	}))

	dummyText := "Dummy text"

	app.Post("/",
		validator.QueryValidator[QueryParams](),
		func(c *fiber.Ctx) error {
			return c.SendString(dummyText)
		})

	req := httptest.NewRequest(http.MethodPost, "/", nil)

	resp, _ := app.Test(req)

	t.Run("the statusCode code should be 400", func(t *testing.T) {
		got := resp.StatusCode
		expected := fiber.StatusBadRequest
		if got != expected {
			t.Errorf("QueryValidator() statusCode code = %v, expected %v", got, expected)
		}
	})

	t.Run("should response the errors", func(t *testing.T) {
		body, _ := io.ReadAll(resp.Body)
		got := string(body)
		expected := string(`{"errors":[{"failedField":"Text","tag":"required","value":null},{"failedField":"Number","tag":"required","value":null},{"failedField":"Boolean","tag":"required","value":null}],"message":"Invalid Parameters","statusCode":400,"success":false}`)
		if got != expected {
			t.Errorf("QueryValidator() response = %v, expected %v", got, expected)
		}
	})
}

func TestQueryValidator_ResponseData(t *testing.T) {
	type QueryParams struct {
		Text    string `validate:"required" query:"text"`
		Number  int    `validate:"required" query:"number"`
		Boolean bool   `validate:"required" query:"boolean"`
	}

	app := fiber.New()

	dummyText := "Dummy text"

	app.Post("/",
		validator.QueryValidator[QueryParams](),
		func(c *fiber.Ctx) error {
			c.Status(fiber.StatusOK)
			return c.SendString(dummyText)
		})

	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.URL.RawQuery = url.Values{
		"text":    {"Text"},
		"number":  {"25"},
		"boolean": {"true"},
	}.Encode()

	resp, _ := app.Test(req)

	t.Run("should pass the query validator middleware and get the response from the handler", func(t *testing.T) {
		body, _ := io.ReadAll(resp.Body)
		got := string(body)
		expected := dummyText
		if got != expected {
			t.Errorf("QueryValidator() response = %v, expected %v", got, expected)
		}
	})
}
