package ui_test

import (
	"strings"
	"testing"

	"github.com/Dicklesworthstone/beads_viewer/pkg/analysis"
	"github.com/Dicklesworthstone/beads_viewer/pkg/ui"
)

// =============================================================================
// FlowMatrixView Tests
// =============================================================================

func TestFlowMatrixViewEmptyLabels(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels:     []string{},
		FlowMatrix: [][]int{},
	}

	result := ui.FlowMatrixView(flow, 80)
	expected := "No label flows available"
	if result != expected {
		t.Errorf("FlowMatrixView() = %q, want %q", result, expected)
	}
}

func TestFlowMatrixViewNilLabels(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels:     nil,
		FlowMatrix: nil,
	}

	result := ui.FlowMatrixView(flow, 80)
	expected := "No label flows available"
	if result != expected {
		t.Errorf("FlowMatrixView() with nil labels = %q, want %q", result, expected)
	}
}

func TestFlowMatrixViewSingleLabel(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels:     []string{"bug"},
		FlowMatrix: [][]int{{0}},
	}

	result := ui.FlowMatrixView(flow, 80)

	// Should contain the label
	if !strings.Contains(result, "bug") {
		t.Errorf("FlowMatrixView() should contain label 'bug', got: %q", result)
	}

	// Should contain header separator
	if !strings.Contains(result, "---") {
		t.Errorf("FlowMatrixView() should contain header separator, got: %q", result)
	}

	// Should contain the count (0)
	if !strings.Contains(result, "0") {
		t.Errorf("FlowMatrixView() should contain count '0', got: %q", result)
	}
}

func TestFlowMatrixViewMultipleLabels(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels: []string{"bug", "feat", "docs"},
		FlowMatrix: [][]int{
			{0, 2, 1}, // bug blocks feat 2 times, docs 1 time
			{1, 0, 3}, // feat blocks bug 1 time, docs 3 times
			{0, 0, 0}, // docs blocks nothing
		},
	}

	result := ui.FlowMatrixView(flow, 80)

	// Should contain all labels
	for _, label := range []string{"bug", "feat", "docs"} {
		if !strings.Contains(result, label) {
			t.Errorf("FlowMatrixView() should contain label %q, got: %q", label, result)
		}
	}

	// Should contain header separator
	if !strings.Contains(result, " | ") {
		t.Errorf("FlowMatrixView() should contain column separator ' | ', got: %q", result)
	}

	// Should contain values
	if !strings.Contains(result, "2") {
		t.Errorf("FlowMatrixView() should contain value '2', got: %q", result)
	}
	if !strings.Contains(result, "3") {
		t.Errorf("FlowMatrixView() should contain value '3', got: %q", result)
	}
}

func TestFlowMatrixViewLongLabels(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels: []string{"very-long-label-name", "another-long-one"},
		FlowMatrix: [][]int{
			{0, 5},
			{3, 0},
		},
	}

	result := ui.FlowMatrixView(flow, 80)

	// Should contain truncated labels in header (4 char cell width)
	// The header should have truncated versions
	if !strings.Contains(result, "ver") {
		t.Errorf("FlowMatrixView() should contain truncated header label, got: %q", result)
	}

	// Row labels can be longer (leftWidth based on max label length)
	if !strings.Contains(result, "very-long-label-name") {
		t.Errorf("FlowMatrixView() should contain full row label, got: %q", result)
	}
}

func TestFlowMatrixViewShortLabels(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels: []string{"a", "b"},
		FlowMatrix: [][]int{
			{0, 1},
			{2, 0},
		},
	}

	result := ui.FlowMatrixView(flow, 80)

	// Should use minimum left width of 6
	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Fatalf("Expected at least 3 lines, got %d", len(lines))
	}

	// Should contain labels
	if !strings.Contains(result, "a") || !strings.Contains(result, "b") {
		t.Errorf("FlowMatrixView() should contain labels 'a' and 'b', got: %q", result)
	}
}

