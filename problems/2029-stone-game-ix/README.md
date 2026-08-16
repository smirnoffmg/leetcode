# 2029. Stone Game IX

- Difficulty: medium
- Link: https://leetcode.com/problems/stone-game-ix/
- Topics: Array, Math, Greedy, Minimax, Counting, Game Theory, Nim Game, Zero-Sum Game
- Acceptance: 51.1%

## Statement

Alice and Bob continue their games with stones. There is a row of n stones, and each stone has an associated value. You are given an integer array `stones`, where `stones[i]` is the **value** of the `i^th` stone.

Alice and Bob take turns, with **Alice** starting first. On each turn, the player may remove any stone from `stones`. The player who removes a stone **loses** if the **sum** of the values of **all removed stones** is divisible by `3`. Bob will win automatically if there are no remaining stones (even if it is Alice's turn).

Assuming both players play **optimally**, return `true` _if Alice wins and_ `false` _if Bob wins_.

**Example 1:**

```
Input: stones = [2,1]
Output: true
Explanation: The game will be played as follows:
- Turn 1: Alice can remove either stone.
- Turn 2: Bob removes the remaining stone.
The sum of the removed stones is 1 + 2 = 3 and is divisible by 3. Therefore, Bob loses and Alice wins the game.
```

**Example 2:**

```
Input: stones = [2]
Output: false
Explanation: Alice will remove the only stone, and the sum of the values on the removed stones is 2.
Since all the stones are removed and the sum of values is not divisible by 3, Bob wins the game.
```

**Example 3:**

```
Input: stones = [5,1,2,4,3]
Output: false
Explanation: Bob will always win. One possible way for Bob to win is shown below:
- Turn 1: Alice can remove the second stone with value 1. Sum of removed stones = 1.
- Turn 2: Bob removes the fifth stone with value 3. Sum of removed stones = 1 + 3 = 4.
- Turn 3: Alices removes the fourth stone with value 4. Sum of removed stones = 1 + 3 + 4 = 8.
- Turn 4: Bob removes the third stone with value 2. Sum of removed stones = 1 + 3 + 4 + 2 = 10.
- Turn 5: Alice removes the first stone with value 5. Sum of removed stones = 1 + 3 + 4 + 2 + 5 = 15.
Alice loses the game because the sum of the removed stones (15) is divisible by 3. Bob wins the game.
```

**Constraints:**

- `1 <= stones.length <= 10^5`
- `1 <= stones[i] <= 10^4`

## Hints

<details><summary>Hint 1</summary>

There are limited outcomes given the current sum and the stones remaining.

</details>

<details><summary>Hint 2</summary>

Can we greedily simulate starting with taking a stone with remainder 1 or 2 divided by 3?

</details>

## Similar

- [Stone Game](https://leetcode.com/problems/stone-game/) (medium)
- [Stone Game II](https://leetcode.com/problems/stone-game-ii/) (medium)
- [Stone Game III](https://leetcode.com/problems/stone-game-iii/) (hard)
- [Stone Game IV](https://leetcode.com/problems/stone-game-iv/) (hard)
- [Stone Game V](https://leetcode.com/problems/stone-game-v/) (hard)
- [Stone Game VI](https://leetcode.com/problems/stone-game-vi/) (medium)
- [Stone Game VII](https://leetcode.com/problems/stone-game-vii/) (medium)
- [Stone Game VIII](https://leetcode.com/problems/stone-game-viii/) (hard)
- [Stone Game IX](https://leetcode.com/problems/stone-game-ix/) (medium)

## Idea

Проигрыш определяется делимостью суммы на 3, поэтому от камня важен только остаток
его стоимости по модулю 3 — вся игра сжимается в три счётчика `cnt[0..2]`.

Камень с остатком 0 не меняет сумму: это «пас», передающий ход сопернику. На пас
всегда можно ответить пасом, поэтому пары нулей взаимно уничтожаются — важна лишь
чётность `cnt[0]`.

Без нулей безопасная партия жёстко предопределена первым камнем: после остатка `r`
сумма равна `r`, безопасно добавить только `r`, а дальше остатки обязаны чередоваться
— `r, r, r̄, r, r̄, …`, любое отклонение делает сумму кратной 3. Значит, исходов два:
у ходящего нет камня нужного остатка, но камни ещё есть — он вынужден отклониться и
проигрывает; или камни закончились — тогда по условию выигрывает Боб.

**`cnt[0]` чётно.** Если какой-то из остатков 1/2 отсутствует, выиграть Алиса не
может: цепочка `r, r, r, …` обрывается либо концом камней (Боб), либо третьим камнем,
делающим сумму кратной 3, — а первый и третий ходы всегда её. Если есть оба остатка,
Алиса начинает с меньшинства — остатка `r` с `cnt[r] <= cnt[r̄]`. В цепочке
`r, r, r̄, r, r̄, …` чётные позиции (все ходы Боба) требуют именно `r`, так что камни
`r` кончаются на ходе Боба, пока камни `r̄` ещё остались, — отклоняться приходится
ему. Итого: Алиса выигрывает ⇔ `cnt[1] > 0 && cnt[2] > 0`.

**`cnt[0]` нечётно.** Один пас остаётся неспаренным, и зажатый игрок ровно один раз
уворачивается, пасуя вместо рокового хода, — сжатие из чётного случая перестаёт
работать. Алисе остаётся выигрыш «на истощение»: она начинает с большинства `r`,
цепочка съедает все `cnt[r̄]` камней меньшинства и `cnt[r̄] + 2` камня большинства, и
чтобы после потраченного паса зажатым снова оказался Боб, запас должен быть строгим:
`cnt[r] > cnt[r̄] + 2`, то есть `|cnt[1] - cnt[2]| > 2`.

Разбор случаев легко испортить ошибкой на единицу, поэтому в тестах рядом лежит
эталонный минимакс от определения игры (состояние — счётчики, сумма mod 3, чей ход)
и сверка на 2000 случайных входов.

- Time: O(n) — один проход подсчёта остатков
- Space: O(1) — три счётчика
