CREATE TABLE orders (
  order_id            SERIAL PRIMARY KEY NOT NULL,
  customer_id         INTEGER    NOT NULL,
  amount              NUMERIC(20,4)  CHECK(amount >= 0)     NOT NULL,
  currency            VARCHAR(3) NOT NULL
);