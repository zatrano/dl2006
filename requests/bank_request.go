package requests

import (
	"github.com/gofiber/fiber/v2"
)

type BankRequest struct {
	Name     string `form:"name" validate:"required,min=2"`
	Iban     string `form:"iban" validate:"required,min=10"`
	IsActive string `form:"is_active"`
}

func ValidateBankRequest(c *fiber.Ctx) error {
	var req BankRequest
	errorMessages := map[string]string{
		"Name_required": "Banka adı zorunludur",
		"Name_min":      "Banka adı en az 2 karakter olmalıdır",
		"Iban_required": "IBAN zorunludur",
		"Iban_min":      "IBAN en az 10 karakter olmalıdır",
	}
	if err := validateRequest(c, &req, errorMessages, "/dashboard/banks/create"); err != nil {
		return err
	}
	c.Locals("bankRequest", req)
	return c.Next()
}
