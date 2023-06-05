package es

import (
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/m4n5ter/pstash/stash/config"
)

func GetWriterIndexerFromESConf(esConf config.ElasticSearchConf) (writer *Writer, indexer *Index) {
	if esConf.Hosts == nil || len(esConf.Hosts) == 0 {
		return nil, nil
	}

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(esConf.Hosts...),
		elastic.SetBasicAuth(esConf.Username, esConf.Password),
	)
	logx.Must(err)

	writer, err = NewWriter(esConf)
	logx.Must(err)

	var loc *time.Location
	if len(esConf.TimeZone) > 0 {
		loc, err = time.LoadLocation(esConf.TimeZone)
		logx.Must(err)
	} else {
		loc = time.Local
	}
	indexer = NewIndex(client, esConf.Index, loc)

	return writer, indexer
}