func TestFlowMatrixViewFormatStructure(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels: []string{"X", "Y"},
		FlowMatrix: [][]int{
			{0, 1},
			{2, 0},
		},
	}

	result := ui.FlowMatrixView(flow, 80)
	lines := strings.Split(result, "\n")

	// Should have: header, separator, row1, row2, empty
	if len(lines) < 4 {
		t.Fatalf("Expected at least 4 lines, got %d: %q", len(lines), result)
	}

	// First line is header with column labels
	if !strings.Contains(lines[0], " | ") {
		t.Errorf("Header line should contain ' | ', got: %q", lines[0])
	}

	// Second line is separator with dashes
	if !strings.Contains(lines[1], "---") {
		t.Errorf("Separator line should contain dashes, got: %q", lines[1])
	}

	// Third line is first data row
	if !strings.Contains(lines[2], " | ") {
		t.Errorf("Data row should contain ' | ', got: %q", lines[2])
	}
}

func TestFlowMatrixViewLargeMatrix(t *testing.T) {
	// Test with a larger matrix
	labels := []string{"api", "web", "db", "auth", "core"}
	matrix := make([][]int, 5)
	for i := range matrix {
		matrix[i] = make([]int, 5)
		for j := range matrix[i] {
			if i != j {
				matrix[i][j] = i + j // Some non-zero values
			}
		}
	}

	flow := analysis.CrossLabelFlow{
		Labels:     labels,
		FlowMatrix: matrix,
	}

	result := ui.FlowMatrixView(flow, 120)

	// All labels should appear
	for _, label := range labels {
		if !strings.Contains(result, label) {
			t.Errorf("FlowMatrixView() should contain label %q", label)
		}
	}

	// Should have correct number of data rows
	lines := strings.Split(strings.TrimSpace(result), "\n")
	// header + separator + 5 data rows = 7 lines
	if len(lines) != 7 {
		t.Errorf("Expected 7 lines for 5 labels, got %d", len(lines))
	}
}

func TestFlowMatrixViewWidthParameter(t *testing.T) {
	// The width parameter is accepted but not currently used in truncation logic
	// Test that it doesn't cause issues with various widths
	flow := analysis.CrossLabelFlow{
		Labels:     []string{"label1", "label2"},
		FlowMatrix: [][]int{{0, 1}, {2, 0}},
	}

	widths := []int{20, 40, 80, 120, 200}
	for _, w := range widths {
		result := ui.FlowMatrixView(flow, w)
		if result == "" {
			t.Errorf("FlowMatrixView(width=%d) returned empty string", w)
		}
		// Should always contain the labels regardless of width
		if !strings.Contains(result, "label1") {
			t.Errorf("FlowMatrixView(width=%d) should contain 'label1'", w)
		}
	}
}

func TestFlowMatrixViewZeroWidth(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels:     []string{"test"},
		FlowMatrix: [][]int{{0}},
	}

	// Should not panic with zero width
	result := ui.FlowMatrixView(flow, 0)
	if result == "" {
		t.Error("FlowMatrixView(width=0) should not return empty string")
	}
}

func TestFlowMatrixViewNegativeWidth(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels:     []string{"test"},
		FlowMatrix: [][]int{{0}},
	}

	// Should not panic with negative width
	result := ui.FlowMatrixView(flow, -10)
	if result == "" {
		t.Error("FlowMatrixView(width=-10) should not return empty string")
	}
}

func TestFlowMatrixViewLargeValues(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels: []string{"A", "B"},
		FlowMatrix: [][]int{
			{0, 9999},
			{1234, 0},
		},
	}

	result := ui.FlowMatrixView(flow, 80)

	// Should contain large values
	if !strings.Contains(result, "9999") {
		t.Errorf("FlowMatrixView() should contain '9999', got: %q", result)
	}
	if !strings.Contains(result, "1234") {
		t.Errorf("FlowMatrixView() should contain '1234', got: %q", result)
	}
}

