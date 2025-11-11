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
	"path"
	"strings"
)

// Command executes a git command in the directory
// with the provided args.
//
// The directory parameter is optional. If set to
// the empty string the command will be run in the
// current working directory.
//
// Stdout and Stderr will be printed to the os provided
// stdout and stderr.
// Returns an error if the git command fails.
func Command(directory string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if directory != "" {
		cmd.Dir = directory
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git command failed: %w", err)
	}

	return nil
}

// CommandOutput executes a git command and returns its output.
//
// The directory parameter is optional. If set to
// the empty string the command will be run in the
// current working directory.
//
// Stderr will be printed to the os provided stderr.
// Returns the stdout output as a string.
func CommandOutput(directory string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	if directory != "" {
		cmd.Dir = directory
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git command failed: %w", err)
	}

	return string(output), nil
}

// CloneBareRepo clones a git repository as a bare
// repository into the specified directory and
// configures upstream tracking for all branches.
//
// Returns an error if the clone or configuration
// fails.
func CloneBareRepo(url string, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
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

// GetWorktreeBranch returns the branch name associated with a worktree.
//
// Returns an empty string if the worktree is in a detached HEAD state.
// Returns an error if the worktree cannot be found or if the git command fails.
func GetWorktreeBranch(worktreeName string) (string, error) {
	output, err := CommandOutput("", "worktree", "list", "--porcelain")
	if err != nil {
		return "", fmt.Errorf("failed to list worktrees: %w", err)
	}

	for block := range strings.SplitSeq(strings.TrimSpace(output), "\n\n") {
		lines := strings.Split(strings.TrimSpace(block), "\n")

		if len(lines) == 0 {
			continue
		}

		worktreePath, found := strings.CutPrefix(lines[0], "worktree ")
		if !found {
			continue
		}

		currentName := path.Base(worktreePath)
		if currentName != worktreeName {
			continue
		}

		if len(lines) < 3 {
			return "", fmt.Errorf("unexpected worktree format for %s", worktreeName)
		}

		if strings.HasPrefix(lines[2], "detached") {
			return "", nil // Detached HEAD state
		}

		branchRef := strings.TrimPrefix(lines[2], "branch ")
		return path.Base(branchRef), nil
	}

	return "", fmt.Errorf("worktree %s not found", worktreeName)
}

// AddWorktree adds a git worktree to the specified path in the
// filesystem basing it on the specified commitIsh.
//
// A new branch will be created named after newBranchName and will
// be checked out in the worktree.
//
// If newBranchName is set to the empty string no new branch will
// be created and the commitIsh will be checked out in the worktree.
//
// Returns an error if the worktree cannot be added or the git
// command fails.
func AddWorktree(path string, commitIsh string, newBranchName string) error {
	gitArgs := []string{"worktree", "add", path, commitIsh}
	if newBranchName != "" {
		gitArgs = append(gitArgs, "-b", newBranchName)
	}

	if err := Command("", gitArgs...); err != nil {
		return fmt.Errorf("failed to add git worktree at %s: %w", path, err)
	}

	return nil
}

// RemoveWorktree removes a git worktree from the filesystem.
//
// Returns an error if the worktree cannot be deleted or the
// git command fails.
func RemoveWorktree(worktreeName string, force bool) error {
	args := []string{"worktree", "remove", worktreeName}
	if force {
		args = append(args, "-f")
	}

	if err := Command("", args...); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	return nil
}

// RemoveBranch removes a git branch from the local repository.
//
// Returns an error if the branch cannot be deleted or the git
// command fails.
func RemoveBranch(branchName string, force bool) error {
	args := []string{"branch", "-d", branchName}
	if force {
		args = []string{"branch", "-D", branchName}
	}

	if err := Command("", args...); err != nil {
		return fmt.Errorf("failed to remove branch (branch will need deleted manually): %w", err)
	}

	return nil
}
