package tools

import "github.com/YitzhakMizrahi/bootstrap-cli/internal/install"

// CoreTools returns a list of essential development tools
func CoreTools() []*install.Tool {
	return []*install.Tool{
		{
			Name:        "Git",
			PackageName: "git",
			PackageNames: &install.PackageMapping{
				Default: "git",
				APT:     "git",
				DNF:     "git",
				Pacman:  "git",
				Brew:    "git",
			},
			VerifyCommand: "git --version",
		},
		{
			Name:        "cURL",
			PackageName: "curl",
			PackageNames: &install.PackageMapping{
				Default: "curl",
				APT:     "curl",
				DNF:     "curl",
				Pacman:  "curl",
				Brew:    "curl",
			},
			VerifyCommand: "curl --version",
		},
		{
			Name:        "Wget",
			PackageName: "wget",
			PackageNames: &install.PackageMapping{
				Default: "wget",
				APT:     "wget",
				DNF:     "wget",
				Pacman:  "wget",
				Brew:    "wget",
			},
			VerifyCommand: "wget --version",
		},
		{
			Name:        "Build Essential",
			PackageName: "build-essential",
			PackageNames: &install.PackageMapping{
				Default: "build-essential",
				APT:     "build-essential",
				DNF:     "gcc gcc-c++ make",
				Pacman:  "base-devel",
				Brew:    "gcc make",
			},
			VerifyCommand: "gcc --version && make --version",
		},
		{
			Name:        "Ripgrep",
			PackageName: "ripgrep",
			PackageNames: &install.PackageMapping{
				Default: "ripgrep",
				APT:     "ripgrep",
				DNF:     "ripgrep",
				Pacman:  "ripgrep",
				Brew:    "ripgrep",
			},
			VerifyCommand: "rg --version",
		},
		{
			Name:        "Bat",
			PackageName: "bat",
			PackageNames: &install.PackageMapping{
				Default: "bat",
				APT:     "bat",
				DNF:     "bat",
				Pacman:  "bat",
				Brew:    "bat",
			},
			VerifyCommand: "bat --version",
		},
		{
			Name:        "lsd",
			PackageName: "lsd",
			PackageNames: &install.PackageMapping{
				Default: "lsd",
				APT:     "lsd",
				DNF:     "lsd",
				Pacman:  "lsd",
				Brew:    "lsd",
			},
			VerifyCommand: "lsd --version",
		},
		{
			Name:        "fzf",
			PackageName: "fzf",
			PackageNames: &install.PackageMapping{
				Default: "fzf",
				APT:     "fzf",
				DNF:     "fzf",
				Pacman:  "fzf",
				Brew:    "fzf",
			},
			VerifyCommand: "fzf --version",
		},
		{
			Name:        "htop",
			PackageName: "htop",
			PackageNames: &install.PackageMapping{
				Default: "htop",
				APT:     "htop",
				DNF:     "htop",
				Pacman:  "htop",
				Brew:    "htop",
			},
			VerifyCommand: "htop --version",
		},
		{
			Name:        "tree",
			PackageName: "tree",
			PackageNames: &install.PackageMapping{
				Default: "tree",
				APT:     "tree",
				DNF:     "tree",
				Pacman:  "tree",
				Brew:    "tree",
			},
			VerifyCommand: "tree --version",
		},
	}
} 