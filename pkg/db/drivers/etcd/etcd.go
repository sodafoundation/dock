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
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sodafoundation/dock/pkg/utils/config"

	log "github.com/golang/glog"
	c "github.com/sodafoundation/dock/pkg/context"
	"github.com/sodafoundation/dock/pkg/model"
	"github.com/sodafoundation/dock/pkg/utils"
	"github.com/sodafoundation/dock/pkg/utils/constants"
	"github.com/sodafoundation/dock/pkg/utils/urls"
	uuid "github.com/satori/go.uuid"
)

const (
	defaultSortKey          = "ID"
	defaultBlockProfileName = "default_block"
	defaultFileProfileName  = "default_file"
	typeBlock               = "block"
	typeFile                = "file"
)

var validKey = []string{"limit", "offset", "sortDir", "sortKey"}

const (
	typeFileShares         string = "FileShares"
	typeFileShareSnapshots string = "FileShareSnapshots"
	typeDocks              string = "Docks"
	typePools              string = "Pools"
	typeProfiles           string = "Profiles"
	typeVolumes            string = "Volumes"
	typeAttachments        string = "Attachments"
	typeVolumeSnapshots    string = "VolumeSnapshots"
)

var sortableKeysMap = map[string][]string{
	typeFileShares:         {"ID", "NAME", "STATUS", "AVAILABILITYZONE", "PROFILEID", "TENANTID", "SIZE", "POOLID", "DESCRIPTION"},
	typeFileShareSnapshots: {"ID", "NAME", "VOLUMEID", "STATUS", "USERID", "TENANTID", "SIZE"},
	typeDocks:              {"ID", "NAME", "STATUS", "ENDPOINT", "DRIVERNAME", "DESCRIPTION"},
	typePools:              {"ID", "NAME", "STATUS", "AVAILABILITYZONE", "DOCKID"},
	typeProfiles:           {"ID", "NAME", "DESCRIPTION"},
	typeVolumes:            {"ID", "NAME", "STATUS", "AVAILABILITYZONE", "PROFILEID", "TENANTID", "SIZE", "POOLID", "DESCRIPTION", "GROUPID"},
	typeAttachments:        {"ID", "VOLUMEID", "STATUS", "USERID", "TENANTID", "SIZE"},
	typeVolumeSnapshots:    {"ID", "NAME", "VOLUMEID", "STATUS", "USERID", "TENANTID", "SIZE"},
}

func IsAdminContext(ctx *c.Context) bool {
	return ctx.IsAdmin
}

func AuthorizeProjectContext(ctx *c.Context, tenantId string) bool {
	return ctx.TenantId == tenantId
}

// NewClient
func NewClient(etcd *config.Database) *Client {
	return &Client{
		clientInterface: Init(etcd),
	}
}

// Client
type Client struct {
	clientInterface
}

//Parameter
type Parameter struct {
	beginIdx, endIdx int
	sortDir, sortKey string
}

//IsInArray
func (c *Client) IsInArray(e string, s []string) bool {
	for _, v := range s {
		if strings.EqualFold(e, v) {
			return true
		}
	}
	return false
}

func (c *Client) SelectOrNot(m map[string][]string) bool {
	for key := range m {
		if !utils.Contained(key, validKey) {
			return true
		}
	}
	return false
}

//Get parameter limit
func (c *Client) GetLimit(m map[string][]string) int {
	var limit int
	var err error
	v, ok := m["limit"]
	if ok {
		limit, err = strconv.Atoi(v[0])
		if err != nil || limit < 0 {
			log.Warning("Invalid input limit:", limit, ",use default value instead:50")
			return constants.DefaultLimit
		}
	} else {
		log.Warning("The parameter limit is not present,use default value instead:50")
		return constants.DefaultLimit
	}
	return limit
}

