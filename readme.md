# ⚙️ ATMC — язык конфигурации

ATMC — минималистичный, но мощный язык конфигурации, созданный для удобного описания системных настроек.  
Он сочетает простоту синтаксиса, поддержку импортов, переопределений и гибкое объединение конфигов.

## 📚 Оглавление

- [🚀 Возможности](#возможности)
- [💻 Примеры использования](#примеры-использования)
    - [🧠 Все типы](#все-типы)
    - [📁 Импорты и переменные](#импорты-и-переменные)
    - [🧩 Spread (встраивание)](#встраивание-spread)
    - [🔁 Переопределение](#переопределение)
    - [🧠 Переопределение + рекурсивное слияние](#переопределение--рекурсивное-слияние)
- [🧬 Архитектура ATMC](#архитектура-atmc)

## 🚀 Возможности

- 📦 Импорт конфигов
- 🔀 Слияние (merge) конфигов
- 🔗 Доступ к импортированным данным
- 🧩 Переопределение и расширение полей
- 🌀 Встраивание объектов и массивов (spread)
- 🌿 Поддержка переменных среды (`env`)
- 💬 Поддержка однострочных комментариев
- 🧠 Удобный минималистичный синтаксис:
    - без лишних кавычек, как в JSON
    - со скобками вместо отступов (в отличие от YAML)
    - допускает свободное форматирование
    - запятые опциональны
- 🧾 Поддержка базовых типов:
    - int
    - string
    - float
    - bool
    - array
    - object

## 💻 Пример использования

### Все типы

📄File: `config.atmc`

```atmc
{
  db: {
    postgres: {
      database: "postgres"
      port: 5432
      username: $POSTGRES_USERNAME
      password: $POSTGRES_PASSWORD
    }
    clickhouse: {
      database: "clickhouse"
      port: 6000
      username: $CLICKHOUSE_USERNAME
      password: $CLICKHOUSE_PASSWORD
    } 
  }
  logging: {
    level: ["info", "warn", "error"] 
  }
  outbox: {
    enabled: true
    worker_count: 10  
  }
  tariffs: [
    {id: 1 price: 100.005 description: "fiest tariff"}
    {id: 2 price: 200.659 description: "second tariff"}
  ]
}
```

### Импорты и переменные

📄File: `config.atmc`

```atmc
postgres ./postgres.atmc

{
  db: {
    postgres: postgres.credentials // from imported config
  }
  logging: {
    level: ["info", "warn", "error"] 
  }
}
```

File: postgres.atmc

```atmc
{
  credentials: {
    database: "postgres"
    port:     5432
    username: $POSTGRES_USERNAME
    password: $POSTGRES_PASSWORD
  }

  settings: {
    max_connections: 100
    isolation:       "repeatable read"
    encoding:        "UTF8"
  }
}
```

Output

```atmc
{
  db: {
    postgres: {
      database: "postgres"
      port:     5432
      username: $POSTGRES_USERNAME
      password: $POSTGRES_PASSWORD
    }
  }
  logging: {
    level: ["info", "warn", "error"] 
  }
}
```

### Встраивание (spread)

📄File: `config.atmc`

```atmc
postgres ./postgres.atmc

{
  db: {
    postgres: {
      postgres.credentials... // spread operator
      log_level: "error"
    }
    clickhouse: {
      port: 6000
    }
  }
  logging: {
    level: ["info", "warn", "error"] 
  }
}
```

📄File: `postgres.atmc`

```atmc
{
  credentials: {
    database: "postgres"
    port:     5432
    username: $POSTGRES_USERNAME
    password: $POSTGRES_PASSWORD
  }
  
  settings: {
    max_connections: 100
    isolation:       "repeatable read"
    encoding:        "UTF8"
  }
}
```

Output

```atmc
{
  db: {
    postgres: {
      database: "postgres"
      port:     5432
      username: "username"
      password: "password"
      log_level: "error"
    }
    clickhouse: {
      port: 6000
    }
  }
  logging: {
    level: ["info", "warn", "error"] 
  }
}
```

### Переопределение

📄File: `prod.atmc`

```atmc
common ./common.atmc

{
  common... // встраиваем общий конфиг
  log_level: ["warn", "error"] 
}
```

📄File: `stage.atmc`

```atmc
common ./common.atmc

{
  common... // встраиваем общий конфиг
  log_level: ["info", "warn", "error"] 
}
```

📄File: `common.atmc`

```atmc
{
  outbox: {
    enabled: true
    worker_count: 10  
  }
  log_level: ["error"] 
}
```

Output Prod

```atmc
{
  outbox: {
    enabled: true
    worker_count: 10  
  }
  log_level: ["warn", "error"] // переопределено 
}
```

Output Stage

```atmc
{
  outbox: {
    enabled: true
    worker_count: 10  
  }
  log_level: ["info", "warn", "error"] // переопределено
}
```

### Переопределение + рекурсивное слияние

📄File: `prod.atmc`

```atmc
common ./common.atmc

{
  common... // встраиваем общий конфиг
  logging: {
    level: ["warn", "error"]
    enabled_tracing: true
  }
}
```

📄File: `stage.atmc`

```atmc
common ./common.atmc

{
  common... // встраиваем общий конфиг
  logging: {
    level: ["info", "warn", "error"] 
  }
}
```

📄File: `common.atmc`

```atmc
{
  outbox: {
    enabled: true
    worker_count: 10  
  }
  logging: {
    enabled: true
    level: ["error"] 
  }
}
```

Output Prod

```atmc
{
  // получено из общего конфига без изменений
  outbox: {
    enabled: true
    worker_count: 10  
  }
  logging: {
    enabled: true // вложенное поле получено из общего конфига без изменений
    level: ["warn", "error"] // переопределено
    enabled_tracing: true // добавлено в prod конифге
  }
}
```

Output Stage

```atmc
{
  // получено из общего конфига без изменений
  outbox: {
    enabled: true
    worker_count: 10  
  }
  logging: {
    enabled: true // вложенное поле получено из общего конфига без изменений
    level: ["info", "warn", "error"] // переопределено
  }
}
```

## 🚀 Под капотом

Язык конфигурации работает благодаря множеству компонентов:

- lexer
    - преобразует код в токены
- parser
    - преобразует токены в AST
- analyzer
    - проверяет семантику в рамках одного AST
    - проверяет наличие неиспользованных переменных
    - проверяет использование неопределенных переменных
- linker
    - резолвит значения переменных из всех связанных AST
    - резолвит значения переменных среды
    - выдает один итоговый AST
- processor
    - получает на вход путь до файла с конфигом
    - запускает все необходимые компоненты для обработки всех связанных файлов
    - отдает полученный итоговый AST после линковки
- compiler
    - компилирует итоговый AST
    - компиляторы могут быть разными (в map, в struct)

