# alertmanager-webhook-feishu

## 自定义机器人

首先你需要对飞书机器人有点[了解](https://open.feishu.cn/document/ukTMukTMukTM/uATM04CMxQjLwEDN) 。当前使用 [自定义机器人](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN) 转发来自 AlertManager 的告警。

自定义机器人相比其他机器人的好处是免应用创建，不需要维护 app_id/secret。

[其他区别](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN) @20210729：

1. 自定义机器人只能用于在群聊中自动发送通知，不能响应用户@机器人的消息，不能获得任何的用户、租户信息。 
2. 自定义机器人可以被添加至外部群使用，机器人应用只能在内部群使用。
3. 自定义机器人发送的消息卡片，只支持跳转链接的交互方式，无法发送包含回传交互元素的消息卡片，不能通过消息卡片进行信息收集
4. 在消息卡片中如果要@提到某用户，请注意：自定义机器人仅支持通过 open_id 的方式实现，暂不支持email、user_id等其他方式。

数据链路：Prometheus -> AlertManager -> alertmanager-webhook-feishu -> Feishu。

## 功能列表

- [x] 支持多个飞书机器人（一个飞书机器人对应一个群）
- [ ] [签名校验](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#348211be)
- [ ] 配置可 reload
- [ ] 自定义飞书模板
- [x] @所有人 支持
- [x] @某人 支持 open_id
- [x] @某人 支持 email  @20210729 [官方没有直接支持](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#4996824a) ，我们可以申请一个飞书应用，通过 email 获取 open_id。

### open_id 获取方式

[open_id](https://open.feishu.cn/document/home/user-identity-introduction/open-id) 理论上用于标识一个 User ID 在具体某一应用中的身份。但是经过实践发现，随便某个应用下的 open_id 可以用于自定义机器人 @某人。很奇怪的设计。



