// Copyright 2017 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
This module implements the etcd database operation of data structure
defined in api module.

*/

package etcd

import (
	"reflect"
	"strings"
	"testing"

	c "github.com/sodafoundation/dock/pkg/context"
	"github.com/sodafoundation/dock/pkg/model"
	. "github.com/sodafoundation/dock/testutils/collection"
)

type fakeClientCaller struct{}

func (*fakeClientCaller) Create(req *Request) *Response {
	return &Response{
		Status: "Success",
	}
}

func (*fakeClientCaller) Get(req *Request) *Response {
	var resp []string

	if strings.Contains(req.Url, "docks") {
		resp = append(resp, StringSliceDocks[0])
	}
	if strings.Contains(req.Url, "pools") {
		resp = append(resp, StringSlicePools[0])
	}
	if strings.Contains(req.Url, "volumes") {
		resp = append(resp, StringSliceVolumes[0])
	}
	if strings.Contains(req.Url, "attachments") {
		resp = append(resp, StringSliceAttachments[0])
	}
	if strings.Contains(req.Url, "snapshots") {
		resp = append(resp, StringSliceSnapshots[0])
	}
	if strings.Contains(req.Url, "replications") {
		resp = append(resp, StringSliceReplications[0])
	}
	if strings.Contains(req.Url, "acls") {
		resp = append(resp, ByteFileShareAcl)
	}
	return &Response{
		Status:  "Success",
		Message: resp,
	}
}

func (*fakeClientCaller) List(req *Request) *Response {
	var resp []string

	if strings.Contains(req.Url, "docks") {
		resp = StringSliceDocks
	}
	if strings.Contains(req.Url, "pools") {
		resp = StringSlicePools
	}
	if strings.Contains(req.Url, "volumes") {
		resp = StringSliceVolumes
	}
	if strings.Contains(req.Url, "attachments") {
		resp = StringSliceAttachments
	}
	if strings.Contains(req.Url, "snapshots") {
		resp = StringSliceSnapshots
	}
	if strings.Contains(req.Url, "replications") {
		resp = StringSliceReplications
	}
	return &Response{
		Status:  "Success",
		Message: resp,
	}
}

func (*fakeClientCaller) Update(req *Request) *Response {
	return &Response{
		Status: "Success",
	}
}

func (*fakeClientCaller) Delete(req *Request) *Response {
	return &Response{
		Status: "Success",
	}
}

var fc = &Client{
	clientInterface: &fakeClientCaller{},
}

func TestCreateDock(t *testing.T) {
	if _, err := fc.CreateDock(c.NewAdminContext(), &model.DockSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create dock failed:", err)
	}
}

func TestCreatePool(t *testing.T) {
	if _, err := fc.CreatePool(c.NewAdminContext(), &model.StoragePoolSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create pool failed:", err)
	}
}

func TestGetDock(t *testing.T) {
	dck, err := fc.GetDock(c.NewAdminContext(), "")
	if err != nil {
		t.Error("Get dock failed:", err)
	}

	var expected = &SampleDocks[0]
	if !reflect.DeepEqual(dck, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, dck)
	}
}

func TestGetPool(t *testing.T) {
	pol, err := fc.GetPool(c.NewAdminContext(), "")
	if err != nil {
		t.Error("Get pool failed:", err)
	}

	var expected = &SamplePools[0]
	if !reflect.DeepEqual(pol, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, pol)
	}
}

/*func TestGetVolume(t *testing.T) {
	vol, err := fc.GetVolume(c.NewAdminContext(), "")
	if err != nil {
		t.Error("Get volume failed:", err)
	}

	var expected = &SampleVolumes[0]
	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, vol)
	}
}*/

func TestListDocks(t *testing.T) {
	m := map[string][]string{
		"offset":     {"0"},
		"limit":      {"732"},
		"sortDir":    {"desc"},
		"sortKey":    {"id"},
		"Name":       {"sample"},
		"DriverName": {"sample"},
	}

	dcks, err := fc.ListDocksWithFilter(c.NewAdminContext(), m)
	if err != nil {
		t.Error("List docks failed:", err)
	}

	var expected []*model.DockSpec
	expected = append(expected, &SampleDocks[0])
	if !reflect.DeepEqual(dcks, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, dcks)
	}
}

