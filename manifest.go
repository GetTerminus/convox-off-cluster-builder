package convox_off_cluster_builder

import "github.com/convox/rack/manifest1"

type Manifest interface {
	Processes() []string
}

func LoadFile(path string) (Manifest, error) {
	return manifest1.LoadFile(path)
}
