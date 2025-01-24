---
sidebar_position: 4
pagination_next: null
pagination_prev: null
---

# Battlesnake

Описание Червячка (как твоего, так и соперников)

```json
{
  "id": "totally-unique-snake-id",
  "name": "Sneky McSnek Face",
  "health": 54,
  "body": [
    {"x": 0, "y": 0},
    {"x": 1, "y": 0},
    {"x": 2, "y": 0}
  ],
  "latency": "123",
  "head": {"x": 0, "y": 0},
  "length": 3,
  "shout": "why are we shouting??",
  "squad": "1",
  "customizations":{
    "color":"#26CF04",
    "head":"smile",
    "tail":"bolt"
  }
}
```

| **Свойство**       | **Тип** | **Описание**  | **Пример** |
| ------------------ | -------- | ---------- | ------------ |
| **id**             | string   | Уникальный идентификатор Червячка.  | "totally-unique-snake-id" |
| **name**           | string   | Имя Червячка, данное ему автором.   | "Sneky McSnek Face"     |
| **health**         | integer  | Здоровье, от 0 до 100.              | 54                      |
| **latency**        | string   | Время предыдущего ответа Червячка, либо величина таймаута (<code>game.timeout</code>), если он был превышен. | 423 |
| **head**           | object   | Координаты головы Червячка. Эквивалентны первому элементу массива body. |   |
| **length**         | integer  | Длина Червячка. Эквивалентна длине массива `body`. | 3 |
| **shout**          | string   | Высказывание Червячка на предыдущем ходу. | "why are we shouting??" |
| **customizations** | object   | Настройки внешнего вида Червячка.|  |
