# Сервис мерч магазина

## Содержание


- [Введение](#введение)
- [Начало работы](#начало-работы)
    - [Установка](#установка)
    - [Использование](#использование)

## Введение
В рамках задания требовалось разработать инструмент, который обеспечит необходимый функционал для внутреннего мерч магазина.

Были выполнены все основные задания.

Использованы следующие технологии:
- Golang
- PostgreSQL (в качестве хранилища данных)
- pq (драйвер для работы с PostgreSQL)
- cleanenv (для настройки конфигурации)
- bcrypto (для хеширования паролей в базу данных)
- jwt (для выдачи и валидации токенов)
- mockery и testify (для тестирования уровня бизнес логики)
- Docker и Docker compose (для запуска сервиса и e2e тестирования)

## Начало работы

### Установка

Перед началом, убедитесь, что у вас установлен `Docker` и `Docker compose`

Для запуска сервиса выполните следующие шаги:

1. Склонируйте этот репозиторий.
2. Перейдите в каталог проекта.
3. При необходимости, отредактируйте конфигурационный файл config.yaml или docker-compose.yaml(environment).
4. Для запуска сервиса выполните:

    ```bash
    make docker_up
    ```

5. Для тестирования бизнес логики выполните:

    ```bash
    make test
    ```

6. Для E2E тестирования аутентификации, покупки итема и отправки монеток выполните:

    ```bash
    make e2e
    ```

    Тест запускается в отдельных докер-контейнерах с отличающимися портами.

### Использование

Чтобы взаимодействовать с сервисом, вы можете использовать различные API-эндпоинты, согласно документации API [`schema.yaml`](schema.yaml)