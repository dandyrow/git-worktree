/*
Package cmd implements the command line interface the user interacts with.

The package is built using the Cobra CLI library and organizes commands
in a hierarchical structure. Each command is implemented in its own file
and registered with the root command during package initialization.

# Copyright © 2025 Daniel Lowry <devlopment@daniellowry.co.uk>

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
	"fmt"

	"dandyrow/git-worktree/internal/git"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <branch name> [<base branch>] [flags]",
	Short: "Add the specified branch as a worktree",
	Long: `Adds the specified branch as a worktree.

The base branch optional argument causes a new branch to be created based upon the specified base branch.`,
	Args: cobra.RangeArgs(1, 2),

	RunE: func(cmd *cobra.Command, args []string) error {
		branchName := args[0]

		path, err := constructPath(cmd, branchName)
		if err != nil {
			return err
		}

		if len(args) == 1 {
			return runWorktreeAdd(path, branchName, "")
		}

		baseBranch := args[1]
		return runWorktreeAdd(path, baseBranch, branchName)
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

// Adds a git worktree to the path based on the commitIsh.
//
// The newBranchName paramter is optional. If set a new branch will be created
// branched off the commitIsh and will be checked out in the created worktree.
// If set to the empty string, the worktree will be created with the commitIsh
// checked out.
func runWorktreeAdd(path string, commitIsh string, newBranchName string) error {
	gitArgs := []string{"worktree", "add", path, commitIsh}
	if newBranchName != "" {
		gitArgs = append(gitArgs, "-b", newBranchName)
	}

	if err := git.Command("", gitArgs...); err != nil {
		return fmt.Errorf("failed to add git worktree at %s: %w", path, err)
	}

	return nil
}
