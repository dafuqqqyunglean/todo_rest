# To Do Project
Для сборки проекта:

```bash
docker-compose build
```
Для запуска проекта:
```bash
docker-compose up
```

Миграции:
```bash
docker compose exec app migrate -path ./schema -database 'postgres://todo_user:todo_password@postgres:5432/postgres?sslmode=disable' up
```
## Работа с БД и запросы
![image](https://github.com/user-attachments/assets/1aedde5c-c642-4ea4-b5f8-ac653b5358a6)

### Аутентификация:
![Screenshot_1](https://github.com/user-attachments/assets/08af6656-0ea6-4a05-8239-0e2a4a2fb891)


### Post:
![image](https://github.com/user-attachments/assets/3cbde40a-ec31-423c-a730-9a35fe9c2f7f)


### Get:
![image](https://github.com/user-attachments/assets/ad526445-cd7e-4864-8e65-713d4c26426a)


### Put:
![Screenshot_4](https://github.com/user-attachments/assets/022277bf-dad3-4f6e-ba36-b201ac67e6dc)


### Delete:
![Screenshot_5](https://github.com/user-attachments/assets/c440febe-1a83-4e9f-ab67-1038936d9f00)

### Логи в MongoDB
![Screenshot_6](https://github.com/user-attachments/assets/32b60280-a6ab-4169-b475-dc93b12ec925)

### Кэш в Redis 
![image](https://github.com/user-attachments/assets/08ee8300-c304-4e88-8133-4467cb76ad49)
