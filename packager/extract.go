package packager

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/lunchbox/utils"
)

// Extract extracts the app content if argument is an archive, or does nothing if a dir.
// It returns effective app name, and cleanup function
func Extract(appname string) (string, func(), error) {
	// try verbatim first
	s, err := os.Stat(appname)
	if err != nil {
		// try appending our extension
		appname = utils.DirNameFromAppName(appname)
		s, err = os.Stat(appname)
	}
	if err != nil {
		return "", func() {}, err
	}
	if s.IsDir() {
		// directory: already decompressed
		return appname, func() {}, nil
	}
	// not a dir: probably a tarball package, extract that in a temp dir
	tempDir, err := ioutil.TempDir("", "dockerapp")
	if err != nil {
		return "", func() {}, err
	}
	err = extract(appname, tempDir)
	if err != nil {
		return "", func() {}, err
	}
	return tempDir, func() { os.RemoveAll(tempDir) }, nil
}

func extract(appname, outputDir string) error {
	f, err := os.Open(appname)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(f)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		outputDir = outputDir + "/"
		switch header.Typeflag {
		case tar.TypeDir: // = directory
			os.Mkdir(outputDir+header.Name, 0755)
		case tar.TypeReg: // = regular file
			data := make([]byte, header.Size)
			_, err := tarReader.Read(data)
			if err != nil {
				return err
			}
			ioutil.WriteFile(outputDir+header.Name, data, 0644)
		}
	}
	return nil
}