# Transform Array Using Pair Operations

- Difficulty: medium
- Link: https://leetcode.com/problems/transform-array-using-pair-operations/
- Contest: Q2, номер задачи на момент решения ещё не опубликован

## Statement

You are given two integer arrays `source` and `target`.

In one **operation**, you may choose two **distinct** indices `i` and `j` in `source`, along with any integer `delta`.
Then update `source` as follows:

- `source[i] = source[i] + source[j] - delta`
- `source[j] = delta`

Return `true` if it is possible to make `source` equal to `target` after performing the operation **any** (including
zero) number of times. Otherwise, return `false`.

**Example 1:**

**Input:** source = [1,2,3], target = [0,2,4]

**Output:** true

**Example 2:**

**Input:** source = [-5,-5], target = [-15,5]

**Output:** true

**Example 3:**

**Input:** source = [1,2,1], target = [0,2,5]

**Output:** false

**Constraints:**

- `2 <= source.length == target.length <= 10^5`
- `-10^9 <= source[i], target[i] <= 10^9`

## Idea

Операция — это перераспределение внутри пары: `source[i] + source[j]` до и после одинакова,
остальные элементы не трогаются. Значит, **сумма массива — инвариант**, и при разных суммах
ответ `false`.

Главное — что инвариант **полный**: равенства сумм уже достаточно. Доказательство
конструктивное. Назначим последний индекс «копилкой» и для каждого `j < n − 1` сделаем
операцию `i = n − 1`, `delta = target[j]`. Элемент `j` сразу становится нужным, а избыток
или недостача стекает в копилку. Когда все `j` выставлены, в копилке остаётся
`сумма − Σ target[j]`, а это ровно `target[n − 1]`. Нужно `n ≥ 2` — ограничения это
гарантируют.

Разность сумм копится в одном цикле. Переполнения нет: `|Σ| ≤ 10^5 · 2·10^9 = 2·10^14`.

Тест проверяет достаточность не формулой, а буквально: применяет операции из
конструкции и сравнивает результат с `target`.

- Time: O(n)
- Space: O(1)
