package handlers

import (
	"net/http"

	"davet.link/configs/logconfig"
	"davet.link/pkg/renderer"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DashboardHomeHandler struct {
	userService services.IUserService
}

func NewDashboardHomeHandler() *DashboardHomeHandler {
	return &DashboardHomeHandler{
		userService: services.NewUserService(),
	}
}

func (h *DashboardHomeHandler) HomePage(c *fiber.Ctx) error {
	userCount, userErr := h.userService.GetUserCount()
	if userErr != nil {
		logconfig.Log.Error("Anasayfa: Kullanıcı sayısı alınamadı", zap.Error(userErr))
		userCount = 0
	}

	mapData := fiber.Map{
		"Title":     "Dashboard",
		"UserCount": userCount,
	}
	return renderer.Render(c, "dashboard/home/home", "layouts/dashboard", mapData, http.StatusOK)
}
