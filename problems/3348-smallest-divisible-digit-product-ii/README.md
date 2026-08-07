# 3348. Smallest Divisible Digit Product II

- Difficulty: hard
- Link: https://leetcode.com/problems/smallest-divisible-digit-product-ii/
- Topics: Math, String, Backtracking, Greedy, Number Theory
- Acceptance: 31.8%

## Statement

You are given a string `num` which represents a **positive** integer, and an integer `t`.

A number is called **zero-free** if _none_ of its digits are 0.

Return a string representing the **smallestzero-free** number greater than or equal to `num` such that the **product of its digits** is divisible by `t`. If no such number exists, return `"-1"`.

**Example 1:**

**Input:** num = "1234", t = 256

**Output:** "1488"

**Explanation:**

The smallest zero-free number that is greater than 1234 and has the product of its digits divisible by 256 is 1488, with the product of its digits equal to 256.

**Example 2:**

**Input:** num = "12355", t = 50

**Output:** "12355"

**Explanation:**

12355 is already zero-free and has the product of its digits divisible by 50, with the product of its digits equal to 150.

**Example 3:**

**Input:** num = "11111", t = 26

**Output:** "-1"

**Explanation:**

No number greater than 11111 has the product of its digits divisible by 26.

**Constraints:**

- `2 <= num.length <= 2 * 10^5`
- `num` consists only of digits in the range `['0', '9']`.
- `num` does not contain leading zeros.
- `1 <= t <= 10^14`

## Hints

<details><summary>Hint 1</summary>

`t` should only have 2, 3, 5 and 7 as prime factors.

</details>

<details><summary>Hint 2</summary>

Find the shortest suffix that must be changed.

</details>

<details><summary>Hint 3</summary>

Try to form the string greedily.

</details>

## Similar

- [Smallest Number With Given Digit Product](https://leetcode.com/problems/smallest-number-with-given-digit-product/) (medium)

## Idea

### Наблюдения

Произведение цифр 1–9 раскладывается только на простые 2, 3, 5, 7. Если у `t` есть другой простой делитель, ответа нет — `-1`. Иначе `t` — это вектор степеней `(e2, e3, e5, e7)`, и «произведение делится на `t`» означает покомпонентное покрытие этого вектора суммой векторов цифр.

Ключевая величина — `minLen(need)`: минимальное число ненулевых цифр, покрывающих остаток `need`. Пятёрки и семёрки дают только цифры 5 и 7 — ровно по одной степени, их количество фиксировано. Для части `2^a * 3^b` жадная упаковка: `a/3` восьмёрок, `b/2` девяток, остатки `a%3` и `b%2` закрываются одной цифрой (2, 4, 3 или 6), и только пара остатков `2^2 * 3` требует двух. Одна цифра уменьшает `minLen` максимум на 1, поэтому жадная оценка точна.

### Перебор кандидатов

Любое число `>= num` той же длины — это либо само `num`, либо «префикс `num` + увеличенная цифра в позиции `i` + произвольный хвост». Чем позже расходимся, тем меньше число, поэтому:

1. Если `num` без нулей и его произведение уже покрывает `t` — ответ `num`.
1. Идём по `i` справа налево (не правее первого нуля — префикс обязан быть zero-free). Для каждой цифры `d > num[i]` считаем остаток требования после префикса и `d`; если `minLen(остаток)` влезает в `n-1-i` свободных позиций — нашли точку расхождения.
1. Если длина `n` недостижима, любое число из `n+1` ненулевых цифр уже больше `num`: берём длину `max(n+1, minLen(t))` и минимальное заполнение.

Остаток после префикса считается за O(1) по префиксным суммам степеней.

### Минимальный хвост

Хвост фиксированной длины минимален так: как можно больше единиц спереди (меньше «значащих» цифр — меньше строка), затем цифры покрытия по возрастанию. Пятёрки и семёрки предопределены, а для части `2^a * 3^b` цифры выбираются жадно от меньшей к большей: цифра годится, если после неё остаток всё ещё укладывается в оставшиеся позиции. Это даёт, например, `26` вместо `34` для остатка `2^2 * 3` — оба минимальной длины, но строка меньше. Общая сортировка в конце сливает единицы, пятёрки, семёрки и цифры покрытия.

- Time: O(n) — префиксные суммы, проход по позициям с 9 цифрами и O(1)-проверкой; сортировка хвоста — по факту подсчётная на алфавите из 9 цифр, здесь `slices.Sort` на O(n log n) от длины хвоста
- Space: O(n) — префиксные суммы степеней и буфер ответа
