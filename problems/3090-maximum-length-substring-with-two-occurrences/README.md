# 3090. Maximum Length Substring With Two Occurrences

- Difficulty: easy
- Link: https://leetcode.com/problems/maximum-length-substring-with-two-occurrences/
- Topics: Hash Table, String, Sliding Window
- Acceptance: 73.0%

## Statement

Given a string `s`, return the **maximum** length of a substring such that it contains _at most two occurrences_ of each character.

**Example 1:**

**Input:** s = "bcbbbcba"

**Output:** 4

**Explanation:**

The following substring has a length of 4 and contains at most two occurrences of each character: `"bcbbbcba"`.

**Example 2:**

**Input:** s = "aaaa"

**Output:** 2

**Explanation:**

The following substring has a length of 2 and contains at most two occurrences of each character: `"aaaa"`.

**Constraints:**

- `2 <= s.length <= 100`
- `s` consists only of lowercase English letters.

## Hints

<details><summary>Hint 1</summary>

We can try all substrings by brute-force since the constraints are very small.

</details>

## Idea

Свойство «каждая буква встречается не более двух раз» антимонотонно по
вложенности: подстрока хорошей строки тоже хорошая. Значит, работает скользящее
окно — та же схема, что в 2958 (там «не более k раз» для произвольного k):
правая граница на каждом шаге двигается на символ вперёд, после добавления
окно сжимается слева ровно до восстановления свойства, и длина снова валидного
окна — кандидат в ответ. Замер строго после восстановления: единственная точка,
где окно гарантированно хорошее.

Состояние окна — массив из 26 счётчиков вместо хеш-таблицы: алфавит фиксирован,
аллокаций нет. Валидность проверяется проходом по всем 26 счётчикам, хотя
нарушиться после добавления может только счётчик добавленной буквы — при
`n <= 100` эта константа неощутима.

Перебор всех подстрок с проверкой каждой (O(n²) с инкрементальным подсчётом)
при `n <= 100` тоже прошёл бы — на него и намекает условие; здесь он оставлен
эталоном в тестах, а окно даёт линейность и единственную точку замера.

- Time: O(n) — каждый символ входит в окно и выходит из него по разу; на каждый
  шаг — проверка 26 счётчиков.
- Space: O(1) — массив из 26 счётчиков.
