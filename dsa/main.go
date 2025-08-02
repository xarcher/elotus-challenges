package main

import "fmt"

// Solutions Gray Code
// 1. Generate Gray codes using the formula: G(n) = n ^ (n >> 1).
// 2. For n bits, the Gray code sequence has 2^n elements.
// 3. The result is a slice of integers representing the Gray code sequence.
// Time Complexity: O(2^n) for generating the sequence.
// Space Complexity: O(2^n) for storing the sequence.
func grayCode(n int) []int {
	if n <= 0 {
		return []int{}
	}

	result := make([]int, 1<<n)
	for i := 0; i < (1 << n); i++ {
		result[i] = i ^ (i >> 1)
	}

	return result
}

// Solutions Longest Common Subarray
//  1. Iterate through every possible starting index pair (i, j).
//  2. For each pair, extend the subarray while elements are equal.
//  3. Track the maximum length found.
//
// Time Complexity: O(n * m * k) where
//
//	n = len(nums1), m = len(nums2), k = average length of matching subarrays.
//
// Space Complexity: O(1)
func findLength(nums1 []int, nums2 []int) int {
	maxLen := 0
	for i := 0; i < len(nums1); i++ {
		for j := 0; j < len(nums2); j++ {
			k := 0
			for i+k < len(nums1) && j+k < len(nums2) && nums1[i+k] == nums2[j+k] {
				k++
			}

			if k > maxLen {
				maxLen = k
			}
		}
	}

	return maxLen
}

// Solutions Sum of Distances in Tree
// 1. Use DFS to calculate the size of each subtree and the sum of distances from the root.
// 2. For each node, calculate the distance to all its children and propagate the results up the tree.
// 3. The result is a slice where each index corresponds to a node and contains the sum of distances from that node to all other nodes.
// Time Complexity: O(n) for traversing the tree.
// Space Complexity: O(n) for storing the size of subtrees and distances.
func sumOfDistancesInTree(n int, edges [][]int) []int {
	tree := make([][]int, n)
	for _, edge := range edges {
		a, b := edge[0], edge[1]
		tree[a] = append(tree[a], b)
		tree[b] = append(tree[b], a)
	}

	res := make([]int, n)
	count := make([]int, n)
	for i := 0; i < n; i++ {
		count[i] = 1
	}

	var dfs1 func(node, parent int)
	dfs1 = func(node, parent int) {
		for _, child := range tree[node] {
			if child == parent {
				continue
			}
			dfs1(child, node)
			count[node] += count[child]
			res[node] += res[child] + count[child]
		}
	}

	var dfs2 func(node, parent int)
	dfs2 = func(node, parent int) {
		for _, child := range tree[node] {
			if child == parent {
				continue
			}
			res[child] = res[node] - count[child] + (n - count[child])
			dfs2(child, node)
		}
	}

	dfs1(0, -1)
	dfs2(0, -1)
	return res
}

func main() {
	// Example usage
	n := 2
	grayCodes := grayCode(n)
	fmt.Printf("gray codes %v\n", grayCodes)

	nums1 := []int{1, 2, 3, 2, 1}
	nums2 := []int{3, 2, 1, 4, 7}
	length := findLength(nums1, nums2)
	fmt.Printf("Length of the longest common subarray: %d\n", length)

	nums1 = []int{0, 0, 0, 0, 0}
	nums2 = []int{0, 0, 0, 0, 0}
	length = findLength(nums1, nums2)
	fmt.Printf("Length of the longest common subarray: %d\n", length)

	nTree := 6
	edges := [][]int{{0, 1}, {0, 2}, {2, 3}, {2, 4}, {2, 5}}
	distances := sumOfDistancesInTree(nTree, edges)
	fmt.Printf("Sum of distances in tree: %v\n", distances)

	nTree = 1
	edges = [][]int{}
	distances = sumOfDistancesInTree(nTree, edges)
	fmt.Printf("Sum of distances in tree: %v\n", distances)

	nTree = 2
	edges = [][]int{{1, 0}}
	distances = sumOfDistancesInTree(nTree, edges)
	fmt.Printf("Sum of distances in tree: %v\n", distances)
}
