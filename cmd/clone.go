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
	"path"
	"strings"

	"dandyrow/git-worktree/internal/git"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <repository> [<directory>]",
	Short: "Clones the specified git repo as a blank repository for use with worktrees",
	Long:  `Clones the specified git repo as a blank repository for use with worktrees.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoURL := args[0]

		directory := strings.TrimSuffix(path.Base(repoURL), ".git")
		if len(args) > 1 {
			directory = args[1]
		}

		return git.CloneBareRepo(repoURL, directory)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
