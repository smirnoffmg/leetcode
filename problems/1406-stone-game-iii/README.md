# 1406. Stone Game III

- Difficulty: hard
- Link: https://leetcode.com/problems/stone-game-iii/
- Topics: Array, Math, Dynamic Programming, Game Theory
- Acceptance: 64.6%

## Statement

Alice and Bob continue their games with piles of stones. There are several stones **arranged in a row**, and each stone has an associated value which is an integer given in the array `stoneValue`.

Alice and Bob take turns, with Alice starting first. On each player's turn, that player can take `1`, `2`, or `3` stones from the **first** remaining stones in the row.

The score of each player is the sum of the values of the stones taken. The score of each player is `0` initially.

The objective of the game is to end with the highest score, and the winner is the player with the highest score and there could be a tie. The game continues until all the stones have been taken.

Assume Alice and Bob **play optimally**.

Return `"Alice"`_ if Alice will win,_ `"Bob"`_ if Bob will win, or_ `"Tie"`_ if they will end the game with the same score_.

**Example 1:**

```
Input: stoneValue = [1,2,3,7]
Output: "Bob"
Explanation: Alice will always lose. Her best move will be to take three piles and the score become 6. Now the score of Bob is 7 and Bob wins.
```

**Example 2:**

```
Input: stoneValue = [1,2,3,-9]
Output: "Alice"
Explanation: Alice must choose all the three piles at the first move to win and leave Bob with negative score.
If Alice chooses one pile her score will be 1 and the next move Bob's score becomes 5. In the next move, Alice will take the pile with value = -9 and lose.
If Alice chooses two piles her score will be 3 and the next move Bob's score becomes 3. In the next move, Alice will take the pile with value = -9 and also lose.
Remember that both play optimally so here Alice will choose the scenario that makes her win.
```

**Example 3:**

```
Input: stoneValue = [1,2,3,6]
Output: "Tie"
Explanation: Alice cannot win this game. She can end the game in a draw if she decided to choose all the first three piles, otherwise she will lose.
```

**Constraints:**

- `1 <= stoneValue.length <= 5 * 10^4`
- `-1000 <= stoneValue[i] <= 1000`

## Hints

<details><summary>Hint 1</summary>

The game can be mapped to minmax game. Alice tries to maximize the total score and Bob tries to minimize it.

</details>

<details><summary>Hint 2</summary>

Use dynamic programming to simulate the game. If the total score was 0 the game is "Tie", and if it has positive value then "Alice" wins, otherwise "Bob" wins.

</details>

## Similar

