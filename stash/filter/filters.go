package filter

import "github.com/m4n5ter/pstash/stash/config"

const (
	filterDrop         = "drop"
	filterRemoveFields = "remove_field"
	filterTransfer     = "transfer"
	opAnd              = "and"
	opOr               = "or"
	typeContains       = "contains"
	typeMatch          = "match"
)

type FilterFunc func(map[string]any) map[string]any

func CreateFilters(p config.Cluster) []FilterFunc {
	var filters []FilterFunc

	for _, f := range p.Filters {
		switch f.Action {
		case filterDrop:
			filters = append(filters, DropFilter(f.Conditions))
		case filterRemoveFields:
			filters = append(filters, RemoveFieldFilter(f.Fields))
		case filterTransfer:
			filters = append(filters, TransferFilter(f.Field, f.Target))
		}
	}

	return filters
}
