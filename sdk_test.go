package zoho_mail_go_sdk

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	clientID     = ""
	clientSecret = ""
	code         = ""
	sendTo       = ""
	content      = ""
)

func TestSdk(t *testing.T) {
	client, err := New(clientID, clientSecret, code)
	if err != nil {
		t.Fatal(err)
	}

	err = client.SendMail(&Mail{
		ToAddress: sendTo,
		Subject:   "sdk Test",
		Content:   content,
	})

	if err != nil {
		t.Fatal(err)
	}

	marshal, err := json.Marshal(client)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("client.json", marshal, 0644)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}

func TestSendMail(t *testing.T) {
	file, err := os.ReadFile("client.json")
	if err != nil {
		t.Fatal(err)
	}

	var client Client
	err = json.Unmarshal(file, &client)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		marshal, err := json.Marshal(client)
		if err != nil {
			t.Fatal(err)
		}

		err = os.WriteFile("client.json", marshal, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = client.SendMail(&Mail{
		ToAddress:  sendTo,
		Subject:    "sdk Test",
		Content:    content,
		AskReceipt: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}
