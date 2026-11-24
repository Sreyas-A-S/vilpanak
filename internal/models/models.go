package models

import (
	"time"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleAdmin      UserRole = "admin"
	RoleAgent      UserRole = "agent"
)

type PropertyStatus string

const (
	StatusDraft     PropertyStatus = "draft"
	StatusPublished PropertyStatus = "published"
	StatusSold      PropertyStatus = "sold"
	StatusRented    PropertyStatus = "rented"
	StatusArchived  PropertyStatus = "archived"
)

type ListingType string

const (
	ListingSale  ListingType = "sale"
	ListingRent  ListingType = "rent"
	ListingLease ListingType = "lease"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	Phone        string    `json:"phone"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type Amenity struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type PropertyImage struct {
	ID           int64     `json:"id"`
	PropertyID   int64     `json:"property_id"`
	ImageURL     string    `json:"image_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	SortOrder    int       `json:"sort_order"`
	IsCover      bool      `json:"is_cover"`
	CreatedAt    time.Time `json:"created_at"`
}

type Property struct {
	ID             int64          `json:"id"`
	AgentID        *int64         `json:"agent_id"`
	AgentName      string         `json:"agent_name,omitempty"`
	CategoryID     int            `json:"category_id"`
	CategoryName   string         `json:"category_name,omitempty"`
	Title          string         `json:"title"`
	Slug           string         `json:"slug"`
	Description    string         `json:"description"`
	ListingType    ListingType    `json:"listing_type"`
	Status         PropertyStatus `json:"status"`
	Price          float64        `json:"price"`
	PriceFrequency string         `json:"price_frequency"`
	Bedrooms       int            `json:"bedrooms"`
	Bathrooms      int            `json:"bathrooms"`
	AreaSqFt       float64        `json:"area_sqft"`
	YearBuilt      *int           `json:"year_built"`
	Furnished      bool           `json:"furnished"`
	ParkingSpaces  int            `json:"parking_spaces"`
	Address        string         `json:"address"`
	City           string         `json:"city"`
	State          string         `json:"state"`
	PostalCode     string         `json:"postal_code"`
	Latitude       *float64       `json:"latitude"`
	Longitude      *float64       `json:"longitude"`
	IsFeatured     bool           `json:"is_featured"`
	ViewCount      int            `json:"view_count"`
	Images         []PropertyImage `json:"images,omitempty"`
	CoverImage     string         `json:"cover_image,omitempty"`
	Amenities      []Amenity      `json:"amenities,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type Lead struct {
	ID            int64     `json:"id"`
	PropertyID    *int64    `json:"property_id"`
	PropertyTitle string    `json:"property_title,omitempty"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Message       string    `json:"message"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type PropertyFilter struct {
	ListingType string
	CategoryID  int
	City        string
	MinPrice    float64
	MaxPrice    float64
	Bedrooms    int
	SearchQuery string
	SortBy      string // newest, price_asc, price_desc
	Page        int
	Limit       int
}
