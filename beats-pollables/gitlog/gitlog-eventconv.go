package gitlog

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/libgit2/git2go"

	pollcommon "github.com/pestophagous/hackybeat/beats-pollables/common"
)

func init() {
	pollcommon.RegisterPollable(newPolledGitLog(pollcommon.Logger, receiveItem))
}

// type polledFeed struct calls here when an item is ready. this method converts and forwards the item to libbeat publisher
func receiveItem(item *git.Commit) {

	// At a minimum, the event object must contain a @timestamp field and a type field. Beyond that, events can contain
	// any additional fields, and they can be created as often as necessary.
	event := common.MapStr{
		"@timestamp":   common.Time(item.Committer().When),
		"type":         "hackybeat-gitlog-testing",
		"gitauthor":    item.Author().Email,
		"gitcommitter": item.Committer().Email,
		"gitsummary":   item.Summary(),
		"gitoid":       item.TreeId().String(),
	}

	pollcommon.BeatsPublish(event)
}
