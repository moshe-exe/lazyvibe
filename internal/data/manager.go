package data

import (
	"sync"
	"time"
)

// Cache TTL constants
const (
	VMTTL       = 2 * time.Second
	SessionsTTL = 10 * time.Second
	StatsTTL    = 30 * time.Second
)

// cacheEntry represents a cached data item with timestamp.
type cacheEntry[T any] struct {
	data      T
	timestamp time.Time
}

func (c *cacheEntry[T]) isValid(ttl time.Duration) bool {
	return time.Since(c.timestamp) < ttl
}

// Manager manages data fetching with caching.
type Manager struct {
	mu sync.RWMutex

	vmCache       *cacheEntry[VMStatus]
	sessionsCache *cacheEntry[[]SessionEntry]
	statsCache    *cacheEntry[[]DailyActivity]
	projectsCache *cacheEntry[[]ProjectSummary]
}

// NewManager creates a new data manager.
func NewManager() *Manager {
	return &Manager{}
}

// GetVMStatus returns the VM status with caching.
func (m *Manager) GetVMStatus(forceRefresh bool) VMStatus {
	m.mu.RLock()
	if !forceRefresh && m.vmCache != nil && m.vmCache.isValid(VMTTL) {
		status := m.vmCache.data
		m.mu.RUnlock()
		return status
	}
	m.mu.RUnlock()

	status := GetVMStatus()

	m.mu.Lock()
	m.vmCache = &cacheEntry[VMStatus]{data: status, timestamp: time.Now()}
	m.mu.Unlock()

	return status
}

// GetSessions returns sessions with caching.
func (m *Manager) GetSessions(forceRefresh bool) []SessionEntry {
	m.mu.RLock()
	if !forceRefresh && m.sessionsCache != nil && m.sessionsCache.isValid(SessionsTTL) {
		sessions := m.sessionsCache.data
		m.mu.RUnlock()
		return sessions
	}
	m.mu.RUnlock()

	sessions := ParseSessions()

	m.mu.Lock()
	m.sessionsCache = &cacheEntry[[]SessionEntry]{data: sessions, timestamp: time.Now()}
	m.projectsCache = nil // Invalidate projects cache when sessions change
	m.mu.Unlock()

	return sessions
}

// GetDailyActivity returns daily activity with caching.
func (m *Manager) GetDailyActivity(forceRefresh bool) []DailyActivity {
	m.mu.RLock()
	if !forceRefresh && m.statsCache != nil && m.statsCache.isValid(StatsTTL) {
		activity := m.statsCache.data
		m.mu.RUnlock()
		return activity
	}
	m.mu.RUnlock()

	activity := ParseStatsCache()

	m.mu.Lock()
	m.statsCache = &cacheEntry[[]DailyActivity]{data: activity, timestamp: time.Now()}
	m.mu.Unlock()

	return activity
}

// GetProjects returns project summaries with caching.
func (m *Manager) GetProjects(forceRefresh bool) []ProjectSummary {
	m.mu.RLock()
	if !forceRefresh && m.projectsCache != nil && m.projectsCache.isValid(SessionsTTL) {
		projects := m.projectsCache.data
		m.mu.RUnlock()
		return projects
	}
	m.mu.RUnlock()

	sessions := m.GetSessions(forceRefresh)
	projects := AggregateProjects(sessions)

	m.mu.Lock()
	m.projectsCache = &cacheEntry[[]ProjectSummary]{data: projects, timestamp: time.Now()}
	m.mu.Unlock()

	return projects
}

// GetDashboardData returns all dashboard data.
func (m *Manager) GetDashboardData(forceRefresh bool) DashboardData {
	return DashboardData{
		VMStatus:      m.GetVMStatus(forceRefresh),
		Sessions:      m.GetSessions(forceRefresh),
		DailyActivity: m.GetDailyActivity(forceRefresh),
		Projects:      m.GetProjects(forceRefresh),
	}
}

// RefreshAll forces a refresh of all data.
func (m *Manager) RefreshAll() DashboardData {
	return m.GetDashboardData(true)
}
