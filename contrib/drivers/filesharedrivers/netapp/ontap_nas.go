// Copyright 2020 The SODA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package netapp

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ghodss/yaml"
	log "github.com/golang/glog"
	uuid "github.com/satori/go.uuid"

	"github.com/netapp/trident/storage"
	sa "github.com/netapp/trident/storage_attribute"
	drivers "github.com/netapp/trident/storage_drivers"
	"github.com/netapp/trident/storage_drivers/ontap"
	"github.com/netapp/trident/utils"

	. "github.com/sodafoundation/dock/contrib/drivers/utils/config"
	"github.com/sodafoundation/dock/pkg/model"
	pb "github.com/sodafoundation/dock/pkg/model/proto"
	"github.com/sodafoundation/dock/pkg/utils/config"
)

func (d *NASDriver) GetVolumeConfig(name string, size int64) (volConfig *storage.VolumeConfig) {
	volConfig = &storage.VolumeConfig{
		Version:      VolumeVersion,
		Name:         name,
		InternalName: name,
		Size:         strconv.FormatInt(size*bytesGiB, 10),
		Protocol:     d.nasStorageDriver.GetProtocol(),
		AccessMode:   accessMode,
		AccessInfo:   utils.VolumeAccessInfo{},
	}
	return volConfig
}

// This function reads the config file for the driver in /etc/opensds/driver
// and creates a connection with the storage
func (d *NASDriver) Setup() error {
	// Read NetApp ONTAP config file
	d.conf = &ONTAPConfig{}

	p := config.CONF.OsdsDock.Backends.NetappOntapNas.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	if _, err := Parse(d.conf, p); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("unable to instantiate ontap backend.")
		}
	}()

	empty := ""
	config := &drivers.OntapStorageDriverConfig{
		CommonStorageDriverConfig: &drivers.CommonStorageDriverConfig{
			Version:           d.conf.Version,
			StorageDriverName: StorageDriverName,
			StoragePrefixRaw:  json.RawMessage("{}"),
			StoragePrefix:     &empty,
		},
		ManagementLIF: d.conf.ManagementLIF, // This is Storage VM (SVM) IP
		DataLIF:       d.conf.DataLIF,       // This is the data LIF for NFS/CIFS access
		SVM:           d.conf.Svm,
		Username:      d.conf.Username,
		Password:      d.conf.Password,
	}
	marshaledJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatal("unable to marshal ONTAP config:  ", err)
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

	// Create the NAS Storage driver instance from the trident driver
	d.nasStorageDriver = &ontap.NASStorageDriver{
		Config: *config,
	}

	// Initialize the driver.
	if err = d.nasStorageDriver.Initialize(driverContext, configJSON, commonConfig); err != nil {
		log.Errorf("could not initialize storage driver (%s). failed: %v", commonConfig.StorageDriverName, err)
		return err
	}
	log.Infof("Storage driver (%s) initialized successfully.", commonConfig.StorageDriverName)

	return nil
}

func (d *NASDriver) Unset() error {
	//driver to clean up and stop any ongoing operations.
	d.nasStorageDriver.Terminate()
	return nil
}

// Method to create a fileshare on the filer
func (d *NASDriver) CreateFileShare(opt *pb.CreateFileShareOpts) (vol *model.FileShareSpec, err error) {
	var name = opt.GetName()
	volConfig := d.GetVolumeConfig(name, opt.GetSize())
	storagePool := &storage.Pool{
		Name:               opt.GetPoolName(),
		StorageClasses:     make([]string, 0),
		Attributes:         make(map[string]sa.Offer),
		InternalAttributes: make(map[string]string),
	}

	// Here the storage pool is the aggregate on the Netapp Filer
	// This Create call will create a volume and mount to be used as
	// the export/share
	err = d.nasStorageDriver.Create(volConfig, storagePool, make(map[string]sa.Request))
	if err != nil {
		log.Errorf("create nas fileshare (%s) failed: %v", opt.GetId(), err)
		return nil, err
	}
	var exportLocation []string
	server := d.conf.DataLIF
	location := server + ":/" + name
	exportLocation = append(exportLocation, location)

	return &model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		Protocols:        []string{NFSProtocol},
		AvailabilityZone: opt.GetAvailabilityZone(),
		PoolId:           opt.GetPoolId(),
		ExportLocations:  exportLocation,
		Metadata: map[string]string{
			KFileshareName: name,
		},
	}, nil
}

