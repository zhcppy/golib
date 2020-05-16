package ssh

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func NewClient(IP string, privKey []byte) (*ssh.Client, *ssh.Session, error) {
	auth, err := GetAuth(privKey)
	if err != nil {
		return nil, nil, err
	}

	cfg := &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            "ubuntu",
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		BannerCallback:  ssh.BannerDisplayStderr(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", IP, 22), cfg)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}
	return client, session, nil
}

func GetAuth(pemKey []byte) (ssh.AuthMethod, error) {
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(pemKey)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}
