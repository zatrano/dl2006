package requests

import (
	"github.com/gofiber/fiber/v2"
)

type InvitationParticipantRequest struct {
	Title       string `form:"title" validate:"required,min=2"`
	PhoneNumber string `form:"phone_number" validate:"required,min=10"`
	GuestCount  int    `form:"guest_count" validate:"required,min=1"`
}

func ValidateInvitationParticipantRequest(c *fiber.Ctx) error {
	var req InvitationParticipantRequest
	errorMessages := map[string]string{
		"Title_required":       "Ad Soyad zorunludur",
		"Title_min":            "Ad Soyad en az 2 karakter olmalıdır",
		"PhoneNumber_required": "Telefon numarası zorunludur",
		"PhoneNumber_min":      "Telefon numarası en az 10 karakter olmalıdır",
		"GuestCount_required":  "Kişi sayısı zorunludur",
		"GuestCount_min":       "Kişi sayısı en az 1 olmalıdır",
	}
	if err := validateRequest(c, &req, errorMessages, ""); err != nil {
		return err
	}
	c.Locals("invitationParticipantRequest", req)
	return nil
}
