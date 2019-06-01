package plugin

import (
	"regexp"
	"strings"

	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

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
