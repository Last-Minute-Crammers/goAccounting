package aiService

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
}

// 生成随机字符串
func genNonce(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// 生成标准查询字符串
func genCanonicalQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		escapedKey := url.QueryEscape(k)
		escapedValue := url.QueryEscape(params[k])
		parts = append(parts, fmt.Sprintf("%s=%s", escapedKey, escapedValue))
	}

	return strings.Join(parts, "&")
}

// 生成HMAC-SHA256签名
func genSignature(appSecret, signingString string) string {
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(signingString))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 生成认证头
func GenerateAuthHeaders(method, uri string, queryParams map[string]string, appID, appKey string) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := genNonce(8)
	canonicalQueryString := genCanonicalQueryString(queryParams)

	signedHeadersString := fmt.Sprintf("x-ai-gateway-app-id:%s\nx-ai-gateway-timestamp:%s\nx-ai-gateway-nonce:%s",
		appID, timestamp, nonce)

	signingString := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		strings.ToUpper(method),
		uri,
		canonicalQueryString,
		appID,
		timestamp,
		signedHeadersString)

	signature := genSignature(appKey, signingString)

	log.Printf("认证信息生成:")
	log.Printf("  Timestamp: %s", timestamp)
	log.Printf("  Nonce: %s", nonce)
	log.Printf("  Canonical Query: %s", canonicalQueryString)
	log.Printf("  Signed Headers: %s", signedHeadersString)
	log.Printf("  Signing String: %s", signingString)
	log.Printf("  Signature: %s", signature)

	return map[string]string{
		"X-AI-GATEWAY-APP-ID":         appID,
		"X-AI-GATEWAY-TIMESTAMP":      timestamp,
		"X-AI-GATEWAY-NONCE":          nonce,
		"X-AI-GATEWAY-SIGNED-HEADERS": "x-ai-gateway-app-id;x-ai-gateway-timestamp;x-ai-gateway-nonce",
		"X-AI-GATEWAY-SIGNATURE":      signature,
	}
}