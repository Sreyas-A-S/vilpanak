-- Schema Migration: 000001_init_schema.up.sql

CREATE TYPE user_role AS ENUM ('super_admin', 'admin', 'agent');
CREATE TYPE property_status AS ENUM ('draft', 'published', 'sold', 'rented', 'archived');
CREATE TYPE listing_type AS ENUM ('sale', 'rent', 'lease');

-- 1. Users Table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'agent',
    phone VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Categories
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Amenities
CREATE TABLE IF NOT EXISTS amenities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50) DEFAULT 'check-circle'
);

-- 4. Properties
CREATE TABLE IF NOT EXISTS properties (
    id BIGSERIAL PRIMARY KEY,
    agent_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    category_id INT REFERENCES categories(id) ON DELETE RESTRICT,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(280) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    listing_type listing_type NOT NULL DEFAULT 'sale',
    status property_status NOT NULL DEFAULT 'draft',
    price NUMERIC(14, 2) NOT NULL,
    price_frequency VARCHAR(20) DEFAULT 'total',
    
    -- Specific Attributes
    bedrooms SMALLINT DEFAULT 0,
    bathrooms SMALLINT DEFAULT 0,
    area_sqft NUMERIC(10, 2) NOT NULL DEFAULT 0,
    year_built INT,
    furnished BOOLEAN DEFAULT FALSE,
    parking_spaces SMALLINT DEFAULT 0,
    
    -- Location & Geo
    address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20),
    latitude NUMERIC(10, 7),
    longitude NUMERIC(10, 7),
    
    -- Flags & Counters
    is_featured BOOLEAN DEFAULT FALSE,
    view_count INT DEFAULT 0,
    
    -- Full-text search vector
    search_vector TSVECTOR,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Property Amenities Pivot
CREATE TABLE IF NOT EXISTS property_amenities (
    property_id BIGINT REFERENCES properties(id) ON DELETE CASCADE,
    amenity_id INT REFERENCES amenities(id) ON DELETE CASCADE,
    PRIMARY KEY(property_id, amenity_id)
);

-- 6. Property Images
CREATE TABLE IF NOT EXISTS property_images (
    id BIGSERIAL PRIMARY KEY,
    property_id BIGINT REFERENCES properties(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500) NOT NULL,
    sort_order INT DEFAULT 0,
    is_cover BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 7. Customer Leads & Inquiries
CREATE TABLE IF NOT EXISTS leads (
    id BIGSERIAL PRIMARY KEY,
    property_id BIGINT REFERENCES properties(id) ON DELETE SET NULL,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'new',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_properties_slug ON properties(slug);
CREATE INDEX IF NOT EXISTS idx_properties_status ON properties(status);
CREATE INDEX IF NOT EXISTS idx_properties_filter ON properties(status, listing_type, category_id, city, price);
CREATE INDEX IF NOT EXISTS idx_property_images_prop ON property_images(property_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_leads_prop ON leads(property_id);
