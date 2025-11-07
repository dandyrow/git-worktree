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
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
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

		err := os.MkdirAll(directory, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", directory, err)
		}

		gitCommand := exec.Command("git", "clone", "--bare", "--single-branch", repoURL, directory+"/.git")
		gitCommand.Stdout = os.Stdout
		gitCommand.Stderr = os.Stderr

		err = gitCommand.Run()
		if err != nil {
			return fmt.Errorf("git command failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
