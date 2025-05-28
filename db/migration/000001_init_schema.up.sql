CREATE TABLE shops (
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    description varchar(255) NOT NULL,
    logo_url varchar(255),
    website_url varchar(255),
    email varchar(255),
    whatsapp_phone varchar(255) UNIQUE,
    address varchar(255) NOT NULL,
    city varchar(255) NOT NULL,
    state varchar(255) NOT NULL,
    zip_code varchar(255) NOT NULL,
    country varchar(255) NOT NULL,
    latitude float NOT NULL,
    longitude float NOT NULL,
    is_active boolean NOT NULL DEFAULT TRUE,
    slug varchar(255) UNIQUE DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    shop_id INTEGER NOT NULL REFERENCES shops(id) ON DELETE CASCADE,

    email varchar(255),
    unconfirmed_email varchar(255),
    phone varchar(255),
    unconfirmed_phone varchar(255),

    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    slug UUID UNIQUE DEFAULT gen_random_uuid(),

    UNIQUE (shop_id, email),
    UNIQUE (shop_id, phone),

    -- Enforce either email OR phone, not both, not neither
    CHECK (
        (email IS NULL AND phone IS NOT NULL) OR
        (email IS NOT NULL AND phone IS NULL)
    )
);

CREATE TABLE user_profiles (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    phone TEXT,
    first_name TEXT,
    last_name TEXT,
    address TEXT,
    city TEXT,
    country TEXT,
    postal_code TEXT,
    UNIQUE (user_id)
);

-- Roles: e.g. 'admin', 'manager', 'staff', 'customer'
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Role assignments (many-to-many)
CREATE TABLE user_roles (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Categories: Scoped to shop, hierarchical
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    shop_id INTEGER NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    parent_id INTEGER REFERENCES categories(id),
    UNIQUE (shop_id, name)
);

-- Products: Belong to a shop, and (optionally) a category
CREATE TABLE products (
    id varchar PRIMARY KEY,
    shop_id INTEGER NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id),
    
    name TEXT NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Orders: Created by users in a shop
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    shop_id INTEGER NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id),
    
    total DECIMAL(10,2),
    status TEXT CHECK (status IN ('pending', 'paid', 'shipped', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order items: Line items in an order
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id varchar REFERENCES products(id),
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL
);

-- user login otp
create table user_login_otp(
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	otp char(6) NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    expires_at TIMESTAMP DEFAULT (now() + INTERVAL '5 minutes')

);