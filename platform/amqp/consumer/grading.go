package consumer

import (
	"encoding/json"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform"
	payload "github.com/codern-org/codern/platform/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type gradingConsumer struct {
	logger            *zap.Logger
	ch                *amqp.Channel
	wsHub             *platform.WebSocketHub
	influxDb          *platform.InfluxDb
	assignmentUsecase domain.AssignmentUsecase
}

func NewGradingConsumer(
	logger *zap.Logger,
	rabbitmq *platform.RabbitMq,
	wsHub *platform.WebSocketHub,
	influxDb *platform.InfluxDb,
	assignmentUsecase domain.AssignmentUsecase,
) domain.GradingConsumer {
	return &gradingConsumer{
		logger:            logger,
		ch:                rabbitmq.Ch,
		wsHub:             wsHub,
		influxDb:          influxDb,
		assignmentUsecase: assignmentUsecase,
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
				results = append(results, domain.SubmissionResult{
					SubmissionId: submissionId,
					TestcaseId:   message.Metadata.TestcaseIds[i],
					IsPassed:     message.Results[i].Pass,
					Status:       message.Status,
					MemoryUsage:  &message.Results[i].Memory,
					TimeUsage:    &message.Results[i].Time,
				})
			}

			if err := c.assignmentUsecase.CreateSubmissionResults(
				submissionId,
				message.CompileOutput,
				results,
			); err != nil {
				delivery.Reject(true)
				c.logger.Error("Cannot create submission results", zap.Error(err))
				continue
			}

			submission, err := c.assignmentUsecase.GetSubmission(submissionId)
			if err != nil || submission == nil {
				delivery.Reject(false)
				c.logger.Error("Cannot get submission data when consuming submission result", zap.Error(err))
				continue
			}

			if err := c.wsHub.SendMessage(submission.UserId, "onSubmissionUpdate", submission); err != nil {
				delivery.Reject(false)
				c.logger.Error("Cannot send websocket message after consuming submission result", zap.Error(err))
				continue
			}

			c.logger.Info("Consumed submission result", zap.Int("submission_id", submissionId))
			delivery.Ack(true)

			c.influxDb.WritePoint(
				"gradingLatency", map[string]string{
					"language": submission.Language,
				},
				map[string]interface{}{
					"userId":       submission.UserId,
					"assignmentId": submission.AssignmentId,
					"latency":      time.Since(message.Metadata.StartTime).Nanoseconds(),
				},
			)
		}
	}()

	return nil
}
