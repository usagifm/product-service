CREATE TABLE products(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    description TEXT NOT NULL,
    quantity INT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz DEFAULT (now())
);
