# 628. Maximum Product of Three Numbers

- Difficulty: easy
- Link: https://leetcode.com/problems/maximum-product-of-three-numbers/
- Topics: Array, Math, Sorting
- Acceptance: 46.8%

## Statement

Given an integer array `nums`, _find three numbers whose product is maximum and return the maximum product_.

**Example 1:**

```
Input: nums = [1,2,3]
Output: 6
```

**Example 2:**

```
Input: nums = [1,2,3,4]
Output: 24
```

**Example 3:**

```
Input: nums = [-1,-2,-3]
Output: -6
```

**Constraints:**

- `3 <= nums.length <= 10^4`
- `-1000 <= nums[i] <= 1000`

## Similar

- [Maximum Product Subarray](https://leetcode.com/problems/maximum-product-subarray/) (medium)
- [Maximum Product of Three Elements After One Replacement](https://leetcode.com/problems/maximum-product-of-three-elements-after-one-replacement/) (medium)

## Idea

Произведение трёх чисел максимально в одном из двух случаев. Если среди трёх
взятых нет отрицательных или их два, знак плюс, и выгоднее всего взять самые
большие по модулю числа нужного знака. Отсюда ровно два кандидата:

- три наибольших: `hi1 · hi2 · hi3`;
- два наименьших (самые большие по модулю отрицательные) и наибольшее:
  `lo1 · lo2 · hi1`.

Вариант с одним или тремя отрицательными множителями выигрывает, только когда
положительных вариантов нет вовсе (например, все числа отрицательны), — и тогда
он совпадает с первым кандидатом: три наибольших из отрицательных дают
наименьший по модулю результат.

Сортировка находит эти пять чисел за O(n log n), но знать весь порядок не
нужно — нужны только крайние порядковые статистики. Кормен начинает главу о
них с того же наблюдения: «We can solve the selection problem in O(n lg n)
time, since we can sort the numbers… This chapter presents faster algorithms»
(*Introduction to Algorithms*, 3rd ed., с. 213), а §9.1 показывает, что минимум
и максимум находятся одним проходом с текущими рекордами. Здесь рекордов пять:
три сверху, два снизу; каждый элемент сдвигает цепочку рекордов своего края.

Второе следствие отказа от сортировки — вход не меняется. `slices.Sort(nums)`
сортировал срез вызывающего; это закреплено тестом
`TestMaximumProductDoesNotMutateInput`.

Решение сверено с перебором всех троек на случайных входах длины 3–9 со
значениями в `[-5, 5]` — узкий диапазон часто даёт нули и равные значения.
На `n = 10⁴` линейная версия быстрее сортировки на 95.5% (361 → 16.1 мкс,
`-count=10`, `benchstat`, `p = 0.000`).

- Time: O(n)
- Space: O(1)
