package models

type DashboardStats struct {
	TotalClients    int `json:"total_clients"`
	ActivePrograms  int `json:"active_programs"`
	UpcomingSession int `json:"upcoming_sessions"`
}
