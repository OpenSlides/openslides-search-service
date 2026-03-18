package meta

import (
	log "github.com/sirupsen/logrus"
)

// Member is part of the meta model.
type Member struct {
	Type       string
	Required   bool
	Searchable bool
	Analyzer   *string
	Relation   *CollectionRelation
	Order      int32
}

// Clone returns a deep copy.
func (m *Member) Clone() *Member {
	return &Member{
		Type:     m.Type,
		Required: m.Required,
		Order:    m.Order,
	}
}

// RetainStrings returns a function which keeps string type fields in [Retain].
func RetainStrings() func(string, string, *Member) bool {
	return func(k, fk string, f *Member) bool {
		switch f.Type {
		case "string", "HTMLStrict", "text", "HTMLPermissive":
			return true
		default:
			log.Tracef("removing non-string %s.%s\n", k, fk)
			return false
		}
	}
}
