/*
Git-worktree is a CLI for working with git worktrees.

It simplifies working with git worktrees by providing commands to:
- Clone repositories as bare repositories suitable for worktrees
- Add new worktrees for existing remote branches
- Add new worktrees containing a new branch off of an existing branch
- Manage multiple working directories from a single git repository

# Copyright © 2025 Daniel Lowry

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
package main

import "dandyrow/git-worktree/cmd"

func main() {
	cmd.Execute()
}
