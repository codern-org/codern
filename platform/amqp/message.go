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

// Accept only string
type GradeMetadataMessage struct {
	Id          string   `json:"id"`
	TestcaseIds []string `json:"testcaseIds"`
}

type GradeResponseMessage struct {
	CompileOutput string                       `json:"compileOutput"`
	Result        []GradeResponseResultMessage `json:"result"`
	Metadata      GradeMetadataMessage         `json:"metadata"`
}

type GradeResponseResultMessage struct {
	Pass bool `json:"pass"`
	Time int  `json:"time"`
}
