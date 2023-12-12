package pkg

import "testing"

func TestSendEmail(t *testing.T) {
	config, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	email := NewEmailNotify(config)
	content := &EmailContent{
		Subject: "测试",
		Content: "嘿嘿",
		Images: []string{
			"/Users/beer/Downloads/IMG_0018.png",
		},
	}
	err = email.Send(config.To, content)
	if err != nil {
		t.Fatal(err)
	}
}
