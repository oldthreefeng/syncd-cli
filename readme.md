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

```cgo
# ./syncd-cli -h                                                                                                  [12:18:53]
syncd-cli version:1.0.0
Usage syncd-cli [-aupginsh] 

example: 1) syncd-cli -u syncd -p 111111  -g 2 -i 192.168.1.1,test.example.com -n test01,test02 -s 9527,22
         2) syncd-cli --user syncd --password 111111 --groupId 2 --ips 192.168.1.1 --names test01 --sshPort 9527

Options:
  -g, --groupId int       groupId for cluster, must be needed (default 1)
  -h, --help              this help
  -a, --hostApi string    sycnd server addr api (default "http://127.0.0.1:8878/")
  -i, --ips strings       set ip/hostname to the cluster, use ',' to split
  -n, --names strings     set names to the cluster, use ',' to split
  -p, --password string   password for syncd tools (default "111111")
  -s, --sshPort ints      set sshPort to the cluster, use ',' to split
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
```
添加如下:
![](https://pic.fenghong.tech/syncd-cli.png)

## TODO
- [x] add server
- [ ] add user
- [ ] list server
- [ ] list user
- [ ] list project
