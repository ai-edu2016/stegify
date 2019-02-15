# stegify
[![GoDoc](https://godoc.org/github.com/DimitarPetrov/stegify?status.svg)](https://godoc.org/github.com/DimitarPetrov/stegify)
[![Go Report Card](https://goreportcard.com/badge/github.com/DimitarPetrov/stegify)](https://goreportcard.com/report/github.com/DimitarPetrov/stegify)
[![cover.run](https://cover.run/go/github.com/DimitarPetrov/stegify.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2FDimitarPetrov%2Fstegify)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  


## 说明
这是一个用于将文件隐藏于图片的工具。
这种技术称为LSB (Least Significant Bit) [steganography](https://en.wikipedia.org/wiki/steganography) 

## 演示

| Carrier                                | Data                                | Result                                               |
| ---------------------------------------| ------------------------------------|------------------------------------------------------|
| ![Original File](examples/street.jpeg) | ![Encoded File](examples/lake.jpeg) | ![Encoded File](examples/test_decode.jpeg) |

The `Result` file contains the `Data` file hidden in it. And as you can see it is fully transparent.

## 安装
```
$ go get -u github.com/DimitarPetrov/stegify
```

## 用法

### 作为命令行工具
```
$ stegify -op encode -carrier <file-name> -data <file-name> -result <file-name>
$ stegify -op decode -carrier <file-name> -result <file-name>
```
When encoding, the file with name given to flag `-data` is hidden inside the file with name given to flag
`-carrier` and the resulting file is saved in new file in the current working directory under the
name given to flag `-result`. The file extension of result file is inherited from the carrier file and must not be specified
explicitly in the `-result` flag.

When decoding, given a file name of a carrier file with previously encoded data in it, the data is extracted
and saved in new file in the current working directory under the name given to flag `-result`.
The result file won't have any file extension and therefore it should be specified explicitly in `-result` flag.

In both cases the flag `-result` could be omitted and it will be used the default file name: `result`

### 在代码中使用

你可以看看 [godoc](https://godoc.org/github.com/DimitarPetrov/stegify) 中的详细信息。

## 说明
如果载体文件是jpeg或者jpg格式，在编码后返回的文件为png格式，尽管扩展名为jpeg或者jpg。
