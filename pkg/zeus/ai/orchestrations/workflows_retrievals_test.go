package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
)

func (t *ZeusWorkerTestSuite) TestRetrievalsWorkflowTask() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	act := NewZeusAiPlatformActivities()
	rets, err := act.SelectRetrievalTask(ctx, t.Ou, 1706767039731058000)
	t.Require().Nil(err)
	t.Require().NotEmpty(rets)
	ret := rets[0]

	ret.RetrievalPlatform = apiApproval
	cp := &MbChildSubProcessParams{
		WfID:         uuid.New().String(),
		Ou:           t.Ou,
		WfExecParams: artemis_orchestrations.WorkflowExecParams{},
		Oj: artemis_orchestrations.OrchestrationJob{Orchestrations: artemis_autogen_bases.Orchestrations{
			OrchestrationID: 1706767039731058000,
		}},
		Window: artemis_orchestrations.Window{},
		Wsr: artemis_orchestrations.WorkflowStageReference{
			WorkflowRunID: 1704069081079680000,
			ChildWfID:     "TestRetrievalsWorkflow-" + uuid.New().String(),
		},
		Tc: TaskContext{
			TaskID:    1706842030247904000,
			Retrieval: ret,
		},
	}

	cp, err = ZeusAiPlatformWorker.ExecuteRetrievalsWorkflow(ctx, cp)
	t.Require().Nil(err)
	t.Require().NotZero(cp.Wsr.InputID)
}

