# 3310. Remove Methods From Project

- Difficulty: medium
- Link: https://leetcode.com/problems/remove-methods-from-project/
- Topics: Depth-First Search, Breadth-First Search, Graph Theory
- Acceptance: 57.3%

## Statement

You are maintaining a project that has `n` methods numbered from `0` to `n - 1`.

You are given two integers `n` and `k`, and a 2D integer array `invocations`, where `invocations[i] = [a_i, b_i]` indicates
that method `a_i` invokes method `b_i`.

There is a known bug in method `k`. Method `k`, along with any method invoked by it, either **directly** or **indirectly**,
are considered **suspicious** and we aim to remove them.

A group of methods can only be removed if no method **outside** the group invokes any methods **within** it.

Return an array containing all the remaining methods after removing all the **suspicious** methods. You may return the
answer in _any order_. If it is not possible to remove **all** the suspicious methods, **none** should be removed.

**Example 1:**

**Input:** n = 4, k = 1, invocations = [[1,2],[0,1],[3,2]]

**Output:** [0,1,2,3]

**Explanation:**

Method 2 and method 1 are suspicious, but they are directly invoked by methods 3 and 0, which are not suspicious. We return all elements without removing anything.

**Example 2:**

**Input:** n = 5, k = 0, invocations = [[1,2],[0,2],[0,1],[3,4]]

**Output:** [3,4]

**Explanation:**

Methods 0, 1, and 2 are suspicious and they are not directly invoked by any other method. We can remove them.

**Example 3:**

**Input:** n = 3, k = 2, invocations = [[1,2],[0,1],[2,0]]

**Output:** []

**Explanation:**

All methods are suspicious. We can remove them.

**Constraints:**

- `1 <= n <= 10^5`
- `0 <= k <= n - 1`
- `0 <= invocations.length <= 2 * 10^5`
- `invocations[i] == [a_i, b_i]`
- `0 <= a_i, b_i <= n - 1`
- `a_i != b_i`
- `invocations[i] != invocations[j]`

## Hints

<details><summary>Hint 1</summary>

Use DFS from node `k`.

</details>

<details><summary>Hint 2</summary>

Mark all the nodes visited from node `k`, and then check if they can be visited from the other nodes.

</details>

## Idea

### Постановка

Методы — вершины орграфа, `invocations` — рёбра `a → b` («`a` вызывает `b`»). Подозрительное множество `S` — это в точности вершины, достижимые из `k`, включая саму `k`.

Условие удаления «никакой метод вне группы не вызывает метод внутри неё» на языке графа означает: не существует ребра `a → b`, где `a ∉ S`, а `b ∈ S`. Рёбра в обратную сторону (`a ∈ S`, `b ∉ S`) не мешают: удаляемый метод может вызывать остающийся.

Группа фиксирована — это всё `S` целиком, выбирать подмножество нельзя («если нельзя удалить все подозрительные, не удаляем ни одного»). Поэтому задача не оптимизационная: достаточно построить `S` и одной проверкой решить, удалять его или нет.

### Алгоритм

1. Список смежности по исходящим рёбрам.
2. Обход от `k`, отмечающий `suspicious[v]` — это `S`.
3. Один проход по `invocations`: если нашлось ребро `a → b` с `!suspicious[a] && suspicious[b]`, группу удалять нельзя — обнуляем `suspicious`.
4. Ответ — вершины с `!suspicious[i]`.

Второй обход графа на шаге 3 не нужен: условие проверяется локально на каждом ребре, а `suspicious` к этому моменту уже посчитан полностью.

Шаги 3 и 4 не разветвляются на два разных построения ответа: «удалять нельзя» и «удалять нечего» — это один и тот же ответ `0..n-1`, поэтому запрет выражается обнулением `suspicious`.

Обход итеративный: при `n = 10^5` цепочка вызовов может быть линейной, и рекурсия на таком графе даёт глубину порядка `n`. Порядок посещения не важен — нужно только само множество достижимых вершин, — поэтому вместо стека с push/pop берётся очередь-курсор: слайс `queue` только растёт, а `i` пробегает по нему, попутно давая готовый список вершин `S`.

Цикл внутри `S` (пример 3) обрабатывается сам собой — отметка `suspicious` до укладки в очередь не даёт зайти в вершину дважды. Ребро, входящее в `S` из цикла снаружи (`n = 4, k = 0, [[0,1],[1,0],[3,1]]`), ловится шагом 3 как обычное.

### Хранение графа

Список смежности — не `[][]int`, а CSR (compressed sparse row): две плоские таблицы, где соседи `v` лежат в `targets[starts[v]:starts[v+1]]`. `[][]int` на `n = 10^5` — это `n` независимых `append`-слайсов, то есть порядка `n` аллокаций и заголовок в 24 байта на вершину; CSR — две аллокации и последовательная память под обход.

Строится в три прохода без временного массива курсоров:

1. `starts[a+1]++` по каждому ребру — исходящие степени со сдвигом на единицу;
2. префиксные суммы `starts[v+1] += starts[v]` — границы блоков;
3. раскладка рёбер, где `starts[a]` служит курсором и постинкрементится.

После шага 3 каждый `starts[v]` продвинулся ровно на длину своего блока, то есть стал прежним `starts[v+1]`; сдвиг массива вправо на одну позицию возвращает начала блоков на место.

Замер на цепочке `n = 10^5` (M2 Pro): 1.89 мс и 100 003 аллокации у `[][]int` против 0.93 мс и 5 аллокаций у CSR.

- Time: O(n + m), где `m = len(invocations)` — построение графа, обход и финальные проходы линейны
- Space: O(n + m) — `starts` и `targets` плюс `suspicious`, `queue` и ответ на O(n)
