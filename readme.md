## 介绍

> syncd-cli 是自动化部署工具 [syncd](https://github.com/dreamans/syncd) 的一个命令行客户端，用于批量添加server，实现一键自动添加，提高开发效率。
>
> Syncd是一款开源的代码部署工具，它具有简单、高效、易用等特点，可以提高团队的工作效率。

## 安装

required要求

- [X] go1.8+
- [X] syncd2.0+

`go get` 方式

```shell
$ go get gogs.wangke.co/go/syncd-cli
$ syncd-cli -h
```

`git clone` 方式

```shell
$ git clone https://gogs.wangke.co/go/syncd-cli.git
$ cd syncd-cli && go build -o syncd-cli syncd-cli.go
$ ./syncd-cli -h
```

## Usage

```shell
# ./syncd-cli -h                                                                                                  [12:18:53]
syncd-cli version:1.0.0
Usage syncd-cli <command> [-aupginsh]

command [--add] [--list]  [user|server]

add server example:
        1) syncd-cli -d server -g 2 -i 192.168.1.1,test.example.com -n test01,test02 -s 9527,22
        2) syncd-cli --add server  --roleGroupId 2 --ipEmail 192.168.1.1 --names test01 --sshPort 9527
add user example:
        1) syncd-cli --add user  --ipEmail text@wangke.co --names test01
        2) syncd-cli  -d user  -i text@wangke.co -n test01
list server and user example:
        1) syncd-cli -l user
        2) syncd-cli -l server 
        3) syncd-cli --list 
        4) syncd-cli --list server 

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

### @v1.1.0

从命令行读取批量文件的信息,感觉太冗余了, 不方便创建,还是以文件的方式创建批量容易一点,
所用想了一个从文件里面读取信息,然后根据信息创建相应资源的方法.

```
@v1.1.0  
$ go run syncd-cli.go -h
syncd-cli version:1.1.0
Usage syncd-cli <command> [-afhpu]

command <apply> <get>  [user|server] <?-f files>

add server example:
        1) syncd-cli apply user -f files
        2) syncd-cli apply server -f files
list server and user example:
        1) syncd-cli get user
        2) syncd-cli get server

Options:
  -f, --file string       add server/user from files
  -h, --help              this help
  -a, --hostApi string    sycnd server addr api (default "http://127.0.0.1:8878")
  -p, --password string   password for syncd tools (default "111111")
  -u, --user string       user for syncd tools (default "syncd")
  

```


**file文件的格式**

`testserver`,`testuser`等文件名称无要求, 但是对文件的格式要求.以一个空格进行分割.

```shell
# 第一列是groupid,对应的集群id;
# 第二列是name, 对应的是server的名字
# 第三列是ip/hostname, 对应server的ip或者域名
# 第四列是sshport, 对应的是server的ssh端口
$ cat testserver
1 test01 test.wangke.co 22
1 test02 test01.wangke.co 9527
1 test03 test02.wangke.co 6822
```
`testuser`文件内容以空格区分,共四列.(后续可以添加至6列,源码的user还有电话号码,真实姓名等,非必须 批量创建的默认密码为`111111`)

```shell
# 第一列是role_id, 对应的是角色, 比如1是管理员
# 第二列是name, 对应的是用户名
# 第三列是email, 对应的是用户的邮箱
# 第四列是status, 对应的是用户能否登陆.
$ cat testuser
1 test01 test01@wangke.co 1
1 test02 test02@wangke.co 1

```

因为`testuser`和`testserver`的的文件格式和数据类型是一样的, 所用到的方法是一样的, 唯一的区分就是利用`apply user`还是`apply server`

```go
type server struct {
    id int
    name string
    ip string
    port int
}

type user struct {
	id int
	name string
	email string
	status int
}

```
### 重要提醒
```

方法是一样的. 所以标志位很重要, 不然创建错了就是连环错误了.

$ syncd-cli apply user -f testuser
$ syncd-cli apply server -f testserver
```



## example
### @v1.0.0

```shell
root@master-louis: ~/go/src/github.com/oldthreefeng/syncd-cli master ⚡
# ./syncd-cli -i 192.168.1.2,text.example.com -n test1,texte -s 9527,22              [12:18:58]
INFO[0000] your token is under .syncd-token             
INFO[0000] group_id=1&name=test1&ip=192.168.1.2&ssh_port=9527  
INFO[0000] {"code":0,"message":"success"}               
INFO[0000] group_id=1&name=texte&ip=text.example.com&ssh_port=22  
INFO[0000] {"code":0,"message":"success"}

