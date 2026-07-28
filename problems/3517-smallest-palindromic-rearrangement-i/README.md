# 3517. Smallest Palindromic Rearrangement I

- Difficulty: medium
- Link: https://leetcode.com/problems/smallest-palindromic-rearrangement-i/
- Topics: String, Sorting, Counting Sort
- Acceptance: 71.0%

## Statement

You are given a **palindromic** string `s`.

Return the **lexicographically smallest** palindromic permutation of `s`.

**Example 1:**

**Input:** s = "z"

**Output:** "z"

**Explanation:**

A string of only one character is already the lexicographically smallest palindrome.

**Example 2:**

**Input:** s = "babab"

**Output:** "abbba"

**Explanation:**

Rearranging `"babab"` → `"abbba"` gives the smallest lexicographic palindrome.

**Example 3:**

**Input:** s = "daccad"

**Output:** "acddca"

**Explanation:**

Rearranging `"daccad"` → `"acddca"` gives the smallest lexicographic palindrome.

**Constraints:**

- `1 <= s.length <= 10^5`
- `s` consists of lowercase English letters.
- `s` is guaranteed to be palindromic.

## Hints

<details><summary>Hint 1</summary>

Consider a palindrome as composed of two mirror-image halves.

</details>

<details><summary>Hint 2</summary>

Construct one half (using `s`), and then the other half is its reverse to obtain the lexicographically smallest permutation.

</details>

## Similar

- [Shortest Palindrome](https://leetcode.com/problems/shortest-palindrome/) (hard)

## Idea

TODO

- Time: O(?)
- Space: O(?)
