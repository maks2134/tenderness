# Tenderness - Магазин товаров 18+

Веб-приложение магазина товаров для взрослых, построенное на Fiber v3 (Go) и React.

## Технологии

### Backend
- **Go 1.25** с **Fiber v3**
- **PostgreSQL 17.4** для базы данных
- **sqlx** для работы с SQL
- **goose** для миграций базы данных

### Frontend
- **React 18** с **Vite**
- Современный адаптивный дизайн

## Структура проекта

```
tenderness/
├── server/          # Go backend
│   ├── cmd/         # Точка входа приложения
│   ├── internal/    # Внутренние пакеты
│   │   ├── app/     # Инициализация приложения
│   │   ├── configs/ # Конфигурация
│   │   ├── domain/  # Модели и storage
│   │   ├── handlers/# HTTP handlers
│   │   ├── repository/# Репозитории для работы с БД
│   │   ├── routes/  # Маршруты API
│   │   └── services/# Бизнес-логика
│   └── migrations/  # Миграции goose
└── client/          # React frontend
    └── src/         # Исходный код React
```

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Go 1.25+ (для локальной разработки)
- Node.js 20+ (для локальной разработки фронтенда)

### Запуск через Docker Compose

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd tenderness
```

2. Запустите все сервисы:
```bash
docker-compose up --build
```

Это запустит:
- PostgreSQL 17.4 на порту 5432
- Go backend на порту 3000 (внутренний)
- React frontend на порту 80 (http://localhost)

3. Откройте браузер и перейдите на http://localhost

### Локальная разработка

#### Backend

1. Перейдите в директорию server:
```bash
cd server
```

2. Установите зависимости:
```bash
go mod download
go mod tidy
```

3. Создайте файл `.env`:
```env
PORT=:3000
DB_HOST=localhost
DB_PORT=5432
DB_USER=tenderness
DB_PASSWORD=tenderness123
DB_NAME=tenderness_db
DB_SSLMODE=disable
```

4. Запустите PostgreSQL через Docker:
```bash
docker-compose up postgres -d
```

5. Запустите приложение:
```bash
go run cmd/tenderness/main.go
```

#### Frontend

1. Перейдите в директорию client:
```bash
cd client
```

2. Установите зависимости:
```bash
npm install
```

3. Запустите dev сервер:
```bash
npm run dev
```

Frontend будет доступен на http://localhost:5173

## API Endpoints

### Products
- `GET /api/products` - Получить список товаров (с пагинацией)
  - Query params: `page`, `limit`
- `GET /api/products/featured` - Получить рекомендуемые товары
  - Query params: `limit`
- `GET /api/products/search?q={query}` - Поиск товаров
  - Query params: `q`, `page`, `limit`
- `GET /api/products/category/:category` - Товары по категории
  - Query params: `page`, `limit`
- `GET /api/products/:id` - Получить товар по ID

### Categories
- `GET /api/categories` - Получить список категорий

### Health
- `GET /health` - Проверка здоровья сервиса

## Миграции

Миграции выполняются автоматически при запуске приложения. Для ручного управления миграциями используйте goose:

```bash
# В контейнере или локально
goose -dir migrations postgres "host=localhost user=tenderness password=tenderness123 dbname=tenderness_db sslmode=disable" up
```

## Разработка

### Добавление новых миграций

1. Создайте новый файл в `server/migrations/`:
```bash
touch server/migrations/004_your_migration.sql
```

2. Используйте формат goose:
```sql
-- +goose Up
-- +goose StatementBegin
-- Ваш SQL код
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Откат миграции
-- +goose StatementEnd
```

### Структура базы данных

#### Таблица products
- `id` - SERIAL PRIMARY KEY
- `created_at`, `updated_at` - TIMESTAMP
- `name` - VARCHAR(255) NOT NULL
- `description` - TEXT
- `price` - DECIMAL(10, 2) NOT NULL
- `image_url` - VARCHAR(500)
- `category` - VARCHAR(100)
- `in_stock` - BOOLEAN DEFAULT true
- `rating` - DECIMAL(3, 2) DEFAULT 0.00
- `views` - INTEGER DEFAULT 0

#### Таблица categories
- `id` - SERIAL PRIMARY KEY
- `created_at`, `updated_at` - TIMESTAMP
- `name` - VARCHAR(255) NOT NULL UNIQUE
- `description` - TEXT
- `image_url` - VARCHAR(500)

