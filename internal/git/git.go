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
