package consumer

import (
	"encoding/json"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform"
	payload "github.com/codern-org/codern/platform/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type gradingConsumer struct {
	logger           *zap.Logger
	ch               *amqp.Channel
	wsHub            *platform.WebSocketHub
	workspaceUsecase domain.WorkspaceUsecase
}

func NewGradingConsumer(
	logger *zap.Logger,
	rabbitmq *platform.RabbitMq,
	wsHub *platform.WebSocketHub,
	workspaceUsecase domain.WorkspaceUsecase,
) domain.GradingConsumer {
	return &gradingConsumer{
		logger:           logger,
		ch:               rabbitmq.Ch,
		wsHub:            wsHub,
		workspaceUsecase: workspaceUsecase,
	}
}

func (c *gradingConsumer) ConsumeSubmssionResult() error {
	messages, err := c.ch.Consume("grading_response", "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range messages {
			var message payload.GradeResponseMessage

			if err := json.Unmarshal(delivery.Body, &message); err != nil {
				delivery.Reject(false)
				c.logger.Error("Cannot unmarshal GradeResponseMessage", zap.Error(err))
				continue
			}

			submissionId := message.Metadata.SubmissionId
			results := make([]domain.SubmissionResult, 0)

			for i := range message.Results {
				status := domain.SubmissionResultError
				if message.Results[i].Pass {
					status = domain.SubmissionResultDone
				}
				results = append(results, domain.SubmissionResult{
					SubmissionId: submissionId,
					TestcaseId:   message.Metadata.TestcaseIds[i],
					Status:       status,
					StatusDetail: nil,
					MemoryUsage:  nil,
					TimeUsage:    &message.Results[i].Time,
				})
			}

			err = c.workspaceUsecase.UpdateSubmissionResults(submissionId, message.CompileOutput, results)
			if err != nil {
				delivery.Reject(false)
				c.logger.Error("Cannot update submission results", zap.Error(err))
				continue
			}

			submission, err := c.workspaceUsecase.GetSubmission(submissionId)
			if err != nil || submission == nil {
				delivery.Reject(false)
				c.logger.Error("Cannot get submission data when consuming submission result", zap.Error(err))
				continue
			}
			err = c.wsHub.SendMessage(submission.UserId, "onSubmissionUpdate", submission)
			if err != nil {
				delivery.Reject(false)
				c.logger.Error("Cannot send websocket message after consuming submission result", zap.Error(err))
				continue
			}

			c.logger.Info("Consumed submission result", zap.Int("submission_id", submissionId))
			delivery.Ack(false)
		}
	}()

	return nil
}
