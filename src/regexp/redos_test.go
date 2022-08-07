package regexp

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestReDoS(t *testing.T) {
	t.Parallel()
	fmt.Println("Test /a*z/.match('a'*65535 + 'b')")
	fmt.Printf("* goRegexp took %s\n", goRegexp())
	fmt.Printf("* nodeRegexp took %s\n", nodeRegexp())
}

func goRegexp() time.Duration {
	// Go の正規表現は O(mn)
	r := regexp.MustCompile(`a*z`)
	n := time.Now()
	r.MatchString(strings.Repeat("a", 65535) + "b")
	return time.Since(n)
}

func nodeRegexp() time.Duration {
	n := time.Now()
	err := exec.Command(
		"which",
		"node",
	).Run()
	if err != nil {
		fmt.Println("node not found")
		return 0
	}
	// node の正規表現はバックトラックが多発すると遅い
	// テスト時のバージョンは v16.14.0
	err = exec.Command(
		"node",
		"-e",
		"/a*z/.exec('a'.repeat(65535) + 'b') === null",
	).Run()
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return time.Since(n)
}
