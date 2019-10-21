/*
Copyright 2019 louis.
@Time : 2019/10/20 10:00
@Author : louis
@File : syncd-cli
@Software: GoLand

*/

package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	tokenFile = ".syncd-token"
	agent     = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:68.0) Gecko/20100101 Firefox/68.0"
)

var (
	//host     = "https://syncd.fenghong.tech/"
	host     string
	list     string
	user     string
	password string
	GroupId  int
	Names    []string
	Ips      []string
	SSHPort  []int
	add      string
	files    string
	h        bool
)

var _token string

func TokenFail() {
	RemoveToken()
	panic(fmt.Sprintf("login faild, please set the right password"))
}

func RemoveToken() {
	if err := os.Remove(tokenFile); err != nil {
		log.Infoln("remove .token failed")
	}
}

func SetToken(token string) {
	err := ioutil.WriteFile(tokenFile, []byte(token), 0644)
	if err != nil {
		log.Fatalln(err)
	}
	_token = token
}

func GetToken() string {
	if _token == "" {
		tokenByte, err := ioutil.ReadFile(tokenFile)
		if err != nil {
			log.Fatalln("need login")
		}

		_token = string(tokenByte)
	}
	return _token
}

func md5s(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

type RespData map[string]interface{}

type Response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    RespData `json:"data"`
}

func listServerDetail(res RespData) {
	for _, v := range res {
		fmt.Println(v)
	}
}

func ParseResponse(respBody string) (RespData, error) {
	response := Response{}
	err := json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		panic(err)
	}

	if response.Code == 1005 {
		TokenFail()
	}

	if response.Code != 0 {
		return nil, errors.New(response.Message)
	}

	return response.Data, nil
}

func login(user, password string) {
	url := host + "api/login"
	_, _, errs := gorequest.New().
		Post(url).
		Type("form").
		AppendHeader("Accept", "application/json").
		Send(fmt.Sprintf("username=%s&password=%s", user, md5s(password))).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != 200 {
				panic(fmt.Sprintf("%s", errs))
			}

			respData, err := ParseResponse(body)
			if err != nil {
				panic(err)
			}

			//respData
			SetToken(respData["token"].(string))
		})

	if errs != nil {
		log.Fatalf("%s", errs)
	}
	log.Infof("your token is under %s\n", tokenFile)
}

func userAdd(roleId int, userName, email string, status int) {
	url := host + "api/user/add"
	pass := "111111"
	_, body, errs := gorequest.New().Post(url).
		AppendHeader("Accept", "application/json").
		AppendHeader("User-Agent", agent).
		AddCookie(authCookie()).
		Send(fmt.Sprintf("role_id=%d&username=%s&password=%s&email=%s&status=%d",
			roleId, userName, md5s(pass), email, status)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != 200 {
				panic(errs)
			}
		})
	if errs != nil {
		log.Fatalln(errs)
	}
	log.Infof("role_id=%d&username=%s&password=%s&email=%s&status=%d",
		roleId, userName, pass, email, status)
	log.Infoln(body)
}

func serverAdd(groupId int, name, ip string, sshPort int) {
	url := host + "api/server/add"
	_, body, errs := gorequest.New().Post(url).
		AppendHeader("Accept", "application/json").
		AppendHeader("User-Agent", agent).
		AddCookie(authCookie()).
		Send(fmt.Sprintf("group_id=%d&name=%s&ip=%s&ssh_port=%d",
			groupId, name, ip, sshPort)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != 200 {
				panic(errs)
			}
		})
	if errs != nil {
		log.Fatalln(errs)
	}
	log.Infof("group_id=%d&name=%s&ip=%s&ssh_port=%d\n",
		groupId, name, ip, sshPort)
	log.Infoln(body)
}

type QueryBind struct {
	Keyword string `form:"keyword"`
	Offset  int    `form:"offset"`
	Limit   int    `form:"limit" binding:"required,gte=1,lte=999"`
}

func List(api string) {
	url := host + api
	_, body, errs := gorequest.New().Get(url).Query(QueryBind{Keyword: "", Offset: 0, Limit: 7}).
		AppendHeader("Accept", "application/json").
		AppendHeader("User-Agent", agent).
		AddCookie(authCookie()).
		End()
	if errs != nil {
		log.Fatalln(errs)
	}
	var serverBody Response
	err := json.Unmarshal([]byte(body), &serverBody)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Infoln(serverBody)
	listServerDetail(serverBody.Data)
}

