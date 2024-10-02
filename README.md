GoNews API Gateway
==================
Точка входа в бэкенд агрегатора новостей.

## Требования

-   golang 1.22

-   docker >=23.0

# Начало работы

Просмотр конфигурации по-умолчанию:

        $ make build
        $ ./bin/go-news-api-gw -print-config

---

Для быстрого запуска с конфигом по-умолчанию:

        $ make run

Сервер будет запущен на `127.0.0.1:8080`. 

---

Запуск сервера с конфигурацией из файла `config.yaml`:

        $ ./bin/go-news-api-gw -config config.yaml

---

Остановить сервер:

        $ make clean

---

Показать версию сборки:

        $ ./bin/go-news-api-gw -version

## Docker

