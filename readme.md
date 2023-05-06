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