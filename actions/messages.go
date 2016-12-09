package actions

import (
	"github.com/zballs/comit/forms"
)

// Error messages

const (
	TmspRequestFailure  = "Failed to write TMSP request"
	TmspResponseFailure = "Failed to read TMSP request"

	HttpRequestFailure  = "Failed to read HTTP request"
	HttpResponseFailure = "Failed to write HTTP response"

	WebsocketFailure = "Failed to open websocket"

	InvalidPublicKey  = "Invalid public key"
	InvalidPrivateKey = "Invalid private key"
	InvalidFormID     = "Invalid form ID"
)

type CreateAccount struct {
	Result     string `json:"create account"`
	PubKeystr  string `json:"public key, omitempty"`
	PrivKeystr string `json:"private key, omitempty"`
}

type RemoveAccount struct {
	Result string `json:"remove account"`
}

type CreateAdmin struct {
	Result     string `json:"create admin"`
	PubKeystr  string `json:"public key, omitempty"`
	PrivKeystr string `json:"private key, omitempty"`
}

type RemoveAdmin struct {
	Result string `json:"remove admin"`
}

type Connect struct {
	Result string `json:"login"`
	Type   string `json:"type, omitempty"`
}

type SubmitForm struct {
	Result string `json:"submit form"`
	FormID string `json:"form ID, omitempty"`
}

type ResolveForm struct {
	Result string `json:"resolve form"`
}

type FindForm struct {
	Result string      `json:"find form"`
	Form   *forms.Form `json:"form, omitempty"`
}

var select_option = `<option value="%s">%s</option>`
