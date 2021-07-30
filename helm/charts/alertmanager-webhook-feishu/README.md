## 配置说明

目录结构如下：

```
$ tree -L 1
.
├── Chart.yaml
├── README.md
├── charts
├── config
├── templates
└── values.yaml

```

config 目录存放 config.yml 及自定义模板。config/config.example.yml 为示例文件。

config 目录会被映射到容器目录：`/etc/alertmanager-webhook-feishu`

应用默认参数为：`--config=/etc/alertmanager-webhook-feishu/config.yml  ` `--email`

详见 templates/deployment.yaml。

## 安装



```
helm upgrade -i alertmanager-webhook-feishu ./ -n monitoring --create-namespace
```

