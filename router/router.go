package router

import (
	"Web/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	var r = gin.Default()
	r.POST("/register", controller.RegisterHandler) //注册用户为什么用POST
	user := r.Group("/user")
	{
		user.GET("/token", controller.LoginHandler)
		user.GET("/token/refresh", controller.RefreshTokenHandler) //刷新refresh token
		user.PUT("/password", controller.PasswordChangeHandler)
		user.PUT("/info", controller.InfoChangeHandler)
		user.GET("/info/:user_id", controller.GetUserInfoHandler)
	}
	book := r.Group("/book")
	{
		book.GET("/list", controller.GetBookListHandler)
		book.GET("/search", controller.SearchBookHandler)
		book.PUT("/star", controller.BookStarHandler)
		book.GET("/label", controller.GetLabledBookHandler)
	}
	comment := r.Group("comment")
	{
		comment.GET("/:book_id", controller.GetBookCommentHandler)
		comment.POST("/:book_id", controller.WriteBookCommentHandler)
		comment.DELETE("/:comment_id", controller.DeleteBookCommentHandler)
		comment.PUT("/:comment_id", controller.ChangeBookCommentHandler)
	}
	operate := r.Group("/operate")
	{
		operate.PUT("/praise", controller.ThumbsUpHandler)
		operate.GET("/collect/list", controller.GetUserCollectListHandler)
		operate.PUT("/focus", controller.UserFocusHandler)
	}
	r.Run()
}
