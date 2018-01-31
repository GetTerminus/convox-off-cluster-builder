package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/convox/rack/manifest"
	"github.com/convox/rack/manifest1"
	"github.com/convox/rack/structs"
	"gopkg.in/yaml.v2"
)

func main() {

	var composeFileName = flag.String("compose-file", "", "the path to the docker-compose file to build")
	var appNameFlag = flag.String("app", "", "the name of the service")
	var buildIDFlag = flag.String("build-id", "", "the build ID to identify this build")
	var descriptionFlag = flag.String("description", "", "(optional) a description of this build")

	flag.Parse()

	if descriptionFlag != nil {

	}
	if *composeFileName == "" {
		flag.Usage()
		log.Fatal("missing compose-file")
	}

	if *appNameFlag == "" {
		flag.Usage()
		log.Fatal("missing app")
	}

	if *buildIDFlag == "" {
		flag.Usage()
		log.Fatal("missing app")
	}

	data, err := ioutil.ReadFile(*composeFileName)

	if err != nil {
		log.Fatal(err)
	}

	testM, err := manifest.Load(data, manifest.Environment{"UNDERPANTS_CLIENT_ID": "underpants", "UNDERPANTS_CLIENT_SECRET": "unknown", "JWT_SIGNING_SECRET": "unknown"})
	if testM != nil {

	}
	m, err := manifest1.LoadFile(*composeFileName)

	if err != nil {
		log.Fatal(err)
	}

	appName := *appNameFlag
	buildID := *buildIDFlag
	description := *descriptionFlag

	output := manifest1.NewOutput(false)

	testOpt := manifest.BuildOptions{}
	testM.Build("local-build", appName, testOpt)
	opt := manifest1.BuildOptions{}
	buildStream := output.Stream("local-build")
	err = m.Build(".", appName, buildStream, opt)
	if err != nil {
		log.Fatal(err)
	}

	tagStream := output.Stream("tag")
	err = TagForExport(m, appName, tagStream, buildID)
	if err != nil {
		log.Fatal(err)
	}

	testDocker, err := testM.BuildDockerfile(".", appName)
	if testDocker != nil {

	}
	testBuildSources, err := testM.BuildSources(".", appName)
	if testBuildSources != nil {

	}
	exportStream := output.Stream("export")
	imagesFile, err := Export(m, appName, exportStream, buildID)
	if err != nil {
		log.Fatal(err)
	}

	testJsonBytes, err := GenerateBuildJsonFile2(testM, appName, buildID, description)
	if err != nil {
		log.Fatal(err)
	}
	if testJsonBytes != nil {

	}

	buildJsonBytes, err := GenerateBuildJsonFile(m, appName, buildID, description)
	if err != nil {
		log.Fatal(err)
	}

	err = MakeFinalPackage(imagesFile, buildJsonBytes, appName, buildID)
	if err != nil {
		log.Fatal(err)
	}
}

func MakeFinalPackage(imagesFile *os.File, buildJsonBytes []byte, appName, buildID string) error {

	buf, err := ioutil.TempFile(".", "output_tar")
	if err != nil {
		return err
	}

	defer func() {
		buf.Close()
		os.Remove(buf.Name())

		imagesFile.Close()
		os.Remove(imagesFile.Name())
	}()

	tw := tar.NewWriter(buf)

	imagesStat, err := imagesFile.Stat()
	if err != nil {
		return err
	}

	hdr := &tar.Header{
		Name:     imagesFile.Name(),
		Mode:     int64(imagesStat.Mode()),
		Size:     imagesStat.Size(),
		Typeflag: tar.TypeReg,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	_, err = io.Copy(tw, imagesFile)
	if err != nil {
		return err
	}

	hdr = &tar.Header{
		Name:     "build.json",
		Mode:     int64(imagesStat.Mode()),
		Size:     int64(len(buildJsonBytes)),
		Typeflag: tar.TypeReg,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	if _, err := tw.Write(buildJsonBytes); err != nil {
		return err
	}

	tw.Close()
	buf.Sync()
	buf.Seek(0, 0)

	gzfile, err := os.Create(fmt.Sprintf("%s-%s.tgz", appName, buildID))
	if err != nil {
		return err
	}
	zw := gzip.NewWriter(gzfile)
	_, err = io.Copy(zw, buf)
	if err != nil {
		return err
	}

	zw.Flush()
	zw.Close()

	log.Println("Output: " + gzfile.Name())
	gzfile.Sync()
	gzfile.Close()

	return nil
}

func Export(m *manifest1.Manifest, appName string, s manifest1.Stream, buildID string) (*os.File, error) {
	fileName := fmt.Sprintf("%s.%s.tar", appName, buildID)
	params := []string{"save", "-o", fileName}

	for _, service := range m.Services {
		params = append(params, fmt.Sprintf("%s:%s.%s", appName, service.Name, buildID))
	}

	if err := manifest1.DefaultRunner.Run(s, manifest1.Docker(params...), manifest1.RunnerOptions{Verbose: true}); err != nil {
		return nil, fmt.Errorf("export error: %s", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func TagForExport(m *manifest1.Manifest, appName string, s manifest1.Stream, buildID string) error {
	for _, service := range m.Services {
		if err := manifest1.DefaultRunner.Run(s, manifest1.Docker("tag", fmt.Sprintf("%s/%s:latest", appName, service.Name), fmt.Sprintf("%s:%s.%s", appName, service.Name, buildID)), manifest1.RunnerOptions{Verbose: true}); err != nil {
			return fmt.Errorf("tag for export error: %s", err)
		}
	}

	return nil
}

func GenerateBuildJsonFile2(m *manifest.Manifest, appName, buildID, description string) ([]byte, error) {
	if m == nil {
		return nil, fmt.Errorf("m cannot be nil")
	}

	b := &structs.Build{
		Id:          buildID,
		App:         appName,
		Description: description,
		Release:     buildID,
	}

	manifestYaml, err := yaml.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	b.Manifest = string(manifestYaml)

	rez, err := json.Marshal(b)
	if err != nil {
		return []byte{}, err
	}

	return rez, nil
}

func GenerateBuildJsonFile(m *manifest1.Manifest, appName, buildID, description string) ([]byte, error) {
	if m == nil {
		return nil, fmt.Errorf("m cannot be nil")
	}

	b := &structs.Build{
		Id:          buildID,
		App:         appName,
		Description: description,
		Release:     buildID,
	}

	manifestYaml, err := yaml.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	b.Manifest = string(manifestYaml)

	rez, err := json.Marshal(b)
	if err != nil {
		return []byte{}, err
	}

	return rez, nil
}
