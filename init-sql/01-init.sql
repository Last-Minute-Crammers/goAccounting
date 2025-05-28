-- 使用预先创建的数据库
USE Accounting;

-- 插入测试用户数据
INSERT INTO user (username, email, password, created_at, updated_at) 
VALUES 
('testuser', 'test@example.com', '$2a$10$1ozLELaQcXN1zukiLyyptuuJ06FdN9CWKh60Y7hhJaLDP3PlY1ldq', NOW(), NOW());

-- 插入基础分类数据
INSERT INTO category (name, icon, income_expense, created_at, updated_at)
VALUES 
('餐饮', 'food', 1, NOW(), NOW()),
('交通', 'transport', 1, NOW(), NOW()),
('工资', 'income', 0, NOW(), NOW());

-- 插入示例交易记录
INSERT INTO transaction (user_id, category_id, amount, description, transaction_time, created_at, updated_at)
VALUES 
(1, 1, 35.50, '午餐', NOW(), NOW(), NOW()),
(1, 3, 5000.00, '月薪', NOW(), NOW(), NOW());