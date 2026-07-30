# 3014. Minimum Number of Pushes to Type Word I

- Difficulty: easy
- Link: https://leetcode.com/problems/minimum-number-of-pushes-to-type-word-i/
- Topics: Math, String, Greedy
- Acceptance: 73.1%

## Statement

You are given a string `word` containing **distinct** lowercase English letters.

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
Input: word = "xycdefghij"
Output: 12
Explanation: The remapped keypad given in the image provides the minimum cost.
"x" -> one push on key 2
"y" -> two pushes on key 2
"c" -> one push on key 3
"d" -> two pushes on key 3
"e" -> one push on key 4
"f" -> one push on key 5
"g" -> one push on key 6
"h" -> one push on key 7
"i" -> one push on key 8
"j" -> one push on key 9
Total cost is 1 + 2 + 1 + 2 + 1 + 1 + 1 + 1 + 1 + 1 = 12.
It can be shown that no other mapping can provide a lower cost.
```

**Constraints:**

- `1 <= word.length <= 26`
- `word` consists of lowercase English letters.
- All letters in `word` are distinct.

## Hints

<details><summary>Hint 1</summary>

We have 8 keys in total. We can type 8 characters with one push each, 8 different characters with two pushes each, and so on.

</details>

<details><summary>Hint 2</summary>

The optimal way is to map letters to keys evenly.

</details>

## Similar

- [Letter Combinations of a Phone Number](https://leetcode.com/problems/letter-combinations-of-a-phone-number/) (medium)

## Idea

Все буквы в слове различны, поэтому частоты не нужны — ответ зависит только от длины слова.

Клавиш 8: первые 8 букв можно посадить на «первое нажатие» (цена 1), следующие 8 — на «второе» (цена 2) и т.д. То есть каждая буква стоит минимум 1 нажатие, а каждый полный слой из 8 букв добавляет всем последующим буквам ещё +1.

Отсюда трюк: прибавляем к ответу текущий остаток длины и сдвигаемся на 8 — каждая буква оказывается посчитана ровно столько раз, сколько стоит. Например, для `l = 10`: `10 + 2 = 12`.

- Time: O(n)
- Space: O(1)
