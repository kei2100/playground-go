package main

import (
	"fmt"
	"sort"
)

// 重み付き区間スケジューリング問題

type job struct {
	begin  int
	end    int
	weight int
}

func main() {
	jobs := []*job{
		{begin: 0, end: 5, weight: 1},
		{begin: 4, end: 7, weight: 4},
		{begin: 3, end: 8, weight: 6},
		{begin: 5, end: 8, weight: 2},
		{begin: 6, end: 10, weight: 5},
		{begin: 9, end: 12, weight: 1},
		{begin: 11, end: 16, weight: 7},
		{begin: 14, end: 17, weight: 2},
		{begin: 13, end: 18, weight: 9},
		{begin: 15, end: 19, weight: 3},
		{begin: 16, end: 20, weight: 8},
	}
	sum := resolve(jobs)
	fmt.Println(sum)
}

func resolve(jobs []*job) int {
	sort.Slice(jobs, func(i, j int) bool {
		if jobs[i].end != jobs[j].end {
			return jobs[i].end < jobs[j].end
		}
		return jobs[i].begin < jobs[j].begin
	})
	head := &job{begin: 0, end: 0, weight: 0}
	jobs = append([]*job{head}, jobs...)
	n := len(jobs)
	dp := make([]int, n)
	for i := 0; i < n; i++ {
		if i == 0 {
			dp[i] = jobs[i].weight
		} else {
			prev := prevJobIndex(i, jobs)
			dp[i] = max(jobs[i].weight+dp[prev], dp[i-1])
		}
	}
	for i, v := range dp {
		fmt.Printf("dp[%d]: %d\n", i, v)
	}
	return dp[n-1]
}

func prevJobIndex(currentIndex int, jobs []*job) int {
	cur := jobs[currentIndex]
	for i := currentIndex - 1; i >= 0; i-- {
		prev := jobs[i]
		if cur.begin >= prev.end {
			return i
		}
	}
	panic("unexpected")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
