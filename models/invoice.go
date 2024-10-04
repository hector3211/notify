package models

import "time"

type JobStatus string

const (
	JOBPENDING     JobStatus = "pending"
	JOBCUTTING     JobStatus = "cutting"
	JOBFABRICATING JobStatus = "fabricating"
	JOBDONE        JobStatus = "done"
)

func (j JobStatus) String() string {
	switch j {
	case JOBCUTTING:
		return "cutting"
	case JOBFABRICATING:
		return "fabricating"
	case JOBDONE:
		return "done"
	default:
		return "pending"
	}
}

func StatusArray(currStatus JobStatus) []JobStatus {
	var filteredStatuses []JobStatus
	statuses := []JobStatus{
		JOBPENDING,
		JOBFABRICATING,
		JOBCUTTING,
		JOBDONE,
	}

	for _, status := range statuses {
		if status != currStatus {
			filteredStatuses = append(filteredStatuses, status)
		}
	}

	return filteredStatuses
}

func MatchJobStatus(s string) JobStatus {
	switch s {
	case "cutting":
		return JOBCUTTING
	case "fabricating":
		return JOBFABRICATING
	case "done":
		return JOBDONE
	default:
		return JOBPENDING
	}
}

type Invoice struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Invoice   string    `json:"invoice"`
	Status    JobStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
