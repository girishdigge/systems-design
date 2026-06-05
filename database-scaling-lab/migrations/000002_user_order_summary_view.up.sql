CREATE VIEW user_order_summary AS
SELECT
    u.id,
    COUNT(o.id)::BIGINT AS total_orders
FROM users u
LEFT JOIN orders o
    ON u.id = o.user_id
GROUP BY u.id;