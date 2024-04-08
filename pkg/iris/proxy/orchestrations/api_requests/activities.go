package iris_api_requests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	hera_reddit "github.com/zeus-fyi/olympus/pkg/hera/reddit"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

type IrisApiRequestsActivities struct {
}

func NewIrisApiRequestsActivities() IrisApiRequestsActivities {
	return IrisApiRequestsActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (i *IrisApiRequestsActivities) GetActivities() ActivitiesSlice {
	return []interface{}{i.RelayRequest, i.InternalSvcRelayRequest, i.ExtLoadBalancerRequest, i.UpdateOrgRoutingTable,
		i.SelectSingleOrgGroupsRoutingTables, i.SelectOrgGroupRoutingTable, i.SelectAllRoutingTables,
		i.DeleteOrgRoutingTable, i.ExtToAnvilInternalSimForkRequest,
	}
}

func (i *IrisApiRequestsActivities) RelayRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)
	if pr.IsInternal {
		r.SetAuthToken(artemis_orchestration_auth.Bearer)
	}
	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}
	resp, err := r.R().SetBody(&pr.Payload).SetResult(&pr.Response).Post(pr.Url)
	if err != nil {
		log.Err(err).Msg("IrisApiRequestsActivities: failed to relay api request")
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Msg("IrisApiRequestsActivities failed to relay api request")
		return nil, fmt.Errorf("failed to relay api request: status code %d", resp.StatusCode())
	}
	return pr, err
}

func (i *IrisApiRequestsActivities) InternalSvcRelayRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)
	if pr.IsInternal {
		r.SetAuthToken(artemis_orchestration_auth.Bearer)
	}
	resp, err := r.R().SetBody(&pr.Payload).SetResult(&pr.Response).Post(pr.Url)
	if err != nil {
		log.Err(err).Interface("statusCode", resp.StatusCode()).Msg("InternalSvcRelayRequest: failed to relay api request")
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("statusCode", resp.StatusCode()).Msg("Failed to relay api request")
		return nil, fmt.Errorf("InternalSvcRelayRequest: failed to relay api request: status code %d", resp.StatusCode())
	}
	return pr, err
}

func (i *IrisApiRequestsActivities) ExtToAnvilInternalSimForkRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	pr.IsInternal = true
	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}
	return i.ExtLoadBalancerRequest(ctx, pr)
}

const (
	flowsSecretsOrgID = 1710298581127603000
)

