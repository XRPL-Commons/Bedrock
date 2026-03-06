package tester

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// LoadSnapshot reads a gas snapshot from a file
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to read snapshot: %w", err)
	}

	snapshot := &Snapshot{
		Entries: make(map[string]int64),
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		value, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			continue
		}

		snapshot.Entries[name] = value
	}

	return snapshot, nil
}

// SaveSnapshot writes a gas snapshot to a file
func SaveSnapshot(path string, entries []GasEntry) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Gas Snapshot - %s\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("# Function: Gas Used\n\n")

	// Sort entries by function name for deterministic output
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Function < entries[j].Function
	})

	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%s: %d\n", e.Function, e.GasUsed))
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	return nil
}

// DiffSnapshots compares two snapshots and returns the differences
func DiffSnapshots(old, new *Snapshot) *SnapshotDiff {
	diff := &SnapshotDiff{
		Added:   make(map[string]int64),
		Removed: make(map[string]int64),
		Changed: make(map[string][2]int64),
	}

	// Check for changed and removed entries
	for name, oldVal := range old.Entries {
		if newVal, ok := new.Entries[name]; ok {
			if oldVal != newVal {
				diff.Changed[name] = [2]int64{oldVal, newVal}
			}
		} else {
			diff.Removed[name] = oldVal
		}
	}

	// Check for added entries
	for name, newVal := range new.Entries {
		if _, ok := old.Entries[name]; !ok {
			diff.Added[name] = newVal
		}
	}

	return diff
}

// HasChanges returns true if the diff contains any changes
func (d *SnapshotDiff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// FormatDiff returns a human-readable representation of the diff
func (d *SnapshotDiff) FormatDiff() string {
	if !d.HasChanges() {
		return "No changes"
	}

	var sb strings.Builder

	if len(d.Added) > 0 {
		sb.WriteString("Added:\n")
		for name, val := range d.Added {
			sb.WriteString(fmt.Sprintf("  + %s: %d\n", name, val))
		}
	}

	if len(d.Removed) > 0 {
		sb.WriteString("Removed:\n")
		for name, val := range d.Removed {
			sb.WriteString(fmt.Sprintf("  - %s: %d\n", name, val))
		}
	}

	if len(d.Changed) > 0 {
		sb.WriteString("Changed:\n")
		for name, vals := range d.Changed {
			delta := vals[1] - vals[0]
			sign := "+"
			if delta < 0 {
				sign = ""
			}
			sb.WriteString(fmt.Sprintf("  ~ %s: %d -> %d (%s%d)\n", name, vals[0], vals[1], sign, delta))
		}
	}

	return sb.String()
}

// SnapshotFromEntries creates a Snapshot from gas entries
func SnapshotFromEntries(entries []GasEntry) *Snapshot {
	snapshot := &Snapshot{
		Entries:   make(map[string]int64),
		CreatedAt: time.Now(),
	}
	for _, e := range entries {
		snapshot.Entries[e.Function] = e.GasUsed
	}
	return snapshot
}
