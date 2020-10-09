package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestInitCmd(t *testing.T) {
	dir, err := ioutil.TempDir("", "gong-init")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	tests := []struct {
		name    string
		flag    string
		dirName string
	}{
		{
			name:    "should create empty Git repository to current working directory",
			flag:    "",
			dirName: "",
		},
		{
			name:    "should create empty Git repository to <directory> and create the directory if it does not exist",
			flag:    "",
			dirName: "gong",
		},
		{
			name:    "should create an empty Git repository to (current working) <directory> with initial-branch name as <branchname>.",
			flag:    "default-branch",
			dirName: "gong",
		},
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.flag != "" {
				initCmd.Flags().Set(tt.flag, "dev")
			}

			err := initCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			gitDir := ".git"
			expectedPaths := []string{"HEAD", "objects", "refs/heads", "refs/tags"}

			if flag := initCmd.Flags().Lookup("default-branch"); flag != nil && flag.Value.String() != "" {
				expectedPaths = []string{"refs/heads/dev"}
			} else {
				expectedPaths = []string{"refs/heads/main"}
			}

			for _, p := range expectedPaths {
				generatedPath := fmt.Sprintf("%s/%s/%s", tt.dirName, gitDir, p)
				if _, err := os.Stat(generatedPath); err != nil {
					if os.IsNotExist(err) {
						t.Fatal(err)
					}
				}
			}
		})
	}
}
