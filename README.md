# localmap

有时候，需要紧急临时暴露一个本地端口到公网中进行调试，但是又不方便使用 vpn，ddns 等东西。
而此时你手上正好有一台 vps，你就想能不能通过这台服务器做一次转发呢？
`localmap` 就是这么诞生的。

You need to debug some programs in public network but they are running in local network.
It's troublesome to use vpn or ddns.
You want to do port transfer through a vps that you have.
Yeah. `localmap` do that.

# Feature

* 哪里都可以使用，不管是校园网还是公共wifi，只要能连接到转发服务器的网络都可以使用。
* 基于 tcp 的端口转发，支持 http, https, ws 等
* 足够轻量，可以使用自己的服务器做转发

* use everywhere. campus network, public area wifi.
* base on tcp transfer. support http, https, ws, etc...
* very lite, use yourself server to do transfer.

# Install

From [![CircleCI](https://circleci.com/gh/XGHeaven/localmap/tree/master.svg?style=svg)](https://circleci.com/gh/XGHeaven/localmap/tree/master) [CircleCI Artifacts](https://circleci.com/gh/XGHeaven/localmap/tree/master) or from source

```go
go get github.com/xgheaven/localmap
```

# Usage

安装之后，会有一个 `localmap` 的二进制文件，有两种工作模式，一直是作为服务端使用，另外一种是作为客户端使用。
`localmap` will install your computer. They have two work mode, one as server, the other as client.

## Start as Server

当作为服务端使用的时候，需要设置服务器监听的端口，服务器在这个端口上接受客户端的连接。
默认情况下会自动监听 8000 端口来提供服务。
Need to set which port to listen when start as server. listen to port 8000 default.

```bash
# localmap -sport 8000
```

## Start as Client

当作为客户端使用的时候，需要制定服务器的地址，服务器的端口，和客户端需要转发的端口。
Need to set which server port, server address and client port to transfer.

```
# localmap -addr yourserver.com -sport 8000 -cport 80
```

执行之后，服务器会自动生成一个新的端口，你只需要通过新分配的端口即可访问本地端口。
you just use address that alloc from server to browse your website.

如果你不想自己搭建服务器，可以使用我免费提供的服务器 `localmap.xgheaven.com`，不过带宽有限，建议还是通过自己的服务器。
if you want not to build a server, you can use `localmap.xgheaven.com` which I provide.

# Thanks

欢迎大家来提供反馈。
Welcome to support issue.
