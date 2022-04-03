

```shell

sudo apt-get install make g++ libz-dev libsnappy-dev

git submodule update --init"

# 编译rocksdb 静态库
cd rocksdb && make static_lib

# rocksdb 测试

cd rocksdb_performance

base) lv@lv:rocksdb_performance$ ./test_basic -t 100000
total record number is 100000
value size is 1000
100000 records put in 317554 usec, 3.17554 usec average
100000 records get in 145347 usec, 1.45347 usec average

(base) lv@lv:rocksdb_performance$ ./ingest_data -t 100000 -s 1000
total record number is 100000
value size is 1000
100000 total records set in 296197 usec,2.96197 usec average, throughput 337.613 MB/s, rps is 337613

base) lv@lv:rocksdb_performance$ ./ingest_data -t 1000000 -s 1000
total record number is 1000000
value size is 1000
1000000 total records set in 2912944 usec,2.91294 usec average, throughput 343.295 MB/s, rps is 343295

(base) lv@lv:rocksdb_performance$ ./ingest_data -t 10000000 -s 1000
total record number is 10000000
value size is 1000
10000000 total records set in 30143994 usec,3.0144 usec average, throughput 331.741 MB/s, rps is 331741
free(): invalid pointer
已放弃 ???

(base) lv@lv:rocksdb_performance$ dd bs=1000 count=100000 if=/dev/zero of=tmpfile
记录了100000+0 的读入
记录了100000+0 的写出
100000000 bytes (100 MB, 95 MiB) copied, 0.199428 s, 501 MB/s

(base) lv@lv:rocksdb_performance$ dd bs=1000 count=1000000 if=/dev/zero of=tmpfile
记录了1000000+0 的读入
记录了1000000+0 的写出
1000000000 bytes (1.0 GB, 954 MiB) copied, 1.64298 s, 609 MB/s

(base) lv@lv:rocksdb_performance$ dd bs=1000 count=10000000 if=/dev/zero of=tmpfile
记录了10000000+0 的读入
记录了10000000+0 的写出
10000000000 bytes (10 GB, 9.3 GiB) copied, 15.9707 s, 626 MB/s
```
