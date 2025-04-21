package routes

import (
	"hostel-management/internal/dependencies"

	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(r *gin.RouterGroup, deps *dependencies.Dependencies) {
	r.GET("/support", deps.FaqHandler.SupportHandler)

	r.GET("/", deps.HomeHandler.HomeHandler)
	r.POST("/upload_banner", deps.HomeHandler.UploadBannerHandler)
	r.POST("/delete_banner", deps.HomeHandler.DeleteBannerHandler)

	r.GET("/news", deps.NewsHandler.News)
	r.GET("/news/:id", deps.NewsHandler.NewsInfoHandler)
	r.GET("/create_news", deps.NewsHandler.CreateNewsPageHandler)
	r.POST("/create_news/add", deps.NewsHandler.CreateNewsHandler)
	r.POST("/news/:id/delete", deps.NewsHandler.DeleteNewsHandler)

	r.GET("/notices", deps.NoticeHandler.Notices)
	r.GET("/notices/:id", deps.NoticeHandler.NoticeInfoHandler)
	r.GET("/create_notice", deps.NoticeHandler.CreateNoticePageHandler)
	r.POST("/create_notice/add", deps.NoticeHandler.CreateNoticeHandler)
	r.POST("/notices/:id/delete", deps.NoticeHandler.DeleteNoticeHandler)
}
