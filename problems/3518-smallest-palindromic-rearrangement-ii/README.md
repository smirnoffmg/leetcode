# 3518. Smallest Palindromic Rearrangement II

- Difficulty: hard
- Link: https://leetcode.com/problems/smallest-palindromic-rearrangement-ii/
- Topics: Hash Table, Math, String, Combinatorics, Counting
- Acceptance: 31.7%

## Statement

You are given a **palindromic** string `s` and an integer `k`.

Return the **k-thlexicographically smallest** palindromic permutation of `s`. If there are fewer than `k` distinct palindromic permutations, return an empty string.

**Note:** Different rearrangements that yield the same palindromic string are considered identical and are counted once.

**Example 1:**

**Input:** s = "abba", k = 2

**Output:** "baab"

**Explanation:**

- The two distinct palindromic rearrangements of `"abba"` are `"abba"` and `"baab"`.
- Lexicographically, `"abba"` comes before `"baab"`. Since `k = 2`, the output is `"baab"`.

**Example 2:**

**Input:** s = "aa", k = 2

**Output:** ""

**Explanation:**

- There is only one palindromic rearrangement: `"aa"`.
- The output is an empty string since `k = 2` exceeds the number of possible rearrangements.

**Example 3:**

**Input:** s = "bacab", k = 1

**Output:** "abcba"

**Explanation:**

- The two distinct palindromic rearrangements of `"bacab"` are `"abcba"` and `"bacab"`.
- Lexicographically, `"abcba"` comes before `"bacab"`. Since `k = 1`, the output is `"abcba"`.

**Constraints:**

- `1 <= s.length <= 10^4`
- `s` consists of lowercase English letters.
- `s` is guaranteed to be palindromic.
- `1 <= k <= 10^6`

## Hints

<details><summary>Hint 1</summary>

Only build `floor(n / 2)` characters (the rest are determined by symmetry).

</details>

<details><summary>Hint 2</summary>

Count character frequencies and use half the counts for construction.

</details>

<details><summary>Hint 3</summary>

Incrementally choose each character (from smallest to largest) and calculate how many valid arrangements result if that character is chosen at the current index.

</details>

<details><summary>Hint 4</summary>

If the count is at least `k`, fix that character; otherwise, subtract the count from `k` and try the next candidate.

</details>

<details><summary>Hint 5</summary>

Use combinatorics to compute the number of permutations at each step.

</details>

## Idea

Палиндром полностью определяется своей левой половиной: серединный символ (при нечётной
длине) — единственный с нечётной частотой, правая половина — зеркало левой. Значит задача
сводится к поиску `k`-й в лексикографическом порядке перестановки мультимножества из
половинных частот `cnt[c] / 2`.

Идём по позициям левой половины и на каждой пробуем символы от `'a'` к `'z'`: если число
различных перестановок оставшегося мультимножества `>= k`, символ фиксируем, иначе вычитаем
это число из `k` и берём следующий кандидат. Как только `k == 1`, остаток дописываем
отсортированным. Если общее число перестановок меньше исходного `k` — ответа нет.

Число перестановок мультимножества считается как произведение биномиальных коэффициентов
`C(total, cnt[c])` и обрезается потолком `1e6 + 1`: точные значения астрономические, а `k`
не превышает `1e6`. Обрезание же делает подсчёт быстрым — биномиальный цикл упирается в
потолок за ~20 итераций, а когда потолок недостижим, сам `min(r, n - r)` мал.

- Time: O(n · A² · log(limit)), где A = 26 — размер алфавита
- Space: O(n)
