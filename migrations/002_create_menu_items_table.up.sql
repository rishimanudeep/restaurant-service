CREATE TABLE IF NOT EXISTS menu_items (
    restaurant_id SERIAL REFERENCES restaurants(id),
    items JSONB
);
