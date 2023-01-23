package ssh

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

type SSHAccount struct {
	Username string
	Password string
}

func NewSSHAccount(username, password string) *SSHAccount {
	return &SSHAccount{
		Username: username,
		Password: password,
	}
}

func (s *SSHAccount) Connector(host, port string) (*ssh.Client, error) {
	authMethods := []ssh.AuthMethod{}

	keyboardInteractiveChallege := func(username, instruction string, questions []string, echos []bool) ([]string, error) {
		if len(questions) == 0 {
			return []string{}, nil
		}

		return []string{s.Password}, nil

	}

	authMethods = append(authMethods, ssh.KeyboardInteractive(keyboardInteractiveChallege))
	authMethods = append(authMethods, ssh.Password(s.Password))

	sshConfig := &ssh.ClientConfig{
		User:            s.Username,
		Auth:            authMethods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), sshConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}
