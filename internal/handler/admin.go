package handler

import (
	"html/template"
	"net/http"

	"vilpanak/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) ShowLogin(c *gin.Context) {
	csrfToken, _ := c.Get("csrf_token")
	tmpl, err := template.ParseFiles("templates/admin/login.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Template error: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(c.Writer, gin.H{
		"CSRFToken": csrfToken,
		"Error":     c.Query("error"),
	})
}

func (h *AdminHandler) PostLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Demo/Default super admin verification (in production check hash against DB)
	if email == "admin@vilpanak.com" && password == "admin123" {
		session, _ := middleware.Store.Get(c.Request, "admin_session")
		session.Values["user_id"] = int64(1)
		session.Values["user_name"] = "Super Administrator"
		session.Values["user_role"] = "super_admin"
		_ = session.Save(c.Request, c.Writer)

		c.Redirect(http.StatusSeeOther, "/admin/dashboard")
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/login?error=Invalid+email+or+password")
}

func (h *AdminHandler) Logout(c *gin.Context) {
	session, _ := middleware.Store.Get(c.Request, "admin_session")
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1
	_ = session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusSeeOther, "/admin/login")
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	csrfToken, _ := c.Get("csrf_token")
	userName, _ := c.Get("user_name")

	stats := gin.H{
		"TotalProperties": 12,
		"ForSaleCount":    8,
		"ForRentCount":    4,
		"NewLeadsCount":   3,
	}

	tmpl, err := template.ParseFiles("templates/layouts/admin_base.html", "templates/admin/dashboard.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Template error: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(c.Writer, gin.H{
		"CSRFToken": csrfToken,
		"UserName":  userName,
		"Stats":     stats,
	})
}
