package probe

import (
	"bufio"
	"context"
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
			if !t.ping() {
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

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "pong") {
				return true, nil
			}
		}

		if err := cmd.Wait(); err != nil {
			return false, nil
		}

		return false, nil
	})

	if err != nil {
		return false
	}

	return result.(bool)
}
