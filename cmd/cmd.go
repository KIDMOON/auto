package cmd

import (
	ssh_server "awesomeProject/src"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/pkgms/go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"math/rand"
	"os"
)

var config string

func Start() {
	var cmd = &cobra.Command{
		Use:          "auto",
		Long:         `自动上传文件，自动部署启动`,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("auto Publish")
			run()
		}}
	version(cmd)
	cmd.AddCommand(tomcat)
	cmd.AddCommand(upload)
	tomcat.Flags().StringVarP(&flag, "stop", "s", "", "关闭tomcat")
	cmd.Flags().StringVarP(&config, "config", "c", "", "设置配置文件路径")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".") // 还可以在工作目录中查找配置
	if config == "" {
		str, _ := os.Getwd()
		config = str + "/config.yml"
	}
	err := server.ParseConfig(config, &Auto)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	cmd.Execute()
}

var Auto ssh_server.Auto

func run() {
	connection := ssh_server.SSh_Connection(&Auto)
	defer connection.Close()
	remoteFilePath := sftp.Join(Auto.ProjectPathOnLine, fmt.Sprintf("shang%d.war", rand.Int()))
	Auto.WarPathLocal = mvnPackage()
	cloneFile(connection, remoteFilePath, Auto.WarPathLocal)
	rmUnzip(connection, remoteFilePath)
	cpProfile(connection)
	startTomcat(connection)
}

// 字节的单位转换 保留两位小数
func formatFileSize(s int64) (size string) {
	if s < 1024 {
		return fmt.Sprintf("%.2fB", float64(s)/float64(1))
	} else if s < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(s)/float64(1024))
	} else if s < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(s)/float64(1024*1024))
	} else if s < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(s)/float64(1024*1024*1024))
	} else if s < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(s)/float64(1024*1024*1024*1024))
	} else { //if s < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(s)/float64(1024*1024*1024*1024*1024))
	}
}

func rmUnzip(connection *ssh.Client, remoteFileName string) {
	ssh_server.OpenSessionRunCmd(connection, "rm -rf "+Auto.ProjectPathOnLine+"/WEB-INF")
	ssh_server.OpenSessionRunCmd(connection, "rm -rf  "+Auto.ProjectPathOnLine+"/org")
	ssh_server.OpenSessionRunCmd(connection, "rm -rf  "+Auto.ProjectPathOnLine+"/META-INF")
	log.Printf("原删除成功代码")
	ssh_server.OpenSessionRunCmd(connection, "unzip "+remoteFileName+" -d "+Auto.ProjectPathOnLine+"/")
	log.Printf("解压代码")
	ssh_server.OpenSessionRunCmd(connection, "rm -rf "+remoteFileName)
	log.Printf("部署成功")
}
func cloneFile(connection *ssh.Client, remoteFilePath string, localFileName string) {
	sftpClient := ssh_server.OpenFtpClient(connection)
	defer sftpClient.Close()
	//获取当前目录
	cwd, err := sftpClient.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("线上当前目录：", cwd)
	//上传文件(将本地file.dat文件通过sftp传到远程服务器)
	remoteFileName := fmt.Sprintf("shang%d.war", rand.Int())
	fmt.Sprintf("upload %s", remoteFileName)
	remoteFile, err := sftpClient.Create(remoteFilePath)
	log.Println("文件名：", remoteFile.Name())
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer remoteFile.Close()
	//打开本地文件shang
	localFile, err := os.Open(localFileName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer localFile.Close()

	//本地文件流拷贝到上传文件流
	log.Println("文件拷贝" + localFileName + "到服务器")
	log.Println("文件拷贝上传中。。。。。。")
	n, err := io.Copy(remoteFile, localFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	//获取本地文件大小
	localFileInfo, err := os.Stat(localFileName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("文件上传成功[%s->%s]本地文件大小：%s，上传文件大小：%s", localFileName, remoteFile.Name(), formatFileSize(localFileInfo.Size()), formatFileSize(n))
}
