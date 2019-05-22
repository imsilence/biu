

# 介绍 #

参考[biu-framework](https://github.com/0xbug/Biu-framework) 思路使用实现Go语言版本

# 改进 #

1. plugin 格式调整并支持不同方式提交数据

2. 高并发

3. 单二进制文件+plugin描述文件部署

# 使用 #

1. 编译 `go build`

2. 拷贝 biu.exe和plugins目录到运行环境

3. 命令行运行 `biu.exe -h`

4. web api运行 `biu.exe api -h`

    a. 浏览器访问: `http://ip:port/task/`

    b. curl访问: `curl -XPOST "127.0.0.1:8080/task/" -d "H=127.0.0.1:8888&c=192.168.1.1/24&p=all"`