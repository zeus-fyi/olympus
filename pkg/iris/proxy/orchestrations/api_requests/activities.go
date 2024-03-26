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
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	hera_reddit "github.com/zeus-fyi/olympus/pkg/hera/reddit"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
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

func (i *IrisApiRequestsActivities) ExtLoadBalancerRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	if pr.Url == "" {
		err := fmt.Errorf("error: URL is required")
		log.Err(err).Msg("ExtLoadBalancerRequest: URL is required")
		return pr, err
	}
	var bearer string
	var user, pw string
	if pr.SecretNameRef != "" {
		ps, err := aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), org_users.NewOrgUserWithID(pr.OrgID, pr.UserID), pr.SecretNameRef)
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
	for k, v := range pr.RequestHeaders {
		switch k {
		case "Authorization":
			if !pr.IsInternal {
				continue
			}
		}
		r.SetHeader(k, strings.Join(v, ", ")) // Joining all values with a comma
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
		log.Err(err).Msg("ExtLoadBalancerRequest: Failed to relay api request")
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
			err = json.Unmarshal(resp.Body(), &pr.Response)
			if err != nil {
				log.Warn().Err(err).Msg("sendRequest: failed to unmarshal response body")
				tv := resp.String()
				tv = strings.TrimPrefix(tv, "\"")
				tv = strings.TrimSuffix(tv, "\"")
				tv = unescapeUnicode(tv)
				tv = html.UnescapeString(tv)
				//doc, err := goquery.NewDocumentFromReader(strings.NewReader(tv))
				//if err != nil {
				//	log.Err(err).Msg("sendRequest: failed to parse response body")
				//	return resp, err
				//}
				//// Find all meta tags and print their attributes
				//doc.Find("meta").Each(func(i int, s *goquery.Selection) {
				//	// For each item found, get the html
				//	name, _ := s.Attr("name")
				//	property, _ := s.Attr("property")
				//	content, _ := s.Attr("content")
				//	fmt.Printf("Meta tag found: name=%s, property=%s, content=%s\n", name, property, content)
				//})
				//pr.Response = echo.Map{
				//	"msg_body": mb,
				//}
				err = nil
			}
		}
		pr.PayloadSizeMeter.Add(resp.Size())
		pr.StatusCode = resp.StatusCode()
		if resp.StatusCode() >= 400 || pr.Response == nil {
			pr.RawResponse = resp.Body()
		}
		pr.StatusCode = resp.StatusCode()
		pr.ResponseHeaders = filterHeaders(resp.RawResponse.Header)
		pr.ReceivedAt = resp.ReceivedAt()
		pr.Latency = resp.Time()
		if pr.IsInternal {
			pr.RawResponse = resp.Body()
		}
		if pr.RegexFilters != nil {
			tv := resp.String()
			tv = strings.TrimPrefix(tv, "\"")
			tv = strings.TrimSuffix(tv, "\"")
			tv = unescapeUnicode(tv)
			tv = html.UnescapeString(tv)
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tv))
			if err != nil {
				log.Err(err).Msg("sendRequest: failed to parse response body")
				return resp, err
			}
			tmp, rerr := ExtractParams(pr.RegexFilters, []byte(doc.Text()))
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
		if key[:2] == "X-" {
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

func Unescape(value string, isBytes bool) (string, error) {
	// All strings normalize newlines to the \n representation.
	value = newlineNormalizer.Replace(value)
	n := len(value)

	// Nothing to unescape / decode.
	if n < 2 {
		return value, fmt.Errorf("unable to unescape string")
	}

	// Raw string preceded by the 'r|R' prefix.
	isRawLiteral := false
	if value[0] == 'r' || value[0] == 'R' {
		value = value[1:]
		n = len(value)
		isRawLiteral = true
	}

	// Quoted string of some form, must have same first and last char.
	if value[0] != value[n-1] || (value[0] != '"' && value[0] != '\'') {
		return value, fmt.Errorf("unable to unescape string")
	}

	// Normalize the multi-line CEL string representation to a standard
	// Go quoted string.
	if n >= 6 {
		if strings.HasPrefix(value, "'''") {
			if !strings.HasSuffix(value, "'''") {
				return value, fmt.Errorf("unable to unescape string")
			}
			value = "\"" + value[3:n-3] + "\""
		} else if strings.HasPrefix(value, `"""`) {
			if !strings.HasSuffix(value, `"""`) {
				return value, fmt.Errorf("unable to unescape string")
			}
			value = "\"" + value[3:n-3] + "\""
		}
		n = len(value)
	}
	value = value[1 : n-1]
	// If there is nothing to escape, then return.
	if isRawLiteral || !strings.ContainsRune(value, '\\') {
		return value, nil
	}

	// Otherwise the string contains escape characters.
	// The following logic is adapted from `strconv/quote.go`
	var runeTmp [utf8.UTFMax]byte
	buf := make([]byte, 0, 3*n/2)
	for len(value) > 0 {
		c, encode, rest, err := unescapeChar(value, isBytes)
		if err != nil {
			return "", err
		}
		value = rest
		if c < utf8.RuneSelf || !encode {
			buf = append(buf, byte(c))
		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf = append(buf, runeTmp[:n]...)
		}
	}
	return string(buf), nil
}

// unescapeChar takes a string input and returns the following info:
//
//	value - the escaped unicode rune at the front of the string.
//	encode - the value should be unicode-encoded
//	tail - the remainder of the input string.
//	err - error value, if the character could not be unescaped.
//
// When encode is true the return value may still fit within a single byte,
// but unicode encoding is attempted which is more expensive than when the
// value is known to self-represent as a single byte.
//
// If isBytes is set, unescape as a bytes literal so octal and hex escapes
// represent byte values, not unicode code points.
func unescapeChar(s string, isBytes bool) (value rune, encode bool, tail string, err error) {
	// 1. Character is not an escape sequence.
	switch c := s[0]; {
	case c >= utf8.RuneSelf:
		r, size := utf8.DecodeRuneInString(s)
		return r, true, s[size:], nil
	case c != '\\':
		return rune(s[0]), false, s[1:], nil
	}

	// 2. Last character is the start of an escape sequence.
	if len(s) <= 1 {
		err = fmt.Errorf("unable to unescape string, found '\\' as last character")
		return
	}

	c := s[1]
	s = s[2:]
	// 3. Common escape sequences shared with Google SQL
	switch c {
	case 'a':
		value = '\a'
	case 'b':
		value = '\b'
	case 'f':
		value = '\f'
	case 'n':
		value = '\n'
	case 'r':
		value = '\r'
	case 't':
		value = '\t'
	case 'v':
		value = '\v'
	case '\\':
		value = '\\'
	case '\'':
		value = '\''
	case '"':
		value = '"'
	case '`':
		value = '`'
	case '?':
		value = '?'

	// 4. Unicode escape sequences, reproduced from `strconv/quote.go`
	case 'x', 'X', 'u', 'U':
		n := 0
		encode = true
		switch c {
		case 'x', 'X':
			n = 2
			encode = !isBytes
		case 'u':
			n = 4
			if isBytes {
				err = fmt.Errorf("unable to unescape string")
				return
			}
		case 'U':
			n = 8
			if isBytes {
				err = fmt.Errorf("unable to unescape string")
				return
			}
		}
		var v rune
		if len(s) < n {
			err = fmt.Errorf("unable to unescape string")
			return
		}
		for j := 0; j < n; j++ {
			x, ok := unhex(s[j])
			if !ok {
				err = fmt.Errorf("unable to unescape string")
				return
			}
			v = v<<4 | x
		}
		s = s[n:]
		if !isBytes && v > utf8.MaxRune {
			err = fmt.Errorf("unable to unescape string")
			return
		}
		value = v

	// 5. Octal escape sequences, must be three digits \[0-3][0-7][0-7]
	case '0', '1', '2', '3':
		if len(s) < 2 {
			err = fmt.Errorf("unable to unescape octal sequence in string")
			return
		}
		v := rune(c - '0')
		for j := 0; j < 2; j++ {
			x := s[j]
			if x < '0' || x > '7' {
				err = fmt.Errorf("unable to unescape octal sequence in string")
				return
			}
			v = v*8 + rune(x-'0')
		}
		if !isBytes && v > utf8.MaxRune {
			err = fmt.Errorf("unable to unescape string")
			return
		}
		value = v
		s = s[2:]
		encode = !isBytes

		// Unknown escape sequence.
	default:
		err = fmt.Errorf("unable to unescape string")
	}

	tail = s
	return
}

func unhex(b byte) (rune, bool) {
	c := rune(b)
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

var (
	newlineNormalizer = strings.NewReplacer("\r\n", "\n", "\r", "\n")
)
