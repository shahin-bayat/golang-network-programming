package records

import (
	"fmt"
)

type RecordStore interface {
	Get(name string) ([4]byte, error)
}

type InMemoryRecordStore struct {
	Records map[string][4]byte
}

func (m *InMemoryRecordStore) Get(name string) ([4]byte, error) {
	ip, ok := m.Records[name]
	if !ok {
		return ip, fmt.Errorf("record %s not found in memory", name)
	}
	return ip, nil
}