func TestListPools(t *testing.T) {
	m := map[string][]string{
		"offset":  {"0"},
		"limit":   {"-5"},
		"sortDir": {"desc"},
		"sortKey": {"DockId"},
		"Name":    {"sample-pool-01"},
	}
	pols, err := fc.ListPoolsWithFilter(c.NewAdminContext(), m)
	if err != nil {
		t.Error("List pools failed:", err)
	}

	var expected []*model.StoragePoolSpec
	expected = append(expected, &SamplePools[0])
	if !reflect.DeepEqual(pols, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, pols)
	}
}

/*func TestUpdateVolume(t *testing.T) {
	var vol = model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		Name:        "Test Name",
		Description: "Test Description",
	}

	result, err := fc.UpdateVolume(c.NewAdminContext(), &vol)
	if err != nil {
		t.Error("Update volumes failed:", err)
	}

	if result.Id != "bd5b12a8-a101-11e7-941e-d77981b584d8" {
		t.Errorf("Expected %+v, got %+v\n", "bd5b12a8-a101-11e7-941e-d77981b584d8", result.Id)
	}

	if result.Name != "Test Name" {
		t.Errorf("Expected %+v, got %+v\n", "Test Name", result.Name)
	}

	if result.Description != "Test Description" {
		t.Errorf("Expected %+v, got %+v\n", "Test Description", result.Description)
	}

	if result.PoolId != "084bf71e-a102-11e7-88a8-e31fe6d52248" {
		t.Errorf("Expected %+v, got %+v\n", "084bf71e-a102-11e7-88a8-e31fe6d52248", result.PoolId)
	}
}*/

func TestDeleteDock(t *testing.T) {
	if err := fc.DeleteDock(c.NewAdminContext(), ""); err != nil {
		t.Error("Delete dock failed:", err)
	}
}

func TestDeletePool(t *testing.T) {
	if err := fc.DeletePool(c.NewAdminContext(), ""); err != nil {
		t.Error("Delete pool failed:", err)
	}
}

/*func TestDeleteVolume(t *testing.T) {
	if err := fc.DeleteVolume(c.NewAdminContext(), ""); err != nil {
		t.Error("Delete volume failed:", err)
	}
}*/

func TestFilterAndSort(t *testing.T) {
	// Use a specific type(pool) to test unit test
	type test struct {
		input    []*model.StoragePoolSpec
		param    map[string][]string
		expected []*model.StoragePoolSpec
	}
	tests := []test{
		// select by storage type
		{
			input: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
				&SamplePools[2],
			},
			param: map[string][]string{
				"storageType": {"block"},
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
			},
		},
		// sort by name asc
		{
			input: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
				&SamplePools[2],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"asc"},
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[2],
				&SamplePools[0],
				&SamplePools[1],
			},
		},
		// sort by name desc
		{
			input: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
				&SamplePools[2],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"desc"},
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[1],
				&SamplePools[0],
				&SamplePools[2],
			},
		},
		// limit is 2
		{
			input: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
				&SamplePools[2],
			},
			param: map[string][]string{
				"limit":  {"2"},
				"offset": {"1"},
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[1],
				&SamplePools[2],
			},
		},
	}
	for _, testcase := range tests {
		ret := fc.FilterAndSort(testcase.input, testcase.param, sortableKeysMap[typePools])
		var res = []*model.StoragePoolSpec{}
		for _, data := range ret.([]interface{}) {
			res = append(res, data.(*model.StoragePoolSpec))
		}
		if !reflect.DeepEqual(res, testcase.expected) {
			var expected []model.StoragePoolSpec
			for _, value := range testcase.expected {
				expected = append(expected, *value)
			}
			var got []model.StoragePoolSpec
			for _, value := range res {
				got = append(got, *value)
			}
			t.Errorf("Expected %+v\n", expected)
			t.Errorf("Got %+v\n", got)
		}
	}
}