func authCookie() *http.Cookie {
	cookie := http.Cookie{}
	cookie.Name = "_syd_identity"
	cookie.Value = GetToken()
	return &cookie
}

func usages() {
	_, _ = fmt.Fprintf(os.Stderr, `syncd-cli version:1.1.0
Usage syncd-cli <command> [-afhpu] 

command <apply|get>  <user|server> [?-f files]

add server example: 
	1) syncd-cli apply user -f files
	2) syncd-cli apply server -f files
list server and user example:
	1) syncd-cli get user
	2) syncd-cli get server

Options:
`)
	flag.PrintDefaults()
}

func init() {
	flag.StringVarP(&host, "hostApi", "a", "http://127.0.0.1:8878/", "sycnd server addr api")
	//flag.StringVarP(&host, "hostApi", "a", "https://syncd.fenghong.tech/", "sycnd server addr api")
	flag.StringVarP(&user, "user", "u", "syncd", "user for syncd tools")
	flag.StringVarP(&password, "password", "p", "111111", "password for syncd tools")

	flag.StringVarP(&add, "add", "d", "", "add user or server(deprecated)")
	flag.StringVarP(&files, "file", "f", "", "add server/user from files")
	flag.StringVarP(&list, "list", "l", "", "list server and user(deprecated)")
	flag.IntVarP(&GroupId, "roleGroupId", "g", 1, "group_id for cluster // or role_id for user, must be needed(deprecated)")
	flag.StringSliceVarP(&Ips, "ipEmail", "i", []string{""}, "set ip/hostname to the cluster with names // or email for add user, use ',' to split(deprecated)")
	flag.StringSliceVarP(&Names, "names", "n", []string{""}, "set names to the cluster with ips, use ',' to split(deprecated)")
	flag.IntSliceVarP(&SSHPort, "sshPort", "s", []int{}, "set sshPort to the cluster, use ',' to split(deprecated)")
	flag.BoolVarP(&h, "help", "h", false, "this help")
	flag.Usage = usages
}

func useV100() {
	//是否列出server,user

	switch list {
	case "user":
		List("api/user/list")
	case "server":
		List("api/server/list")
	default:
		fmt.Println("use `syncd-cli get [user | server]` instead")
	}

	if Ips[0] == "" || Names[0] == "" || SSHPort[0] == 0 {
		return
	}
	switch add {
	case "user":
		// userAdd() //easy to add
		for k, v := range Ips {
			userAdd(GroupId, Names[k], v, 1)
		}
	case "server":
		// Ips未指定,则返回
		for k, v := range Ips {
			serverAdd(GroupId, Names[k], v, SSHPort[k])
		}
	}
}

type server struct {
	gid  int
	name string
	ip   string
	port int
}

func readFromServerFile(file string) []server {
	openFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("%v",err)
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
			return nil
		}
		newserver = append(newserver, server{
			gid:  gid,
			name: name,
			ip:   ip,
			port: port,
		})
	}
	return newserver
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}

	// 登录认证
	login(user, password)

	// 使用v1.0.0
	if list != "" {
		useV100()
	}

	switch os.Args[1] {
	case "apply":
		switch os.Args[2] {
		case "server":
			ser := readFromServerFile(files)
			for _, v := range ser {
				serverAdd(v.gid, v.name, v.ip, v.port)
			}
		case "user":
			usr := readFromServerFile(files)
			for _, v := range usr {
				// 偷个懒, 数据类型一样,结构一样, 所以从文件读取的是一样的
				// v.gid ==> roleId
				// v.name==> username
				// v.ip  ==> email
				// v.port==> status
				userAdd(v.gid, v.name, v.ip, v.port)
			}
		default:
			fmt.Println("syncd-cli apply [user|server] -f files")
		}

	case "get":
		switch os.Args[2] {
		case "user":
			List("api/user/list")
		case "server":
			List("api/server/list")
		default:
			fmt.Println("syncd-cli get [user | server]")
		}
	default:
		fmt.Println()
		fmt.Println("	Use syncd-cli@v1.1.0 instead")
		fmt.Println("syncd-cli <get|apply> <user|server> [?-f filename>]")
		fmt.Println("syncd-cli <get|apply> <user|server> [?--file filename]")
	}

}
