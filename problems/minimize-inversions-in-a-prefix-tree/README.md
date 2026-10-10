# Minimize Inversions in a Prefix Tree

- Difficulty: hard
- Link: https://leetcode.com/problems/minimize-inversions-in-a-prefix-tree/
- Contest: Q4, номер задачи на момент решения ещё не опубликован

## Statement

You are given an array `words` of `n` **unique** strings containing only lowercase English letters from `'a'` to `'j'`.
You are also given an integer array `target` containing each index from 0 to `n − 1` exactly once.

Build a **prefix tree** from `words`:

- The root represents the empty string `""`.
- Each other node represents a different non-empty prefix of a word. The parent of a node is the prefix obtained by
  removing its last character.
- The node representing the complete word `words[i]` is marked with index `i`. A marked node may also have children.

You may choose the order of the children of **each node independently**. You cannot change which node is a child's
parent.

Traverse the tree starting at the root. When visiting a node, first record its index if it is marked. Then visit its
children in the chosen order, completing the traversal of each child before moving to the next. This is a **pre-order
traversal**.

The recorded indices form an array `traversal` of length `n`. Let `pos[i]` be the position of index `i` in `target`.

An **inversion** is a pair of positions `(a, b)` such that `0 <= a < b < n` and `pos[traversal[a]] > pos[traversal[b]]`.

Return the **minimum possible number of inversions** over all choices of child orderings.

A **prefix** of a string is obtained by removing zero or more characters from the end of the string.

**Example 1:**

**Input:** words = ["bad","bag","fig","fed"], target = [3,1,0,2]

**Output:** 2

**Example 2:**

**Input:** words = ["ab","ac","ba","bc"], target = [0,3,1,2]

**Output:** 1

**Constraints:**

- `1 <= n == words.length <= 10^4`
- `1 <= words[i].length <= 20`
- `words[i]` contains only lowercase English letters from `'a'` to `'j'`
- All `words[i]` are **unique**
- `0 <= target[i] <= n - 1`
- `target` is a permutation of integers from 0 to `n - 1`

## Idea

Заменим каждое слово его позицией в `target`. Тогда ответ — число инверсий в
последовательности позиций, которую выдаёт прямой обход.

**Инверсии раскладываются по узлам.** Возьмём любую пару слов и их наименьшего общего
предка `v` в дереве. Есть ровно три случая:

1. Одно слово — метка самого `v`, другое — его потомок. Метка записывается раньше всех
   потомков при любом порядке детей, поэтому вклад фиксирован: число потомков с позицией
   меньше позиции `v`.
1. Слова лежат в поддеревьях разных детей `A` и `B` узла `v`. Обход выписывает поддерево
   целиком, поэтому, если `A` идёт раньше `B`, пара инвертирована при `pos a > pos b`.
   Сумма по всем таким парам — `cross(A, B)`. Важно, что она зависит только от того,
   **кто из двух детей раньше**, а не от порядков внутри них.
1. Оба слова в поддереве одного ребёнка — пара решается ниже, в этом ребёнке.

Значит, выбор порядка детей в `v` влияет только на слагаемые второго вида, и решения в
разных узлах независимы. Ответ — сумма по узлам: «метка против потомков» плюс минимум
по перестановкам детей суммы `cross` по упорядоченным парам.

**Порядок детей.** Детей не больше 10 (буквы `a..j`), поэтому порядок подбирается DP по
подмножествам: `dp[mask]` — лучшая цена, если дети из `mask` уже стоят первыми в каком-то
порядке. Добавление ребёнка `b` следующим стоит `Σ cross[x][b]` по `x ∈ mask`, и не важно,
в каком порядке стоят дети внутри `mask`. Это та же схема, что в динамике Хелда — Карпа
для коммивояжёра, только цена перехода зависит от множества, а не от последнего элемента.
Итого `O(2^c · c)` на узел. Таблица «цена поставить `b` после `mask`» строится
инкрементально по младшему биту маски.

**Почему не сортировка.** Хочется отсортировать детей компаратором
`cross(A,B) < cross(B,A)`, как в обменном рассуждении. Перебор ловит контрпример:
`words = ["da","b","a","cc","ba","c"]`, `target = [1,0,5,2,4,3]`, жадность даёт 5,
оптимум 4. Направленного цикла среди детей корня тут нет. Ломается **равенство**:
`a ~ b` и `b ~ d` (у обоих порядков одна цена), но `d` строго раньше `a`. Отношение
«не хуже» не транзитивно, и сортировка вправе вернуть `a, b, d, c`. Кроме того,
доказательства, что попарный минимум даёт глобальный, тут нет: цены `cross` — не
разности весов, и правило Смита к ним не применяется.

**Как считать `cross`.** Каждое поддерево возвращает отсортированный список своих позиций.
`cross(A, B)` — это число пар `x ∈ A, y ∈ B` с `x > y`, и считается оно одним проходом
двух указателей по отсортированным спискам. Это шаг слияния из подсчёта инверсий
сортировкой слиянием (CLRS, 3rd ed., Problem 2-4, с. 42 (PDF 63)). Потом списки детей
сливаются в список узла. Глубина дерева не больше 20, поэтому каждая позиция проходит
не больше 20 слияний и при каждом участвует не больше чем в 10 подсчётах `cross`.

Проверка: эталон по определению перебирает все порядки детей во всех узлах, строит каждый
обход и считает инверсии в лоб. Сверка идёт на 3000 случайных маленьких входах и на 30
входах с семью детьми у корня. Мутант «сортировка компаратором» сверку не проходит.

- Time: O(n · L · σ + m · 2^σ · σ), где `L ≤ 20` — длина слова, `σ = 10` — алфавит,
  `m` — число узлов с двумя и более детьми. На полном 10-арном дереве из `10^4` слов — около 60 мс
- Space: O(n · L · σ) на бор, плюс `O(σ · 2^σ)` на таблицы DP в узле
