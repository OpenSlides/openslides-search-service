// Package meta implements handling of the meta data model.
package meta

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

var (
	modelNum  atomic.Int32
	fieldNum  atomic.Int32
	filterNum atomic.Int32
)

// Fields is part of the meta model.
type Fields struct {
	Type string `yaml:"type"`
	To   string `yaml:"to"`
}

// MemberTo is part of the meta model.
type MemberTo struct {
	Collections []string `yaml:"collections"`
	Field       string   `yaml:"field"`
}

// Member is part of the meta model.
type Member struct {
	Type                  string    `yaml:"type"`
	Description           string    `yaml:"description"`
	To                    *MemberTo `yaml:"to"`
	Fields                *Fields   `yaml:"fields"`
	ReplacementCollection string    `yaml:"replacement_collection"`
	ReplacementEnum       []string  `yaml:"replacement_enum"`
	RestrictionMode       string    `yaml:"restriction_mode"`
	Required              bool      `yaml:"required"`
	Order                 int32     `yaml:"-"`
}

// Collection is part of the meta model.
type Collection struct {
	Fields map[string]*Member
	Order  int32
}

// Collections is part of the meta model.
type Collections map[string]*Collection

// Filter is part of the meta model.
type Filter struct {
	Name  string
	Items []string
}

// FilterKey is part of the meta model.
type FilterKey struct {
	Name  string
	Order int32
}

// Filters is a list of filters.
type Filters []Filter

func load[T any](r io.Reader) (T, error) {
	dec := yaml.NewDecoder(r)
	var t T
	if err := dec.Decode(&t); err != nil {
		var n T
		return n, err
	}
	return t, nil
}

func fetchRemote[T any](path string) (T, error) {
	resp, err := http.Get(path)
	if err != nil {
		var n T
		return n, err
	}
	if resp.StatusCode != http.StatusOK {
		var n T
		return n, fmt.Errorf("%s failed: %s (%d)",
			path, http.StatusText(resp.StatusCode), resp.StatusCode)
	}
	defer resp.Body.Close()
	return load[T](resp.Body)
}

func fetchLocal[T any](path string) (T, error) {
	in, err := os.Open(path)
	if err != nil {
		var n T
		return n, err
	}
	defer in.Close()
	return load[T](in)
}

// Fetch loads a meta model.
func Fetch[T any](path string) (T, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return fetchRemote[T](path)
	}
	return fetchLocal[T](path)
}

// UnmarshalYAML implements [gopkg.in/yaml.v3.Unmarshaler].
func (fs *Fields) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err == nil {
		fs.Type = s
		return nil
	}
	var field struct {
		Type string `yaml:"type"`
		To   string `yaml:"to"`
	}
	if err := value.Decode(&field); err != nil {
		return fmt.Errorf("fields object without type: %w", err)
	}
	fs.Type = field.Type
	fs.To = field.To
	return nil
}

// UnmarshalYAML implements [gopkg.in/yaml.v3.Unmarshaler].
func (mt *MemberTo) UnmarshalYAML(value *yaml.Node) error {
	// 1. string
	var s string
	if err := value.Decode(&s); err == nil {
		mt.Field = s
		return nil
	}

	// 2. List of strings
	var collections []string
	if err := value.Decode(&collections); err == nil {
		mt.Collections = collections
		return nil
	}

	// 3. struct
	var memberTo struct {
		Collections []string `yaml:"collections"`
		Field       string   `yaml:"field"`
	}
	if err := value.Decode(&memberTo); err != nil {
		return fmt.Errorf("memberTo object without field: %w", err)
	}
	mt.Field = memberTo.Field
	mt.Collections = memberTo.Collections
	return nil
}

// UnmarshalYAML implements [gopkg.in/yaml.v3.Unmarshaler].
func (m *Member) UnmarshalYAML(value *yaml.Node) error {
	m.Order = fieldNum.Add(1)
	var s string
	if err := value.Decode(&s); err == nil {
		m.Type = s
		return nil
	}
	var member struct {
		Type                  string    `yaml:"type"`
		Description           string    `yaml:"description"`
		To                    *MemberTo `yaml:"to"`
		Fields                *Fields   `yaml:"fields"`
		ReplacementCollection string    `yaml:"replacement_collection"`
		ReplacementEnum       []string  `yaml:"replacement_enum"`
		RestrictionMode       string    `yaml:"restriction_mode"`
		Required              bool      `yaml:"required"`
	}
	if err := value.Decode(&member); err != nil {
		return fmt.Errorf("member object without type: %w", err)
	}
	m.Type = member.Type
	m.Description = member.Description
	m.To = member.To
	m.Fields = member.Fields
	m.ReplacementCollection = member.ReplacementCollection
	m.ReplacementEnum = member.ReplacementEnum
	m.RestrictionMode = member.RestrictionMode
	m.Required = member.Required
	return nil
}

// UnmarshalYAML implements [gopkg.in/yaml.v3.Unmarshaler].
func (m *Collection) UnmarshalYAML(value *yaml.Node) error {
	m.Order = modelNum.Add(1)
	return value.Decode(&m.Fields)
}

