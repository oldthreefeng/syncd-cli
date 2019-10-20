> syncd-cli 是自动化部署工具 [syncd](https://github.com/dreamans/syncd) 的一个命令行客户端，用于批量添加`server`，实现一键自动添加，提高开发效率。 

## 安装

`go get` 方式

```shell
$ go get gogs.wangke.co/go/syncd-cli
```
`git clone` 方式

```cgo
$ git clone gogs.wangke.co/go/syncd-cli
$ cd syncd-cli && go build -o syncd-cli syncd-cli.go
```

## Usage

```shell
# ./syncd-cli -h                                                                                                  [12:18:53]
syncd-cli version:1.0.0
Usage syncd-cli <command> [-aupginsh]

command [--add] [--list]  [user|server]

add server example:
        1) syncd-cli -d server -u syncd -p 111111  -g 2 -i 192.168.1.1,test.example.com -n test01,test02 -s 9527,22
        2) syncd-cli --add server  --user syncd --password 111111 --roleGroupId 2 --ipEmail 192.168.1.1 --names test01 --sshPort 9527
add user example:
        1) syncd-cli --add user --user syncd --password 111111  --roleGroupId 1 --ipEmail text@wangke.co --names test01
        2) syncd-cli  -d user -u syncd -p 111111 -g 1 -i text@wangke.co -n test01
list server and user example:
        1) syncd-cli -l user -u syncd -p 111111
        2) syncd-cli -l server -u syncd -p 111111
        3) syncd-cli --list user --user syncd --password 111111
        4) syncd-cli --list server --user syncd --password 111111

Options:
  -d, --add string        add user or server
  -h, --help              this help
  -a, --hostApi string    sycnd server addr api (default "http//127.0.0.1:8878/")
  -i, --ipEmail strings   set ip/hostname to the cluster with names // or email for add user, use ',' to split
  -l, --list string       list server or user
  -n, --names strings     set names to the cluster with ips, use ',' to split
  -p, --password string   password for syncd tools (default "111111")
  -g, --roleGroupId int   group_id for cluster // or role_id for user, must be needed (default 1)
  -s, --sshPort ints      set sshPort to the cluster server, use ',' to split
  -u, --user string       user for syncd tools (default "syncd")

```

## example
``` 
root@master-louis: ~/go/src/github.com/oldthreefeng/syncd-cli master ⚡
# ./syncd-cli -u syncd -p 111111 -g 1 -i 192.168.1.2,text.example.com -n test1,texte -s 9527,22              [12:18:58]
INFO[0000] your token is under .syncd-token             
INFO[0000] group_id=1&name=test1&ip=192.168.1.2&ssh_port=9527  
INFO[0000] {"code":0,"message":"success"}               
INFO[0000] group_id=1&name=texte&ip=text.example.com&ssh_port=22  
INFO[0000] {"code":0,"message":"success"}

$./syncd-cli -d user -g 1 -i text@wangke.co -n test01
time="2019-10-20T17:59:08+08:00" level=info msg="your token is under .syncd-token\n"

time="2019-10-20T17:59:08+08:00" level=info msg="role_id=1&username=test01&password=1111111&email=text@wangke.co&status=1"
time="2019-10-20T17:59:08+08:00" level=info msg="{\"code\":0,\"message\":\"success\"}"

```
添加如下:
![](https://pic.fenghong.tech/syncd-cli.png)
![](https://pic.fenghong.tech/syncd-cli-add-user.png)

## TODO
- [x] add server
- [x] add user
- [x] list server
- [x] list user
- [ ] list project
