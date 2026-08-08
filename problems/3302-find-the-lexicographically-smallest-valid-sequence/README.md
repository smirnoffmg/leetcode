# 3302. Find the Lexicographically Smallest Valid Sequence

- Difficulty: medium
- Link: https://leetcode.com/problems/find-the-lexicographically-smallest-valid-sequence/
- Topics: Two Pointers, String, Dynamic Programming, Greedy
- Acceptance: 44.9%

## Statement

You are given two strings `word1` and `word2`.

A string `x` is called **almost equal** to `y` if you can change **at most** one character in `x` to make it _identical_ to `y`.

A sequence of indices `seq` is called **valid** if:

- The indices are sorted in **ascending** order.
- _Concatenating_ the characters at these indices in `word1` in **the same** order results in a string that is **almost equal**
  to `word2`.

Return an array of size `word2.length` representing the lexicographically smallest **valid** sequence of indices.
If no such sequence of indices exists, return an **empty** array.

**Note** that the answer must represent the _lexicographically smallest array_, **not** the corresponding string formed
by those indices.

**Example 1:**

**Input:** word1 = "vbcca", word2 = "abc"

**Output:** [0,1,2]

**Explanation:**

The lexicographically smallest valid sequence of indices is `[0, 1, 2]`:

- Change `word1[0]` to `'a'`.
- `word1[1]` is already `'b'`.
- `word1[2]` is already `'c'`.

**Example 2:**

**Input:** word1 = "bacdc", word2 = "abc"

**Output:** [1,2,4]

**Explanation:**

The lexicographically smallest valid sequence of indices is `[1, 2, 4]`:

- `word1[1]` is already `'a'`.
- Change `word1[2]` to `'b'`.
- `word1[4]` is already `'c'`.

**Example 3:**

**Input:** word1 = "aaaaaa", word2 = "aaabc"

**Output:** []

**Explanation:**

There is no valid sequence of indices.

**Example 4:**

**Input:** word1 = "abc", word2 = "ab"

**Output:** [0,1]

**Constraints:**

- `1 <= word2.length < word1.length <= 3 * 10^5`
- `word1` and `word2` consist only of lowercase English letters.

## Hints

<details><summary>Hint 1</summary>

Let `dp[i]` be the longest suffix of `word2` that exists as a subsequence of suffix of the substring of `word1`
starting at index `i`.

</details>

<details><summary>Hint 2</summary>

If `dp[i + 1]  and `word1[i] == word2\[m - dp[i + 1] - 1\]`,`dp[i] = dp[i + 1] + 1`. Otherwise, `dp[i] = dp[i + 1]\`.

</details>

<details><summary>Hint 3</summary>

For each index `i`, greedily select characters using the `dp` array to know whether a solution exists.

</details>

## Similar

- [Smallest K-Length Subsequence With Occurrences of a Letter](https://leetcode.com/problems/smallest-k-length-subsequence-with-occurrences-of-a-letter/) (hard)

## Idea

### Что даёт «почти равно»

Замену можно потратить ровно один раз. Значит любое валидное решение — это набор индексов, где все символы совпадают с `word2`, кроме максимум одного. Пока замена не потрачена, у нас есть свобода; как только потрачена — остаток `word2` обязан быть подпоследовательностью остатка `word1` буквально.

### Правый жадный матчинг

Нужен предикат «остаток `word2[j:]` укладывается в `word1[i:]`». Его даёт жадный матчинг справа налево: `last[j]` — самый правый индекс `word1`, на который можно поставить `word2[j]`, если всё, что правее, уже разложено оптимально; `last[m] = n`, недостижимые позиции остаются `-1`. Жадность корректна по тому же аргументу, что и обычная проверка подпоследовательности: брать самое правое вхождение никогда не хуже.

Тогда предикат — это просто `i <= last[j]`, а условие «после замены на позиции `i` хватит остатка» — `i < last[j+1]`.

Индексация по позициям `word2`, а не `word1`, здесь не косметика: массив короче (`m < n`), и вся содержательная работа делается одним проходом двух указателей вместо посимвольного пересчёта длины суффикса. По benchstat (`-count=10`) это −23% на случае без замены и −6% на остальных.

`int32` вместо `int` — отдельный размен, не ускорение: изолированный замер даёт −25% памяти (4.58 МиБ → 3.44 МиБ) ценой +4% времени (p=0.000). Кеш-выигрыша нет, оба массива и так помещаются в L2, а сужение добавляет конверсий. Оставлено ради памяти, которую джадж тоже оценивает.

### Жадный проход слева

Идём по `word1` слева направо с указателем `j` по `word2`:

- совпадение — берём индекс безусловно, замену не тратим;
- несовпадение при неистраченной замене — берём, если `i < last[j+1]`;
- иначе пропускаем.

Почему это лексикографический минимум: массивы индексов сравниваются поэлементно, поэтому первым делом надо минимизировать самый ранний индекс, потом следующий и т.д. Взять индекс `i`, если это не ломает достижимость хвоста, всегда строго лучше, чем пропустить его. При совпадении достижимость не ухудшается никогда (замена остаётся в кармане, а раннее вхождение символа — стандартный жадный оптимум для подпоследовательности), при несовпадении её проверяет `last`. Если после прохода `j < m`, ответа нет.

Корректность проверена случайным перебором против brute force на алфавите из трёх букв.

- Time: O(n) — проход справа для `last` и жадный проход слева
- Space: O(m) — массив `last` из `int32`

Ещё ~14% даёт вынос `last` и буфера ответа в глобальные массивы фиксированного размера (0 аллокаций), но тогда функция возвращает срез разделяемого состояния и перестаёт быть реентерабельной — на джадже это безопасно, в репозитории с тестами нет.
