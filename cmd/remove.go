/*
Package cmd implements the command line interface the user interacts with.

The package is built using the Cobra CLI library and organizes commands
in a hierarchical structure. Each command is implemented in its own file
and registered with the root command during package initialization.

# Copyright © 2025 Daniel Lowry <development@daniellowry.co.uk>

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

var removeCmd = &cobra.Command{
	Use:   "remove <worktree>",
	Short: "Removes a worktree and it's associated branch",
	Long:  `Removes a worktree and it's associated branch.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		worktreeName := args[0]
		force, _ := cmd.Flags().GetBool("force")

		branchName, err := git.GetWorktreeBranch(worktreeName)
		if err != nil {
			return fmt.Errorf("failed to get branch for worktree: %w", err)
		}

		if err := git.RemoveWorktree(worktreeName, force); err != nil {
			return err
		}

		if branchName != "" {
			if err := git.RemoveBranch(branchName, force); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	removeCmd.Flags().BoolP("force", "f", false, "force removal of the worktree and associated branch")
	rootCmd.AddCommand(removeCmd)
}
