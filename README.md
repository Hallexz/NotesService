## Стек технологий

- Go
- PostgreSQL
- Docker

## Установка
1. Клонирование репозитория:

   ```
   git clone https://github.com/ваш-логин/NotesService.git
   cd NotesService
   ```
   
2. Собрать и запустить контейнеры:
   ```
   docker-compose up --build
   ```
   
Сервис использует REST API: 
1. Аутентификация: ```POST /auth```
   Тело запроса должно содержать JSON с полями username и password.
2. Создание заметки: ```POST /notes```
   Требует аутентификации (JWT). Тело запроса должно содержать JSON с данными заметки.
3. Получение заметок: ```GET /getNotes```
   Требует аутентификации (JWT). Возвращает список заметок пользователя.
4. Логи
   Логи приложения записываются в файл app.log. Убедитесь, что у вас есть права на запись в этой папке.
