package ai_platform_service_orchestrations

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-resty/resty/v2"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (t *ZeusWorkerTestSuite) TestApiCallRequestTask() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	act := NewZeusAiPlatformActivities()

	rets, err := act.SelectRetrievalTask(ctx, t.Ou, 1706767039731058000)
	t.Require().Nil(err)
	t.Require().NotEmpty(rets)
	ret := rets[0]
	t.Require().Equal(apiApproval, ret.RetrievalPlatform)
	t.Require().NotNil(ret.WebFilters)
	t.Require().NotNil(ret.WebFilters.RoutingGroup)
	tmp := "https://api.twitter.com/2/users/"
	r := RouteTask{
		Ou:        t.Ou,
		Retrieval: ret,
		RouteInfo: iris_models.RouteInfo{
			RoutePath: aws.ToString(&tmp),
		},
	}
	td, err := act.ApiCallRequestTask(ctx, r, nil)
	t.Require().Nil(err)
	t.Require().NotEmpty(td)
}

func (t *ZeusWorkerTestSuite) TestApiCallRequestTask2() {
	//urlv := "https://matokeholdings.com/"
	//
	//v, err := fetchURL(urlv)
	//t.Require().Nil(err)
	//t.Require().NotEmpty(v)

	tmm := "\"\\u003c!doctype html\\u003e\\n\\u003chtml lang=\\\"en-GB\\\"\\u003e\\n\\u003chead\\u003e\\u003cmeta charset=\\\"UTF-8\\\"\\u003e\\u003cscript\\u003eif(navigator.userAgent.match(/MSIE|Internet Explorer/i)||navigator.userAgent.match(/Trident\\\\/7\\\\..*?rv:11/i)){var href=document.location.href;if(!href.match(/[?\\u0026]nowprocket/)){if(href.indexOf(\\\"?\\\")==-1){if(href.indexOf(\\\"#\\\")==-1){d\"\n"

	//gg, err := Unescape(tmm, false)
	//t.Require().NotEmpty(gg)
	b, err := unescapeUnicodeSequences(tmm)
	t.Require().NotEmpty(b)
	t.Require().Nil(err)
}

func fetchURL(url string) (string, error) {
	client := resty.New()

	resp, err := client.R().Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %w", err)
	}

	// Convert body to a string
	bodyStr := string(resp.Body())

	// Unescape Unicode characters
	unescapedBody, err := unescapeUnicodeSequences(bodyStr)
	if err != nil {
		return "", fmt.Errorf("error unescaping Unicode: %w", err)
	}

	// Finally, unescape HTML entities
	return html.UnescapeString(unescapedBody), nil
}

// unescapeUnicode attempts to unquote a string that contains escaped Unicode.
// This is necessary because direct unquoting might fail if the string isn't properly quoted.
func unescapeUnicodeSequences(s string) (string, error) {
	// strconv.Unquote requires the string to be in double quotes.
	// Backslashes in the original string also need to be escaped.
	quoted := `"` + strings.ReplaceAll(s, `\`, `\\`) + `"`

	unescaped, err := strconv.Unquote(quoted)
	if err != nil {
		return "", err
	}

	return unescaped, nil
}
