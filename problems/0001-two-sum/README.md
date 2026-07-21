# 1. Two Sum

- Difficulty: easy
- Link: https://leetcode.com/problems/two-sum/
- Topics: Array, Hash Table
- Acceptance: 57.8%

## Statement

Given an array of integers `nums` and an integer `target`, return _indices of the two numbers such that they add up to `target`_.

You may assume that each input would have **_exactly_ one solution**, and you may not use the _same_ element twice.

You can return the answer in any order.

**Example 1:**

```
Input: nums = [2,7,11,15], target = 9
Output: [0,1]
Explanation: Because nums[0] + nums[1] == 9, we return [0, 1].
```

**Example 2:**

```
Input: nums = [3,2,4], target = 6
Output: [1,2]
```

**Example 3:**

```
Input: nums = [3,3], target = 6
Output: [0,1]
```

**Constraints:**

- `2 <= nums.length <= 10^4`
- `-10^9 <= nums[i] <= 10^9`
- `-10^9 <= target <= 10^9`
- **Only one valid answer exists.**

**Follow-up:** Can you come up with an algorithm that is less than `O(n^2)` time complexity?

## Hints

<details><summary>Hint 1</summary>

A really brute force way would be to search for all possible pairs of numbers but that would be too slow. Again, it's best to try out brute force solutions just for completeness. It is from these brute force solutions that you can come up with optimizations.

</details>

<details><summary>Hint 2</summary>

So, if we fix one of the numbers, say `x`, we have to scan the entire array to find the next number `y` which is `value - x` where value is the input parameter. Can we change our array somehow so that this search becomes faster?

</details>

<details><summary>Hint 3</summary>

The second train of thought is, without changing the array, can we use additional space somehow? Like maybe a hash map to speed up the search?

</details>

## Similar

- [3Sum](https://leetcode.com/problems/3sum/) (medium)
- [4Sum](https://leetcode.com/problems/4sum/) (medium)
- [Two Sum II - Input Array Is Sorted](https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/) (medium)
- [Two Sum III - Data structure design](https://leetcode.com/problems/two-sum-iii-data-structure-design/) (easy)
- [Subarray Sum Equals K](https://leetcode.com/problems/subarray-sum-equals-k/) (medium)
- [Two Sum IV - Input is a BST](https://leetcode.com/problems/two-sum-iv-input-is-a-bst/) (easy)
- [Two Sum Less Than K](https://leetcode.com/problems/two-sum-less-than-k/) (easy)
- [Max Number of K-Sum Pairs](https://leetcode.com/problems/max-number-of-k-sum-pairs/) (medium)
- [Count Good Meals](https://leetcode.com/problems/count-good-meals/) (medium)
- [Count Number of Pairs With Absolute Difference K](https://leetcode.com/problems/count-number-of-pairs-with-absolute-difference-k/) (easy)
- [Number of Pairs of Strings With Concatenation Equal to Target](https://leetcode.com/problems/number-of-pairs-of-strings-with-concatenation-equal-to-target/) (medium)
- [Find All K-Distant Indices in an Array](https://leetcode.com/problems/find-all-k-distant-indices-in-an-array/) (easy)
- [First Letter to Appear Twice](https://leetcode.com/problems/first-letter-to-appear-twice/) (easy)
- [Number of Excellent Pairs](https://leetcode.com/problems/number-of-excellent-pairs/) (hard)
- [Number of Arithmetic Triplets](https://leetcode.com/problems/number-of-arithmetic-triplets/) (easy)
- [Node With Highest Edge Score](https://leetcode.com/problems/node-with-highest-edge-score/) (medium)
- [Check Distances Between Same Letters](https://leetcode.com/problems/check-distances-between-same-letters/) (easy)
- [Find Subarrays With Equal Sum](https://leetcode.com/problems/find-subarrays-with-equal-sum/) (easy)
- [Largest Positive Integer That Exists With Its Negative](https://leetcode.com/problems/largest-positive-integer-that-exists-with-its-negative/) (easy)
- [Number of Distinct Averages](https://leetcode.com/problems/number-of-distinct-averages/) (easy)
- [Count Pairs Whose Sum is Less than Target](https://leetcode.com/problems/count-pairs-whose-sum-is-less-than-target/) (easy)

## Idea

TODO

- Time: O(?)
- Space: O(?)
