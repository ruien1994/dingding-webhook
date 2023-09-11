# dingding-webhook
用于接收{"text":"数据内容"}格式的告警信息，把这个告警转给钉钉告警
## 用法
curl -X POST -H "Content-Type: application/json" -d '{"text":"S1:告警内容"}' http://127.0.0.1:8080/webhook/钉钉的YOUR_ACCESS_TOKEN
