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
	rabbitMq          *platform.RabbitMq
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
) error {
	consumer := &gradingConsumer{
		logger:            logger,
		rabbitMq:          rabbitmq,
		wsHub:             wsHub,
		influxDb:          influxDb,
		assignmentUsecase: assignmentUsecase,
	}
	return consumer.startConsumers()
}

func (c *gradingConsumer) startConsumers() error {
	if err := c.rabbitMq.Consume("grading_response", c.readSubmssionResult); err != nil {
		return err
	}
	return nil
}

func (c *gradingConsumer) readSubmssionResult(delivery amqp.Delivery) {
	var message payload.GradeResponseMessage

	if err := json.Unmarshal(delivery.Body, &message); err != nil {
		delivery.Reject(false)
		c.logger.Error("Cannot unmarshal GradeResponseMessage", zap.Error(err))
		return
	}

	assignmentId := message.Metadata.AssignmentId
	submissionId := message.Metadata.SubmissionId
	results := make([]domain.SubmissionResult, 0)

	assignment, err := c.assignmentUsecase.Get(assignmentId)
	if err != nil {
		delivery.Reject(false)
		c.logger.Error("Cannot get assignment when consuming submission result", zap.Error(err))
		return
	}

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
		assignment,
		submissionId,
		message.CompileOutput,
		results,
	); err != nil {
		delivery.Reject(true)
		c.logger.Error("Cannot create submission results", zap.Error(err))
		return
	}

	submission, err := c.assignmentUsecase.GetSubmission(submissionId)
	if err != nil || submission == nil {
		delivery.Reject(false)
		c.logger.Error("Cannot get submission when consuming submission result", zap.Error(err))
		return
	}

	c.influxDb.WritePoint(
		"gradingLatency", map[string]string{
			"language": submission.Language,
		},
		map[string]interface{}{
			"userId":       submission.SubmitterId,
			"assignmentId": submission.AssignmentId,
			"latency":      time.Since(message.Metadata.StartTime).Nanoseconds(),
		},
	)

	if err := c.wsHub.SendMessage(submission.SubmitterId, "onSubmissionUpdate", submission); err != nil {
		delivery.Reject(false)
		c.logger.Error("Cannot send websocket message after consuming submission result", zap.Error(err))
		return
	}

	c.logger.Info("Consumed submission result", zap.Int("submission_id", submissionId))
	delivery.Ack(true)
}
