# Sneaker shop backend

Gin, Gorm, PostgreSQL, NATS


## Setup DB

```PostgreSQL
CREATE DATABASE sneaker_shop;
CREATE USER sneaker_shop_user WITH PASSWORD 'your_password_here';
GRANT CONNECT ON DATABASE sneaker_shop TO sneaker_shop_user;
\c sneaker_shop
GRANT USAGE ON SCHEMA public TO sneaker_shop_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO sneaker_shop_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO sneaker_shop_user;
```