---
title: Сетевое взаимодействие
pagination_next: null
pagination_prev: null
---

# Сетевое взаимодействие

## Общая информация

### Запросы

Движок отправляет `http-запросы` к веб-серверам игроков. Эти запросы могут содержать тело в формате [JSON](https://www.json.org/). 

### Ответы

Чтобы ответ сервера был воспринят движком как корректный, он должен соответствовать следующим правилам:

- содержимое: `JSON-encoded` строка,
- `Content-Type: application/json`,
- код ответа `HTTP 200 OK`,
- время отклика `500 мс` (включая транспортную задержку).

Если ответ воспринят движком как некорректный, **ход будет сделан за тебя**. Если это первый ход в партии, это будет движение *вверх*. Далее - будет повторяться предыдущий ход.

## `GET /`

Пустой GET запрос по корню веб-сервера. Служит для проверки соединения и получения настроек кастомизации.

### Схема запроса

Тело запроса отсутствует

### Схема ответа

| **Свойство**  | **Тип**            | **Описание**                                                                                                                                           |
| -------------- | ------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **author**     | string _(optional)_ | Имя автора (разработчика Червячка). <em>Example: "HubsyHubs"</em>                   |
| **color**      | string _(optional)_ | Шестнадцатеричный код цвета Червячка. <em>Example: "#888888"</em>                                |
| **head**       | string _(optional)_ | Отображение головы. <em>Example: "default"</em>                |
| **tail**       | string _(optional)_ | Отображение хвоста. <em>Example: "default"</em>                |

### Пример ответа

```json title="200 OK"
{
  "author": "MyUsername",
  "color": "#888888",
  "head": "default",
  "tail": "default",
}
```

## `POST /start`

Это запрос о старте игры, который рассылается по готовности всех участвующих Червячков. 

### Схема запроса

| **Параметр**                       | **Тип**  | **Описание**                                                                  |
| ---------------------------------- | -------- | ----------------------------------------------------------------------------- |
| **game**                           | object   | [Game Object](objects/game)                 |
| **turn**                           | integer  | Номер хода (0)                                     |
| **board**                          | object   | [Board Object](objects/board) |
| **you**                            | object   | [Battlesnake Object](objects/battlesnake)        |

### Схема ответа

Responses to this request are ignored by the game engine.

## `POST /move`

Этот запрос отправляется на каждом ходу. В теле запроса содержится информация о состоянии игрового поля и всех Червячков. 

### Схема запроса

| **Параметр**                      | **Тип** | **Описание**                                                               |
| ---------------------------------- | -------- | ----------------------------------------------------------------------------- |
| **game**                           | object   | [Game Object](objects/game)                 |
| **turn**                           | integer  | Номер хода                                                    |
| **board**                          | object   | [Board Object](objects/board) |
| **you**                            | object   | [Battlesnake Object](objects/battlesnake)        |


### Схема ответа

| **Свойство**  | **Тип**            | **Описание**                                                                                      |
| ------------- | ------------------- | ---------------------------------------------------------------------------------------------------- |
| **move**      | string              | Твой ход ("up", "down", "left" или "right")       |
| **shout**     | string _(optional)_ | Опциональное сообщение, которое будет отправлено всем соперникам (до 256 символов) |

### Пример запроса

```json title="POST /move"
{
  "game": {
    "id": "totally-unique-game-id",
    "ruleset": {
      "name": "standard",
      "settings": {
        "foodSpawnChance": 15,
      }
    },
    "map": "standard",
    "source": "league",
    "timeout": 500
  },
  "turn": 14,
  "board": {
    "height": 11,
    "width": 11,
    "food": [
      {"x": 5, "y": 5},
      {"x": 9, "y": 0},
      {"x": 2, "y": 6}
    ],
    "snakes": [
      {
        "id": "snake-508e96ac-94ad-11ea-bb37",
        "name": "My Snake",
        "health": 54,
        "body": [
          {"x": 0, "y": 0},
          {"x": 1, "y": 0},
          {"x": 2, "y": 0}
        ],
        "latency": "111",
        "head": {"x": 0, "y": 0},
        "length": 3,
        "shout": "why are we shouting??",
        "customizations":{
          "color":"#FF0000",
          "head":"pixel",
          "tail":"pixel"
        }
      },
      {
        "id": "snake-b67f4906-94ae-11ea-bb37",
        "name": "Another Snake",
        "health": 16,
        "body": [
          {"x": 5, "y": 4},
          {"x": 5, "y": 3},
          {"x": 6, "y": 3},
          {"x": 6, "y": 2}
        ],
        "latency": "222",
        "head": {"x": 5, "y": 4},
        "length": 4,
        "shout": "I'm not really sure...",
        "customizations":{
          "color":"#26CF04",
          "head":"silly",
          "tail":"curled"
        }
      }
    ]
  },
  "you": {
    "id": "snake-508e96ac-94ad-11ea-bb37",
    "name": "My Snake",
    "health": 54,
    "body": [
      {"x": 0, "y": 0},
      {"x": 1, "y": 0},
      {"x": 2, "y": 0}
    ],
    "latency": "111",
    "head": {"x": 0, "y": 0},
    "length": 3,
    "shout": "why are we shouting??",
    "customizations": {
      "color":"#FF0000",
      "head":"pixel",
      "tail":"pixel"
    }
  }
}
```


### Пример ответа

```json title="200 OK"
{
  "move": "up",
  "shout": "I guess I'll go up then."
}
```

## `POST /end`

Червячок получает этот запрос по завершению игры. В нем содержится информация об итоговом количестве ходов, состоянии игрового поля и результате.

### Схема запроса

| **Параметр**                      | **Тип** | **Описание**                                                               |
| ---------------------------------- | -------- | ----------------------------------------------------------------------------- |
| **game**                           | object   | [Game Object](objects/game)                 |
| **turn**                           | integer  | Количество ходов                                        |
| **board**                          | object   | [Board Object](objects/board) |
| **you**                            | object   | [Battlesnake Object](objects/battlesnake)        |

### Схема ответа

Тело ответа пустое.