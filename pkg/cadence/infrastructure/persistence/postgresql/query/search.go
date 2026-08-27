package query

import (
	"strings"

	"github.com/nightnoryu/go-kita/maybe"
)

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// likePattern turns a raw search query into an ILIKE substring pattern, escaping LIKE
// wildcards in the input so fuzzy search matches the literal query
func likePattern(searchQuery maybe.Maybe[string]) maybe.Maybe[string] {
	value, ok := maybe.JustValid(searchQuery)
	if !ok {
		return searchQuery
	}
	return maybe.NewJust("%" + likeEscaper.Replace(value) + "%")
}
