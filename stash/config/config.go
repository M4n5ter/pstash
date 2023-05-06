package config

import (
	"time"

	"github.com/zeromicro/go-zero/core/service"
)

type (
	Condition struct {
		Key   string
		Value string
		Type  string `json:",default=match,options=match|contains"`
		Op    string `json:",default=and,options=and|or"`
	}

	ElasticSearchConf struct {
		Hosts         []string
		Index         string
		DocType       string `json:",default=doc"`
		TimeZone      string `json:",optional"`
		MaxChunkBytes int    `json:",default=15728640"` // default 15M
		Compress      bool   `json:",default=false"`
		Username      string `json:",optional"`
		Password      string `json:",optional"`
	}

	ZincObserveConf struct {
		Schema        string
		Host          string
		Port          int
		Username      string
		Password      string
		Organization  string `json:",default=default"`
		Stream        string `json:",default=default"`
		IngestionType string `json:",default=_multi,options=_json|_bulk|_multi"`
	}

	Filter struct {
		Action     string      `json:",options=drop|remove_field|transfer"`
		Conditions []Condition `json:",optional"`
		Fields     []string    `json:",optional"`
		Field      string      `json:",optional"`
		Target     string      `json:",optional"`
	}

	KafkaConf struct {
		service.ServiceConf
		Brokers    []string
		Group      string
		Topics     []string
		Offset     string `json:",options=first|last,default=last"`
		Conns      int    `json:",default=1"`
		Consumers  int    `json:",default=8"`
		Processors int    `json:",default=8"`
		MinBytes   int    `json:",default=10240"`    // 10K
		MaxBytes   int    `json:",default=10485760"` // 10M
		Username   string `json:",optional"`
		Password   string `json:",optional"`
	}

	NatsConf struct {
		Urls []string
		// 连接的名字，推荐使用，默认为 Pstash
		ConnectionName string `json:",default=Pstash"`
		// 认证类型: Password/Token/NKey/Credentials
		AuthType string `json:",default=none,options=|none|password|token|nkey|creds"`
		// 认证内容:
		// password 	-> Password like : <USER>@<PASSWORD>  admin@$2a$11$qbtrnb0mSG2eV55xoyPqHOZx/lLBlryHRhU3LK2oOPFRwGF/5rtGK
		// token    	-> <YOUR TOKEN>
		// nkey			-> <YOUR NKey Seed File Path>
		// creds		-> <YOUR Credentials File Path>
		// See docs in https://docs.nats.io/using-nats/developer/connecting/userpass ,
		// https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/nkey_auth ,
		// https://docs.nats.io/using-nats/developer/connecting/creds
		AuthContent        string `json:",optional"`
		EnableJetStream    bool   `json:",default=false,options=true|false"`
		Subject            string `json:",default=log"`
		QueueSubscribeName string `json:",optional"`
		NeedAck            bool   `json:",default=false,options=true|false"`
		//TODO:support TLS connection
		//TODO:support reconnection in https://docs.nats.io/using-nats/developer/connecting/reconnect
	}

	TcpInputConf struct {
		IP   string `json:",default=0.0.0.0"`
		Port int    `json:",default=17171"`
		// milliseconds
		Timeout int `json:",default=10000"`
	}

	Cluster struct {
		Input struct {
			Kafka KafkaConf    `json:",optional"`
			Nats  NatsConf     `json:",optional"`
			Tcp   TcpInputConf `json:",optional"`
		}
		Filters []Filter `json:",optional"`
		Output  struct {
			ElasticSearch ElasticSearchConf `json:",optional"`
			ZincObserve   ZincObserveConf   `json:",optional"`
		}
	}

	Config struct {
		Clusters    []Cluster
		GracePeriod time.Duration `json:",default=10s"`
	}
)
