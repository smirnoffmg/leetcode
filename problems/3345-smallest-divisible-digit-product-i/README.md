# 3345. Smallest Divisible Digit Product I

- Difficulty: easy
- Link: https://leetcode.com/problems/smallest-divisible-digit-product-i/
- Topics: Math, Enumeration
- Acceptance: 70.8%

## Statement

You are given two integers `n` and `t`. Return the **smallest** number greater than or equal to `n` such that
the **product of its digits** is divisible by `t`.

**Example 1:**

**Input:** n = 10, t = 2

**Output:** 10

**Explanation:**

The digit product of 10 is 0, which is divisible by 2, making it the smallest number greater than or equal to 10 that
satisfies the condition.

**Example 2:**

**Input:** n = 15, t = 3

**Output:** 16

**Explanation:**

The digit product of 16 is 6, which is divisible by 3, making it the smallest number greater than or equal to 15
that satisfies the condition.

**Constraints:**

- `1 <= n <= 100`
- `1 <= t <= 10`

## Hints

<details><summary>Hint 1</summary>

You have to check at most 10 numbers.

</details>

<details><summary>Hint 2</summary>

Apply a brute-force approach by checking each possible number.

</details>

## Similar

- [Smallest Number With Given Digit Product](https://leetcode.com/problems/smallest-number-with-given-digit-product/) (medium)

## Idea

Перебор от `n` вверх до первого числа, у которого произведение цифр делится на `t`.

Цикл гарантированно короткий: среди любых 10 подряд идущих чисел есть оканчивающееся
на 0, а его произведение цифр равно 0 и делится на любое `t`. Значит, проверяем не
более 10 чисел.

- Time: O(1) — ≤ 10 кандидатов, у каждого O(len(digits)) на произведение цифр
- Space: O(1)