func TestListFileSharesWithFilter(t *testing.T) {
	type test struct {
		input    []*model.FileShareSpec
		param    map[string][]string
		expected []*model.FileShareSpec
	}
	tests := []test{
		// select by storage type
		{
			input: []*model.FileShareSpec{
				&SampleFileShares[0],
				&SampleFileShares[1],
			},
			param: map[string][]string{
				"poolId": {"a5965ebe-dg2c-434t-b28e-f373746a71ca"},
			},
			expected: []*model.FileShareSpec{
				&SampleFileShares[0],
			},
		},
		// sort by name asc
		{
			input: []*model.FileShareSpec{
				&SampleFileShares[0],
				&SampleFileShares[1],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"asc"},
			},
			expected: []*model.FileShareSpec{
				&SampleFileShares[0],
				&SampleFileShares[1],
			},
		},
		// sort by name desc
		{
			input: []*model.FileShareSpec{
				&SampleFileShares[0],
				&SampleFileShares[1],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"desc"},
			},
			expected: []*model.FileShareSpec{
				&SampleFileShares[1],
				&SampleFileShares[0],
			},
		},
		// limit is 1
		{
			input: []*model.FileShareSpec{
				&SampleFileShares[0],
				&SampleFileShares[1]},
			param: map[string][]string{
				"limit":  {"1"},
				"offset": {"1"},
			},
			expected: []*model.FileShareSpec{
				&SampleFileShares[1],
			},
		},
	}
	for _, testcase := range tests {
		ret := fc.FilterAndSort(testcase.input, testcase.param, sortableKeysMap[typeFileShares])
		var res = []*model.FileShareSpec{}
		for _, data := range ret.([]interface{}) {
			res = append(res, data.(*model.FileShareSpec))
		}
		if !reflect.DeepEqual(res, testcase.expected) {
			var expected []model.FileShareSpec
			for _, value := range testcase.expected {
				expected = append(expected, *value)
			}
			var got []model.FileShareSpec
			for _, value := range res {
				got = append(got, *value)
			}
			t.Errorf("Expected %+v\n", expected)
			t.Errorf("Got %+v\n", got)
		}
	}
}
func TestListFileShareSnapshotsWithFilter(t *testing.T) {
	type test struct {
		input    []*model.FileShareSnapshotSpec
		param    map[string][]string
		expected []*model.FileShareSnapshotSpec
	}
	tests := []test{
		// select by storage type
		{
			input: []*model.FileShareSnapshotSpec{
				&SampleFileShareSnapshots[0],
				&SampleFileShareSnapshots[1],
			},
			param: map[string][]string{
				"status": {"error"},
			},
			expected: []*model.FileShareSnapshotSpec{},
		},
		// sort by name asc
		{
			input: []*model.FileShareSnapshotSpec{
				&SampleFileShareSnapshots[0],
				&SampleFileShareSnapshots[1],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"asc"},
			},
			expected: []*model.FileShareSnapshotSpec{
				&SampleFileShareSnapshots[0],
				&SampleFileShareSnapshots[1],
			},
		},
		// sort by name desc
		{
			input: []*model.FileShareSnapshotSpec{
				&SampleFileShareSnapshots[0],
				&SampleFileShareSnapshots[1],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"desc"},
			},
			expected: []*model.FileShareSnapshotSpec{
				&SampleFileShareSnapshots[1],
				&SampleFileShareSnapshots[0],
			},
		},
		// limit is 1
		{
			input: []*model.FileShareSnapshotSpec{
				&SampleFileShareSnapshots[0],
				&SampleFileShareSnapshots[1]},
			param: map[string][]string{
				"limit":  {"1"},
				"offset": {"1"},
			},
			expected: []*model.FileShareSnapshotSpec{
				&SampleFileShareSnapshots[1],
			},
		},
	}
	for _, testcase := range tests {
		ret := fc.FilterAndSort(testcase.input, testcase.param, sortableKeysMap[typeFileShareSnapshots])
		var res = []*model.FileShareSnapshotSpec{}
		for _, data := range ret.([]interface{}) {
			res = append(res, data.(*model.FileShareSnapshotSpec))
		}
		if !reflect.DeepEqual(res, testcase.expected) {
			var expected []model.FileShareSnapshotSpec
			for _, value := range testcase.expected {
				expected = append(expected, *value)
			}
			var got []model.FileShareSnapshotSpec
			for _, value := range res {
				got = append(got, *value)
			}
			t.Errorf("Expected %+v\n", expected)
			t.Errorf("Got %+v\n", got)
		}
	}
}
func TestListDocksWithFilter(t *testing.T) {
	type test struct {
		input    []*model.DockSpec
		param    map[string][]string
		expected []*model.DockSpec
	}
	tests := []test{
		// select by storage type
		{
			input: []*model.DockSpec{
				&SampleMultiDocks[0],
				&SampleMultiDocks[1],
				&SampleMultiDocks[2],
			},
			param: map[string][]string{
				"driverName": {"cinder"},
			},
			expected: []*model.DockSpec{
				&SampleMultiDocks[1],
			},
		},
		// sort by name asc
		{
			input: []*model.DockSpec{
				&SampleMultiDocks[0],
				&SampleMultiDocks[1],
				&SampleMultiDocks[2],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"asc"},
			},
			expected: []*model.DockSpec{
				&SampleMultiDocks[0],
				&SampleMultiDocks[2],
				&SampleMultiDocks[1],
			},
		},
		// sort by name desc
		{
			input: []*model.DockSpec{
				&SampleMultiDocks[0],
				&SampleMultiDocks[1],
				&SampleMultiDocks[2],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"desc"},
			},
			expected: []*model.DockSpec{
				&SampleMultiDocks[1],
				&SampleMultiDocks[2],
				&SampleMultiDocks[0],
			},
		},
		// limit is 1
		{
			input: []*model.DockSpec{
				&SampleMultiDocks[0],
				&SampleMultiDocks[1]},
			param: map[string][]string{
				"limit":  {"1"},
				"offset": {"1"},
			},
			expected: []*model.DockSpec{
				&SampleMultiDocks[1],
			},
		},
	}
	for _, testcase := range tests {
		ret := fc.FilterAndSort(testcase.input, testcase.param, sortableKeysMap[typeFileShareSnapshots])
		var res = []*model.DockSpec{}
		for _, data := range ret.([]interface{}) {
			res = append(res, data.(*model.DockSpec))
		}
		if !reflect.DeepEqual(res, testcase.expected) {
			var expected []model.DockSpec
			for _, value := range testcase.expected {
				expected = append(expected, *value)
			}
			var got []model.DockSpec
			for _, value := range res {
				got = append(got, *value)
			}
			t.Errorf("Expected %+v\n", expected)
			t.Errorf("Got %+v\n", got)
		}
	}
}

