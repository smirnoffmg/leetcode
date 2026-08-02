# 877. Stone Game

- Difficulty: medium
- Link: https://leetcode.com/problems/stone-game/
- Topics: Array, Math, Dynamic Programming, Game Theory
- Acceptance: 74.6%

## Statement

Alice and Bob play a game with piles of stones. There are an **even** number of piles arranged in a row, and each pile
has a **positive** integer number of stones `piles[i]`.

The objective of the game is to end with the most stones. The **total** number of stones across all the piles is **odd**,
so there are no ties.

Alice and Bob take turns, with **Alice starting first**. Each turn, a player takes the entire pile of stones either from
the **beginning** or from the **end** of the row. This continues until there are no more piles left, at which point the person
with the **most stones wins**.

Assuming Alice and Bob play optimally, return `true`_ if Alice wins the game, or_ `false`_ if Bob wins_.

**Example 1:**

```
Input: piles = [5,3,4,5]
Output: true
Explanation:
Alice starts first, and can only take the first 5 or the last 5.
Say she takes the first 5, so that the row becomes [3, 4, 5].
If Bob takes 3, then the board is [4, 5], and Alice takes 5 to win with 10 points.
If Bob takes the last 5, then the board is [3, 4], and Alice takes 4 to win with 9 points.
This demonstrated that taking the first 5 was a winning move for Alice, so we return true.
```

**Example 2:**

```
Input: piles = [3,7,2,3]
Output: true
```

**Constraints:**

- `2 <= piles.length <= 500`
- `piles.length` is **even**.
- `1 <= piles[i] <= 500`
- `sum(piles[i])` is **odd**.

## Similar

- [Stone Game V](https://leetcode.com/problems/stone-game-v/) (hard)
- [Stone Game VI](https://leetcode.com/problems/stone-game-vi/) (medium)
- [Stone Game VII](https://leetcode.com/problems/stone-game-vii/) (medium)
- [Stone Game VIII](https://leetcode.com/problems/stone-game-viii/) (hard)
- [Stone Game IX](https://leetcode.com/problems/stone-game-ix/) (medium)
- [Strictly Palindromic Number](https://leetcode.com/problems/strictly-palindromic-number/) (medium)
- [Visit Array Positions to Maximize Score](https://leetcode.com/problems/visit-array-positions-to-maximize-score/) (medium)

## Idea

### Постановка

Игра конечная, с полной информацией и нулевой суммой: сумма очков обоих игроков равна `sum(piles)` и не зависит от хода игры. Ходы разрешены только с концов ряда, поэтому любая достижимая позиция — непрерывный отрезок `piles[i..j]`, и различных позиций `n(n+1)/2`, а не `2ⁿ`.

Это [486. Predict the Winner](https://leetcode.com/problems/predict-the-winner/) с двумя дополнительными ограничениями: длина чётна, а сумма нечётна. Первое из них делает задачу тривиальной (см. ниже), второе исключает ничью, поэтому критерий ответа — строгое `> 0`, а не `>= 0`.

### Функция ценности

Для отрезка `piles[i..j]`

```
f(i, j) = максимум по стратегиям игрока, делающего ход,
          величины (его очки − очки соперника) на этом отрезке
          при условии оптимальной игры соперника
```

Определение не содержит номера игрока: правила симметричны, позиция характеризуется только парой `(i, j)`. Поэтому в состоянии не нужно хранить ни накопленные очки, ни чью очередь ходить.

### Рекуррентное соотношение

Взяв `piles[i]`, ходящий оставляет сопернику позицию `[i+1..j]`, значение которой `f(i+1, j)` есть перевес соперника; относительно ходившего это та же величина с обратным знаком. Симметрично для `piles[j]`. Отсюда

```
f(i, i) = piles[i]
f(i, j) = max( piles[i] − f(i+1, j),  piles[j] − f(i, j−1) ),   i < j
```

Смена знака при переходе хода — negamax: минимизирующий игрок не выделяется в отдельную ветвь, поскольку `min(x) = −max(−x)`.

Ответ — `f(0, n−1) > 0`.

### Почему жадность неверна

Правило «брать больший конец» оценивает только немедленный выигрыш и игнорирует то, что ход открывает сопернику. Контрпример — `piles = [1,1,5,2]` (тест `greedy trap`): жадный первый ход `2` приводит к поражению 3:6, тогда как оптимальный ход `1` даёт 6:3.

### Проверка на примере 1

`piles = [5,3,4,5]`:

| отрезок     | вычисление      | `f`        |
| ----------- | --------------- | ---------- |
| одиночные   | база            | 5, 3, 4, 5 |
| `[5,3]`     | `max(5−3, 3−5)` | 2          |
| `[3,4]`     | `max(3−4, 4−3)` | 1          |
| `[4,5]`     | `max(4−5, 5−4)` | 1          |
| `[5,3,4]`   | `max(5−1, 4−2)` | 4          |
| `[3,4,5]`   | `max(3−1, 5−1)` | 4          |
| `[5,3,4,5]` | `max(5−4, 5−4)` | **1**      |

`1 > 0` → `true`; при оптимальной игре обоих счёт 9:8, что согласуется с разбором в условии.

### Реализация

Значения заполняются по возрастанию длины отрезка: внешний цикл идёт по убыванию `i`, внутренний — по возрастанию `j`, то есть `f(i, ·)` вычисляется после `f(i+1, ·)`.

Строка `i` зависит только от строки `i+1`, поэтому двумерная таблица сворачивается в один массив с инвариантом: перед присваиванием `dp[j]` хранит `f(i+1, j)` (значение с предыдущей итерации внешнего цикла), а `dp[j−1]` — уже `f(i, j−1)` (обновлён на текущей). Начальное состояние `dp = piles` соответствует базе `f(i, i)`.

### Наблюдение о чётности: ответ всегда true

Раскрасим кучи по чётности индекса. Пусть Алиса берёт `piles[0]`; тогда у Боба остаются края с индексами `1` и `n−1`, оба нечётные (`n` чётно), и он вынужден взять нечётный. После его хода у Алисы снова оба края чётные — и так до конца. Симметрично, взяв `piles[n−1]`, Алиса заберёт все нечётные индексы.

То есть Алиса может заранее выбрать любую из двух групп целиком, а нечётность суммы исключает их равенство. Значит она берёт бо́льшую группу и всегда выигрывает: `return true` за O(1).

Это наблюдение — артефакт ограничений, а не решение игры: при нечётной длине (тесты `odd length *`) или при чётной сумме (`tie *`) оно ломается, и нужен DP. Поэтому в решении оставлен DP, а чётностный аргумент вынесен в тест `TestStoneGameAliceAlwaysWinsUnderConstraints` как проверяемое свойство.

### Замечания

Эквивалентная формулировка — нисходящая рекурсия с мемоизацией по `(i, j)`; она ближе к определению `f` и не требует явно задавать порядок обхода. При `n <= 20` проходит перебор без мемоизации (`2ⁿ` вызовов); в тестах он используется как независимый эталон.

Тот же шаблон `max(взятое − ценность остатка для соперника)` решает семейство Stone Game; отличается только описание позиции.

- Time: O(n²)
- Space: O(n)
