package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var bufChan chan bool

func test() {
	bufChan = make(chan bool, 5)
	go func() {
		for _ = range bufChan {
			// what am I going to do, log this?
			fmt.Printf("：：dasdsads")
		}

	}()

}

func test1() {
	select {
	case bufChan <- true:
		time.Sleep(time.Second)
	default:
		fmt.Println("资源已满，请稍后再试")
		time.Sleep(time.Second)
	}
}

func main() {
	k := url.Values{}
	k.Add("pwdcode", "111")
	k.Add("account", "gvemf")
	k.Add("phone", "16521580951")
	rep, err := http.PostForm("http://123.57.205.127/tq/api/json_consume.php", k)
	if err == nil {
		body, err := ioutil.ReadAll(rep.Body)
		if err == nil {
			log.Print(string(body))
		}
	}

}
