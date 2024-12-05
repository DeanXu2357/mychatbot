package probe

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sync/singleflight"
)

type TailscaleNotifier struct {
	target string
	group  singleflight.Group
}

func NewTailscaleNotifier(target string) *TailscaleNotifier {
	return &TailscaleNotifier{target: target}
}

func (t *TailscaleNotifier) NotOK() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		for {
			ping := t.ping()
			if !ping {
				ch <- struct{}{}
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return ch
}

func (t *TailscaleNotifier) OK() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		for {
			if t.ping() {
				ch <- struct{}{}
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return ch
}

func (t *TailscaleNotifier) ping() bool {
	result, err, _ := t.group.Do("ping", func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "tailscale", "ping", t.target)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return false, nil
		}

		if err := cmd.Start(); err != nil {
			return false, nil
		}

		var ok bool
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "pong") {
				ok = true
				break
			}
		}

		if err := cmd.Wait(); err != nil {
			if _, isSignalKilled := err.(*exec.ExitError); isSignalKilled {
				return ok, nil
			}

			fmt.Println("err:", err)
			return ok, err
		}

		return ok, nil
	})

	if err != nil {
		return false
	}

	return result.(bool)
}