func (i *IrisApiRequestsActivities) ExtLoadBalancerRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	if pr.Url == "" {
		err := fmt.Errorf("error: URL is required")
		log.Err(err).Msg("ExtLoadBalancerRequest: URL is required")
		return pr, err
	}

	to := pr.OrgID
	if pr.IsFlowRequest && pr.SecretNameRef == "api-iris" {
		to = flowsSecretsOrgID
	}
	var bearer string
	var user, pw string
	if pr.SecretNameRef != "" {
		ps, err := aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), org_users.NewOrgUserWithID(to, pr.UserID), pr.SecretNameRef)
		if ps != nil && ps.BearerToken != "" {
			bearer = ps.BearerToken
		} else if err != nil {
			log.Err(err).Msg("ProcessRpcLoadBalancerRequest: failed to get mockingbird secrets")
			return pr, err
		}
		if len(ps.TwillioAccount) > 0 {
			user = ps.TwillioAccount
		}
		if len(ps.TwillioAuth) > 0 {
			pw = ps.TwillioAuth
		}
		if strings.HasPrefix(pr.Url, "https://oauth.reddit.com") {
			ps, err = aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), org_users.NewOrgUserWithID(pr.OrgID, pr.UserID), "reddit")
			if err != nil {
				log.Err(err).Msg("ProcessRpcLoadBalancerRequest: failed to get mockingbird secrets")
				return pr, err
			}
			rc, rrr := hera_reddit.InitOrgRedditClient(ctx, ps.OAuth2Public, ps.OAuth2Secret, ps.Username, ps.Password)
			if rrr != nil {
				log.Err(rrr).Msg("ProcessRpcLoadBalancerRequest: failed to get reddit client")
				return pr, rrr
			}
			switch pr.PayloadTypeREST {
			case "GET", "get":
				resp, rerr := rc.GetRedditReq(ctx, pr.ExtRoutePath, &pr.Response, pr.QueryParams)
				if rerr != nil {
					log.Err(rerr).Interface("resp", resp).Msg("ProcessRpcLoadBalancerRequest: failed to get reddit request")
					return pr, rerr
				}
				return pr, nil
			case "POST", "post":
				resp, rerr := rc.PostRedditReq(ctx, pr.ExtRoutePath, pr.Payload, &pr.Response)
				if rerr != nil {
					log.Err(rerr).Interface("resp", resp).Msg("ProcessRpcLoadBalancerRequest: failed to post reddit request")
					return pr, rerr
				}
				return pr, nil
			case "PUT", "put":
				resp, rerr := rc.PutRedditReq(ctx, pr.ExtRoutePath, pr.Payload, &pr.Response)
				if rerr != nil {
					log.Err(rerr).Interface("resp", resp).Msg("ProcessRpcLoadBalancerRequest: failed to put reddit request")
					return pr, rerr
				}
				return pr, nil
			case "DELETE", "delete":
				resp, rerr := rc.DeleteRedditReq(ctx, pr.ExtRoutePath, pr.Payload, &pr.Response)
				if rerr != nil {
					log.Err(rerr).Interface("resp", resp).Msg("ProcessRpcLoadBalancerRequest: failed to delete reddit request")
					return pr, rerr
				}
				return pr, nil
			default:
				return pr, errors.New("ProcessRpcLoadBalancerRequest: invalid payload type for supported reddit requests")
			}
		}
	}

	var r *resty.Client
	r = resty.New()
	r.SetBaseURL(pr.Url)
	if pr.MaxTries > 0 {
		r.SetRetryCount(pr.MaxTries)
	}
	if len(bearer) > 0 {
		r.SetAuthToken(bearer)
		log.Info().Msg("ExtLoadBalancerRequest: setting bearer token")
	}
	if len(user) > 0 && len(pw) > 0 {
		r.SetBasicAuth(user, pw)
		log.Info().Msg("ExtLoadBalancerRequest: setting basic auth")
	}
	parsedURL, err := url.Parse(pr.Url)
	if err != nil {
		log.Err(err).Msg("ExtLoadBalancerRequest: failed to parse url")
		return pr, err
	}

	if pr.OrgID == 7138983863666903883 {
		// for internal
	} else if pr.IsInternal {
		log.Info().Interface("pr.URL", pr.Url).Msg("ExtLoadBalancerRequest: anvil request")
	} else {
		if parsedURL.Scheme != "https" {
			err = fmt.Errorf("error: URL must be an HTTPS URL")
			log.Err(err).Msg("ExtLoadBalancerRequest: http request unauthorized")
			return pr, err
		}
	}

	if len(pr.QueryParams) > 0 {
		r.QueryParam = pr.QueryParams
	}
	if pr.RequestHeaders != nil {
		for k, v := range pr.RequestHeaders {
			switch k {
			case "Authorization":
				if !pr.IsInternal {
					continue
				}
			}
			r.SetHeader(k, strings.Join(v, ", ")) // Joining all values with a comma
		}
	}
	if len(pr.Referrers) > 0 {
		r.SetHeader("Referer", strings.Join(pr.Referrers, ", ")) // Joining all values with a comma
	}

	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}

	var resp *resty.Response
	switch pr.PayloadTypeREST {
	case "GET", "get":
		resp, err = sendRequest(r.R(), pr, "GET")
	case "PUT", "put":
		resp, err = sendRequest(r.R(), pr, "PUT")
	case "DELETE", "delete":
		resp, err = sendRequest(r.R(), pr, "DELETE")
	case "POST", "post":
		resp, err = sendRequest(r.R(), pr, "POST")
	case "OPTIONS", "options":
		resp, err = sendRequest(r.R(), pr, "OPTIONS")
	default:
		resp, err = sendRequest(r.R(), pr, "POST")
	}
	if err != nil {
		log.Err(err).Interface("pr.Url", pr.Url).Interface("pr.ExtRoutePath", pr.ExtRoutePath).Msg("ExtLoadBalancerRequest: Failed to relay api request")
		return pr, fmt.Errorf("failed to relay api request")
	}
	if pr.StatusCode >= 400 && !skipErrorOnStatusCodes(pr.StatusCode, pr.SkipErrorOnStatusCodes) {
		log.Err(err).Msg("IrisApiRequestsActivities: failed to relay api request")
		return pr, fmt.Errorf("failed to relay api request: status code %d", resp.StatusCode())
	}
	return pr, err
}

