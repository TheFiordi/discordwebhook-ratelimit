package discordwebhook

import "testing"

func TestSendMessage(t *testing.T) {

	username := "TEST"
	content := "TESTCONTENT"

	type args struct {
		webhook string
		message Message
		proxy   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test", args: struct {
			webhook string
			message Message
			proxy   string
		}{webhook: "https://discord.com/api/webhooks/1232705275266596976/AkBkm3G11sU3RFUrWbKyj2a9f5Zo8Phq-jrC-CjDtNgrN6uJYfc158biFj7UzTEGb4dg", message: Message{
			Username: &username,
			Content:  &content,
		}, proxy: ""}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendMessage(tt.args.webhook, tt.args.message, tt.args.proxy); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
