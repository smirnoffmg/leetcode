# 2213. Longest Substring of One Repeating Character

- Difficulty: hard
- Link: https://leetcode.com/problems/longest-substring-of-one-repeating-character/
- Topics: Array, String, Segment Tree, Ordered Set
- Acceptance: 54.7%

## Statement

You are given a **0-indexed** string `s`. You are also given a **0-indexed** string `queryCharacters` of length `k` and a **0-indexed** array of integer **indices** `queryIndices` of length `k`, both of which are used to describe `k` queries.

The `i^th` query updates the character in `s` at index `queryIndices[i]` to the character `queryCharacters[i]`.

Return _an array_ `lengths` _of length_ `k`\_ where\_ `lengths[i]` _is the **length** of the **longest substring** of_ `s`\_ consisting of **only one repeating** character **after** the\_ `i^th` _query is performed._

**Example 1:**

```
Input: s = "babacc", queryCharacters = "bcb", queryIndices = [1,3,3]
Output: [3,3,4]
Explanation:
- 1st query updates s = "bbbacc". The longest substring consisting of one repeating character is "bbb" with length 3.
- 2nd query updates s = "bbbccc".
  The longest substring consisting of one repeating character can be "bbb" or "ccc" with length 3.
- 3rd query updates s = "bbbbcc". The longest substring consisting of one repeating character is "bbbb" with length 4.
Thus, we return [3,3,4].
```

**Example 2:**

```
Input: s = "abyzz", queryCharacters = "aa", queryIndices = [2,1]
Output: [2,3]
Explanation:
- 1st query updates s = "abazz". The longest substring consisting of one repeating character is "zz" with length 2.
- 2nd query updates s = "aaazz". The longest substring consisting of one repeating character is "aaa" with length 3.
Thus, we return [2,3].
```

**Constraints:**

- `1 <= s.length <= 10^5`
- `s` consists of lowercase English letters.
- `k == queryCharacters.length == queryIndices.length`
- `1 <= k <= 10^5`
- `queryCharacters` consists of lowercase English letters.
- `0 <= queryIndices[i] < s.length`

## Hints

<details><summary>Hint 1</summary>

Use a segment tree to perform fast point updates and range queries.

</details>

<details><summary>Hint 2</summary>

We need each segment tree node to store the length of the longest substring of that segment consisting of only 1 repeating character.

</details>

<details><summary>Hint 3</summary>

We will also have each segment tree node store the leftmost and rightmost character of the segment, the max length of a prefix substring consisting of only 1 repeating character, and the max length of a suffix substring consisting of only 1 repeating character.

</details>

<details><summary>Hint 4</summary>

Use this information to properly merge the two segment tree nodes together.

</details>

## Similar

- [Merge Intervals](https://leetcode.com/problems/merge-intervals/) (medium)
- [Longest Repeating Character Replacement](https://leetcode.com/problems/longest-repeating-character-replacement/) (medium)
- [Consecutive Characters](https://leetcode.com/problems/consecutive-characters/) (easy)
- [Create Sorted Array through Instructions](https://leetcode.com/problems/create-sorted-array-through-instructions/) (hard)
- [Longest Increasing Subsequence II](https://leetcode.com/problems/longest-increasing-subsequence-ii/) (hard)

## Idea

Наивный пересчёт после каждого запроса — O(n) на запрос, итого O(n·k) ≈ 10^10 — не проходит.
Нужна структура, которая после точечной замены символа пересчитывает ответ за O(log n):
дерево отрезков.

Ключевое наблюдение — какую сводку хранить в узле, чтобы она сливалась из детей.
Одной длины лучшего блока недостаточно: при слиянии двух отрезков лучший блок может
пересекать границу. Поэтому узел хранит пять полей: длину отрезка, крайние символы,
длины однобуквенных префикса и суффикса и лучший блок внутри отрезка. Слияние:

- лучший блок — максимум детских плюс, если правый символ левого ребёнка равен левому
  символу правого, склейка «суффикс левого + префикс правого»;
- префикс — префикс левого ребёнка, а если тот однобуквенный целиком и символы на стыке
  совпали, он продлевается префиксом правого (симметрично для суффикса).

Это стандартный приём «сводка с крайними значениями для склейки на стыке» — тот же,
что в задачах про максимальную сумму подотрезка на дереве отрезков.

Точечное обновление меняет лист и пересобирает сводки вверх по пути; ответ на запрос —
сводка корня.

Реализационная тонкость: сводки узлов разложены по плоским массивам (`pre`, `suf`,
`maxLen`, `leftChar`, `rightChar`) вместо среза структур. На n = k = 10⁵ это −27%
времени и −33% памяти (benchstat, p=0.000): слияние читает только нужные поля, не
копирует структуры целиком, а массивы символов занимают байт на узел вместо слова.

- Time: O(n + k·log n)
- Space: O(n)
