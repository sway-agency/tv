package tv

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	commandTimeout = 15 * time.Second
	maxOutputBytes = 16 << 20
)

var ErrCommandFailed = errors.New("the TV did not accept the command")

type ADB struct {
	path   string
	serial string
	mu     sync.Mutex
}

func NewADB(path string, serial string) *ADB {
	return &ADB{path: path, serial: serial}
}

func (a *ADB) run(ctx context.Context, args ...string) ([]byte, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()
	command := exec.CommandContext(ctx, a.path, append([]string{"-s", a.serial}, args...)...)
	command.Env = []string{"HOME=" + os.Getenv("HOME")}
	stdout := &limitedBuffer{limit: maxOutputBytes}
	command.Stdout = stdout
	err := command.Run()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCommandFailed, err)
	}
	if stdout.overflow {
		return nil, fmt.Errorf("%w: output too large", ErrCommandFailed)
	}
	return stdout.Bytes(), nil
}

func (a *ADB) shell(ctx context.Context, args ...string) (string, error) {
	output, err := a.run(ctx, append([]string{"shell"}, args...)...)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

type limitedBuffer struct {
	bytes.Buffer
	limit    int
	overflow bool
}

func (b *limitedBuffer) Write(data []byte) (int, error) {
	if b.Len()+len(data) > b.limit {
		b.overflow = true
		return len(data), nil
	}
	return b.Buffer.Write(data)
}