func TestListPoolsWithFilter(t *testing.T) {
	type test struct {
		input    []*model.StoragePoolSpec
		param    map[string][]string
		expected []*model.StoragePoolSpec
	}
	tests := []test{
		// select by storage type
		{
			input: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
				&SamplePools[2],
			},
			param: map[string][]string{
				"storageType": {"block"},
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
			},
		},
		// sort by name asc
		{
			input: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
				&SamplePools[2],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"asc"},
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[2],
				&SamplePools[0],
				&SamplePools[1],
			},
		},
		// sort by name desc
		{
			input: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
				&SamplePools[2],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"desc"},
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[1],
				&SamplePools[0],
				&SamplePools[2],
			},
		},
		// limit is 2
		{
			input: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
				&SamplePools[2],
			},
			param: map[string][]string{
				"limit":  {"2"},
				"offset": {"1"},
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[1],
				&SamplePools[2],
			},
		},
	}
	for _, testcase := range tests {
		ret := fc.FilterAndSort(testcase.input, testcase.param, sortableKeysMap[typePools])
		var res = []*model.StoragePoolSpec{}
		for _, data := range ret.([]interface{}) {
			res = append(res, data.(*model.StoragePoolSpec))
		}
		if !reflect.DeepEqual(res, testcase.expected) {
			var expected []model.StoragePoolSpec
			for _, value := range testcase.expected {
				expected = append(expected, *value)
			}
			var got []model.StoragePoolSpec
			for _, value := range res {
				got = append(got, *value)
			}
			t.Errorf("Expected %+v\n", expected)
			t.Errorf("Got %+v\n", got)
		}
	}
}