// OrderedKeys returns the keys in document order.
func (ms Collections) OrderedKeys() []string {
	keys := make([]string, 0, len(ms))
	for k := range ms {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return ms[keys[i]].Order < ms[keys[j]].Order
	})
	return keys
}

// UnmarshalYAML implements [gopkg.in/yaml.v3.Unmarshaler].
func (fk *FilterKey) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	*fk = FilterKey{
		Order: filterNum.Add(1),
		Name:  s,
	}
	return nil
}

// UnmarshalYAML implements [gopkg.in/yaml.v3.Unmarshaler].
func (fs *Filters) UnmarshalYAML(value *yaml.Node) error {
	var fsm map[FilterKey][]string
	if err := value.Decode(&fsm); err != nil {
		return err
	}
	sorted := make([]FilterKey, 0, len(fsm))
	for k := range fsm {
		sorted = append(sorted, k)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Order < sorted[j].Order
	})

	*fs = make(Filters, 0, len(sorted))
	for _, s := range sorted {
		*fs = append(*fs, Filter{
			Name:  s.Name,
			Items: fsm[s],
		})
	}
	return nil
}

// RetainStrings returns a function which keeps string type fields in [Retain].
func RetainStrings(verbose bool) func(string, string, *Member) bool {
	return func(k, fk string, f *Member) bool {
		switch f.Type {
		case "string", "HTMLStrict", "text", "HTMLPermissive":
			return true
		default:
			if verbose {
				log.Printf("removing non-string %s.%s\n", k, fk)
			}
			return false
		}
	}
}

func copyStrings(s []string) []string {
	if s == nil {
		return nil
	}
	t := make([]string, len(s))
	copy(t, s)
	return t
}

// Clone returns a deep copy.
func (fs *Fields) Clone() *Fields {
	if fs == nil {
		return nil
	}
	return &Fields{
		Type: fs.Type,
		To:   fs.To,
	}
}

// Clone returns a deep copy.
func (mt *MemberTo) Clone() *MemberTo {
	if mt == nil {
		return nil
	}
	return &MemberTo{
		Collections: copyStrings(mt.Collections),
		Field:       mt.Field,
	}
}

// Clone returns a deep copy.
func (m *Member) Clone() *Member {
	return &Member{
		Type:                  m.Type,
		Description:           m.Description,
		To:                    m.To.Clone(),
		Fields:                m.Fields.Clone(),
		ReplacementCollection: m.ReplacementCollection,
		ReplacementEnum:       copyStrings(m.ReplacementEnum),
		RestrictionMode:       m.RestrictionMode,
		Required:              m.Required,
		Order:                 m.Order,
	}
}

// Clone returns a deep copy.
func (m *Collection) Clone() *Collection {
	var fields map[string]*Member
	if m.Fields != nil {
		fields = make(map[string]*Member)
		for k, v := range m.Fields {
			fields[k] = v.Clone()
		}
	}
	return &Collection{
		Fields: fields,
		Order:  m.Order,
	}
}

// Clone returns a deep copy.
func (ms Collections) Clone() Collections {
	cp := make(Collections, len(ms))
	for k, v := range ms {
		cp[k] = v.Clone()
	}
	return cp
}

// Retain removes members that are not marked to be kept by the keep function.
func (ms Collections) Retain(keep func(string, string, *Member) bool) {
	for k, m := range ms {
		for kf, f := range m.Fields {
			if !keep(k, kf, f) {
				delete(m.Fields, kf)
			}
		}
		if len(m.Fields) == 0 {
			// log.Printf("throw away collection '%s'.\n", k)
			delete(ms, k)
		}
	}
}

// OrderedKeys returns the keys in document order.
func (m *Collection) OrderedKeys() []string {
	fields := make([]string, 0, len(m.Fields))
	for f := range m.Fields {
		fields = append(fields, f)
	}
	sort.Slice(fields, func(i, j int) bool {
		return m.Fields[fields[i]].Order < m.Fields[fields[j]].Order
	})
	return fields
}

// AsFilters converts a collection into a filter.
func (ms Collections) AsFilters() Filters {
	keys := ms.OrderedKeys()
	fs := make(Filters, 0, len(keys))
	for _, k := range keys {
		items := ms[k].OrderedKeys()
		fs = append(fs, Filter{Name: k, Items: items})
	}
	return fs
}

func (fs Filters) Write(w io.Writer) error {
	b := bufio.NewWriter(w)
	for i := range fs {
		if i > 0 {
			fmt.Fprintln(b)
		}
		fmt.Fprintf(b, "%s:\n", fs[i].Name)
		for _, item := range fs[i].Items {
			fmt.Fprintf(b, "  - %s\n", item)
		}
	}
	return b.Flush()
}

// Retain returns a keep function for [Retain].
func (fs Filters) Retain(verbose bool) func(string, string, *Member) bool {
	type key struct {
		rel   string
		field string
	}
	keep := map[key]struct{}{}
	for _, m := range fs {
		for _, f := range m.Items {
			keep[key{rel: m.Name, field: f}] = struct{}{}
		}
	}
	return func(rk, fk string, _ *Member) bool {
		_, ok := keep[key{rel: rk, field: fk}]
		if !ok && verbose {
			log.Printf("removing filtered %s.%s\n", rk, fk)
		}
		return ok
	}
}
