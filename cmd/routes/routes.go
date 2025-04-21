package routes

import (
	"hostel-management/internal/dependencies"
	"hostel-management/pkg/auth"
	"hostel-management/storage/redis"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, redisCache *redis.RedisCache) error {
	// Настраиваем шаблоны
	r.LoadHTMLGlob("web/templates/*.html")
	r.Static("/static", "./web/static")

	deps := dependencies.BuildDependencies(redisCache)

	// Публичные маршруты
	public := r.Group("/")
	RegisterPublicRoutes(public, deps)

	// Защищенные маршруты
	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	RegisterProtectedRoutes(protected, deps)

	// Пользователи
	user := r.Group("/")
	user.Use(auth.AuthMiddleware())
	RegisterUserRoutes(user, deps)

	// Коменданты
	headman := r.Group("/headman")
	headman.Use(auth.HeadmanMiddleware())
	RegisterHeadmanRoutes(headman, deps)

	// TODO: Сделать добавление админов и комендантов нормальным
	// Админы
	admin := r.Group("/admin")
	admin.Use(auth.AdminMiddleware())
	RegisterAdminRoutes(admin, deps)

	return nil
}
