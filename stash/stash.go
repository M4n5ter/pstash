package main

import (
	"flag"
	"github.com/m4n5ter/pstash/stash/config"
	"github.com/m4n5ter/pstash/stash/es"
	"github.com/m4n5ter/pstash/stash/filter"
	"github.com/m4n5ter/pstash/stash/handler"
	"github.com/m4n5ter/pstash/stash/natsinput"
	"github.com/m4n5ter/pstash/stash/tcpinput"
	"github.com/m4n5ter/pstash/stash/zo"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/config.yaml", "Specify the config file")

func toKqConf(c config.KafkaConf) []kq.KqConf {
	var ret []kq.KqConf

	for _, topic := range c.Topics {
		ret = append(ret, kq.KqConf{
			ServiceConf: c.ServiceConf,
			Brokers:     c.Brokers,
			Group:       c.Group,
			Topic:       topic,
			Offset:      c.Offset,
			Conns:       c.Conns,
			Consumers:   c.Consumers,
			Processors:  c.Processors,
			MinBytes:    c.MinBytes,
			MaxBytes:    c.MaxBytes,
			Username:    c.Username,
			Password:    c.Password,
		})
	}

	return ret
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	proc.SetTimeToForceQuit(c.GracePeriod)
	logx.DisableStat()

	group := service.NewServiceGroup()
	defer group.Stop()

	for _, processor := range c.Clusters {
		// OUTPUT
		ew, indexer := es.GetWriterIndexerFromESConf(processor.Output.ElasticSearch)
		zw := zo.NewWriter(processor.Output.ZincObserve)
		// must have at least one output
		if ew == nil && zw == nil {
			panic("no output")
		}

		// FILTER
		filters := filter.CreateFilters(processor)
		hdl := handler.NewHandler(&handler.Writer{ESWriter: ew, ZOWriter: zw}, indexer)
		hdl.AddFilters(filters...)
		hdl.AddFilters(filter.AddUriFieldFilter("url", "uri"))

		// INPUT
		// must have at least one input
		if processor.Input.Nats.Urls == nil && processor.Input.Kafka.Brokers == nil && processor.Input.Tcp.IP == "" {
			panic("no input")
		} else if processor.Input.Nats.Urls != nil {
			group.Add(natsinput.NewNats(processor.Input.Nats, hdl))
		} else if processor.Input.Tcp.IP != "" {
			group.Add(tcpinput.NewTcpInput(processor.Input.Tcp, hdl))
		} else if processor.Input.Kafka.Brokers != nil {
			for _, k := range toKqConf(processor.Input.Kafka) {
				group.Add(kq.MustNewQueue(k, hdl))
			}
		}
	}

	group.Start()
}