# 将test01邮箱为text@wangke.co加入管理员,默认密码为111111
$./syncd-cli -d user -i text@wangke.co -n test01   
time="2019-10-20T17:59:08+08:00" level=info msg="your token is under .syncd-token\n"
time="2019-10-20T17:59:08+08:00" level=info msg="role_id=1&username=test01&password=1111111&email=text@wangke.co&status=1"
time="2019-10-20T17:59:08+08:00" level=info msg="{\"code\":0,\"message\":\"success\"}"
```

### @v1.1.0

```
@v1.1.0
# 采用从文件读取方式创建server
$ go run syncd-cli.go apply server -f test.log
time="2019-10-21T00:25:03+08:00" level=info msg="your token is under .syncd-token\n"
time="2019-10-21T00:25:03+08:00" level=info msg="group_id=1&name=test01&ip=test.wangke.co&ssh_port=22\n"
time="2019-10-21T00:25:03+08:00" level=info msg="{\"code\":0,\"message\":\"success\"}"
time="2019-10-21T00:25:03+08:00" level=info msg="group_id=1&name=test02&ip=test01.wangke.co&ssh_port=9527\n"
time="2019-10-21T00:25:03+08:00" level=info msg="{\"code\":0,\"message\":\"success\"}"
time="2019-10-21T00:25:04+08:00" level=info msg="group_id=1&name=test03&ip=test02.wangke.co&ssh_port=6822\n"
time="2019-10-21T00:25:04+08:00" level=info msg="{\"code\":0,\"message\":\"success\"}"

$ go run syncd-cli.go apply user -f testuser.log
time="2019-10-21T00:27:03+08:00" level=info msg="your token is under .syncd-token\n"
time="2019-10-21T00:27:03+08:00" level=info msg="role_id=1&username=test01&password=111111&email=test01@wangke.co&status=1"
time="2019-10-21T00:27:03+08:00" level=info msg="{\"code\":0,\"message\":\"success\"}"
time="2019-10-21T00:27:03+08:00" level=info msg="role_id=1&username=test02&password=111111&email=test02@wangke.co&status=1"
time="2019-10-21T00:27:03+08:00" level=info msg="{\"code\":0,\"message\":\"success\"}"

$ go run syncd-cli.go get user
[map[ctime:0 email:louis@wangke.co id:2 last_login_ip: last_login_time:0 mobile: password: role_id:1 role_name:管理员 salt: status:1 truename: username:louis] map[ctime:0 email: id:1 last_login_ip: last_login_time:0 mobile: password: role_id:1 role_name:管理员 salt: status:1 truename: username:syncd]]
2

$ go run syncd-cli.go get server
[map[ctime:1.571545269e+09 group_id:1 group_name:aliyun id:5 ip:text.example.com name:texte ssh_port:22] map[ctime:1.571545269e+09 group_id:1 group_name:aliyun id:4 ip:192.168.1.2 name:test1 ssh_port:9527] map[ctime:1.571541399e+09 group_id:1 group_name:aliyun id:3 ip:192.168.1.1 name:test01 ssh_port:22] map[ctime:1.57140845e+09 group_id:3 group_name:vrtul id:2 ip:vps.wangke.co name:vrtest ssh_port:9527] map[ctime:1.571408398e+09 group_id:1 group_name:aliyun id:1 ip:gogs.wangke.co name:alitest ssh_port:9527]]
5
```
添加如下:
![](https://pic.fenghong.tech/syncd-cli.png)
![](https://pic.fenghong.tech/syncd-cli-add-user.png)

## 算法思路 ##

本来想开发和`kubectl`,`go`,`kubeadm`等类似的管理cli. 奈何时间水平有限. 

脑子里想的是这样的@v1.1.0实现的比较简陋

```cgo
$ syncd get user 
$ syncd get server 
$ syncd apply -f adduser.yaml
```
实际上...

```cgo
$ syncd-cli --list user
$ syncd-cli --list server

$ syncd-cli --add user -i test@wangke.co -n test01 
```

> 整体上, 利用`http`的`GET`还有`POST`完成显示和添加动作的. [gorequest](https://github.com/parnurzeal/gorequest)的`GET/POST`的确好用,可以试试. 
>
> 记录日志当然是用的[logrus](https://github.com/sirupsen/logrus), 当时用的`go mod`学习教程就是用的这个模板, 日志的格式也可以, 方便阅读.
>
> 命令行的开发主要就是用的[pflag](https://github.com/spf13/pflag), 看了`kubernetes`和`docker`源码相关, `kubectl`等命令行管理工具也是基于这个开发的.

首先, 登录验证, 获取token, 将token存入当前目录下的`.syncd-token`, 其次, 获取user/server列表或者添加user/server, 逻辑都是一样的,发送`POST`请求, 同时携带cookie, 将cookie的`name`和`value`封装成`http.cookie`, 每次需要用到,直接调用即可. 

## TODO
- [x] add server from cli
- [x] add user from cli
- [x] list server
- [x] list user
- [ ] list project
- [x] add server info from file
- [x] add user info from file

