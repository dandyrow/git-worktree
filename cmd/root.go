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
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-worktree",
	Short: "A CLI to make Git Worktrees more dev friendly",
	Long:  `A CLI to make Git Worktrees more dev friendly`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
