INSERT INTO platform_configs (platform, shop_id, shop_name, enabled, app_key, api_url, sync_enabled) VALUES
('shopee', 'SHOP001', 'Shopee旗舰店', 1, 'test_shopee_key', 'https://partner.shopeemobile.com', 1),
('aliexpress', 'AE001', '速卖通官方店', 1, 'test_ae_key', 'https://api.aliexpress.com', 1),
('tiktok', 'TK001', 'TikTok Shop', 1, 'test_tk_key', 'https://open-api.tiktokglobalshop.com', 1);

INSERT INTO products (platform, product_id, shop_id, name, description, category, status, price, original_price, currency, stock, sku, main_image, sold_count, view_count, rating, review_count) VALUES
('shopee', 'SP10001', 'SHOP001', '无线蓝牙耳机 TWS降噪入耳式', '高品质蓝牙5.0，主动降噪，超长续航', '电子产品/耳机', 'active', 29.99, 59.99, 'USD', 500, 'BT-EAR-001', 'https://example.com/images/earphone.jpg', 1250, 8500, 4.8, 320),
('shopee', 'SP10002', 'SHOP001', '智能手表运动手环', '心率监测，睡眠追踪，50米防水', '电子产品/智能穿戴', 'active', 45.50, 89.99, 'USD', 300, 'WATCH-001', 'https://example.com/images/watch.jpg', 680, 5200, 4.6, 180),
('aliexpress', 'AE20001', 'AE001', 'USB充电数据线 3条装', '快充支持，耐用编织线材，1米/2米可选', '电子配件/数据线', 'active', 8.99, 15.99, 'USD', 1000, 'CABLE-USB-001', 'https://example.com/images/cable.jpg', 3200, 15000, 4.9, 850),
('aliexpress', 'AE20002', 'AE001', '手机支架车载磁吸式', '360度旋转，强力磁吸，通用型', '汽车用品/手机支架', 'active', 12.50, 24.99, 'USD', 450, 'HOLDER-MAG-001', 'https://example.com/images/holder.jpg', 920, 6800, 4.7, 240),
('tiktok', 'TK30001', 'TK001', '便携式迷你风扇', 'USB充电，三档风速，静音设计', '家居用品/小家电', 'active', 15.99, 29.99, 'USD', 600, 'FAN-MINI-001', 'https://example.com/images/fan.jpg', 1580, 12000, 4.8, 420);


INSERT INTO product_variations (product_id, platform, variation_id, name, sku, price, stock, attributes) VALUES
(1, 'shopee', 'VAR-001', '黑色', 'BT-EAR-001-BLK', 29.99, 250, '{"color":"黑色"}'),
(1, 'shopee', 'VAR-002', '白色', 'BT-EAR-001-WHT', 29.99, 250, '{"color":"白色"}'),
(2, 'shopee', 'VAR-003', '黑色-标准版', 'WATCH-001-BLK-STD', 45.50, 150, '{"color":"黑色","version":"标准版"}'),
(2, 'shopee', 'VAR-004', '银色-运动版', 'WATCH-001-SLV-SPT', 49.99, 150, '{"color":"银色","version":"运动版"}');


INSERT INTO inventories (platform, sku, product_id, stock, reserved_stock, available_stock, warehouse_code, location) VALUES
('shopee', 'BT-EAR-001-BLK', 'SP10001', 250, 20, 230, 'WH-SG-01', 'A-01-05'),
('shopee', 'BT-EAR-001-WHT', 'SP10001', 250, 15, 235, 'WH-SG-01', 'A-01-06'),
('shopee', 'WATCH-001-BLK-STD', 'SP10002', 150, 10, 140, 'WH-SG-01', 'A-02-03'),
('aliexpress', 'CABLE-USB-001', 'AE20001', 1000, 50, 950, 'WH-CN-01', 'B-05-12'),
('aliexpress', 'HOLDER-MAG-001', 'AE20002', 450, 25, 425, 'WH-CN-01', 'B-06-08'),
('tiktok', 'FAN-MINI-001', 'TK30001', 600, 30, 570, 'WH-US-01', 'C-03-15');


