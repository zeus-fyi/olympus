package ai_platform_service_orchestrations

const socialMediaJsonFormat = `[
  {
    "messageID": "1704249180",
    "messageBody": "message content text"
  }
]`

/*
	if taskInst.AnalysisResponseFormat == socialMediaEngagementResponseFormat {
		switch taskInst.RetrievalPlatform {
		case twitterPlatform:
			content = hera_search.FormatSearchResultsV3(sr)
			prompt = "\n" + "The messages you will be reading and writing must be in the following format" + socialMediaJsonFormat + "\nIf you are replying to any messages," +
				" you must return a JSON in the same format with messageID being the messageID you are responding to and messageBody being the text of your response\n" + prompt
		}
	}
*/

// TODO: add json extraction
