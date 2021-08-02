# alertmanager-webhook-feishu

## 自定义机器人

首先你需要对飞书机器人有点[了解](https://open.feishu.cn/document/ukTMukTMukTM/uATM04CMxQjLwEDN) 。我们使用 [自定义机器人](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN) 转发来自 AlertManager 的告警。

自定义机器人相比应用机器人的好处是免应用创建，不需要维护 app_id/secret。

[其他区别](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN) @20210729：

1. 自定义机器人只能用于在群聊中自动发送通知，不能响应用户@机器人的消息，不能获得任何的用户、租户信息。 
2. 自定义机器人可以被添加至外部群使用，机器人应用只能在内部群使用。
3. 自定义机器人发送的消息卡片，只支持跳转链接的交互方式，无法发送包含回传交互元素的消息卡片，不能通过消息卡片进行信息收集。
4. 在消息卡片中如果要@提到某用户，请注意：自定义机器人仅支持通过 open_id 的方式实现，暂不支持email、user_id等其他方式。

## 使用

基于参考配置文件 [config.example.yml](config.example.yml) 调整。运行：

```bash
$ alertmanager-webhook-feishu server --config=/path/to/config.yml
```

数据链路：Prometheus -> AlertManager -> alertmanager-webhook-feishu -> Feishu。AlertManager 配置文件中添加：

```yaml
receivers:
  - name: 'you_name_it'
    webhook_configs:
    - url: 'http://alertmanager-webhook-feishu.monitoring/hook/your_group_name'
      send_resolved: true
```

完整配置参数：

```bash
$ alertmanager-webhook-feishu server -h
start webhook server

Usage:
  alertmanager-webhook-feishu server [flags]

Flags:
  -c, --config string   config file for bot webhook
  -e, --email           if email supported, need feishu appid/secret for enabling
  -h, --help            help for server
  -p, --port int        server port (default 8000)

Global Flags:
  -v, --verbose   show verbose log

```



### helm chart

见 [helm/charts/alertmanager-webhook-feishu](helm/charts/alertmanager-webhook-feishu)

## 功能列表

- [x] 支持多个飞书机器人
- [ ] [签名校验](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#348211be)
- [ ] 配置可 reload（配置更新，K8s 上可滚动部署，reload 功能不是必须的）
- [x] 自定义飞书模板
- [x] @所有人 支持
- [x] @某人 支持 open_id
- [x] @某人 支持 email  @20210729 [官方没有直接支持](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#4996824a) ，我们可以申请一个飞书应用，通过 email 获取 open_id。

### open_id 获取方式

[open_id](https://open.feishu.cn/document/home/user-identity-introduction/open-id) 理论上用于标识一个 User ID 在具体某一应用中的身份。但是经过实践发现，随便某个应用下的 open_id 可以用于自定义机器人 @某人。很奇怪的设计。

目前，为了能通过 email 得到 open_id，需要创建飞书「企业自建应用」。需要注意如下两点：

1. 应用可用性：可用成员为「全部成员」。
2. 最小权限：「通过手机号或邮箱获取用户 ID」。

### 飞书消息模板

项目本质上是将数据结构 [model.WebhookMessage](model/model.go) 生成为 [飞书消息格式](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#8b0f2a1b) 发送到飞书服务器。

#### 默认模板

当前的默认消息模板为「消息卡片」，见 [tmpl/templates/default.tmpl](tmpl/templates/default.tmpl)。其中，为每个 Alert 对象模板化输出了 markdown 文本，[飞书 markdown 只支持部分标签](https://open.feishu.cn/document/ukTMukTMukTM/uADOwUjLwgDM14CM4ATN)。（为 Alert 对象输出 markdown 文本使用了[独立的模板](tmpl/templates/default_alert.tmpl)，发现一个模板搞不定 - -!）

默认模板适配 [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus) 项目。比如，在模板中会使用如下 [Alert](https://github.com/prometheus-operator/kube-prometheus/blob/main/manifests/kubernetes-prometheusRule.yaml#L15) 的 labels 和 annotations。

```yaml
spec:
  groups:
  - name: kubernetes-apps
    rules:
    - alert: KubePodCrashLooping
      annotations:
        description: Pod {{ $labels.namespace }}/{{ $labels.pod }} ({{ $labels.container }}) is restarting {{ printf "%.2f" $value }} times / 10 minutes.
        runbook_url: https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubepodcrashlooping
        summary: Pod is crash looping.
      expr: |
        increase(kube_pod_container_status_restarts_total{job="kube-state-metrics"}[10m]) > 0
        and
        kube_pod_container_status_waiting{job="kube-state-metrics"} == 1
      for: 15m
      labels:
        severity: warning

```

主要是如下字段：

1. annotations.description 描述中已包含了指标 labels 的模板化输出。
2. annotations.runbook_url 操作手册地址。（kube-prometheus 项目存在操作手册地址不存在或是更新不及时的情况，需要自己微调）
3. annotations.summary 总结，信息量上，与 alertname 有点重合。
4. labels.severity 严重程度。因为指标 labels 都有在 annotations.description 中有体现，其他 labels 没有必要在消息模板中进一步处理。

#### 自定义模板

为了满足个性化需求，可以自定义模板。模板基于 go [text/template](https://golang.org/pkg/text/template/)。参考这个[链接](https://stackoverflow.com/questions/55170279/go-text-template-syntax-highlighting-in-goland)，配置编辑器的语法提示，提高配置速度。

使用自定义模板：

```yaml
bots:
  webhook:
    url: https://open.feishu.cn/open-apis/bot/v2/hook/xxx
    template:
    	custom_path: /path/to/custom/path

```

