-- +goose Up
-- +goose StatementBegin
INSERT INTO categories (name, description, image_url) VALUES
('Игрушки', 'Разнообразные игрушки для взрослых', 'https://via.placeholder.com/300'),
('Белье', 'Сексуальное белье и аксессуары', 'https://via.placeholder.com/300'),
('Аксессуары', 'Различные аксессуары', 'https://via.placeholder.com/300')
ON CONFLICT (name) DO NOTHING;

INSERT INTO products (name, description, price, image_url, category, in_stock, rating, views) VALUES
('Игрушка премиум класса', 'Высококачественная игрушка из безопасных материалов', 2999.99, 'https://via.placeholder.com/400', 'Игрушки', true, 4.8, 150),
('Кружевное белье', 'Элегантное кружевное белье премиум качества', 1999.99, 'https://via.placeholder.com/400', 'Белье', true, 4.9, 200),
('Набор аксессуаров', 'Комплект аксессуаров для особых моментов', 3499.99, 'https://via.placeholder.com/400', 'Аксессуары', true, 4.7, 120),
('Игрушка стандарт', 'Классическая игрушка отличного качества', 1499.99, 'https://via.placeholder.com/400', 'Игрушки', true, 4.5, 180),
('Комплект белья', 'Стильный комплект белья из шелка', 2499.99, 'https://via.placeholder.com/400', 'Белье', true, 4.6, 160),
('Премиум аксессуар', 'Эксклюзивный аксессуар ручной работы', 4999.99, 'https://via.placeholder.com/400', 'Аксессуары', true, 5.0, 90),
('Игрушка компакт', 'Компактная версия для путешествий', 999.99, 'https://via.placeholder.com/400', 'Игрушки', true, 4.4, 140),
('Боди с кружевом', 'Сексуальное боди с кружевными вставками', 1799.99, 'https://via.placeholder.com/400', 'Белье', true, 4.7, 170)
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM products;
DELETE FROM categories;
-- +goose StatementEnd