- [Stone Game V](https://leetcode.com/problems/stone-game-v/) (hard)
- [Stone Game VI](https://leetcode.com/problems/stone-game-vi/) (medium)
- [Stone Game VII](https://leetcode.com/problems/stone-game-vii/) (medium)
- [Stone Game VIII](https://leetcode.com/problems/stone-game-viii/) (hard)
- [Stone Game IX](https://leetcode.com/problems/stone-game-ix/) (medium)

## Idea

### Постановка

Игра конечная, с полной информацией и нулевой суммой: сумма очков обоих игроков равна `sum(stoneValue)` и не зависит от хода игры, поэтому выигрыш одного ровно равен проигрышу другого. Значение позиции вычисляется обратной индукцией от терминальной позиции.

Камни берутся только с начала ряда, поэтому любая достижимая позиция — суффикс `stoneValue[i:]`. Различных позиций `n + 1`, а не `3ⁿ`: разные последовательности ходов, снявшие одинаковое число камней, приводят к одной и той же позиции.

Отличие от [486. Predict the Winner](https://leetcode.com/problems/predict-the-winner/) и [877. Stone Game](https://leetcode.com/problems/stone-game/) — позиция описывается одним индексом вместо пары, а ветвление три вместо двух. Второе отличие существеннее: значения могут быть отрицательными, поэтому «взять больше» не обязательно лучше, и ничья возможна.

### Функция ценности

Определим для суффикса `stoneValue[i:]`

```
f(i) = максимум по стратегиям игрока, делающего ход,
       величины (его очки − очки соперника) на этом суффиксе
       при условии оптимальной игры соперника
```

Определение не содержит номера игрока: правила симметричны, и позиция характеризуется только индексом `i`. Это и есть причина, по которой в состоянии не требуется хранить ни накопленные очки, ни чья очередь хода.

### Рекуррентное соотношение

Игрок, делающий ход на `[i:]`, имеет не более трёх альтернатив. Взяв `k` камней (`k = 1, 2, 3`), он получает `stoneValue[i] + … + stoneValue[i+k−1]` и оставляет сопернику позицию `[i+k:]`, значение которой `f(i+k)` есть перевес соперника; относительно ходившего это та же величина с обратным знаком. Отсюда

```
f(n) = 0
f(i) = max( (stoneValue[i] + … + stoneValue[i+k−1]) − f(i+k) ),  k = 1..3, i+k <= n
```

Смена знака при переходе хода — стандартный приём negamax: минимизирующий игрок не выделяется в отдельную ветвь, поскольку `min(x) = −max(−x)`.

Правая часть зависит только от бо́льших индексов, поэтому соотношение корректно определяет `f` индукцией по `n − i`.

**Критерий ответа.** `f(0) = S_Alice − S_Bob`, поэтому `f(0) > 0` → `"Alice"`, `f(0) < 0` → `"Bob"`, `f(0) = 0` → `"Tie"`.

### Почему жадность неверна

Правило «взять максимум доступного» оценивает только немедленный выигрыш и игнорирует то, что ход открывает сопернику. Контрпример — `stoneValue = [7,1,1,9,1,1]` (тест `greedy trap`): жадный ход `7+1+1 = 9` оставляет сопернику позицию ценой `11` и даёт разницу `−2`, тогда как лучший ход `7+1` даёт `−1`. Пример 2 из условия показывает обратную крайность: там оптимально забрать все три камня, чтобы не отдать выбор сопернику.

### Проверка на примере 1

`stoneValue = [1,2,3,7]`:

| суффикс     | вычисление                  | `f`    |
| ----------- | --------------------------- | ------ |
| `[]`        | база                        | 0      |
| `[7]`       | `7−0`                       | 7      |
| `[3,7]`     | `max(3−7, 10−0)`            | 10     |
| `[2,3,7]`   | `max(2−10, 5−7, 12−0)`      | 12     |
| `[1,2,3,7]` | `max(1−12, 3−10, 6−7)`      | **−1** |

`−1 < 0` → `"Bob"`. Сумма всех камней `13`, значит счёт `6:7` — что согласуется с разбором в условии. Последняя строка показывает механику выбора: все три хода Алисы проигрышны, она лишь минимизирует отставание, забирая `1+2+3`.

### Реализация

Значения заполняются по убыванию `i`, накопленная сумма `taken` наращивается внутренним циклом, поэтому префиксные суммы не нужны. Внутренний цикл ограничен тремя итерациями, так что общее время линейно.

База `f(n) = 0` задаётся нулевым хвостом массива `dp` длины `n + 1`, а инициализация `best` случаем `k = 1` избавляет от `math.MinInt`: при `i < n` хотя бы один ход всегда доступен.

Случай `n = 1` обрабатывается без ветвления: единственная итерация даёт `dp[0] = stoneValue[0]`.

### Замечания

`f(i)` зависит только от `f(i+1)`, `f(i+2)`, `f(i+3)`, поэтому массив сворачивается в три переменные и память становится O(1). Здесь это оставлено как есть: при `n <= 5·10⁴` массив стоит 400 КБ, а версия с ротацией трёх переменных заметно менее читаема.

Эквивалентная формулировка — нисходящая рекурсия с мемоизацией по `i`; она ближе к определению `f`, но при `n = 5·10⁴` требует внимания к глубине стека. При `n <= 12` проходит и перебор без мемоизации; в тестах он используется как независимый эталон для сверки.

Ряд из единиц даёт периодическую по `n` последовательность `f`: `1, 2, 3, 2, 1, 0` с периодом `6`, то есть ничья ровно при `n`, кратном `6` (тесты `six ones`, `seven ones`, `max length`).

Тот же шаблон `max(взятое − ценность остатка для соперника)` решает семейство задач Stone Game; отличается только описание позиции.

- Time: O(n)
- Space: O(n)
