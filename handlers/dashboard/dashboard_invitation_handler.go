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

type DashboardInvitationHandler struct {
	invitationService services.IInvitationService
	userService       services.IUserService
	categoryService   services.IInvitationCategoryService
}

func NewDashboardInvitationHandler() *DashboardInvitationHandler {
	return &DashboardInvitationHandler{
		invitationService: services.NewInvitationService(),
		userService:       services.NewUserService(),
		categoryService:   services.NewInvitationCategoryService(),
	}
}

func (h *DashboardInvitationHandler) ListInvitations(c *fiber.Ctx) error {
	var params queryparams.ListParams
	if err := c.QueryParser(&params); err != nil {
		logconfig.Log.Warn("Davetiyeler: Query parametreleri parse edilemedi, varsayılanlar kullanılıyor.", zap.Error(err))
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
	paginatedResult, dbErr := h.invitationService.GetAllInvitations(params)
	renderData := fiber.Map{
		"Title":  "Davetiyeler",
		"Result": paginatedResult,
		"Params": params,
	}
	if dbErr != nil {
		logconfig.Log.Error("Davetiyeler listesi DB Hatası", zap.Error(dbErr))
		renderData[renderer.FlashErrorKeyView] = "Davetiyeler getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.Invitation{},
			Meta: queryparams.PaginationMeta{
				CurrentPage: params.Page, PerPage: params.PerPage,
			},
		}
	}
	return renderer.Render(c, "dashboard/invitations/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *DashboardInvitationHandler) ShowCreateInvitation(c *fiber.Ctx) error {
	usersResult, _ := h.userService.GetAllUsers(queryparams.ListParams{PerPage: 1000})
	categoriesResult, _ := h.categoryService.GetAllCategories(queryparams.ListParams{PerPage: 1000})
	return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
		"Title":      "Yeni Davetiye Oluştur",
		"Users":      usersResult.Data,
		"Categories": categoriesResult.Data,
	}, http.StatusOK)
}

func (h *DashboardInvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	if err := requests.ValidateInvitationRequest(c); err != nil {
		return err
	}
	req := c.Locals("invitationRequest").(requests.InvitationRequest)
	invitation := &models.Invitation{
		InvitationKey: req.InvitationKey,
		UserID:        req.UserID,
		CategoryID:    req.CategoryID,
		Template:      req.Template,
		Type:          req.Type,
		Title:         req.Title,
		Image:         req.Image,
		Description:   req.Description,
		Venue:         req.Venue,
		Address:       req.Address,
		Location:      req.Location,
		Link:          req.Link,
		Telephone:     req.Telephone,
		Note:          req.Note,
		IsConfirmed:   req.IsConfirmed == "true",
		IsParticipant: req.IsParticipant == "true",
	}
	invitation.InvitationDetail = &models.InvitationDetail{
		Title:  req.DetailTitle,
		Person: req.DetailPerson,
	}
	// Katılımcı ekleme kaldırıldı, sadece website tarafından eklenir
	if err := h.invitationService.CreateInvitation(c.UserContext(), invitation); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Davetiye oluşturulamadı")
	}
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}

func (h *DashboardInvitationHandler) ShowUpdateInvitation(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	invitation, err := h.invitationService.GetInvitationByID(c.UserContext(), uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Davetiye bulunamadı")
	}
	usersResult, _ := h.userService.GetAllUsers(queryparams.ListParams{PerPage: 1000})
	categoriesResult, _ := h.categoryService.GetAllCategories(queryparams.ListParams{PerPage: 1000})
	return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
		"Title":      "Davetiye Düzenle",
		"Invitation": invitation,
		"Users":      usersResult.Data,
		"Categories": categoriesResult.Data,
	}, http.StatusOK)
}

func (h *DashboardInvitationHandler) UpdateInvitation(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := requests.ValidateInvitationRequest(c); err != nil {
		return err
	}
	req := c.Locals("invitationRequest").(requests.InvitationRequest)
	invitation := &models.Invitation{
		InvitationKey: req.InvitationKey,
		UserID:        req.UserID,
		CategoryID:    req.CategoryID,
		Template:      req.Template,
		Type:          req.Type,
		Title:         req.Title,
		Image:         req.Image,
		Description:   req.Description,
		Venue:         req.Venue,
		Address:       req.Address,
		Location:      req.Location,
		Link:          req.Link,
		Telephone:     req.Telephone,
		Note:          req.Note,
		IsConfirmed:   req.IsConfirmed == "true",
		IsParticipant: req.IsParticipant == "true",
	}
	invitation.InvitationDetail = &models.InvitationDetail{
		Title:  req.DetailTitle,
		Person: req.DetailPerson,
	}
	// Katılımcı ekleme kaldırıldı, sadece website tarafından eklenir
	if err := h.invitationService.UpdateInvitation(c.UserContext(), uint(id), invitation); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Davetiye güncellenemedi")
	}
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}

func (h *DashboardInvitationHandler) DeleteInvitation(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := h.invitationService.DeleteInvitation(c.UserContext(), uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Davetiye silinemedi")
	}
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}

// Katılımcı listesi (dashboard)
func (h *DashboardInvitationHandler) ListParticipants(c *fiber.Ctx) error {
	invID, _ := c.ParamsInt("id")
	participants, err := h.invitationService.GetParticipantsByInvitationID(uint(invID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Katılımcılar getirilemedi")
	}
	return renderer.Render(c, "dashboard/invitations/participants", "layouts/dashboard", fiber.Map{
		"Participants": participants,
	}, http.StatusOK)
}

// Katılımcı güncelleme (dashboard)
func (h *DashboardInvitationHandler) UpdateParticipant(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := requests.ValidateInvitationParticipantRequest(c); err != nil {
		return err
	}
	req := c.Locals("invitationParticipantRequest").(requests.InvitationParticipantRequest)
	participant := &models.InvitationParticipant{
		Title:       req.Title,
		PhoneNumber: req.PhoneNumber,
		GuestCount:  req.GuestCount,
	}
	if err := h.invitationService.UpdateParticipant(uint(id), participant); err != nil {
		return c.Status(500).SendString("Katılımcı güncellenemedi")
	}
	return c.Redirect("/dashboard/invitations/participants/"+c.Query("invitation_id"), 302)
}

// Katılımcı silme (dashboard)
func (h *DashboardInvitationHandler) DeleteParticipant(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	invID := c.Query("invitation_id")
	if err := h.invitationService.DeleteParticipant(uint(id)); err != nil {
		return c.Status(500).SendString("Katılımcı silinemedi")
	}
	return c.Redirect("/dashboard/invitations/participants/"+invID, 302)
}
