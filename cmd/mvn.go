package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var mvn = &cobra.Command{
	Use:          "package",
	Long:         `打包当前目录`,
	Short:        "打包当前目录",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("mvn打包")
		mvnPackage()
	}}

/**
  打包返回打包后的文件
*/
func mvnPackage() string {
	s, err := os.Stat(Auto.WarPathLocal)
	if err != nil {
		log.Fatal(err)
	}
	if !s.IsDir() {
		log.Println("非文件夹，不自动打包")

		if strings.HasSuffix(s.Name(), ".war") {
			log.Fatalf("不是文件夹，也不是war包")
		}
		return Auto.WarPathLocal
	}
	mvnTest := "-Dmaven.test.skip=true"
	cmdExec := exec.Command("mvn", "clean", "package", mvnTest)
	cmdExec.Dir = Auto.WarPathLocal
	log.Println(cmdExec.String())
	log.Println("打包中")
	out, err := cmdExec.CombinedOutput()
	if err != nil {
		log.Printf("combined out:\n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err.Error())
	}
	log.Printf("打包完成")
	files, err := ioutil.ReadDir(Auto.WarPathLocal + "/target")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".war") {
			filename := Auto.WarPathLocal + "/target/" + f.Name()
			log.Printf(filename)
			return filename
		}
	}
	log.Fatalf("打包失败")
	return ""
}
