package internal

import "embed"

//go:embed defaults/config/tools/schema.yaml
//go:embed defaults/config/tools/modern/*.yaml
//go:embed defaults/config/dotfiles/schema.yaml
//go:embed defaults/config/dotfiles/shell/*.yaml
//go:embed defaults/config/fonts/schema.yaml
//go:embed defaults/config/fonts/monospace/*.yaml
//go:embed defaults/config/languages/schema.yaml
//go:embed defaults/config/languages/interpreted/*.yaml
var configFiles embed.FS 