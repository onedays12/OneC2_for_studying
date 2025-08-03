# 介绍

一个用于教学目的C2项目，完整项目要等几个月后了，功能极其不完善，请勿用于实战！

# 使用

来到Beacon目录，构建一个可执行文件
```
go build main.go
```

来到OneServer目录构建一个可执行文件
```
go build main.go
```

如果遇到依赖问题，则使用
```
go mod tidy
```
`C2.postman_collection` 导入到postman，直接拖拽文件即可

![](https://images-of-oneday.oss-cn-guangzhou.aliyuncs.com/images/2025/08/03/22-10-03-9e03b1d9b331044d09065e4054cdcb82-20250803221003-3fe7e2.png)

运行OneServer和Beacon的可执行文件

用postman发送请求来测试

**创建监听器**

![](https://images-of-oneday.oss-cn-guangzhou.aliyuncs.com/images/2025/08/03/22-12-56-5deabe0848dacb962e1a19ffe25a9c88-20250803221255-9c4dd6.png)

**创建任务**（一定要Beacon完成上线才可以创建任务，BeaconID是硬编码的，其实很多配置都是硬编码的）

![](https://images-of-oneday.oss-cn-guangzhou.aliyuncs.com/images/2025/08/03/22-14-06-2bb6b9e43142880f1a5c3c66f509bf56-20250803221406-b037f8.png)

**输出任务结果**

![](https://images-of-oneday.oss-cn-guangzhou.aliyuncs.com/images/2025/08/03/23-37-42-05f079bc91dedda5233f1fa4eabc33a2-20250803233741-d5e2bf.png)

请不用测试 `/api/beacon/generate`，因为它生成的文件并不是beacon，只是用于验证能否patch配置！因为sleep参数我改成了整数而不是字符串，所以生成的可执行文件不能反序列化，如果想成功反序列化，请参考我写的文章！

![](https://images-of-oneday.oss-cn-guangzhou.aliyuncs.com/images/2025/08/03/22-21-41-82aa694859da3bb5d0cc76db5df2f5c5-20250803222141-9dcca8.png)

# 更多细节

如果你对一步一步构建C2感兴趣或者想了解更多细节，请阅读我写的这篇文章：[从零开始手搓C2框架-先知社区](https://xz.aliyun.com/news/18564)

我的旧博客：[关于这个博客 | onedaybook](https://oneday.gitbook.io/onedaybook)

我的新博客（还在弄，过一段时间）：

**欢迎各位师傅交换友链！**

# 参考资料

1、[Adaptix-Framework/AdaptixC2](https://github.com/Adaptix-Framework/AdaptixC2?tab=readme-ov-file)

2、[sliver/server at master · BishopFox/sliver](https://github.com/BishopFox/sliver/tree/master/server)
3、[HavocFramework/Havoc: The Havoc Framework](https://github.com/HavocFramework/Havoc)

4、[mai1zhi2/SharpBeacon: CobaltStrike Beacon written in .Net 4 用.net重写了stager及Beacon，其中包括正常上线、文件管理、进程管理、令牌管理、结合SysCall进行注入、原生端口转发、关ETW等一系列功能](https://github.com/mai1zhi2/SharpBeacon)

5、[testxxxzzz/geacon_pro: 重构了Cobaltstrike Beacon，行为对国内主流杀软免杀，支持4.1以上的版本。 A cobaltstrike Beacon bypass anti-virus, supports 4.1+ version.](https://github.com/testxxxzzz/geacon_pro)

6、[chainreactors/malefic: IoM implant, C2 Framework and Infrastructure](https://github.com/chainreactors/malefic)