package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCloneCmnd(t *testing.T) {
	tests := []struct {
		name, repository, url, directory string
	}{
		{
			name:       "should create a clone of a repository into a newly created directory named after the repository",
			repository: "git-gong",
			url:        "https://github.com/erikjuhani/git-gong.git",
			directory:  "",
		},
		{
			name:       "should create a clone of a repository into a directory",
			repository: "git-gong",
			url:        "https://github.com/erikjuhani/git-gong.git",
			directory:  "gong-git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "gong-clone")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			if err := os.Chdir(dir); err != nil {
				t.Fatal(err)
			}

			args := []string{cloneCmd.Name(), tt.url}

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			clone := tt.repository
			if tt.directory != "" {
				clone = tt.directory
			}

			if _, err := os.Stat(clone); err != nil {
				if os.IsNotExist(err) {
					t.Fatal(err)
				}
			}
		})
	}
}
