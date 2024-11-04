# Полезные команды для работы с Docker, S3 и Docker Compose

## Команды Docker и MinIO
Вы можете выполнять следующие команды для управления docker-compose, docker и MinIO

```bash
# Войти в контейнер MinIO
docker exec -it minio sh

# Просмотреть содержимое бакета
mc ls local/bucket

# Показать содержание файла
mc cat local/bucket/4a440613-d285-43d0-a948-99e001f1677a

# Удалить все содержимое бакета
mc rm --recursive local/bucket --force

# Запустить Docker Compose в фоновом режиме
docker-compose up -d

# Остановить Docker Compose
docker-compose stop

# Остановить и удалить Docker Compose
docker-compose down
```

# Выполнение команд
Все команды выполняются из корня репозитория

# Сборка клиентов и сервера
Будут собраны клиенты под OS:
- darwin/amd64
- darwin/arm64
- linux/amd64
- windows/amd64
```bash
make build-all version=0.0.1 # значение версии должно быть без пробелов
```
