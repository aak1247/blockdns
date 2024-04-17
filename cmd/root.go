package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Execute() {
	err := execute()
	if err != nil {
		fmt.Println("failed: " + err.Error())
		os.Exit(1)
	}
}

func execute() error {
	var args = os.Args[1:]
	if len(args) == 0 {
		return errors.New("Please input the url")
	}
	var path = args[0]
	remove := false
	if path == "-r" {
		remove = true
		if len(args) == 1 {
			return errors.New("Please input the url")
		}
		path = args[1]
	}
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Println(`Easily block a url by dns spoofing, add corresponding record to etcd 
Usage: blockdns -r [url]

eg: blockdns baidu.com`)
	}
	cmd := fmt.Sprint("etcdctl put " + convert(path) + " '{\"host\":\"127.0.0.1\",\"ttl\":60000}'")
	if remove {
		cmd = fmt.Sprint("etcdctl del " + convert(path))
	}
	// 执行
	fmt.Println("execute command " + cmd)
	err := exec.Command("bash", "-c", cmd).Run()
	if err != nil {
		return err
	}
	fmt.Println("add dns block to etcd success")
	return nil
}

// reverse url and convert "." to "/"
func convert(url string) string {
	words := strings.Split(url, ".")
	reversed := make([]string, len(words))
	// reverse it
	for i, word := range words {
		reversed[len(words)-i-1] = word
	}
	return "/coredns/" + strings.Join(reversed, "/")
}
