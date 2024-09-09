package ast

import (
	gcslc "github.com/PlayerR9/go-commons/slices"
)

var (
	// filter_non_nil_noders is a filter that filters out Noders that are nil.
	filter_non_nil_noders gcslc.PredicateFilter[Noder]
)

func init() {
	filter_non_nil_noders = func(child Noder) bool {
		return child != nil
	}
}
