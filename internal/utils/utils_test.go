package utils

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {

	root, _ := filepath.Abs("../../")

	req := require.New(t)

	cases := map[string]struct {
		configPath string
		isError    bool
		res        map[string]string
	}{
		"wrong path": {
			configPath: root + "/config.txt",
			isError:    true,
		},

		"empty path": {
			configPath: "",
			isError:    true,
		},
		"success": {
			configPath: root + "/test/config/test1",
			isError:    true,
			res: map[string]string{
				"sourcepath": "c:/temp/master",
				"synchpath":  "c:/temp/slave",
				"loglevel":   "INFO",
			},
		},

		"values with spaces": {
			configPath: root + "/test/config/test2",
			isError:    true,
			res: map[string]string{
				"sourcepath": "c:/temp/master",
				"synchpath":  "c:/temp/slave",
				"loglevel":   "INFO",
			},
		},
	}

	for name, cs := range cases {

		t.Run(name, func(t *testing.T) {

			ConfigPath = cs.configPath

			fmt.Println(ConfigPath)

			res, err := GetConfig()

			if cs.isError {
				req.Error(err)
			} else {
				req.NoError(err)
				req.Equal(cs.res, res)
			}

		})

	}

}
