package agentmgr

import (
	"io/ioutil"
	"os"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/partition/mbr"
)

// AgentManager is an abstraction to update the AgentEnv into a VM
type AgentManager interface {
	Update(apiv1.AgentEnv) error
	ToBytes() ([]byte, error)
}

// NewAgentManager will initalize a new config drive for AgentEnv settings
func NewAgentManager() (AgentManager, error) {
	// BUG: Not handling error
	name, err := tempFileName()
	if err != nil {
		return nil, bosherr.WrapError(err, "unable to generate agent config disk temp file")
	}
	mgr := configDriveManager{
		diskFileName: name,
	}
	err = mgr.createDisk()
	if err != nil {
		return nil, bosherr.WrapError(err, "unable to create config disk")
	}
	return mgr, nil
}

// NewAgentManagerFromData will allow AgentEnv settings updates on an existing config drive
func NewAgentManagerFromData(data []byte) (AgentManager, error) {
	// BUG: Not handling error
	name, err := tempFileName()
	if err != nil {
		return nil, bosherr.WrapError(err, "unable to generate agent config disk temp file")
	}
	err = ioutil.WriteFile(name, data, 0666)
	if err != nil {
		return nil, bosherr.WrapError(err, "unable to store config disk to temp file")
	}
	mgr := configDriveManager{
		diskFileName: name,
	}
	return mgr, nil
}

type configDriveManager struct {
	diskFileName string
}

const configPath = "/ec2/latest"
const metaDataPath = configPath + "/meta-data.json"
const userDataPath = configPath + "/user-data"

func tempFileName() (string, error) {
	f, err := ioutil.TempFile("", "config-disk")
	if err != nil {
		return "", err
	}
	name := f.Name()
	err = f.Close()
	if err != nil {
		return "", err
	}
	err = os.Remove(name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (c configDriveManager) Update(agentEnv apiv1.AgentEnv) error {
	disk, err := diskfs.Open(c.diskFileName)
	if err != nil {
		return err
	}

	// The partition table doesn't appear to get populated? Manually populating it.
	table, err := disk.GetPartitionTable()
	if err != nil {
		return err
	}
	disk.Table = table

	fs, err := disk.GetFilesystem(1)
	if err != nil {
		return err
	}

	// The AgentEnv appears to be what goes into userdata
	userDataContent, err := agentEnv.AsBytes()
	if err != nil {
		return err
	}

	err = c.writeFile(fs, metaDataPath, []byte("{}"))
	if err != nil {
		return err
	}

	err = c.writeFile(fs, userDataPath, userDataContent)
	if err != nil {
		return err
	}

	return nil
}

func (c configDriveManager) writeFile(fs filesystem.FileSystem, path string, contents []byte) error {
	rw, err := fs.OpenFile(path, os.O_CREATE|os.O_RDWR)
	if err != nil {
		return err
	}

	_, err = rw.Write(contents)
	if err != nil {
		return err
	}

	return nil
}

func (c configDriveManager) ToBytes() ([]byte, error) {
	f, err := os.Open(c.diskFileName)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func (c configDriveManager) createDisk() error {
	// Note that the sizes are rough guesstimates.
	// diskSize = 35MB; minimim size is 32MB but...
	// partition start of 2048 is ~1MB into disk.
	// partition size of 68000 is about 33.25MB.

	diskSize := uint64(35 * 1024 * 1024)
	image, err := diskfs.Create(c.diskFileName, int64(diskSize), diskfs.Raw)
	if err != nil {
		return err
	}

	table := &mbr.Table{
		LogicalSectorSize:  512,
		PhysicalSectorSize: 512,
		Partitions: []*mbr.Partition{
			{
				Bootable: false,
				Type:     mbr.Fat32LBA,
				Start:    2048,
				Size:     68000,
			},
		},
	}
	err = image.Partition(table)

	fs, err := image.CreateFilesystem(disk.FilesystemSpec{
		Partition:   1,
		FSType:      filesystem.TypeFat32,
		VolumeLabel: "config-2",
	})
	if err != nil {
		return err
	}

	return fs.Mkdir(configPath)
}
