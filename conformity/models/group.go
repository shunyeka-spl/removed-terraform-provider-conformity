package models

type GroupAccount struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

type Group struct {
	Type       string `json:"type,omitempty"`
	ID         string `json:"id,omitempty"`
	Attributes struct {
		Name             string   `json:"name"`
		Tags             []string `json:"tags"`
		CreatedDate      int64    `json:"created-date,omitempty"`
		LastModifiedDate int64    `json:"last-modified-date,omitempty"`
	} `json:"attributes"`
	Relationships struct {
		Organisation struct {
			Data struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"data"`
		} `json:"organisation,omitempty"`
		Accounts struct {
			Data []GroupAccount `json:"data,omitempty"`
		} `json:"accounts,omitempty"`
	} `json:"relationships,omitempty"`
}
