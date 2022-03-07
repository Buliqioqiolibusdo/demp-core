package plugin

import (
	"path"

	"github.com/buliqioqiolibusdo/demp-core/config"
)

const DefaultPluginFsPathBase = "plugins"
const DefaultPluginDirName = "plugins"

var DefaultPluginDirPath = path.Join(config.DefaultConfigDirPath, DefaultPluginDirName)
