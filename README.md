Ad Hoc Analytics
==

高性能即席分析引擎：每秒处理TB量级消费者画像数据

使用方法：

```
cd tools
go run lookup_server.go --data_files data.csv --tag_option_file tag_option.csv
```

其中 data.csv 的格式为（参见testdata/data.txt）

```
日期,用户id,tag1:option1,tag2:option2,...
```

其中日期为数据采集的日期，用户 id 为 uint64 整数，tag和 option 也都是整数

tag_option.csv 格式为（参见 data/tag_option.csv ）

```
tag id,tag name,option id,option name,category name
```

category, tag, option 的关系如下

```
category 1:
  tag 1: option 1
  tag 1: option 2
  tag 2: option 3
  tag 3: option 4
category 2:
  tag 4: option 5
  tag 4: option 6
```

其中跨 category 的 tag id 不会发生重复，跨 tag 的 option id 不会发生重复。
