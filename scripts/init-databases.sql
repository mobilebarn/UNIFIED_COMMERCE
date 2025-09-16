-- Database initialization script for Retail OS Platform
-- Creates separate databases for each microservice following the database-per-service pattern

-- Identity Service Database
CREATE DATABASE identity_service;
CREATE USER identity_user WITH ENCRYPTED PASSWORD 'identity_pass';
GRANT ALL PRIVILEGES ON DATABASE identity_service TO identity_user;

-- Merchant Account Service Database
CREATE DATABASE merchant_account_service;
CREATE USER merchant_user WITH ENCRYPTED PASSWORD 'merchant_pass';
GRANT ALL PRIVILEGES ON DATABASE merchant_account_service TO merchant_user;

-- Inventory Service Database
CREATE DATABASE inventory_service;
CREATE USER inventory_user WITH ENCRYPTED PASSWORD 'inventory_pass';
GRANT ALL PRIVILEGES ON DATABASE inventory_service TO inventory_user;

-- Order Service Database
CREATE DATABASE order_service;
CREATE USER order_user WITH ENCRYPTED PASSWORD 'order_pass';
GRANT ALL PRIVILEGES ON DATABASE order_service TO order_user;

-- Cart & Checkout Service Database
CREATE DATABASE cart_checkout_service;
CREATE USER cart_user WITH ENCRYPTED PASSWORD 'cart_pass';
GRANT ALL PRIVILEGES ON DATABASE cart_checkout_service TO cart_user;

-- Payments Service Database
CREATE DATABASE payments_service;
CREATE USER payments_user WITH ENCRYPTED PASSWORD 'payments_pass';
GRANT ALL PRIVILEGES ON DATABASE payments_service TO payments_user;

-- Promotions Service Database
CREATE DATABASE promotions_service;
CREATE USER promotions_user WITH ENCRYPTED PASSWORD 'promotions_pass';
GRANT ALL PRIVILEGES ON DATABASE promotions_service TO promotions_user;

-- Show created databases
\l