# 3016. Minimum Number of Pushes to Type Word II

- Difficulty: medium
- Link: https://leetcode.com/problems/minimum-number-of-pushes-to-type-word-ii/
- Topics: Hash Table, String, Greedy, Sorting, Counting
- Acceptance: 81.5%

## Statement

You are given a string `word` containing lowercase English letters.

Telephone keypads have keys mapped with **distinct** collections of lowercase English letters, which can be used to form words by pushing them. For example, the key `2` is mapped with `["a","b","c"]`, we need to push the key one time to type `"a"`, two times to type `"b"`, and three times to type `"c"` _._

It is allowed to remap the keys numbered `2` to `9` to **distinct** collections of letters. The keys can be remapped to **any** amount of letters, but each letter **must** be mapped to **exactly** one key. You need to find the **minimum** number of times the keys will be pushed to type the string `word`.

Return _the **minimum** number of pushes needed to type_ `word` _after remapping the keys_.

An example mapping of letters to keys on a telephone keypad is given below. Note that `1`, `*`, `#`, and `0` do **not** map to any letters.

**Example 1:**

```
Input: word = "abcde"
Output: 5
Explanation: The remapped keypad given in the image provides the minimum cost.
"a" -> one push on key 2
"b" -> one push on key 3
"c" -> one push on key 4
"d" -> one push on key 5
"e" -> one push on key 6
Total cost is 1 + 1 + 1 + 1 + 1 = 5.
It can be shown that no other mapping can provide a lower cost.
```

**Example 2:**

```
Input: word = "xyzxyzxyzxyz"
Output: 12
Explanation: The remapped keypad given in the image provides the minimum cost.
"x" -> one push on key 2
"y" -> one push on key 3
"z" -> one push on key 4
Total cost is 1 * 4 + 1 * 4 + 1 * 4 = 12
It can be shown that no other mapping can provide a lower cost.
Note that the key 9 is not mapped to any letter: it is not necessary to map letters to every key, but to map all the letters.
```

**Example 3:**

```
Input: word = "aabbccddeeffgghhiiiiii"
Output: 24
Explanation: The remapped keypad given in the image provides the minimum cost.
"a" -> one push on key 2
"b" -> one push on key 3
"c" -> one push on key 4
"d" -> one push on key 5
"e" -> one push on key 6
"f" -> one push on key 7
"g" -> one push on key 8
"h" -> two pushes on key 9
"i" -> one push on key 9
Total cost is 1 * 2 + 1 * 2 + 1 * 2 + 1 * 2 + 1 * 2 + 1 * 2 + 1 * 2 + 2 * 2 + 6 * 1 = 24.
It can be shown that no other mapping can provide a lower cost.
```

**Constraints:**

- `1 <= word.length <= 10^5`
- `word` consists of lowercase English letters.

## Hints

<details><summary>Hint 1</summary>

We have 8 keys in total. We can type 8 characters with one push each, 8 different characters with two pushes each, and so on.

</details>

<details><summary>Hint 2</summary>

The optimal way is to map letters to keys evenly.

</details>

<details><summary>Hint 3</summary>

Sort the letters by frequencies in the word in non-increasing order.

</details>

## Similar

- [Letter Combinations of a Phone Number](https://leetcode.com/problems/letter-combinations-of-a-phone-number/) (medium)

## Idea

В отличие от [3014](../3014-minimum-number-of-pushes-to-type-word-i/), буквы могут повторяться, поэтому важна не длина слова, а частоты букв.

Клавиш 8, и раскладка свободная: 8 букв можно посадить на «первое нажатие» (цена 1), следующие 8 — на «второе» (цена 2) и т.д. Цена буквы зависит только от её **ранга** по частоте: буква ранга `r` (с нуля) стоит `r/8 + 1` нажатий. Жадность очевидна — самые частые буквы должны получить самую дешёвую позицию, поэтому сортируем частоты по убыванию.

Ответ — это `сумма частот × цена`, но считать цену явно не нужно. Множитель `r/8 + 1` означает, что счётчик буквы входит в сумму по одному разу за каждый пройденный ярус из 8 клавиш, то есть

```
ответ = (сумма всех частот) + (сумма частот с ранга ≥ 8) + (≥ 16) + (≥ 24)
```

Это тот же приём со сдвигом на 8, что и в 3014, только суммируются не остатки длины, а суффиксы отсортированного массива частот. Первое слагаемое как раз равно `len(word)`.

Например, для `"aabbccddeeffgghhiiiiii"` частоты по убыванию — `[6,2,2,2,2,2,2,2,2,0,...]`: `22 + 2 + 0 + 0 = 24`.

Подсчёт частот идёт побайтово, а не через `range` по строке: алфавит по условию — только `a`–`z`, а `range` на каждом шаге декодирует UTF-8, что здесь лишняя работа.

Сортируются всегда ровно 26 чисел, поэтому её вклад в сложность константный.

- Time: O(n)
- Space: O(1)
