package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type EmailService interface {
	SendMessage(
		ctx context.Context,
		src string,
		addrs []string,
		subject string,
		message string,
	) error
}

type emailService struct {
	sesClient *ses.Client
}

func NewEmailService(sc *ses.Client) EmailService {
	return &emailService{sc}
}

func (s *emailService) SendMessage(
	ctx context.Context,
	src string,
	addrs []string,
	subject string,
	message string,
) error {
	sendEmailInput := &ses.SendEmailInput{
		Source: aws.String(src),
		Destination: &types.Destination{
			ToAddresses: addrs,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(subject),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data: aws.String(message),
				},
			},
		},
	}

	_, err := s.sesClient.SendEmail(ctx, sendEmailInput)
	return err
}
