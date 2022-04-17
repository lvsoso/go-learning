代码来自 https://github.com/Q1mi/golang-unit-test-demo


```shell
go test -v
go test -run=XXX
go test -cover


go test -cover -coverprofile=c.out
go tool cover -html=c.out

go test -bench=Split
go test -bench=Split -benchmem

go test -bench=BenchmarkFib
go test -bench=Fib40 -benchtime=20s

go test -bench=. -cpu 1
```
