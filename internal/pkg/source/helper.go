package source

import (
	"unicode"
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
}

// Implement the sort.Interface for LogEntrySlice
func (s LogEntrySlice) Len() int { return len(s.Entries) }
func (s LogEntrySlice) Less(i, j int) bool {
	switch s.SortByField {
	case "fieldName1":
		return s.Entries[i].Fields[0] < s.Entries[j].Fields[0]
	case "fieldName2":
		return s.Entries[i].Fields[1] < s.Entries[j].Fields[1]
	default:
		return s.Entries[i].Fields[1] < s.Entries[j].Fields[1]
	}
}
func (s LogEntrySlice) Swap(i, j int) { s.Entries[i], s.Entries[j] = s.Entries[j], s.Entries[i] }
