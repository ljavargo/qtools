```
prtool version: 1.0.0
Options:
  -closeAll
    	关闭当前用户指定project的所有PR
  -g string
    	group名称
  -h	帮助
  -p string
    	必填，project名称
  -s string
    	source branch
  -t string
    	target branch, 支持两种模式，批量：'branch1,branch2,branch3' 链式：'branch1>branch2>branch3' (注意：要加单引号防止>被处理为重定向符号)  (default "master")
  -tk string
    	必填，private token
  -tt string
    	pr title (default "create pr by prtool")
  -u string
    	必填，gitlab主页url地址 (default "https://git.example.com")
```
