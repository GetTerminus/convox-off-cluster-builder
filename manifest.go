package convox_off_cluster_builder

import "github.com/convox/rack/manifest1"

type Manifest interface {
	Processes() []string
	ProcessHasBuild(process string) bool
	ProcessBuildPath(process string) (string, error)
}

type V1Manifest manifest1.Manifest

func LoadFileV1(path string) (Manifest, error) {
	man, err := manifest1.LoadFile(path)

	if err != nil {
		return nil, err
	}

	v1man := V1Manifest(*man) //deref and cast the struct to my owned copy of the type
	return &v1man, nil
}

func (m *V1Manifest) Processes() []string {
	procs := []string{}
	for k, _ := range m.Services {
		procs = append(procs, k)
	}

	return procs
}

func (m *V1Manifest) ProcessHasBuild(process string) bool {
	build := m.Services[process].Build

}
