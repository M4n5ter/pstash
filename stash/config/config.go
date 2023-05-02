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

	TcpInputConf struct {
		IP   string `json:",default=0.0.0.0"`
		Port int    `json:",default=17171"`
		// milliseconds
		Timeout int `json:",default=10000"`
	}

	Cluster struct {
		Input struct {
			Kafka KafkaConf    `json:",optional"`
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
