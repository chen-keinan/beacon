package utils

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

//PluginSourceSubFolder plugin source folder
const PluginSourceSubFolder = "plugins/source"

//CompilePluginSubFolder plugins complied folder
const CompilePluginSubFolder = "plugins/compile"

//FolderMgr defines the interface for kube-knark folder
//fileutil.go
//go:generate mockgen -destination=./mocks/mock_FolderMgr.go -package=mocks . FolderMgr
type FolderMgr interface {
	CreateFolder(folderName string) error
	GetHomeFolder() (string, error)
}

//KFolder kube-knark folder object
type KFolder struct {
}

//NewKFolder return KFolder instance
func NewKFolder() FolderMgr {
	return &KFolder{}
}

//CreateFolder create new kube beacon folder
func (kf KFolder) CreateFolder(folderName string) error {
	_, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderName, 0750)
		if errDir != nil {
			return err
		}
	}
	return nil
}

//GetHomeFolder return kube-knark home folder
func (kf KFolder) GetHomeFolder() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	// User can set a custom KUBE_KNARK_HOME from environment variable
	usrHome := GetEnv(common.BeaconHomeEnvVar, usr.HomeDir)
	return path.Join(usrHome, ".beacon"), nil
}

//GetPluginSourceSubFolder return plugins source folder path
func GetPluginSourceSubFolder(fm FolderMgr) (string, error) {
	folder, err := fm.GetHomeFolder()
	if err != nil {
		return "", err
	}
	return path.Join(folder, PluginSourceSubFolder), nil
}

//GetCompilePluginSubFolder return plugin compiled folder path
func GetCompilePluginSubFolder(fm FolderMgr) (string, error) {
	folder, err := fm.GetHomeFolder()
	if err != nil {
		return "", err
	}
	return path.Join(folder, CompilePluginSubFolder), nil
}

//CreatePluginsCompiledFolderIfNotExist create plugins compiled folder if not exist
func CreatePluginsCompiledFolderIfNotExist(fm FolderMgr) error {
	ebpfFolder, err := GetCompilePluginSubFolder(fm)
	if err != nil {
		return err
	}
	return fm.CreateFolder(ebpfFolder)
}

//CreatePluginsSourceFolderIfNotExist plugins source folder if not exist
func CreatePluginsSourceFolderIfNotExist(fm FolderMgr) error {
	pluginfFolder, err := GetPluginSourceSubFolder(fm)
	if err != nil {
		return err
	}
	return fm.CreateFolder(pluginfFolder)
}

//GetHomeFolder return beacon home folder
func GetHomeFolder() string {
	usr, err := user.Current()
	if err != nil {
		panic("Failed to fetch user home folder")
	}
	// User can set a custom BEACON_HOME from environment variable
	usrHome := GetEnv(common.BeaconHomeEnvVar, usr.HomeDir)
	return path.Join(usrHome, ".beacon")
}

//CreateHomeFolderIfNotExist create beacon home folder if not exist
func CreateHomeFolderIfNotExist(fm FolderMgr) error {
	beaconFolder, err := fm.GetHomeFolder()
	if err != nil {
		return err
	}
	_, err = os.Stat(beaconFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(beaconFolder, 0750)
		if errDir != nil {
			return fmt.Errorf("failed to create beacon home folder at %s", beaconFolder)
		}
	}
	return nil
}

//GetBenchmarkFolder return benchmark folder
func GetBenchmarkFolder(spec, version string) string {
	return filepath.Join(GetHomeFolder(), fmt.Sprintf("benchmarks/%s/%s/", spec, version))
}

//CreateBenchmarkFolderIfNotExist create beacon benchmark folder if not exist
func CreateBenchmarkFolderIfNotExist(spec, version string) error {
	benchmarkFolder := GetBenchmarkFolder(spec, version)
	_, err := os.Stat(benchmarkFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(benchmarkFolder, 0750)
		if errDir != nil {
			return fmt.Errorf("failed to create beacon benchmark folder folder at %s", benchmarkFolder)
		}
	}
	return nil
}

//GetK8sBenchAuditFiles return k8s benchmark file
func GetK8sBenchAuditFiles(spec, version string) ([]FilesInfo, error) {
	filesData := make([]FilesInfo, 0)
	folder := GetBenchmarkFolder(spec, version)
	filesInfo, err := ioutil.ReadDir(filepath.Join(folder))
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range filesInfo {
		filePath := filepath.Join(GetBenchmarkFolder(spec, version), filepath.Clean(fileInfo.Name()))
		fData, err := ioutil.ReadFile(filepath.Clean(filePath))
		if err != nil {
			return nil, err
		}
		filesData = append(filesData, FilesInfo{fileInfo.Name(), string(fData)})
	}
	return filesData, nil
}

//FilesInfo file data
type FilesInfo struct {
	Name string
	Data string
}

//GetEnv Get Environment Variable value or return default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
