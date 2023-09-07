package consumer

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform"
	payload "github.com/codern-org/codern/platform/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
)

type gradingConsumer struct {
	ch               *amqp.Channel
	workspaceUsecase domain.WorkspaceUsecase
}

func NewGradingConsumer(
	rabbitmq *platform.RabbitMq,
	workspaceUsecase domain.WorkspaceUsecase,
) domain.GradingConsumer {
	return &gradingConsumer{
		ch:               rabbitmq.Ch,
		workspaceUsecase: workspaceUsecase,
	}
}

func (c *gradingConsumer) ConsumeSubmssionResult() error {
	messages, err := c.ch.Consume("grading", "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// TODO: improve failover strategy
	go func() {
		for delivery := range messages {
			var message payload.GradeResultMessage
			if err := json.Unmarshal(delivery.Body, &message); err != nil {
				delivery.Reject(true)
				return
			}

			ids := strings.Split(message.Id, ".")
			submissionId, err := strconv.Atoi(ids[0])
			if err != nil {
				delivery.Reject(true)
				return
			}
			testcaseId, err := strconv.Atoi(ids[1])
			if err != nil {
				delivery.Reject(true)
				return
			}

			c.workspaceUsecase.UpdateSubmissionResult(&domain.SubmissionResult{
				SubmissionId:   submissionId,
				TestcaseId:     testcaseId,
				Status:         "DONE", // TODO: correctness
				StatusDetail:   &message.Status,
				MemoryUsage:    &message.Metadata.MemoryUsage,
				TimeUsage:      &message.Metadata.TimeUsage,
				CompilationLog: &message.Metadata.CompilationLog,
			})
		}
	}()

	return nil
}
