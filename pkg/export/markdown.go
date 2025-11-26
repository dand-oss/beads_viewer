package export

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"beads_viewer/pkg/model"
)

// GenerateMarkdown creates a comprehensive markdown report of all issues
func GenerateMarkdown(issues []model.Issue, title string) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC1123)))

	// Table of Contents / Summary
	sb.WriteString("## Summary\n\n")
	
	open := 0
	blocked := 0
	closed := 0
	
	for _, i := range issues {
		switch i.Status {
		case model.StatusOpen: open++
		case model.StatusInProgress: open++
		case model.StatusBlocked: blocked++
		case model.StatusClosed: closed++
		}
	}
		sb.WriteString(fmt.Sprintf("- **Total**: %d\n", len(issues)))
	sb.WriteString(fmt.Sprintf("- **Open**: %d\n", open))
	sb.WriteString(fmt.Sprintf("- **Blocked**: %d\n", blocked))
	sb.WriteString(fmt.Sprintf("- **Closed**: %d\n\n", closed))

	sb.WriteString("## Table of Contents\n\n")
	for _, i := range issues {
		link := fmt.Sprintf("#%s", strings.ToLower(i.ID)) // This is heuristic, markdown anchors vary by renderer
		sb.WriteString(fmt.Sprintf("- [%s %s](%s) (%s)\n", i.ID, i.Title, link, i.Status))
	}
	sb.WriteString("\n---\n\n")

	// Issues
	for _, i := range issues {
		sb.WriteString(fmt.Sprintf("## %s %s\n\n", i.ID, i.Title))
		
		// Metadata Table
		sb.WriteString("| Type | Priority | Status | Assignee | Created |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		sb.WriteString(fmt.Sprintf("| %s | %d | %s | %s | %s |\n\n", 
			i.IssueType, i.Priority, i.Status, i.Assignee, i.CreatedAt.Format("2006-01-02")))

		if i.Description != "" {
			sb.WriteString("### Description\n\n")
			sb.WriteString(i.Description + "\n\n")
		}

		if i.AcceptanceCriteria != "" {
			sb.WriteString("### Acceptance Criteria\n\n")
			sb.WriteString(i.AcceptanceCriteria + "\n\n")
		}
		
		if len(i.Dependencies) > 0 {
			sb.WriteString("### Dependencies\n\n")
			for _, dep := range i.Dependencies {
				sb.WriteString(fmt.Sprintf("- **%s**: %s\n", dep.Type, dep.DependsOnID))
			}
			sb.WriteString("\n")
		}

		if len(i.Comments) > 0 {
			sb.WriteString("### Comments\n\n")
			for _, c := range i.Comments {
				sb.WriteString(fmt.Sprintf("> **%s** (%s):\n> %s\n\n", c.Author, c.CreatedAt.Format("2006-01-02"), strings.ReplaceAll(c.Text, "\n", "\n> ")))
			}
		}
		
		sb.WriteString("---\n\n")
	}

	return sb.String(), nil
}

// SaveMarkdownToFile writes the generated markdown to a file
func SaveMarkdownToFile(issues []model.Issue, filename string) error {
	// Sort issues for the report
	sort.Slice(issues, func(i, j int) bool {
		// Sort logic: Open first, then priority, then date
		// (Same as UI)
		iClosed := issues[i].Status == model.StatusClosed
		jClosed := issues[j].Status == model.StatusClosed
		if iClosed != jClosed {
			return !iClosed
		}
		if issues[i].Priority != issues[j].Priority {
			return issues[i].Priority < issues[j].Priority
		}
		return issues[i].CreatedAt.After(issues[j].CreatedAt)
	})

	content, err := GenerateMarkdown(issues, "Beads Export")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(content), 0644)
}
