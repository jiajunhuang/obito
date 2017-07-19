# iOS notification service

[简体中文README](./README.zh.md)

This project connect to APNs directly, it use [Golang](https://golang.org/) and [Redis](https://redis.io/) only.

> I'm writing unittests, hold-on :)

Table of Contents
=================

   * [iOS notification service](#ios-notification-service)
      * [Usage](#usage)
         * [Server side](#server-side)
         * [iOS side](#ios-side)
      * [API](#api)
         * [Tag your device(so you can push to lots of devices just use a exsiting tag)](#tag-your-deviceso-you-can-push-to-lots-of-devices-just-use-a-exsiting-tag)
         * [Send a notification to a device](#send-a-notification-to-a-device)
         * [Clear badge number stored in server side](#clear-badge-number-stored-in-server-side)
         * [Push to lots of device](#push-to-lots-of-device)
         * [Set the relationship between your device and device token](#set-the-relationship-between-your-device-and-device-token)
      * [FAQ](#faq)

## Usage

### Server side

```bash
$ git clone https://github.com/jiajunhuang/obito
$ cd obito
$ go build
$ # set a crontab job, too, use crontab to remove expired keys
$ cd cron
$ go build
```

> currently we use UUID + project name as key, it may waste some memory,
> for example, UUID is 36 characters, and device token get from APNs
> is 64 characters long, so if we have 100,0000 users, we need about 100M
> memory.

### iOS side

- Your application should create a UUID and store it in sandbox, use the UUID to identify device
- Everytime the application is launched, check if the device token changed, if it changed, then post a request to tell the server side to refresh relationship between device token and UUID
- Everytime the application is launched, post a request to server side to tell your server to clear the badge

## API

first we have to define a universal response, if the request is succeed, it will return:

```json
{
    "message": "",
    "result": {}
}
```

with http status code set to 200, if it fails, it will return:

```json
{
    "message": "why it fails",
    "result": {}
}
```

with http status code set to, maybe 400, or 500...etc.

### Tag your device(so you can push to lots of devices just use a exsiting tag)

`POST /tag`

```json
{
    "uuid": "12345678-1730-4414-9728-95616073fe82",
    "tag_list": ["tag1", "tag2"]
}
```

### Send a notification to a device

`POST /push`

```json
{
    "uuid": "12345678-xxxx-1234-1234-abcdefghijkl",
    "content": "Hello, this is a test with very long somewhat"
}
```

### Clear badge number stored in server side

`PUT /badge/clear`

```json
{
    "uuid": "12345678-1730-4414-9728-95616073fe82"
}
```

### Push to lots of device

`POST /tag/push`

```json
{
    "tag": "tag1",
    "content": "push_by_tag"
}
```

### Set the relationship between your device and device token

`POST /report`

```json
{
    "uuid": "f8206027-8752-4850-bb20-3ae0f23d082e",
    "device_token": "12345678-1730-4414-9728-95616073fe82"
}
```

## FAQ

- when I post a bad json, the response content type is `content/text`? why?

    it's a bug of [gin](https://github.com/gin-gonic/gin/issues/633), I'm trying to find out the reason and fix it, hold on.

-------------------

MIT License
