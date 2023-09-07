package payload

type GradeMessage struct {
	Id       string                `json:"id"`
	Type     string                `json:"language"`
	Settings *GradeSettingsMessage `json:"settings"`
	Files    []GradeFileMessage    `json:"files"`
}

type GradeSettingsMessage struct {
	LimitMemory int `json:"softLimitMemory"`
	LimitTime   int `json:"softLimitTime"`
}

type GradeFileMessage struct {
	Name       string `json:"name"`
	SourceType string `json:"sourceType"`
	Source     string `json:"source"`
}

type GradeResultMessage struct {
	Id       string                     `json:"id"`
	Status   string                     `json:"status"`
	Metadata GradeResultMetadataMessage `json:"metadata"`
}

type GradeResultMetadataMessage struct {
	MemoryUsage    int    `json:"memory"`
	TimeUsage      int    `json:"time"`
	CompilationLog string `json:"compilationLog"`
}
