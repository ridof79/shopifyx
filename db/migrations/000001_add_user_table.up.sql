DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS bank_accounts CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS payments_counter CASCADE;

-- Create the extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL CHECK (LENGTH(username) >= 5 AND LENGTH(username) <= 15),
    name VARCHAR(100) NOT NULL CHECK (LENGTH(name) >= 5 AND LENGTH(name) <= 50),
    password VARCHAR(200) NOT NULL CHECK (LENGTH(password) >= 5)
);

CREATE INDEX idx_users_id ON users (id);

-- Products table
CREATE TABLE products (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(100) NOT NULL CHECK (LENGTH(name) >= 5 AND LENGTH(name) <= 60),
    price NUMERIC NOT NULL CHECK (price >= 0),
    image_url TEXT NOT NULL CHECK (image_url ~* '^https?://'),
    stock INTEGER NOT NULL CHECK (stock >= 0),
    condition VARCHAR(10) NOT NULL CHECK (condition IN ('new', 'second')),
    tags TEXT[] NOT NULL DEFAULT '{}',
    is_purchaseable BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_id UUID REFERENCES users(id)
);

CREATE INDEX idx_products_id ON products (id);

-- Bank accounts table
CREATE TABLE bank_accounts (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    bank_name VARCHAR(50) NOT NULL CHECK (LENGTH(bank_name) >= 5 AND LENGTH(bank_name) <= 15),
    bank_account_name VARCHAR(100) NOT NULL CHECK (LENGTH(bank_account_name) >= 5 AND LENGTH(bank_account_name) <= 15),
    bank_account_number VARCHAR(30) NOT NULL CHECK (LENGTH(bank_account_number) >= 5 AND LENGTH(bank_account_number) <= 15),
    user_id UUID REFERENCES users(id)
);

CREATE INDEX idx_bank_accounts_id ON bank_accounts (id);

-- Payments table
CREATE TABLE payments (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    bank_account_id UUID NOT NULL REFERENCES bank_accounts(id),
    payment_proof_image_url TEXT NOT NULL CHECK (payment_proof_image_url ~* '^https?://'),
    buyer_id UUID NOT NULL REFERENCES users(id),
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity >= 1)
);

-- View: Total products sold
CREATE VIEW total_product_sold AS
SELECT 
    product_id,
    COALESCE(SUM(quantity), 0) AS total_sold
FROM payments
GROUP BY product_id;

-- View: Total products sold by users
CREATE VIEW total_users_sold AS
SELECT u.id AS user_id,
       u.username AS username,
       COALESCE(SUM(py.quantity), 0) AS total_sold
FROM users u
LEFT JOIN products p ON u.id = p.user_id
LEFT JOIN payments py ON p.id = py.product_id
GROUP BY u.id, u.username;

-- View: Total seller_id from bank_accounts
CREATE VIEW seller_bank_account AS
SELECT 
    p.id AS product_id,
    u.id AS seller_id,
    ba.id AS bank_account_id,
    p.is_purchaseable 
FROM products p
JOIN users u ON p.user_id = u.id
JOIN bank_accounts ba ON u.id = ba.user_id;
