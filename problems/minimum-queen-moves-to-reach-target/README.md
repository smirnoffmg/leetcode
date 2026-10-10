# Minimum Queen Moves to Reach Target

- Difficulty: easy
- Link: https://leetcode.com/problems/minimum-queen-moves-to-reach-target/
- Contest: Q1, номер задачи на момент решения ещё не опубликован

## Statement

There is an `8 x 8` empty chessboard with **1-indexed** rows and columns.

You are given an array `source = [sr, sc]` representing the starting position of a **queen**, and an array
`target = [tr, tc]` representing the target position.

In one move, the queen travels one or more squares along a single **row**, **column**, or **diagonal**, staying within
the board.

Return the **minimum** number of moves for the queen to land **exactly** on `target`.

**Example 1:**

**Input:** source = [8,1], target = [1,8]

**Output:** 1

**Example 2:**

**Input:** source = [4,2], target = [1,3]

**Output:** 2

**Explanation:**

The queen moves from `(4, 2)` to `(4, 3)`, then from `(4, 3)` to `(1, 3)`, reaching the target in 2 moves.

## Idea

Ответ никогда не больше двух, поэтому поиск не нужен — достаточно отличить 0, 1 и 2.

- **0** — клетки совпадают.
- **1** — клетки на одной линии ферзя. Строка и столбец сравниваются напрямую, диагонали —
  через инварианты: вдоль диагонали `\` постоянна разность `r − c`, вдоль `/` — сумма `r + c`.
- **2** — всё остальное. Верхняя граница конструктивна: сначала по своему столбцу в клетку
  `(tr, sc)`, затем по строке `tr` в цель. Доска пустая, мешать нечему. Нижняя граница — от
  противного: один ход означал бы общую линию, а её нет.

Наивное решение — BFS по 64 клеткам — верно, но прячет главное наблюдение: у ферзя на пустой
доске диаметр графа ходов равен двум. Тест сверяет формулу именно с таким BFS на всех
64 × 64 парах.

- Time: O(1)
- Space: O(1)
