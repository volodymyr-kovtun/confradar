// Package differ provides .env file comparison functionality.
package differ

// DiffResult holds the comparison between two .env files.
type DiffResult struct {
	LeftPath  string     `json:"left_path"`
	RightPath string     `json:"right_path"`
	OnlyLeft  []DiffKey  `json:"only_left"`
	OnlyRight []DiffKey  `json:"only_right"`
	Changed   []DiffPair `json:"changed"`
	Common    []DiffPair `json:"common"`
}

// DiffKey is a key that exists in only one file.
type DiffKey struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// DiffPair is a key that exists in both files, possibly with different values.
type DiffPair struct {
	Key       string `json:"key"`
	LeftValue string `json:"left_value"`
	RightValue string `json:"right_value"`
}

// Stats returns a summary of the diff.
func (d *DiffResult) Stats() (onlyLeft, onlyRight, changed, common int) {
	return len(d.OnlyLeft), len(d.OnlyRight), len(d.Changed), len(d.Common)
}

// HasDifferences returns true if there are any differences.
func (d *DiffResult) HasDifferences() bool {
	return len(d.OnlyLeft) > 0 || len(d.OnlyRight) > 0 || len(d.Changed) > 0
}
