# 1882. Process Tasks Using Servers

- Difficulty: medium
- Link: https://leetcode.com/problems/process-tasks-using-servers/
- Topics: Array, Heap (Priority Queue)
- Acceptance: 48.8%

## Statement

You are given two **0-indexed** integer arrays `servers` and `tasks` of lengths `n`​​​​​​ and `m`​​​​​​ respectively. `servers[i]` is the **weight** of the `i^​​​​​​th`​​​​ server, and `tasks[j]` is the **time needed** to process the `j^​​​​​​th`​​​​ task **in seconds**.

Tasks are assigned to the servers using a **task queue**. Initially, all servers are free, and the queue is **empty**.

At second `j`, the `j^th` task is **inserted** into the queue (starting with the `0^th` task being inserted at second `0`). As long as there are free servers and the queue is not empty, the task in the front of the queue will be assigned to a free server with the **smallest weight**, and in case of a tie, it is assigned to a free server with the **smallest index**.

If there are no free servers and the queue is not empty, we wait until a server becomes free and immediately assign the next task. If multiple servers become free at the same time, then multiple tasks from the queue will be assigned **in order of insertion** following the weight and index priorities above.

A server that is assigned task `j` at second `t` will be free again at second `t + tasks[j]`.

Build an array `ans`​​​​ of length `m`, where `ans[j]` is the **index** of the server the `j^​​​​​​th` task will be assigned to.

Return _the array_ `ans`​​​​.

**Example 1:**

```
Input: servers = [3,3,2], tasks = [1,2,3,2,1,2]
Output: [2,2,0,2,1,2]
Explanation: Events in chronological order go as follows:
- At second 0, task 0 is added and processed using server 2 until second 1.
- At second 1, server 2 becomes free. Task 1 is added and processed using server 2 until second 3.
- At second 2, task 2 is added and processed using server 0 until second 5.
- At second 3, server 2 becomes free. Task 3 is added and processed using server 2 until second 5.
- At second 4, task 4 is added and processed using server 1 until second 5.
- At second 5, all servers become free. Task 5 is added and processed using server 2 until second 7.
```

**Example 2:**

```
Input: servers = [5,1,4,3,2], tasks = [2,1,2,4,5,2,1]
Output: [1,4,1,4,1,3,2]
Explanation: Events in chronological order go as follows:
- At second 0, task 0 is added and processed using server 1 until second 2.
- At second 1, task 1 is added and processed using server 4 until second 2.
- At second 2, servers 1 and 4 become free. Task 2 is added and processed using server 1 until second 4.
- At second 3, task 3 is added and processed using server 4 until second 7.
- At second 4, server 1 becomes free. Task 4 is added and processed using server 1 until second 9.
- At second 5, task 5 is added and processed using server 3 until second 7.
- At second 6, task 6 is added and processed using server 2 until second 7.
```

**Constraints:**

- `servers.length == n`
- `tasks.length == m`
- `1 <= n, m <= 2 * 10^5`
- `1 <= servers[i], tasks[j] <= 2 * 10^5`

## Hints

<details><summary>Hint 1</summary>

You can maintain a Heap of available Servers and a Heap of unavailable servers

</details>

<details><summary>Hint 2</summary>

Note that the tasks will be processed in the input order so you just need to find the x-th server that will be available according to the rules

</details>

## Similar

- [Parallel Courses III](https://leetcode.com/problems/parallel-courses-iii/) (hard)

## Idea

Задачи назначаются строго в порядке очереди, поэтому «кому достанется задача `j`» определяется только моментом её назначения и состоянием серверов к этому моменту. Значит, можно идти по задачам, а не по секундам. Посекундная симуляция не годится: один сервер и `m` длинных задач дают `4 * 10^10` секунд. Поиск сервера перебором тоже: `n * m` даёт те же `4 * 10^10` сравнений.

Это событийная симуляция: нужен минимум по двум разным порядкам, и оба множества всё время меняются. Отсюда две кучи:

- `free` — свободные серверы, порядок `(вес, индекс)`: вершина и есть сервер, которого требует условие;
- `busy` — занятые, порядок `(время освобождения, индекс)`: вершина — ближайшее событие.

Перед назначением задачи `j`:

1. `now = max(now, j)`: задача не может уйти раньше, чем попала в очередь, и раньше предыдущей задачи.
1. Если `free` пуста, время прыгает на вершину `busy`. Секунды между событиями ничего не меняют.
1. Все серверы с `freeAt <= now` переходят в `free`. Именно все, а не только один: свежеосвободившийся сервер может оказаться легче текущей вершины `free`, и без этого шага задача уйдёт не туда.
1. Вершина `free` получает задачу и уходит в `busy` с `freeAt = now + tasks[j]`.

Случай «несколько серверов освободились разом» отдельной обработки не требует: все они попадают в `free`, а задачи из очереди забирают их по одной в порядке `(вес, индекс)`.

В шаге 2 не нужен `max` с текущим `now`: после назначения все в `busy` освобождаются позже прошлого `now` (длительности `>= 1`, а всё `<= now` уже выпущено), и новая задача приходит не позже `now + 1`. Значит, вершина `busy` не меньше нового `now`, и время назад не идёт.

Ключ упакован в один `int`: `приоритет << 18 | индекс`. Индекс и вес меньше `2^18`, время освобождения не больше `2 * 10^5 + 4 * 10^10 < 2^36`, всего 54 бита. Тест на предельных значениях сверяет это с вариантом без упаковки.

Какая куча лучше, решил бенчмарк: `benchstat`, `-count=10`, Apple M2 Pro, один бинарник, `n = m = 2 * 10^5` (в `few-servers` `n = 16`).

| Вход                | своя куча, упакованный `int` | `container/heap`, упакованный `int` | `container/heap`, структура |
| ------------------- | ---------------------------- | ----------------------------------- | --------------------------- |
| `deep-heaps`        | 24.1 мс                      | 59.3 мс (+146%)                     | 66.5 мс (+176%)             |
| `short-tasks`       | 9.8 мс                       | 42.4 мс ± 21% (+331%)               | 46.9 мс (+377%)             |
| `few-servers`       | 4.8 мс                       | 16.4 мс (+239%)                     | 18.6 мс (+284%)             |
| аллокаций на прогон | 3                            | 600–800 тыс.                        | 600–800 тыс.                |

Все разницы со своей кучей значимы при `p = 0.000`. Главный выигрыш — отказ от `container/heap`: каждый `Push` и `Pop` упаковывает значение в `any` и аллоцирует, это 4 аллокации на задачу с освобождением. Упаковка ключа сама по себе даёт меньше: +12–13% к времени у структуры на `deep-heaps` и `few-servers` при `p = 0.000`, а на `short-tasks` разница незначима (`p = 0.063`). Память при этом различается на 80–85%.

- Time: O((n + m) log n). Построение `free` стоит O(n). На каждую задачу по одному извлечению из `free` и вставке в `busy`, а выпусков из `busy` суммарно не больше `m`, потому что каждый оплачен более ранней вставкой.
- Space: O(n) сверх ответа на `m` элементов.
