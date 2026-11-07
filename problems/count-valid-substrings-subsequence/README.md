# Count Valid Substrings (subsequence)

- Link: задача контестная; название и слаг условные, функция называется `countValidSubstrings`, как и в варианте с подстрокой
- Contest: номер, сложность и ссылка на момент решения не опубликованы

## Statement

You are given a string `s` consisting only of the characters `'a'`, `'b'`, and `'c'`.

A string is called **valid** if it does not contain `"abc"` as a **subsequence**.

Return the number of **valid substrings** in `s`. Substrings with different starting or ending indices are counted separately.

A **substring** is a contiguous **non-empty** sequence of characters within a string.

A **subsequence** is a string that can be derived from another string by deleting some or no characters without changing the order of the remaining characters.

**Example 1:**

```
Input: s = "abac"
Output: 9
```

**Example 2:**

```
Input: s = "abc"
Output: 5
```

**Example 3:**

```
Input: s = "ac"
Output: 3
```

**Constraints:**

- `1 <= s.length <= 10^5`
- `s` only consists of `'a'`, `'b'`, and `'c'` characters.

## Idea

Если подстрока содержит `abc` как подпоследовательность, то содержит и любая её надстрока. Значит, при фиксированном правом конце хорошие левые концы образуют хвост, и ответ — сумма `right - L(right)`, где `L(right)` — самый правый старт, с которого `abc` ещё укладывается в `s[..right]`.

В отличие от варианта с подстрокой, буквы вхождения могут быть разнесены, поэтому одного запомненного индекса мало. Нужна цепочка «самый правый старт префикса длины k»:

- `a` в позиции `right`: префикс `a` теперь начинается в `right`;
- `b`: префикс `ab` начинается там же, где последний префикс `a`;
- `c`: префикс `abc` начинается там же, где последний префикс `ab`.

Обновлять надо именно присваиванием, без `max`: источник перехода не убывает, поэтому новое значение не меньше старого. Значения `-1` до первого вхождения дают в формуле все `right + 1` подстрок.

Все три перехода — одна формула: `starts[k+1] = starts[k]`, где `k` — номер буквы, а `starts[0]` — старт пустого префикса, то есть сам `right`. Поэтому цикл обходится без ветвлений.

Сверка с перебором всех подстрок идёт на 20 000 случайных строк длиной до 12 символов.

Варианты замерены в одном бинарнике (`benchstat`, `-benchtime=100x`, `-count=5`, Apple M2 Pro, `n = 10^5`, `p = 0.008`):

| Вход      | таблица без ветвлений | `switch`  |
| --------- | --------------------- | --------- |
| `random`  | 120.2 мкс             | 438.8 мкс |
| `all-abc` | 118.9 мкс             | 130.7 мкс |

На случайной строке таблица быстрее на 73%, на предсказуемой на 9%. Разрыв на случайном входе согласуется с ошибками предсказания ветвлений, но счётчики процессора не снимались, так что это гипотеза. Замер короткий, поэтому оценка разброса у него грубая.

- Time: O(n)
- Space: O(1)
