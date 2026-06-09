package scpsdk

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	ClientType string = "Openapi"

	// HTTP request header parameters
	HeaderLanguage   string = "Accept-Language"
	HeaderAccessKey  string = "Scp-AccessKey"
	HeaderAccountId  string = "Scp-AccountId"
	HeaderSignature  string = "Scp-Signature"
	HeaderTimestamp  string = "Scp-Timestamp"
	HeaderClientType string = "Scp-ClientType"
	HeaderAuthToken  string = "X-Auth-Token"
)

func ConvertTimestampsToUTC(b []byte) ([]byte, error) {
	var result map[string]interface{}
	err := json.Unmarshal(b, &result)
	if err != nil {
		fmt.Println(err)
	}
	convertDateFields(result)
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return jsonBytes, nil
}

func convertDateFields(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if str, ok := value.(string); ok {
				parsedTime, err := time.Parse("2006-01-02T15:04:05", str)
				if err == nil {
					v[key] = parsedTime.UTC().Format(time.RFC3339)
				}
			}
			convertDateFields(value)
		}
	case []interface{}:
		for i := range v {
			convertDateFields(v[i])
		}
	}
}
