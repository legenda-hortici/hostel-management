package handlers

import (
	"hostel-management/internal/services"
	"hostel-management/pkg/helpers"
	"hostel-management/pkg/middlewares"
	handlers "hostel-management/pkg/validation"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	userService services.UserService
}

func NewProfileHandler(userService services.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

func (h *ProfileHandler) Profile(c *gin.Context) {
	const op = "handlers.ProfileHandler.ProfileHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	email, err := handlers.ValidateUserByEmail(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		log.Printf("Access denied: %v", err)
		return
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить пользователя")
		log.Printf("Failed to get user: %v: %v", err, op)
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":    "profile",
		"Role":    role,
		"User":    user,
		"Flashes": flashes,
	})
}

func (h *ProfileHandler) UpdateProfileHandler(c *gin.Context) {
	const op = "handlers.ProfileHandler.UpdateProfileHandler"

	c.Request.ParseMultipartForm(10 << 20) // Ограничение на 10 MB

	email, err := handlers.ValidateUserByEmail(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		log.Printf("Access denied: %v", err)
		return
	}

	if c.Request.Method != "POST" {
		log.Printf(" Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	name := c.PostForm("username")
	surname := c.PostForm("surname")
	password := c.PostForm("password")

	avatarFile, _, err := c.Request.FormFile("avatar")
	var avatarPath string
	if err == nil && avatarFile != nil {
		avatarPath, err = helpers.SaveAvatar(avatarFile)
		if err != nil {
			log.Printf("Failed to save avatar: %v", err)
			middlewares.HandleError(c, 500, "Ошибка: не удалось сохранить аватар")
			return
		}
	}

	if avatarPath == "" {
		avatarPath = "Не указана"
	}

	err = h.userService.UpdateUserByEmail(email, name, surname, password, avatarPath)
	if err != nil {
		log.Printf("Failed to update user: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось обновить данные пользователя")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Данные успешно обновлены")
	session.Save()

	c.Redirect(303, "/profile")
}
