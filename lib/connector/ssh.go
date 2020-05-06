package connector

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/schidstorm/go-remote/lib"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
	lssh "golang.org/x/crypto/ssh"
)

type Ssh struct {
	Address   string
	Config    *ssh.ClientConfig
	sshClient *ssh.Client
	os        lib.OS
}

func SshNormal(address string, user string) (*Ssh, error) {
	return SshWithPassphrase(address, "", user)
}

func SshWithPassphrase(address string, passphrase string, user string) (*Ssh, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	privateKeyContent, err := ioutil.ReadFile(path.Join(homeDir, ".ssh", "id_rsa"))
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if passphrase == "" {
		signer, err = ssh.ParsePrivateKey(privateKeyContent)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKeyContent, []byte(passphrase))
	}

	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return &Ssh{
		Address: address,
		Config:  config,
	}, nil
}

func (ssh *Ssh) GuessOs() lib.OS {
	if ssh.doesCommandRun("cmd /c") {
		return lib.Windows
	}

	return lib.Linux
}

func (ssh *Ssh) doesCommandRun(command string) bool {
	session, err := ssh.sshClient.NewSession()
	if err == nil {
		return false
	}

	defer session.Close()
	err = session.Run(command)
	if err != nil {
		return false
	}

	return true
}

// Connect creates a ssh client connection
func (ssh *Ssh) Connect() (io.ReadWriteCloser, error) {
	log.Println("Dialing Ssh " + ssh.Address)
	sshClient, err := lssh.Dial("tcp", ssh.Address, ssh.Config)
	ssh.sshClient = sshClient
	if err != nil {
		return nil, err
	}

	ssh.os = ssh.GuessOs()

	var localFileName string
	if ssh.os == lib.Linux {
		log.Println("Detecting linux")
		localFileName = strings.ReplaceAll(path.Base(os.Args[0]), ".exe", "")
	} else {
		log.Println("Detecting windows")
		localFileName = strings.ReplaceAll(path.Base(os.Args[0]), ".exe", "") + ".exe"
	}
	localPath := path.Join(path.Dir(os.Args[0]), localFileName)

	remotePath, err := ssh.transferBinary(localPath)
	if err != nil {
		return nil, err
	}

	rwc, err := ssh.executeBinary(remotePath, []string{"--listen"})
	if err != nil {
		return nil, err
	}

	return rwc, nil
}

func (ssh *Ssh) transferBinary(localPath string) (string, error) {
	remotePath := path.Join(os.TempDir(), "go-remote-"+lib.VERSION)
	log.Println("Transfer binary to " + remotePath)
	fileCopySession, err := ssh.sshClient.NewSession()
	if err != nil {
		return "", err
	}

	err = scp.CopyPath(localPath, remotePath, fileCopySession)
	if err != nil {
		return "", err
	}
	fileCopySession.Close()

	if ssh.os == lib.Linux {
		chmodSession, err := ssh.sshClient.NewSession()
		if err != nil {
			return "", err
		}
		defer chmodSession.Close()
		log.Println("Make file executable")
		err = chmodSession.Run("chmod +x " + remotePath)
		if err != nil {
			return "", err
		}
	}

	return remotePath, nil
}

func (ssh *Ssh) executeBinary(remotePath string, arguments []string) (io.ReadWriteCloser, error) {
	serverSession, err := ssh.sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	log.Println("Attach to remote connection")
	stdoutPipe, err := serverSession.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderrPipe, err := serverSession.StderrPipe()
	if err != nil {
		return nil, err
	}

	go io.Copy(os.Stderr, stderrPipe)

	stdinPipe, err := serverSession.StdinPipe()
	if err != nil {
		return nil, err
	}

	err = serverSession.Start(remotePath + " " + strings.Join(arguments, " "))
	if err != nil {
		return nil, err
	}

	return &lib.StdInOutPipe{
		stdoutPipe,
		stdinPipe,
		serverSession,
	}, nil
}
