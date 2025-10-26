package pyai

import (
	"bufio"
	"os/exec"
	"sync"
)

type PyAI struct {
	cmd    *exec.Cmd
	stdin  *bufio.Writer
	stdout *bufio.Reader
	mu     sync.Mutex
}

func NewPythonAI() (*PyAI, error) {
	cmd := exec.Command("python3", "path/to/your_script.py")

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &PyAI{
		cmd:    cmd,
		stdin:  bufio.NewWriter(stdinPipe),
		stdout: bufio.NewReader(stdoutPipe),
	}, nil
}

// Send request and get response
func (p *PyAI) Query(request string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Write request + newline
	_, err := p.stdin.WriteString(request + "\n")
	if err != nil {
		return "", err
	}
	err = p.stdin.Flush()
	if err != nil {
		return "", err
	}

	// Read response line (assuming response is single line)
	response, err := p.stdout.ReadString('\n')
	if err != nil {
		return "", err
	}

	return response[:len(response)-1], nil // strip newline
}
