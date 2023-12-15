package model

type Image struct {
	ID         string   `json:"id"`
	S3URL      string   `json:"s3_url"`
	Seed       int      `json:"seed"`
	Prompt     string   `json:"prompt"`
	Data       []byte   `json:"data"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
	Styles     []string `json:"styles"`
	CreatedAt  int64    `json:"created_at"`
	CreatedBy  string   `json:"created_by"`
}
