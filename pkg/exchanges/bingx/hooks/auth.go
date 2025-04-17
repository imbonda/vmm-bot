package hooks

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/imbonda/vmm-bot/pkg/utils"
)

func GetSigAuthBeforeRequestHook(client *resty.Client, creds *utils.Credentials) resty.RequestMiddleware {
	return func(client *resty.Client, request *resty.Request) error {
		return authenticate(request, creds)
	}
}

func authenticate(request *resty.Request, creds *utils.Credentials) error {
	if request.Method != http.MethodPost {
		return nil
	}
	sign(request, creds)
	return nil
}

func sign(request *resty.Request, creds *utils.Credentials) {
	request.Header.Set("X-BX-APIKEY", creds.APIKey)
	timestamp := time.Now().UnixNano() / 1e6
	request.SetFormData(map[string]string{"timestamp": fmt.Sprint(timestamp)})
	signature := generateSignature(request.FormData, creds)
	request.SetFormData(map[string]string{"signature": signature})
}

func generateSignature(formData url.Values, creds *utils.Credentials) string {
	// Step 1: Sort the form data by key
	var keys []string
	for key := range formData {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Step 2: Create the query string from form data
	var queryString strings.Builder
	for i, key := range keys {
		values := formData[key]
		queryString.WriteString(fmt.Sprintf("%s=%s", key, strings.Join(values, ",")))
		if i < len(keys)-1 {
			queryString.WriteString("&")
		}
	}

	// Step 3: Generate signature
	sig := utils.HMAC256(queryString.String(), creds.APISecret)

	return sig
}
