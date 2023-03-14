package internal

import "embed"

//go:embed view
// Assets
//这段代码是使用 Go 1.16 中新增的 embed 包定义了一个名为 Assets 的变量，它的类型是 embed.FS。
//该变量使用了 //go:embed 指令，表示将 view 目录中的所有文件和子目录嵌入到编译后的二进制文件中。这些文件可以通过 Assets 变量来访问。
//通过这种方式，我们可以将程序所需的资源文件打包到二进制文件中，而无需在运行时读取外部文件，从而简化了程序的部署和分发。
var Assets embed.FS
