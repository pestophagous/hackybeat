package gitlog

import (
	"time"

	"github.com/libgit2/git2go"

	lpkg "github.com/pestophagous/hackybeat/util/logger"
)

const gitRepoPath = "/Users/someone/temp/testrepo"
const longDuration = time.Minute * 10
const shortDuration = time.Second * 60
const maxCommitsPerPoll = 100

type gitLog struct {
	logger   *lpkg.LogWithNilCheck
	callback func(item *git.Commit)
	counter  uint
}

func newPolledGitLog(log *lpkg.LogWithNilCheck, conv func(item *git.Commit)) *gitLog {
	p := &gitLog{
		logger:   log,
		callback: conv,
		counter:  0,
	}
	return p
}

// type RevWalkIterator func(commit *Commit) bool <-- declared in git2go
func (this *gitLog) onVisitCommit(commit *git.Commit) bool {
	this.counter++
	if this.counter > maxCommitsPerPoll {
		return false
	}

	this.callback(commit)
	return true
}

func (this *gitLog) OnShutdown() {
	this.logger.ReleaseLog()
}

// convenience function if you're already inside a block with a proven non-nil error:
func (this *gitLog) logFailureOf(what string, e error) {
	this.logger.Err("%s failed on %v. %v", what, gitRepoPath, e)
}

// convenience function when an error may or may not be nil, but you only want to log when it's non-nil:
func (this *gitLog) logPossibleFailureOf(what string, e error) {
	if e != nil {
		this.logFailureOf(what, e)
	}
}

func (this *gitLog) DoPoll() time.Duration {
	this.counter = 0
	var err error
	var repo *git.Repository
	var walk *git.RevWalk

	repo, err = git.OpenRepository(gitRepoPath)
	if err != nil {
		this.logFailureOf("OpenRepository", err)
		return longDuration
	}
	defer repo.Free()

	walk, err = repo.Walk()
	if err != nil {
		this.logFailureOf("Repository.Walk", err)
		return longDuration
	}
	defer walk.Free()

	err = walk.PushRef("HEAD")
	this.logPossibleFailureOf("RevWalk.PushRef(\"HEAD\")", err)

	// the fun happens here. this is how we visit commits:
	err2 := walk.Iterate(this.onVisitCommit)
	this.logPossibleFailureOf("RevWalk.Iterate", err2)

	if err != nil || err2 != nil {
		return longDuration
	}

	// specify when we wish to be polled again
	return shortDuration
}
