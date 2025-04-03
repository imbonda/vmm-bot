package hooks

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/imbonda/bybit-vmm-bot/pkg/utils"
)

func GetSigAuthBeforeRequestHook(client *resty.Client, creds *utils.Credentials) resty.RequestMiddleware {
	return func(client *resty.Client, request *resty.Request) error {
		return addSignatureAuthentication(request, creds)
	}
}

func addSignatureAuthentication(request *resty.Request, creds *utils.Credentials) error {
	if request.Method != http.MethodPost {
		return nil
	}
	request.SetFormData(map[string]string{
		"api_key": creds.APIKey,
	})
	signature := generateSignature(request.FormData, creds)
	request.SetFormData(map[string]string{
		"sign": signature,
	})
	return nil
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
	for _, key := range keys {
		values := formData[key]
		queryString.WriteString(fmt.Sprintf("%s=%s&", key, strings.Join(values, ",")))
	}

	// Step 3: Append the secret_key to the query string
	queryString.WriteString(fmt.Sprintf("secret_key=%s", creds.APISecret))

	// Step 5: Generate signature
	sig := utils.MD5(queryString.String())

	// Step 6: Uppercase
	return strings.ToUpper(sig)
}
