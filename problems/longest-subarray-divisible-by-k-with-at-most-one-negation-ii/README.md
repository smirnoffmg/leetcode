# Longest Subarray Divisible by K with At Most One Negation II

- Difficulty: hard
- Link: https://leetcode.com/problems/longest-subarray-divisible-by-k-with-at-most-one-negation-ii/
- Contest: Q4, номер задачи на момент решения ещё не опубликован

## Statement

You are given an integer array `nums` and an integer `k`.

A subarray is **valid** if its sum is divisible by `k`, or can become divisible by `k` by **negating one element within
that subarray**.

Negating an element means replacing its value `x` with `-x`.

Return the **length of the longest valid subarray**. If no valid subarray exists, return 0.

A **subarray** is a contiguous, non-empty sequence of elements within an array.

**Example 1:**

**Input:** nums = [4,1,2], k = 3

**Output:** 3

**Explanation:** Negating `nums[2] = 2` changes the sum to `4 + 1 - 2 = 3`, which is divisible by `k`.

**Example 2:**

**Input:** nums = [5,3,4], k = 7

**Output:** 2

**Explanation:** The subarray `[3, 4]` has a sum of 7, divisible by `k = 7` without any negation.

**Example 3:**

**Input:** nums = [2,2,5], k = 6

**Output:** 2

**Explanation:** Negating either element of `[2, 2]` gives a sum of 0.

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^5 <= nums[i] <= 10^5`
- `1 <= k <= 3000`

## Idea

Отрезок `(l, r]` через префиксные суммы валиден, если

- `P[l] ≡ P[r]` — без инверсии, или
- `P[l] + 2·a_i ≡ P[r]` для некоторого `l ≤ i < r` — с инверсией `a_i`.

Первое условие — классика «самый длинный подмассив, делящийся на `k`»: для каждого остатка
хранить **самый ранний** префикс `first[u]`, потому что для длины нужен минимальный `l`.

Второе сводится к тому же, если для каждого остатка `c` хранить `best[c]` — самый ранний `l`,
который вместе с каким-то `i ≥ l` даёт `P[l] + 2a_i ≡ c`. Для фиксированной пары
«остаток префикса `u`, сдвиг `d = 2a_i mod k`» лучший `l` — это `first[u]`, и от того,
какой именно `i` дал сдвиг `d`, ответ не зависит — нужно лишь `first[u] ≤ i`. Поэтому
пару `(u, d)` достаточно учесть **один раз**: при первом `i` со сдвигом `d`, когда `u` уже
встречался. Если `u` появится позже, пару учтёт следующий элемент с тем же `d`.

Отсюда `order` — остатки в порядке первого появления, и `applied[d]` — сколько из них уже
пропущено через сдвиг `d`. Пар не больше `k²`, так что вся работа — `O(n + k²)`.

Лобовое решение пересчитывает все `k` остатков на каждом шаге: `O(n · k) = 3·10^8`. Оно тоже
может уложиться, но держится на запасе по константе; `O(k²) ≈ 9·10^6` от `n` не зависит.

Ловушки, на которых ломалась ручная перепечатка:

- «не встречался» — это `-1`, а `0` — законная позиция пустого префикса. Проверка `> 0`
  вместо `>= 0` теряет все отрезки от начала массива: пример 3 даёт 0.

- Цикл по `order` закрывается сразу после обновления `best`; всё остальное — один раз на
  элемент. Если внутренний цикл проглотит обновление `p`, сумма сдвигается столько раз,
  сколько в нём итераций.

- Time: O(n + k²)

- Space: O(k)
