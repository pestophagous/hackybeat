package gitlog

import (
	"fmt"
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

func (this *gitLog) InstanceIdForLogging() string {
	return fmt.Sprintf("gitLog %v", gitRepoPath)
}

func (this *gitLog) DoPoll() time.Duration {
	this.counter = 0
	var err error
	var repo *git.Repository
	var walk *git.RevWalk

	repo, err = git.OpenRepository(gitRepoPath)
	if err != nil {
		this.logger.LogFailureOf("OpenRepository", this, err)
		return longDuration
	}
	defer repo.Free()

	walk, err = repo.Walk()
	if err != nil {
		this.logger.LogFailureOf("Repository.Walk", this, err)
		return longDuration
	}
	defer walk.Free()

	err = walk.PushRef("HEAD")
	this.logger.LogPossibleFailureOf("RevWalk.PushRef(\"HEAD\")", this, err)

	// the fun happens here. this is how we visit commits:
	err2 := walk.Iterate(this.onVisitCommit)
	this.logger.LogPossibleFailureOf("RevWalk.Iterate", this, err2)

	if err != nil || err2 != nil {
		return longDuration
	}

	// specify when we wish to be polled again
	return shortDuration
}
