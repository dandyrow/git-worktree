# git-worktree

A developer-friendly CLI tool that simplifies working with Git worktrees, making it easier to manage multiple working directories from a single Git repository.

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.25.3-00ADD8?logo=go)](https://go.dev/)

If you find this tool useful, please consider giving it a ⭐ star on GitHub!

## Overview

Git worktrees allow you to check out multiple branches simultaneously in different directories, all linked to a single Git repository. This is incredibly useful for:

- Working on multiple features or bug fixes in parallel
- Quickly switching between branches without losing work-in-progress
- Running different versions of your code simultaneously for comparison or testing
- Avoiding the overhead of multiple repository clones

**git-worktree** makes working with Git worktrees effortless by providing intuitive commands to clone repositories as bare repos (optimized for worktrees), add new worktrees, and manage them efficiently.

## Features

- 🚀 **Easy Setup**: Clone repositories configured for worktrees with a single command
- 🌿 **Branch Management**: Add worktrees for existing branches or create new ones
- 🗑️ **Clean Removal**: Remove worktrees and their associated branches safely
- ⚙️ **Flexible Options**: Customize worktree names and target directories
- 🛠️ **Built on Git**: Works seamlessly with your existing Git workflow

## Installation

### Pre-built Binaries

Pre-built binaries will be available in the [Releases](../../releases) section (coming soon).

### From Source

#### Prerequisites

- Go 1.25.3 or later
- Git installed and available in your PATH

#### Build and Install

```bash
# Clone the repository
git clone https://github.com/dandyrow/git-worktree
cd git-worktree

# Build the binary
go build -o git-worktree main.go

# Optionally, install to your GOPATH/bin
go install
```

The binary will be available as `git-worktree`.

## Usage

### Clone a Repository for Worktrees

Clone a repository as a bare repository, which is ideal for managing multiple worktrees:

```bash
git-worktree clone <repository-url> [directory]
```

**Examples:**

```bash
# Clone to a directory named after the repository
git-worktree clone https://github.com/user/repo.git

# Clone to a custom directory
git-worktree clone https://github.com/user/repo.git my-project
```

This creates a bare repository structure optimized for working with worktrees.

### Add a Worktree

Add a worktree for an existing branch or create a new branch:

```bash
git-worktree add <branch-name> [base-branch] [flags]
```

**Examples:**

```bash
# Add a worktree for an existing remote branch
git-worktree add feature/new-feature

# Create a new branch based on another branch
git-worktree add feature/my-branch main

# Add a worktree with a custom directory name
git-worktree add feature/long-branch-name --name custom-name

# Add a worktree in a specific directory
git-worktree add feature/api-v2 --target-directory ~/projects/worktrees
```

**Flags:**

- `-d, --target-directory <path>`: Target directory where the worktree should be created (default: current directory)
- `-n, --name <name>`: Set a custom name for the worktree directory (default: branch name)

### Remove a Worktree

Remove a worktree and its associated branch:

```bash
git-worktree remove <worktree-name> [flags]
```

**Examples:**

```bash
# Remove a worktree
git-worktree remove feature/old-feature

# Force removal even if there are uncommitted changes
git-worktree remove feature/old-feature --force
```

**Flags:**

- `-f, --force`: Force removal of the worktree and associated branch, even with uncommitted changes

### Get Help

```bash
# General help
git-worktree --help

# Command-specific help
git-worktree add --help
git-worktree clone --help
git-worktree remove --help
```

## Typical Workflow

Here's a typical workflow using git-worktree:

```bash
# 1. Clone your repository for worktree usage
git-worktree clone https://github.com/user/my-project.git
cd my-project

# 2. Add a worktree for the main branch
git-worktree add main

# 3. Create a new feature branch and worktree
git-worktree add feature/new-api main

# 4. Work in different directories simultaneously
cd main           # Work on main branch
cd ../feature/new-api  # Work on feature branch

# 5. When done, remove the feature worktree
cd ..
git-worktree remove feature/new-api
```

## FAQ

### Why use worktrees instead of cloning the repository multiple times?

Worktrees share the same `.git` directory, which means:
- Less disk space used (no duplicate objects)
- Faster branch switching
- Shared configuration and hooks
- Single source of truth for all branches

### Can I use this with an existing Git repository?

Yes! You can use `git-worktree add` in any existing Git repository to add worktrees. The `git-worktree clone` command is just a convenience for starting fresh with a bare repository structure.

### What's the difference between a bare repository and a normal repository?

A bare repository doesn't have a working directory - it only contains the Git database. This makes it ideal for worktrees because each worktree gets its own working directory, and the bare repository serves as the central storage for all of them.

### Can I use regular git commands after using git-worktree?

Absolutely! `git-worktree` is just a wrapper around standard Git commands. You can use `git worktree`, `git branch`, and any other Git commands as usual.

## Project Structure

```
git-worktree/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command and CLI setup
│   ├── clone.go           # Clone command
│   ├── add.go             # Add command
│   └── remove.go          # Remove command
├── internal/
│   └── git/               # Git operations
│       ├── clone.go       # Clone functionality
│       ├── worktree.go    # Worktree management
│       └── commandWrapper.go  # Git command execution
├── main.go                # Application entry point
├── go.mod                 # Go module definition
└── README.md              # This file
```

## Contributing

Contributions are welcome! Here's how you can help:

1. **Report Bugs**: Open an issue describing the bug and how to reproduce it
2. **Suggest Features**: Open an issue describing the feature and its use case
3. **Submit Pull Requests**:
   - Fork the repository
   - Create a feature branch (`git checkout -b feature/amazing-feature`)
   - Commit your changes (`git commit -m 'Add amazing feature'`)
   - Push to the branch (`git push origin feature/amazing-feature`)
   - Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Maintain the existing code style
- Add comments for exported functions and complex logic
- Test your changes thoroughly

### Building from Source

#### Development Build

```bash
# Clone the repository
git clone <repository-url>
cd git-worktree/develop

# Install dependencies
go mod download

# Build
go build -o git-worktree main.go

# Run
./git-worktree --help
```

#### Production Build

For a production-ready build with optimizations:

```bash
go build -ldflags="-s -w" -o git-worktree main.go
```

This will produce a smaller, optimized binary.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

This means you are free to:
- Use this software for any purpose
- Change the software to suit your needs
- Share the software with others
- Share your changes with others

Under the following conditions:
- You must include the original copyright notice and license
- You must make your modified source code available under the same license

## Author

**Daniel Lowry**

- Email: development@daniellowry.co.uk
- Copyright © 2025 Daniel Lowry

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) - A powerful CLI framework for Go
- Inspired by the power and flexibility of Git worktrees
