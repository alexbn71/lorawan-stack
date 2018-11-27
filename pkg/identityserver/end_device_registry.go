// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
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

package identityserver

import (
	"context"
	"strconv"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/pkg/events"
	"go.thethings.network/lorawan-stack/pkg/identityserver/blacklist"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	evtCreateEndDevice = events.Define("end_device.create", "Create end_device")
	evtUpdateEndDevice = events.Define("end_device.update", "Update end_device")
	evtDeleteEndDevice = events.Define("end_device.delete", "Delete end_device")
)

func (is *IdentityServer) createEndDevice(ctx context.Context, req *ttnpb.CreateEndDeviceRequest) (dev *ttnpb.EndDevice, err error) {
	err = rights.RequireApplication(ctx, req.EndDeviceIdentifiers.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_WRITE)
	if err != nil {
		return nil, err
	}
	if err := blacklist.Check(ctx, req.DeviceID); err != nil {
		return nil, err
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		devStore := store.GetEndDeviceStore(db)
		dev, err = devStore.CreateEndDevice(ctx, &req.EndDevice)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtCreateEndDevice(ctx, req.EndDeviceIdentifiers, nil))
	return dev, nil
}

func (is *IdentityServer) getEndDevice(ctx context.Context, req *ttnpb.GetEndDeviceRequest) (dev *ttnpb.EndDevice, err error) {
	err = rights.RequireApplication(ctx, req.EndDeviceIdentifiers.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_READ)
	if err != nil {
		return nil, err
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		devStore := store.GetEndDeviceStore(db)
		dev, err = devStore.GetEndDevice(ctx, &req.EndDeviceIdentifiers, &req.FieldMask)
		return err
	})
	if err != nil {
		return nil, err
	}
	return dev, nil
}

func (is *IdentityServer) listEndDevices(ctx context.Context, req *ttnpb.ListEndDevicesRequest) (devs *ttnpb.EndDevices, err error) {
	err = rights.RequireApplication(ctx, req.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_READ)
	if err != nil {
		return nil, err
	}
	var total uint64
	ctx = store.SetTotalCount(ctx, &total)
	defer func() {
		if err == nil {
			grpc.SetHeader(ctx, metadata.Pairs("x-total-count", strconv.FormatUint(total, 10)))
		}
	}()
	devs = new(ttnpb.EndDevices)
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		devStore := store.GetEndDeviceStore(db)
		devs.EndDevices, err = devStore.ListEndDevices(ctx, &req.ApplicationIdentifiers, &req.FieldMask)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return devs, nil
}

func (is *IdentityServer) updateEndDevice(ctx context.Context, req *ttnpb.UpdateEndDeviceRequest) (dev *ttnpb.EndDevice, err error) {
	err = rights.RequireApplication(ctx, req.EndDeviceIdentifiers.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_WRITE)
	if err != nil {
		return nil, err
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		devStore := store.GetEndDeviceStore(db)
		dev, err = devStore.UpdateEndDevice(ctx, &req.EndDevice, &req.FieldMask)
		return err
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtUpdateEndDevice(ctx, req.EndDeviceIdentifiers, req.FieldMask.Paths))
	return dev, nil
}

func (is *IdentityServer) deleteEndDevice(ctx context.Context, ids *ttnpb.EndDeviceIdentifiers) (*types.Empty, error) {
	err := rights.RequireApplication(ctx, ids.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_WRITE)
	if err != nil {
		return nil, err
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		devStore := store.GetEndDeviceStore(db)
		err = devStore.DeleteEndDevice(ctx, ids)
		return err
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtDeleteEndDevice(ctx, ids, nil))
	return ttnpb.Empty, nil
}

type endDeviceRegistry struct {
	*IdentityServer
}

func (dr *endDeviceRegistry) Create(ctx context.Context, req *ttnpb.CreateEndDeviceRequest) (*ttnpb.EndDevice, error) {
	return dr.createEndDevice(ctx, req)
}
func (dr *endDeviceRegistry) Get(ctx context.Context, req *ttnpb.GetEndDeviceRequest) (*ttnpb.EndDevice, error) {
	return dr.getEndDevice(ctx, req)
}
func (dr *endDeviceRegistry) List(ctx context.Context, req *ttnpb.ListEndDevicesRequest) (*ttnpb.EndDevices, error) {
	return dr.listEndDevices(ctx, req)
}
func (dr *endDeviceRegistry) Update(ctx context.Context, req *ttnpb.UpdateEndDeviceRequest) (*ttnpb.EndDevice, error) {
	return dr.updateEndDevice(ctx, req)
}
func (dr *endDeviceRegistry) Delete(ctx context.Context, req *ttnpb.EndDeviceIdentifiers) (*types.Empty, error) {
	return dr.deleteEndDevice(ctx, req)
}