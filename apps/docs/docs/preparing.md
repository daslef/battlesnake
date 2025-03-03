---
title: Подготовка к турниру
description: Коллекция советов и полезных алгоритмов для подготовки к финалу.
draft: false
---

## Рекомендуемые алгоритмы

### Flood Fill

Flood fill это алгоритм для определения границ и размера закрытых пространств. Известный пример применения алгоритма - инструмент 'заливка' в графических программах. Алгоритм помогает избегать ограниченных пространств, из которых будет сложно или невозможно выбраться. Например, можно определить опасность движения Червячка в некотором направлении, статистически подсчитав, сколько ходов он сможет сделать далее, пока не утратит возможность движения.

- [Flood Fill - Wikipedia](https://en.wikipedia.org/wiki/Flood_fill)

### A\* Pathfinding

A\* алгоритм для построения кратчайшего пути между двумя точками сетки. Он поможет при поиске пути до ближайшей еды или ближайшего соперника для его атаки. 

- [An intro to A\* pathfinding in the context of games](https://www.redblobgames.com/pathfinding/a-star/introduction.html)
- [A JavaScript library implementing A\* and some other pathfinding algorithms](https://github.com/qiao/PathFinding.js/)
- [An older Battlesnake project implementing A\* pathfinding](https://github.com/sockbot/battlesnake-api)

### Предостережение

Опытные разработчики могут задуматься об использовании алгоритмов *классического машинного обучения*, *обучения с подкреплением* или *нейронных сетей*. Однако, чтобы наше состязание было справедливым и шанс на победу получили все участники, **их использование в этом сезоне запрещено**. Предлагаем вместо этого сделать упор на стратегию и креативность.

## Бренд и внешний вид

### Запоминающееся имя

Вряд ли червячки с именами _"test"_ или _"snake-1"_ запомнятся надолго. Можно использовать название команды, а если участвуешь в одиночку - придумать корректное имя, которое притом было бы забавным и креативным.

### Кастомизация

Ты можешь настроить внешний вид Червячка. Обрати внимание на [доступные головы и хвосты](https://play.battlesnake.com/customizations), а также подбери запоминающийся цвет, например, используя инструмент [ColorPicker](https://colorpicker.me/).

```json title="customizations.json"
{
	"color": "#736CCB",
	"head": "beluga",
	"tail": "curled"
}
```

![Sample Snake](/img/samplesnake.png)
