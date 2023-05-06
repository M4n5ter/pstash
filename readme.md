## pstash

**pstash** 基于 [go-stash](https://github.com/kevwan/go-stash) 进行开发，目的是为了支持更多的 Input 和 Output

目前支持的 Input 有：
* Kafka
* Tcp
* [NATS](https://nats.io/)

目前支持的 Output 有：
* ElasticSearch
* [ZincObserve](https://github.com/zinclabs/zincobserve)

## RoadMap
* [x] 支持 ZincObserve Output
* [x] 支持 Tcp Input
* [x] 支持更多的 Input，比如 [NATS](https://nats.io/)
* [ ] 剥离 go-zero 相关依赖


## Configuration

### Nats Input
```yaml
Nats:
  Urls:
      - nats://127.0.0.1:4222 #required
      - nats://1.1.1.1:4222
  ConnectionName: "pstash" # optional, default is Pstash
  AuthType: "none" # default=none,options=|none|password|token|nkey|creds"
  AuthContent: "auth.txt" # optional
  EnableJetStream: true
  Subject: "test-log"
  QueueSubscribeName: "" # optional, default is empty
  NeedAck: true
  TlsCertPath: "cert.pem" # optional
  TlsKeyPath: "key.pem" # optional
  TlsCaPath: "ca.pem" # optional
```

### Tcp Input
```yaml
Tcp:
  IP: "0.0.0.0" # default=0.0.0.0
  Port: 17171 # default=17171
  Timeout: 10000 # default=10000 (millisecond)
```

### ZincObserve Output
```yaml
Output:
  ZincObserve:
    Schema: http
    Host: 127.0.0.1
    Port: 5080
    Username: admin@m4n5ter.email
    Password: "123456"
    Organization: "m4n5ter" # default=default
    Stream: "test" # default=default
    IngestionType: "_multi" # default=_multi,options=_json|_bulk|_multi
```