package apollo_status

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func StatusPageUpdate(apiKey, pageId, metricId string) {
	url := fmt.Sprintf("%s/pages/%s/metrics/%s/data.json", StatusApi, pageId, metricId)

	// need at least 1 data point for every 5 minutes
	// submit random data for the whole day
	totalPoints := 60 / 5 * 24
	epochInSeconds := int32(time.Now().Unix())

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < totalPoints; i++ {
		jsonStr := fmt.Sprintf(`{"data":{"timestamp":"%d","value":"%d"}}`,
			epochInSeconds-int32(i)*5*60,
			rand.Intn(1000))

		jsonBytes := []byte(jsonStr)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", apiKey))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		err = resp.Body.Close()
		if err != nil {
			panic(err)
		}

		if resp.StatusCode == 401 {
			genericError := "Error encountered. Please ensure that your page code and authorization key are correct."
			panic(genericError)
		}

		fmt.Printf("Submitted point %d of %d\n", i+1, totalPoints)
		time.Sleep(time.Second)
	}
}