func skipErrorOnStatusCodes(statusCode int, skipErrorOnStatusCodes []int) bool {
	for _, code := range skipErrorOnStatusCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}

func sendRequest(request *resty.Request, pr *ApiProxyRequest, method string) (*resty.Response, error) {
	var resp *resty.Response
	var err error

	ext := ""
	if pr.ExtRoutePath != "" {
		ext = pr.ExtRoutePath
	}
	if len(pr.Payload) == 0 {
		pr.Payload = nil
	}

	if pr.Payload != nil || pr.Payloads != nil {
		if pr.Payloads != nil && pr.Payload == nil {
			switch method {
			case "GET", "get":
				resp, err = request.SetBody(&pr.Payloads).Get(ext)
			case "OPTIONS", "options":
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Options(ext)
			case "PUT", "put":
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Put(ext)
			case "DELETE", "delete":
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Delete(ext)
			case "POST-form", "post-form":
				fbp := make(map[string]string)
				for k, v := range pr.Payload {
					fbp[k] = fmt.Sprintf("%v", v)
				}
				resp, err = request.SetFormData(fbp).SetResult(&pr.Response).Post(ext)
			case "POST", "post":
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Post(ext)
			default:
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Post(ext)
			}
		} else {
			switch method {
			case "GET", "get":
				resp, err = request.SetBody(&pr.Payload).Get(ext)
			case "OPTIONS", "options":
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Options(ext)
			case "PUT", "put":
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Put(ext)
			case "DELETE", "delete":
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Delete(ext)
			case "POST", "post":
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Post(ext)
			default:
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Post(ext)
			}
		}
	} else {
		switch method {
		case "OPTIONS", "options":
			resp, err = request.SetResult(&pr.Response).Options(ext)
		case "GET", "get":
			resp, err = request.Get(ext)
		case "PUT", "put":
			resp, err = request.SetResult(&pr.Response).Put(ext)
		case "DELETE", "delete":
			resp, err = request.SetResult(&pr.Response).Delete(ext)
		case "POST", "post":
			resp, err = request.SetResult(&pr.Response).Post(ext)
		default:
			resp, err = request.SetResult(&pr.Response).Post(ext)
		}
	}
	if err != nil {
		if resp != nil {
			log.Err(err).Int("statusCode", resp.StatusCode()).Interface("resp", resp.String()).Interface("url", pr.Url).Interface("pr.ExtRoutePath", ext).Msg("sendRequest: failed to relay api request")
		} else {
			log.Err(err).Interface("url", pr.Url).Interface("pr.ExtRoutePath", ext).Msg("sendRequest: failed to relay api request")
		}
	}

	if resp != nil {
		if pr.PayloadSizeMeter == nil {
			pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
		}
		if method == "GET" {
			if pr.RequestHeaders != nil && pr.RequestHeaders["X-Scrape-Html"] != nil {
				tv := resp.String()
				tv = strings.TrimPrefix(tv, "\"")
				tv = strings.TrimSuffix(tv, "\"")
				tv = unescapeUnicode(tv)
				tv = html.UnescapeString(tv)
				doc, herr := goquery.NewDocumentFromReader(strings.NewReader(tv))
				if herr != nil {
					log.Err(herr).Msg("sendRequest: failed to parse response body")
					return resp, herr
				}
				plMap, perr := extractAndRespond(doc)
				if perr != nil {
					return nil, perr
				}
				pr.Response = plMap
			}
			if pr.Response == nil && resp.Body() != nil {
				err = json.Unmarshal(resp.Body(), &pr.Response)
				if err != nil {
					log.Warn().Err(err).Msg("sendRequest: failed to unmarshal response body")
					err = nil
				}
			}
		} else if method == "POST" && pr.RequestHeaders["X-Scrape-Html"] != nil {
			tv := resp.String()
			tv = strings.TrimPrefix(tv, "\"")
			tv = strings.TrimSuffix(tv, "\"")
			tv = unescapeUnicode(tv)
			tv = html.UnescapeString(tv)
			doc, herr := goquery.NewDocumentFromReader(strings.NewReader(tv))
			if herr != nil {
				log.Err(herr).Msg("sendRequest: failed to parse response body")
				return resp, herr
			}
			plMap, rerr := extractAndRespond(doc)
			if rerr != nil {
				log.Err(rerr).Msg("sendRequest: extractAndRespond")
				return resp, rerr
			}
			pr.Response = plMap
		}
		pr.PayloadSizeMeter.Add(resp.Size())
		pr.StatusCode = resp.StatusCode()
		if resp.StatusCode() >= 400 || pr.Response == nil {
			pr.RawResponse = resp.Body()
		}
		pr.StatusCode = resp.StatusCode()
		if pr.RequestHeaders != nil && resp != nil && resp.RawResponse != nil && resp.RawResponse.Header != nil {
			pr.ResponseHeaders = filterHeaders(resp.RawResponse.Header)
		}
		pr.ReceivedAt = resp.ReceivedAt()
		pr.Latency = resp.Time()
		if pr.IsInternal {
			pr.RawResponse = resp.Body()
		}
		if pr.RegexFilters != nil {
			br := resp.Body()
			if pr.RequestHeaders != nil && pr.RequestHeaders["X-Scrape-Html"] != nil {
				tv := resp.String()
				tv = strings.TrimPrefix(tv, "\"")
				tv = strings.TrimSuffix(tv, "\"")
				tv = unescapeUnicode(tv)
				tv = html.UnescapeString(tv)
				doc, herr := goquery.NewDocumentFromReader(strings.NewReader(tv))
				if herr != nil {
					log.Err(herr).Msg("sendRequest: failed to parse response body")
					return resp, herr
				}
				pr.Response, err = extractAndRespond(doc)
				if err != nil {
					return nil, err
				}
				brs, berr := json.Marshal(pr.Response)
				if berr != nil {
					log.Err(err).Msg("sendRequest: failed to marshal response body")
				} else {
					br = brs
				}
			}
			tmp, rerr := ExtractParams(pr.RegexFilters, br)
			if rerr != nil {
				log.Err(rerr).Msg("sendRequest: failed to extract params")
				return resp, rerr
			}
			pr.RawResponse = []byte(strings.Join(tmp, ", "))
		}
	}
	return resp, nil
}

