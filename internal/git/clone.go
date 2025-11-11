/*
Package git contains utility functions for working with git

This is an internal package intended only for use within this project.

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
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const directoryPermission os.FileMode = 0o755

// CloneBareRepo clones a git repository as a bare
// repository into the specified directory and
// configures upstream tracking for all branches.
//
// Returns an error if the clone or configuration
// fails.
func CloneBareRepo(url string, dir string) error {
	if err := os.MkdirAll(dir, directoryPermission); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := Command("", "clone", "--bare", "--single-branch", url, dir+"/.git"); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	const remoteFetchConfig string = "+refs/heads/*:refs/remotes/origin/*"
	if err := Command(dir, "config", "remote.origin.fetch", remoteFetchConfig); err != nil {
		return fmt.Errorf("failed to configure remote fetch: %w", err)
	}

	if err := Command(dir, "fetch", "--quiet"); err != nil {
		return fmt.Errorf("failed to fetch remote branches: %w", err)
	}

	if err := setupUpstreamTracking(dir); err != nil {
		return fmt.Errorf("failed to setup upstream tracking: %w", err)
	}

	return nil
}

// setupUpstreamTracking configures upstream tracking
// for all local branches in the repository located
// in dir as git does it on normal repository clones.
func setupUpstreamTracking(dir string) error {
	cmd := exec.Command("git", "for-each-ref", "--format=%(refname:short)", "refs/heads")
	cmd.Dir = dir

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	for branch := range strings.FieldsSeq(string(output)) {
		if branch == "" {
			continue
		}

		upstreamBranch := "origin/" + branch
		if err := Command(dir, "branch", "--set-upstream-to="+upstreamBranch, branch); err != nil {
			return fmt.Errorf("failed to set upstream for branch %s: %w", branch, err)
		}
	}

	return nil
}
