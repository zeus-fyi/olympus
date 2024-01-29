package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

func (t *ZeusWorkerTestSuite) TestJsonModelOutputActivity() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	act := NewZeusAiPlatformActivities()

	td, err := act.SelectTaskDefinition(ctx, ou, 1701313525731432000)
	t.Require().Nil(err)
	t.Require().NotEmpty(td)

	var schemas []*artemis_orchestrations.JsonSchemaDefinition
	for _, task := range td {
		t.Require().NotEmpty(task.Schemas)
		schemas = append(schemas, task.Schemas...)
	}

	fd := artemis_orchestrations.ConvertToFuncDef(schemas)
	t.Require().NotNil(fd)
	t.Require().NotNil(fd.Name)
	t.Require().NotNil(fd.Parameters)
	resp, err := act.CreateJsonOutputModelResponse(ctx, ou, hera_openai.OpenAIParams{
		Prompt:             JsonBooksInputExample,
		Model:              Gpt3JsonModel,
		FunctionDefinition: fd,
	})
	t.Require().Nil(err)
	//var m any
	//if len(resp.Response.Choices) > 0 && len(resp.Response.Choices[0].Message.ToolCalls) > 0 {
	//	m, err = UnmarshallOpenAiJsonInterfaceSlice(fd.Name, resp)
	//	t.Require().Nil(err)
	//
	//} else {
	//	m, err = UnmarshallOpenAiJsonInterface(fd.Name, resp)
	//	t.Require().Nil(err)
	//}
	//jsd := artemis_orchestrations.ConvertToJsonSchema(fd)
	//resp.JsonResponseResults = artemis_orchestrations.AssignMapValuesMultipleJsonSchemasSlice(jsd, m)
	t.Require().NotNil(resp)
	t.Require().NotNil(resp.JsonResponseResults)
	//
	//for _, res := range resp.JsonResponseResults {
	//
	//	for _, v := range res {
	//		t.Require().NotNil(v)
	//
	//		for _, f := range v.Fields {
	//			switch f.FieldName {
	//			case "title":
	//				t.Require().NotNil(f.StringValue)
	//				fmt.Println("title", *f.StringValue)
	//			case "score":
	//				t.Require().NotNil(f.NumberValue)
	//				fmt.Println("score", *f.NumberValue)
	//			}
	//		}
	//	}
	//	t.Require().NotNil(res)
	//}
}

const JsonBooksInputExample = `{
  "books": [
    {
      "title": "Stars Beyond Reach",
      "summary": "In a distant future, humanity discovers an uninhabited planet teeming with life. The novel explores the challenges of colonizing this new world."
    },
    {
      "title": "Galactic Echoes",
      "summary": "This novel follows a group of explorers who encounter a mysterious alien artifact that challenges their understanding of the universe."
    },
    {
      "title": "The Last Starship",
      "summary": "In a war-torn galaxy, the crew of the last surviving starship fights to protect the remnants of human civilization."
    },
    {
      "title": "Quantum Paradox",
      "summary": "A scientist uncovers a time travel conspiracy that could unravel the fabric of reality."
    },
    {
      "title": "Neptune's Secret",
      "summary": "Set on a floating city in the clouds of Neptune, this story revolves around a detective solving a crime that could spark a planetary revolution."
    },
    {
      "title": "Whispers of the Past",
      "summary": "A historical fiction novel set in the Victorian era, unraveling a tale of forbidden love and societal pressures."
    },
    {
      "title": "Journey Through The Sahara",
      "summary": "An adventurous travelogue detailing the experiences of a group traversing the Sahara desert."
    },
    {
      "title": "Murder at the Manor",
      "summary": "A classic whodunit mystery set in a sprawling English manor during the 1920s."
    },
    {
      "title": "The Baker Street Heist",
      "summary": "A fast-paced crime thriller about a daring robbery in the heart of London."
    },
    {
      "title": "Serenity Falls",
      "summary": "A heartwarming romance novel set in a small, picturesque mountain town."
    }
  ]
}
`
