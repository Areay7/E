CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id VARCHAR(100) NOT NULL COMMENT '平台订单ID',
    platform VARCHAR(50) NOT NULL COMMENT '平台名称',
    shop_id VARCHAR(100) COMMENT '店铺ID',
    order_sn VARCHAR(100) COMMENT '订单编号',
    status VARCHAR(50) COMMENT '订单状态',
    payment_status VARCHAR(50) COMMENT '支付状态',
    shipping_status VARCHAR(50) COMMENT '物流状态',

    buyer_username VARCHAR(200) COMMENT '买家用户名',
    buyer_email VARCHAR(200) COMMENT '买家邮箱',

    recipient_name VARCHAR(200) COMMENT '收件人姓名',
    recipient_phone VARCHAR(50) COMMENT '收件人电话',
    shipping_address TEXT COMMENT '收货地址',
    country VARCHAR(100) COMMENT '国家',
    province VARCHAR(100) COMMENT '省份',
    city VARCHAR(100) COMMENT '城市',
    district VARCHAR(100) COMMENT '区县',
    zipcode VARCHAR(20) COMMENT '邮编',

    total_amount DECIMAL(15,2) DEFAULT 0 COMMENT '订单总额',
    currency VARCHAR(10) DEFAULT 'USD' COMMENT '货币',
    shipping_fee DECIMAL(15,2) DEFAULT 0 COMMENT '运费',
    discount_amount DECIMAL(15,2) DEFAULT 0 COMMENT '优惠金额',
    actual_amount DECIMAL(15,2) DEFAULT 0 COMMENT '实付金额',

    tracking_number VARCHAR(100) COMMENT '物流单号',
    shipping_carrier VARCHAR(100) COMMENT '物流公司',
    shipped_at DATETIME COMMENT '发货时间',
    delivered_at DATETIME COMMENT '签收时间',

    buyer_message TEXT COMMENT '买家留言',
    seller_note TEXT COMMENT '卖家备注',

    raw_data TEXT COMMENT '原始数据JSON',

    order_time DATETIME NOT NULL COMMENT '下单时间',
    payment_time DATETIME COMMENT '支付时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME COMMENT '软删除时间',

    INDEX idx_platform (platform),
    INDEX idx_shop_id (shop_id),
    INDEX idx_order_sn (order_sn),
    INDEX idx_status (status),
    INDEX idx_country (country),
    INDEX idx_tracking_number (tracking_number),
    INDEX idx_order_time (order_time),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX idx_platform_order (platform, order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',

    platform VARCHAR(50) COMMENT '平台名称',
    item_id VARCHAR(100) COMMENT '商品项ID',
    product_id VARCHAR(100) COMMENT '商品ID',
    variation_id VARCHAR(100) COMMENT '变体ID',

    product_name VARCHAR(500) COMMENT '商品名称',
    variation_name VARCHAR(500) COMMENT '变体名称',
    sku VARCHAR(100) COMMENT 'SKU',

    quantity INT DEFAULT 0 COMMENT '数量',
    unit_price DECIMAL(15,2) DEFAULT 0 COMMENT '单价',
    total_price DECIMAL(15,2) DEFAULT 0 COMMENT '总价',
    discount_amount DECIMAL(15,2) DEFAULT 0 COMMENT '优惠金额',

    image_url VARCHAR(500) COMMENT '商品图片',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME COMMENT '软删除时间',

    INDEX idx_order_id (order_id),
    INDEX idx_product_id (product_id),
    INDEX idx_sku (sku),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单明细表';


CREATE TABLE IF NOT EXISTS products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) NOT NULL COMMENT '平台名称',
    product_id VARCHAR(100) NOT NULL COMMENT '商品ID',
    shop_id VARCHAR(100) COMMENT '店铺ID',

    name VARCHAR(500) COMMENT '商品名称',
    description TEXT COMMENT '商品描述',
    category VARCHAR(200) COMMENT '分类',
    brand VARCHAR(200) COMMENT '品牌',

    status VARCHAR(50) COMMENT '状态',
    price DECIMAL(15,2) DEFAULT 0 COMMENT '价格',
    original_price DECIMAL(15,2) DEFAULT 0 COMMENT '原价',
    currency VARCHAR(10) DEFAULT 'USD' COMMENT '货币',

    stock INT DEFAULT 0 COMMENT '库存',
    sku VARCHAR(100) COMMENT 'SKU',

    main_image VARCHAR(500) COMMENT '主图',
    images TEXT COMMENT '图片列表JSON',

    weight DECIMAL(10,2) DEFAULT 0 COMMENT '重量(kg)',
    length DECIMAL(10,2) DEFAULT 0 COMMENT '长度(cm)',
    width DECIMAL(10,2) DEFAULT 0 COMMENT '宽度(cm)',
    height DECIMAL(10,2) DEFAULT 0 COMMENT '高度(cm)',

    sold_count INT DEFAULT 0 COMMENT '销量',
    view_count INT DEFAULT 0 COMMENT '浏览量',
    rating DECIMAL(3,2) DEFAULT 0 COMMENT '评分',
    review_count INT DEFAULT 0 COMMENT '评论数',

    raw_data TEXT COMMENT '原始数据JSON',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME COMMENT '软删除时间',

    INDEX idx_platform (platform),
    INDEX idx_shop_id (shop_id),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_sku (sku),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX idx_platform_product (platform, product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';


CREATE TABLE IF NOT EXISTS product_variations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',

    platform VARCHAR(50) COMMENT '平台名称',
    variation_id VARCHAR(100) COMMENT '变体ID',

    name VARCHAR(500) COMMENT '变体名称',
    sku VARCHAR(100) COMMENT 'SKU',

    price DECIMAL(15,2) DEFAULT 0 COMMENT '价格',
    stock INT DEFAULT 0 COMMENT '库存',

    attributes TEXT COMMENT '属性JSON',
    image_url VARCHAR(500) COMMENT '图片',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME COMMENT '软删除时间',

    INDEX idx_product_id (product_id),
    INDEX idx_variation_id (variation_id),
    INDEX idx_sku (sku),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品变体表';


CREATE TABLE IF NOT EXISTS inventories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) NOT NULL COMMENT '平台名称',
    sku VARCHAR(100) NOT NULL COMMENT 'SKU',
    product_id VARCHAR(100) COMMENT '商品ID',
    variation_id VARCHAR(100) COMMENT '变体ID',

    stock INT DEFAULT 0 COMMENT '总库存',
    reserved_stock INT DEFAULT 0 COMMENT '预留库存',
    available_stock INT DEFAULT 0 COMMENT '可用库存',

    warehouse_code VARCHAR(50) COMMENT '仓库代码',
    location VARCHAR(200) COMMENT '库位',

    last_sync_at DATETIME COMMENT '最后同步时间',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME COMMENT '软删除时间',

    INDEX idx_platform (platform),
    INDEX idx_product_id (product_id),
    INDEX idx_warehouse_code (warehouse_code),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX idx_platform_sku (platform, sku)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存表';

