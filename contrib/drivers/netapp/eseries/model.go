package eseries

import (
	"github.com/netapp/trident/storage_drivers/eseries"
	. "github.com/sodafoundation/dock/contrib/drivers/utils/config"
)

//type BackendOptions struct {
//	Version           int    `yaml:"version"`
//	StorageDriverName string `yaml:"storageDriverName"`
//	ManagementLIF     string `yaml:"managementLIF"`
//	DataLIF           string `yaml:"dataLIF"`
//	Svm               string `yaml:"svm"`
//	IgroupName        string `yaml:"igroupName"`
//	Username          string `yaml:"username"`
//	Password          string `yaml:"password"`
//}
//
//type ONTAPConfig struct {
//	BackendOptions `yaml:"backendOptions"`
//	Pool           map[string]PoolProperties `yaml:"pool,flow"`
//}
//
//type SANDriver struct {
//	sanStorageDriver *ontap.SANStorageDriver
//	conf             *ONTAPConfig
//}
//
//type Pool struct {
//	PoolId        int   `json:"poolId"`
//	TotalCapacity int64 `json:"totalCapacity"`
//	AllocCapacity int64 `json:"allocatedCapacity"`
//	UsedCapacity  int64 `json:"usedCapacity"`
//}

type BackendOptions struct {
	Version               int    `yaml:"version"`
	DriverName            string `yaml:"DriverName"`
	DebugTraceFlags       int    `yaml:"debugTraceFlags"`
	DisableDelete         string `yaml:"disableDelete"`
	StoragePrefix         string `yaml:"disableDelete"`
	AccessGroup           string `yaml:"accessGroup"`
	HostType              string `yaml:"hostType"`
	PoolNameSearchPattern string `yaml:"poolNameSearchPattern"`
	ControllerA           string `yaml:"controllerA"`
	ControllerB           string `yaml:"controllerB"`
	WebProxyHostname      string `yaml:"webProxyHostname"`
	WebProxyVerifyTLS     bool   `yaml:"webProxyVerifyTLS"`
	WebProxyPort          string `yaml:"webProxyPort"`
	WebProxyUseHTTP       bool   `yaml:"webProxyUseHTTP"`
	PasswordArray         string `yaml:"passwordArray"`
	Size                  int    `yaml:"size"`
	Username              string `yaml:"username"`
	Password              string `yaml:"password"`
	Protocol              string `yaml:"protocol"`
	ConfigVersion         string `yaml:"configVersion"`
	HostDataIP            string `yaml:"hostDataIP"`
	Telemetry             string
}

type EseiresConfig struct {
	BackendOptions `yaml:"backendOptions"`
	Pool           map[string]PoolProperties `yaml:"pool,flow"`
}
type SANDriverE struct {
	sanStorageDriver *eseries.SANStorageDriver
	conf             *EseiresConfig
}
type Pool struct {
	PoolId        int   `json:"poolId"`
	TotalCapacity int64 `json:"totalCapacity"`
	AllocCapacity int64 `json:"allocatedCapacity"`
	UsedCapacity  int64 `json:"usedCapacity"`
}
