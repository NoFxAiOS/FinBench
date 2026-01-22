package benchmark

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"FinBench/market"
)

// CaptureSnapshot captures real-time market data and creates a snapshot
func CaptureSnapshot(symbol, interval string, klineCount int) (*Snapshot, error) {
	klines, err := market.GetKlines(symbol, interval, klineCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}

	if len(klines) < klineCount {
		return nil, fmt.Errorf("insufficient klines: got %d, need %d", len(klines), klineCount)
	}

	now := time.Now()
	snapshot := &Snapshot{
		ID:        fmt.Sprintf("%s_%s_%s", now.Format("20060102_150405"), symbol, interval),
		Symbol:    symbol,
		Interval:  interval,
		Timestamp: now.UnixMilli(),
		Klines:    klines,
	}

	return snapshot, nil
}

// SaveSnapshot saves a snapshot to the specified directory
func SaveSnapshot(snapshot *Snapshot, dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filename := fmt.Sprintf("%s.json", snapshot.ID)
	filepath := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	return nil
}

// LoadSnapshot loads a snapshot from a file
func LoadSnapshot(filepath string) (*Snapshot, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot: %w", err)
	}

	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return &snapshot, nil
}

// LoadSnapshots loads all snapshots from a directory
func LoadSnapshots(dir string) ([]*Snapshot, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var snapshots []*Snapshot
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filepath := filepath.Join(dir, entry.Name())
		snapshot, err := LoadSnapshot(filepath)
		if err != nil {
			continue // Skip invalid files
		}
		snapshots = append(snapshots, snapshot)
	}

	// Sort by timestamp (newest first)
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Timestamp > snapshots[j].Timestamp
	})

	return snapshots, nil
}

// SnapshotIndex holds metadata about available snapshots
type SnapshotIndex struct {
	UpdatedAt time.Time           `json:"updated_at"`
	Snapshots []SnapshotMetadata  `json:"snapshots"`
}

// SnapshotMetadata holds metadata for a single snapshot
type SnapshotMetadata struct {
	ID        string `json:"id"`
	Symbol    string `json:"symbol"`
	Interval  string `json:"interval"`
	Timestamp int64  `json:"timestamp"`
	Filepath  string `json:"filepath"`
}

// UpdateIndex updates the snapshot index file
func UpdateIndex(dir string) error {
	snapshots, err := LoadSnapshots(dir)
	if err != nil {
		return err
	}

	index := &SnapshotIndex{
		UpdatedAt: time.Now(),
		Snapshots: make([]SnapshotMetadata, len(snapshots)),
	}

	for i, s := range snapshots {
		index.Snapshots[i] = SnapshotMetadata{
			ID:        s.ID,
			Symbol:    s.Symbol,
			Interval:  s.Interval,
			Timestamp: s.Timestamp,
			Filepath:  fmt.Sprintf("%s.json", s.ID),
		}
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "index.json"), data, 0644)
}
