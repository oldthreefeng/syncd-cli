/*
Copyright 2019 louis.
@Time : 2019/10/20 10:00
@Author : louis
@File : api-cli
@Software: GoLand

*/

package main

import (
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
	user     string
	password string
	GroupId  int
	Names    []string
	Ips      []string
	SSHPort  []int
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

func authCookie() *http.Cookie {
	cookie := http.Cookie{}
	cookie.Name = "_syd_identity"
	cookie.Value = GetToken()
	return &cookie
}

func usages() {
	_, _ = fmt.Fprintf(os.Stderr, `syncd-cli version:1.0.0
Usage syncd-cli [-aupginsh] 

example: 1) syncd-cli -u syncd -p 111111  -g 2 -i 192.168.1.1,test.example.com -n test01,test02 -s 9527,22
         2) syncd-cli --user syncd --password 111111 --groupId 2 --ips 192.168.1.1 --names test01 --sshPort 9527

Options:
`)
	flag.PrintDefaults()
}


func init() {
	flag.StringVarP(&host, "hostApi", "a", "http://127.0.0.1:8878/", "sycnd server addr api")
	flag.StringVarP(&user, "user", "u", "syncd", "user for syncd tools")
	flag.StringVarP(&password, "password", "p", "111111", "password for syncd tools")
	flag.IntVarP(&GroupId, "groupId", "g", 1, "groupId for cluster, must be needed")
	flag.StringSliceVarP(&Ips, "ips", "i", []string{""}, "set ip/hostname to the cluster, use ',' to split")
	flag.StringSliceVarP(&Names, "names", "n", []string{""}, "set names to the cluster, use ',' to split")
	flag.IntSliceVarP(&SSHPort, "sshPort", "s", []int{}, "set sshPort to the cluster, use ',' to split")
	flag.BoolVarP(&h, "help", "h", false, "this help")
	flag.Usage = usages
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	login(user, password)
	for k,v := range Ips{
		serverAdd(GroupId,Names[k],v,SSHPort[k])
	}
}
