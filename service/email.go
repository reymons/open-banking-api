package service

import (
	"context"
	"fmt"
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
	sesClient   *ses.Client
	email       string
	titlePrefix string
}

func NewEmailService(sc *ses.Client, email string, titlePrefix string) EmailService {
	return &emailService{sc, email, titlePrefix}
}

func (s *emailService) SendMessage(
	ctx context.Context,
	addrs []string,
	subject string,
	message string,
) error {
	sendEmailInput := &ses.SendEmailInput{
		Source: aws.String(s.email),
		Destination: &types.Destination{
			ToAddresses: addrs,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(fmt.Sprintf("%s - %s", s.titlePrefix, subject)),
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
