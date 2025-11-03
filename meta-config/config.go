package metaconfig

import (
	metamodule "meta/meta-module"
	metanode "meta/meta-node"
)

type Config struct {
	Module metamodule.Config `yaml:"module"`
	Node   metanode.Config   `yaml:"node"`
}
