package payload

type GradeRequestMessage struct {
	Language  string               `json:"language"`
	SourceUrl string               `json:"sourceUrl"`
	Test      []GradeTestMessage   `json:"test"`
	Metadata  GradeMetadataMessage `json:"metadata"`
}

type GradeTestMessage struct {
	InputUrl  string `json:"input"`
	OutputUrl string `json:"output"`
}

type GradeMetadataMessage struct {
	SubmissionId int   `json:"submissionId"`
	TestcaseIds  []int `json:"testcaseIds"`
}

type GradeResponseMessage struct {
	CompileOutput string                       `json:"compileOutput"`
	Results       []GradeResponseResultMessage `json:"results"`
	Metadata      GradeMetadataMessage         `json:"metadata"`
}

type GradeResponseResultMessage struct {
	Hash string `json:"hash"`
	Pass bool   `json:"pass"`
	Time int    `json:"time"`
}
