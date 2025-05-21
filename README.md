# Shortener
`Shortener` - сервис для сокращения ссылок, написаный на стеке: `Golang`, `Echo`, `GORM`, `PostgreSQL`, `Swagger`, `Machinery`, `Redis`, `RabbitMQ`, `Docker`.
## Установка
```bash
git clone https://github.com/involved-entity/shortener
cd shortener
```
## Запуск
Настройте конфиг `config/local.docker.yml` по примеру `config/local.docker.example.yml`.
```bash
docker compose up
```
## Функционал
Реализован функционал авторизации и регистрации (подтверждение почты, сброс пароля, JWT), сокращения ссылок, просмотра кликов по ссылке (данные о браузере, IP, реферере, языке браузера). Для ручек подготовлены `End-To-End` тесты. Присутствует Swagger документация.
