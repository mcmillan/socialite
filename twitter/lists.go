package twitter

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

const (
	maxDepth = 15
)

type ListStatusesResponse []Status

func GetListStatuses(sinceID string) (ListStatusesResponse, error) {
	var (
		combinedResponses ListStatusesResponse
		maxID             string
	)

	for currentDepth := 0; currentDepth < maxDepth; currentDepth++ {
		logFields := log.Fields{
			"package":      "twitter",
			"maxID":        maxID,
			"sinceID":      sinceID,
			"currentDepth": currentDepth,
			"maxDepth":     maxDepth,
		}

		log.WithFields(logFields).Info("Loading page")

		var listStatusesResponse ListStatusesResponse

		err := request("lists/statuses", "GET", map[string]string{
			"slug":              os.Getenv("TWITTER_LIST_SLUG"),
			"owner_screen_name": os.Getenv("TWITTER_LIST_OWNER_NAME"),
			"since_id":          sinceID,
			"max_id":            maxID,
			"include_rts":       "true",
		}, &listStatusesResponse)

		if err != nil {
			return nil, err
		}

		if len(listStatusesResponse) == 0 {
			log.WithFields(logFields).Info("No more tweets to load")
			break
		}

		maxID, err = listStatusesResponse[len(listStatusesResponse)-1].PrevID()

		if err != nil {
			log.WithFields(logFields).Error("Unable to ascertain new maxID, cautiously breaking")
			break
		}

		combinedResponses = append(combinedResponses, listStatusesResponse...)
	}

	return combinedResponses, nil
}
