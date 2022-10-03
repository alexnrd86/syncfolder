package synch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChekMasterFolder(t *testing.T) {

	root, _ := filepath.Abs("../../")

	req := require.New(t)

	cases := map[string]struct {
		masterFolder string
		slavePath    string
		isError      bool
		errMsg       string
	}{
		"Master folder exists": {
			masterFolder: root + "/test/temp/master",
			slavePath:    root + "/test/temp/slave",
			isError:      false,
		},

		"No Master folder": {
			masterFolder: root + "/test/temp2/master",
			slavePath:    root + "/test/temp/slave",
			isError:      true,
			errMsg:       "no such file or directory",
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {

			err := CheckMasterFolder(cs.masterFolder, cs.slavePath)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errMsg)
			} else {
				req.NoError(err)
			}
		})
	}

}

func TestChekFile(t *testing.T) {

	root, _ := filepath.Abs("../../")

	req := require.New(t)

	folder, _ := os.ReadDir(root + "/test/temp/master")
	file := folder[0]

	cases := map[string]struct {
		masterFolder string
		slavePath    string
		isError      bool
		errMsg       string
	}{
		"Slave folder exists": {
			masterFolder: root + "/test/temp/master",
			slavePath:    root + "/test/temp/slave",
			isError:      false,
		},

		"No Slave folder": {
			masterFolder: root + "/test/temp/master",
			slavePath:    root + "/test/temp2/slave",
			isError:      true,
			errMsg:       "no such file or directory",
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {

			err := checkFile(file, cs.masterFolder, cs.slavePath)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errMsg)
			} else {
				req.NoError(err)
			}
		})
	}

}

func TestCopyFile(t *testing.T) {

	root, _ := filepath.Abs("../../")

	req := require.New(t)

	_ = os.Remove(root + "/test/temp/slave/file1")

	cases := map[string]struct {
		inPath  string
		outPath string
		isError bool
		errMsg  string
	}{
		"Master file exists": {
			inPath:  root + "/test/temp/master/file1",
			outPath: root + "/test/temp/slave/file1",
			isError: false,
		},

		"No Master file": {
			inPath:  root + "/test/temp/master/file11",
			outPath: root + "/test/temp/file11",
			isError: true,
			errMsg:  "no such file or directory",
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {

			err := copyFile(cs.inPath, cs.outPath)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errMsg)
			} else {
				req.NoError(err)
				_, err = os.Open(root + "/test/temp/slave/file1")
				req.NoError(err)
			}
		})

	}
	_ = os.Remove(root + "/test/temp/slave/file1")
}

func TestCheckFolder(t *testing.T) {

	root, _ := filepath.Abs("../../")

	req := require.New(t)

	cases := map[string]struct {
		name      string
		slavePath string
		isError   bool
		errMsg    string
	}{
		"Slave folder exists": {
			name:      "testdir",
			slavePath: root + "/test/temp/slave",
			isError:   false,
		},

		"No Slave folder": {
			name:      "testdir",
			slavePath: root + "/test/temp/slave2",
			isError:   true,
			errMsg:    "no such file or directory",
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {

			err := checkFolder(cs.name, cs.slavePath, 0755)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errMsg)
			} else {
				req.NoError(err)
				_, err = os.ReadDir(root + "/test/temp/slave/testdir")
				req.NoError(err)
			}
		})

	}
	_ = os.Remove(root + "/test/temp/slave/testdir")

}

func TestChekSlaveFolder(t *testing.T) {

	root, _ := filepath.Abs("../../")

	req := require.New(t)

	cases := map[string]struct {
		slavePath string
		isError   bool
		errMsg    string
	}{
		"Slave folder exists": {

			slavePath: root + "/test/temp/slave",
			isError:   false,
		},

		"No Slave folder": {

			slavePath: root + "/test/temp2/slave",
			isError:   true,
			errMsg:    "no such file or directory",
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {

			err := CheckSlaveFolder(root+"/test/temp/master", cs.slavePath)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errMsg)
			} else {
				req.NoError(err)
			}
		})
	}

}

