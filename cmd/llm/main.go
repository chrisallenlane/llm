// Package main encapsulates the executable
package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/chrisallenlane/llm/internal/config"
	"github.com/chrisallenlane/llm/internal/message"
	"github.com/chrisallenlane/llm/internal/session"
	"github.com/docopt/docopt-go"
)

const version = "0.0.0"

// docopt options string
const optstring = `
Usage:
	llm [options] db (path|destroy)
	llm [options] (q|question) (add|ask) [<msg>]
	llm [options] (q|question) (edit|rm) <id>
	llm [options] (q|question) send
	llm [options] (s|session) (new|edit|view|rm|use) <name>
	llm [options] (s|session) cp <orig> <name> [--all]
	llm [options] (s|session) log [<num>]
	llm [options] (s|session) ls
	llm [options] (s|session) rewind <id>
	llm [options] (s|session) search <string>
	llm [options] (s|session) truncate

Options:
	-f --forget          Do not read/write to/from message log
	-d --db=<path>       DB path`

func main() {
	// initialize options
	opts, err := docopt.ParseArgs(optstring, nil, version)
	if err != nil {
		// panic here, because this should never happen
		panic(fmt.Errorf("docopt failed to parse: %v", err))
	}

	// establish a SQLite3 connection
	dbPath, err := config.DBPath(opts, runtime.GOOS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to construct database path: %s: %v\n", dbPath, err)
		os.Exit(1)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to database %s: %v\n", dbPath, err)
		os.Exit(1)
	}

	// VACUUM the database to keep it tidy
	db.Exec("VACUUM;")

	// migrate the schema
	db.AutoMigrate(
		&config.Config{},
		&message.Message{},
		&session.Session{},
	)

	// seed the database
	if err := session.Seed(db); err != nil {
		fmt.Fprintf(os.Stderr, "failed to seed database: %v\n", err)
		os.Exit(1)
	}

	// get the name of the current session
	var sessName string
	var conf config.Config
	result := db.Where("name = ?", "session").Limit(1).Find(&conf)
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to load session config: %v\n", result.Error)
		os.Exit(1)
	}

	if result.RowsAffected == 0 {
		sessName = "assistant"
	} else {
		sessName = conf.Value
	}

	// get the current session
	var sess session.Session
	result = db.Preload("Messages").
		Where("name = ?", sessName).
		Limit(1).
		Find(&sess)
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to load current session: %v\n", result.Error)
		os.Exit(1)
	}
	// read and validate `<num>` if passed
	// TODO: properly compute the log read-behind
	num := math.MaxInt32
	if opts["<num>"] != nil {
		var err error
		num, err = strconv.Atoi(opts["<num>"].(string))
		if err != nil || num < 1 {
			fmt.Fprintf(os.Stderr, "invalid value for <num>: %v\n", err)
			os.Exit(1)
		}
	}

	sess.Messages = sess.Limit(num)

	// determine which command to execute
	var cmd func(map[string]interface{}, session.Session, *gorm.DB)

	// DRY out some of the command-routing logic
	cmdQ := opts["q"].(bool) || opts["question"].(bool)
	cmdS := opts["s"].(bool) || opts["session"].(bool)

	switch {
	//case opts["config"].(bool) && opts["list"].(bool):
	//cmd = cmdConfigList

	//case opts["config"].(bool) && opts["set"].(bool):
	//cmd = cmdConfigSet

	case opts["db"].(bool) && opts["destroy"].(bool):
		cmd = cmdDBDestroy

	case opts["db"].(bool) && opts["path"].(bool):
		cmd = cmdDBPath

	case cmdQ && opts["add"].(bool):
		cmd = cmdQuestionAdd

	case cmdQ && opts["ask"].(bool):
		cmd = cmdQuestionAsk

	case cmdQ && opts["rm"].(bool):
		cmd = cmdQuestionRemove

	case cmdQ && opts["send"].(bool):
		cmd = cmdQuestionSend

	case cmdS && opts["cp"].(bool):
		cmd = cmdSessionCopy

	case cmdS && opts["edit"].(bool):
		cmd = cmdSessionEdit

	case cmdS && opts["log"].(bool):
		cmd = cmdSessionLog

	case cmdS && opts["ls"].(bool):
		cmd = cmdSessionList

	case cmdS && opts["new"].(bool):
		cmd = cmdSessionNew

	case cmdS && opts["rm"].(bool):
		cmd = cmdSessionRemove

	case cmdS && opts["search"].(bool):
		cmd = cmdSessionSearch

	case cmdS && opts["use"].(bool):
		cmd = cmdSessionUse

	case cmdS && opts["view"].(bool):
		cmd = cmdSessionView

	default:
		fmt.Println(optstring)
		os.Exit(0)
	}

	// execute the command
	cmd(opts, sess, db)
}
