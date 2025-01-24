---
sidebar_position: 1
pagination_next: null
pagination_prev: null
---

# Game

Информация о карте, правилах игры, формате и уникальном идентификаторе.

```json
{
  "id": "totally-unique-game-id",
  "ruleset": {
    "name": "standard",
  },
  "map": "standard",
  "timeout": 500,
  "source": "league"
}
```


| **Свойство** | **Тип** | **Описание**    |
| ------------ | ------------------------ | ------ |
| **id**       | string                   | Уникальный идентификатор игры. |
| **ruleset**  | object                   | Информация о версии правил игры.    |
| **map**      | string                   | Название карты.      |
| **timeout**  | integer _(ms)_           | Допустимая задержка ответа.  |
| **source**   | string                   | Одно из значений: *tournament*, *league*, *arena*, *challenge*, *custom*. |
