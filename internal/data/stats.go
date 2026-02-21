package data

import (
	"encoding/json"
	"os"
	"sort"
)

const statsCachePath = "~/.claude/stats-cache.json"

// statsCacheFile represents the structure of stats-cache.json
type statsCacheFile struct {
	DailyActivity []dailyActivityJSON `json:"dailyActivity"`
}

// dailyActivityJSON represents the JSON structure of daily activity
type dailyActivityJSON struct {
	Date          string `json:"date"`
	MessageCount  int    `json:"messageCount"`
	SessionCount  int    `json:"sessionCount"`
	ToolCallCount int    `json:"toolCallCount"`
}

// ParseStatsCache parses stats-cache.json and returns daily activity data.
func ParseStatsCache() []DailyActivity {
	fpath := expandPath(statsCachePath)

	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil
	}

	var cacheFile statsCacheFile
	if err := json.Unmarshal(data, &cacheFile); err != nil {
		return nil
	}

	result := make([]DailyActivity, 0, len(cacheFile.DailyActivity))
	for _, day := range cacheFile.DailyActivity {
		// Estimate tokens: rough heuristic based on messages and tool calls
		// Average message ~500 chars, tool call ~200 chars, 4 chars per token
		estimatedTokens := (day.MessageCount*500 + day.ToolCallCount*200) / 4
		activity := DailyActivity{
			Date:          day.Date,
			MessageCount:  day.MessageCount,
			SessionCount:  day.SessionCount,
			ToolCallCount: day.ToolCallCount,
			TokenCount:    estimatedTokens,
		}
		result = append(result, activity)
	}

	// Sort by date (most recent last for sparkline display)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return result
}
