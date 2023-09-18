package publisher

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform"
	payload "github.com/codern-org/codern/platform/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
)

type gradingPublisher struct {
	ch *amqp.Channel
}

func NewGradingPublisher(rabbitmq *platform.RabbitMq) domain.GradingPublisher {
	return &gradingPublisher{
		ch: rabbitmq.Ch,
	}
}

func (p *gradingPublisher) Grade(assignment *domain.Assignment, submission *domain.Submission) error {
	testcaseIds := make([]string, 0)
	testcases := make([]payload.GradeTestMessage, 0)

	for i := range assignment.Testcases {
		testcases = append(testcases, payload.GradeTestMessage{
			InputUrl:  assignment.Testcases[i].InputFileUrl,
			OutputUrl: assignment.Testcases[i].OutputFileUrl,
		})
		testcaseIds = append(testcaseIds, strconv.Itoa(assignment.Testcases[i].Id))
	}

	message := &payload.GradeRequestMessage{
		Language:  submission.Language,
		SourceUrl: submission.FileUrl,
		Test:      testcases,
		Metadata: payload.GradeMetadataMessage{
			Id:          strconv.Itoa(submission.Id),
			TestcaseIds: testcaseIds,
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
