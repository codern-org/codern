package domain

type GradingPublisher interface {
	Grade(assignment *AssignmentWithStatus, submission *Submission) error
}

type GradingConsumer interface {
	ConsumeSubmssionResult() error
}
