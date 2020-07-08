package eseries

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	tridentconfig "github.com/netapp/trident/config"
	"github.com/netapp/trident/storage"
	sa "github.com/netapp/trident/storage_attribute"
	drivers "github.com/netapp/trident/storage_drivers"
	"github.com/netapp/trident/storage_drivers/eseries"
	"github.com/netapp/trident/utils"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	odu "github.com/sodafoundation/dock/contrib/drivers/utils"
	. "github.com/sodafoundation/dock/contrib/drivers/utils/config"
	"github.com/sodafoundation/dock/pkg/model"
	pb "github.com/sodafoundation/dock/pkg/model/proto"
	"github.com/sodafoundation/dock/pkg/utils/config"
	"strconv"
	"strings"
)

func (d *SANDriverE) backendName() string {
	if d.sanStorageDriver.Config.BackendName == "" {
		// Use the old naming scheme if no name is specified
		return "eseries_" + d.sanStorageDriver.Config.HostDataIP
	} else {
		return d.sanStorageDriver.Config.BackendName
	}
}

// poolName constructs the name of the pool reported by this driver instance
func (d *SANDriverE) poolName(name string) string {

	return fmt.Sprintf("%s_%s",
		d.backendName(),
		strings.Replace(name, "-", "", -1))
}

func (d *SANDriverE) Protocol() string {
	return "iscsi"
}

func getVolumeName(id string) string {
	r := strings.NewReplacer("-", "")
	return volumePrefix + r.Replace(id)
}

func getSnapshotName(id string) string {
	r := strings.NewReplacer("-", "")
	return snapshotPrefix + r.Replace(id)
}

func (d *SANDriverE) GetVolumeConfig(name string, size int64) (volConfig *storage.VolumeConfig) {

	volConfig = &storage.VolumeConfig{
		Version:         tridentconfig.OrchestratorAPIVersion,
		Name:            name,
		InternalName:    name,
		Size:            strconv.FormatInt(size*bytesGiB, 10),
		Protocol:        tridentconfig.Block,
		SnapshotPolicy:  "",
		ExportPolicy:    "",
		SnapshotDir:     "false",
		UnixPermissions: "",
		StorageClass:    "",
		AccessMode:      tridentconfig.ReadWriteOnce,
		AccessInfo:      utils.VolumeAccessInfo{},
		BlockSize:       "",
		FileSystem:      "",
	}
	return volConfig
}
func (d *SANDriverE) GetSnapshotConfig(snapName string, volName string) (snapConfig *storage.SnapshotConfig) {
	snapConfig = &storage.SnapshotConfig{
		Version:            SnapshotVersion,
		Name:               snapName,
		InternalName:       snapName,
		VolumeName:         volName,
		VolumeInternalName: volName,
	}
	return snapConfig
}

