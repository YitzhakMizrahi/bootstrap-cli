package internal

import "embed"

//go:embed config/defaults/tools/modern/*.yaml
//go:embed config/defaults/tools/schema.yaml
//go:embed config/defaults/dotfiles/schema.yaml
//go:embed config/defaults/dotfiles/shell/*.yaml
var configFiles embed.FS 