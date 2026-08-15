# 3702. Longest Subsequence With Non-Zero Bitwise XOR

- Difficulty: medium
- Link: https://leetcode.com/problems/longest-subsequence-with-non-zero-bitwise-xor/
- Topics: Array, Bit Manipulation
- Acceptance: 47.3%

## Statement

You are given an integer array `nums`.

Return the length of the **longest subsequence** in `nums` whose bitwise **XOR** is **non-zero**.
If no such **subsequence** exists, return 0.

**Example 1:**

**Input:** nums = [1,2,3]

**Output:** 2

**Explanation:**

One longest subsequence is `[2, 3]`. The bitwise XOR is computed as `2 XOR 3 = 1`, which is non-zero.

**Example 2:**

**Input:** nums = [2,3,4]

**Output:** 3

**Explanation:**

The longest subsequence is `[2, 3, 4]`. The bitwise XOR is computed as `2 XOR 3 XOR 4 = 5`, which is non-zero.

**Constraints:**

- `1 <= nums.length <= 10^5`
- `0 <= nums[i] <= 10^9`

## Hints

<details><summary>Hint 1</summary>

What happens if you take the entire array?

</details>

<details><summary>Hint 2</summary>

If the XOR of the entire array is 0, can removing one element help?

</details>

<details><summary>Hint 3</summary>

What if all elements are 0?

</details>

## Idea

Слово «подпоследовательность» здесь — ловушка: раз порядок на XOR не влияет, мы просто
выбираем подмножество. А значит, вопрос не «какое окно взять», а «сколько элементов
можно оставить».

Начинаем с жадного максимума — всего массива — и смотрим, чего стоит его починить:

- XOR всего массива ненулевой → ответ `n`, ничего чинить не надо.
- XOR нулевой, но есть ненулевой элемент `x` → выкидываем его один: XOR остатка равен
  `0 ^ x = x`, то есть ненулевой. Ответ `n - 1`, и лучше не бывает — `n` уже отпало.
- Все элементы нули → XOR любого подмножества нулевой, ответ `0`.

Ключевое наблюдение: удаление элемента из множества меняет XOR предсказуемо — на сам
этот элемент. Поэтому «цена ремонта» либо ноль, либо ровно один элемент, либо ремонт
невозможен; перебирать длины между `n` и `0` не нужно.

Наивный соблазн — скользящее окно — ломается дважды: он решает задачу про подмассив
(`[0,1,1,0]`: лучшее окно — 2, лучший ответ — 3), и сам приём требует монотонности,
которой у XOR нет.

- Time: O(n)
- Space: O(1)
