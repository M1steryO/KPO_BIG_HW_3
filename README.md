# 🧩 Микросервисное приложение для заказов и платежей

Этот проект реализует микросервисную архитектуру с тремя основными сервисами:

- **API Gateway** — единая точка входа в систему.
- **Order Service** — отвечает за создание и хранение заказов.
- **Payments Service** — обрабатывает счета пользователей, пополнение баланса и списание при оплате.

## 📦 Стек технологий

- Go
- Kafka
- PostgreSQL
- Docker / Docker Compose
- Swagger

## 🚀 Быстрый старт

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/your_username/kpo_big_hw_3.git 
```
### 2.  Переходите в папку проекта
```bash
cd kpo_big_hw_3
```

### 3. Запускаете через docker-compose
```bash
docker compose up --build
```


### 4. После запуска спецификация API
будет доступна по адресу:
```bash
http://localhost:8083/swagger/index.html
```



