# Longest Subarray Divisible by K with At Most One Negation I

- Difficulty: medium
- Link: https://leetcode.com/problems/longest-subarray-divisible-by-k-with-at-most-one-negation-i/
- Contest: Q3, номер задачи на момент решения ещё не опубликован

## Statement

You are given an integer array `nums` and an integer `k`.

A subarray is **valid** if its sum is divisible by `k`, or can become divisible by `k` by **negating one element within
that subarray**.

Negating an element means replacing its value `x` with `-x`.

Return the **length of the longest valid subarray**. If no valid subarray exists, return 0.

A **subarray** is a contiguous, non-empty sequence of elements within an array.

**Example 1:**

**Input:** nums = [4,1,2], k = 3

**Output:** 3

**Explanation:** Negating `nums[2] = 2` changes the sum to `4 + 1 - 2 = 3`, which is divisible by `k`.

**Example 2:**

**Input:** nums = [5,3,4], k = 7

**Output:** 2

**Explanation:** The subarray `[3, 4]` has a sum of 7, divisible by `k = 7` without any negation.

**Example 3:**

**Input:** nums = [2,2,5], k = 6

**Output:** 2

**Explanation:** Negating either element of `[2, 2]` gives a sum of 0.

**Constraints:**

- `1 <= nums.length <= 1000`
- `-10^5 <= nums[i] <= 10^5`
- `1 <= k <= 10^5`

## Idea

Инверсия `x` вычитает из суммы `2x`. Поэтому отрезок с суммой `S` валиден, если
`S ≡ 0 (mod k)` или внутри есть элемент с `2x ≡ S (mod k)`. «Можно ли что-то инвертировать»
сводится к вопросу о **принадлежности остатка множеству** — никакого перебора элементов.

При `n ≤ 1000` хватает всех `n(n+1)/2 ≈ 5·10^5` отрезков: фиксируем `l`, двигаем `r`,
ведём остаток суммы и множество остатков `2x`. Проверка отрезка — один поиск.

Две тонкости:

- Множество — массив меток размера `k`, а не `map`: метка `l + 1` значит «встречен при
  текущем `l`», и массив не надо чистить. Честная очистка стоила бы `O(n · k) = 10^8`.
- Метку элемента `r` надо ставить **до** проверки отрезка `[l, r]`: он уже внутри и тоже
  может быть инвертирован. Обратный порядок теряет случай «инвертировать последний» —
  пример 1 даёт 2 вместо 3.

`%` в Go для отрицательного делимого отрицателен, поэтому `mod` нормализует остаток.

Вариант с большим `n` — в [II](../longest-subarray-divisible-by-k-with-at-most-one-negation-ii/README.md).

- Time: O(n²)
- Space: O(k)
