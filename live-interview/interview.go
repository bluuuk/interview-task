package main

import "fmt"

/*
Write a function LongestSubarrayWithDistinctEntries that takes an array and
returns the length of the longest subarray with the property that all its
elements are distinct.

Input: [1, 2, 1, 3, 1, 2, 1, 4]
Output: 3
Explanation: the solution is 3 because there is no subarray with
			more than 3 elements with unique entries and there
			is at least one subarray of size 3 with this property: [2, 1, 3].
*/

func main() {
	test := LongestSubarrayWithDistinctEntries_fast([]int{1, 2, 1, 3, 1, 2, 1, 4})
	expected := 3
	if expected == test {
		fmt.Println("Test passed")
	} else {
		fmt.Printf("%d instead of %d\n", test, expected)
	}
}

func LongestSubarrayWithDistinctEntries_fast(array []int) int {

	var current = 0
	window := make(map[int]int) // element[index]

	for i, n := range array {

		j, contains := window[n]

		if contains {
			// update the window by removing old elements
			for k, index := range window {
				if index < j {
					delete(window, k)
				}
			}

			// update the duplicating element to the current one
			window[n] = i
		} else {
			// extend the window
			window[n] = i
		}

		// update the current score
		if len(window) > current {
			current = len(window)
		}
	}

	return current
}

func LongestSubarrayWithDistinctEntries_slow(array []int) int {

	var current = 0

	for i, _ := range array {

		// build hashset to count distinct elements
		var hashset = make(map[int]bool)
		currentScore := 0

		for _, m := range array[i:] {

			// check if we have a duplicate
			if _, contains := hashset[m]; contains {
				break
			} else {
				// otherwise, put the element into the hashset and increse the counter
				hashset[m] = true
				currentScore += 1
			}
		}

		// update current
		if currentScore > current {
			current = currentScore
		}

	}

	return current
}
