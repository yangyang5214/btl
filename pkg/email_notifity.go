package pkg

import (
	"bytes"
	"encoding/json"
	"github.com/go-gomail/gomail"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

type EmailConfig struct {
	To       []string `json:"to,omitempty"`
	From     string   `json:"from,omitempty"`
	Password string   `json:"password,omitempty"`
	SmtpHost string   `json:"smtp_host,omitempty"`
	SmtpPort int      `json:"smtp_port,omitempty"`
}

type EmailContent struct {
	Subject string
	Content string
	Images  []string
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
	gm     *gomail.Message
}

func NewEmailNotify(config *EmailConfig) *EmailNotify {
	return &EmailNotify{
		config: config,
		gm:     gomail.NewMessage(),
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

func (e *EmailNotify) Send(to []string, content *EmailContent) error {
	e.gm.SetHeader("From", e.config.From)
	e.gm.SetHeader("To", to...)
	e.gm.SetHeader("Subject", content.Subject)
	e.gm.SetBody("text/plain", content.Content)

	for _, image := range content.Images {
		e.gm.Attach(image)
	}

	d := gomail.NewDialer(e.config.SmtpHost, e.config.SmtpPort, e.config.From, e.config.Password)
	if err := d.DialAndSend(e.gm); err != nil {
		return err
	}
	return nil
}
