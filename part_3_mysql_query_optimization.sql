-- ============================================================================
-- Part 3: Complex MySQL Query Optimization
-- ============================================================================

-- ============================================================================
-- 1. RELATIONAL SCHEMA
-- ============================================================================

CREATE TABLE `categories` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name`       VARCHAR(100)    NOT NULL,
    `created_at` TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `products` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name`        VARCHAR(255)    NOT NULL,
    `category_id` BIGINT UNSIGNED NOT NULL,
    `price`       DECIMAL(12,2)   NOT NULL DEFAULT 0.00,
    `created_at`  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_products_category` FOREIGN KEY (`category_id`)
        REFERENCES `categories`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `customers` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name`       VARCHAR(255)    NOT NULL,
    `email`      VARCHAR(255)    NOT NULL,
    `created_at` TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_customers_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `orders` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `customer_id` BIGINT UNSIGNED NOT NULL,
    `product_id`  BIGINT UNSIGNED NOT NULL,
    `quantity`    INT UNSIGNED    NOT NULL DEFAULT 1,
    `total_price` DECIMAL(12,2)   NOT NULL DEFAULT 0.00,
    `created_at`  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_orders_customer` FOREIGN KEY (`customer_id`)
        REFERENCES `customers`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT `fk_orders_product` FOREIGN KEY (`product_id`)
        REFERENCES `products`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `order_items` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `order_id`   BIGINT UNSIGNED NOT NULL,
    `product_id` BIGINT UNSIGNED NOT NULL,
    `quantity`    INT UNSIGNED    NOT NULL DEFAULT 1,
    `unit_price`  DECIMAL(12,2)   NOT NULL DEFAULT 0.00,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_order_items_order` FOREIGN KEY (`order_id`)
        REFERENCES `orders`(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_order_items_product` FOREIGN KEY (`product_id`)
        REFERENCES `products`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- 2. INDEXES FOR OPTIMIZATION
-- ============================================================================

CREATE INDEX `idx_products_category_id` ON `products`(`category_id`);

CREATE INDEX `idx_orders_customer_id` ON `orders`(`customer_id`);

CREATE INDEX `idx_orders_product_id` ON `orders`(`product_id`);

CREATE INDEX `idx_orders_created_at` ON `orders`(`created_at`);

CREATE INDEX `idx_orders_customer_total` ON `orders`(`customer_id`, `total_price`);

CREATE INDEX `idx_order_items_order_id`   ON `order_items`(`order_id`);

CREATE INDEX `idx_order_items_product_id` ON `order_items`(`product_id`);

CREATE INDEX `idx_order_items_product_qty` ON `order_items`(`product_id`, `quantity`);

-- ============================================================================
-- 3. OPTIMIZED QUERIES
-- ============================================================================

-- --------------------------------------------------------------------------
-- Query 1: Retrieve all products with their category details
--           and total sold quantities
-- --------------------------------------------------------------------------
-- Strategy:
--   • LEFT JOIN to categories so every product appears even without a category.
--   • LEFT JOIN to a pre-aggregated subquery on order_items to get the total
--     sold quantity per product.  Pre-aggregating avoids doing a GROUP BY over
--     the full product × order_items cross product.
--   • The subquery benefits from idx_order_items_product_qty (covering index).
-- --------------------------------------------------------------------------

SELECT
    p.id            AS product_id,
    p.name          AS product_name,
    p.price         AS product_price,
    c.id            AS category_id,
    c.name          AS category_name,
    COALESCE(oi_agg.total_sold, 0) AS total_sold_quantity
FROM `products` p
JOIN `categories` c
    ON c.id = p.category_id
LEFT JOIN (
    SELECT
        `product_id`,
        SUM(`quantity`) AS total_sold
    FROM `order_items`
    GROUP BY `product_id`
) oi_agg
    ON oi_agg.product_id = p.id
ORDER BY p.id;

-- --------------------------------------------------------------------------
-- Query 2: Top 10 customers by total amount spent across all orders
-- --------------------------------------------------------------------------
-- Strategy:
--   • Direct GROUP BY on orders.customer_id with SUM(total_price).
--   • The covering index idx_orders_customer_total lets MySQL compute
--     the SUM entirely from the index (Index-Only Scan / "Using index").
--   • LIMIT 10 + ORDER BY DESC stops the filesort early.
-- --------------------------------------------------------------------------

SELECT
    c.id                    AS customer_id,
    c.name                  AS customer_name,
    c.email                 AS customer_email,
    SUM(o.total_price)      AS total_spent
FROM `customers` c
INNER JOIN `orders` o
    ON o.customer_id = c.id
GROUP BY c.id, c.name, c.email
ORDER BY total_spent DESC
LIMIT 10;

-- --------------------------------------------------------------------------
-- Query 3: Order history with related products and customers
--          (efficient indexing and joins, avoids N+1)
-- --------------------------------------------------------------------------
-- Strategy:
--   • Single query with JOINs — NOT a loop of queries per order (no N+1).
--   • orders → customers via idx_orders_customer_id + PK on customers.
--   • orders → order_items via idx_order_items_order_id.
--   • order_items → products via idx_order_items_product_id + PK on products.
--   • products → categories via idx_products_category_id + PK on categories.
--   • Optional: filter by date range or customer for pagination.
-- --------------------------------------------------------------------------


SELECT
    o.id                          AS order_id,
    o.created_at                  AS order_date,
    o.total_price                 AS order_total_price,

    cu.id                         AS customer_id,
    cu.name                       AS customer_name,
    cu.email                      AS customer_email,

    oi.id                         AS item_id,
    oi.quantity                   AS item_quantity,
    oi.unit_price                 AS item_unit_price,

    p.id                          AS product_id,
    p.name                        AS product_name,
    cat.name                      AS category_name
FROM `orders` o
INNER JOIN `customers` cu
    ON cu.id = o.customer_id
INNER JOIN `order_items` oi
    ON oi.order_id = o.id
INNER JOIN `products` p
    ON p.id = oi.product_id
INNER JOIN `categories` cat
    ON cat.id = p.category_id
ORDER BY o.created_at DESC, o.id, oi.id
LIMIT 100;

-- For paginated access (keyset / cursor pagination, more efficient than OFFSET):
-- WHERE o.created_at < '2026-01-01 00:00:00'   -- cursor from previous page
-- LIMIT 100;

-- ============================================================================
-- 4. OPTIMIZATION EXPLANATIONS & CONSIDERATIONS
-- ============================================================================

/*
┌─────────────────────────────────────────────────────────────────────────────┐
│  A. INDEXING STRATEGY                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  • Every FOREIGN KEY column has an index. MySQL InnoDB requires this for    │
│    FK constraints anyway, and it speeds up JOIN lookups.                    │
│                                                                             │
│  • COVERING INDEXES (composite indexes that include all columns needed      │
│    by a query) allow "Index-Only Scans" — MySQL reads the index B-tree     │
│    without touching the clustered (primary) index at all.                   │
│      – idx_orders_customer_total covers (customer_id, total_price)          │
│        → Query 2 can SUM(total_price) GROUP BY customer_id from index.     │
│      – idx_order_items_product_qty covers (product_id, quantity)            │
│        → Query 1's subquery can SUM(quantity) GROUP BY product_id.         │
│                                                                             │
│  • idx_orders_created_at enables efficient range scans for date filters     │
│    and ORDER BY created_at DESC (avoids filesort).                          │
│                                                                             │
│  • Unique index on customers.email prevents duplicates and speeds up        │
│    lookups by email.                                                        │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│  B. AVOIDING THE N+1 PROBLEM                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  The N+1 problem occurs when the application executes 1 query to fetch      │
│  a list of N parent rows, then N additional queries to fetch related        │
│  data for each row.                                                         │
│                                                                             │
│  We avoid this by:                                                          │
│  1. Using JOINs to fetch all related data in ONE query (see Query 3).       │
│  2. Using subquery pre-aggregation instead of per-row lookups (Query 1).    │
│  3. In application code (Go), use a single SQL query or at most two         │
│     queries (one for parents, one for all children with WHERE IN (...))     │
│     rather than a loop.                                                     │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│  C. LARGE DATASET OPTIMIZATION TECHNIQUES                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. PARTITIONING                                                            │
│     – Partition `orders` by RANGE on `created_at` (e.g., monthly).          │
│       This lets MySQL prune partitions when filtering by date range,        │
│       dramatically reducing I/O for historical queries.                     │
│                                                                             │
│       ALTER TABLE `orders` PARTITION BY RANGE (UNIX_TIMESTAMP(created_at))  │
│       (                                                                     │
│           PARTITION p2025 VALUES LESS THAN (UNIX_TIMESTAMP('2026-01-01')),   │
│           PARTITION p2026 VALUES LESS THAN (UNIX_TIMESTAMP('2027-01-01')),   │
│           PARTITION pmax  VALUES LESS THAN MAXVALUE                          │
│       );                                                                    │
│                                                                             │
│  2. MATERIALIZED VIEWS (simulated in MySQL with summary tables)             │
│     – For Query 2 (top customers), maintain a `customer_spending_summary`   │
│       table that is updated via triggers or periodic batch jobs.             │
│     – This turns an expensive GROUP BY + SUM over millions of rows into     │
│       a simple SELECT ... ORDER BY ... LIMIT 10.                            │
│                                                                             │
│       CREATE TABLE `customer_spending_summary` (                            │
│           `customer_id`  BIGINT UNSIGNED PRIMARY KEY,                        │
│           `total_spent`  DECIMAL(14,2) NOT NULL DEFAULT 0.00,               │
│           `order_count`  INT UNSIGNED  NOT NULL DEFAULT 0,                  │
│           `last_updated` TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP   │
│       );                                                                    │
│                                                                             │
│  3. QUERY CACHING (Application Level)                                       │
│     – Cache expensive aggregate results in Redis with a TTL.                │
│     – Example: cache the top-10-customers report for 5 minutes.             │
│                                                                             │
│  4. PAGINATION                                                              │
│     – Use KEYSET (cursor-based) pagination instead of OFFSET.               │
│       OFFSET forces MySQL to scan and discard rows; keyset uses indexed     │
│       WHERE conditions to jump directly to the next page.                   │
│                                                                             │
│  5. EXPLAIN & ANALYZE                                                       │
│     – Always run EXPLAIN ANALYZE on queries to verify that indexes are      │
│       being used and no full-table scans occur.                             │
│                                                                             │
│       EXPLAIN ANALYZE SELECT ... ;                                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

-- ============================================================================
-- 5. VERIFY QUERY PLANS (run these to confirm index usage)
-- ============================================================================

-- Check Query 1 plan
EXPLAIN SELECT
    p.id, p.name, p.price, c.id, c.name,
    COALESCE(oi_agg.total_sold, 0)
FROM `products` p
JOIN `categories` c ON c.id = p.category_id
LEFT JOIN (
    SELECT product_id, SUM(quantity) AS total_sold
    FROM order_items GROUP BY product_id
) oi_agg ON oi_agg.product_id = p.id
ORDER BY p.id;

-- Check Query 2 plan
EXPLAIN
SELECT c.id, c.name, c.email, SUM(o.total_price) AS total_spent
FROM customers c
    INNER JOIN orders o ON o.customer_id = c.id
GROUP BY
    c.id,
    c.name,
    c.email
ORDER BY total_spent DESC
LIMIT 10;

-- Check Query 3 plan
EXPLAIN
SELECT o.id, o.created_at, o.total_price, cu.id, cu.name, cu.email, oi.id, oi.quantity, oi.unit_price, p.id, p.name, cat.name
FROM
    orders o
    INNER JOIN customers cu ON cu.id = o.customer_id
    INNER JOIN order_items oi ON oi.order_id = o.id
    INNER JOIN products p ON p.id = oi.product_id
    INNER JOIN categories cat ON cat.id = p.category_id
ORDER BY o.created_at DESC, o.id, oi.id
LIMIT 100;