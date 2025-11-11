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
	"path"
	"strings"
)

type worktreeInfo struct {
	Path     string
	Name     string
	Branch   string
	Detached bool
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

	worktrees, err := parseWorktreeList(output)
	if err != nil {
		return "", fmt.Errorf("failed to parse worktree list: %w", err)
	}

	worktree, err := findWorktreeByName(worktrees, worktreeName)
	if err != nil {
		return "", err
	}

	if worktree.Detached {
		return "", nil
	}

	return worktree.Branch, nil
}

func parseWorktreeList(list string) ([]worktreeInfo, error) {
	var worktrees []worktreeInfo

	for block := range strings.SplitSeq(strings.TrimSpace(list), "\n\n") {
		if strings.TrimSpace(block) == "" {
			continue
		}

		info, err := parseWorktreeBlock(block)
		if err != nil {
			continue
		}

		worktrees = append(worktrees, *info)
	}

	return worktrees, nil
}

func parseWorktreeBlock(block string) (*worktreeInfo, error) {
	lines := strings.Split(strings.TrimSpace(block), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty worktree block")
	}

	worktreePath, found := strings.CutPrefix(lines[0], "worktree ")
	if !found {
		return nil, fmt.Errorf("invalid worktree format: missing 'worktree' prefix")
	}

	info := &worktreeInfo{
		Path: worktreePath,
		Name: path.Base(worktreePath),
	}

	branchLine := lines[2]
	if strings.HasPrefix(branchLine, "detached") {
		info.Detached = true
		return info, nil
	}

	if strings.HasPrefix(branchLine, "branch") {
		branchRef := strings.TrimPrefix(branchLine, "branch ")
		info.Branch = path.Base(branchRef)
		return info, nil
	}

	return nil, fmt.Errorf("no branch information found for worktree %s", info.Name)
}

func findWorktreeByName(worktrees []worktreeInfo, name string) (*worktreeInfo, error) {
	for _, worktree := range worktrees {
		if worktree.Name == name {
			return &worktree, nil
		}
	}

	return nil, fmt.Errorf("worktree %s not found", name)
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
