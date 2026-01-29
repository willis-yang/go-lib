package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHManager struct {
	client  *ssh.Client
	session *ssh.Session
	stdout  io.Reader
	stderr  io.Reader
	stdin   io.WriteCloser
}

func loadPrivateKey(keyData []byte, passphrase string) (ssh.AuthMethod, error) {
	var signer ssh.Signer
	var err error
	if passphrase == "" {
		signer, err = ssh.ParsePrivateKey(keyData)
		if err != nil {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(""))
		}
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(passphrase))
	}

	if err != nil {
		return nil, fmt.Errorf("decode key error: %w", err)
	}

	return ssh.PublicKeys(signer), nil
}

func NewSSHManager(host, port, user string, keyData []byte, timeout time.Duration) (*SSHManager, error) {
	authMethod, err := loadPrivateKey(keyData, "")
	if err != nil {
		return nil, fmt.Errorf("load private key error: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, fmt.Errorf("connect ssh failed: %w", err)
	}

	return &SSHManager{client: client}, nil
}

// ExecuteCommand
// 设置超时时间
func (m *SSHManager) ExecuteCommand(cmd string, timeout time.Duration) (string, error) {
	var err error
	m.session, err = m.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("create session failed: %w", err)
	}
	defer func() {
		if m.session != nil {
			m.session.Close()
			m.session = nil
		}
	}()
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 定义结果通道
	type result struct {
		output []byte
		err    error
	}
	ch := make(chan result, 1)

	go func() {
		output, err := m.session.CombinedOutput(cmd)
		ch <- result{output, err}
	}()

	// 等待命令完成或超时
	select {
	case <-ctx.Done():
		// 超时触发，关闭会话终止命令
		m.session.Close()
		return "", fmt.Errorf("command execution timed out after %v", timeout)
	case res := <-ch:
		if res.err != nil {
			return "", fmt.Errorf("execute command failed: %w", res.err)
		}
		return string(res.output), nil
	}
}

// ExecuteCommandWithStream
func (m *SSHManager) ExecuteCommandWithStream(cmd string) error {
	var err error
	m.session, err = m.client.NewSession()
	if err != nil {
		return fmt.Errorf("create session failed: %w", err)
	}
	m.stdout, err = m.session.StdoutPipe()
	if err != nil {
		m.Close()
		return fmt.Errorf("get stdout failed: %w", err)
	}

	m.stderr, err = m.session.StderrPipe()
	if err != nil {
		m.Close()
		return fmt.Errorf("get stderr failed: %w", err)
	}

	m.stdin, err = m.session.StdinPipe()
	if err != nil {
		m.Close()
		return fmt.Errorf("get stdin failed: %w", err)
	}

	if err := m.session.Start(cmd); err != nil {
		m.Close()
		return fmt.Errorf("start command failed: %w", err)
	}

	return nil
}

func (m *SSHManager) Close() error {
	var errs []error
	if m.stdin != nil {
		if err := m.stdin.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭 stdin 失败: %w", err))
		}
		m.stdin = nil
	}
	if m.session != nil {
		if err := m.session.Signal(ssh.SIGTERM); err != nil && err != io.EOF {
			errs = append(errs, fmt.Errorf("发送 SIGTERM 失败: %w", err))
		}
		if err := m.session.Wait(); err != nil && err != io.EOF {
			var exitErr *ssh.ExitError
			if !errors.As(err, &exitErr) || exitErr.ExitStatus() != 0 {
				errs = append(errs, fmt.Errorf("等待会话结束失败: %w", err))
			}
		}
		if err := m.session.Close(); err != nil && err != io.EOF {
			errs = append(errs, fmt.Errorf("关闭会话失败: %w", err))
		}
		m.session = nil
	}
	if m.client != nil {
		if err := m.client.Close(); err != nil && err != io.EOF {
			errs = append(errs, fmt.Errorf("关闭客户端失败: %w", err))
		}
		m.client = nil
	}
	m.stdout = nil
	m.stderr = nil
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
