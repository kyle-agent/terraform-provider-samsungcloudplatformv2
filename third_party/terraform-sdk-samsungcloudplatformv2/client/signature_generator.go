package scpsdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

// MakeHmacSignature HMAC signature generator for SCP OpenAPI
func MakeHmacSignature(method string, requestUrl string, timestamp string, accessKey string, secretKey string, accountId string, clientType string) string {
	message := strings.Join([]string{method, requestUrl, timestamp, accessKey, accountId, clientType}, "")
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	raw := h.Sum(nil)
	base64Signature := base64.StdEncoding.EncodeToString(raw)
	return base64Signature
}
