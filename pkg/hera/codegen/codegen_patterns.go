package hera_v1_codegen

import (
	"context"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

/*
insertMsgCtx := workflow.WithActivityOptions(ctx, ao)
var sq *hera_search.TwitterSearchQuery
err = workflow.ExecuteActivity(insertMsgCtx, h.SelectTwitterSearchQuery, ou, groupName).Get(insertMsgCtx, &sq)

	if err != nil {
		logger.Error("failed to execute SelectTwitterSearchQuery", "Error", err)
		// You can decide if you want to return the error or continue monitoring.
		return err
	}
*/

func AddActivityToWorkflow(f filepaths.Path) string {
	actInst := `add an activity to the AiIngestReddit workflow`
	bins := &BuildAiInstructions{
		Path:               f,
		PromptInstructions: actInst,
		OrderedInstructions: []BuildAiFileInstruction{
			{
				DirIn:                DbSchemaDir,
				FileName:             "603_ai_twitter.sql",
				FileLevelInstruction: "Use the example SQL table definitions and indexes in this file to create the new SQL table definitions",
			},
			{
				DirIn:    HeraDbModelsDir + "/search",
				FileName: "twitter.go",
				FileLevelInstruction: `Create the lowercase functions using the exact styling shown and assign it to q.RawQuery
					INSERT INTO "public"."ai_incoming_tweets" ("search_id", "tweet_id", "message_text")
					VALUES ($1, $2, $3)
					ON CONFLICT ("tweet_id")
					DO UPDATE SET
						"message_text" = EXCLUDED."message_text"
					RETURNING "tweet_id"
`,
				OrderedFileFunctionInstructions: []FunctionInstruction{
					{
						FunctionInstruction: "Reference 1",
						FunctionInfo: FunctionInfo{
							Name: "insertTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 1",
						FunctionInfo: FunctionInfo{
							Name: "InsertTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "selectTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "SelectTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "insertIncomingTweets",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "InsertIncomingTweets",
						},
					},
				},
			},
		},
	}
	prompt := GenerateInstructions(context.Background(), bins)
	return prompt
}

func GenerateSqlTableFromExample(f filepaths.Path) string {
	actInst := `write: two new SQL table definitions called ai_reddit_search_query and ai_reddit_incoming_posts. we index this data and search the date range using the primary key id which
				is a unix timestamp, and filter on text from post body, subreddit, and title for the search query derive a solution based on how we used indexing for tweets
				the final solution should allow us to insert the below data into the new tables and then query the data using the below data as an example
		        group your generated answers using the filepath of the file that contains most likely the reference for the answer 
				`
	bins := &BuildAiInstructions{
		Path:               f,
		PromptInstructions: actInst,
		OrderedInstructions: []BuildAiFileInstruction{
			{
				DirIn:                DbSchemaDir,
				FileName:             "603_ai_twitter.sql",
				FileLevelInstruction: "Use the example SQL table definitions and indexes in this file to create the new SQL table definitions",
			},
			{
				DirIn:    HeraDbModelsDir + "/search",
				FileName: "twitter.go",
				FileLevelInstruction: `Create the lowercase functions using the exact styling shown and assign it to q.RawQuery
					INSERT INTO "public"."ai_incoming_tweets" ("search_id", "tweet_id", "message_text")
					VALUES ($1, $2, $3)
					ON CONFLICT ("tweet_id")
					DO UPDATE SET
						"message_text" = EXCLUDED."message_text"
					RETURNING "tweet_id"
`,
				OrderedFileFunctionInstructions: []FunctionInstruction{
					{
						FunctionInstruction: "Reference 1",
						FunctionInfo: FunctionInfo{
							Name: "insertTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 1",
						FunctionInfo: FunctionInfo{
							Name: "InsertTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "selectTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "SelectTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "insertIncomingTweets",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "InsertIncomingTweets",
						},
					},
				},
			},
		},
	}
	prompt := GenerateInstructions(context.Background(), bins)
	return prompt
}