func (t *ZeusWorkerTestSuite) TestRetrievalsWorkflow() {
	t.initWorker()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	act := NewZeusAiPlatformActivities()
	rets, err := act.SelectRetrievalTask(ctx, t.Ou, 1706487709357339000)
	t.Require().Nil(err)
	t.Require().NotEmpty(rets)
	ret := rets[0]

	ret.RetrievalPlatform = apiApproval

	//t.Require().Equal(apiApproval, ret.RetrievalPlatform)
	t.Require().NotNil(ret.WebFilters)
	t.Require().NotNil(ret.WebFilters.RoutingGroup)

	cp := &MbChildSubProcessParams{
		WfID:         uuid.New().String(),
		Ou:           t.Ou,
		WfExecParams: artemis_orchestrations.WorkflowExecParams{},
		Oj:           artemis_orchestrations.OrchestrationJob{},
		Window:       artemis_orchestrations.Window{},
		Wsr: artemis_orchestrations.WorkflowStageReference{
			WorkflowRunID:  0,
			ChildWfID:      "TestRetrievalsWorkflow-" + uuid.New().String(),
			RunCycle:       0,
			IterationCount: 0,
			ChunkOffset:    0,
		},
		Tc: TaskContext{
			TaskID:    0,
			EvalID:    1704066747085827000,
			Retrieval: ret,
			TriggerActionsApproval: artemis_orchestrations.TriggerActionsApproval{
				TriggerAction:    apiApproval,
				ApprovalID:       1706566091973007000,
				EvalID:           1704066747085827000,
				TriggerID:        1706487755984811000,
				WorkflowResultID: 1706421224945827000,
			},
			AIWorkflowTriggerResultApiResponse: artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{
				ResponseID:  1706566091973014000,
				ApprovalID:  1706566091973007000,
				TriggerID:   1706487755984811000,
				RetrievalID: 1706487709357339000,
				ReqPayloads: []echo.Map{
					{
						"key1": "value1",
					},
					{
						"key2": "value2",
					},
				},
			},
		},
	}

	_, err = ZeusAiPlatformWorker.ExecuteRetrievalsWorkflow(ctx, cp)
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestRetrievalsExtract() {
	m := echo.Map{
		"q":        "best book",
		"category": "fiction books",
	}

	route := "customsearch/h/v1?q={q}&category={category}"
	ps, _, err := ReplaceAndPassParams(route, m)
	t.Require().Nil(err)

	expected := "customsearch/h/v1?q=best+book&category=fiction+books" // Expect spaces to be replaced with '+'
	t.Require().Equal(expected, ps)
	fmt.Println(ps)

}

//func (t *ZeusWorkerTestSuite) TestRetrievalsExtractStrReg() {
//	fp := filepaths.Path{
//		PackageName: "",
//		DirIn:       "/Users/alex/PycharmProjects/scratchPad/scrape",
//		DirOut:      "",
//		FnIn:        "google_search2.txt",
//		FnOut:       "",
//		Env:         "",
//		FilterFiles: nil,
//	}
//	b := fp.ReadFileInPath()
//	act := NewZeusAiPlatformActivities()
//	rets, err := act.SelectRetrievalTask(ctx, t.Ou, 1708579569890359000)
//	t.Require().Nil(err)
//	t.Require().NotEmpty(rets)
//	ret := rets[0]
//	t.Require().NotNil(ret.WebFilters)
//	t.Require().NotNil(ret.WebFilters.RegexPatterns)
//	for ri, rp := range ret.WebFilters.RegexPatterns {
//		fmt.Println("RegexPattern:", rp, "ind", ri)
//		ret.WebFilters.RegexPatterns[ri] = FixRegexInput(rp)
//	}
//	params, perr := ExtractParams(ret.WebFilters.RegexPatterns, b)
//	t.Require().Nil(perr)
//	t.Require().NotEmpty(params)
//	fmt.Println("Extracted parameters:", strings.Join(params, ", "))
//}
//
//func (t *ZeusWorkerTestSuite) TestRetrievalsExtractStrReg2() {
//	fp := filepaths.Path{
//		PackageName: "",
//		DirIn:       "/Users/alex/PycharmProjects/scratchPad/scrape",
//		DirOut:      "",
//		FnIn:        "google_search.txt",
//		FnOut:       "",
//		Env:         "",
//		FilterFiles: nil,
//	}
//	b := fp.ReadFileInPath()
//
//	params, err := ExtractParams([]string{`"https?:\/\/[^\"]+"`}, b)
//	t.Require().Nil(err)
//	t.Require().NotEmpty(params)
//	fmt.Println("Extracted parameters:", strings.Join(params, ", "))
//}

func (t *ZeusWorkerTestSuite) TestRetrievalsExtractStrReg3() {
	ep := "customsearch/v1?q={q}&cx=sdffs"

	pl := echo.Map{
		"q": "Alex George Zeusfyi",
	}
	pp, qp, err := ReplaceAndPassParams(ep, pl)
	t.Require().Nil(err)
	t.Require().NotEmpty(pp)
	t.Require().NotEmpty(qp)
}

const ex = `
{
  "context": {
    "title": "zeusfyi"
  },
  "items": [
    {
      "cacheId": "9bHXwQFfbF8J",
      "displayLink": "www.zeus.fyi",
      "formattedUrl": "https://www.zeus.fyi/",
      "htmlFormattedUrl": "https://www.zeus.fyi/",
      "htmlSnippet": "<b>Zeusfyi</b> handles cloud resource provisioning, app deployment, and configuration even with complex requirements like GPUs, NVMe setups.",
      "htmlTitle": "<b>Zeusfyi</b> | The Cloud Platform for Microservices &amp; SOA Development",
      "kind": "customsearch#result",
      "link": "https://www.zeus.fyi/",
      "pagemap": {
        "cse_image": [
          {
            "src": "https://static.wixstatic.com/media/c837a6_786cc680f77b40d19fc7bfd75f7a49e0~mv2.png/v1/fill/w_640,h_400,al_b,q_85,usm_0.66_1.00_0.01,enc_auto/c837a6_786cc680f77b40d19fc7bfd75f7a49e0~mv2.png"
          }
        ],
        "cse_thumbnail": [
          {
            "height": "177",
            "src": "https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcRs2-21Bx0qEanPae8W6lSgJ7Yjo9QdWoJ6RXTvoM-0S7tPuHLc_MEsXeAI",
            "width": "284"
          }
        ],
        "metatags": [
          {
            "format-detection": "telephone=no",
            "og:description": "Zeusfyi handles cloud resource provisioning, app deployment, and configuration even with complex requirements like GPUs, NVMe setups. We even help you architect the solution you need. PaaS",
            "og:site_name": "Zeusfyi",
            "og:title": "Zeusfyi | The Cloud Platform for Microservices & SOA Development",
            "og:type": "website",
            "og:url": "https://www.zeus.fyi",
            "skype_toolbar": "skype_toolbar_parser_compatible",
            "twitter:card": "summary_large_image",
            "twitter:description": "Zeusfyi handles cloud resource provisioning, app deployment, and configuration even with complex requirements like GPUs, NVMe setups. We even help you architect the solution you need. PaaS",
            "twitter:title": "Zeusfyi | The Cloud Platform for Microservices & SOA Development",
            "viewport": "width=320, user-scalable=yes"
          }
        ]
      },
      "snippet": "Zeusfyi handles cloud resource provisioning, app deployment, and configuration even with complex requirements like GPUs, NVMe setups.",
      "title": "Zeusfyi | The Cloud Platform for Microservices & SOA Development"
    },
    {
      "cacheId": "RIT81T8w3xsJ",
      "displayLink": "medium.zeus.fyi",
      "formattedUrl": "https://medium.zeus.fyi/",
      "htmlFormattedUrl": "https://medium.zeus.fyi/",
      "htmlSnippet": "We develop cutting edge cloud infra technology in the Kubernetes space. Easier than heroku, compatible with all cloud providers, hybrid, and on-premises.",
      "htmlTitle": "<b>Zeusfyi</b>",
      "kind": "customsearch#result",
      "link": "https://medium.zeus.fyi/",
      "pagemap": {
        "cse_image": [
          {
            "src": "https://cdn-images-1.medium.com/max/1200/1*HYX1fSo-KLX9RKNvgBFj1A.png"
          }
        ],
        "cse_thumbnail": [
          {
            "height": "196",
            "src": "https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcSaK90DSuwXzTWQO3Rq3NFRDWPOo5FWjlomqpYgXvJqnsJfsBCaod1uzKU",
            "width": "202"
          }
        ],
        "metatags": [
          {
            "al:android:app_name": "Medium",
            "al:android:package": "com.medium.reader",
            "al:android:url": "medium://zeusfyi",
            "al:ios:app_name": "Medium",
            "al:ios:app_store_id": "828256236",
            "al:ios:url": "medium://zeusfyi",
            "al:web:url": "https://medium.zeus.fyi/",
            "fb:app_id": "542599432471018",
            "medium-com:creator": "https://medium.zeus.fyi/@zeusfyi",
            "og:description": "We develop cutting edge cloud infra technology in the Kubernetes space. Easier than heroku, compatible with all cloud providers, hybrid, and on-premises. https://linktr.ee/zeusfyi.",
            "og:image": "https://cdn-images-1.medium.com/max/1200/1*HYX1fSo-KLX9RKNvgBFj1A.png",
            "og:site_name": "Zeusfyi",
            "og:title": "Zeusfyi",
            "og:type": "website",
            "og:url": "https://medium.zeus.fyi/",
            "referrer": "always",
            "theme-color": "#000000",
            "title": "Zeusfyi",
            "twitter:app:id:iphone": "828256236",
            "twitter:app:name:iphone": "Medium",
            "twitter:app:url:iphone": "medium://zeusfyi",
            "twitter:card": "summary_large_image",
            "twitter:creator": "@ctrl_alt_lulz",
            "twitter:description": "We develop cutting edge cloud infra technology in the Kubernetes space. Easier than heroku, compatible with all cloud providers, hybrid, and on-premises. https://linktr.ee/zeusfyi.",
            "twitter:image:src": "https://cdn-images-1.medium.com/max/1200/1*HYX1fSo-KLX9RKNvgBFj1A.png",
            "twitter:site": "@zeus_fyi",
            "twitter:title": "Zeusfyi",
            "viewport": "width=device-width, initial-scale=1.0, viewport-fit=contain"
          }
        ]
      },
      "snippet": "We develop cutting edge cloud infra technology in the Kubernetes space. Easier than heroku, compatible with all cloud providers, hybrid, and on-premises.",
      "title": "Zeusfyi"
    },
    {
      "cacheId": "HSh_jXV1FU4J",
      "displayLink": "github.com",
      "formattedUrl": "https://github.com/zeus-fyi/zeus",
      "htmlFormattedUrl": "https://github.com/zeus-fyi/zeus",
      "htmlSnippet": "Overview &middot; Automates translation of kubernetes yaml configurations into representative SQL models &middot; Users upload these infrastructure configurations via API&nbsp;...",
      "htmlTitle": "zeus-fyi/zeus: Zeus + SciFi = the power of the gods, meets ... - GitHub",
      "kind": "customsearch#result",
      "link": "https://github.com/zeus-fyi/zeus",
      "pagemap": {
        "cse_image": [
          {
            "src": "https://opengraph.githubassets.com/1ba1f8b2f54d32b988e84a9cc1423130b904f6b82bce864fdac78d44decd7971/zeus-fyi/zeus"
          }
        ],
        "cse_thumbnail": [
          {
            "height": "75",
            "src": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTT_HUpJOEpp6ovshvaQDmF11G-exXUoFVKnIkS692IhUyIzX8c3puqPQ",
            "width": "150"
          }
        ],
        "metatags": [
          {
            "analytics-location": "/<user-name>/<repo-name>",
            "apple-itunes-app": "app-id=1477376905, app-argument=https://github.com/zeus-fyi/zeus",
            "browser-errors-url": "https://api.github.com/_private/browser/errors",
            "browser-stats-url": "https://api.github.com/_private/browser/stats",
            "color-scheme": "light dark",
            "current-catalog-service-hash": "82c569b93da5c18ed649ebd4c2c79437db4611a6a1373e805a3cb001c64130b7",
            "expected-hostname": "github.com",
            "fb:app_id": "1401488693436528",
            "github-keyboard-shortcuts": "repository,copilot",
            "go-import": "github.com/zeus-fyi/zeus git https://github.com/zeus-fyi/zeus.git",
            "hostname": "github.com",
            "hovercard-subject-tag": "repository:564491293",
            "html-safe-nonce": "2913d63208e37c5021206497211c56243f4e75c9af58b3519502fea841c0383a",
            "octolytics-dimension-repository_id": "564491293",
            "octolytics-dimension-repository_is_fork": "false",
            "octolytics-dimension-repository_network_root_id": "564491293",
            "octolytics-dimension-repository_network_root_nwo": "zeus-fyi/zeus",
            "octolytics-dimension-repository_nwo": "zeus-fyi/zeus",
            "octolytics-dimension-repository_public": "true",
            "octolytics-dimension-user_id": "108896502",
            "octolytics-dimension-user_login": "zeus-fyi",
            "octolytics-url": "https://collector.github.com/github/collect",
            "og:description": "Zeus + SciFi = the power of the gods, meets the power of science fiction. Designing wisdom into intelligence, through intelligent design. - zeus-fyi/zeus",
            "og:image": "https://opengraph.githubassets.com/1ba1f8b2f54d32b988e84a9cc1423130b904f6b82bce864fdac78d44decd7971/zeus-fyi/zeus",
            "og:image:alt": "Zeus + SciFi = the power of the gods, meets the power of science fiction. Designing wisdom into intelligence, through intelligent design. - zeus-fyi/zeus",
            "og:image:height": "600",
            "og:image:width": "1200",
            "og:site_name": "GitHub",
            "og:title": "GitHub - zeus-fyi/zeus: Zeus + SciFi = the power of the gods, meets the power of science fiction. Designing wisdom into intelligence, through intelligent design.",
            "og:type": "object",
            "og:url": "https://github.com/zeus-fyi/zeus",
            "request-id": "F211:7746:462D29:60B5B7:65CC2EFA",
            "route-action": "disambiguate",
            "route-controller": "files",
            "route-pattern": "/:user_id/:repository",
            "theme-color": "#1e2327",
            "turbo-body-classes": "logged-out env-production page-responsive",
            "turbo-cache-control": "no-preview",
            "twitter:card": "summary_large_image",
            "twitter:description": "Zeus + SciFi = the power of the gods, meets the power of science fiction. Designing wisdom into intelligence, through intelligent design. - zeus-fyi/zeus",
            "twitter:image:src": "https://opengraph.githubassets.com/1ba1f8b2f54d32b988e84a9cc1423130b904f6b82bce864fdac78d44decd7971/zeus-fyi/zeus",
            "twitter:site": "@github",
            "twitter:title": "GitHub - zeus-fyi/zeus: Zeus + SciFi = the power of the gods, meets the power of science fiction. Designing wisdom into intelligence, through intelligent design.",
            "viewport": "width=device-width",
            "visitor-hmac": "8d923f54a3de7f969e192a292828ff40b015e3d5f3034648862063af7c579785",
            "visitor-payload": "eyJyZWZlcnJlciI6IiIsInJlcXVlc3RfaWQiOiJGMjExOjc3NDY6NDYyRDI5OjYwQjVCNzo2NUNDMkVGQSIsInZpc2l0b3JfaWQiOiI1NjIwOTk2Mzk4MzIzMDg1MDUwIiwicmVnaW9uX2VkZ2UiOiJpYWQiLCJyZWdpb25fcmVuZGVyIjoiaWFkIn0="
          }
        ],
        "softwaresourcecode": [
          {
            "author": "zeus-fyi",
            "name": "zeus",
            "text": "Documentation https://docs.zeus.fyi zK8s == Kubernetes + Zeus Here we overview the core concepts needed to understand how you can build, deploy, configure K8s apps using Zeus, with a full walkthrou..."
          }
        ]
      },
      "snippet": "Overview · Automates translation of kubernetes yaml configurations into representative SQL models · Users upload these infrastructure configurations via API ...",
      "title": "zeus-fyi/zeus: Zeus + SciFi = the power of the gods, meets ... - GitHub"
    },
    {
      "displayLink": "www.linkedin.com",
      "formattedUrl": "https://www.linkedin.com/company/zeusfyi",
      "htmlFormattedUrl": "https://www.linkedin.com/company/<b>zeusfyi</b>",
      "htmlSnippet": "<b>Zeusfyi</b> | 97 followers on LinkedIn. The Cloud Platform for Microservices &amp; SOA Development | Zeus + SciFi = the power of the gods, meets the power of&nbsp;...",
      "htmlTitle": "<b>Zeusfyi</b> | LinkedIn",
      "kind": "customsearch#result",
      "link": "https://www.linkedin.com/company/zeusfyi",
      "pagemap": {
        "cse_image": [
          {
            "src": "https://media.licdn.com/dms/image/D560BAQFOPTy4xndkdg/company-logo_200_200/0/1694541234347/zeus_fyi_logo?e=2147483647&v=beta&t=7DmMZgwnXXV2KEPMxgJBaBCLKbQBYdL45ucOmwMf_tc"
          }
        ],
        "cse_thumbnail": [
          {
            "height": "200",
            "src": "https://encrypted-tbn1.gstatic.com/images?q=tbn:ANd9GcQY74dFQu0EA2BA5oSr8ZQ0SbYF07F_FGMw_kXuJAK60utiOI75w2JvtQs",
            "width": "200"
          }
        ],
        "metatags": [
          {
            "al:android:app_name": "LinkedIn",
            "al:android:package": "com.linkedin.android",
            "al:android:url": "https://www.linkedin.com/company/zeusfyi",
            "al:ios:app_name": "LinkedIn",
            "al:ios:app_store_id": "288429040",
            "al:ios:url": "https://www.linkedin.com/company/zeusfyi",
            "bingbot": "nocache",
            "clientsideingraphs": "1",
            "linkedin:pagetag": "noncanonical_subdomain=control",
            "locale": "en_US",
            "og:description": "Zeusfyi | 97 followers on LinkedIn. The Cloud Platform for Microservices &amp; SOA Development | Zeus + SciFi &#61; the power of the gods, meets the power of science fiction. Designing wisdom into intelligence, through intelligent design.",
            "og:image": "https://media.licdn.com/dms/image/D560BAQFOPTy4xndkdg/company-logo_200_200/0/1694541234347/zeus_fyi_logo?e=2147483647&v=beta&t=7DmMZgwnXXV2KEPMxgJBaBCLKbQBYdL45ucOmwMf_tc",
            "og:title": "Zeusfyi | LinkedIn",
            "og:type": "article",
            "og:url": "https://www.linkedin.com/company/zeusfyi",
            "pagekey": "p_org_guest_company_overview",
            "twitter:card": "summary",
            "twitter:description": "Zeusfyi | 97 followers on LinkedIn. The Cloud Platform for Microservices &amp; SOA Development | Zeus + SciFi &#61; the power of the gods, meets the power of science fiction. Designing wisdom into intelligence, through intelligent design.",
            "twitter:image": "https://media.licdn.com/dms/image/D560BAQFOPTy4xndkdg/company-logo_200_200/0/1694541234347/zeus_fyi_logo?e=2147483647&v=beta&t=7DmMZgwnXXV2KEPMxgJBaBCLKbQBYdL45ucOmwMf_tc",
            "twitter:site": "@linkedin",
            "twitter:title": "Zeusfyi | LinkedIn",
            "viewport": "width=device-width, initial-scale=1.0"
          }
        ]
      },
      "snippet": "Zeusfyi | 97 followers on LinkedIn. The Cloud Platform for Microservices & SOA Development | Zeus + SciFi = the power of the gods, meets the power of ...",
      "title": "Zeusfyi | LinkedIn"
    },
    {
      "cacheId": "9Crv1Rr5-FsJ",
      "displayLink": "docs.zeus.fyi",
      "formattedUrl": "https://docs.zeus.fyi/docs/mockingbird/intro",
      "htmlFormattedUrl": "https://docs.zeus.fyi/docs/mockingbird/intro",
      "htmlSnippet": "Mockingbird is a time series controlled AI systems coordinator &amp; workflow executor system, data indexer and searcher, that builds control loops and hierarchical&nbsp;...",
      "htmlTitle": "Introduction | <b>Zeusfyi</b>",
      "kind": "customsearch#result",
      "link": "https://docs.zeus.fyi/docs/mockingbird/intro",
      "pagemap": {
        "cse_image": [
          {
            "src": "https://docs.zeus.fyi/img/icon.png"
          }
        ],
        "listitem": [
          {
            "name": "Introduction",
            "position": "1"
          }
        ],
        "metatags": [
          {
            "docsearch:docusaurus_tag": "docs-default-current",
            "docsearch:language": "en",
            "docsearch:version": "current",
            "docusaurus_locale": "en",
            "docusaurus_tag": "docs-default-current",
            "docusaurus_version": "current",
            "og:description": "Mockingbird is a time series controlled AI systems coordinator & workflow executor system, data indexer and searcher,",
            "og:image": "https://docs.zeus.fyi/img/icon.png",
            "og:title": "Introduction | Zeusfyi",
            "og:url": "https://docs.zeus.fyi/docs/mockingbird/intro",
            "twitter:card": "summary_large_image",
            "twitter:image": "https://docs.zeus.fyi/img/icon.png",
            "viewport": "width=device-width,initial-scale=1"
          }
        ]
      },
      "snippet": "Mockingbird is a time series controlled AI systems coordinator & workflow executor system, data indexer and searcher, that builds control loops and hierarchical ...",
      "title": "Introduction | Zeusfyi"
    },
    {
      "displayLink": "www.linkedin.com",
      "formattedUrl": "https://www.linkedin.com/in/alexandersgeorge",
      "htmlFormattedUrl": "https://www.linkedin.com/in/alexandersgeorge",
      "htmlSnippet": "Taking creativity to the art of cloud computing to lower the technical skill barrier to… | Learn more about Alex George&#39;s work experience, education,&nbsp;...",
      "htmlTitle": "Alex George - <b>Zeusfyi</b> | LinkedIn",
      "kind": "customsearch#result",
      "link": "https://www.linkedin.com/in/alexandersgeorge",
      "pagemap": {
        "Person": [
          {}
        ],
        "cse_image": [
          {
            "src": "https://media.licdn.com/dms/image/D5603AQGlm0O2wdq1Kw/profile-displayphoto-shrink_800_800/0/1672093557250?e=2147483647&v=beta&t=H33ZwMZVDP7RS2wo41-TyQtE6ghyh50kHFHxhCtfbKo"
          }
        ],
        "cse_thumbnail": [
          {
            "height": "225",
            "src": "https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcQtbcVDtQDoQ4J-9_iLKlRqt4X99-7FmXTJ9CCwiNaJHRCsVt2wi3u09T8",
            "width": "225"
          }
        ],
        "metatags": [
          {
            "al:android:app_name": "LinkedIn",
            "al:android:package": "com.linkedin.android",
            "al:android:url": "https://www.linkedin.com/in/alexandersgeorge",
            "al:ios:app_name": "LinkedIn",
            "al:ios:app_store_id": "288429040",
            "al:ios:url": "https://www.linkedin.com/in/alexandersgeorge",
            "litmsprofilename": "public-profile-frontend",
            "locale": "en_US",
            "og:description": "Taking creativity to the art of cloud computing to lower the technical skill barrier to… | Learn more about Alex George's work experience, education, connections & more by visiting their profile on LinkedIn",
            "og:image": "https://media.licdn.com/dms/image/D5603AQGlm0O2wdq1Kw/profile-displayphoto-shrink_800_800/0/1672093557250?e=2147483647&v=beta&t=H33ZwMZVDP7RS2wo41-TyQtE6ghyh50kHFHxhCtfbKo",
            "og:title": "Alex George - Zeusfyi | LinkedIn",
            "og:type": "profile",
            "og:url": "https://www.linkedin.com/in/alexandersgeorge",
            "pagekey": "public_profile_v3_mobile",
            "platform": "https://static.licdn.com/aero-v1/sc/h/cz36ewtx1ig2xy0w9zx0c4ww1",
            "platform-worker": "https://static.licdn.com/aero-v1/sc/h/7nirg34a8ey4y2l4rw7xgwxx4",
            "profile:first_name": "Alex",
            "profile:last_name": "George",
            "twitter:card": "summary",
            "twitter:description": "Taking creativity to the art of cloud computing to lower the technical skill barrier to… | Learn more about Alex George's work experience, education, connections & more by visiting their profile on LinkedIn",
            "twitter:image": "https://media.licdn.com/dms/image/D5603AQGlm0O2wdq1Kw/profile-displayphoto-shrink_800_800/0/1672093557250?e=2147483647&v=beta&t=H33ZwMZVDP7RS2wo41-TyQtE6ghyh50kHFHxhCtfbKo",
            "twitter:site": "@Linkedin",
            "twitter:title": "Alex George - Zeusfyi | LinkedIn",
            "ubba": "https://static.licdn.com/aero-v1/sc/h/5v7lpzb1fp10xxmwvevpgbj5h",
            "viewport": "width=device-width, initial-scale=1.0"
          }
        ]
      },
      "snippet": "Taking creativity to the art of cloud computing to lower the technical skill barrier to… | Learn more about Alex George's work experience, education, ...",
      "title": "Alex George - Zeusfyi | LinkedIn"
    },
    {
      "cacheId": "Jk8RsCil5csJ",
      "displayLink": "status.zeus.fyi",
      "formattedUrl": "https://status.zeus.fyi/",
      "htmlFormattedUrl": "https://status.zeus.fyi/",
      "htmlSnippet": "Welcome to <b>Zeusfyi&#39;s</b> home for real-time and historical data on system performance.",
      "htmlTitle": "<b>Zeusfyi</b> Status",
      "kind": "customsearch#result",
      "link": "https://status.zeus.fyi/",
      "pagemap": {
        "metatags": [
          {
            "_globalsign-domain-verification": "y_VzfckMy4iePo5oDJNivyYIjh8LffYa4jzUndm_bZ",
            "handheldfriendly": "True",
            "issued": "1708359505",
            "mobileoptimized": "320",
            "viewport": "width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0"
          }
        ]
      },
      "snippet": "Welcome to Zeusfyi's home for real-time and historical data on system performance.",
      "title": "Zeusfyi Status"
    },
    {
      "displayLink": "www.linkedin.com",
      "formattedUrl": "https://www.linkedin.com/.../alexandersgeorge_serverless-evm-zeusfyi-acti...",
      "htmlFormattedUrl": "https://www.linkedin.com/.../alexandersgeorge_serverless-evm-<b>zeusfyi</b>-acti...",
      "htmlSnippet": "Nov 1, 2023 <b>...</b> Introducing... Serverless EVM Simulation Environments Overview Serverless EVM execution environments that you can use QuickNode,&nbsp;...",
      "htmlTitle": "Alex George on LinkedIn: Serverless EVM | <b>Zeusfyi</b>",
      "kind": "customsearch#result",
      "link": "https://www.linkedin.com/posts/alexandersgeorge_serverless-evm-zeusfyi-activity-7125531815360020482-BPsb",
      "pagemap": {
        "cse_image": [
          {
            "src": "https://static.licdn.com/aero-v1/sc/h/c45fy346jw096z9pbphyyhdz7"
          }
        ],
        "cse_thumbnail": [
          {
            "height": "170",
            "src": "https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcReZZwU-Nj7Psst_HLlTt1mnnJPzL2RWdkIan0Ee9L-u1pereMyh2NU5lc",
            "width": "297"
          }
        ],
        "metatags": [
          {
            "al:android:app_name": "LinkedIn",
            "al:android:package": "com.linkedin.android",
            "al:android:url": "https://www.linkedin.com/posts/alexandersgeorge_serverless-evm-zeusfyi-activity-7125531815360020482-BPsb",
            "al:ios:app_name": "LinkedIn",
            "al:ios:app_store_id": "288429040",
            "al:ios:url": "https://www.linkedin.com/posts/alexandersgeorge_serverless-evm-zeusfyi-activity-7125531815360020482-BPsb",
            "lnkd:url": "https://www.linkedin.com/feed/update/urn:li:activity:7125531815360020482",
            "locale": "en_US",
            "og:description": "Introducing...\nServerless EVM Simulation Environments\n\nOverview\n\nServerless EVM execution environments that you can use QuickNode, self-managed, or other…",
            "og:image": "https://static.licdn.com/aero-v1/sc/h/c45fy346jw096z9pbphyyhdz7",
            "og:title": "Alex George on LinkedIn: Serverless EVM | Zeusfyi",
            "og:type": "article",
            "og:url": "https://www.linkedin.com/posts/alexandersgeorge_serverless-evm-zeusfyi-activity-7125531815360020482-BPsb",
            "pagekey": "p_public_post",
            "twitter:card": "summary_large_image",
            "twitter:description": "Introducing...\nServerless EVM Simulation Environments\n\nOverview\n\nServerless EVM execution environments that you can use QuickNode, self-managed, or other…",
            "twitter:image": "https://static.licdn.com/aero-v1/sc/h/c45fy346jw096z9pbphyyhdz7",
            "twitter:site": "@linkedin",
            "twitter:title": "Alex George on LinkedIn: Serverless EVM | Zeusfyi",
            "twitter:url": "https://www.linkedin.com/posts/alexandersgeorge_serverless-evm-zeusfyi-activity-7125531815360020482-BPsb",
            "viewport": "width=device-width, initial-scale=1.0"
          }
        ]
      },
      "snippet": "Nov 1, 2023 ... Introducing... Serverless EVM Simulation Environments Overview Serverless EVM execution environments that you can use QuickNode, ...",
      "title": "Alex George on LinkedIn: Serverless EVM | Zeusfyi"
    },
    {
      "cacheId": "o-pZ4Zyp7PwJ",
      "displayLink": "twitter.com",
      "formattedUrl": "https://twitter.com/ctrl_alt_lulz",
      "htmlFormattedUrl": "https://twitter.com/ctrl_alt_lulz",
      "htmlSnippet": "... any LLM application. medium.<b>zeus.fyi</b>. Unveiling the Next Generation of AI-Powered Workflow Automation. Mockingbird AI Systems Coordinator from <b>Zeusfyi</b>. 1.",
      "htmlTitle": "ctrl-alt-lulz (@ctrl_alt_lulz) / X",
      "kind": "customsearch#result",
      "link": "https://twitter.com/ctrl_alt_lulz",
      "pagemap": {
        "collection": [
          {
            "name": "Profile posts"
          }
        ],
        "creativework": [
          {
            "name": "Expanded Tweet URLs",
            "url": "https://t.co/zSIR4Zo4J3"
          },
          {
            "name": "Expanded Tweet URLs",
            "url": "https://t.co/DFMkbJcRQC"
          },
          {
            "name": "Expanded Tweet URLs",
            "url": "https://t.co/Zmm3iw8yhf"
          }
        ],
        "cse_image": [
          {
            "src": "https://pbs.twimg.com/profile_images/1607461193626832896/f7-Ccew8_200x200.jpg"
          }
        ],
        "cse_thumbnail": [
          {
            "height": "200",
            "src": "https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcS-9z5_3M34UpkTtdqRw4qWvY5lyJdkNVU4LmycvvYGyz43Ge-1Qbpj8dN-",
            "width": "200"
          }
        ],
        "imageobject": [
          {
            "caption": "Image",
            "contenturl": "https://pbs.twimg.com/media/GFcADXpaUAAo-4J.jpg",
            "thumbnailurl": "https://pbs.twimg.com/media/GFcADXpaUAAo-4J?format=jpg&name=thumb",
            "width": "1925"
          },
          {
            "caption": "Image",
            "contenturl": "https://pbs.twimg.com/media/GFm7EVRaQAAZexe.jpg",
            "thumbnailurl": "https://pbs.twimg.com/media/GFm7EVRaQAAZexe?format=jpg&name=thumb",
            "width": "1202"
          },
          {
            "caption": "Image",
            "contenturl": "https://pbs.twimg.com/media/GFmTMvzaYAAJWLM.png",
            "thumbnailurl": "https://pbs.twimg.com/media/GFmTMvzaYAAJWLM?format=jpg&name=thumb",
            "width": "494"
          }
        ],
        "interactioncounter": [
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1752482994361954462/likes",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1752482994361954462/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1752482994361954462/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1752482994361954462",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755453721998471520/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755453721998471520/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755453721998471520/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755453721998471520",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755148829098385550/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755148829098385550/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755148829098385550/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755148829098385550",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755119465862443216/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755119465862443216/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755119465862443216/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755119465862443216",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754713727088419120/likes",
            "userinteractioncount": "2"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754713727088419120/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754713727088419120/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754713727088419120",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754706273063543151/likes",
            "userinteractioncount": "2"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754706273063543151/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754706273063543151/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754706273063543151",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754642780121620809/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754642780121620809/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754642780121620809/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754642780121620809",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754639623048065441/likes",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754639623048065441/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754639623048065441/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754639623048065441",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754600862641795321/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754600862641795321/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754600862641795321/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754600862641795321",
            "userinteractioncount": "3"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754602781531689306/likes",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754602781531689306/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754602781531689306/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754602781531689306",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754602431932190852/likes",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754602431932190852/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754602431932190852/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754602431932190852",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754594675363852543/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754594675363852543/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754594675363852543/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754594675363852543",
            "userinteractioncount": "3"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754598413768053026/likes",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754598413768053026/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754598413768053026/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754598413768053026",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754596684309700731/likes",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754596684309700731/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754596684309700731/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754596684309700731",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754595449095524767/likes",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754595449095524767/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754595449095524767/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754595449095524767",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754594112479916413/likes",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754594112479916413/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754594112479916413/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754594112479916413",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754589242557440472/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754589242557440472/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754589242557440472/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754589242557440472",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754585874145456588/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754585874145456588/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754585874145456588/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754585874145456588",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754571613453135979/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754571613453135979/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754571613453135979/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754571613453135979",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754570012818649416/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754570012818649416/retweets",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754570012818649416/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754570012818649416",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/LikeAction",
            "name": "Likes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754565335829958741/likes",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Retweets",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754565335829958741/retweets",
            "userinteractioncount": "1"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Quotes",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754565335829958741/retweets/with_comments",
            "userinteractioncount": "0"
          },
          {
            "interactiontype": "https://schema.org/InteractAction",
            "name": "Replies",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754565335829958741",
            "userinteractioncount": "0"
          }
        ],
        "metatags": [
          {
            "al:android:app_name": "X",
            "al:android:package": "com.twitter.android",
            "al:android:url": "twitter://user?screen_name=ctrl_alt_lulz",
            "al:ios:app_name": "X",
            "al:ios:app_store_id": "333903271",
            "al:ios:url": "twitter://user?screen_name=ctrl_alt_lulz",
            "apple-mobile-web-app-status-bar-style": "white",
            "apple-mobile-web-app-title": "Twitter",
            "facebook-domain-verification": "x6sdcc8b5ju3bh8nbm59eswogvg6t1",
            "fb:app_id": "2231777543",
            "mobile-web-app-capable": "yes",
            "og:description": "We’re building, wisdom, then intelligence, into AI systems, using intelligent design. Our WAI system is called Mockingbird.",
            "og:image": "https://pbs.twimg.com/profile_images/1607461193626832896/f7-Ccew8_200x200.jpg",
            "og:site_name": "X (formerly Twitter)",
            "og:title": "ctrl-alt-lulz (@ctrl_alt_lulz) on X",
            "og:type": "profile",
            "og:url": "https://twitter.com/ctrl_alt_lulz",
            "theme-color": "#FFFFFF",
            "title": "ctrl-alt-lulz (@ctrl_alt_lulz) on X",
            "viewport": "width=device-width,initial-scale=1,maximum-scale=1,user-scalable=0,viewport-fit=cover"
          }
        ],
        "person": [
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          },
          {
            "additionalname": "ctrl_alt_lulz",
            "givenname": "ctrl-alt-lulz",
            "identifier": "1520852589067194373"
          }
        ],
        "socialmediaposting": [
          {
            "articlebody": "Mockingbird AI systems coordinator now in Beta. Auto-Evals Approval Triggers for API calls JSON schema dict Auto-Prompt Chunking https://medium.zeus.fyi/unveiling-the-next-generation-of-ai-powered-...",
            "commentcount": "1",
            "datecreated": "2024-01-31T00:04:40.000Z",
            "datepublished": "2024-01-31T00:04:40.000Z",
            "identifier": "1752482994361954462",
            "position": "1",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1752482994361954462"
          },
          {
            "articlebody": "Interesting idea for an AI experiment you can run on Mockingbird. Ask it to come up with a rubric for anything, then have it score something by providing just the 1-5 score, an the other have...",
            "commentcount": "0",
            "datecreated": "2024-02-08T04:49:16.000Z",
            "datepublished": "2024-02-08T04:49:16.000Z",
            "identifier": "1755453721998471520",
            "position": "2",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755453721998471520"
          },
          {
            "articlebody": "It’s not whey bro, it’s wai. 📡🐦‍⬛",
            "commentcount": "0",
            "datecreated": "2024-02-07T08:37:44.000Z",
            "datepublished": "2024-02-07T08:37:44.000Z",
            "identifier": "1755148829098385550",
            "position": "3",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755148829098385550"
          },
          {
            "articlebody": "I was thinking, as I was adding in some Helm chart API coding, that I could auto-eval helm chart outputs to auto-k8s, and basically build anything via iterative workflow auto-evals and feedback lol.",
            "commentcount": "0",
            "datecreated": "2024-02-07T06:41:03.000Z",
            "datepublished": "2024-02-07T06:41:03.000Z",
            "identifier": "1755119465862443216",
            "position": "4",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1755119465862443216"
          },
          {
            "articlebody": "Just thought up a quick hack for a sick AI powered reliable UI navigator. Iterative cycled workflows that auto eval themselves on ui scripts like applescripts to drive a UI. Could probably...",
            "commentcount": "0",
            "datecreated": "2024-02-06T03:48:48.000Z",
            "datepublished": "2024-02-06T03:48:48.000Z",
            "identifier": "1754713727088419120",
            "position": "5",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754713727088419120"
          },
          {
            "articlebody": "The takeaway from our Strong Supervisor research excerpt preview is that these generated knowledge trees can be saved, copied, and built, and extended over time, making other subsequent “hard”...",
            "commentcount": "0",
            "datecreated": "2024-02-06T03:19:11.000Z",
            "datepublished": "2024-02-06T03:19:11.000Z",
            "identifier": "1754706273063543151",
            "isbasedon": "https://twitter.com/ctrl_alt_lulz/status/1753870642225770682",
            "position": "6",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754706273063543151"
          },
          {
            "articlebody": "Wisdom of the crowd, also applies to models.",
            "commentcount": "0",
            "datecreated": "2024-02-05T23:06:53.000Z",
            "datepublished": "2024-02-05T23:06:53.000Z",
            "identifier": "1754642780121620809",
            "position": "7",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754642780121620809"
          },
          {
            "articlebody": "What we're close to wrapping up. Fluent external Kubernetes clusters integration in like 30s. Upload kubeconfigs, set service account values in the platform secrets. Let Zeus handle everything...",
            "commentcount": "0",
            "datecreated": "2024-02-05T22:54:20.000Z",
            "datepublished": "2024-02-05T22:54:20.000Z",
            "identifier": "1754639623048065441",
            "position": "8",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754639623048065441"
          },
          {
            "articlebody": "Anyone else out there building an AI systems controller that's capable of controlling, building, maintaining, and planning cloud infrastructure? Or is it just me who's crazy enough to do that....",
            "commentcount": "3",
            "datecreated": "2024-02-05T20:20:19.000Z",
            "datepublished": "2024-02-05T20:20:19.000Z",
            "identifier": "1754600862641795321",
            "position": "9",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754600862641795321"
          },
          {
            "articlebody": "Also, you can 1-click/ one API call rollout changes to fleets across multi-cloud, multi-region right now on our platform. Simple enough for people + AI.",
            "commentcount": "0",
            "datecreated": "2024-02-05T20:27:56.000Z",
            "datepublished": "2024-02-05T20:27:56.000Z",
            "identifier": "1754602781531689306",
            "ispartof": "https://twitter.com/ctrl_alt_lulz/status/1754600862641795321",
            "position": "10",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754602781531689306"
          },
          {
            "articlebody": "If you work in devops, or related fields in Kubernetes and basic cloud app infra. This system will deprecate that in the not so distant future. Guess who it's been learning from?",
            "commentcount": "0",
            "datecreated": "2024-02-05T20:26:33.000Z",
            "datepublished": "2024-02-05T20:26:33.000Z",
            "identifier": "1754602431932190852",
            "ispartof": "https://twitter.com/ctrl_alt_lulz/status/1754600862641795321",
            "position": "11",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754602431932190852"
          },
          {
            "articlebody": "Mockingbird today is already a semi-automatic rifle, and I haven't even connected it to all the platforms we index. In a month it goes full machine gun kelly.",
            "commentcount": "3",
            "datecreated": "2024-02-05T19:55:44.000Z",
            "datepublished": "2024-02-05T19:55:44.000Z",
            "identifier": "1754594675363852543",
            "position": "12",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754594675363852543"
          },
          {
            "articlebody": "Rich user engagement stats check, self reliant AI economic models via Ad revenue, check, adaptive metrics system an control feedbacks, check check and check it out.",
            "commentcount": "0",
            "datecreated": "2024-02-05T20:10:35.000Z",
            "datepublished": "2024-02-05T20:10:35.000Z",
            "identifier": "1754598413768053026",
            "ispartof": "https://twitter.com/ctrl_alt_lulz/status/1754594675363852543",
            "position": "13",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754598413768053026"
          },
          {
            "articlebody": "Platforms we've been indexing for months: Reddit, Twitter, Discord, Telegram",
            "commentcount": "0",
            "datecreated": "2024-02-05T20:03:43.000Z",
            "datepublished": "2024-02-05T20:03:43.000Z",
            "identifier": "1754596684309700731",
            "ispartof": "https://twitter.com/ctrl_alt_lulz/status/1754594675363852543",
            "position": "14",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754596684309700731"
          },
          {
            "articlebody": "423k orchestrations served and counting",
            "commentcount": "0",
            "datecreated": "2024-02-05T19:58:48.000Z",
            "datepublished": "2024-02-05T19:58:48.000Z",
            "identifier": "1754595449095524767",
            "ispartof": "https://twitter.com/ctrl_alt_lulz/status/1754594675363852543",
            "position": "15",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754595449095524767"
          },
          {
            "articlebody": "Imagine if you could instantly think of 1M ways to approach a problem, an then aggregate and bubble the top 3, and then continue this process until you solve any problem. You'd be smart as f huh.",
            "commentcount": "0",
            "datecreated": "2024-02-05T19:53:29.000Z",
            "datepublished": "2024-02-05T19:53:29.000Z",
            "identifier": "1754594112479916413",
            "position": "16",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754594112479916413"
          },
          {
            "articlebody": "Probably the biggest thing I’m working on is auto-generating model heuristics for generalized problem solving with adaptive metrics and task evolvers, executive workflow planners and multi-model...",
            "commentcount": "0",
            "datecreated": "2024-02-05T19:34:08.000Z",
            "datepublished": "2024-02-05T19:34:08.000Z",
            "identifier": "1754589242557440472",
            "position": "17",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754589242557440472"
          },
          {
            "articlebody": "Everyone should know that love grows, not falls. This man needs that lesson.",
            "commentcount": "0",
            "datecreated": "2024-02-05T19:20:45.000Z",
            "datepublished": "2024-02-05T19:20:45.000Z",
            "identifier": "1754585874145456588",
            "isbasedon": "https://twitter.com/ESYudkowsky/status/1754324150267855266",
            "position": "18",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754585874145456588"
          },
          {
            "articlebody": "If you’re a ML researcher looking for new topics, my own barometer is binary: Does this get me significantly closer to ASI or not? If not, you aren’t solving the right problems.",
            "commentcount": "0",
            "datecreated": "2024-02-05T18:24:05.000Z",
            "datepublished": "2024-02-05T18:24:05.000Z",
            "identifier": "1754571613453135979",
            "position": "19",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754571613453135979"
          },
          {
            "articlebody": "“All you need is attention” should have been called “All you need is feedback and control systems”",
            "commentcount": "0",
            "datecreated": "2024-02-05T18:17:44.000Z",
            "datepublished": "2024-02-05T18:17:44.000Z",
            "identifier": "1754570012818649416",
            "position": "20",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754570012818649416"
          },
          {
            "datecreated": "2024-02-05T18:00:40.000Z",
            "datepublished": "2024-02-05T18:00:40.000Z",
            "identifier": "1754565717901603208",
            "position": "21"
          },
          {
            "articlebody": "how is continual learning different from building and refining model generated score cards from env + adaptive stats? a more accurate rating rubric is no different than learned knowledge. try...",
            "commentcount": "0",
            "datecreated": "2024-02-05T17:59:08.000Z",
            "datepublished": "2024-02-05T17:59:08.000Z",
            "identifier": "1754565335829958741",
            "ispartof": "https://twitter.com/omarsar0/status/1754509803060171262",
            "url": "https://twitter.com/ctrl_alt_lulz/status/1754565335829958741"
          }
        ]
      },
      "snippet": "... any LLM application. medium.zeus.fyi. Unveiling the Next Generation of AI-Powered Workflow Automation. Mockingbird AI Systems Coordinator from Zeusfyi. 1.",
      "title": "ctrl-alt-lulz (@ctrl_alt_lulz) / X"
    },
    {
      "cacheId": "06q4DWHQEaMJ",
      "displayLink": "marketplace.quicknode.com",
      "formattedUrl": "https://marketplace.quicknode.com/add-on/zeusfyi-4",
      "htmlFormattedUrl": "https://marketplace.quicknode.com/add-on/<b>zeusfyi</b>-4",
      "htmlSnippet": "Sep 8, 2023 <b>...</b> Adaptive RPC Load Balancer. by <b>Zeusfyi</b> 4Supported chains 4Plans ... Curl Example: Procedures. &quot;procedure&quot;: &quot;eth_maxBlockAggReduce&quot; curl --location&nbsp;...",
      "htmlTitle": "Adaptive RPC Load Balancer - QuickNode Marketplace",
      "kind": "customsearch#result",
      "link": "https://marketplace.quicknode.com/add-on/zeusfyi-4",
      "pagemap": {
        "cse_image": [
          {
            "src": "https://marketplace.quicknode.com/rails/active_storage/representations/redirect/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaHBBb1liIiwiZXhwIjpudWxsLCJwdXIiOiJibG9iX2lkIn19--4141b9fc16d8742e1b125e77b6a78f843252eeac/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaDdCem9MWm05eWJXRjBTU0lJY0c1bkJqb0dSVlE2RTNKbGMybDZaVjkwYjE5bWFXeHNXd2RwQXJBRWFRSjJBZz09IiwiZXhwIjpudWxsLCJwdXIiOiJ2YXJpYXRpb24ifX0=--b90ee6b8f6859f22e40c1be1e1930bb8ed7123dc/Illustration(5).png"
          }
        ],
        "metatags": [
          {
            "csrf-param": "authenticity_token",
            "csrf-token": "LS6_c_7nDCVHX_FnNK6d9dAdMo9PEc5JPXD643ceol9kgPCmf_MsNqQsJSsmK7glC7y4zAEc14onoxbDybE4kA",
            "google": "notranslate",
            "og:description": "Load balance RPC volume at scale across QuickNode, other infrastructure RPC providers, and even your own endpoints if you run your own infrastructure.",
            "og:image": "https://marketplace.quicknode.com/rails/active_storage/representations/redirect/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaHBBb1liIiwiZXhwIjpudWxsLCJwdXIiOiJibG9iX2lkIn19--4141b9fc16d8742e1b125e77b6a78f843252eeac/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaDdCem9MWm05eWJXRjBTU0lJY0c1bkJqb0dSVlE2RTNKbGMybDZaVjkwYjE5bWFXeHNXd2RwQXJBRWFRSjJBZz09IiwiZXhwIjpudWxsLCJwdXIiOiJ2YXJpYXRpb24ifX0=--b90ee6b8f6859f22e40c1be1e1930bb8ed7123dc/Illustration(5).png",
            "og:site_name": "Adaptive RPC Load Balancer - QuickNode Marketplace",
            "og:title": "Adaptive RPC Load Balancer - QuickNode Marketplace",
            "og:type": "website",
            "og:url": "https://marketplace.quicknode.com/add-on/zeusfyi-4",
            "twitter:card": "summary_large_image",
            "twitter:creator": "@QuickNode",
            "twitter:description": "Load balance RPC volume at scale across QuickNode, other infrastructure RPC providers, and even your own endpoints if you run your own infrastructure.",
            "twitter:image:src": "https://marketplace.quicknode.com/rails/active_storage/representations/redirect/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaHBBb1liIiwiZXhwIjpudWxsLCJwdXIiOiJibG9iX2lkIn19--4141b9fc16d8742e1b125e77b6a78f843252eeac/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaDdCem9MWm05eWJXRjBTU0lJY0c1bkJqb0dSVlE2RTNKbGMybDZaVjkwYjE5bWFXeHNXd2RwQXJBRWFRSjJBZz09IiwiZXhwIjpudWxsLCJwdXIiOiJ2YXJpYXRpb24ifX0=--b90ee6b8f6859f22e40c1be1e1930bb8ed7123dc/Illustration(5).png",
            "twitter:site": "@QuickNode",
            "twitter:title": "Adaptive RPC Load Balancer - QuickNode Marketplace",
            "viewport": "width=device-width, initial-scale=1"
          }
        ]
      },
      "snippet": "Sep 8, 2023 ... Adaptive RPC Load Balancer. by Zeusfyi 4Supported chains 4Plans ... Curl Example: Procedures. \"procedure\": \"eth_maxBlockAggReduce\" curl --location ...",
      "title": "Adaptive RPC Load Balancer - QuickNode Marketplace"
    }
  ],
  "kind": "customsearch#search",
  "queries": {
    "nextPage": [
      {
        "count": 10,
        "cx": "71b870adea1674892",
        "inputEncoding": "utf8",
        "outputEncoding": "utf8",
        "safe": "off",
        "searchTerms": "zeusfyi",
        "startIndex": 11,
        "title": "Google Custom Search - zeusfyi",
        "totalResults": "611"
      }
    ],
    "request": [
      {
        "count": 10,
        "cx": "71b870adea1674892",
        "inputEncoding": "utf8",
        "outputEncoding": "utf8",
        "safe": "off",
        "searchTerms": "zeusfyi",
        "startIndex": 1,
        "title": "Google Custom Search - zeusfyi",
        "totalResults": "611"
      }
    ]
  },
  "searchInformation": {
    "formattedSearchTime": "0.22",
    "formattedTotalResults": "611",
    "searchTime": 0.215921,
    "totalResults": "611"
  },
  "spelling": {
    "correctedQuery": "zeus fyi",
    "htmlCorrectedQuery": "<b><i>zeus fyi</i></b>"
  },
  "url": {
    "template": "https://www.googleapis.com/customsearch/v1?q={searchTerms}&num={count?}&start={startIndex?}&lr={language?}&safe={safe?}&cx={cx?}&sort={sort?}&filter={filter?}&gl={gl?}&cr={cr?}&googlehost={googleHost?}&c2coff={disableCnTwTranslation?}&hq={hq?}&hl={hl?}&siteSearch={siteSearch?}&siteSearchFilter={siteSearchFilter?}&exactTerms={exactTerms?}&excludeTerms={excludeTerms?}&linkSite={linkSite?}&orTerms={orTerms?}&dateRestrict={dateRestrict?}&lowRange={lowRange?}&highRange={highRange?}&searchType={searchType}&fileType={fileType?}&rights={rights?}&imgSize={imgSize?}&imgType={imgType?}&imgColorType={imgColorType?}&imgDominantColor={imgDominantColor?}&alt=json",
    "type": "application/json"
  }
}
`
