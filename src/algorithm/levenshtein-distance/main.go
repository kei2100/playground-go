package main

func main() {
	println(lsd([]rune("input"), []rune("import")))
	println(lsd([]rune("kitten"), []rune("sitting")))
	println(lsd([]rune("こんにちは世界"), []rune("こんにちは世界")))
	println(lsd([]rune("こにゃにゃちわ世界"), []rune("こんにちは世界")))
	println(lsd([]rune("こにゃちわ世界"), []rune("こんにちは世界")))
	println(lsd([]rune("こんばんは世界"), []rune("こんにちは世界")))
}

func lsd(s []rune, t []rune) int {
	dp := make([][]int, len(s)+1)
	for i := 0; i <= len(s); i++ {
		dp[i] = make([]int, len(t)+1)
		for j := 0; j <= len(t); j++ {
			if i == 0 {
				dp[i][j] = j
				continue
			}
			if j == 0 {
				dp[i][j] = i
				continue
			}
			a := dp[i-1][j] + 1
			b := dp[i][j-1] + 1
			c := dp[i-1][j-1]
			if s[i-1] != t[j-1] {
				c += 1
			}
			dp[i][j] = min(a, b, c)
		}
	}
	return dp[len(s)][len(t)]
}

func min(a, b, c int) int {
	ret := a
	if ret > b {
		ret = b
	}
	if ret > c {
		ret = c
	}
	return ret
}
