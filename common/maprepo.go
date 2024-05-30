package common

import (
	"fmt"
)

type MapRepo struct {
	entries map[string][]byte
}

func NewMapRepo() *MapRepo {
	return &MapRepo{entries: make(map[string][]byte)}
}

func (m *MapRepo) DoesEntryExist(entryId string) (bool, error) {
	_, exists := m.entries[entryId]
	return exists, nil
}

func (m *MapRepo) GetEntry(entryId string) ([]byte, error) {
	data, exists := m.entries[entryId]
	if !exists {
		return nil, fmt.Errorf("entry does not exist")
	}
	return data, nil
}

func (m *MapRepo) CreateEntry(entryId string, data []byte) error {
	_, exists := m.entries[entryId]
	if exists {
		return fmt.Errorf("entry already exists")
	}
	m.entries[entryId] = data
	return nil
}

func (m *MapRepo) UpdateEntry(entryId string, data []byte) error {
	_, exists := m.entries[entryId]
	if !exists {
		return fmt.Errorf("entry does not exist")
	}
	m.entries[entryId] = data
	return nil
}

func (m *MapRepo) DeleteEntry(entryId string) error {
	_, exists := m.entries[entryId]
	if !exists {
		return fmt.Errorf("entry does not exist")
	}
	delete(m.entries, entryId)
	return nil
}

func (m *MapRepo) ListEntryIds() ([]string, error) {
	ids := make([]string, 0)
	for id := range m.entries {
		ids = append(ids, id)
	}
	return ids, nil
}
