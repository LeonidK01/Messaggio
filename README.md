## Контекст задачи
Менеджер может отправлять сообщения через REST API.


## Доступные REST запросы

**POST** /app/v1/send
```json
{
    "to": "uuid",
    "from": "uuid",
    "text": "string"
}
```