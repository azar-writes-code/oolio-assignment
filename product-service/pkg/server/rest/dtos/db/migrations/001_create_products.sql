-- Active: 1774208111445@@localhost@5432@invoice_db
-- Migrations: 001_create_products.sql
-- Create the products table for the online food ordering system

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'product_image') THEN
        CREATE TYPE product_image AS (
            thumbnail TEXT,
            mobile TEXT,
            tablet TEXT,
            desktop TEXT
        );
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    -- TODO: add restaurant table to map products (one restaurant can have multiple products)
    -- restaurant_id UUID NOT NULL REFERENCES restaurants(id), 
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category TEXT [] NOT NULL,
    image product_image,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE products IS 'Products table for the online food ordering system';