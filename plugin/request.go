package plugin

import (
	"fmt"
	"strings"
)

type Header map[string]string

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Params []Param

type Request struct {
	Method    string
	UserAgent string
	Header    Header
	Auth      Auth
	Queryset  Params
	Body      Params
	File      Params
	Json      interface{}
}

func (r Request) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Method: %q", r.Method))
	if r.UserAgent != "" {
		builder.WriteString(fmt.Sprintf(", User-Agent: %q", r.UserAgent))
	}

	if r.Auth.Username != "" && r.Auth.Password != "" {
		builder.WriteString(fmt.Sprintf(", Auth: %q:%s", r.Auth.Username, r.Auth.Password))
	}

	if len(r.Header) != 0 {
		builder.WriteString(fmt.Sprintf(", Header: %#v", r.Header))
	}

	if r.Queryset != nil {
		builder.WriteString(fmt.Sprintf(", Queryset: %#v", r.Queryset))
	}

	if r.Body != nil {
		builder.WriteString(fmt.Sprintf(", Body: %#v", r.Body))
	}

	if r.File != nil {
		builder.WriteString(fmt.Sprintf(", File: %#v", r.File))
	}

	if r.Json != nil {
		builder.WriteString(fmt.Sprintf(", Json: %#v", r.Json))
	}

	return builder.String()
}
