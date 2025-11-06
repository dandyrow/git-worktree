/*
Copyright © 2025 Daniel Lowry

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <branch name>",
	Short: "Add the specified branch as a worktree",
	Long:  `Adds the specified branch as a worktree`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var stderr bytes.Buffer
		branchName := args[0]

		path, err := constructPath(cmd, branchName)
		if err != nil {
			return err
		}

		gitCommand := exec.Command("git", "worktree", "add", path, branchName)
		gitCommand.Stderr = &stderr
		err = gitCommand.Run()
		if err != nil {
			return fmt.Errorf("git error: %v", stderr.String())
		}

		return nil
	},
}

func init() {
	addCmd.Flags().StringP("target-directory", "d", "./", "target directory where the worktree should be created")
	addCmd.Flags().StringP("name", "n", "", "set name of worktree to something different than branch name")
	rootCmd.AddCommand(addCmd)
}

func constructPath(cmd *cobra.Command, branchName string) (string, error) {
	targetDirectory, err := cmd.Flags().GetString("target-directory")
	if err != nil {
		return "", fmt.Errorf("constructing path: %v", err)
	}

	worktreeName, err := cmd.Flags().GetString("name")
	if err != nil {
		return "", fmt.Errorf("constructing path: %v", err)
	}

	if worktreeName == "" {
		worktreeName = branchName
	}

	return targetDirectory + worktreeName, nil
}
