package main

func main() {
	println(do(10, 7, 5, 3))
	// true
	println(do(6, 9, 7))
	// false
	println(do(19, 7, 18, 5, 4, 8))
	// true
	seq := []int{3590, 1260, 2560, 510, 1780, 2710, 120, 610, 2410, 2620, 1250, 1910, 50, 4130, 2760, 190, 720, 1560, 2590, 2400, 2090, 3590, 650, 4320, 4420}
	println(do(50000, seq...))
	println(do(50610, seq...))
	println(do(50190, seq...))
	println(do(50800, seq...))
	// true
	println(do(50600, seq...))
	println(do(50790, seq...))
	// false
}

func do(want int, values ...int) bool {
	dp := make([][]bool, len(values))
	for i := 0; i < len(values); i++ {
		dp[i] = make([]bool, want+1)
		v := values[i]
		for j := 1; j <= want; j++ {
			if j == v {
				dp[i][j] = true
			} else if i > 0 {
				dp[i][j] = dp[i-1][j]
				if jj := j - v; jj > 0 {
					if dp[i-1][jj] {
						dp[i][j] = true
					}
				}
			}
		}
		if dp[i][want] {
			return true
		}
	}
	return false
}
