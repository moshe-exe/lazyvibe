package data

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const sessionsIndexPattern = "~/.claude/projects/*/sessions-index.json"

// sessionsIndexFile represents the structure of sessions-index.json
type sessionsIndexFile struct {
	Entries []sessionEntryJSON `json:"entries"`
}

// sessionEntryJSON represents the JSON structure of a session entry
type sessionEntryJSON struct {
	SessionID    string  `json:"sessionId"`
	ProjectPath  string  `json:"projectPath"`
	Summary      string  `json:"summary"`
	FirstPrompt  string  `json:"firstPrompt"`
	MessageCount int     `json:"messageCount"`
	Created      string  `json:"created"`
	Modified     string  `json:"modified"`
	GitBranch    *string `json:"gitBranch"`
}

// ParseSessions parses all sessions-index.json files and returns session entries.
func ParseSessions() []SessionEntry {
	pattern := expandPath(sessionsIndexPattern)
	indexFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}

	var sessions []SessionEntry

	for _, fpath := range indexFiles {
		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}

		// Try to parse as struct with entries field first
		var indexFile sessionsIndexFile
		if err := json.Unmarshal(data, &indexFile); err != nil {
			// Try parsing as array directly
			var entries []sessionEntryJSON
			if err := json.Unmarshal(data, &entries); err != nil {
				continue
			}
			indexFile.Entries = entries
		}

		for _, entry := range indexFile.Entries {
			summary := entry.Summary
			if summary == "" {
				summary = entry.FirstPrompt
			}
			if len(summary) > 100 {
				summary = summary[:100]
			}

			projectName := filepath.Base(entry.ProjectPath)
			if projectName == "" || projectName == "." {
				projectName = "Unknown"
			}

			session := SessionEntry{
				SessionID:    entry.SessionID,
				ProjectPath:  entry.ProjectPath,
				ProjectName:  projectName,
				Summary:      summary,
				MessageCount: entry.MessageCount,
				Created:      parseTimestamp(entry.Created),
				Modified:     parseTimestamp(entry.Modified),
				GitBranch:    entry.GitBranch,
			}
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// AggregateProjects aggregates sessions into project summaries.
func AggregateProjects(sessions []SessionEntry) []ProjectSummary {
	projects := make(map[string]*ProjectSummary)

	for _, session := range sessions {
		key := session.ProjectPath
		if existing, ok := projects[key]; ok {
			existing.SessionCount++
			existing.TotalMessages += session.MessageCount
			if session.Modified.After(existing.LastActivity) {
				existing.LastActivity = session.Modified
			}
		} else {
			projects[key] = &ProjectSummary{
				ProjectName:   session.ProjectName,
				ProjectPath:   session.ProjectPath,
				SessionCount:  1,
				TotalMessages: session.MessageCount,
				LastActivity:  session.Modified,
			}
		}
	}

	// Convert to slice and sort by last activity
	result := make([]ProjectSummary, 0, len(projects))
	for _, p := range projects {
		result = append(result, *p)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].LastActivity.After(result[j].LastActivity)
	})

	return result
}

// parseTimestamp parses an ISO timestamp string to time.Time.
func parseTimestamp(ts string) time.Time {
	if ts == "" {
		return time.Now()
	}

	// Handle ISO format with Z suffix
	if strings.HasSuffix(ts, "Z") {
		ts = ts[:len(ts)-1] + "+00:00"
	}

	// Try various formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, ts); err == nil {
			return t
		}
	}

	return time.Now()
}

// expandPath expands ~ to the user's home directory.
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
