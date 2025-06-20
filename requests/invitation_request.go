package requests

import (
	"github.com/gofiber/fiber/v2"
)

type InvitationRequest struct {
	InvitationKey     string   `form:"invitation_key" validate:"required,min=2"`
	UserID            uint     `form:"user_id" validate:"required,gt=0"`
	CategoryID        uint     `form:"category_id" validate:"required,gt=0"`
	Template          string   `form:"template"`
	Type              string   `form:"type"`
	Title             string   `form:"title" validate:"required,min=2"`
	Image             string   `form:"image"`
	Description       string   `form:"description"`
	Venue             string   `form:"venue"`
	Address           string   `form:"address"`
	Location          string   `form:"location"`
	Link              string   `form:"link"`
	Telephone         string   `form:"telephone"`
	Note              string   `form:"note"`
	Date              string   `form:"date"`
	Time              string   `form:"time"`
	IsConfirmed       string   `form:"is_confirmed"`
	IsParticipant     string   `form:"is_participant"`
	DetailTitle       string   `form:"detail_title"`
	DetailPerson      string   `form:"detail_person"`
	ParticipantTitles []string `form:"participant_titles[]"`
	ParticipantPhones []string `form:"participant_phones[]"`
	ParticipantCounts []int    `form:"participant_counts[]"`
}

func ValidateInvitationRequest(c *fiber.Ctx) error {
	var req InvitationRequest
	errorMessages := map[string]string{
		"InvitationKey_required": "Davetiye anahtarı zorunludur",
		"InvitationKey_min":      "Davetiye anahtarı en az 2 karakter olmalıdır",
		"UserID_required":        "Kullanıcı seçimi zorunludur",
		"UserID_gt":              "Kullanıcı seçimi zorunludur",
		"CategoryID_required":    "Kategori seçimi zorunludur",
		"CategoryID_gt":          "Kategori seçimi zorunludur",
		"Title_required":         "Başlık zorunludur",
		"Title_min":              "Başlık en az 2 karakter olmalıdır",
	}
	if err := validateRequest(c, &req, errorMessages, "/dashboard/invitations/create"); err != nil {
		return err
	}
	c.Locals("invitationRequest", req)
	return c.Next()
}
