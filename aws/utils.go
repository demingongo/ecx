package aws

import (
	"encoding/json"
	"os/exec"
	"time"
)

func execAWS[T any](args []string, resp *T) ([]byte, error) {
	cmd := exec.Command("aws", args...)
	stdout, err := cmd.Output()
	if err != nil {
		return stdout, err
	}
	if len(stdout) > 0 {
		err = json.Unmarshal(stdout, resp)
	}
	return stdout, err
}

func sleep(seconds time.Duration) {
	time.Sleep(seconds * time.Second)
}
