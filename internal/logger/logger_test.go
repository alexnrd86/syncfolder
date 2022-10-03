package logger

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetLogLevel(t *testing.T) {
	req := require.New(t)

	cases := map[string]struct {
		logLevel string
		isError  bool
		errMsg   string
	}{
		"success INFO": {
			logLevel: "INFO",
		},

		"success ERROR": {
			logLevel: "ERROR",
		},

		"success CRITICAL": {
			logLevel: "CRITICAL",
		},

		"wrong log level": {
			logLevel: "SOMETHING",
			isError:  true,
			errMsg:   "log level is not set in config. Default log level ERROR",
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {

			err := SetLogLevel(cs.logLevel)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errMsg)
			} else {
				req.NoError(err)
			}
		})
	}

}

func TestLog(t *testing.T) {

	req := require.New(t)
	root, _ := filepath.Abs("../../")

	cases := map[string]struct {
		path    string
		isError bool
	}{
		"success": {

			path: root + "/test/logger/testlog.txt",
		},

		"empty path": {
			path:    "",
			isError: true,
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {

			LogPath = cs.path
			err := log("TEST", "testfunc", "test message")

			if cs.isError {
				req.Error(err)

			} else {
				req.NoError(err)
			}
		})
	}

}