func unescapeUnicode(input string) string {
	// Create a new strings.Builder for efficient string concatenation
	var builder strings.Builder

	// Iterate over each rune in the input string
	for i := 0; i < len(input); {
		// Check if the current substring matches the escape sequences and replace them
		if i+6 <= len(input) && (input[i:i+6] == "\\u003c" || input[i:i+6] == "\\u003e") {
			if input[i:i+6] == "\\u003c" {
				builder.WriteString("<")
			} else if input[i:i+6] == "\\u003e" {
				builder.WriteString(">")
			}
			i += 6 // Skip past the escape sequence
		} else {
			// Write the current rune to the builder and move to the next rune
			builder.WriteRune(rune(input[i]))
			i++
		}
	}

	// Return the constructed string
	return builder.String()
}

func filterHeaders(headers http.Header) http.Header {
	filteredHeaders := make(http.Header)
	for key, values := range headers {
		if len(key) > 2 && key[:2] == "X-" {
			filteredHeaders[key] = values
		}
	}
	return filteredHeaders
}
func ExtractParams(regexStrs []string, strContent []byte) ([]string, error) {
	// Split the regexStr into individual patterns
	var combinedParams []string
	for _, pattern := range regexStrs {
		// Compile and execute each pattern
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err // Return an error if the regular expression compilation fails
		}
		// Find all matches and extract the parameter names
		matches := re.FindAll(strContent, -1)
		for _, match := range matches {
			combinedParams = append(combinedParams, string(match))
		}
	}
	return combinedParams, nil
}

