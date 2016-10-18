package deduper

import (
	"fmt"
	"testing"
	"time"

	lpkg "github.com/pestophagous/hackybeat/util/logger"
)

var counterFalseCheck uint = 0
var counterTrueCheck uint = 0

func ___wantApproval(t *testing.T, approval bool) {
	if approval == false {
		t.Errorf("Did not expect to be blocked/forbidden %d", counterFalseCheck)
	}
	counterFalseCheck++
}

func expectForbidden(t *testing.T, approval bool) {
	if approval == true {
		t.Errorf("Uh-oh. Where was the expected block/disapproval? %d", counterTrueCheck)
	}
	counterTrueCheck++
}

func verifyPurgeAtRestored(t *testing.T) {
	if purgeAtAge != defaultPurgeAt {
		t.Error("Somebody forgot to restore defaultPurgeAt.")
	}
}

func TestDeduperToolBasic(t *testing.T) {
	var errorMsgs []string
	logger := lpkg.NewNoopLogWithNilCheck()

	logger.L.Info = func(format string, v ...interface{}) {
		t.Logf(format, v...)
	}

	logger.L.Err = func(format string, v ...interface{}) {
		errorMsgs = append(errorMsgs, fmt.Sprintf(format, v...))
	}

	dedupeTool := NewDeduperTool("TestDeduperToolBasic", logger)

	dedupeAggressivePurgingScenario(t, dedupeTool)
	verifyPurgeAtRestored(t)

	dedupeModeratePurgingScenario(t, dedupeTool)
	verifyPurgeAtRestored(t)

	if len(errorMsgs) > 0 {
		t.Errorf("Errors while working with deduper.Tool: %v", errorMsgs)
	}
}

func dedupeAggressivePurgingScenario(t *testing.T, dedupeTool *Tool) {
	p := purgeAtAge
	defer func() { purgeAtAge = p }()

	time1 := time.Now()
	time0 := time1.Add(-1 * time.Minute)

	// first do a quick batch where things are NOT purged:
	___wantApproval(t, dedupeTool.IsGrantingApproval(time1, "dedupeAggressivePurgingScenario", 6353))
	___wantApproval(t, dedupeTool.IsGrantingApproval(time0, "dedupeAggressivePurgingScenario", 6353))
	expectForbidden(t, dedupeTool.IsGrantingApproval(time0, "dedupeAggressivePurgingScenario", 6353))

	// set purge (package var) so low that everything is always purged and we never find priors:
	purgeAtAge = 0

	// purge setting is so strict that everything is "stale" and too-old, and therefore presumed duplicate, and forbidden
	expectForbidden(t, dedupeTool.IsGrantingApproval(time1, "dedupeAggressivePurgingScenario", 6353))
	expectForbidden(t, dedupeTool.IsGrantingApproval(time0, "dedupeAggressivePurgingScenario", 6353))
	expectForbidden(t, dedupeTool.IsGrantingApproval(time0, "dedupeAggressivePurgingScenario", 6353))
}

func dedupeModeratePurgingScenario(t *testing.T, dedupeTool *Tool) {
	p := purgeAtAge
	defer func() { purgeAtAge = p }()

	halfDayOld := time.Now().Add(-12 * time.Hour)
	oneHourOld := time.Now().Add(-1 * time.Hour)

	// put the items in, and show that they stayed in:
	___wantApproval(t, dedupeTool.IsGrantingApproval(halfDayOld, "dedupeModeratePurgingScenario", 6353))
	___wantApproval(t, dedupeTool.IsGrantingApproval(oneHourOld, "dedupeModeratePurgingScenario", 6353))
	expectForbidden(t, dedupeTool.IsGrantingApproval(halfDayOld, "dedupeModeratePurgingScenario", 6353))
	expectForbidden(t, dedupeTool.IsGrantingApproval(oneHourOld, "dedupeModeratePurgingScenario", 6353))

	// set purge (package var) so that some are purged:
	purgeAtAge = -3 * time.Hour

	// this one is presumed duplicate due to stale age:
	expectForbidden(t, dedupeTool.IsGrantingApproval(halfDayOld, "dedupeModeratePurgingScenario", 6353))
	// this one should still be in the DB, and will be detected as a dupe because of that reason:
	expectForbidden(t, dedupeTool.IsGrantingApproval(oneHourOld, "dedupeModeratePurgingScenario", 6353))
}
