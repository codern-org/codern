package publisher

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/platform"
	payload "github.com/codern-org/codern/platform/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
)

type gradingPublisher struct {
	cfg *config.Config
	ch  *amqp.Channel
}

func NewGradingPublisher(cfg *config.Config, rabbitmq *platform.RabbitMq) domain.GradingPublisher {
	return &gradingPublisher{
		cfg: cfg,
		ch:  rabbitmq.Ch,
	}
}

func (p *gradingPublisher) Grade(assignment *domain.Assignment, submission *domain.Submission) error {
	testcaseIds := make([]int, 0)
	testcases := make([]payload.GradeTestMessage, 0)

	for i := range assignment.Testcases {
		testcases = append(testcases, payload.GradeTestMessage{
			InputUrl:  assignment.Testcases[i].InputFileUrl,
			OutputUrl: assignment.Testcases[i].OutputFileUrl,
		})
		testcaseIds = append(testcaseIds, assignment.Testcases[i].Id)
	}

	// TODO: hardcoded filer url
	sourceUrl, err := url.JoinPath(p.cfg.Client.SeaweedFs.FilerUrls[1], submission.FileUrl)
	if err != nil {
		// TODO: domain error
		return err
	}

	message := &payload.GradeRequestMessage{
		Language:  submission.Language,
		SourceUrl: sourceUrl,
		Test:      testcases,
		Metadata: payload.GradeMetadataMessage{
			SubmissionId: submission.Id,
			TestcaseIds:  testcaseIds,
		},
	}
	body, err := json.Marshal(message)
	if err != nil {
		// TODO: domain error
		return err
	}

	err = p.ch.PublishWithContext(context.Background(), "grading", "request", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		// TODO: domain error
		return err
	}

	return nil
}
