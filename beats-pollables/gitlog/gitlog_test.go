package gitlog

import (
	"fmt"
	"testing"

	"github.com/libgit2/git2go"

	lpkg "github.com/pestophagous/hackybeat/util/logger"
)

func onItem(item *git.Commit, t *testing.T, counter *uint) {
	*counter++
	//t.Logf("sha %v\n", item.TreeId().String())
}

func TestGitLogBasic(t *testing.T) {
	var errorMsgs []string
	logger := lpkg.NewNoopLogWithNilCheck()

	logger.L.Info = func(format string, v ...interface{}) {
		t.Logf(format, v...)
	}

	logger.L.Err = func(format string, v ...interface{}) {
		errorMsgs = append(errorMsgs, fmt.Sprintf(format, v...))
	}

	// counter is intentionally captured by visitorFunc
	var counter uint = 0
	visitorFunc := func(item *git.Commit) { onItem(item, t, &counter) }

	// The 'logger' will manipulate our local var 'errorMsgs'.
	// The 'visitorFunc' will manipulate out local var 'counter'.
	// Therefore, it is via logger and visitorFunc that we can sense whether 'poller' behaves properly.
	poller := newPolledGitLog(logger, visitorFunc)

	// no matter how many times we poll, counter should always stop at maxCommitsPerPoll
	for i := 0; i < 3; i++ {
		counter = 0
		poller.DoPoll()
		if counter != maxCommitsPerPoll {
			t.Errorf("After gitLog.DoPoll, expected counter of %d, but got %d. (Maybe %v is too small a repo?)",
				maxCommitsPerPoll, counter, gitRepoPath)
		}
		t.Logf("counter %v\n", counter)
	}

	if len(errorMsgs) > 0 {
		t.Errorf("Errors while invoking polledFeed.DoPoll: %v", errorMsgs)
	}
}
