package sms

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	appconfig "github.com/your-org/jvairv2/configs"
)

// AWSSNSSender implementa SMSSender usando AWS SNS
type AWSSNSSender struct {
	client *sns.Client
}

// NewAWSSNSSender crea un nuevo cliente AWS SNS usando la configuración S3 existente
func NewAWSSNSSender(cfg *appconfig.S3Config) (*AWSSNSSender, error) {
	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
	}

	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("error al cargar configuración AWS para SNS: %w", err)
	}

	slog.Info("[AWSSNSSender] cliente inicializado",
		slog.String("region", cfg.Region),
		slog.Bool("conCredenciales", cfg.AccessKey != ""),
	)

	return &AWSSNSSender{
		client: sns.NewFromConfig(awsCfg),
	}, nil
}

// SendSMS envía un SMS via AWS SNS a un número en formato E.164 (+1XXXXXXXXXX)
func (s *AWSSNSSender) SendSMS(ctx context.Context, to string, body string) error {
	slog.InfoContext(ctx, "[AWSSNSSender] publicando mensaje SNS",
		slog.String("to", to),
		slog.Int("bodyLen", len(body)),
	)

	output, err := s.client.Publish(ctx, &sns.PublishInput{
		PhoneNumber: aws.String(to),
		Message:     aws.String(body),
	})
	if err != nil {
		slog.ErrorContext(ctx, "[AWSSNSSender] error en Publish",
			slog.String("to", to),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("error al enviar SMS via AWS SNS a %s: %w", to, err)
	}
	if output.MessageId == nil || *output.MessageId == "" {
		slog.WarnContext(ctx, "[AWSSNSSender] SNS no retornó MessageId", slog.String("to", to))
		return fmt.Errorf("AWS SNS no retornó MessageId para %s", to)
	}

	slog.InfoContext(ctx, "[AWSSNSSender] SMS publicado correctamente",
		slog.String("to", to),
		slog.String("messageId", *output.MessageId),
	)
	return nil
}