func TestRemoveFolder(t *testing.T) {

	root, _ := filepath.Abs("../../")

	req := require.New(t)

	_ = os.Mkdir(root+"/test/temp/master/testdir", 0755)
	_ = os.Mkdir(root+"/test/temp/slave/testdir", 0755)
	_ = os.Mkdir(root+"/test/temp/slave/testdir2", 0755)

	cases := map[string]struct {
		name       string
		masterPath string
		slavePath  string
		isError    bool
		isDeleted  bool
		errMsg     string
	}{
		"Wrong Master Path": {
			name:       "testdir",
			masterPath: root + "/test/temp/wrongmaster",
			slavePath:  root + "/test/temp/slave",
			isError:    true,
			errMsg:     "no such file or directory",
		},

		"Folder exists in master": {
			name:       "testdir",
			masterPath: root + "/test/temp/master",
			slavePath:  root + "/test/temp/slave",
			isError:    false,
			isDeleted:  false,
		},

		"Folder doesn't exists in master": {
			name:       "testdir2",
			masterPath: root + "/test/temp/master",
			slavePath:  root + "/test/temp/slave",
			isError:    false,
			isDeleted:  true,
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {

			deleted, err := removeFolder(cs.name, cs.masterPath, cs.slavePath)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errMsg)
			} else if cs.isDeleted {
				req.NoError(err)
				req.True(deleted)
				_, err = os.ReadDir(cs.slavePath + "/" + name)
				req.Error(err)
				req.Contains(err.Error(), "no such file or directory")

			} else {
				req.NoError(err)
				req.False(deleted)
			}
		})

	}
	_ = os.Remove(root + "/test/temp/slave/testdir")
	_ = os.Remove(root + "/test/temp/master/testdir")

}

func TestDeleteFile(t *testing.T) {

	root, _ := filepath.Abs("../../")

	req := require.New(t)

	var file1, file2 os.DirEntry

	_ = copyFile(root+"/test/temp/master/file1", root+"/test/temp/slave/file1")
	_ = copyFile(root+"/test/temp/master/file1", root+"/test/temp/slave/file2")

	folder, _ := os.ReadDir(root + "/test/temp/slave")
	for _, file := range folder {

		if file.Name() == "file1" {
			file1 = file
		} else {
			file2 = file
		}
	}

	cases := map[string]struct {
		entry      os.DirEntry
		masterPath string
		slavePath  string
		isError    bool
		errMsg     string
	}{
		"Wrong Master Path": {
			entry:      file1,
			masterPath: root + "/test/temp/wrongmaster",
			slavePath:  root + "/test/temp/slave",
			isError:    true,
			errMsg:     "no such file or directory",
		},

		"File exists in master": {
			entry:      file1,
			masterPath: root + "/test/temp/master",
			slavePath:  root + "/test/temp/slave",
			isError:    false,
		},

		"File doesn't exists in master": {
			entry:      file2,
			masterPath: root + "/test/temp/master",
			slavePath:  root + "/test/temp/slave",
			isError:    false,
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {

			err := deleteFile(cs.entry, cs.masterPath, cs.slavePath)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errMsg)
			} else {
				req.NoError(err)

				if cs.entry.Name() == "file1" {
					_, err = os.Open(cs.slavePath + "/" + cs.entry.Name())
					req.NoError(err)
				} else {
					_, err = os.Open(cs.slavePath + "/" + cs.entry.Name())
					req.Error(err)
					req.Contains(err.Error(), "no such file or directory")
				}
			}
		})

	}
	_ = os.Remove(root + "/test/temp/slave/file1")

}

func TestPurgeFolder(t *testing.T) {

	root, _ := filepath.Abs("../../")

	req := require.New(t)

	_ = copyFile(root+"/test/temp/master/file1", root+"/test/temp/slave/file1")

	cases := map[string]struct {
		slavePath string
		isError   bool
		errMsg    string
	}{
		"Wrong Slave Path": {
			slavePath: root + "/test/temp/slave2",
			isError:   true,
			errMsg:    "no such file or directory",
		},

		"File exists in slave": {

			slavePath: root + "/test/temp/slave",
			isError:   false,
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {

			err := purgeFolder(cs.slavePath)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errMsg)
			} else {
				req.NoError(err)

				_, err = os.Open(cs.slavePath + "/file1")
				req.Error(err)
				req.Contains(err.Error(), "no such file or directory")
			}

		})

	}
	_ = os.Remove(root + "/test/temp/slave/file1")

}

func BenchmarkCopyFile(b *testing.B) {

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		root, _ := filepath.Abs("../../")
		inPath := root + "/test/temp/master/file1"
		outPath := root + "/test/temp/slave/file1"
		b.StartTimer()
		_ = copyFile(inPath, outPath)

		b.StopTimer()
		_ = os.Remove(outPath)
		b.StartTimer()
	}

}
