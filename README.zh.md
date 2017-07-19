# iOS推送服务

我们直连APNs. 极光送达率真是感人...

本项目使用Redis + Golang打造，API请求延迟在1ms以内。

Table of Contents
=================

   * [iOS推送服务](#ios推送服务)
      * [使用指南](#使用指南)
      * [API](#api)
         * [给设备打tag](#给设备打tag)
         * [给某个设备发送推送](#给某个设备发送推送)
         * [清除badge](#清除badge)
         * [按tag推送](#按tag推送)
         * [上报设备信息](#上报设备信息)


使用方法：

```bash
$ git clone https://github.com/jiajunhuang/obito
$ cd obito
$ go build
$ # 另外需要crontab配合使用，定期清除已过期的key
$ cd cron
$ go build
```

**开发工程中感谢iOS客户端童鞋的大力支持，谢谢！**

> redis中的key都比较长，后续如果内存吃紧，可以考虑使用更加短的方案，但是目前来说
> 还在可接受范围内

## 使用指南

首先，需要iOS端配合做好以下工作：

- app在首次启动之后在沙箱内生成一个 `uuid` 作为唯一识别码
- 每次打开app的时候获取 `device_token` 并且和上一次对比，如果发生了变化，则需要
  调用 `/report` 上报新的 `device_token` ，完成 `UUID` 和 `device_token` 配对的更新
- 每次打开app需要调用 `/badge/clear` 接口清除服务端存储的badge数量。

## API

所有操作，如果成功，则返回：

状态码为200且json为

```json
{
    "message": "",
    "result": {}
}
```

如果失败，则返回：

状态码为400，500等，且json为

```json
{
    "message": "具体错误提示",
    "result": {}
}
```

### 给设备打tag

`POST /tag`

```json
{
    "uuid": "12345678-1730-4414-9728-95616073fe82",
    "tag_list": ["tag1", "tag2"]
}
```

### 给某个设备发送推送

`POST /push`

```json
{
    "uuid": "12345678-1730-4414-9728-95616073fe82",
    "content": "Hello, this is a test with very long somewhat"
}
```

### 清除badge

`PUT /badge/clear`

```json
{
    "uuid": "12345678-1730-4414-9728-95616073fe82"
}
```

### 按tag推送

`POST /tag/push`

```json
{
    "tag": "tag1",
    "content": "push_by_tag"
}
```

### 上报设备信息

`POST /report`

```json
{
    "uuid": "f8206027-8752-4850-bb20-3ae0f23d082e",
    "device_token": "12345678-1730-4414-9728-95616073fe82"
}
```

-------------------

MIT License
