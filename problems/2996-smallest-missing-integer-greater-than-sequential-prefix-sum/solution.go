// Package smallestmissingintegergreaterthansequentialprefixsum — решение leetcode 2996. Smallest Missing Integer Greater Than Sequential Prefix Sum.
package smallestmissingintegergreaterthansequentialprefixsum

// maxValue — верхняя граница значений по условию.
const maxValue = 50

func missingInteger(nums []int) int {
	prefix := nums[0]

	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1]+1 {
			break
		}

		prefix += nums[i]
	}

	// различных значений не больше maxValue, поэтому свободный слот найдётся
	// не позже maxValue+1 — за этой границей проверять уже нечего
	vals := make([]bool, maxValue+2)
	for _, v := range nums {
		vals[v] = true
	}

	for prefix < len(vals) && vals[prefix] {
		prefix++
	}

	return prefix
}
