package main

import (
	"flag"
	"github.com/m4n5ter/pstash/stash/zo"
	"time"

	"github.com/m4n5ter/pstash/stash/config"
	"github.com/m4n5ter/pstash/stash/es"
	"github.com/m4n5ter/pstash/stash/filter"
	"github.com/m4n5ter/pstash/stash/handler"
	"github.com/olivere/elastic/v7"
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

	group := service.NewServiceGroup()
	defer group.Stop()

	for _, processor := range c.Clusters {
		ew, indexer := getWriterIndexerFromESConf(processor.Output.ElasticSearch)
		zw := zo.NewWriter(processor.Output.ZincObserve)
		filters := filter.CreateFilters(processor)
		handle := handler.NewHandler(&handler.Writer{ESWriter: ew, ZOWriter: zw}, indexer)
		handle.AddFilters(filters...)
		handle.AddFilters(filter.AddUriFieldFilter("url", "uri"))
		for _, k := range toKqConf(processor.Input.Kafka) {
			group.Add(kq.MustNewQueue(k, handle))
		}
	}

	group.Start()
}

func getWriterIndexerFromESConf(esConf config.ElasticSearchConf) (writer *es.Writer, indexer *es.Index) {
	if esConf.Hosts == nil || len(esConf.Hosts) == 0 {
		return nil, nil
	}

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(esConf.Hosts...),
		elastic.SetBasicAuth(esConf.Username, esConf.Password),
	)
	logx.Must(err)

	writer, err = es.NewWriter(esConf)
	logx.Must(err)

	var loc *time.Location
	if len(esConf.TimeZone) > 0 {
		loc, err = time.LoadLocation(esConf.TimeZone)
		logx.Must(err)
	} else {
		loc = time.Local
	}
	indexer = es.NewIndex(client, esConf.Index, loc)

	return writer, indexer
}
