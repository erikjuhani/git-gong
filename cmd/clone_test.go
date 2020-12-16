package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCloneCmd(t *testing.T) {
	r, err := createTestRepo()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupTestRepo(r)

	parts := strings.Split(r.GitPath, "/")
	repoName := fmt.Sprintf("%s/%s", parts[len(parts)-3], strings.TrimSuffix(parts[len(parts)-1], ".git"))

	tests := []struct {
		name, repository, url, directory string
	}{
		{
			name:       "should create a clone of a repository into a newly created directory named after the repository",
			repository: repoName,
			url:        r.GitPath,
			directory:  "",
		},
		{
			name:       "should create a clone of a repository into a directory",
			repository: repoName,
			url:        r.GitPath,
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

			clone := tt.repository
			if tt.directory != "" {
				args = append(args, tt.directory)
				clone = tt.directory
			}

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			if _, err := os.Stat(clone); err != nil {
				if os.IsNotExist(err) {
					t.Fatal(err)
				}
			}
		})
	}
}
