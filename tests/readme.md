> 参考：
> [近期一个Golang项目的测试实践全记录](https://mp.weixin.qq.com/s/cXLH5Wr4qTocrEYmwUPe9Q)

## suites
存放测试套件

### suites/xxx
这里存放测试套件，测试套件文件夹需要包含下列文件：

before.go存放有
- SetUp() 函数，这个函数在Suite运行之前会运行
- Before() 函数，这个函数在所有 Case 运行之前运行

after.go存放有
- TearDown() 函数，这个函数在Suite运行之后会运行
- After() 函数，这个函数在 Case 运行之后运行

### TestSuite初始化
把需要初始化的DB结构使用sql文件导出，放在目录中。这样，每个人想要跑这一套测试用例，只需要搭建一个mysql数据库，倒入sql文件，就可以搭建好数据库环境了。其他的初始化数据等都在TestSuite初始化的SetUp函数中调用。

关于保存测试数据环境，我这里有个小贴士，在SetUp函数中实现 清空数据库+初始化数据库 ，在TearDown函数中不做任何事情。这样如果你要单独运行某个TestSuite，能保持最后的测试数据环境，有助于我们进行测试数据环境测试。
### TestCase编写
在集成测试环境中，TestCase编写调用HTTP请求就是使用正常的 httptest包，其使用方式没有什么特别的。

## envionment
初始化测试环境的工具

当前我这里面存放了初始化环境的配置文件和db的建表文件。

## report
存放报告的地址

在tester目录下运行：`sh coverage.sh`

会在 report 下生成 coverage.out 和 coverage.html，

## goconvey

```shell
$ goconvey # 启动网页
http://127.0.0.1:8080/composer.html
http://127.0.0.1:8080/
```


