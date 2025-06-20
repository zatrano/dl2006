package requests

import (
	"github.com/gofiber/fiber/v2"
)

type SocialMediaRequest struct {
	Name     string `form:"name" validate:"required,min=2"`
	Icon     string `form:"icon" validate:"required"`
	IsActive string `form:"is_active"`
}

func ValidateSocialMediaRequest(c *fiber.Ctx) error {
	var req SocialMediaRequest
	errorMessages := map[string]string{
		"Name_required": "Sosyal medya adı zorunludur",
		"Name_min":      "Sosyal medya adı en az 2 karakter olmalıdır",
		"Icon_required": "İkon zorunludur",
	}
	if err := validateRequest(c, &req, errorMessages, "/dashboard/social-media/create"); err != nil {
		return err
	}
	c.Locals("socialMediaRequest", req)
	return c.Next()
}
