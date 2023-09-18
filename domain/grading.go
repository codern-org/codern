package domain

type GradingPublisher interface {
	Grade(assignment *Assignment, submission *Submission) error
}

type GradingConsumer interface {
	ConsumeSubmssionResult() error
}
