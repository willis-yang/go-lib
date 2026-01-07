package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
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

func loadPrivateKey(keyPath, passphrase string) (ssh.AuthMethod, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败: %w", err)
	}

	var signer ssh.Signer

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

func NewSSHManager(host, port, user, keyPath string) (*SSHManager, error) {
	authMethod, err := loadPrivateKey(keyPath, "")
	if err != nil {
		return nil, fmt.Errorf("load private key error: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, fmt.Errorf("连接失败: %w", err)
	}

	return &SSHManager{client: client}, nil
}

// ExecuteCommand
func (m *SSHManager) ExecuteCommand(cmd string) (string, error) {
	var err error

	m.session, err = m.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建会话失败: %w", err)
	}
	defer func() {
		if m.session != nil {
			m.session.Close()
			m.session = nil
		}
	}()
	output, err := m.session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("执行失败: %w", err)
	}

	return string(output), nil
}

// ExecuteCommandWithStream
func (m *SSHManager) ExecuteCommandWithStream(cmd string) error {
	var err error
	m.session, err = m.client.NewSession()
	if err != nil {
		return fmt.Errorf("创建会话失败: %w", err)
	}
	m.stdout, err = m.session.StdoutPipe()
	if err != nil {
		m.Close()
		return fmt.Errorf("获取 stdout 失败: %w", err)
	}

	m.stderr, err = m.session.StderrPipe()
	if err != nil {
		m.Close()
		return fmt.Errorf("获取 stderr 失败: %w", err)
	}

	m.stdin, err = m.session.StdinPipe()
	if err != nil {
		m.Close()
		return fmt.Errorf("获取 stdin 失败: %w", err)
	}

	if err := m.session.Start(cmd); err != nil {
		m.Close()
		return fmt.Errorf("启动命令失败: %w", err)
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

//
//func safeSSHOperation(host, port, user, keyPath string) error {
//	manager, err := NewSSHManager(host, port, user, keyPath)
//	if err != nil {
//		return err
//	}
//	defer manager.Close()
//	output, err := manager.ExecuteCommand("ls -la")
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("输出: %s\n", output)
//	return nil
//}
