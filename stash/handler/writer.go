package handler

import (
	"fmt"
	"github.com/m4n5ter/pstash/stash/es"
	"github.com/m4n5ter/pstash/stash/zo"
)

type Writer struct {
	ESWriter *es.Writer
	ZOWriter *zo.Writer
}

func (w *Writer) Write(index, val string) error {
	var ese, zoe error
	if w.ESWriter != nil {
		ese = w.ESWriter.Write(index, val)
	}
	if w.ZOWriter != nil {
		zoe = w.ZOWriter.Write(index, val)
	}

	if ese != nil && zoe != nil {
		return fmt.Errorf("es write error: %v, zo write error: %v", ese, zoe)
	}

	if ese != nil && zoe == nil {
		return fmt.Errorf("es write error: %w", ese)
	}

	if ese == nil && zoe != nil {
		return fmt.Errorf("zo write error: %w", zoe)
	}
	return nil
}
