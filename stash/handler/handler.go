package handler

import (
	json "github.com/json-iterator/go"
	"github.com/m4n5ter/pstash/stash/es"
	"github.com/m4n5ter/pstash/stash/filter"
	"unsafe"
)

type MessageHandler struct {
	writer  *Writer
	indexer *es.Index
	filters []filter.FilterFunc
}

func NewHandler(writer *Writer, indexer *es.Index) *MessageHandler {
	return &MessageHandler{
		writer:  writer,
		indexer: indexer,
	}
}

func (mh *MessageHandler) AddFilters(filters ...filter.FilterFunc) {
	for _, f := range filters {
		mh.filters = append(mh.filters, f)
	}
}

func (mh *MessageHandler) Consume(_, val string) error {
	var m map[string]any
	var index string
	if err := json.Unmarshal([]byte(val), &m); err != nil {
		return err
	}

	if mh.writer.ESWriter != nil {
		index = mh.indexer.GetIndex(m)
	}

	for _, proc := range mh.filters {
		if m = proc(m); m == nil {
			return nil
		}
	}

	bs, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return mh.writer.Write(index, bytes2string(bs))
}

func bytes2string(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