// This method lists all the associated aggregates on the
// provided SVM of the filer (SVM is provided as part of config)
func (d *NASDriver) ListPools() ([]*model.StoragePoolSpec, error) {
	log.Info("Entered to list pools for netapp nas")

	var pools []*model.StoragePoolSpec
	aggregates, err := d.nasStorageDriver.API.VserverGetAggregateNames()

	if err != nil {
		msg := fmt.Sprintf("list pools failed: %v", err)
		log.Error(msg)
		return nil, err
	}

	c := d.conf
	for _, aggr := range aggregates {
		if _, ok := c.Pool[aggr]; !ok {
			continue
		}
		aggregate, _ := d.nasStorageDriver.API.AggregateCommitment(aggr)
		aggregateCapacity := aggregate.AggregateSize / bytesGiB
		aggregateAllocatedCapacity := aggregate.TotalAllocated / bytesGiB

		pool := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, aggr).String(),
			},
			Name:             aggr,
			TotalCapacity:    int64(aggregateCapacity),
			FreeCapacity:     int64(aggregateCapacity) - int64(aggregateAllocatedCapacity),
			ConsumedCapacity: int64(aggregateAllocatedCapacity),
			StorageType:      c.Pool[aggr].StorageType,
			Extras:           c.Pool[aggr].Extras,
			AvailabilityZone: c.Pool[aggr].AvailabilityZone,
		}
		if pool.AvailabilityZone == "" {
			pool.AvailabilityZone = DefaultAZ
		}
		pools = append(pools, pool)
	}

	log.Info("List pools successfully:", pools)
	return pools, nil
}

// Function to delete the Fileshare
// This function basically calls to delete the corresponding volume
// on the filer and unmount it
func (d *NASDriver) DeleteFileShare(opts *pb.DeleteFileShareOpts) error {
	var name = opts.GetMetadata()[KFileshareName]

	err := d.nasStorageDriver.Destroy(name)
	if err != nil {
		log.Errorf("delete nas fileshare (%s) failed: %v", name, err)
		return err
	}

	log.Info("Deleted fileshare:", name)
	return nil
}

// Function to create a fileshare snapshot for th specified fileshare (volume in backend)
func (d *NASDriver) CreateFileShareSnapshot(opts *pb.CreateFileShareSnapshotOpts) (*model.FileShareSnapshotSpec, error) {
	snapName := opts.GetName()
	fileshareName := opts.GetMetadata()[KFileshareName]
	log.Infof("Creating snapshot for fileshare %s", fileshareName)

	snapConfig := &storage.SnapshotConfig{
		Version:            SnapshotVersion,
		Name:               snapName,
		InternalName:       snapName,
		VolumeName:         fileshareName,
		VolumeInternalName: fileshareName,
	}
	snapshot, err := d.nasStorageDriver.CreateSnapshot(snapConfig)

	if err != nil {
		log.Errorf("create snapshot for %s with snapshot name %s failed with err %s",
			fileshareName, snapName, err)
		return nil, err
	}

	log.Infof("Snapshot %s created successfully for fileshare %s", snapName, fileshareName)
	return &model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opts.GetId(),
		},
		Name:         opts.GetName(),
		SnapshotSize: snapshot.SizeBytes / bytesGiB,
		Description:  opts.GetDescription(),
		Metadata: map[string]string{
			KFileshareSnapName: snapName,
			KFileshareSnapID:   opts.GetId(),
			KFileshareName:     fileshareName,
		},
	}, nil
}

// Function to delete a fileshare snapshot for th specified fileshare (volume in backend)
func (d *NASDriver) DeleteFileShareSnapshot(opts *pb.DeleteFileShareSnapshotOpts) error {
	snapName := opts.GetMetadata()[KFileshareSnapName]
	fileshareName := opts.GetMetadata()[KFileshareName]

	snapConfig := &storage.SnapshotConfig{
		Version:            SnapshotVersion,
		Name:               snapName,
		InternalName:       snapName,
		VolumeName:         fileshareName,
		VolumeInternalName: fileshareName,
	}

	err := d.nasStorageDriver.DeleteSnapshot(snapConfig)
	if err != nil {
		log.Errorf("delete snapshot for %s with snapshot name %s failed with err %s",
			fileshareName, snapName, err)
		return err
	}
	log.Infof("Successfully deleted snapshot %s of fileshare %s", snapName, fileshareName)
	return nil
}

func (d *NASDriver) CreateFileShareAcl(opt *pb.CreateFileShareAclOpts) (*model.FileShareAclSpec, error) {
	return nil, &model.NotImplementError{"method is not implemented"}
}

func (d *NASDriver) DeleteFileShareAcl(opt *pb.DeleteFileShareAclOpts) error {
	return &model.NotImplementError{"method is not implemented"}
}
