# Приложение для тестового задания

Выполняет очень простой функционал: складывает два целых числа, вычитает из одного целого числа другое и выводит результаты

## Модели

### Модель вычисления

Модель создаётся для каждого вычисления и записывается в базу, в таблицу `calculations`.

Типы вычислений:
* `addition` - сложение
* `subtraction` - вычитание

Описание модели:

| Атрибут    | Тип данных | Описание                           |
|------------|------------|------------------------------------|
| id         | int        | Идентификатор данных в таблице     |
| created_at | datetime   | Дата и время проведения вычисления |
| a          | int        | Первое число для вычисления        |
| b          | int        | Второе число для вычисления        |
| result     | int        | Результат вычисления               |
| type       | string     | Тип вычисления                     |

## API

В приложении есть 3 метода API:
* addition - сложение
* subtraction - вычитание
* results - вывод списка результатов

### addition

Вызывается по пути `/addition`, метод `POST`.

Принимает JSON. Складывает число `a` и число `b`. Результат записывается в базу.

Запрос:

| Атрибут | Тип данных | Описание                    |
|---------|------------|-----------------------------|
| a       | int        | Первое число для вычисления |
| b       | int        | Второе число для вычисления |

Ответ: модель вычисления в виде JSON.

Пример запроса:
```json
{
    "a": 500,
    "b": 70
}
```

Пример ответа:
```json
{
    "id": 7,
    "created_at": "2024-02-16T14:08:17.343694737+03:00",
    "a": 500,
    "b": 70,
    "result": 570,
    "type": "addition"
}
```

### subtraction

Вызывается по пути `/subtraction`, метод `POST`.

Принимает JSON. Вычитает из числа `a` число `b`. Результат записывается в базу.

Запрос:

| Атрибут | Тип данных | Описание                    |
|---------|------------|-----------------------------|
| a       | int        | Первое число для вычисления |
| b       | int        | Второе число для вычисления |

Ответ: модель вычисления в виде JSON.

Пример запроса:
```json
{
    "a": 900,
    "b": 80
}
```

Пример ответа:
```json
{
    "id": 9,
    "created_at": "2024-02-16T14:08:50.624367176+03:00",
    "a": 900,
    "b": 80,
    "result": 820,
    "type": "subtraction"
}
```

### results

Вызывается по пути `/results`, метод `POST`.

Принимает JSON. Выводит результаты с указанным типом вычисления, ограниченные по лимиту и офсету.

Запрос: 

| Атрибут | Тип данных | Описание                 |
|---------|------------|--------------------------|
| type    | string     | Тип вычисления           |
| limit   | int        | Сколько записей выводить |
| offset  | int        | С какой записи выводить  |

Ответ:

| Атрибут | Тип данных | Описание                                    |
|---------|------------|---------------------------------------------|
| count   | int        | Общее количество записей с переданным типом |
| data    | array      | Массив моделей вычисления                   |

Пример запроса:
```json
{
    "type": "addition",
    "limit": 10,
    "offset": 0
}
```

Пример ответа:
```json
{
    "count": 5,
    "data": [
        {
            "id": 7,
            "created_at": "2024-02-16T14:08:17.343694+03:00",
            "a": 500,
            "b": 70,
            "result": 570,
            "type": "addition"
        },
        {
            "id": 6,
            "created_at": "2024-02-16T14:07:53.246089+03:00",
            "a": 500,
            "b": 70,
            "result": 570,
            "type": "addition"
        },
        {
            "id": 5,
            "created_at": "2024-02-16T13:55:45.854305+03:00",
            "a": 500,
            "b": 70,
            "result": 570,
            "type": "addition"
        },
        {
            "id": 4,
            "created_at": "2024-02-16T13:51:19.956744+03:00",
            "a": 100,
            "b": 3,
            "result": 103,
            "type": "addition"
        },
        {
            "id": 3,
            "created_at": "2024-02-16T13:40:25.259206+03:00",
            "a": 100,
            "b": 3,
            "result": 103,
            "type": "addition"
        }
    ]
}
```

## Логгирование

Логи форматируются в JSON. 

Пример лога:
```json
{"level":"WARNING","message":"can't get env variable \"DB_MAX_CONNECTIONS\" value","base_message":"strconv.Atoi: parsing \"\": invalid syntax","data":null,"trace":["/home/dshirokov/code/test/simple-app/internal/config/config.go:46"],"messages":["can't get env variable \"DB_MAX_CONNECTIONS\" value","strconv.Atoi: parsing \"\": invalid syntax"],"time":"2024-02-16 13:51:15 +03:00"}
```

И пример красивого лога:
```json
{
  "level": "WARNING",
  "message": "can't get env variable \"DB_MAX_CONNECTIONS\" value",
  "base_message": "strconv.Atoi: parsing \"\": invalid syntax",
  "data": null,
  "trace": [
    "/home/dshirokov/code/test/simple-app/internal/config/config.go:46"
  ],
  "messages": [
    "can't get env variable \"DB_MAX_CONNECTIONS\" value",
    "strconv.Atoi: parsing \"\": invalid syntax"
  ],
  "time": "2024-02-16 13:51:15 +03:00"
}
```

## Зависимости приложения

* PostgreSQL - СУБД для хранения результатов вычисления

## Переменные окружения

Требуются для работы приложения

| Переменная         | Описание                                                                    | Значение по умолчанию | Тип значения                                                                                                          |
|--------------------|-----------------------------------------------------------------------------|-----------------------|-----------------------------------------------------------------------------------------------------------------------|
| LOG_DEBUG          | Писать ли DEBUG логи: 0 - не писать, 1 - писать                             | 0                     | Целое число, 0 или 1                                                                                                  |
| DB_URL             | DSN подключения к PostgreSQL                                                | Нет                   | Валидный DSN подключения к PostgreSQL, например, `postgresql://developer:developer@main.database.res:5432/simple_app` |
| DB_MIN_CONNECTIONS | Минимальное количество соединений для пула соединений к каждой базе данных  | 1                     | Целое число больше 0                                                                                                  |
| DB_MAX_CONNECTIONS | Максимальное количество соединений для пула соединений к каждой базе данных | 4                     | Целое число больше 0 и больше DB_MIN_CONNECTIONS                                                                      |
| QUERY_TIMEOUT      | Максимальное время выполнения SQL запроса                                   | 15                    | Целое число больше 0                                                                                                  |
| HTTP_PORT          | HTTP порт, с которым будет работать сервер                                  | 80                    | Валидное значение порта                                                                                               |

