package model

type LogData struct {
	ID     int64  `meddler:"log_id,pk"`
	ProcID int64  `meddler:"log_job_id"`
	Data   []byte `meddler:"log_data"`
}
