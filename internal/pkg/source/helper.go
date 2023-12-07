package source

import (
	"fmt"
	"unicode"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

func normalizeJSON(input []byte) []byte {
	out := make([]byte, 0, len(input))

	for _, r := range string(input) {
		if unicode.IsPrint(r) {
			out = append(out, []byte(string(r))...)
		}
	}

	return out
}

// Custom LogEntrySlice to implement sort.Interface to support sort
type LogEntrySlice struct {
	Entries     LogEntries
	SortByField string
	cfg         *config.Config
}

// Implement the sort.Interface for LogEntrySlice
func (s LogEntrySlice) Len() int { return len(s.Entries) }
func (s LogEntrySlice) Less(i, j int) bool {
	for fieldIndex, field := range s.cfg.Fields {
		if field.Title == s.SortByField {
			return s.Entries[i].Fields[fieldIndex] < s.Entries[j].Fields[fieldIndex]
		}
	}
	//default
	return s.Entries[i].Fields[0] < s.Entries[j].Fields[0]
}
func (s LogEntrySlice) Swap(i, j int) { s.Entries[i], s.Entries[j] = s.Entries[j], s.Entries[i] }

func getFieldFromConfigByIndex(index int, cfg *config.Config) (config.Field, error) {
	for i, field := range cfg.Fields {
		if i == index {
			return field, nil
		}
	}
	return config.Field{}, fmt.Errorf("no field found with index %d", index)
}
