/*
Copyright 2019 louis.
@Time : 2019/10/21 0:08
@Author : louis
@File : syncd_test
@Software: GoLand

*/

package main

import (
	"bufio"
	"fmt"
	"os"
)

type server struct {
	gid  int
	name string
	ip   string
	port int
}

func readFromFile(file string)  {
	openFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer openFile.Close()
	var newserver []server
	scanner := bufio.NewScanner(openFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			break
		}
		var gid, port int
		var name, ip string
		_, err := fmt.Sscanf(line, "%d %s %s %d", &gid, &name, &ip, &port)
		if err != nil {
			return
		}
		newserver = append(newserver, server{
			gid:  gid,
			name: name,
			ip:   ip,
			port: port,
		})
	}
	fmt.Println(newserver)
}

func main() {
	readFromFile("testserver.log")
	readFromFile("testuser.log")
}
