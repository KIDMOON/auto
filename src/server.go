package ssh_server

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
)

func SSh_Connection(auto *Auto) *ssh.Client {

	var clientConfig = ssh.ClientConfig{
		User:            auto.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(auto.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var address = auto.Address
	client, err := ssh.Dial("tcp", address, &clientConfig)
	if err != nil {
		log.Fatal("连接失败" + auto.Address + auto.Username + auto.Password)
	}
	return client
}

func OpenSessionRunCmd(client *ssh.Client, cmd string) {
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer session.Close()
	log.Println(cmd)

	if err := session.Run(cmd); err != nil {
		log.Fatal(err.Error())
	}
}

func OpenSessionRunCmdString(client *ssh.Client, cmd string) string {
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer session.Close()
	log.Println(cmd)

	if buff, err := session.CombinedOutput(cmd); err != nil {
		log.Fatal(err.Error())

	} else {
		return string(buff)
	}
	return ""
}

func OpenFtpClient(client *ssh.Client) *sftp.Client {
	ftp, err := sftp.NewClient(client)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return ftp
}

type Auto struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	Address           string `json:"address"`
	WarPathLocal      string `json:"warPathLocal"`
	ProjectPathOnLine string `json:"projectPathOnLine"`
	TomcatPath        string `json:"tomcatPath"`
	ProfileAddress    string `json:"profileAddress"`
}
