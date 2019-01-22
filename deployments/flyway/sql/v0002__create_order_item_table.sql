CREATE TABLE order_items (
  order_item_id SERIAL PRIMARY KEY NOT NULL,
  order_id      INTEGER    NOT NULL,
  product_id    INTEGER    NOT NULL,
  quantity      INTEGER    NOT NULL,
  price         NUMERIC(20,4) CHECK(price >= 0)      NOT NULL,
  currency      VARCHAR(3) NOT NULL
);