//Get parameter offset
func (c *Client) GetOffset(m map[string][]string, size int) int {

	var offset int
	var err error
	v, ok := m["offset"]
	if ok {
		offset, err = strconv.Atoi(v[0])

		if err != nil || offset > size {
			log.Warning("Invalid input offset or input offset is out of bounds:", offset, ", use largest offset: size")

			return size
		}
		if offset < 0 {
			log.Warning("input offset is less than zero:", offset, ",use offset value instead:0")

			return constants.DefaultOffset
		}

	} else {
		log.Warning("The parameter offset is not present,use default value instead:0")
		return constants.DefaultOffset
	}
	return offset
}

//Get parameter sortDir
func (c *Client) GetSortDir(m map[string][]string) string {
	var sortDir string
	v, ok := m["sortDir"]
	if ok {
		sortDir = v[0]
		if !strings.EqualFold(sortDir, "desc") && !strings.EqualFold(sortDir, "asc") {
			log.Warning("Invalid input sortDir:", sortDir, ",use default value instead:desc")
			return constants.DefaultSortDir
		}
	} else {
		log.Warning("The parameter sortDir is not present,use default value instead:desc")
		return constants.DefaultSortDir
	}
	return sortDir
}

//Get parameter sortKey
func (c *Client) GetSortKey(m map[string][]string, sortKeys []string) string {
	var sortKey string
	v, ok := m["sortKey"]
	if ok {
		sortKey = strings.ToUpper(v[0])
		if !c.IsInArray(sortKey, sortKeys) {
			log.Warning("Invalid input sortKey:", sortKey, ",use default value instead:ID")
			return defaultSortKey
		}

	} else {
		log.Warning("The parameter sortKey is not present,use default value instead:ID")
		return defaultSortKey
	}
	return sortKey
}

func (c *Client) FilterAndSort(src interface{}, params map[string][]string, sortableKeys []string) interface{} {
	var ret interface{}
	ret = utils.Filter(src, params)
	if len(params["sortKey"]) > 0 && utils.ContainsIgnoreCase(sortableKeys, params["sortKey"][0]) {
		ret = utils.Sort(ret, params["sortKey"][0], c.GetSortDir(params))
	}
	ret = utils.Slice(ret, c.GetOffset(params, reflect.ValueOf(src).Len()), c.GetLimit(params))
	return ret
}

//ParameterFilter
func (c *Client) ParameterFilter(m map[string][]string, size int, sortKeys []string) *Parameter {

	limit := c.GetLimit(m)
	offset := c.GetOffset(m, size)

	beginIdx := offset
	endIdx := limit + offset

	// If use not specified the limit return all the items.
	if limit == constants.DefaultLimit || endIdx > size {
		endIdx = size
	}

	sortDir := c.GetSortDir(m)
	sortKey := c.GetSortKey(m, sortKeys)

	return &Parameter{beginIdx, endIdx, sortDir, sortKey}
}

