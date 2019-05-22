package plugin

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

var requests map[string]func(string, *grequests.RequestOptions) (*grequests.Response, error)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Value struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Request struct {
	Protocols []string          `json:"protocols"`
	Method    string            `json:"method"`
	Ports     []int             `json:"ports"`
	Suffixes  []string          `json:"suffixes"`
	UserAgent string            `json:"user-agent"`
	Headers   map[string]string `json:"headers"`
	Auth      Auth              `json:"auth"`
	Queryset  []Value           `json:"queryset"`
	Bodys     []Value           `json:"bodys"`
	JSON      interface{}       `json:"json"`
	Files     []Value           `json:"files"`
}

type StatusMatcher int

func (matcher StatusMatcher) Match(response *grequests.Response) bool {
	return int(matcher) == response.StatusCode
}

type TextMatcher struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (matcher TextMatcher) Match(response *grequests.Response) bool {
	value := response.String()
	switch matcher.Type {
	case "contains":
		return strings.Contains(value, matcher.Value)
	case "icontains":
		return strings.Contains(strings.ToLower(value), strings.ToLower(matcher.Value))
	case "regexp":
		matched, err := regexp.MatchString(matcher.Value, value)
		if err == nil {
			return matched
		} else {
			logrus.WithFields(logrus.Fields{
				"matcher": matcher,
				"error":   err,
			}).Errorf("error header matcher")
		}
	}
	return false
}

type HeaderMatcher struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (matcher HeaderMatcher) Match(response *grequests.Response) bool {
	value := response.Header.Get(matcher.Key)
	switch matcher.Type {
	case "contains":
		return strings.Contains(value, matcher.Value)
	case "icontains":
		return strings.Contains(strings.ToLower(value), strings.ToLower(matcher.Value))
	case "regexp":
		matched, err := regexp.MatchString(matcher.Value, value)
		if err == nil {
			return matched
		} else {
			logrus.WithFields(logrus.Fields{
				"matcher": matcher,
				"error":   err,
			}).Errorf("error header matcher")
		}
	}
	return false
}

type Matcher struct {
	Status  []StatusMatcher `json:"status"`
	Headers []HeaderMatcher `json:"headers"`
	Texts   []TextMatcher   `json:"texts"`
}

func (matcher Matcher) Match(response *grequests.Response) bool {
	for _, m := range matcher.Status {
		if m.Match(response) {
			return true
		}
	}

	for _, m := range matcher.Texts {
		if m.Match(response) {
			return true
		}
	}

	for _, m := range matcher.Headers {
		if m.Match(response) {
			return true
		}
	}
	return false
}

type Plugin struct {
	Id             string                    `json:"id"`
	Name           string                    `json:"name"`
	Desc           string                    `json:"desc"`
	Author         string                    `json:"author"`
	Version        string                    `json:"version"`
	References     []string                  `json:"references"`
	Request        Request                   `json:"request"`
	Matcher        Matcher                   `json:"matcher"`
	RequestOptions *grequests.RequestOptions `json:"-"`
	Timeout        time.Duration             `json:"-"`
}

func (plugin *Plugin) BuildRequestOptions() {
	if plugin.RequestOptions == nil {
		data := map[string]string{}
		params := map[string]string{}

		for _, v := range plugin.Request.Queryset {
			params[v.Key] = v.Value
		}

		for _, v := range plugin.Request.Bodys {
			data[v.Key] = v.Value
		}

		var auth []string
		if plugin.Request.Auth.Username != "" && plugin.Request.Auth.Password != "" {
			auth = []string{plugin.Request.Auth.Username, plugin.Request.Auth.Password}
		}

		plugin.RequestOptions = &grequests.RequestOptions{
			Data:               data,
			Params:             params,
			JSON:               plugin.Request.JSON,
			Headers:            plugin.Request.Headers,
			Auth:               auth,
			RequestTimeout:     plugin.Timeout,
			InsecureSkipVerify: true,
			UserAgent:          plugin.Request.UserAgent,
		}
	}
}

func (plugin *Plugin) Execute(target Target) Result {
	plugin.BuildRequestOptions()
	switch target.Type {
	case URL:
		if request, ok := requests[plugin.Request.Method]; ok {
			response, err := request(target.Raw, plugin.RequestOptions)
			logrus.WithFields(logrus.Fields{
				"target":   target.Raw,
				"plugin":   plugin.Name,
				"request":  *plugin.RequestOptions,
				"response": *response,
				"error":    err,
			}).Debug("request target")
			if plugin.Matcher.Match(response) {
				return NewUnSafety(target, plugin)
			} else {
				return NewSafety(target, plugin)
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"plugin": plugin.Name,
				"method": plugin.Request.Method,
			}).Error("error plugin request method")
		}
	}
	return NewSafety(target, plugin)
}

func init() {
	requests = make(map[string]func(string, *grequests.RequestOptions) (*grequests.Response, error))
	requests[http.MethodHead] = grequests.Head
	requests[http.MethodGet] = grequests.Get
	requests[http.MethodPost] = grequests.Post
}
