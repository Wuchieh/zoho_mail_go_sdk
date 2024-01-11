package zoho_mail_go_sdk

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Data struct {
	AccountId string `json:"accountId"`
	Name      string `json:"firstName"`
	Email     string `json:"mailboxAddress"`
}

type accountResp struct {
	Status Status `json:"status"`
	Data   []Data `json:"data"`
}

type emailResp struct {
	Status Status `json:"status"`
}