CREATE TABLE IF NOT EXISTS inventory_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    inventory_id BIGINT UNSIGNED NOT NULL COMMENT '库存ID',

    type VARCHAR(50) COMMENT '变动类型',
    quantity INT DEFAULT 0 COMMENT '变动数量',
    before_stock INT DEFAULT 0 COMMENT '变动前库存',
    after_stock INT DEFAULT 0 COMMENT '变动后库存',

    reason VARCHAR(500) COMMENT '变动原因',
    related_order_id VARCHAR(100) COMMENT '关联订单ID',
    operator_id BIGINT UNSIGNED COMMENT '操作人ID',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_inventory_id (inventory_id),
    INDEX idx_type (type),
    INDEX idx_related_order_id (related_order_id),
    FOREIGN KEY (inventory_id) REFERENCES inventories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存变动日志表';


CREATE TABLE IF NOT EXISTS logistics (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',

    platform VARCHAR(50) COMMENT '平台名称',
    tracking_number VARCHAR(100) COMMENT '物流单号',
    carrier VARCHAR(100) COMMENT '物流公司',

    status VARCHAR(50) COMMENT '物流状态',
    current_location VARCHAR(500) COMMENT '当前位置',

    shipped_at DATETIME COMMENT '发货时间',
    in_transit_at DATETIME COMMENT '运输中时间',
    delivered_at DATETIME COMMENT '签收时间',

    estimated_delivery DATETIME COMMENT '预计送达时间',

    tracking_events TEXT COMMENT '物流轨迹JSON',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME COMMENT '软删除时间',

    INDEX idx_order_id (order_id),
    INDEX idx_tracking_number (tracking_number),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物流信息表';

CREATE TABLE IF NOT EXISTS platform_configs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) NOT NULL COMMENT '平台名称',
    shop_id VARCHAR(100) COMMENT '店铺ID',
    shop_name VARCHAR(200) COMMENT '店铺名称',

    enabled TINYINT(1) DEFAULT 1 COMMENT '是否启用',

    app_key VARCHAR(200) COMMENT 'App Key',
    app_secret VARCHAR(500) COMMENT 'App Secret',
    partner_id VARCHAR(200) COMMENT 'Partner ID',
    partner_key VARCHAR(500) COMMENT 'Partner Key',

    access_token TEXT COMMENT 'Access Token',
    refresh_token TEXT COMMENT 'Refresh Token',
    token_expires_at DATETIME COMMENT 'Token过期时间',

    api_url VARCHAR(500) COMMENT 'API地址',

    sync_enabled TINYINT(1) DEFAULT 1 COMMENT '是否启用同步',
    last_sync_at DATETIME COMMENT '最后同步时间',

    settings TEXT COMMENT '其他设置JSON',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME COMMENT '软删除时间',

    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX idx_platform (platform)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台配置表';

CREATE TABLE IF NOT EXISTS api_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) COMMENT '平台名称',
    api_name VARCHAR(200) COMMENT 'API名称',
    method VARCHAR(10) COMMENT '请求方法',
    url VARCHAR(1000) COMMENT '请求URL',

    request_headers TEXT COMMENT '请求头',
    request_body TEXT COMMENT '请求体',

    response_status INT COMMENT '响应状态码',
    response_headers TEXT COMMENT '响应头',
    response_body TEXT COMMENT '响应体',

    duration BIGINT COMMENT '耗时(ms)',
    success TINYINT(1) DEFAULT 0 COMMENT '是否成功',
    error_message TEXT COMMENT '错误信息',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_platform (platform),
    INDEX idx_api_name (api_name),
    INDEX idx_success (success),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API调用日志表';

CREATE TABLE IF NOT EXISTS sync_tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) COMMENT '平台名称',
    task_type VARCHAR(50) COMMENT '任务类型',
    status VARCHAR(50) COMMENT '任务状态',

    start_time DATETIME COMMENT '开始时间',
    end_time DATETIME COMMENT '结束时间',

    total_count INT DEFAULT 0 COMMENT '总数',
    success_count INT DEFAULT 0 COMMENT '成功数',
    fail_count INT DEFAULT 0 COMMENT '失败数',

    error_message TEXT COMMENT '错误信息',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_platform (platform),
    INDEX idx_task_type (task_type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步任务表';
