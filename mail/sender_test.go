package mail

import (
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
)

func TestGmailSender_SendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
		<h1>Hello World</h1>
		<p>This is a test email message from <a href="devbrian.dev.app">Brian Dev</a></p>
	`
	to := []string{"brnmwas@gmail.com"}
	attachFiles := []string{"../ReadMe.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)

}
