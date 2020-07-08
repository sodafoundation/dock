package eseries

const (
	defaultConfPath   = "/etc/opensds/driver/netapp_eseries_storage.yaml"
	DefaultAZ         = "default"
	StorageDriverName = "eseries-san"
	volumePrefix      = "opensds_"
	snapshotPrefix    = "opensds_snapshot_"
	driverContext     = "csi"
	VolumeVersion     = "1"
	SnapshotVersion   = "1"
	accessMode        = ""
	volumeMode        = "Block"
	bytesGiB          = 1073741824
)
const (
	DefaultHostType        = "linux_dm_mp"
	MinimumVolumeSizeBytes = 1048576 // 1 MiB

	// Constants for internal pool attributes
	Size   = "size"
	Region = "region"
	Zone   = "zone"
	Media  = "media"
)

// ConfigVersion is the expected version specified in the config file
const ConfigVersion = 1

// Default storage prefix
const DefaultDockerStoragePrefix = "netappdvp_"
const DefaultTridentStoragePrefix = "trident_"

// Default SAN igroup / host group names
const DefaultDockerIgroupName = "netappdvp"
const DefaultTridentIgroupName = "trident"

// Storage driver names specified in the config file, etc.
const (
	EseriesIscsiStorageDriverName = "eseries-iscsi"
	FakeStorageDriverName         = "fake"
)

// Filesystem types
const (
	FsXfs  = "xfs"
	FsExt3 = "ext3"
	FsExt4 = "ext4"
	FsRaw  = "raw"
)

// Default Filesystem value
const DefaultFileSystemType = FsExt4

const UnsetPool = ""
const DefaultVolumeSize = "1G"
