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
// License for the specific l

package netapp

//  default value for driver
const (
	defaultConfPath    = "/etc/opensds/driver/netapp_ontap_nas.yaml"
	DefaultAZ          = "default"
	fileSharePrefix    = "SODA_"
	snapshotPrefix     = "SODA_snapshot_"
	KFileshareName     = "nfsFileshareName"
	KFileshareID       = "nfsFileshareID"
	KFileshareSnapName = "snapshotName"
	KFileshareSnapID   = "snapshotID"
	KLvPath            = "lunPath"
	KLvIdFormat        = "NAA"
	StorageDriverName  = "ontap-nas"
	driverContext      = "csi"
	VolumeVersion      = "1"
	SnapshotVersion    = "1"
	accessMode         = ""
	bytesGiB           = 1073741824
)
