package common

import (
	"testing"
	"time"

	beatcommon "github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/publisher"
)

func mockPublisher(event beatcommon.MapStr, opts ...publisher.ClientOption) bool {
	return true
}

func makeEvent(timepoint time.Time, text string, num int, snippets []string) beatcommon.MapStr {
	event := beatcommon.MapStr{
		"@timestamp": beatcommon.Time(timepoint),
		"type":       "hackybeat-common_test-testing",
		"somea":      text,
		"someb":      num,
		"somec":      snippets,
	}
	return event
}

func TestPublishWithDeduper(t *testing.T) {
	InstallPublisherFunc(mockPublisher)
	defer StopAllPollers()

	beatsPublish(makeEvent(time.Now(), "a@b.com", 8, []string{"aa", "B!", ":c:"}))
}
