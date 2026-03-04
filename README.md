# Telegram ToDo Bot (Go, in-memory MVP)

Telegram-бот на Go с архитектурой `transport -> service -> repository`:
- in-memory хранение задач и FSM-состояния диалога
- потокобезопасность через `sync.RWMutex`
- генерация `taskID` через `sync/atomic`

## Фичи

- Команды:
  - `/start` — приветствие + меню
  - `/help` — справка
- Reply-кнопки:
  - `➕ Add task`
  - `📋 List tasks`
  - `✅ Done tasks`
  - `🧹 Clear done`
- Inline-кнопки рядом с задачами:
  - `✅ done`
  - `🗑 delete`
- FSM-сценарий добавления задачи:
  - пользователь нажимает `Add task`
  - бот ждёт следующий текст
  - текст валидируется и сохраняется

```

## Структура проекта

```text
todo-bot/
  cmd/
    bot/
      main.go

  internal/
    transport/
      telegram/
        handler.go
        router.go
        keyboards.go
        render.go
    service/
      todo/
        service.go
        types.go
    repository/
      memory/
        task_repo.go
        state_repo.go
    model/
      task.go
    config/
      config.go
    app/
      app.go

  .env.example
  .gitignore
  go.mod
  README.md
```

## Важные детали реализации

- Хранилище задач: `map[userID]map[taskID]Task`
- Состояния диалога: `map[userID]State`
- Сервис валидирует текст:
  - trim
  - не пустой
  - минимум 2 символа
- Сервис сортирует список задач по `CreatedAt`