func (d *SANDriverE) Setup() error {
	// Read NetApp Eseries config file
	d.conf = &EseiresConfig{}

	p := config.CONF.OsdsDock.Backends.NetappEseriesSan.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	if _, err := Parse(d.conf, p); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("unable to instantiate Eseries backend.")
		}
	}()

	empty := ""

	config := &drivers.ESeriesStorageDriverConfig{
		CommonStorageDriverConfig: &drivers.CommonStorageDriverConfig{
			Version:           d.conf.Version,
			StorageDriverName: StorageDriverName,
			StoragePrefixRaw:  json.RawMessage("{}"),
			StoragePrefix:     &empty,
		},
		WebProxyHostname:      d.conf.WebProxyHostname,
		WebProxyPort:          d.conf.WebProxyPort,
		WebProxyUseHTTP:       d.conf.WebProxyUseHTTP,
		WebProxyVerifyTLS:     d.conf.WebProxyVerifyTLS,
		Username:              d.conf.Username,
		Password:              d.conf.Password,
		ControllerA:           d.conf.ControllerA,
		ControllerB:           d.conf.ControllerB,
		PasswordArray:         d.conf.PasswordArray,
		PoolNameSearchPattern: d.conf.PoolNameSearchPattern,
		HostDataIP:            d.conf.HostDataIP,

		AccessGroup: d.conf.AccessGroup,
		HostType:    d.conf.HostType,
	}
	marshaledJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatal("unable to marshal Eseries config:  ", err)
	}
	configJSON := string(marshaledJSON)

	// Convert config (JSON or YAML) to JSON
	configJSONBytes, err := yaml.YAMLToJSON([]byte(configJSON))
	if err != nil {
		err = fmt.Errorf("invalid config format: %v", err)
		return err
	}
	configJSON = string(configJSONBytes)

	// Parse the common config struct from JSON
	commonConfig, err := drivers.ValidateCommonSettings(configJSON)
	if err != nil {
		err = fmt.Errorf("input failed validation: %v", err)
		return err
	}

	d.sanStorageDriver = &eseries.SANStorageDriver{
		Config: *config,
	}

	// Initialize the driver.
	if err = d.sanStorageDriver.Initialize(driverContext, configJSON, commonConfig); err != nil {
		log.Errorf("could not initialize storage driver (%s). failed: %v", commonConfig.StorageDriverName, err)
		return err
	}
	log.Infof("storage driver (%s) initialized successfully.", commonConfig.StorageDriverName)

	return nil
}
func (d *SANDriverE) Unset() error {
	//driver to clean up and stop any ongoing operations.
	d.sanStorageDriver.Terminate()
	return nil
}
func (d *SANDriverE) CreateVolume(opt *pb.CreateVolumeOpts) (vol *model.VolumeSpec, err error) {
	if opt.GetSnapshotId() != "" {
		return d.createVolumeFromSnapshot(opt)
	}

	var name = getVolumeName(opt.GetId())
	volConfig := d.GetVolumeConfig(name, opt.GetSize())

	storagePool := &storage.Pool{
		Name:               opt.GetPoolName(),
		StorageClasses:     make([]string, 0),
		Attributes:         make(map[string]sa.Offer),
		InternalAttributes: make(map[string]string),
	}
	err = d.sanStorageDriver.Create(volConfig, storagePool, make(map[string]sa.Request))
	if err != nil {
		log.Errorf("create volume (%s) failed: %v", opt.GetId(), err)
		return nil, err
	}
	log.Infof("volume (%s) created successfully.", opt.GetId())

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		Metadata:    opt.GetMetadata(),
	}, nil
}
func (d *SANDriverE) createVolumeFromSnapshot(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {

	var snapName = getSnapshotName(opt.GetSnapshotId())
	var volName = opt.GetMetadata()["volume"]
	var name = getVolumeName(opt.GetId())

	volConfig := d.GetVolumeConfig(name, opt.GetSize())
	volConfig.CloneSourceVolumeInternal = volName
	volConfig.CloneSourceSnapshot = volName
	volConfig.CloneSourceSnapshot = snapName
	err := d.sanStorageDriver.CreateClone(volConfig)
	if err != nil {
		log.Errorf("create volume (%s) from snapshot (%s) failed: %v", opt.GetId(), opt.GetSnapshotId(), err)
		return nil, err
	}

	log.Infof("volume (%s) created from snapshot (%s) successfully.", opt.GetId(), opt.GetSnapshotId())

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		Metadata:    opt.GetMetadata(),
	}, err
}
func (d *SANDriverE) PullVolume(volId string) (*model.VolumeSpec, error) {

	return nil, &model.NotImplementError{"method PullVolume has not been implemented yet"}
}
func (d *SANDriverE) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	var name = getVolumeName(opt.GetId())
	err := d.sanStorageDriver.Destroy(name)
	if err != nil {
		msg := fmt.Sprintf("delete volume (%s) failed: %v", opt.GetId(), err)
		log.Error(msg)
		return err
	}
	log.Infof("volume (%s) deleted successfully.", opt.GetId())
	return nil
}
func (d *SANDriverE) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	var name = getVolumeName(opt.GetId())
	volConfig := d.GetVolumeConfig(name, opt.GetSize())

	newSize := uint64(opt.GetSize() * bytesGiB)
	if err := d.sanStorageDriver.Resize(volConfig, newSize); err != nil {
		log.Errorf("extend volume (%s) failed, error: %v", name, err)
		return nil, err
	}
	log.Infof("volume (%s) extended successfully.", opt.GetId())
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		Metadata:    opt.GetMetadata(),
	}, nil
}
func (d *SANDriverE) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {

	var name = getVolumeName(opt.GetVolumeId())
	hostInfo := opt.GetHostInfo()
	initiator := odu.GetInitiatorName(hostInfo.GetInitiators(), opt.GetAccessProtocol())
	hostName := hostInfo.GetHost()

	publishInfo := &utils.VolumePublishInfo{
		HostIQN:  []string{initiator},
		HostIP:   []string{hostInfo.GetIp()},
		HostName: hostName,
	}

	err := d.sanStorageDriver.Publish(name, publishInfo)
	if err != nil {
		msg := fmt.Sprintf("volume (%s) attachment is failed: %v", opt.GetVolumeId(), err)
		log.Errorf(msg)
		return nil, err
	}

	log.Infof("volume (%s) attachment is created successfully", opt.GetVolumeId())

	connInfo := &model.ConnectionInfo{
		DriverVolumeType: opt.GetAccessProtocol(),
		ConnectionData: map[string]interface{}{
			"target_discovered": true,
			"volumeId":          opt.GetVolumeId(),
			"volume":            name,
			"description":       "NetApp Eseries Attachment",
			"hostName":          hostName,
			"initiator":         initiator,
			"targetIQN":         []string{publishInfo.IscsiTargetIQN},
			"targetPortal":      []string{hostInfo.GetIp() + ":3260"},
			"targetLun":         publishInfo.IscsiLunNumber,
			"igroup":            publishInfo.IscsiIgroup,
		},
	}

	log.Infof("initialize connection successfully: %v", connInfo)
	return connInfo, nil
}
func (d *SANDriverE) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error {
	var name = getVolumeName(opt.GetVolumeId())

	// Validate ExtantVolume exists before destroying
	volExists, err := d.sanStorageDriver.API.GetVolume(name)
	if err != nil {
		return fmt.Errorf("error checking for existing volume (%s), error: %v", name, err)
	}
	if volExists.Label == "" {
		log.Infof("volume %s already deleted, skipping destroy.", name)
		return nil
	}
	delete(opt.Metadata, opt.GetId())
	log.Infof("termination connection successfully")
	return nil
}
func (d *SANDriverE) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (snap *model.VolumeSnapshotSpec, err error) {
	var snapName = getSnapshotName(opt.GetId())
	var volName = getVolumeName(opt.GetVolumeId())

	snapConfig := d.GetSnapshotConfig(snapName, volName)

	snapshot, err := d.sanStorageDriver.CreateSnapshot(snapConfig)

	if err != nil {
		msg := fmt.Sprintf("create snapshot %s (%s) failed: %s", opt.GetName(), opt.GetId(), err)
		log.Error(msg)
		return nil, err
	}

	log.Infof("snapshot %s (%s) created successfully.", opt.GetName(), opt.GetId())

	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Size:        opt.GetSize(),
		Metadata: map[string]string{
			"name":         snapName,
			"volume":       volName,
			"creationTime": snapshot.Created,
			"size":         strconv.FormatInt(snapshot.SizeBytes/bytesGiB, 10) + "G",
		},
	}, nil
}
func (d *SANDriverE) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	// not used, do nothing
	return nil, &model.NotImplementError{"method PullSnapshot has not been implemented yet"}
}
func (d *SANDriverE) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {

	var snapName = getSnapshotName(opt.GetId())
	var volName = getVolumeName(opt.GetVolumeId())

	snapConfig := d.GetSnapshotConfig(snapName, volName)

	err := d.sanStorageDriver.DeleteSnapshot(snapConfig)

	if err != nil {
		msg := fmt.Sprintf("delete volume snapshot (%s) failed: %v", opt.GetId(), err)
		log.Error(msg)
		return err
	}
	log.Infof("volume snapshot (%s) deleted successfully", opt.GetId())
	return nil
}

