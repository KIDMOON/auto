package cmd

import (
	ssh_server "awesomeProject/src"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
	"time"
)

var flag string

var tomcat = &cobra.Command{
	Use:          "tomcat",
	Long:         `使用tomcat启动`,
	Short:        "tomcat操作，具体tomcat -h查看",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("tomcatStart")
		connection := ssh_server.SSh_Connection(&Auto)
		defer connection.Close()
		if flag != "" {
			stpTomcat(connection)
		} else {
			startTomcat(connection)
		}
	}}

func startTomcat(connection *ssh.Client) {
	if Auto.TomcatPath == "" {
		return
	}
	stpTomcat(connection)
	ssh_server.OpenSessionRunCmd(connection, Auto.TomcatPath+"/bin/startup.sh")
	log.Println("启动成功,休眠60秒后检测是否启动成功进程")
	time.Sleep(time.Duration(60) * time.Second)
	cmd := "ps -ef |grep " + Auto.TomcatPath + "/" + " |grep -v \"grep\"|awk '{print $2}'"
	buff := ssh_server.OpenSessionRunCmdString(connection, cmd)
	if buff != "" {
		log.Println("启动成功,进程", buff)
	} else {
		buff := ssh_server.OpenSessionRunCmdString(connection, "cat "+Auto.TomcatPath+"/logs/catalina.out  tail -n 1000")
		log.Println("打印错误日志")
		log.Println(buff)
	}

}

func stpTomcat(connection *ssh.Client) {
	cmd := "ps -ef |grep " + Auto.TomcatPath + "/" + " |grep -v \"grep\"|awk '{print $2}'"
	buff := ssh_server.OpenSessionRunCmdString(connection, cmd)
	if buff != "" {
		ssh_server.OpenSessionRunCmd(connection, "kill -9 "+buff)
	}
}

func cpProfile(connection *ssh.Client) {
	if Auto.ProfileAddress == "" {
		return
	}
	profileAddress := Auto.ProfileAddress
	if profileAddress == "" {
		return
	}
	if strings.HasPrefix(profileAddress, "@local:") {
		profileAddress = strings.Replace(profileAddress, "@local:", "", 1)
		cloneFile(connection, Auto.WarPathLocal+"/WEB-INF/classes/application.properties", profileAddress)
	} else if strings.HasPrefix(profileAddress, "@online:") {
		profileAddress = strings.Replace(profileAddress, "@online:", "", 1)
		ssh_server.OpenSessionRunCmdString(connection, "cp "+profileAddress+" "+Auto.ProjectPathOnLine+"/WEB-INF/classes/application.properties")
	} else if strings.HasPrefix(profileAddress, "/") {
		ssh_server.OpenSessionRunCmdString(connection, "cp "+profileAddress+" "+Auto.ProjectPathOnLine+"/WEB-INF/classes/application.properties")
	}
}
