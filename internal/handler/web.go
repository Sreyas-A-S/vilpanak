package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"vilpanak/internal/models"

	"github.com/gin-gonic/gin"
)

type WebHandler struct{}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

func formatINRCurrency(val float64) string {
	if val >= 10000000 {
		cr := val / 10000000
		return fmt.Sprintf("%.2f Cr", cr)
	} else if val >= 100000 {
		lakh := val / 100000
		return fmt.Sprintf("%.2f L", lakh)
	}
	return fmt.Sprintf("%.0f", val)
}

func getSampleProperties() []models.Property {
	return []models.Property{
		{
			ID:             1,
			Title:          "Skyline Waterfront Penthouse",
			Slug:           "skyline-waterfront-penthouse",
			Description:    "Luxury 4 BHK duplex penthouse overlooking Marine Drive Kochi with private deck, soundproof glazing, and Italian marble flooring throughout.",
			ListingType:    models.ListingSale,
			Price:          17500000,
			PriceFrequency: "total",
			Bedrooms:       4,
			Bathrooms:      4,
			AreaSqFt:       3800,
			ParkingSpaces:  2,
			Address:        "Marine Drive, High Court Junction",
			City:           "Ernakulam",
			State:          "Kerala",
			PostalCode:     "682031",
			CoverImage:     "https://images.unsplash.com/photo-1600596542815-ffad4c1539a9?auto=format&fit=crop&w=800&q=80",
		},
		{
			ID:             2,
			Title:          "Greenfield Traditional Luxury Villa",
			Slug:           "greenfield-traditional-luxury-villa",
			Description:    "Kerala architectural style 5 BHK standalone villa near Technopark with teakwood woodwork, courtyard nadumuttam, and lush tropical garden.",
			ListingType:    models.ListingSale,
			Price:          12500000,
			PriceFrequency: "total",
			Bedrooms:       5,
			Bathrooms:      5,
			AreaSqFt:       4200,
			ParkingSpaces:  3,
			Address:        "Kazhakoottam, Near Technopark Phase 3",
			City:           "Thiruvananthapuram",
			State:          "Kerala",
			PostalCode:     "695582",
			CoverImage:     "https://images.unsplash.com/photo-1600585154340-be6161a56a0c?auto=format&fit=crop&w=800&q=80",
		},
		{
			ID:             3,
			Title:          "Malabar Modern Sea-Breeze Apartment",
			Slug:           "malabar-modern-sea-breeze-apartment",
			Description:    "Contemporary 3 BHK sea-facing residence in Calicut with world-class clubhouse amenities, swimming pool, and dedicated EV charging bay.",
			ListingType:    models.ListingRent,
			Price:          45000,
			PriceFrequency: "month",
			Bedrooms:       3,
			Bathrooms:      3,
			AreaSqFt:       1950,
			ParkingSpaces:  1,
			Address:        "Beach Road, Near Kozhikode Beach",
			City:           "Kozhikode",
			State:          "Kerala",
			PostalCode:     "673001",
			CoverImage:     "https://images.unsplash.com/photo-1545324418-cc1a3fa10c00?auto=format&fit=crop&w=800&q=80",
		},
		{
			ID:             4,
			Title:          "Sobha City Luxury Lakeview Apartment",
			Slug:           "sobha-city-luxury-lakeview-apartment",
			Description:    "High-end 3 BHK apartment in gated township with 25-acre manmade lake, jogging track, supermarket, and multi-tier security.",
			ListingType:    models.ListingSale,
			Price:          8900000,
			PriceFrequency: "total",
			Bedrooms:       3,
			Bathrooms:      3,
			AreaSqFt:       2100,
			ParkingSpaces:  2,
			Address:        "Puzhakkal, Guruvayur Road",
			City:           "Thrissur",
			State:          "Kerala",
			PostalCode:     "680553",
			CoverImage:     "https://images.unsplash.com/photo-1600585154526-990dced4db0d?auto=format&fit=crop&w=800&q=80",
		},
	}
}

// Global template helper functions
var FuncMap = template.FuncMap{
	"div": func(a, b float64) float64 {
		if b == 0 {
			return 0
		}
		return a / b
	},
	"formatPrice": formatINRCurrency,
}

// renderView parses files ensuring the root template carries FuncMap to all children
func renderView(c *gin.Context, viewFile string, data gin.H) {
	tmpl := template.New("").Funcs(FuncMap)
	files := []string{"templates/layouts/base.html", filepath.Join("templates/web", viewFile)}
	
	parsedTmpl, err := tmpl.ParseFiles(files...)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parsing error: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	// Execute the base layout
	if err := parsedTmpl.ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.String(http.StatusInternalServerError, "Template execution error: %v", err)
	}
}

func (h *WebHandler) Home(c *gin.Context) {
	featured := getSampleProperties()
	renderView(c, "index.html", gin.H{
		"FeaturedProperties": featured,
	})
}

func (h *WebHandler) Properties(c *gin.Context) {
	city := c.Query("city")
	listingType := c.Query("type")
	all := getSampleProperties()

	var filtered []models.Property
	for _, p := range all {
		if city != "" && p.City != city {
			continue
		}
		if listingType != "" && string(p.ListingType) != listingType {
			continue
		}
		filtered = append(filtered, p)
	}
	if len(filtered) == 0 {
		filtered = all
	}

	renderView(c, "properties.html", gin.H{
		"Properties": filtered,
		"Filter": models.PropertyFilter{
			City:        city,
			ListingType: listingType,
		},
	})
}

func (h *WebHandler) PropertyDetail(c *gin.Context) {
	slug := c.Param("slug")
	all := getSampleProperties()

	var property *models.Property
	for _, p := range all {
		if p.Slug == slug {
			property = &p
			break
		}
	}

	if property == nil {
		property = &all[0]
	}

	csrfToken, _ := c.Get("csrf_token")

	renderView(c, "property_detail.html", gin.H{
		"Property":       property,
		"CSRFToken":      csrfToken,
		"InquirySuccess": c.Query("inquiry") == "success",
	})
}

func (h *WebHandler) PostInquiry(c *gin.Context) {
	slug := c.Param("slug")
	c.Redirect(http.StatusSeeOther, "/properties/"+slug+"?inquiry=success#inquiry")
}