func (d *SANDriverE) ListPools() ([]*model.StoragePoolSpec, error) {
	var pools []*model.StoragePoolSpec
	var opts map[string]string
	mediaType := utils.GetV(opts, "mediaType", "hdd")
	minFreeSpaceBytes := uint64(10)
	// Pool Name fetch from config file
	pln := d.poolName("sample-pool")
	LP, err := d.sanStorageDriver.API.GetVolumePools(mediaType, uint64(minFreeSpaceBytes), pln)
	if err != nil {
		msg := fmt.Sprintf("list pools failed: %v", err)
		log.Error(msg)
		return nil, err
	}
	c := d.conf
	for _, sp := range LP {
		if _, ok := c.Pool[sp.Label]; !ok {
			continue
		}
		poolFreeSpace, err := strconv.ParseUint(sp.FreeSpace, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Could not parse free space for pool: %v", err)
			log.Error(msg)
			return nil, err
		}
		totalcapcity := poolFreeSpace / minFreeSpaceBytes

		pool := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, sp.WorldWideName).String(),
			},
			Name:          sp.Label,
			TotalCapacity: int64(totalcapcity),
			FreeCapacity:  int64(poolFreeSpace),
			StorageType:   c.Pool[sp.DriveMediaType].StorageType,
		}
		if pool.AvailabilityZone == "" {
			pool.AvailabilityZone = DefaultAZ
		}
		pools = append(pools, pool)
	}

	log.Info("list pools successfully")
	return pools, nil
}
func (d *SANDriverE) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error) {

	return nil, &model.NotImplementError{S: "method InitializeSnapshotConnection has not been implemented yet."}
}

func (d *SANDriverE) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {

	return &model.NotImplementError{S: "method TerminateSnapshotConnection has not been implemented yet."}

}

func (d *SANDriverE) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
}

func (d *SANDriverE) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method UpdateVolumeGroup has not been implemented yet"}
}

func (d *SANDriverE) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts) error {
	return &model.NotImplementError{"method DeleteVolumeGroup has not been implemented yet"}
}
