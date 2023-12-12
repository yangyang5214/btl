package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/smtp"
	"os"
	"path"
)

type EmailConfig struct {
	From     string `json:"from,omitempty"`
	Password string `json:"password,omitempty"`
	SmtpHost string `json:"smtp_host,omitempty"`
	SmtpPort int    `json:"smtp_port,omitempty"`
}

type EmailContent struct {
	Subject string
	Content string
	images  []string //todo
}

func (e *EmailContent) String() []byte {
	var message bytes.Buffer
	message.WriteString("Subject: " + e.Subject + "\n")

	message.WriteString("\r\n")
	message.WriteString(e.Content)
	message.WriteString("\r\n")

	return message.Bytes()
}

type EmailNotify struct {
	config *EmailConfig
}

func NewEmailNotify(config *EmailConfig) *EmailNotify {
	return &EmailNotify{
		config: config,
	}
}

func LoadConfigFromEnv() (*EmailConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	jsonData, err := os.ReadFile(path.Join(homeDir, ".email"))
	if err != nil {
		return nil, err
	}
	var config EmailConfig
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		log.Errorf("json unmarshal <%s> error: %v", jsonData, err)
		return nil, err
	}
	return &config, nil
}

func (e *EmailNotify) addr() string {
	return fmt.Sprintf("%s:%d", e.config.SmtpHost, e.config.SmtpPort)
}

func (e *EmailNotify) Send(to []string, content *EmailContent) error {
	auth := smtp.PlainAuth("", e.config.From, e.config.Password, e.config.SmtpHost)
	err := smtp.SendMail(e.addr(), auth, e.config.From, to, content.String())
	if err != nil {
		log.Errorf("send email failed: %+v", err)
		return err
	}
	return nil
}
