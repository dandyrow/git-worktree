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

		return cloneBareRepo(repoURL, directory)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

func cloneBareRepo(url string, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := gitCommand("", "clone", "--bare", "--single-branch", url, dir+"/.git"); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	const remoteFetchConfig string = "+refs/heads/*:refs/remotes/origin/*"
	if err := gitCommand(dir, "config", "remote.origin.fetch", remoteFetchConfig); err != nil {
		return fmt.Errorf("failed to configure remote fetch: %w", err)
	}

	if err := gitCommand(dir, "fetch", "--quiet"); err != nil {
		return fmt.Errorf("failed to fetch remote branches: %w", err)
	}

	if err := setupUpstreamTracking(dir); err != nil {
		return fmt.Errorf("failed to setup upstream tracking: %w", err)
	}

	return nil
}

func setupUpstreamTracking(dir string) error {
	cmd := exec.Command("git", "for-each-ref", "--format=%(refname:short)", "refs/heads")
	cmd.Dir = dir

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	branches := strings.Fields(string(output))

	for _, branch := range branches {
		if branch == "" {
			continue
		}

		upstreamBranch := "origin/" + branch
		if err := gitCommand(dir, "branch", "--set-upstream-to="+upstreamBranch, branch); err != nil {
			return fmt.Errorf("failed to set upstream for branch %s: %w", branch, err)
		}
	}

	return nil
}

func gitCommand(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if dir != "" {
		cmd.Dir = dir
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git command failed: %w", err)
	}

	return nil
}
