package plugin

import (
	"time"
)

type AuthParams []Auth
type QueryParams []Params
type BodyParams []Params
type FileParams []Params
type JsonParams []interface{}

type POCRequest struct {
	Method    string `json:"method"`
	UserAgent string `json:"user-agent"`
	Header    Header `json:"header"`

	Protocols []string `json:"protocols"`
	Ports     []int    `json:"ports"`
	Suffixes  []string `json:"suffixes"`

	Auths     AuthParams  `json:"auths"`
	Querysets QueryParams `json:"querysets"`
	Bodys     BodyParams  `json:"bodys"`
	Jsons     JsonParams  `json:"jsons"`
	Files     FileParams  `json:"files"`
}

type POC struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Desc       string        `json:"desc"`
	Author     string        `json:"author"`
	Version    string        `json:"version"`
	References []string      `json:"references"`
	Request    POCRequest    `json:"request"`
	Matcher    Matcher       `json:"matcher"`
	Timeout    time.Duration `json:"-"`
	isBuild    bool
	plugins    []*Plugin
}

func (poc *POC) BuildPlugins() []*Plugin {
	if poc.isBuild {
		return poc.plugins
	}
	poc.isBuild = true

	auths := make([]Auth, 1)
	if len(poc.Request.Auths) > 0 {
		auths = poc.Request.Auths
	}
	querysets := make(QueryParams, 1)
	if len(poc.Request.Querysets) > 0 {
		querysets = poc.Request.Querysets
	}

	for _, auth := range auths {
		for _, queryset := range querysets {
			plugin := &Plugin{
				POC: poc,
				Request: Request{
					Method:    poc.Request.Method,
					UserAgent: poc.Request.UserAgent,
					Header:    poc.Request.Header,
					Auth:      auth,
					Queryset:  queryset,
				},
			}
			poc.plugins = append(poc.plugins, plugin)
			for _, body := range poc.Request.Bodys {
				plugin := &Plugin{
					POC: poc,
					Request: Request{
						Method:    poc.Request.Method,
						UserAgent: poc.Request.UserAgent,
						Header:    poc.Request.Header,
						Auth:      auth,
						Queryset:  queryset,
						Body:      body,
					},
				}
				poc.plugins = append(poc.plugins, plugin)
			}
			for _, json := range poc.Request.Jsons {
				plugin := &Plugin{
					POC: poc,
					Request: Request{
						Method:    poc.Request.Method,
						UserAgent: poc.Request.UserAgent,
						Header:    poc.Request.Header,
						Auth:      auth,
						Queryset:  queryset,
						Json:      json,
					},
				}
				poc.plugins = append(poc.plugins, plugin)
			}
			for _, file := range poc.Request.Files {
				plugin := &Plugin{
					POC: poc,
					Request: Request{
						Method:    poc.Request.Method,
						UserAgent: poc.Request.UserAgent,
						Header:    poc.Request.Header,
						Auth:      auth,
						Queryset:  queryset,
						File:      file,
					},
				}
				poc.plugins = append(poc.plugins, plugin)
			}
		}
	}
	return poc.plugins
}