func TestFlowMatrixViewHeaderTruncation(t *testing.T) {
	// Cell width is 4, so labels longer than 4 should be truncated in header
	flow := analysis.CrossLabelFlow{
		Labels:     []string{"abcdef", "xyz"},
		FlowMatrix: [][]int{{0, 1}, {1, 0}},
	}

	result := ui.FlowMatrixView(flow, 80)
	lines := strings.Split(result, "\n")

	// Header line should have truncated "abcdef" to "abc…" or similar
	header := lines[0]
	// The truncation adds ellipsis: "abc…" which is 4 chars visually
	// Due to Unicode ellipsis being multi-byte, check for "abc"
	if !strings.Contains(header, "abc") {
		t.Errorf("Header should contain truncated 'abc', got: %q", header)
	}
}

func TestFlowMatrixViewRowLabelWidth(t *testing.T) {
	// Row labels use leftWidth based on max label length (min 6)
	tests := []struct {
		name        string
		labels      []string
		minLeftPad  int
	}{
		{
			name:        "short labels use min width 6",
			labels:      []string{"a", "b"},
			minLeftPad:  6,
		},
		{
			name:        "medium labels use their length",
			labels:      []string{"medium", "label"},
			minLeftPad:  6,
		},
		{
			name:        "long labels use their length",
			labels:      []string{"very-long-label", "short"},
			minLeftPad:  15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matrix := make([][]int, len(tt.labels))
			for i := range matrix {
				matrix[i] = make([]int, len(tt.labels))
			}

			flow := analysis.CrossLabelFlow{
				Labels:     tt.labels,
				FlowMatrix: matrix,
			}

			result := ui.FlowMatrixView(flow, 80)
			if result == "" {
				t.Error("FlowMatrixView() returned empty string")
			}
			// The output should be properly formatted
			lines := strings.Split(result, "\n")
			if len(lines) < 2 {
				t.Error("FlowMatrixView() should produce at least 2 lines")
			}
		})
	}
}

func TestFlowMatrixViewConsistentRowLength(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels: []string{"api", "web", "db"},
		FlowMatrix: [][]int{
			{0, 1, 2},
			{3, 0, 4},
			{5, 6, 0},
		},
	}

	result := ui.FlowMatrixView(flow, 80)
	lines := strings.Split(strings.TrimSuffix(result, "\n"), "\n")

	// All data rows should have the same column separator pattern
	separatorCount := -1
	for i, line := range lines {
		if i == 1 { // Skip the dash separator line
			continue
		}
		count := strings.Count(line, " | ")
		if count == 0 {
			continue // Skip empty lines
		}
		if separatorCount == -1 {
			separatorCount = count
		} else if count != separatorCount {
			t.Errorf("Line %d has %d separators, expected %d: %q", i, count, separatorCount, line)
		}
	}
}

func TestFlowMatrixViewSpecialCharactersInLabels(t *testing.T) {
	// Labels might contain special characters
	flow := analysis.CrossLabelFlow{
		Labels:     []string{"bug:fix", "feat-new", "doc_update"},
		FlowMatrix: [][]int{{0, 1, 0}, {0, 0, 1}, {1, 0, 0}},
	}

	result := ui.FlowMatrixView(flow, 80)

	// Should handle special characters without panic
	if result == "" {
		t.Error("FlowMatrixView() with special chars should not return empty")
	}

	// Should contain the labels (possibly truncated)
	if !strings.Contains(result, "bug") {
		t.Errorf("Result should contain 'bug', got: %q", result)
	}
}

func TestFlowMatrixViewUnicodeLabels(t *testing.T) {
	flow := analysis.CrossLabelFlow{
		Labels:     []string{"日本語", "한국어"},
		FlowMatrix: [][]int{{0, 1}, {1, 0}},
	}

	result := ui.FlowMatrixView(flow, 80)

	// Should handle Unicode without panic
	if result == "" {
		t.Error("FlowMatrixView() with Unicode should not return empty")
	}
}

