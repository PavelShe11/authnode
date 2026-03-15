package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"net/smtp"
	"sync"
	"time"

	"github.com/PavelShe11/authnode/authMicro/internal/config"
	"github.com/PavelShe11/authnode/common/logger"
	"github.com/PavelShe11/authnode/common/translator"
)

type emailJob struct {
	to, code, lang string
}

type SmtpEmailSender struct {
	config     config.SmtpConfig
	translator *translator.Translator
	codeTTL    time.Duration
	logger     logger.Logger
	queue      chan emailJob
	wg         sync.WaitGroup
}

func NewSmtpEmailSender(cfg config.SmtpConfig, trans *translator.Translator, codeTTL time.Duration, log logger.Logger) *SmtpEmailSender {
	s := &SmtpEmailSender{
		config:     cfg,
		translator: trans,
		codeTTL:    codeTTL,
		logger:     log,
		queue:      make(chan emailJob, 100),
	}
	s.wg.Add(1)
	go s.worker()
	return s
}

// SendVerificationCode ставит задачу в очередь и немедленно возвращает управление.
func (s *SmtpEmailSender) SendVerificationCode(_ context.Context, to, code, lang string) error {
	select {
	case s.queue <- emailJob{to, code, lang}:
		return nil
	default:
		return fmt.Errorf("email queue is full")
	}
}

// Close дожидается отправки всех писем из очереди и закрывает соединение.
func (s *SmtpEmailSender) Close() {
	close(s.queue)
	s.wg.Wait()
}

func (s *SmtpEmailSender) worker() {
	defer s.wg.Done()

	var client *smtp.Client
	defer func() {
		if client != nil {
			_ = client.Quit()
		}
	}()

	for job := range s.queue {
		var err error
		for attempt := 1; attempt <= 3; attempt++ {
			client, err = s.ensureConnection(client)
			if err != nil {
				s.logger.Errorf("smtp connect attempt %d failed: %v", attempt, err)
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			if err = s.sendViaClient(client, job); err != nil {
				s.logger.Errorf("smtp send attempt %d failed: %v", attempt, err)
				_ = client.Close()
				client = nil
				continue
			}
			break
		}
		if err != nil {
			s.logger.Errorf("failed to send email to %s after 3 attempts: %v", job.to, err)
		}
	}
}

func (s *SmtpEmailSender) ensureConnection(client *smtp.Client) (*smtp.Client, error) {
	if client != nil {
		if err := client.Noop(); err == nil {
			return client, nil
		}
		_ = client.Close()
	}
	return s.connect()
}

func (s *SmtpEmailSender) connect() (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	client, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("smtp dial: %w", err)
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{ServerName: s.config.Host}
		if err := client.StartTLS(tlsCfg); err != nil {
			_ = client.Close()
			return nil, fmt.Errorf("starttls: %w", err)
		}
	}

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	if err := client.Auth(auth); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("smtp auth: %w", err)
	}

	return client, nil
}

func (s *SmtpEmailSender) sendViaClient(client *smtp.Client, job emailJob) error {
	minutes := int(math.Ceil(s.codeTTL.Minutes()))
	params := map[string]interface{}{"Code": job.code, "Minutes": minutes}

	subject := s.translator.Translate("emailVerificationSubject", nil, job.lang)
	body := s.translator.Translate("emailVerificationBody", params, job.lang)

	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.config.From, job.to, subject, body,
	)

	if err := client.Mail(s.config.From); err != nil {
		return fmt.Errorf("MAIL FROM: %w", err)
	}
	if err := client.Rcpt(job.to); err != nil {
		return fmt.Errorf("RCPT TO: %w", err)
	}
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA: %w", err)
	}
	if _, err = fmt.Fprint(wc, message); err != nil {
		_ = wc.Close()
		return fmt.Errorf("write data: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("close data: %w", err)
	}
	return client.Reset()
}
