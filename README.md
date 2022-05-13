# Go2CAnalysis

本程序使用基于go-ast的go/parser批量检测go仓库是否含有import "C"特征。存在该特征的仓库结果保存在`Go2CDir.json`下。并且对存在import "C"特征的仓库分析其内部的GO-C函数调用情况，以及在注释中include的头文件。结果保存在对应的csv文件。
总的结果放在当前目录下，比如`./func_detail.csv`。每个仓库的结果放在每个仓库对应的目录下。比如`go`仓库的include分析结果在`./Include/go/go.csv`。
结尾带`_detail`的csv文件是详细结果，包括平均值，最大值，最大仓库等。按照引用总次数降序排列。

使用命令行指定仓库存放的路径，以及要输出的信息等

usage :
`-h` help 
`-dumpdetail` 将所有信息写入文件，包括路径等。不指定时不会将路径信息写入文件
`-showall` 输出所有提示信息，不指定时只输出错误信息和分析的仓库名称
`-path` 指定要分析的目录(即要分析的所有仓库的存放目录)，默认为`./`
`-core` 核心条目个数，指数量最多的x(default:10)个条目
`-top` top项目，指拥有核心条目最多的x(default:5)个仓库

示例:假设要分析的仓库均放在`/data/github_go/repos`下，且要输出提示信息，
```shell
$ go build ./
$ ./anatool -h
  Usage of ./anatool:
  -core int
    	number of core item (default 10)
  -dumppath
    	dump item path
  -path string
    	dir path of the repos (default "./")
  -showall
    	show output details
  -top int
    	top x repos having the core item (default 5)
$ ./anatool -path '/data/github_go/repos' -showall
```