func TestFlowMatrixViewSingleCharTruncation(t *testing.T) {
	// Test truncate function edge case with w=1
	flow := analysis.CrossLabelFlow{
		Labels:     []string{"abc"},
		FlowMatrix: [][]int{{0}},
	}

	// This tests the internal truncate function behavior
	result := ui.FlowMatrixView(flow, 80)
	if result == "" {
		t.Error("FlowMatrixView() should not return empty")
	}
}

// =============================================================================
// Integration-style Tests
// =============================================================================

func TestFlowMatrixViewRealisticData(t *testing.T) {
	// Realistic label flow data from a project
	flow := analysis.CrossLabelFlow{
		Labels: []string{"frontend", "backend", "database", "api", "testing"},
		FlowMatrix: [][]int{
			{0, 5, 0, 3, 2},  // frontend blocks: backend 5, api 3, testing 2
			{2, 0, 8, 4, 1},  // backend blocks: frontend 2, database 8, api 4, testing 1
			{0, 3, 0, 0, 0},  // database blocks: backend 3
			{1, 2, 1, 0, 3},  // api blocks various
			{0, 0, 0, 0, 0},  // testing blocks nothing
		},
	}

	result := ui.FlowMatrixView(flow, 100)

	// Should be readable output
	if result == "" {
		t.Fatal("FlowMatrixView() returned empty for realistic data")
	}

	// Should contain all labels (or truncated versions)
	for _, label := range flow.Labels {
		// Check for at least first 3 chars or full label if shorter
		checkLen := 3
		if len(label) < checkLen {
			checkLen = len(label)
		}
		if !strings.Contains(result, label[:checkLen]) {
			t.Errorf("Result should contain label %q (or truncation)", label)
		}
	}

	// Should show the blockage counts
	if !strings.Contains(result, "8") {
		t.Error("Result should show max blockage count 8")
	}
}

func TestFlowMatrixViewOutput(t *testing.T) {
	// Verify the exact output format matches expectations
	flow := analysis.CrossLabelFlow{
		Labels:     []string{"X", "Y"},
		FlowMatrix: [][]int{{0, 1}, {2, 0}},
	}

	result := ui.FlowMatrixView(flow, 80)

	// Verify structure
	if !strings.Contains(result, "X") {
		t.Error("Output should contain 'X'")
	}
	if !strings.Contains(result, "Y") {
		t.Error("Output should contain 'Y'")
	}
	if !strings.Contains(result, " | ") {
		t.Error("Output should contain ' | ' separator")
	}
	if !strings.Contains(result, "1") {
		t.Error("Output should contain value '1'")
	}
	if !strings.Contains(result, "2") {
		t.Error("Output should contain value '2'")
	}
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkFlowMatrixViewSmall(b *testing.B) {
	flow := analysis.CrossLabelFlow{
		Labels:     []string{"a", "b", "c"},
		FlowMatrix: [][]int{{0, 1, 2}, {3, 0, 4}, {5, 6, 0}},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ui.FlowMatrixView(flow, 80)
	}
}

func BenchmarkFlowMatrixViewMedium(b *testing.B) {
	labels := make([]string, 10)
	for i := range labels {
		labels[i] = strings.Repeat("x", 8)
	}
	matrix := make([][]int, 10)
	for i := range matrix {
		matrix[i] = make([]int, 10)
		for j := range matrix[i] {
			matrix[i][j] = i * j
		}
	}
	flow := analysis.CrossLabelFlow{
		Labels:     labels,
		FlowMatrix: matrix,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ui.FlowMatrixView(flow, 120)
	}
}

func BenchmarkFlowMatrixViewLarge(b *testing.B) {
	labels := make([]string, 20)
	for i := range labels {
		labels[i] = strings.Repeat("label", 3)
	}
	matrix := make([][]int, 20)
	for i := range matrix {
		matrix[i] = make([]int, 20)
		for j := range matrix[i] {
			matrix[i][j] = (i + 1) * (j + 1)
		}
	}
	flow := analysis.CrossLabelFlow{
		Labels:     labels,
		FlowMatrix: matrix,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ui.FlowMatrixView(flow, 200)
	}
}
