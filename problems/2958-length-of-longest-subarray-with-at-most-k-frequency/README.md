# 2958. Length of Longest Subarray With at Most K Frequency

- Difficulty: medium
- Link: https://leetcode.com/problems/length-of-longest-subarray-with-at-most-k-frequency/
- Topics: Array, Hash Table, Sliding Window
- Acceptance: 59.8%

## Statement

You are given an integer array `nums` and an integer `k`.

The **frequency** of an element `x` is the number of times it occurs in an array.

An array is called **good** if the frequency of each element in this array is **less than or equal** to `k`.

Return _the length of the **longestgood** subarray of_ `nums`_._

A **subarray** is a contiguous non-empty sequence of elements within an array.

**Example 1:**

```
Input: nums = [1,2,3,1,2,3,1,2], k = 2
Output: 6
Explanation: The longest possible good subarray is [1,2,3,1,2,3] since the values 1, 2, and 3 occur at most twice in this subarray. Note that the subarrays [2,3,1,2,3,1] and [3,1,2,3,1,2] are also good.
It can be shown that there are no good subarrays with length more than 6.
```

**Example 2:**

```
Input: nums = [1,2,1,2,1,2,1,2], k = 1
Output: 2
Explanation: The longest possible good subarray is [1,2] since the values 1 and 2 occur at most once in this subarray. Note that the subarray [2,1] is also good.
It can be shown that there are no good subarrays with length more than 2.
```

**Example 3:**

```
Input: nums = [5,5,5,5,5,5,5], k = 4
Output: 4
Explanation: The longest possible good subarray is [5,5,5,5] since the value 5 occurs 4 times in this subarray.
It can be shown that there are no good subarrays with length more than 4.
```

**Constraints:**

- `1 <= nums.length <= 10^5`
- `1 <= nums[i] <= 10^9`
- `1 <= k <= nums.length`

## Hints

<details><summary>Hint 1</summary>

For each index `i`, find the rightmost index `j >= i` such that the frequency of each element in the subarray `[i, j]` is at most `k`.

</details>

<details><summary>Hint 2</summary>

We can use 2 pointers / sliding window to achieve it.

</details>

## Similar

- [Longest Substring with At Least K Repeating Characters](https://leetcode.com/problems/longest-substring-with-at-least-k-repeating-characters/) (medium)

## Idea

Свойство «каждый элемент встречается не более k раз» монотонно по вложенности:
любой подотрезок хорошего отрезка тоже хороший. Значит, работает скользящее окно:
для каждой правой границы существует минимальная левая, при которой окно хорошее,
и с ростом правой границы она только сдвигается вправо — назад возвращаться незачем.

Держим счётчик частот элементов окна. Добавили `nums[right]` — нарушиться могла
только его частота, поэтому сжимаем окно слева ровно до тех пор, пока она не
вернётся к `k`. Остальные частоты при сжатии лишь убывают, так что окно снова
хорошее, и его длина — кандидат в ответ.

Наивный вариант — перебор всех подотрезков с проверкой каждого — стоит O(n²) даже
с инкрементальным подсчётом частот и на n = 10⁵ не проходит. Окно же двигает оба
указателя только вперёд: каждый элемент входит в окно и выходит из него по разу.

- Time: O(n)
- Space: O(n) — счётчик частот; в худшем случае все элементы различны.
