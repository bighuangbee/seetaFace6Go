# SeetaFace6 CGO
目的是在golang集成seetaFace6，所以重点关注如何封装和调用接口。
从上往下分层:
1. Golang代码，业务逻辑，比如输入图像、获取seetaFace的处理结果
2. C++ Warp封装层，往下封装了SeetaFace6的调用环境，往上可被Go直接调用
3. SeetaFace6官方库

## SeetaFace6编译
SeetaFace6Open文档进行指引，有编译顺序和环境要求。

- ### mac和linux [仓库](https://gitee.com/bighuangbee/seeta-face6-open)
编译比较容易，安装gcc和Make即可。

- ### windows
依赖MSVC，需要一定熟悉程度，否则容易各种报错。
不建议自己编译，可使用[已编译好的动态库](https://github.com/bighuangbee/SeetaFace6OpenBinary)。

**坑点**:
可能由于编译参数的原因，自己编译处理的dll，运行测试代码发现耗时很久，detect执行一次要一秒多。
一开始以为是封装层dll造成的，按以往的经验调用在cgo的效率损失不应该这么大，但还是做下排除法，mac测试了没这样的问题，于是在C代码下测试（test/main.c）,也没这样的问题，所以可以确定是编译seetaFace的dll有问题。

## C++ Warp封装层
- ### 代码封装
    seetaFace6提供二进制接口是C++风格，不能直接在CGO调用，通常需要创建一个C++封装层，将 C++函数封装为C函数。原因：
  1. C++编译器对函数名进行名称修饰，C语言无法找到这些函数，需要使用extern "C"确保函数以C的方式导出。
  2. 很多C++特性（如类、重载）在C中是不存在的
  
- ### 编译参数
  1. Mac和Linux，不是必须编译出C++封装层的动态库。seetaFace6源码、C++ Warp封装层和golang的代码是同一个编译器GCC，只需要在环境变量或编译参数中指定链接的头文件CFLAGS和动态库路径LDFLAGS，见face_cgo.go和Makefile，比较容易。
  2. Windows，C++代码通过MSVC编译出来的DLL，与MinGW等GCC工具链不兼容，MSVC与MinGW的ABI（应用二进制接口）不同，不能混用。只有两种解决办法：
    - 1. 使用Ming编译seetaFace6源码，像Mac那样确保使用同一个编译器，但没找到解决方法。手动修改编译脚本遇到太多错误，对C++不熟悉，遂放弃。
    - 2. 使用MSVC编译CGO代码，遇到错误`cl: 命令行 error D8021 :无效的数值参数“/Werror`，无法搞定编译参数，穗放弃。
    - 3. 使用MSVC编译C++ Warp封装层，导出Extern C接口，生成动态库Warp.dll, 然后在Mingw编译golang，通过cgo链接封装层的Warp.dll，再调用底层seetaFace库。





