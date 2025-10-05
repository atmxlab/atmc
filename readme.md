# ATMC - язык конфигурации

Язык конфигурации с мощными фичами

## 🚀 Фичи

- импорт конфигов
- merge конфигов
- доступ к импортированным данным
- переопределение полей
- встраивание объектов и массивов (spread)
- поддержка переменных среды (env)
- минималистичный простой синтаксис:
    - нет множества кавычек, как в json
    - для наглядности и удобства используются скобки вместо отступов, как в yml
    - свободное форматирование кода
    - запятые не обязательны - по желанию
- поддержка множества типов:
    - int
    - string
    - float
    - bool
    - array
    - object

## 💻 Пример использования

### 💻 Демонстрация всех типов

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

### 💻 Демонстрация импортов и переменных

File: config.atmc

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

### 💻 Демонстрация spread

File: config.atmc

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

### 💻 Демонстрация переопределения

File: prod.atmc

```atmc
common ./common.atmc

{
  common... // встраиваем общий конфиг
  logging: {
    level: ["warn", "error"] 
  }
}
```

File: stage.atmc

```atmc
common ./common.atmc

{
  common... // встраиваем общий конфиг
  logging: {
    level: ["info", "warn", "error"] 
  }
}
```

File: common.atmc

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
  outbox: {
    enabled: true
    worker_count: 10  
  }
  logging: {
    enabled: true
    level: ["warn", "error"] // переопределено
  }
}
```

Output Stage

```atmc
{
  outbox: {
    enabled: true
    worker_count: 10  
  }
  logging: {
    enabled: true
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

