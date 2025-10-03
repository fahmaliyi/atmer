package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Generic storage type that manages JSON file reads/writes
type Storage[T any] struct {
	filePath string
	mu       sync.Mutex
}

// New creates a new storage bound to a file path
func New[T any](filePath string) *Storage[T] {
	return &Storage[T]{filePath: filePath}
}

// Load all records from file
func (s *Storage[T]) Load() ([]T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []T{}, err // empty if file doesn’t exist yet
		}

		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var records []T
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return records, nil
}

// Save all records to file
func (s *Storage[T]) Save(records []T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Add a new record
func (s *Storage[T]) Add(record T) error {
	records, err := s.Load()
	if err != nil {
		return err
	}
	records = append(records, record)
	return s.Save(records)
}

// Update record(s) by matcher function
func (s *Storage[T]) Update(match func(T) bool, updater func(*T)) error {
	records, err := s.Load()
	if err != nil {
		return err
	}

	updated := false
	for i := range records {
		if match(records[i]) {
			updater(&records[i])
			updated = true
		}
	}

	if !updated {
		return fmt.Errorf("no matching record found")
	}

	return s.Save(records)
}

// Delete record(s) by matcher function
func (s *Storage[T]) Delete(match func(T) bool) error {
	records, err := s.Load()
	if err != nil {
		return err
	}

	newRecords := make([]T, 0, len(records))
	for _, r := range records {
		if !match(r) {
			newRecords = append(newRecords, r)
		}
	}

	return s.Save(newRecords)
}
