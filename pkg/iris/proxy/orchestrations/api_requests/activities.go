package iris_api_requests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

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
	if pr.SecretNameRef != "" {
		ps, err := aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), org_users.NewOrgUserWithID(pr.OrgID, pr.UserID), pr.SecretNameRef)
		if ps != nil && ps.BearerToken != "" {
			bearer = ps.BearerToken
		} else if err != nil {
			log.Err(err).Msg("ProcessRpcLoadBalancerRequest: failed to get mockingbird secrets")
			return pr, err
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
			case "GET":
				resp, rerr := rc.GetRedditReq(ctx, pr.ExtRoutePath, &pr.Response, pr.QueryParams)
				if rerr != nil {
					log.Err(rerr).Interface("resp", resp).Msg("ProcessRpcLoadBalancerRequest: failed to get reddit request")
					return pr, rerr
				}
				return pr, nil
			case "POST":
				resp, rerr := rc.PostRedditReq(ctx, pr.ExtRoutePath, pr.Payload, &pr.Response)
				if rerr != nil {
					log.Err(rerr).Interface("resp", resp).Msg("ProcessRpcLoadBalancerRequest: failed to post reddit request")
					return pr, rerr
				}
				return pr, nil
			case "PUT":
				resp, rerr := rc.PutRedditReq(ctx, pr.ExtRoutePath, pr.Payload, &pr.Response)
				if rerr != nil {
					log.Err(rerr).Interface("resp", resp).Msg("ProcessRpcLoadBalancerRequest: failed to put reddit request")
					return pr, rerr
				}
				return pr, nil
			case "DELETE":
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
	case "GET":
		resp, err = sendRequest(r.R(), pr, "GET")
	case "PUT":
		resp, err = sendRequest(r.R(), pr, "PUT")
	case "DELETE":
		resp, err = sendRequest(r.R(), pr, "DELETE")
	case "POST":
		resp, err = sendRequest(r.R(), pr, "POST")
	case "OPTIONS":
		resp, err = sendRequest(r.R(), pr, "OPTIONS")
	default:
		resp, err = sendRequest(r.R(), pr, "POST")
	}
	if err != nil {
		log.Err(err).Msg("ExtLoadBalancerRequest: Failed to relay api request")
		return pr, fmt.Errorf("failed to relay api request")
	}
	if pr.StatusCode >= 400 {
		log.Err(err).Msg("IrisApiRequestsActivities: failed to relay api request")
		return pr, fmt.Errorf("failed to relay api request: status code %d", resp.StatusCode())
	}
	return pr, err
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
			case "GET":
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Get(ext)
			case "OPTIONS":
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Options(ext)
			case "PUT":
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Put(ext)
			case "DELETE":
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Delete(ext)
			case "POST":
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Post(ext)
			default:
				resp, err = request.SetBody(&pr.Payloads).SetResult(&pr.Response).Post(ext)
			}
		} else {
			switch method {
			case "GET":
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Get(ext)
			case "OPTIONS":
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Options(ext)
			case "PUT":
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Put(ext)
			case "DELETE":
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Delete(ext)
			case "POST":
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Post(ext)
			default:
				resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Post(ext)
			}
		}
	} else {
		switch method {
		case "OPTIONS":
			resp, err = request.SetResult(&pr.Response).Options(ext)
		case "GET":
			resp, err = request.SetResult(&pr.Response).Get(ext)
		case "PUT":
			resp, err = request.SetResult(&pr.Response).Put(ext)
		case "DELETE":
			resp, err = request.SetResult(&pr.Response).Delete(ext)
		case "POST":
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
			tmp, rerr := ExtractParams(pr.RegexFilters, resp.Body())
			if rerr != nil {
				log.Err(rerr).Msg("sendRequest: failed to extract params")
				return resp, rerr
			}
			pr.RawResponse = []byte(strings.Join(tmp, ", "))
		}

	}
	return resp, nil
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
