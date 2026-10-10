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
аллокаций нет. После добавления буквы `c` нарушиться может только её счётчик —
остальные при сжатии лишь убывают, — поэтому цикл сжатия идёт по условию
`count[c] > 2`, а не по проверке всех 26 счётчиков. Первая версия проверяла все
26 на каждом шаге; переход на проверку одной буквы дал −89.6% на `n = 100`
(1.60 мкс → 166 нс, `-count=10`, `benchstat`, `p = 0.000`).

Длина окна не хранится отдельной переменной, а выводится как `right - left + 1`:
отдельный счётчик в первой версии дважды давал off-by-one (см. заметку
«Скользящее окно» в vault).

Перебор всех подстрок с проверкой каждой (O(n²) с инкрементальным подсчётом)
при `n <= 100` тоже прошёл бы — на него и намекает условие; здесь он оставлен
эталоном в тестах, а окно даёт линейность и единственную точку замера.

- Time: O(n) — каждый символ входит в окно и выходит из него по разу.
- Space: O(1) — массив из 26 счётчиков.
