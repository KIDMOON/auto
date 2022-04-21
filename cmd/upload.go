package cmd

import (
	ssh_server "awesomeProject/src"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
)

var upload = &cobra.Command{
	Use:          "upload",
	Long:         `upload上传warPathLocal东西`,
	Short:        "upload上传warPathLocal,非war执行打包命令",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("upload start")
		connection := ssh_server.SSh_Connection(&Auto)
		defer connection.Close()
		remoteFilePath := sftp.Join(Auto.ProjectPathOnLine, fmt.Sprintf("shang%d.war", rand.Int()))
		Auto.WarPathLocal = mvnPackage()
		cloneFile(connection, remoteFilePath, Auto.WarPathLocal)
	}}
