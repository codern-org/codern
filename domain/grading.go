package domain

type GradingPublisher interface {
	Grade(assignment *Assignment, submission *Submission)
}

type GradingConsumer interface {
	ConsumeSubmssionResult() error
}
