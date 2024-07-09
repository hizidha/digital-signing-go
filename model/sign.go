package model

type Sign struct {
	UUID           string `json:"uuid"`
	UserUUID       string `json:"user_uuid"`
	TitleDoc       string `json:"title_document"`
	Category       string `json:"category"`
	Description    string `json:"description"`
	Signature      string `json:"signature"`
	LinktoVerify   string `json:"link_verify"`
	SignTimestamp  string `json:"sign_timestamp"`
	StorageAddress string `json:"storage_address"`
}
