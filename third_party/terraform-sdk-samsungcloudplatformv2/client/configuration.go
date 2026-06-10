package scpsdk

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// contextKeys are used to identify the type of value in the context.
// Since these are string, it is possible to get a short description of the
// context key for logging and debugging using key.String().

type contextKey string

func (c contextKey) String() string {
	return "auth " + string(c)
}

type Configuration struct {
	AuthUrl         string
	ServiceType     string
	AllowSDKVersion []string
	AccountId       string
	DefaultRegion   string
	Region          string
	Endpoint        string
	Credentials     *Credentials
	HTTPClient      *http.Client

	Host          string            `json:"host,omitempty"`
	Scheme        string            `json:"scheme,omitempty"`
	DefaultHeader map[string]string `json:"defaultHeader,omitempty"`
	UserAgent     string            `json:"userAgent,omitempty"`
}

func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}

func (c *Configuration) SetupRequestHeader(path string, method string, request *http.Request) {
	request.Header[HeaderClientType] = []string{ClientType}
	request.Header[HeaderLanguage] = []string{"en-US"}
	request.Header[HeaderAccountId] = []string{c.AccountId}

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	if request.URL.RawQuery != "" {
		path = fmt.Sprintf("%s?%s", path, request.URL.RawQuery)
	}

	signature := MakeHmacSignature(method, path, timestamp, c.Credentials.AccessKey, c.Credentials.SecretKey, c.AccountId, ClientType)

	request.Header[HeaderAccessKey] = []string{c.Credentials.AccessKey}
	request.Header[HeaderTimestamp] = []string{timestamp}
	request.Header[HeaderSignature] = []string{signature}

	if c.Credentials.AuthToken != "" {
		request.Header[HeaderAuthToken] = []string{c.Credentials.AuthToken}
	}
}
