package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform"
	payload "github.com/codern-org/codern/platform/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

type gradingPublisher struct {
	ch *amqp.Channel
}

func NewGradingPublisher(rabbitmq *platform.RabbitMq) (domain.GradingPublisher, error) {
	_, err := rabbitmq.Ch.QueueDeclare("grading", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &gradingPublisher{ch: rabbitmq.Ch}, nil
}

func (p *gradingPublisher) Grade(assignment *domain.Assignment, submission *domain.Submission) {
	// TODO: implements retry strategy, currently drop a fail message partially
	eg, egctx := errgroup.WithContext(context.Background())

	for i := range assignment.Testcases {
		testcase := assignment.Testcases[i]
		eg.Go(func() error {
			select {
			case <-egctx.Done():
				return egctx.Err()
			default:
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				// Follow legacy version message
				message := &payload.GradeMessage{
					Id:   fmt.Sprintf("%d.%d", submission.Id, testcase.Id),
					Type: submission.Language,
					Settings: &payload.GradeSettingsMessage{
						LimitMemory: assignment.MemoryLimit,
						LimitTime:   assignment.TimeLimit,
					},
					Files: []payload.GradeFileMessage{
						{
							Name:       "source",
							SourceType: "URL",
							Source:     submission.FileUrl,
						},
						{
							Name:       "testcase.zip",
							SourceType: "URL",
							Source:     testcase.FileUrl,
						},
					},
				}
				body, err := json.Marshal(message)
				if err != nil {
					return err
				}

				return p.ch.PublishWithContext(ctx, "", "grading", false, false, amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				})
			}
		})
	}
}
