CREATE TABLE restaurant_orders (
                                   restaurant_id INT NOT NULL,
                                   order_id INT NOT NULL,
                                   menu_id INT NOT NULL,
                                   status VARCHAR(50) NOT NULL,
                                   created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                   updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                   PRIMARY KEY (restaurant_id, order_id, menu_id)
);
