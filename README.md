# ava
一个去中心化的分布式任务运行平台,最终版应是ci/cd,自动扩容,协议包伪装,运行,监控,分流等等等的xx平台

#### 1. 开发(弃坑)进度
- [x] 双端断线,tcp,ws心跳检测自动重连
- [x] ws任务命令通道
- [x] 反向socks5代理数据通道
- [x] 接收web请求执行cmd命令
- [x] 节点启动注册业务
- [x] 管理端路由分发,解析投送,定点投送
- [x] 内网穿透白名单
- [ ] 根据主机和业务负载分流
- [ ] 通过任务Id管控节点进程
- [ ] 热更新配置
- [ ] ws消息体结构定义
- [ ] pac等节点代理无感知



#### 2. 启动方式
##### D的运行
会监听本机4000端口用于接收web命令,连接多台运行节点

```bash
main config.josn
```

##### H的运行

H监听端口 websocket: 4560, tcp: 4561, socks5: 4562 
```bash
./main
```
#### 3. 配置文件
节点配置,白名单配置文件config.json,仅管理端需要
```json
{
    "nodes": [
        "172.16.102.199",
        "172.16.102.3"
    ],
    "sites": [
            "172.16.102.199",  
            "172.16.102.3"
        ]
}
```


##### 4. 部分api
web状态查看
```bash
http://127.0.0.1:4000/webWsStatus              --连接状态
http://127.0.0.1:4000/webWorkerMapR            --节点<-->任务对应关系
http://127.0.0.1:4000/info                     --主机状态
```
发送命令请求: POST  http://127.0.0.1:4000/exectask
```bash
{
"worker": "gather-register",
"task_id": "qerty",
"params": "eyJtZXRob2QiOiAiZmFrZS5lY2hvIiwgInBhcmFtcyI6IHsiYSI6IDEyM319"
"route": "192.168.169.128"   ---(可选,定点投送)
}
```




##### 4 bug
断开连接后,第一次发送请求,依然会报成功