# 2996. Smallest Missing Integer Greater Than Sequential Prefix Sum

- Difficulty: easy
- Link: https://leetcode.com/problems/smallest-missing-integer-greater-than-sequential-prefix-sum/
- Topics: Array, Hash Table, Sorting
- Acceptance: 43.5%

## Statement

You are given a **0-indexed** array of integers `nums`.

A prefix `nums[0..i]` is **sequential** if, for all `1 <= j <= i`, `nums[j] = nums[j - 1] + 1`. In particular, the
prefix consisting only of `nums[0]` is **sequential**.

Return _the **smallest** integer_ `x` _missing from_ `nums` _such that_ `x` _is greater than or equal to
the sum of the **longest** sequential prefix._

**Example 1:**

```
Input: nums = [1,2,3,2,5]
Output: 6
Explanation: The longest sequential prefix of nums is [1,2,3] with a sum of 6. 6 is not in the array,
therefore 6 is the smallest missing integer greater than or equal to the sum of the longest sequential prefix.
```

**Example 2:**

```
Input: nums = [3,4,5,1,12,14,13]
Output: 15
Explanation: The longest sequential prefix of nums is [3,4,5] with a sum of 12. 12, 13, and 14 belong to the array
while 15 does not. Therefore 15 is the smallest missing integer greater than or equal to the sum of
the longest sequential prefix.
```

**Constraints:**

- `1 <= nums.length <= 50`
- `1 <= nums[i] <= 50`

## Hints

<details><summary>Hint 1</summary>

To find the longest sequential prefix, iterate from left to right. For a fixed `i`, if `nums[i] != nums[i - 1] + 1`
then the longest sequential prefix ends at `i - 1`.

</details>

## Similar

- [Longest Common Prefix](https://leetcode.com/problems/longest-common-prefix/) (easy)
- [First Missing Positive](https://leetcode.com/problems/first-missing-positive/) (hard)
- [Next Greater Element I](https://leetcode.com/problems/next-greater-element-i/) (easy)

## Idea

Задача распадается на два независимых шага, и весь интерес — во втором.

Первый шаг тривиален: последовательный префикс — это ровно `nums[0]` плюс каждый следующий
элемент, пока держится `nums[i] == nums[i-1]+1`. На первом нарушении цикл обрывается;
искать более длинный последовательный участок дальше по массиву не нужно — по определению
префикс начинается с нуля.

Второй шаг — наименьшее отсутствующее значение, не меньшее суммы `prefix`. Наивно: для
каждого кандидата линейно искать его в массиве, O(n) на кандидата. Но домен значений задан
условием и плотный (`1 <= nums[i] <= 50`), поэтому занятость всех значений сразу
записывается в булев массив за один проход, и проверка кандидата становится O(1).

Тонкое место — верхняя граница поиска. Соблазн ограничить цикл размером массива занятости
и вернуть `-1` за ним — ошибка сразу в две стороны:

- сумма префикса не ограничена пятьюдесятью: для `[1..50]` она равна 1275, и ответ равен ей
  же, потому что ни одно значение массива до неё не дотягивается;
- даже при маленькой сумме ответ может выйти за домен: если заняты все значения 1..50,
  наименьшее свободное — 51.

Правильная граница берётся не из размера структуры, а из условия: 51 не может встретиться
в массиве, значит ответ никогда не превышает `max(prefix, 51)`. Отсюда массив на 52 элемента
и подъём от `prefix` до первого свободного слота или до выхода за массив. Ветка «ничего не
нашли» при такой границе недостижима.

- Time: O(n) — три прохода по массиву, поиск свободного слота не длиннее 51 шага
- Space: O(1) — массив занятости фиксированного размера 52
