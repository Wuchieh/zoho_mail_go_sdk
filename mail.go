package zoho_mail_go_sdk

type Mail struct {
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	AskReceipt  bool   `json:"askReceipt"`
}
