package overlay

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverlay(t *testing.T) {
	t.Parallel()
	// go run
	out1, err := goRun()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello", out1)
	// go run (using -overlay option)
	if err := genOverlayJSON(); err != nil {
		t.Fatal(err)
	}
	out2, err := goRun("-overlay", filepath.Join(curDir(), "overlay.json"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Goodbye", out2)
}

func genOverlayJSON() error {
	data := map[string]map[string]string{
		"Replace": {
			filepath.Join(curDir(), "main", "fn", "fn.go"): filepath.Join(curDir(), "main", "fn_replace", "fn_replace.go"),
		},
	}
	f, err := os.Create(filepath.Join(curDir(), "overlay.json"))
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(data); err != nil {
		return err
	}
	return nil
}

func goRun(opts ...string) (string, error) {
	mainDir := filepath.Join(curDir(), "main")
	args := []string{"run"}
	args = append(args, opts...)
	args = append(args, mainDir)
	cmd := exec.Command("go", args...)
	cmd.Dir = mainDir
	out := bytes.NewBuffer(nil)
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func curDir() string {
	_, b, _, _ := runtime.Caller(0)
	curDir := filepath.Dir(b)
	return curDir
}
