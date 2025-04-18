# bwhNotify
搬瓦工VPS钉钉机器人消息通知

## 功能

能够每天统计bwh的用量情况，并每天定时发送给钉钉机器人，实现钉钉消息通知。

## 使用方法

1. 配置文件

   更改`config.yml`中内容

   ```
   # 搬瓦工VPS信息
   BWHosts:
     -
       veid: <veid of your vps_1>
       apiKey: <api_key of your vps_1>
     -
       veid: <veid of your vps_2>
       apiKey: <api_key of your vps_2>
     # ...还可以添加多个VPS信息...
   # 钉钉机器人信息
   DingTalk:
     accessToken: <钉钉机器人的access token>
     secret: <钉钉机器人的密钥，如果启用了加签则需要，否则留空>
   ```

   bwh的veid和apiKey可到 https://bwh81.net/services 中Manage - Open KiwiVM页面，找到 KiwiVM Extras - API，可查看。

   钉钉机器人配置方法详见 https://open.dingtalk.com/document/orgapp/robot-overview 中自定义机器人的配置。

2. 启动

   ```bash
   bwhNotify -c config.yml
   ```

3. 系统服务

   需要将`bwhNotify`和`config.yml`放在/usr/local/bwhNotify/

   ```bash
   sudo cp bwhNotify.service /lib/systemd/system/
   sudo systemctl start bwhNotify.service
   sudo systemctl enable bwhNotify.service
   ```

   
