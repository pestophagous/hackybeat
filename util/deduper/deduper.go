package deduper

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	lpkg "github.com/pestophagous/hackybeat/util/logger"
)

const tableName = "pestophagous_util_deduper_tool_go"

const defaultPurgeAt = -(7 * 24) * time.Hour

var purgeAtAge = defaultPurgeAt

type Tool struct {
	instanceName string
	logger       *lpkg.LogWithNilCheck
	db           *sql.DB // we'll need to call Close  // TODO
}

func NewDeduperTool(uniqueNamespace string, log *lpkg.LogWithNilCheck) *Tool {
	t := &Tool{
		instanceName: uniqueNamespace,
		logger:       log,
		db:           nil,
	}

	db, err := sql.Open("sqlite3", t.instanceName+".db")

	if err != nil {
		t.logFailureOf("sqlite3.Open", err)
		// We can imagine several possible fallback approaches, but have yet to implement any.
		// One approach would be to use the sqlite3 ':memory:' all-in-memory db, so we would have per-session deduping
		// at least.
		// Another approach would be to create a sort of 'null-object' deduper that always says each event is unique,
		// just to be able to run at all.
		panic("No fallback approach has yet been implemented for when sqlite3 db is not available.")
	} else {
		t.db = db
		t.locateOrCreateTable()
	}

	return t
}

func (this *Tool) IsGrantingApproval(eventTime time.Time, eventType string, object interface{}) bool {

	this.purgeOldValues()

	if eventTime.Before(this.cutoffTimeForPurge()) {
		// treat too-old, too-stale events as if we have seen them already.
		// otherwise, the deduper.Tool does nothing to help a simplistic (quick-and-dirty)
		// poller that reads an ENTIRE log start-to-finish each time.
		return false // <------- BAILING OUT
	}

	var blob bytes.Buffer
	enc := gob.NewEncoder(&blob)
	err := enc.Encode(object)
	this.logPossibleFailureOf("enc.Encode blob", err)

	// TODO: use something like an md5 hash/checksum (hence col name 'digest') instead of the blob/gob
	sel := fmt.Sprintf("SELECT EXISTS (SELECT 1 from %s where eventtype = ? and eventtime = ? and digest = ? LIMIT 1) AS existence", tableName)
	if this.existenceTest(sel, eventType, eventTime.UnixNano(), blob.Bytes()) {
		// found a prior.
		// caller should not proceed to consume the event.
		return false // <------- BAILING OUT
	}

	// TODO: use something like an md5 hash/checksum (hence col name 'digest') instead of the blob/gob
	q := fmt.Sprintf("INSERT INTO %s (eventtime, eventtype, digest) values (?, ?, ?)", tableName)

	stmt, err := this.db.Prepare(q) // TODO. call Prepare earlier? so we might fail at launch?
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(eventTime.UnixNano(), eventType, blob.Bytes())
	if err2 != nil {
		panic(err2)
	}

	return true // <-- caller SHOULD treat their event-object as fresh and never-before-seen
}

func (this *Tool) existenceTest(q string, args ...interface{}) bool {
	var rows *sql.Rows
	var err error
	result := false

	rows, err = this.db.Query(q, args...)
	if err != nil {
		this.logFailureOf("SELECT EXISTS", err)
		return result
	}
	defer rows.Close()

	rows.Next()

	var existence bool = false
	err = rows.Scan(&existence)
	this.logPossibleFailureOf("sql.Rows.Scan", err)
	if err == nil {
		result = existence
	}

	return result
}

func (this *Tool) tableExists() bool {
	q := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = '%s' LIMIT 1) AS existence", tableName)
	return this.existenceTest(q)
}

func (this *Tool) locateOrCreateTable() {

	// if it doesn't yet exist, then create it:
	if this.tableExists() == false {
		this.createTable()
	}

	// if it STILL doesn't exist...
	if this.tableExists() == false {
		// The fallbacks described elsewhere in this file would also be possible here.
		panic("Could neither locate nor create table. No fallback is implemented yet.")
	}

	this.sanityCheckTable()
}

func (this *Tool) sanityCheckTable() {
	// TODO. make sure (especially if this was a preexisting table) that it has the columns we expect

	this.purgeOldValues()
}

func (this *Tool) createTable() {
	// TODO: use something like an md5 hash/checksum (hence col name 'digest') instead of the blob/gob
	q := fmt.Sprintf("CREATE TABLE %s (eventtime integer, eventtype text, digest blob)", tableName)
	_, err := this.db.Exec(q)
	this.logPossibleFailureOf("sqlite3.Exec create table", err)
}

func (this *Tool) cutoffTimeForPurge() time.Time {
	return time.Now().Add(purgeAtAge)
}

func (this *Tool) purgeOldValues() {
	q := fmt.Sprintf("delete from %s where eventtime < ?", tableName)

	cutoff := this.cutoffTimeForPurge()

	_, err := this.db.Exec(q, cutoff.UnixNano())
	this.logPossibleFailureOf("sqlite3.Exec delete", err)
}

// convenience function if you're already inside a block with a proven non-nil error:
func (this *Tool) logFailureOf(what string, e error) {
	this.logger.Err("%s failed on %v. %v", what, this.instanceName, e)
}

// convenience function when an error may or may not be nil, but you only want to log when it's non-nil:
func (this *Tool) logPossibleFailureOf(what string, e error) {
	if e != nil {
		this.logFailureOf(what, e)
	}
}