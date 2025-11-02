package local

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/Muhammadhamd/agent-sdk-go/pkg/mcp"
)

// StdioServer wraps an MCP server running as a subprocess communicating via stdio
type StdioServer struct {
	cmd              *exec.Cmd
	stdin            io.WriteCloser
	stdout           *bufio.Scanner
	stdoutReader     io.ReadCloser
	stderr           io.ReadCloser
	base             mcp.BaseTransport
	mu               sync.Mutex
	requests         map[interface{}]chan *mcp.JSONRPCResponse
	requestIDCounter int64
}

// StdioServerConfig configures a stdio MCP server
type StdioServerConfig struct {
	Command string
	Args    []string
	Env     []string
	Dir     string
	Timeout int // timeout in seconds, 0 for no timeout
}

// NewStdioServer creates a new stdio MCP server wrapper
func NewStdioServer(config StdioServerConfig) mcp.Transport {
	cmd := exec.Command(config.Command, config.Args...)
	if len(config.Env) > 0 {
		cmd.Env = config.Env
	}
	if config.Dir != "" {
		cmd.Dir = config.Dir
	}

	return &StdioServer{
		cmd:              cmd,
		requests:         make(map[interface{}]chan *mcp.JSONRPCResponse),
		requestIDCounter: 0,
	}
}

// Connect starts the subprocess and establishes stdio communication
func (s *StdioServer) Connect(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.base.IsConnected() {
		return fmt.Errorf("server already connected")
	}

	// Environment and directory are set during NewStdioServer creation

	// Get stdin pipe
	stdin, err := s.cmd.StdinPipe()
	if err != nil {
		return mcp.NewTransportError(err)
	}
	s.stdin = stdin

	// Get stdout pipe
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return mcp.NewTransportError(err)
	}
	s.stdoutReader = stdout
	s.stdout = bufio.NewScanner(stdout)

	// Get stderr pipe (for error logging)
	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		stdin.Close()
		stdout.Close()
		return mcp.NewTransportError(err)
	}
	s.stderr = stderr

	// Start the process
	if err := s.cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		stderr.Close()
		return mcp.NewConnectionError(err)
	}

	// Start goroutine to read responses
	go s.readResponses()

	// Start goroutine to read stderr
	go s.readStderr()

	s.base.SetConnected(true)
	return nil
}

// IsConnected implements Transport interface
func (s *StdioServer) IsConnected() bool {
	return s.base.IsConnected()
}

// readResponses continuously reads JSON-RPC messages from stdout
func (s *StdioServer) readResponses() {
	for s.stdout.Scan() {
		line := s.stdout.Text()
		if line == "" {
			continue
		}

		var resp mcp.JSONRPCResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			continue // Skip invalid JSON
		}

		s.mu.Lock()
		if ch, ok := s.requests[resp.ID]; ok {
			ch <- &resp
			close(ch)
			delete(s.requests, resp.ID)
		}
		s.mu.Unlock()
	}
}

// readStderr reads error output from stderr
func (s *StdioServer) readStderr() {
	scanner := bufio.NewScanner(s.stderr)
	for scanner.Scan() {
		// Log stderr output (could be enhanced to send to logger)
		_ = scanner.Text()
	}
}

// SendRequest sends a JSON-RPC request and waits for response
func (s *StdioServer) SendRequest(ctx context.Context, req *mcp.JSONRPCRequest) (*mcp.JSONRPCResponse, error) {
	if !s.base.IsConnected() {
		return nil, fmt.Errorf("server not connected")
	}

	// Generate request ID if not provided
	if req.ID == nil {
		s.mu.Lock()
		s.requestIDCounter++
		req.ID = s.requestIDCounter
		s.mu.Unlock()
	}

	// Create response channel
	ch := make(chan *mcp.JSONRPCResponse, 1)

	s.mu.Lock()
	s.requests[req.ID] = ch
	s.mu.Unlock()

	// Marshal and send request
	reqJSON, err := json.Marshal(req)
	if err != nil {
		s.mu.Lock()
		delete(s.requests, req.ID)
		s.mu.Unlock()
		return nil, err
	}

	reqJSON = append(reqJSON, '\n')

	s.mu.Lock()
	_, err = s.stdin.Write(reqJSON)
	s.mu.Unlock()

	if err != nil {
		s.mu.Lock()
		delete(s.requests, req.ID)
		s.mu.Unlock()
		return nil, mcp.NewTransportError(err)
	}

	// Wait for response
	select {
	case resp := <-ch:
		return resp, nil
	case <-ctx.Done():
		s.mu.Lock()
		delete(s.requests, req.ID)
		s.mu.Unlock()
		return nil, ctx.Err()
	}
}

// SendNotification sends a JSON-RPC notification (no response expected)
func (s *StdioServer) SendNotification(ctx context.Context, notif *mcp.JSONRPCNotification) error {
	if !s.base.IsConnected() {
		return fmt.Errorf("server not connected")
	}

	notifJSON, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	notifJSON = append(notifJSON, '\n')

	s.mu.Lock()
	defer s.mu.Unlock()
	_, err = s.stdin.Write(notifJSON)
	if err != nil {
		return mcp.NewTransportError(err)
	}

	return nil
}

// Close stops the subprocess and cleans up
func (s *StdioServer) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.base.SetConnected(false)

	if s.stdin != nil {
		s.stdin.Close()
	}

	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Kill()
		_ = s.cmd.Wait()
	}

	return nil
}