INSERT INTO orders (order_id, platform, shop_id, order_sn, status, payment_status, shipping_status, buyer_username, buyer_email, recipient_name, recipient_phone, shipping_address, country, province, city, zipcode, total_amount, currency, shipping_fee, discount_amount, actual_amount, tracking_number, shipping_carrier, order_time, payment_time) VALUES
('SP2024033101', 'shopee', 'SHOP001', 'SH240331001', 'completed', 'paid', 'delivered', 'john_doe', 'john@example.com', 'John Doe', '+1234567890', '123 Main St, Apt 4B', 'United States', 'California', 'Los Angeles', '90001', 35.98, 'USD', 5.99, 0, 35.98, 'TRACK001', 'DHL', '2024-03-25 10:30:00', '2024-03-25 10:35:00'),
('SP2024033102', 'shopee', 'SHOP001', 'SH240331002', 'processing', 'paid', 'pending', 'mary_smith', 'mary@example.com', 'Mary Smith', '+1987654321', '456 Oak Ave', 'United States', 'New York', 'New York', '10001', 49.99, 'USD', 6.50, 5.00, 51.49, NULL, NULL, '2024-03-30 14:20:00', '2024-03-30 14:22:00'),
('AE2024033103', 'aliexpress', 'AE001', 'AE240331003', 'shipped', 'paid', 'shipped', 'alice_wang', 'alice@example.com', 'Alice Wang', '+8613800138000', '789 Beijing Road', 'China', 'Beijing', 'Beijing', '100000', 21.49, 'USD', 0, 0, 21.49, 'TRACK002', 'China Post', '2024-03-28 09:15:00', '2024-03-28 09:20:00'),
('AE2024033104', 'aliexpress', 'AE001', 'AE240331004', 'pending', 'pending', 'pending', 'bob_lee', 'bob@example.com', 'Bob Lee', '+821012345678', '321 Seoul Street', 'South Korea', 'Seoul', 'Gangnam', '06000', 12.50, 'USD', 3.00, 0, 15.50, NULL, NULL, '2024-03-31 16:45:00', NULL),
('TK2024033105', 'tiktok', 'TK001', 'TK240331005', 'completed', 'paid', 'delivered', 'emma_brown', 'emma@example.com', 'Emma Brown', '+447700900000', '555 London Road', 'United Kingdom', 'England', 'London', 'SW1A 1AA', 21.98, 'USD', 5.99, 0, 21.98, 'TRACK003', 'Royal Mail', '2024-03-26 11:00:00', '2024-03-26 11:05:00'),
('TK2024033106', 'tiktok', 'TK001', 'TK240331006', 'processing', 'paid', 'pending', 'david_kim', 'david@example.com', 'David Kim', '+14155552671', '888 Market St', 'United States', 'California', 'San Francisco', '94102', 31.98, 'USD', 7.50, 3.50, 35.98, NULL, NULL, '2024-03-31 08:30:00', '2024-03-31 08:32:00');


