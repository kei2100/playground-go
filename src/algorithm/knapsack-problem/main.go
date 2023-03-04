package main

type item struct {
	weight int
	value  int
}

func main() {
	items := []item{
		{weight: 3, value: 100},
		{weight: 6, value: 210},
		{weight: 4, value: 130},
		{weight: 2, value: 57},
	}
	result := knapsack(10, items...)
	println(result)
	// 340
}

func knapsack(maxWeight int, items ...item) int {
	n, w := len(items), maxWeight
	dp := make([][]int, n)
	var max int
	for i := 0; i < n; i++ {
		item := items[i]
		dp[i] = make([]int, w+1)
		for j := 1; j <= w; j++ {
			if i == 0 {
				if j == item.weight {
					dp[i][j] = item.value
				} else {
					dp[i][j] = -1
				}
			} else {
				dp[i][j] = dp[i-1][j]
				v := -1
				if jj := j - item.weight; jj > -1 {
					if vv := dp[i-1][jj]; vv > -1 {
						v = vv + item.value
					}
				}
				if dp[i][j] < v {
					dp[i][j] = v
				}
			}
			if max < dp[i][j] {
				max = dp[i][j]
			}
		}
	}
	return max
}
