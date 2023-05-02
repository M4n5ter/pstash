package filter

import (
	"strings"

	"github.com/m4n5ter/pstash/stash/config"
)

func DropFilter(conds []config.Condition) FilterFunc {
	return func(m map[string]any) map[string]any {
		var qualify bool
		for _, cond := range conds {
			var qualifyOnce bool
			switch cond.Type {
			case typeMatch:
				qualifyOnce = cond.Value == m[cond.Key]
			case typeContains:
				if val, ok := m[cond.Key].(string); ok {
					qualifyOnce = strings.Contains(val, cond.Value)
				}
			}

			switch cond.Op {
			case opAnd:
				if !qualifyOnce {
					return m
				} else {
					qualify = true
				}
			case opOr:
				if qualifyOnce {
					qualify = true
				}
			}
		}

		if qualify {
			return nil
		} else {
			return m
		}
	}
}
