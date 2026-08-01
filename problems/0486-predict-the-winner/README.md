# 486. Predict the Winner

- Difficulty: medium
- Link: https://leetcode.com/problems/predict-the-winner/
- Topics: Array, Math, Dynamic Programming, Recursion, Game Theory
- Acceptance: 56.9%

## Statement

You are given an integer array `nums`. Two players are playing a game with this array: player 1 and player 2.

Player 1 and player 2 take turns, with player 1 starting first. Both players start the game with a score of `0`.
At each turn, the player takes one of the numbers from either end of the array (i.e., `nums[0]` or `nums[nums.length - 1]`) which reduces the size of
the array by `1`. The player adds the chosen number to their score. The game ends when there are no more elements in the array.

Return `true` if Player 1 can win the game. If the scores of both players are equal, then player 1 is still the winner, and you should also return `true`.
You may assume that both players are playing optimally.

**Example 1:**

```
Input: nums = [1,5,2]
Output: false
Explanation: Initially, player 1 can choose between 1 and 2.
If he chooses 2 (or 1), then player 2 can choose from 1 (or 2) and 5. If player 2 chooses 5, then player 1 will be left with 1 (or 2).
So, final score of player 1 is 1 + 2 = 3, and player 2 is 5.
Hence, player 1 will never be the winner and you need to return false.
```

**Example 2:**

```
Input: nums = [1,5,233,7]
Output: true
Explanation: Player 1 first chooses 1. Then player 2 has to choose between 5 and 7. No matter which number player 2 choose, player 1 can choose 233.
Finally, player 1 has more score (234) than player 2 (12), so you need to return True representing player1 can win.
```

**Constraints:**

- `1 <= nums.length <= 20`
- `0 <= nums[i] <= 10^7`

## Similar

- [Can I Win](https://leetcode.com/problems/can-i-win/) (medium)
- [Find the Winning Player in Coin Game](https://leetcode.com/problems/find-the-winning-player-in-coin-game/) (easy)
- [Find the Number of Winning Players](https://leetcode.com/problems/find-the-number-of-winning-players/) (easy)
- [Count The Number of Winning Sequences](https://leetcode.com/problems/count-the-number-of-winning-sequences/) (hard)

## Idea

### Постановка

Игра конечная, с полной информацией и нулевой суммой: сумма очков обоих игроков равна `sum(nums)` и не зависит от хода игры, поэтому выигрыш одного ровно равен проигрышу другого. Для таких игр оптимальные стратегии обоих определены однозначно, а значение позиции вычисляется обратной индукцией от терминальных позиций.

Ходы разрешены только с концов массива, следовательно любая достижимая позиция — непрерывный отрезок `nums[i..j]`. Число различных позиций равно `n(n+1)/2`, а не `2^n`: последовательности ходов «слева, затем справа» и «справа, затем слева» приводят к одной и той же позиции.

### Почему жадность неверна

Правило «брать больший конец» оценивает только немедленный выигрыш и игнорирует то, что ход открывает сопернику. Контрпример — `nums = [2,4,55,6,8]`: жадная игра даёт 65:10 в пользу первого игрока, тогда как при оптимальной игре обоих первый проигрывает (см. тест `greedy trap`). Пример 2 из условия, `[1,5,233,7]`, демонстрирует ту же ошибку: оптимальный первый ход — взять `1`, то есть минимальный из доступных элементов.

### Функция ценности

Определим для отрезка `nums[i..j]`

```
f(i, j) = максимум по стратегиям игрока, делающего ход,
          величины (его очки − очки соперника) на этом отрезке
          при условии оптимальной игры соперника
```

Определение не содержит номера игрока: правила симметричны, и позиция характеризуется только парой `(i, j)`. Это и есть причина, по которой в состоянии не требуется хранить ни накопленные очки, ни чья очередь хода.

### Рекуррентное соотношение

Игрок, делающий ход на `[i..j]`, имеет ровно две альтернативы. Взяв `nums[i]`, он оставляет сопернику позицию `[i+1..j]`, значение которой `f(i+1, j)` есть перевес соперника; относительно ходившего это величина с обратным знаком. Симметрично для `nums[j]`. Отсюда

```
f(i, i) = nums[i]
f(i, j) = max( nums[i] − f(i+1, j),  nums[j] − f(i, j−1) ),   i < j
```

Смена знака при переходе хода — стандартный приём negamax: минимизирующий игрок не выделяется в отдельную ветвь, поскольку `min(x) = −max(−x)`.

Так как обе части правой стороны зависят от отрезков меньшей длины, соотношение корректно определяет `f` индукцией по `j − i`.

**Критерий ответа.** `f(0, n−1) = S₁ − S₂`, где `Sₖ` — итоговые очки игрока `k` при оптимальной игре. Условие «ничья засчитывается первому» даёт ответ `f(0, n−1) >= 0`.

### Проверка на примере 2

| отрезок                      | вычисление             | `f`          |
| ---------------------------- | ---------------------- | ------------ |
| `[1]`, `[5]`, `[233]`, `[7]` | база                   | 1, 5, 233, 7 |
| `[1,5]`                      | `max(1−5, 5−1)`        | 4            |
| `[5,233]`                    | `max(5−233, 233−5)`    | 228          |
| `[233,7]`                    | `max(233−7, 7−233)`    | 226          |
| `[1,5,233]`                  | `max(1−228, 233−4)`    | 229          |
| `[5,233,7]`                  | `max(5−226, 7−228)`    | −221         |
| `[1,5,233,7]`                | `max(1−(−221), 7−229)` | **222**      |

`222 >= 0` → `true`, что согласуется с разбором в условии: `234 − 12 = 222`. Последняя строка показывает механику выбора: взятие `7` уступает взятию `1`, потому что остаток `[1,5,233]` стоит для соперника `+229`, а остаток `[5,233,7]` — `−221`.

### Реализация

Значения заполняются по возрастанию длины отрезка. В коде это выражено убыванием `i` при внутреннем цикле по `j`, что эквивалентно: `f(i, ·)` вычисляется после `f(i+1, ·)`.

Строка `i` зависит только от строки `i+1`, поэтому двумерная таблица сворачивается в один массив с инвариантом: перед присваиванием `dp[j]` хранит `f(i+1, j)` (значение с предыдущей итерации внешнего цикла), а `dp[j−1]` — уже `f(i, j−1)` (обновлён на текущей). Начальное состояние `dp = nums` соответствует базе `f(i, i)`.

Случай `n = 1` обрабатывается без ветвления: внешний цикл не выполняется, `dp[0] = nums[0] >= 0` по ограничениям задачи.

### Замечания

Эквивалентная формулировка — нисходящая рекурсия с мемоизацией по `(i, j)`; она ближе к определению `f` и не требует явно задавать порядок обхода. При `n <= 20` проходит и перебор без мемоизации (`2ⁿ` вызовов); в тестах он используется как независимый эталон для сверки.

Тот же шаблон `max(взятое − ценность остатка для соперника)` решает семейство задач Stone Game; отличается только описание позиции.

- Time: O(n²)
- Space: O(n)
