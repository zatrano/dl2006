package handlers

import (
	"net/http"
	"strings"

	"davet.link/configs/logconfig"
	"davet.link/models"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/renderer"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DashboardInvitationCategoryHandler struct {
	categoryService services.IInvitationCategoryService
}

func NewDashboardInvitationCategoryHandler() *DashboardInvitationCategoryHandler {
	svc := services.NewInvitationCategoryService()
	return &DashboardInvitationCategoryHandler{categoryService: svc}
}

func (h *DashboardInvitationCategoryHandler) ListCategories(c *fiber.Ctx) error {
	var params queryparams.ListParams
	if err := c.QueryParser(&params); err != nil {
		logconfig.Log.Warn("Kategori listesi: Query parametreleri parse edilemedi, varsayılanlar kullanılıyor.", zap.Error(err))
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

	paginatedResult, dbErr := h.categoryService.GetAllCategories(params)

	renderData := fiber.Map{
		"Title":  "Davet Kategorileri",
		"Result": paginatedResult,
		"Params": params,
	}
	if dbErr != nil {
		logconfig.Log.Error("Kategori listesi DB Hatası", zap.Error(dbErr))
		renderData[renderer.FlashErrorKeyView] = "Kategoriler getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.InvitationCategory{},
			Meta: queryparams.PaginationMeta{
				CurrentPage: params.Page, PerPage: params.PerPage,
			},
		}
	}
	return renderer.Render(c, "dashboard/invitation-categories/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *DashboardInvitationCategoryHandler) ShowCreateCategory(c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/invitation-categories/create", "layouts/dashboard", fiber.Map{
		"Title": "Yeni Kategori Ekle",
	})
}

func (h *DashboardInvitationCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req struct {
		Name     string `form:"name"`
		Icon     string `form:"icon"`
		Template string `form:"template"`
		IsActive string `form:"is_active"`
	}
	_ = c.BodyParser(&req)

	if req.Name == "" || req.Icon == "" || req.Template == "" {
		return renderCategoryFormError("Yeni Kategori Ekle", req, "Ad, İkon ve Şablon alanları zorunludur.", c)
	}

	category := &models.InvitationCategory{
		Name:     req.Name,
		Icon:     req.Icon,
		Template: req.Template,
		IsActive: req.IsActive == "true",
	}

	if err := h.categoryService.CreateCategory(c.UserContext(), category); err != nil {
		return renderCategoryFormError("Yeni Kategori Ekle", req, "Kategori oluşturulamadı: "+err.Error(), c)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kategori başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/invitation-categories", fiber.StatusFound)
}

func (h *DashboardInvitationCategoryHandler) ShowUpdateCategory(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kategori bulunamadı.")
		return c.Redirect("/dashboard/invitation-categories", fiber.StatusSeeOther)
	}
	return renderer.Render(c, "dashboard/invitation-categories/update", "layouts/dashboard", fiber.Map{
		"Title":    "Kategori Düzenle",
		"Category": category,
	})
}

func (h *DashboardInvitationCategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	categoryID := uint(id)

	var req struct {
		Name     string `form:"name"`
		Icon     string `form:"icon"`
		Template string `form:"template"`
		IsActive string `form:"is_active"`
	}
	_ = c.BodyParser(&req)

	if req.Name == "" || req.Icon == "" || req.Template == "" {
		category, _ := h.categoryService.GetCategoryByID(categoryID)
		return renderer.Render(c, "dashboard/invitation-categories/update", "layouts/dashboard", fiber.Map{
			"Title":                    "Kategori Düzenle",
			renderer.FlashErrorKeyView: "Zorunlu alanlar eksik.",
			renderer.FormDataKey:       req,
			"Category":                 category,
		}, http.StatusBadRequest)
	}

	categoryData := &models.InvitationCategory{
		Name:     req.Name,
		Icon:     req.Icon,
		Template: req.Template,
		IsActive: req.IsActive == "true",
	}

	// Güncelleyen kullanıcı kimliği context'ten alınmalı, örnek: contextUserIDKey
	updatedBy := uint(0)
	if v := c.Locals("user_id"); v != nil {
		if uid, ok := v.(uint); ok {
			updatedBy = uid
		}
	}

	if err := h.categoryService.UpdateCategory(c.UserContext(), categoryID, categoryData, updatedBy); err != nil {
		category, _ := h.categoryService.GetCategoryByID(categoryID)
		return renderer.Render(c, "dashboard/invitation-categories/update", "layouts/dashboard", fiber.Map{
			"Title":                    "Kategori Düzenle",
			renderer.FlashErrorKeyView: "Güncelleme hatası: " + err.Error(),
			renderer.FormDataKey:       req,
			"Category":                 category,
		}, http.StatusInternalServerError)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kategori başarıyla güncellendi.")
	return c.Redirect("/dashboard/invitation-categories", fiber.StatusFound)
}

func (h *DashboardInvitationCategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	categoryID := uint(id)

	if err := h.categoryService.DeleteCategory(c.UserContext(), categoryID); err != nil {
		errMsg := "Kategori silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/dashboard/invitation-categories", fiber.StatusSeeOther)
	}

	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Kategori başarıyla silindi."})
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kategori başarıyla silindi.")
	return c.Redirect("/dashboard/invitation-categories", fiber.StatusFound)
}

func renderCategoryFormError(title string, req any, message string, c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/invitation-categories/create", "layouts/dashboard", fiber.Map{
		"Title":                    title,
		renderer.FlashErrorKeyView: message,
		renderer.FormDataKey:       req,
	}, http.StatusBadRequest)
}
