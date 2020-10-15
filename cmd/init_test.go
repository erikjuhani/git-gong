package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestInitCmd(t *testing.T) {
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
			name:    "should create empty Git repository to specified directory and create the directory if it does not exist",
			flag:    "",
			dirName: "gong",
		},
		{
			name:    "should create an empty Git repository to current working directory with initial-branch name as given branch name.",
			flag:    "default-branch",
			dirName: "gong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "gong-init")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			if err := os.Chdir(dir); err != nil {
				t.Fatal(err)
			}

			args := []string{initCmd.Name()}

			gitDir := ".git"
			expectedPaths := []string{"HEAD", "objects", "refs/heads", "refs/tags"}

			if tt.flag != "" {
				args = append(args, "--default-branch", "dev")
				expectedPaths = append(expectedPaths, "refs/heads/dev")
			} else {
				expectedPaths = append(expectedPaths, "refs/heads/main")
			}

			if tt.dirName != "" {
				gitDir = fmt.Sprintf("%s/%s", tt.dirName, gitDir)
				args = append(args, tt.dirName)
			}

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			for _, p := range expectedPaths {
				generatedPath := fmt.Sprintf("%s/%s/%s", dir, gitDir, p)
				if _, err := os.Stat(generatedPath); err != nil {
					if os.IsNotExist(err) {
						t.Fatal(err)
					}
				}
			}
		})
	}
}