func TestListVolumesWithFilter(t *testing.T) {
	type test struct {
		input    []*model.VolumeSpec
		param    map[string][]string
		expected []*model.VolumeSpec
	}
	tests := []test{
		// select by storage type
		{
			input: []*model.VolumeSpec{
				&SampleMultiVolumes[0],
				&SampleMultiVolumes[1],
			},
			param: map[string][]string{
				"size": {"1"},
			},
			expected: []*model.VolumeSpec{
				&SampleMultiVolumes[1],
			},
		},
		// sort by name asc
		{
			input: []*model.VolumeSpec{
				&SampleMultiVolumes[0],
				&SampleMultiVolumes[1],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"asc"},
			},
			expected: []*model.VolumeSpec{
				&SampleMultiVolumes[0],
				&SampleMultiVolumes[1],
			},
		},
		// sort by name desc
		{
			input: []*model.VolumeSpec{
				&SampleMultiVolumes[0],
				&SampleMultiVolumes[1],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"desc"},
			},
			expected: []*model.VolumeSpec{
				&SampleMultiVolumes[1],
				&SampleMultiVolumes[0],
			},
		},
		// limit is 1
		{
			input: []*model.VolumeSpec{
				&SampleMultiVolumes[0],
				&SampleMultiVolumes[1],
			},
			param: map[string][]string{
				"limit":  {"1"},
				"offset": {"1"},
			},
			expected: []*model.VolumeSpec{
				&SampleMultiVolumes[1],
			},
		},
		// DurableName Filter
		{
			input: []*model.VolumeSpec{
				&SampleVolumeWithDurableName[0],
			},
			param: map[string][]string{
				"DurableName": {"6216b2326e974b5fb0b3d2af5cd6b25b"},
			},
			expected: []*model.VolumeSpec{
				&SampleVolumeWithDurableName[0],
			},
		},
	}
	for _, testcase := range tests {
		ret := fc.FilterAndSort(testcase.input, testcase.param, sortableKeysMap[typeVolumes])
		var res = []*model.VolumeSpec{}
		for _, data := range ret.([]interface{}) {
			res = append(res, data.(*model.VolumeSpec))
		}
		if !reflect.DeepEqual(res, testcase.expected) {
			var expected []model.VolumeSpec
			for _, value := range testcase.expected {
				expected = append(expected, *value)
			}
			var got []model.VolumeSpec
			for _, value := range res {
				got = append(got, *value)
			}
			t.Errorf("Expected %+v\n", expected)
			t.Errorf("Got %+v\n", got)
		}
	}
}
func TestListVolumeAttachmentsWithFilter(t *testing.T) {
	type test struct {
		input    []*model.VolumeAttachmentSpec
		param    map[string][]string
		expected []*model.VolumeAttachmentSpec
	}
	tests := []test{
		// select by storage type
		{
			input: []*model.VolumeAttachmentSpec{
				&SampleMultiAttachments[0],
				&SampleMultiAttachments[1],
			},
			param: map[string][]string{
				"status": {"attached"},
			},
			expected: []*model.VolumeAttachmentSpec{
				&SampleMultiAttachments[1],
			},
		},
		// sort by volume id asc
		{
			input: []*model.VolumeAttachmentSpec{
				&SampleMultiAttachments[0],
				&SampleMultiAttachments[1],
			},
			param: map[string][]string{
				"sortKey": {"volumeId"},
				"sortDir": {"asc"},
			},
			expected: []*model.VolumeAttachmentSpec{
				&SampleMultiAttachments[0],
				&SampleMultiAttachments[1],
			},
		},
		// sort by volume id desc
		{
			input: []*model.VolumeAttachmentSpec{
				&SampleMultiAttachments[0],
				&SampleMultiAttachments[1],
			},
			param: map[string][]string{
				"sortKey": {"volumeId"},
				"sortDir": {"desc"},
			},
			expected: []*model.VolumeAttachmentSpec{
				&SampleMultiAttachments[1],
				&SampleMultiAttachments[0],
			},
		},
		// limit is 1
		{
			input: []*model.VolumeAttachmentSpec{
				&SampleMultiAttachments[0],
				&SampleMultiAttachments[1],
			},
			param: map[string][]string{
				"limit":  {"1"},
				"offset": {"1"},
			},
			expected: []*model.VolumeAttachmentSpec{
				&SampleMultiAttachments[1],
			},
		},
	}
	for _, testcase := range tests {
		ret := fc.FilterAndSort(testcase.input, testcase.param, sortableKeysMap[typeAttachments])
		var res = []*model.VolumeAttachmentSpec{}
		for _, data := range ret.([]interface{}) {
			res = append(res, data.(*model.VolumeAttachmentSpec))
		}
		if !reflect.DeepEqual(res, testcase.expected) {
			var expected []model.VolumeAttachmentSpec
			for _, value := range testcase.expected {
				expected = append(expected, *value)
			}
			var got []model.VolumeAttachmentSpec
			for _, value := range res {
				got = append(got, *value)
			}
			t.Errorf("Expected %+v\n", expected)
			t.Errorf("Got %+v\n", got)
		}
	}
}
func TestListVolumeSnapshotsWithFilter(t *testing.T) {
	type test struct {
		input    []*model.VolumeSnapshotSpec
		param    map[string][]string
		expected []*model.VolumeSnapshotSpec
	}
	tests := []test{
		// select by storage type
		{
			input: []*model.VolumeSnapshotSpec{
				&SampleSnapshots[0],
				&SampleSnapshots[1],
			},
			param: map[string][]string{
				"status": {"available"},
			},
			expected: []*model.VolumeSnapshotSpec{
				&SampleSnapshots[0],
				&SampleSnapshots[1],
			},
		},
		// sort by name asc
		{
			input: []*model.VolumeSnapshotSpec{
				&SampleSnapshots[0],
				&SampleSnapshots[1],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"asc"},
			},
			expected: []*model.VolumeSnapshotSpec{
				&SampleSnapshots[0],
				&SampleSnapshots[1],
			},
		},
		// sort by name desc
		{
			input: []*model.VolumeSnapshotSpec{
				&SampleSnapshots[0],
				&SampleSnapshots[1],
			},
			param: map[string][]string{
				"sortKey": {"name"},
				"sortDir": {"desc"},
			},
			expected: []*model.VolumeSnapshotSpec{
				&SampleSnapshots[1],
				&SampleSnapshots[0],
			},
		},
		// limit is 1
		{
			input: []*model.VolumeSnapshotSpec{
				&SampleSnapshots[0],
				&SampleSnapshots[1],
			},
			param: map[string][]string{
				"limit":  {"1"},
				"offset": {"1"},
			},
			expected: []*model.VolumeSnapshotSpec{
				&SampleSnapshots[1],
			},
		},
	}
	for _, testcase := range tests {
		ret := fc.FilterAndSort(testcase.input, testcase.param, sortableKeysMap[typePools])
		var res = []*model.VolumeSnapshotSpec{}
		for _, data := range ret.([]interface{}) {
			res = append(res, data.(*model.VolumeSnapshotSpec))
		}
		if !reflect.DeepEqual(res, testcase.expected) {
			var expected []model.VolumeSnapshotSpec
			for _, value := range testcase.expected {
				expected = append(expected, *value)
			}
			var got []model.VolumeSnapshotSpec
			for _, value := range res {
				got = append(got, *value)
			}
			t.Errorf("Expected %+v\n", expected)
			t.Errorf("Got %+v\n", got)
		}
	}
}
