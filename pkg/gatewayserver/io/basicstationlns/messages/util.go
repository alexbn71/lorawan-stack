// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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

package messages

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"reflect"
	"strings"

	"go.thethings.network/lorawan-stack/pkg/band"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
)

var (
	errDataRateIndex = errors.DefineInvalidArgument("data_rate_index", "data rate index is out of range")
	errDataRate      = errors.DefineNotFound("data_rate", "data rate not found")
	errUID           = errors.DefineInvalidArgument("uid", "invalid uid `{uid}`")
)

func getInt32AsByteSlice(value int32) ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, value)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func getFCtrlAsUint(fCtrl ttnpb.FCtrl) uint {
	var ret uint
	if fCtrl.GetADR() {
		ret = ret | 0x80
	}
	if fCtrl.GetADRAckReq() {
		ret = ret | 0x40
	}
	if fCtrl.GetAck() {
		ret = ret | 0x20
	}
	if fCtrl.GetFPending() || fCtrl.GetClassB() {
		ret = ret | 0x10
	}
	return ret
}

func getDataRateFromIndex(bandID string, index int) (ttnpb.DataRate, bool, error) {
	band, err := band.GetByID(bandID)
	if err != nil {
		return ttnpb.DataRate{}, false, errDataRateIndex.WithCause(err)
	}
	if index >= len(band.DataRates) {
		return ttnpb.DataRate{}, false, errDataRateIndex
	}

	dr := band.DataRates[index].Rate

	if dr.GetLoRa() != nil {
		return dr, true, nil
	}

	return dr, false, nil
}

func getDataRateIndexFromDataRate(bandID string, DR ttnpb.DataRate) (int, error) {
	if (DR == ttnpb.DataRate{}) {
		return 0, errDataRate
	}
	band, err := band.GetByID(bandID)
	if err != nil {
		return 0, err
	}
	for i, dr := range band.DataRates {
		if reflect.DeepEqual(dr.Rate, DR) {
			return i, nil
		}
	}
	return 0, errDataRate
}

// getDataRatesFromBandID parses the available data rates from the band into DataRates.
func getDataRatesFromBandID(id string) (DataRates, error) {
	band, err := band.GetByID(id)
	if err != nil {
		return DataRates{}, err
	}

	// Set the default values.
	drs := DataRates{}
	for _, dr := range drs {
		dr[0] = -1
		dr[1] = 0
		dr[2] = 0
	}

	i := 0
	for _, dr := range band.DataRates {
		if loraDR := dr.Rate.GetLoRa(); loraDR != nil {
			loraDR.GetSpreadingFactor()
			drs[i][0] = int(loraDR.GetSpreadingFactor())
			drs[i][1] = int(loraDR.GetBandwidth() / 1000)
			i++
		} else if fskDR := dr.Rate.GetFSK(); fskDR != nil {
			drs[i][0] = 0 // must be set to 0 for FSK, the BW field is ignored.
			i++
		}
	}
	return drs, nil
}

// GetEUIfromUID extracts the Gateway EUI from the UID.
func GetEUIfromUID(uid string) (*types.EUI64, error) {
	s := strings.Split(uid, "-")
	if len(s) != 2 || s[0] != "eui" {
		return nil, errUID.WithAttributes("uid", uid)
	}
	euiBytes, err := hex.DecodeString(s[1])
	if err != nil {
		return nil, errUID.WithAttributes("uid", uid)
	}
	var eui types.EUI64
	err = eui.Unmarshal(euiBytes)
	if err != nil {
		return nil, errUID.WithAttributes("uid", uid)
	}
	return &eui, nil
}
