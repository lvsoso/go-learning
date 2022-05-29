## socks5

[RFC 1928](https://datatracker.ietf.org/doc/html/rfc1928)


- 建立链接
- 协商（方法’01‘需要子协商）
- 请求
- 转发


### start server

```shell

cd server

go run cmd/main.go

```

### start request

```shell
curl -v --proxy socks5://admin:123456@localhost:1080 baidu.com
```
