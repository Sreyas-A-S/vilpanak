package main

import (
	"log"
	"net/http"

	"vilpanak/internal/config"
	"vilpanak/internal/database"
	"vilpanak/internal/handler"
	"vilpanak/internal/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize session manager & CSRF token storage
	middleware.InitSession(cfg)

	// Attempt connection to PostgreSQL (fails gracefully if DB not yet booted)
	_, _ = database.Connect(cfg)

	r := gin.Default()

	// Global Security & Performance Middlewares
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimiter(rate.Limit(20), 50)) // 20 requests/sec with burst of 50
	r.Use(middleware.CSRFMiddleware())

	// Static Assets
	r.Static("/static", "./static")

	// Handlers
	webHandler := handler.NewWebHandler()
	adminHandler := handler.NewAdminHandler()

	// Public Routes
	r.GET("/", webHandler.Home)
	r.GET("/properties", webHandler.Properties)
	r.GET("/properties/:slug", webHandler.PropertyDetail)
	r.POST("/properties/:slug/inquiry", webHandler.PostInquiry)

	// Admin Routes (Auth Protected)
	r.GET("/admin/login", adminHandler.ShowLogin)
	r.POST("/admin/login", adminHandler.PostLogin)
	r.POST("/admin/logout", adminHandler.Logout)

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.RequireAuth())
	{
		adminGroup.GET("/dashboard", adminHandler.Dashboard)
	}

	serverAddr := ":" + cfg.Port
	log.Printf("Starting Vilpanak Real Estate server at http://localhost%s (Env: %s)", serverAddr, cfg.Env)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
