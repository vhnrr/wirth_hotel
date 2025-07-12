package main

import (
	"log"
	"wirth_hotel/config"
	"wirth_hotel/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	defer config.GetDB().Close()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	// Public routes
	r.GET("/", handlers.ShowHomePage)
	r.GET("/about", handlers.ShowAboutPage)
	r.GET("/review", handlers.ShowReviewPage)

	r.POST("/confirmation", handlers.HandleBookingConfirmation)

	r.GET("/booking/receipt/:bookingID", handlers.DownloadReceipt)

	r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.HandleLogin)

	r.GET("/register", handlers.ShowRegisterPage)
	r.POST("/register", handlers.HandleRegister)

	r.GET("/logout", handlers.HandleLogout)

	// Protected routes
	auth := r.Group("/")
	auth.Use(handlers.RequireLogin())
	{
		auth.GET("/profile", handlers.ShowProfilePage)
		auth.POST("/booking/confirm", handlers.HandleBookingConfirmation)
		// tambah route booking form dll sesuai kebutuhan
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