// CreateDock
func (c *Client) CreateDock(ctx *c.Context, dck *model.DockSpec) (*model.DockSpec, error) {
	if dck.Id == "" {
		dck.Id = uuid.NewV4().String()
	}

	if dck.CreatedAt == "" {
		dck.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	dckBody, err := json.Marshal(dck)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateDockURL(urls.Etcd, "", dck.Id),
		Content: string(dckBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when create dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return dck, nil
}

// GetDock
func (c *Client) GetDock(ctx *c.Context, dckID string) (*model.DockSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateDockURL(urls.Etcd, "", dckID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when get dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var dck = &model.DockSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), dck); err != nil {
		log.Error("when parsing dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return dck, nil
}

// GetDockByPoolId
func (c *Client) GetDockByPoolId(ctx *c.Context, poolId string) (*model.DockSpec, error) {
	pool, err := c.GetPool(ctx, poolId)
	if err != nil {
		log.Error("Get pool failed in db: ", err)
		return nil, err
	}

	docks, err := c.ListDocks(ctx)
	if err != nil {
		log.Error("List docks failed failed in db: ", err)
		return nil, err
	}
	for _, dock := range docks {
		if pool.DockId == dock.Id {
			return dock, nil
		}
	}
	return nil, errors.New("Get dock failed by pool id: " + poolId)
}

// ListDocks
func (c *Client) ListDocks(ctx *c.Context) ([]*model.DockSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateDockURL(urls.Etcd, ""),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list docks in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var dcks = []*model.DockSpec{}
	if len(dbRes.Message) == 0 {
		return dcks, nil
	}
	for _, msg := range dbRes.Message {
		var dck = &model.DockSpec{}
		if err := json.Unmarshal([]byte(msg), dck); err != nil {
			log.Error("When parsing dock in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		dcks = append(dcks, dck)
	}
	return dcks, nil
}

func (c *Client) ListDocksWithFilter(ctx *c.Context, m map[string][]string) ([]*model.DockSpec, error) {
	docks, err := c.ListDocks(ctx)
	if err != nil {
		log.Error("List docks failed: ", err.Error())
		return nil, err
	}

	tmpDocks := c.FilterAndSort(docks, m, sortableKeysMap[typeDocks])
	var res = []*model.DockSpec{}
	for _, data := range tmpDocks.([]interface{}) {
		res = append(res, data.(*model.DockSpec))
	}
	return res, nil
}

// UpdateDock
func (c *Client) UpdateDock(ctx *c.Context, dckID, name, desp string) (*model.DockSpec, error) {
	dck, err := c.GetDock(ctx, dckID)
	if err != nil {
		return nil, err
	}
	if name != "" {
		dck.Name = name
	}
	if desp != "" {
		dck.Description = desp
	}
	dck.UpdatedAt = time.Now().Format(constants.TimeFormat)

	dckBody, err := json.Marshal(dck)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GenerateDockURL(urls.Etcd, "", dckID),
		NewContent: string(dckBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return dck, nil
}

// DeleteDock
func (c *Client) DeleteDock(ctx *c.Context, dckID string) error {
	dbReq := &Request{
		Url: urls.GenerateDockURL(urls.Etcd, "", dckID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete dock in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// CreatePool
func (c *Client) CreatePool(ctx *c.Context, pol *model.StoragePoolSpec) (*model.StoragePoolSpec, error) {
	if pol.Id == "" {
		pol.Id = uuid.NewV4().String()
	}

	if pol.CreatedAt == "" {
		pol.CreatedAt = time.Now().Format(constants.TimeFormat)
	}
	polBody, err := json.Marshal(pol)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GeneratePoolURL(urls.Etcd, "", pol.Id),
		Content: string(polBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create pol in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return pol, nil
}

func (c *Client) ListPoolsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.StoragePoolSpec, error) {
	pools, err := c.ListPools(ctx)
	if err != nil {
		log.Error("List pools failed: ", err.Error())
		return nil, err
	}

	tmpPools := c.FilterAndSort(pools, m, sortableKeysMap[typePools])
	var res = []*model.StoragePoolSpec{}
	for _, data := range tmpPools.([]interface{}) {
		res = append(res, data.(*model.StoragePoolSpec))
	}
	return res, nil
}

// GetPool
func (c *Client) GetPool(ctx *c.Context, polID string) (*model.StoragePoolSpec, error) {
	dbReq := &Request{
		Url: urls.GeneratePoolURL(urls.Etcd, "", polID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get pool in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var pol = &model.StoragePoolSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), pol); err != nil {
		log.Error("When parsing pool in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return pol, nil
}

// ListPools
func (c *Client) ListPools(ctx *c.Context) ([]*model.StoragePoolSpec, error) {
	dbReq := &Request{
		Url: urls.GeneratePoolURL(urls.Etcd, ""),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list pools in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var pols = []*model.StoragePoolSpec{}
	if len(dbRes.Message) == 0 {
		return pols, nil
	}
	for _, msg := range dbRes.Message {
		var pol = &model.StoragePoolSpec{}
		if err := json.Unmarshal([]byte(msg), pol); err != nil {
			log.Error("When parsing pool in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

// UpdatePool
func (c *Client) UpdatePool(ctx *c.Context, polID, name, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error) {
	pol, err := c.GetPool(ctx, polID)
	if err != nil {
		return nil, err
	}
	if name != "" {
		pol.Name = name
	}
	if desp != "" {
		pol.Description = desp
	}
	pol.UpdatedAt = time.Now().Format(constants.TimeFormat)

	polBody, err := json.Marshal(pol)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GeneratePoolURL(urls.Etcd, "", polID),
		NewContent: string(polBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update pool in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return pol, nil
}

// DeletePool
func (c *Client) DeletePool(ctx *c.Context, polID string) error {
	dbReq := &Request{
		Url: urls.GeneratePoolURL(urls.Etcd, "", polID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete pool in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// GetVolume
func (c *Client) GetVolume(ctx *c.Context, volID string) (*model.VolumeSpec, error) {
	vol, err := c.getVolume(ctx, volID)
	if !IsAdminContext(ctx) || err == nil {
		return vol, err
	}
	vols, err := c.ListVolumes(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range vols {
		if v.Id == volID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("specified volume(%s) can't find", volID)
}

func (c *Client) getVolume(ctx *c.Context, volID string) (*model.VolumeSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeURL(urls.Etcd, ctx.TenantId, volID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vol = &model.VolumeSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vol); err != nil {
		log.Error("When parsing volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return vol, nil
}

// ListVolumes
func (c *Client) ListVolumes(ctx *c.Context) ([]*model.VolumeSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeURL(urls.Etcd, ctx.TenantId),
	}

	// Admin user should get all volumes including the volumes whose tenant is not admin.
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateVolumeURL(urls.Etcd, "")
	}

	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list volumes in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vols = []*model.VolumeSpec{}
	if len(dbRes.Message) == 0 {
		return vols, nil
	}
	for _, msg := range dbRes.Message {
		var vol = &model.VolumeSpec{}
		if err := json.Unmarshal([]byte(msg), vol); err != nil {
			log.Error("When parsing volume in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		vols = append(vols, vol)
	}
	return vols, nil
}

// UpdateVolume ...
func (c *Client) UpdateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	result, err := c.GetVolume(ctx, vol.Id)
	if err != nil {
		return nil, err
	}
	if vol.Name != "" {
		result.Name = vol.Name
	}
	if vol.AvailabilityZone != "" {
		result.AvailabilityZone = vol.AvailabilityZone
	}
	if vol.Description != "" {
		result.Description = vol.Description
	}
	if vol.Metadata != nil {
		result.Metadata = utils.MergeStringMaps(result.Metadata, vol.Metadata)
	}
	if vol.Identifier != nil {
		result.Identifier = vol.Identifier
	}
	if vol.PoolId != "" {
		result.PoolId = vol.PoolId
	}
	if vol.Size != 0 {
		result.Size = vol.Size
	}
	if vol.Status != "" {
		result.Status = vol.Status
	}
	if vol.ReplicationDriverData != nil {
		result.ReplicationDriverData = vol.ReplicationDriverData
	}
	if vol.MultiAttach {
		result.MultiAttach = vol.MultiAttach
	}
	if vol.GroupId != "" {
		result.GroupId = vol.GroupId
	}

	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	body, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateVolumeURL(urls.Etcd, result.TenantId, vol.Id),
		NewContent: string(body),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

// DeleteVolume
func (c *Client) DeleteVolume(ctx *c.Context, volID string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		vol, err := c.GetVolume(ctx, volID)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = vol.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateVolumeURL(urls.Etcd, tenantId, volID),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *Client) filterByName(param map[string][]string, spec interface{}, filterList map[string]interface{}) bool {
	v := reflect.ValueOf(spec)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for key := range param {
		_, ok := filterList[key]
		if !ok {
			continue
		}
		filed := v.FieldByName(key)
		if !filed.IsValid() {
			continue
		}
		paramVal := param[key][0]
		var val string
		switch filed.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val = strconv.FormatInt(filed.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val = strconv.FormatUint(filed.Uint(), 10)
		case reflect.String:
			val = filed.String()
		default:
			return false
		}
		if !strings.EqualFold(paramVal, val) {
			return false
		}
	}

	return true
}

func (c *Client) CreateVolumeGroup(ctx *c.Context, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	vg.TenantId = ctx.TenantId
	vgBody, err := json.Marshal(vg)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId, vg.Id),
		Content: string(vgBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return vg, nil
}

func (c *Client) GetVolumeGroup(ctx *c.Context, vgId string) (*model.VolumeGroupSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId, vgId),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vg = &model.VolumeGroupSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vg); err != nil {
		log.Error("When parsing volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return vg, nil
}

func (c *Client) UpdateVolumeGroup(ctx *c.Context, vgUpdate *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	vg, err := c.GetVolumeGroup(ctx, vgUpdate.Id)
	if err != nil {
		return nil, err
	}
	if vgUpdate.Name != "" && vgUpdate.Name != vg.Name {
		vg.Name = vgUpdate.Name
	}
	if vgUpdate.AvailabilityZone != "" && vgUpdate.AvailabilityZone != vg.AvailabilityZone {
		vg.AvailabilityZone = vgUpdate.AvailabilityZone
	}
	if vgUpdate.Description != "" && vgUpdate.Description != vg.Description {
		vg.Description = vgUpdate.Description
	}
	if vgUpdate.PoolId != "" && vgUpdate.PoolId != vg.PoolId {
		vg.PoolId = vgUpdate.PoolId
	}
	if vg.Status != "" && vgUpdate.Status != vg.Status {
		vg.Status = vgUpdate.Status
	}
	if vgUpdate.PoolId != "" && vgUpdate.PoolId != vg.PoolId {
		vg.PoolId = vgUpdate.PoolId
	}
	if vgUpdate.CreatedAt != "" && vgUpdate.CreatedAt != vg.CreatedAt {
		vg.CreatedAt = vgUpdate.CreatedAt
	}
	if vgUpdate.UpdatedAt != "" && vgUpdate.UpdatedAt != vg.UpdatedAt {
		vg.UpdatedAt = vgUpdate.UpdatedAt
	}

	vgBody, err := json.Marshal(vg)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId, vgUpdate.Id),
		NewContent: string(vgBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return vg, nil
}

func (c *Client) UpdateStatus(ctx *c.Context, in interface{}, status string) error {
	switch in.(type) {
	case *model.VolumeSpec:
		volume := in.(*model.VolumeSpec)
		volume.Status = status
		if _, errUpdate := c.UpdateVolume(ctx, volume); errUpdate != nil {
			log.Error("When update volume status in db:", errUpdate.Error())
			return errUpdate
		}

	case *model.VolumeGroupSpec:
		vg := in.(*model.VolumeGroupSpec)
		vg.Status = status
		if _, errUpdate := c.UpdateVolumeGroup(ctx, vg); errUpdate != nil {
			log.Error("When update volume status in db:", errUpdate.Error())
			return errUpdate
		}

	case []*model.VolumeSpec:
		vols := in.([]*model.VolumeSpec)
		if _, errUpdate := c.VolumesToUpdate(ctx, vols); errUpdate != nil {
			return errUpdate
		}
	}
	return nil
}

func (c *Client) ListVolumesByGroupId(ctx *c.Context, vgId string) ([]*model.VolumeSpec, error) {
	volumes, err := c.ListVolumes(ctx)
	if err != nil {
		return nil, err
	}

	var volumesInSameGroup []*model.VolumeSpec
	for _, v := range volumes {
		if v.GroupId == vgId {
			volumesInSameGroup = append(volumesInSameGroup, v)
		}
	}

	return volumesInSameGroup, nil
}

func (c *Client) VolumesToUpdate(ctx *c.Context, volumeList []*model.VolumeSpec) ([]*model.VolumeSpec, error) {
	var volumeRefs []*model.VolumeSpec
	for _, values := range volumeList {
		v, err := c.UpdateVolume(ctx, values)
		if err != nil {
			return nil, err
		}
		volumeRefs = append(volumeRefs, v)
	}
	return volumeRefs, nil
}
