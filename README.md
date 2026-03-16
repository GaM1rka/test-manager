# Todo API

HTTP сервер для управления задачами.

## Запуск

```bash
# Клонировать репозиторий
git clone https://github.com/GaM1rka/test-manager.git
cd test-manager

# Сборка и запуск
docker build -t server .
docker run -p 8080:8080 server
```

## Api Эндпоинты
`POST /todos/` - Создание новой задачи.
Тело запроса:
```json
{
  "title": "Купить молоко",
  "description": "2 литра"
}
```
`GET /todos/` - Получение всех задач.
Пример запроса: GET /todos/

`GET /todos/{id}` - Получение задачи по id.
Пример запроса: GET /todos/1

`PUT /todos/{id}` - Обновление задачи по id.
Тело запроса:
```json
{
  "title": "Купить молоко (обновлено)",
  "description": "2 литра + хлеб", 
  "completed": true
}
```

`DELETE /todos/{id}` - Удаление задачи по id.
Пример запроса: DELETE /todos/1
