// 注： 首先需要安装go-zero的goctl工具
go install github.com/zeromicro/go-zero/tools/goctl@latest

// 1.根据api文件生成md文档
cd nightclub/desc
goctl api doc -dir nightclub.api

// 2. 根据api文件生成handler和logic

// 2.1 如果是在根项目路径下
goctl api go -api ./nightclub/desc/nightclub.api -dir .\nightclub\ -style goZero

// 2.1 如果是在 cd nightclub/desc 之后
goctl api go -api nightclub.api -dir ../ -style goZero