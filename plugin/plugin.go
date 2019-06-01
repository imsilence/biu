package plugin

import (
	"net/http"

	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

var requests map[string]func(string, *grequests.RequestOptions) (*grequests.Response, error)

type Plugin struct {
	POC            *POC
	Request        Request
	RequestOptions *grequests.RequestOptions
}

func (plugin *Plugin) BuildRequestOptions() {
	if plugin.RequestOptions == nil {
		data := map[string]string{}
		params := map[string]string{}

		for _, v := range plugin.Request.Queryset {
			params[v.Key] = v.Value
		}

		for _, v := range plugin.Request.Body {
			data[v.Key] = v.Value
		}

		var auth []string
		if plugin.Request.Auth.Username != "" && plugin.Request.Auth.Password != "" {
			auth = []string{plugin.Request.Auth.Username, plugin.Request.Auth.Password}
		}

		plugin.RequestOptions = &grequests.RequestOptions{
			Data:               data,
			Params:             params,
			JSON:               plugin.Request.Json,
			Headers:            plugin.Request.Header,
			Auth:               auth,
			RequestTimeout:     plugin.POC.Timeout,
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
				"poc":      plugin.POC.Name,
				"request":  *plugin.RequestOptions,
				"response": *response,
				"error":    err,
			}).Debug("request target")
			if plugin.POC.Matcher.Match(response) {
				return NewUnSafety(target, plugin)
			} else {
				return NewSafety(target, plugin)
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"request": plugin.Request,
				"method":  plugin.Request.Method,
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