func extractAndRespond(doc *goquery.Document) (echo.Map, error) {
	var elements []string

	//// Extract meta tags
	//doc.Find("meta").Each(func(i int, s *goquery.Selection) {
	//	element := make(map[string]string)
	//	if name, exists := s.Attr("name"); exists {
	//		element["name"] = name
	//	}
	//	if property, exists := s.Attr("property"); exists {
	//		element["property"] = property
	//	}
	//	if content, exists := s.Attr("content"); exists {
	//		element["content"] = content
	//	}
	//	element["type"] = "meta"
	//	elements = append(elements, element)
	//})

	// Extract h1-h6 tags
	for _, tag := range []string{"h1", "h2", "h3", "h4", "h5", "h6"} {
		doc.Find(tag).Each(func(i int, s *goquery.Selection) {
			tv := cleanText(s.Text())
			if len(tv) == 0 {
				return
			}
			elements = append(elements, tv)
		})
	}

	// Extract p tags
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		tv := cleanText(s.Text())
		if len(tv) == 0 {
			return
		}
		elements = append(elements, tv)
	})

	if len(elements) == 0 {
		log.Info().Interface("doc", doc.Text())
		return nil, nil
	}

	el, err := LimitTokenResp(elements)
	if err != nil {
		log.Err(err).Interface("elements", elements).Msg("LimitTokenResp")
		return nil, err
	}
	return echo.Map{
		"msg_body": el,
	}, nil
}

func LimitTokenResp(elements []string) ([]string, error) {
	lv := 0

	totalCnt := 0
	for i, item := range elements {
		tokenCount, err := GetTokenCountEstimate(context.Background(), "gpt-3.5", item)
		if err != nil {
			return nil, err
		}
		totalCnt += tokenCount
		if totalCnt <= 9000 {
			lv = i // Update last valid if within the token limit
		} else {
			log.Info().Interface("totalCnt", totalCnt).Msg("LimitTokenResp")
			// Return the last valid range if current range exceeds 10000 tokens
			return elements[:lv], nil
		}
	}
	// If the whole array is under the limit, return it
	return elements[:lv], nil
}

type TokenCountsEstimate struct {
	Count int `json:"count"`
}

func GetTokenCountEstimate(ctx context.Context, model, text string) (int, error) {
	if len(model) == 0 {
		model = "gpt-4"
	}
	if strings.HasPrefix(model, "gpt-4") {
		model = "gpt-4"
	}
	if strings.HasPrefix(model, "gpt-3.5") {
		model = "gpt-3.5-turbo"
	}
	apiReq := &ApiProxyRequest{
		Url:             "https://pandora.zeus.fyi",
		PayloadTypeREST: "POST",
		Payload: echo.Map{
			"model": model,
			"text":  text,
		},
		IsInternal: true,
	}
	var tc TokenCountsEstimate
	res := resty_base.GetBaseRestyClient(apiReq.Url, artemis_orchestration_auth.Bearer)
	resp, err := res.R().SetBody(&apiReq.Payload).SetResult(&tc).Post("tokenize")
	if err != nil {
		// .Interface("&apiReq.Payload)", &apiReq.Payload)
		log.Err(err).Msg("Zeus: GetTokenCountEstimate")
		return -1, err
	}
	if resp != nil && resp.StatusCode() >= 400 {
		if err != nil {
			err = fmt.Errorf("GetTokenCountEstimate: failed to relay api request: status code %d", resp.StatusCode())
		}
		// .Interface("&apiReq.Payload)", &apiReq.Payload)
		//log.Err(err).Msg("Zeus: GetTokenCountEstimate")
		return -1, err
	}
	return tc.Count, nil
}

// cleanText trims spaces and removes escape characters from the text
func cleanText(text string) string {
	// Trim leading and trailing whitespace
	text = strings.TrimSpace(text)
	// Remove \r and \n escape sequences
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\t", "")
	// Remove additional internal spaces
	text = strings.Join(strings.Fields(text), " ")
	return text
}
