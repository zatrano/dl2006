package requests

import (
	"github.com/gofiber/fiber/v2"
)

type CardRequest struct {
	Name           string `form:"name" validate:"required,min=2"`
	Slug           string `form:"slug" validate:"required,min=2"`
	UserID         uint   `form:"user_id" validate:"required,gt=0"`
	Photo          string `form:"photo"`
	Telephone      string `form:"telephone"`
	Email          string `form:"email"`
	Location       string `form:"location"`
	Website        string `form:"website"`
	IsActive       string `form:"is_active"`
	BankIDs        []uint `form:"bank_ids[]"`
	SocialMediaIDs []uint `form:"social_media_ids[]"`
}

func ValidateCardRequest(c *fiber.Ctx) error {
	var req CardRequest
	errorMessages := map[string]string{
		"Name_required":   "Kart adı zorunludur",
		"Name_min":        "Kart adı en az 2 karakter olmalıdır",
		"Slug_required":   "Slug zorunludur",
		"Slug_min":        "Slug en az 2 karakter olmalıdır",
		"UserID_required": "Kullanıcı seçimi zorunludur",
		"UserID_gt":       "Kullanıcı seçimi zorunludur",
	}
	if err := validateRequest(c, &req, errorMessages, "/dashboard/cards/create"); err != nil {
		return err
	}
	c.Locals("cardRequest", req)
	return c.Next()
}
