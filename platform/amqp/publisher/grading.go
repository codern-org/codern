package publisher

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/platform"
	payload "github.com/codern-org/codern/platform/amqp"
)

type gradingPublisher struct {
	cfg      *config.Config
	rabbitMq *platform.RabbitMq
}

func NewGradingPublisher(cfg *config.Config, rabbitmq *platform.RabbitMq) domain.GradingPublisher {
	return &gradingPublisher{
		cfg:      cfg,
		rabbitMq: rabbitmq,
	}
}

func (p *gradingPublisher) Grade(assignment *domain.AssignmentWithStatus, submission *domain.Submission) error {
	testcaseIds := make([]int, 0)
	testcases := make([]payload.GradeTestMessage, 0)

	for i := range assignment.Testcases {
		inputUrl, err := url.JoinPath(
			p.cfg.Client.SeaweedFs.FilerUrls.External,
			assignment.Testcases[i].InputFileUrl,
		)
		if err != nil {
			return errs.New(errs.ErrCreateUrlPath, "invalid testcase input url", err)
		}
		outputUrl, err := url.JoinPath(
			p.cfg.Client.SeaweedFs.FilerUrls.External,
			assignment.Testcases[i].OutputFileUrl,
		)
		if err != nil {
			return errs.New(errs.ErrCreateUrlPath, "invalid testcase output url", err)
		}

		testcases = append(testcases, payload.GradeTestMessage{
			InputUrl:  inputUrl,
			OutputUrl: outputUrl,
		})
		testcaseIds = append(testcaseIds, assignment.Testcases[i].Id)
	}

	sourceUrl, err := url.JoinPath(p.cfg.Client.SeaweedFs.FilerUrls.External, submission.FileUrl)
	if err != nil {
		return errs.New(errs.ErrCreateUrlPath, "invalid submission url", err)
	}

	message := &payload.GradeRequestMessage{
		Language:  submission.Language,
		SourceUrl: sourceUrl,
		Test:      testcases,
		Settings: payload.GradeSettingsMessage{
			TimeLimit:         assignment.TimeLimit,
			MemoryLimit:       assignment.MemoryLimit,
			IsAutoTrimEnabled: assignment.IsAutoTrimEnabled,
		},
		Metadata: payload.GradeMetadataMessage{
			AssignmentId: assignment.Id,
			SubmissionId: submission.Id,
			TestcaseIds:  testcaseIds,
			StartTime:    time.Now(),
		},
	}
	body, err := json.Marshal(message)
	if err != nil {
		return errs.New(errs.ErrGradingRequest, "cannot marshal grading request message", err)
	}

	if err := p.rabbitMq.Publish("grading", "request", body); err != nil {
		return errs.New(errs.ErrGradingRequest, "cannot publish grading request message", err)
	}

	return nil
}
