package analysis_test

import (
	"testing"

	"beads_viewer/pkg/analysis"
	"beads_viewer/pkg/model"
)

func TestImpactScore(t *testing.T) {
	// Chain: A -> B -> C (A depends on B, B depends on C)
	// Edges: A->B, B->C
	// In Graph: A->B, B->C (u -> v)
	// Impact Depth logic:
	// C (Leaf): Should have Impact 1 (It's a root dependency).
	// B: Impact 1 + Impact(C) = 2.
	// A: Impact 1 + Impact(B) = 3.
	// Wait.
	// If A->B->C.
	// B is "upstream" of A?
	// If B breaks, A breaks.
	// If C breaks, B breaks, A breaks.
	// So C has highest impact (3).
	// A has lowest impact (1).
	
	// Let's verify my implementation.
	// Forward iteration of Topo Sort.
	// A->B->C.
	// Sort: A, B, C.
	// Loop:
	// A: To(A)? None. Impact = 1.
	// B: To(B)? A. Impact = 1 + 1 = 2.
	// C: To(C)? B. Impact = 1 + 2 = 3.
	// Correct. C has score 3.
	
	issues := []model.Issue{
		{ID: "A", Dependencies: []*model.Dependency{{DependsOnID: "B"}}},
		{ID: "B", Dependencies: []*model.Dependency{{DependsOnID: "C"}}},
		{ID: "C"},
	}
	
	an := analysis.NewAnalyzer(issues)
	stats := an.Analyze()
	
	if stats.CriticalPathScore["C"] != 3 {
		t.Errorf("Expected C to have score 3, got %f", stats.CriticalPathScore["C"])
	}
	if stats.CriticalPathScore["B"] != 2 {
		t.Errorf("Expected B to have score 2, got %f", stats.CriticalPathScore["B"])
	}
	if stats.CriticalPathScore["A"] != 1 {
		t.Errorf("Expected A to have score 1, got %f", stats.CriticalPathScore["A"])
	}
}
