CREATE TABLE config (
    id SERIAL PRIMARY KEY,
    key_access VARCHAR(100) NOT NULL,
    key_refresh VARCHAR(100) NOT NULL,
    lifetime_access INT NOT NULL,
    lifetime_refresh INT NOT NULL
);