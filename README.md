# Rust Server Updater

Утилита для автоматического обновления серверов Rust.

## Статус сборки

![Build Status](https://github.com/USERNAME/rust_updater/actions/workflows/build.yml/badge.svg)

## CI/CD Pipeline

В проекте настроен автоматический CI/CD pipeline с использованием GitHub Actions:

### Автоматическое создание тегов

При каждом пуше в ветку `main` автоматически создается новый тег версии и инициируется сборка.

### Сборка и Релиз

- Исходный код компилируется для Windows x64
- Проходят все тесты
- При создании тега автоматически публикуется новый релиз с исполняемым файлом

## Использование

1. Загрузите последний релиз из раздела [Releases](https://github.com/USERNAME/rust_updater/releases)
2. Запустите `rust_updater_windows_amd64.exe`

## Разработка

Для локальной сборки:

```bash
go build -v ./cmd/
```

## Настройка окружения для разработки

1. Go 1.20 или выше
2. Рекомендуется использовать IDE с поддержкой Go (GoLand, VS Code с расширением Go) 