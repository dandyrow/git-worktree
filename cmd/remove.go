/*
Package cmd implements the command line interface the user interacts with.

The package is built using the Cobra CLI library and organizes commands
in a hierarchical structure. Each command is implemented in its own file
and registered with the root command during package initialization.

Copyright © 2025 Daniel Lowry <development@daniellowry.co.uk>

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
	"os"
	"os/exec"
	"path"
	"strings"

	"dandyrow/git-worktree/internal/git"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <worktree>",
	Short: "Removes a worktree and it's associated branch",
	Long:  `Removes a worktree and it's associated branch.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			force = false
		}

		branchName, err := getBranchForWorktree(name)
		if err != nil {
			return fmt.Errorf("failed to get branch name for worktree: %w", err)
		}

		worktreeCmdArgs := []string{"worktree", "remove", name}
		if force {
			worktreeCmdArgs = append(worktreeCmdArgs, "-f")
		}

		if err := git.Command("", worktreeCmdArgs...); err != nil {
			return fmt.Errorf("failed to remove git worktree: %w", err)
		}

		if branchName == "" {
			return nil
		}

		branchCmdArgs := []string{"branch", "-d", branchName}
		if force {
			branchCmdArgs = append(branchCmdArgs, "--force")
		}

		if err := git.Command("", branchCmdArgs...); err != nil {
			return fmt.Errorf("failed to remove git branch. Branch will need removed manually: %w", err)
		}

		return nil
	},
}

func init() {
	removeCmd.Flags().BoolP("force", "f", false, "force removal of the worktree and associated branch")
	rootCmd.AddCommand(removeCmd)
}

func getBranchForWorktree(name string) (string, error) {
	worktreeListCmd := exec.Command("git", "worktree", "list", "--porcelain")
	worktreeListCmd.Stderr = os.Stderr
	worktreeList, err := worktreeListCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list worktrees: %w", err)
	}

	for block := range strings.SplitSeq(string(worktreeList), "\n\n") {
		lines := strings.Split(strings.TrimSpace(block), "\n")

		if len(lines) == 0 {
			continue
		}

		worktreePath, found := strings.CutPrefix(lines[0], "worktree ")
		if !found {
			continue
		}

		outputName := path.Base(worktreePath)

		if outputName != name {
			continue
		}

		if strings.HasPrefix(lines[2], "detached") {
			return "", nil
		}

		branchRef := strings.TrimPrefix(lines[2], "branch ")
		return path.Base(branchRef), nil
	}

	return "", fmt.Errorf("failed to find worktree %s", name)
}
