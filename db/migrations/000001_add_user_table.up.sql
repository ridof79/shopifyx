DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS bank_accounts CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS payments_counter CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username VARCHAR(15) UNIQUE NOT NULL CHECK (LENGTH(username) >= 5),
    name VARCHAR(50) NOT NULL CHECK (LENGTH(name) >= 5),
    password VARCHAR(200) NOT NULL CHECK (LENGTH(password) >= 5)
);

CREATE INDEX idx_users_id ON users (id);

CREATE TABLE products (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(60) NOT NULL CHECK (
        LENGTH(name) >= 5
        AND LENGTH(name) <= 60
    ),
    price NUMERIC NOT NULL CHECK (price >= 0),
    image_url TEXT NOT NULL CHECK (image_url ~* '^https?://'),
    stock INTEGER NOT NULL CHECK (stock >= 0),
    condition VARCHAR(10) NOT NULL CHECK (condition IN ('new', 'second')),
    tags TEXT [] NOT NULL DEFAULT '{}',
    is_purchaseable BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_id UUID REFERENCES users(id)
);

CREATE INDEX idx_products_id ON products (id);

CREATE TABLE bank_accounts (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    bank_name VARCHAR(15) NOT NULL CHECK (
        LENGTH(bank_name) >= 5
        AND LENGTH(bank_name) <= 15
    ),
    bank_account_name VARCHAR(15) NOT NULL CHECK (
        LENGTH(bank_account_name) >= 5
        AND LENGTH(bank_account_name) <= 15
    ),
    bank_account_number VARCHAR(15) NOT NULL CHECK (
        LENGTH(bank_account_number) >= 5
        AND LENGTH(bank_account_number) <= 15
    ),
    user_id UUID REFERENCES users(id)
);

CREATE INDEX bank_account_id ON bank_accounts (id);

CREATE TABLE payments (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    bank_account_id UUID NOT NULL REFERENCES bank_accounts(id),
    payment_proof_image_url TEXT NOT NULL CHECK (payment_proof_image_url ~* '^https?://'),
    buyer_id UUID NOT NULL REFERENCES users(id)
);

CREATE TABLE payments_counter (
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity >= 1),
    payment_id UUID NOT NULL REFERENCES payments(id),
    seller_id UUID NOT NULL REFERENCES users(id)
);

CREATE INDEX idx_payments_counter_product_id ON payments_counter (product_id);
CREATE INDEX idx_payments_counter_seller_id ON payments_counter (seller_id);