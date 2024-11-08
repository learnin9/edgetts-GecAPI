# 说明
该代码为自动获取edgetts的Sec-MS-GEC和Sec-MS-GEC-Version的方法，并提供了接口服务给自己的edgetts调用使用。
# 使用方法
1. 使用mitmproxy进行了中间人截取请求，并调用python程序获取到Sec-MS-GEC和Sec-MS-GEC-Version数据
2. 使用html+js进行轮询进行edgetts的文字播放
3. 使用了windows批处理进行edge浏览器的打开html文件和自动关闭edge浏览器
4. 当python代码拿到了sec后，通过接口方式进行http post到本地的gin服务
5. gin服务收到了post请求后，将json数据存储到redis中，并计算ttl值，过期消失
6. gin服务对外提供GET请求接口，返回sec的json结果

# 目录说明
## startGec目录
### index.html
此文件为edge浏览器自动打开播放文字转语音的，此文件将由start.bat进行启动
### server.py
这个是mitmproxy启动时需要调用的server.py程序，他会自动获取到Sec-MS-GEC和Sec-MS-GEC-Version，你可以进行post到gin服务，也可以直接写入到redis中去。自己随意。
### start.bat
鼠标双击执行，不过里面的代码执行的路径需要自己改下路径。

```bash
start "" "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" "c:\gec\index.html"
```
自己edge浏览器安装路径需要自己找了，index.html文件需要新建目录存放。
# 使用方法
1. 先启动mitmproxy，执行命令：
```bash
pip install mitmproxy -i https://pypi.tuna.tsinghua.edu.cn/simple
```
2. 启动server.py，执行 
```bash 
mitmproxy -s c:\gec\server.py
```
3. 自行安装go 1.2x，可以自己编译，执行build.bat自动编译exe可执行文件,也可以使用我编译好的go-edgetts.exe文件. 
4. 设置edge浏览器代理模式指向 127.0.0.1的8080端口 
5. 双击执行 start.bat
# 注意事项
1. 需要自己导入mitmproxy证书到浏览器的信任根服务，具体方法参照 https://blog.csdn.net/qq_36841447/article/details/134012335
2. 惭愧的是，本人不擅长js,所以一直没找到好的方法让在不关闭浏览器的情况下重头进行edgetts调用，我发现即使做了轮询方法，edgetts并不会再次调用请求。所以只能靠bat批处理来进行自动重启edge浏览器重新加载html方式了。如果谁有更好的解决方法，欢迎提出来。
3. 只能使用windows 10以上系统，我目前使用的是腾讯云99元/年的2核2G的windows 2022 server版本。