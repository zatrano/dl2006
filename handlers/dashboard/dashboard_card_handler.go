package handlers

import (
	"net/http"

	"davet.link/configs/logconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/renderer"
	"davet.link/requests"
	"davet.link/services"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

type DashboardCardHandler struct {
	cardService        services.ICardService
	userService        services.IUserService
	bankService        services.IBankService
	socialMediaService services.ISocialMediaService
}

func NewDashboardCardHandler() *DashboardCardHandler {
	return &DashboardCardHandler{
		cardService:        services.NewCardService(),
		userService:        services.NewUserService(),
		bankService:        services.NewBankService(),
		socialMediaService: services.NewSocialMediaService(),
	}
}

func (h *DashboardCardHandler) ListCards(c *fiber.Ctx) error {
	var params queryparams.ListParams
	if err := c.QueryParser(&params); err != nil {
		logconfig.Log.Warn("Kullanıcı listesi: Query parametreleri parse edilemedi, varsayılanlar kullanılıyor.", zap.Error(err))
		params = queryparams.DefaultListParams()
	}

	if params.Page <= 0 {
		params.Page = queryparams.DefaultPage
	}
	if params.PerPage <= 0 || params.PerPage > queryparams.MaxPerPage {
		params.PerPage = queryparams.DefaultPerPage
	}
	if params.SortBy == "" {
		params.SortBy = queryparams.DefaultSortBy
	}
	if params.OrderBy == "" {
		params.OrderBy = queryparams.DefaultOrderBy
	}

	paginatedResult, dbErr := h.cardService.GetAllCards(params)

	renderData := fiber.Map{
		"Title":  "Kartlar",
		"Result": paginatedResult,
		"Params": params,
	}
	if dbErr != nil {
		logconfig.Log.Error("Kart listesi DB Hatası", zap.Error(dbErr))
		renderData[renderer.FlashErrorKeyView] = "Kartlar getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.User{},
			Meta: queryparams.PaginationMeta{
				CurrentPage: params.Page, PerPage: params.PerPage,
			},
		}
	}
	return renderer.Render(c, "dashboard/cards/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *DashboardCardHandler) ShowCreateCard(c *fiber.Ctx) error {
	usersResult, _ := h.userService.GetAllUsers(queryparams.ListParams{PerPage: 1000})
	banksResult, _ := h.bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
	socialMediasResult, _ := h.socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
	return renderer.Render(c, "dashboard/cards/create", "layouts/dashboard", fiber.Map{
		"Title":        "Yeni Kart Oluştur",
		"Users":        usersResult.Data,
		"Banks":        banksResult.Data,
		"SocialMedias": socialMediasResult.Data,
	}, http.StatusOK)
}

func (h *DashboardCardHandler) CreateCard(c *fiber.Ctx) error {
	if err := requests.ValidateCardRequest(c); err != nil {
		return err
	}
	req := c.Locals("cardRequest").(requests.CardRequest)
	card := &models.Card{
		Name:      req.Name,
		Slug:      req.Slug,
		UserID:    req.UserID,
		Photo:     req.Photo,
		Telephone: req.Telephone,
		Email:     req.Email,
		Location:  req.Location,
		Website:   req.Website,
		IsActive:  req.IsActive == "true",
	}
	for _, bankID := range req.BankIDs {
		card.CardBanks = append(card.CardBanks, models.CardBank{CardID: card.ID, BankID: bankID, IBAN: ""})
	}
	for _, smID := range req.SocialMediaIDs {
		card.CardSocialMedia = append(card.CardSocialMedia, models.CardSocialMedia{CardID: card.ID, SocialMediaID: smID, URL: ""})
	}
	if err := h.cardService.CreateCard(c.UserContext(), card); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Kart oluşturulamadı")
	}
	return c.Redirect("/dashboard/cards", http.StatusFound)
}

func (h *DashboardCardHandler) ShowUpdateCard(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	card, err := h.cardService.GetCardByID(c.UserContext(), uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Kart bulunamadı")
	}
	usersResult, _ := h.userService.GetAllUsers(queryparams.ListParams{PerPage: 1000})
	banksResult, _ := h.bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
	socialMediasResult, _ := h.socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
	return renderer.Render(c, "dashboard/cards/update", "layouts/dashboard", fiber.Map{
		"Title":        "Kartı Düzenle",
		"Card":         card,
		"Users":        usersResult.Data,
		"Banks":        banksResult.Data,
		"SocialMedias": socialMediasResult.Data,
	}, http.StatusOK)
}

func (h *DashboardCardHandler) UpdateCard(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := requests.ValidateCardRequest(c); err != nil {
		return err
	}
	req := c.Locals("cardRequest").(requests.CardRequest)
	card := &models.Card{
		Name:      req.Name,
		Slug:      req.Slug,
		UserID:    req.UserID,
		Photo:     req.Photo,
		Telephone: req.Telephone,
		Email:     req.Email,
		Location:  req.Location,
		Website:   req.Website,
		IsActive:  req.IsActive == "true",
	}
	for _, bankID := range req.BankIDs {
		card.CardBanks = append(card.CardBanks, models.CardBank{CardID: uint(id), BankID: bankID, IBAN: ""})
	}
	for _, smID := range req.SocialMediaIDs {
		card.CardSocialMedia = append(card.CardSocialMedia, models.CardSocialMedia{CardID: uint(id), SocialMediaID: smID, URL: ""})
	}
	if err := h.cardService.UpdateCard(c.UserContext(), uint(id), card); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Kart güncellenemedi")
	}
	return c.Redirect("/dashboard/cards", http.StatusFound)
}

func (h *DashboardCardHandler) DeleteCard(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := h.cardService.DeleteCard(c.UserContext(), uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Kart silinemedi")
	}
	return c.Redirect("/dashboard/cards", http.StatusFound)
}
