package models

// AnalysisResult — раздел анализа тендера (пишет микросервис analizator_zakupok).
// Файлы: result/{reg}/analysis/analysis.json и valid_info/analysis.json.
type AnalysisResult struct {
	RegNumber      string         `json:"reg_number"`
	Law            string         `json:"law,omitempty"`
	Status         string         `json:"status"` // pending|running|completed|failed
	ChecklistID    string         `json:"checklist_id,omitempty"`
	ChecklistName  string         `json:"checklist_name,omitempty"`
	Model          string         `json:"model,omitempty"`
	Recommendation string         `json:"recommendation,omitempty"` // participate|caution|skip|unknown
	Score          float64        `json:"score"`
	Summary        string         `json:"summary,omitempty"`
	Items          []AnalysisItem `json:"items,omitempty"`
	Risks          []string       `json:"risks,omitempty"`
	Actions        []string       `json:"actions,omitempty"`
	SourcesUsed    []string       `json:"sources_used,omitempty"`
	ChunksTotal    int            `json:"chunks_total,omitempty"`
	ChunksUsed     int            `json:"chunks_used,omitempty"`
	Error          string         `json:"error,omitempty"`
	StartedAt      string         `json:"started_at,omitempty"`
	AnalyzedAt     string         `json:"analyzed_at,omitempty"`
}

// AnalysisItem — результат пункта чек-листа.
type AnalysisItem struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Status   string   `json:"status"` // ok|warn|fail|unknown
	Score    float64  `json:"score"`
	Findings string   `json:"findings"`
	Evidence []string `json:"evidence,omitempty"`
	ChunkIDs []string `json:"chunk_ids,omitempty"`
}