INSERT INTO order_items (order_id, platform, item_id, product_id, variation_id, product_name, variation_name, sku, quantity, unit_price, total_price, discount_amount, image_url) VALUES
(1, 'shopee', 'ITEM001', 'SP10001', 'VAR-001', '无线蓝牙耳机 TWS降噪入耳式', '黑色', 'BT-EAR-001-BLK', 1, 29.99, 29.99, 0, 'https://example.com/images/earphone.jpg'),
(2, 'shopee', 'ITEM002', 'SP10002', 'VAR-003', '智能手表运动手环', '黑色-标准版', 'WATCH-001-BLK-STD', 1, 45.50, 45.50, 5.00, 'https://example.com/images/watch.jpg'),
(3, 'aliexpress', 'ITEM003', 'AE20001', NULL, 'USB充电数据线 3条装', NULL, 'CABLE-USB-001', 2, 8.99, 17.98, 0, 'https://example.com/images/cable.jpg'),
(3, 'aliexpress', 'ITEM004', 'AE20002', NULL, '手机支架车载磁吸式', NULL, 'HOLDER-MAG-001', 1, 12.50, 12.50, 0, 'https://example.com/images/holder.jpg'),
(4, 'aliexpress', 'ITEM005', 'AE20002', NULL, '手机支架车载磁吸式', NULL, 'HOLDER-MAG-001', 1, 12.50, 12.50, 0, 'https://example.com/images/holder.jpg'),
(5, 'tiktok', 'ITEM006', 'TK30001', NULL, '便携式迷你风扇', NULL, 'FAN-MINI-001', 1, 15.99, 15.99, 0, 'https://example.com/images/fan.jpg'),
(6, 'tiktok', 'ITEM007', 'TK30001', NULL, '便携式迷你风扇', NULL, 'FAN-MINI-001', 2, 15.99, 31.98, 3.50, 'https://example.com/images/fan.jpg');


INSERT INTO logistics (order_id, platform, tracking_number, carrier, status, current_location, shipped_at, delivered_at) VALUES
(1, 'shopee', 'TRACK001', 'DHL', 'delivered', 'Los Angeles, CA', '2024-03-26 09:00:00', '2024-03-29 15:30:00'),
(3, 'aliexpress', 'TRACK002', 'China Post', 'in_transit', 'Beijing Distribution Center', '2024-03-29 10:00:00', NULL),
(5, 'tiktok', 'TRACK003', 'Royal Mail', 'delivered', 'London, UK', '2024-03-27 08:00:00', '2024-03-30 14:20:00');


INSERT INTO inventory_logs (inventory_id, type, quantity, before_stock, after_stock, reason, related_order_id) VALUES
(1, 'decrease', -1, 251, 250, '订单出库', 'SP2024033101'),
(3, 'decrease', -1, 151, 150, '订单出库', 'SP2024033102'),
(4, 'decrease', -2, 1002, 1000, '订单出库', 'AE2024033103'),
(5, 'decrease', -1, 451, 450, '订单出库', 'AE2024033103'),
(6, 'decrease', -1, 601, 600, '订单出库', 'TK2024033105');


INSERT INTO sync_tasks (platform, task_type, status, start_time, end_time, total_count, success_count, fail_count) VALUES
('shopee', 'orders', 'completed', '2024-03-31 00:00:00', '2024-03-31 00:05:30', 150, 150, 0),
('aliexpress', 'orders', 'completed', '2024-03-31 00:10:00', '2024-03-31 00:18:45', 280, 278, 2),
('tiktok', 'products', 'completed', '2024-03-31 01:00:00', '2024-03-31 01:12:20', 520, 520, 0),
('shopee', 'inventory', 'running', '2024-03-31 02:00:00', NULL, 0, 0, 0);


INSERT INTO api_logs (platform, api_name, method, url, response_status, duration, success, created_at) VALUES
('shopee', 'GetOrderList', 'POST', 'https://partner.shopeemobile.com/api/v2/order/get_order_list', 200, 1250, 1, '2024-03-31 00:00:15'),
('shopee', 'GetOrderDetail', 'POST', 'https://partner.shopeemobile.com/api/v2/order/get_order_detail', 200, 850, 1, '2024-03-31 00:01:20'),
('aliexpress', 'GetOrderList', 'GET', 'https://api.aliexpress.com/sync/orders', 200, 2100, 1, '2024-03-31 00:10:30'),
('tiktok', 'GetProducts', 'GET', 'https://open-api.tiktokglobalshop.com/api/products/search', 200, 1680, 1, '2024-03-31 01:00:45'),
('shopee', 'GetOrderList', 'POST', 'https://partner.shopeemobile.com/api/v2/order/get_order_list', 500, 3000, 0, '2024-03-31 02:30:00');
