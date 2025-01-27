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

package networkserver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/mohae/deepcopy"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/band"
	"go.thethings.network/lorawan-stack/pkg/cluster"
	"go.thethings.network/lorawan-stack/pkg/component"
	"go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/crypto"
	"go.thethings.network/lorawan-stack/pkg/encoding/lorawan"
	"go.thethings.network/lorawan-stack/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
	"google.golang.org/grpc"
)

func TestProcessDownlinkTask(t *testing.T) {
	getPaths := []string{
		"frequency_plan_id",
		"last_dev_status_received_at",
		"lorawan_phy_version",
		"mac_settings",
		"mac_state",
		"multicast",
		"pending_mac_state",
		"queued_application_downlinks",
		"recent_downlinks",
		"recent_uplinks",
		"session",
	}

	const appIDString = "process-downlink-test-app-id"
	appID := ttnpb.ApplicationIdentifiers{ApplicationID: appIDString}
	const devID = "process-downlink-test-dev-id"

	devAddr := types.DevAddr{0x42, 0xff, 0xff, 0xff}

	fNwkSIntKey := types.AES128Key{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	nwkSEncKey := types.AES128Key{0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	sNwkSIntKey := types.AES128Key{0x42, 0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	sessionKeys := &ttnpb.SessionKeys{
		FNwkSIntKey: &ttnpb.KeyEnvelope{
			Key: &fNwkSIntKey,
		},
		NwkSEncKey: &ttnpb.KeyEnvelope{
			Key: &nwkSEncKey,
		},
		SNwkSIntKey: &ttnpb.KeyEnvelope{
			Key: &sNwkSIntKey,
		},
	}

	rxMetadata := MakeRxMetadataSlice()
	eu868macParameters := &ttnpb.MACParameters{
		Channels: MakeEU868Channels(&ttnpb.MACParameters_Channel{
			UplinkFrequency:   430000000,
			DownlinkFrequency: 431000000,
			MinDataRateIndex:  ttnpb.DATA_RATE_0,
			MaxDataRateIndex:  ttnpb.DATA_RATE_3,
		}),
		Rx1Delay:          ttnpb.RX_DELAY_3,
		Rx1DataRateOffset: 2,
		Rx2DataRateIndex:  ttnpb.DATA_RATE_1,
		Rx2Frequency:      420000000,
	}

	assertGetRxMetadataGatewayPeers := func(ctx context.Context, getPeerCh <-chan test.ClusterGetPeerRequest, peer124, peer3 cluster.Peer) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()

		a := assertions.New(t)
		return test.AssertClusterGetPeerRequestSequence(ctx, getPeerCh,
			[]cluster.Peer{
				nil,
				peer124,
				peer124,
				peer3,
				peer124,
			},
			func(reqCtx context.Context, role ttnpb.PeerInfo_Role, ids ttnpb.Identifiers) bool {
				return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
					a.So(role, should.Equal, ttnpb.PeerInfo_GATEWAY_SERVER) &&
					a.So(ids, should.Resemble, ttnpb.GatewayIdentifiers{
						GatewayID: "gateway-test-0",
					})
			},
			func(reqCtx context.Context, role ttnpb.PeerInfo_Role, ids ttnpb.Identifiers) bool {
				return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
					a.So(role, should.Equal, ttnpb.PeerInfo_GATEWAY_SERVER) &&
					a.So(ids, should.Resemble, ttnpb.GatewayIdentifiers{
						GatewayID: "gateway-test-1",
					})
			},
			func(reqCtx context.Context, role ttnpb.PeerInfo_Role, ids ttnpb.Identifiers) bool {
				return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
					a.So(role, should.Equal, ttnpb.PeerInfo_GATEWAY_SERVER) &&
					a.So(ids, should.Resemble, ttnpb.GatewayIdentifiers{
						GatewayID: "gateway-test-2",
					})
			},
			func(reqCtx context.Context, role ttnpb.PeerInfo_Role, ids ttnpb.Identifiers) bool {
				return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
					a.So(role, should.Equal, ttnpb.PeerInfo_GATEWAY_SERVER) &&
					a.So(ids, should.Resemble, ttnpb.GatewayIdentifiers{
						GatewayID: "gateway-test-3",
					})
			},
			func(reqCtx context.Context, role ttnpb.PeerInfo_Role, ids ttnpb.Identifiers) bool {
				return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
					a.So(role, should.Equal, ttnpb.PeerInfo_GATEWAY_SERVER) &&
					a.So(ids, should.Resemble, ttnpb.GatewayIdentifiers{
						GatewayID: "gateway-test-4",
					})
			},
		)
	}

	assertScheduleRxMetadataGateways := func(ctx context.Context, authCh <-chan test.ClusterAuthRequest, scheduleDownlink124Ch, scheduleDownlink3Ch <-chan NsGsScheduleDownlinkRequest, payload []byte, makeTxRequest func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest, resp NsGsScheduleDownlinkResponse) (*ttnpb.DownlinkMessage, bool) {
		t := test.MustTFromContext(ctx)
		t.Helper()

		a := assertions.New(t)

		var correlationIDs []string
		if !a.So(AssertAuthNsGsScheduleDownlinkRequest(ctx, authCh, scheduleDownlink124Ch,
			func(ctx context.Context, msg *ttnpb.DownlinkMessage) bool {
				correlationIDs = msg.CorrelationIDs
				return a.So(msg, should.Resemble, &ttnpb.DownlinkMessage{
					CorrelationIDs: correlationIDs,
					RawPayload:     payload,
					Settings: &ttnpb.DownlinkMessage_Request{
						Request: makeTxRequest(
							&ttnpb.DownlinkPath{
								Path: &ttnpb.DownlinkPath_UplinkToken{
									UplinkToken: []byte("token-gtw-1"),
								},
							},
							&ttnpb.DownlinkPath{
								Path: &ttnpb.DownlinkPath_UplinkToken{
									UplinkToken: []byte("token-gtw-2"),
								},
							},
						),
					},
				})
			},
			grpc.EmptyCallOption{},
			NsGsScheduleDownlinkResponse{
				Error: errors.New("test"),
			},
		), should.BeTrue) {
			t.Error("Downlink assertion failed for gateways 1 and 2")
			return nil, false
		}
		t.Logf("Downlink correlation IDs: %v", correlationIDs)

		if !a.So(AssertAuthNsGsScheduleDownlinkRequest(ctx, authCh, scheduleDownlink3Ch,
			func(ctx context.Context, msg *ttnpb.DownlinkMessage) bool {
				return a.So(msg, should.Resemble, &ttnpb.DownlinkMessage{
					CorrelationIDs: correlationIDs,
					RawPayload:     payload,
					Settings: &ttnpb.DownlinkMessage_Request{
						Request: makeTxRequest(
							&ttnpb.DownlinkPath{
								Path: &ttnpb.DownlinkPath_UplinkToken{
									UplinkToken: []byte("token-gtw-3"),
								},
							},
						),
					},
				})
			},
			grpc.EmptyCallOption{},
			NsGsScheduleDownlinkResponse{
				Error: errors.New("test"),
			},
		), should.BeTrue) {
			t.Error("Downlink assertion failed for gateway 3")
			return nil, false
		}

		lastDown := &ttnpb.DownlinkMessage{
			CorrelationIDs: correlationIDs,
			RawPayload:     payload,
			Settings: &ttnpb.DownlinkMessage_Request{
				Request: makeTxRequest(
					&ttnpb.DownlinkPath{
						Path: &ttnpb.DownlinkPath_UplinkToken{
							UplinkToken: []byte("token-gtw-4"),
						},
					},
				),
			},
		}

		if !a.So(AssertAuthNsGsScheduleDownlinkRequest(ctx, authCh, scheduleDownlink124Ch,
			func(ctx context.Context, msg *ttnpb.DownlinkMessage) bool {
				return a.So(msg, should.Resemble, lastDown)
			},
			grpc.EmptyCallOption{},
			resp,
		), should.BeTrue) {
			t.Error("Downlink assertion failed for gateway 4")
			return nil, false
		}
		return lastDown, true
	}

	for _, tc := range []struct {
		Name               string
		DownlinkPriorities DownlinkPriorities
		Handler            func(context.Context, TestEnvironment) bool
		ErrorAssertion     func(*testing.T, error) bool
	}{
		{
			Name: "Class A/windows open/Rx1,2/application downlink/FOpts present/EU868/1.1",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, time.Now())
					}()
				}

				lastUp := &ttnpb.UplinkMessage{
					CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
					DeviceChannelIndex: 3,
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
					ReceivedAt: time.Now().Add(-time.Second),
					RxMetadata: deepcopy.Copy(rxMetadata).([]*ttnpb.RxMetadata),
					Settings: ttnpb.TxSettings{
						DataRateIndex: ttnpb.DATA_RATE_0,
						Frequency:     430000000,
					},
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RxWindowsAvailable: true,
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
							FCnt:           0x42,
							FPort:          0x1,
							FRMPayload:     []byte("testPayload"),
							Priority:       ttnpb.TxSchedulePriority_HIGHEST,
							SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
						},
					},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.Resemble, getPaths)

					go func() {
						dev, sets, err := req.Func(CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetRxMetadataGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0x60,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0x86,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* ResetConf */
								0x01, 0x01,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx1Delay:         ttnpb.RX_DELAY_3,
							Rx1DataRateIndex: ttnpb.DATA_RATE_0,
							Rx1Frequency:     431000000,
							Rx2DataRateIndex: ttnpb.DATA_RATE_1,
							Rx2Frequency:     420000000,
						}
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 5) {
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-2")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-app-down-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-app-down-2")
				}

				setDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						PendingRequests: []*ttnpb.MACCommand{
							{
								CID: ttnpb.CID_DEV_STATUS,
							},
						},
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					RecentDownlinks: []*ttnpb.DownlinkMessage{
						lastDown,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.Resemble, []string{
						"mac_state",
						"queued_application_downlinks",
						"recent_downlinks",
						"session",
					})
					if a.So(resp.Device, should.NotBeNil) &&
						a.So(resp.Device.MACState, should.NotBeNil) &&
						a.So(resp.Device.MACState.LastConfirmedDownlinkAt, should.NotBeNil) {
						a.So([]time.Time{start, *resp.Device.MACState.LastConfirmedDownlinkAt, time.Now().Add(time.Second)}, should.BeChronological)
						setDevice.MACState.LastConfirmedDownlinkAt = resp.Device.MACState.LastConfirmedDownlinkAt
					}
					a.So(resp.Device, should.Resemble, setDevice)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device: setDevice,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		// Adapted from https://github.com/TheThingsNetwork/lorawan-stack/issues/866#issue-461484955.
		{
			Name: "Class A/windows open/Rx1, Rx2 does not fit/application downlink/FOpts present/EU868/1.0.2",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, time.Now())
					}()
				}

				lastUp := &ttnpb.UplinkMessage{
					CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
					DeviceChannelIndex: 3,
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
					ReceivedAt: time.Now().Add(-time.Second),
					RxMetadata: deepcopy.Copy(rxMetadata).([]*ttnpb.RxMetadata),
					Settings: ttnpb.TxSettings{
						DataRateIndex: ttnpb.DATA_RATE_6,
						Frequency:     430000000,
					},
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RxWindowsAvailable: true,
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
							FCnt:           0x42,
							FPort:          0x15,
							FRMPayload:     []byte("AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUU="),
							Priority:       ttnpb.TxSchedulePriority_HIGHEST,
							SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
						},
					},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.Resemble, getPaths)

					go func() {
						dev, sets, err := req.Func(CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetRxMetadataGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0x60,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0x86,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* ResetConf */
								0x01, 0x01,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x15)

						/** FRMPayload **/
						b = append(b, []byte("AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUU=")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx1Delay:         ttnpb.RX_DELAY_3,
							Rx1DataRateIndex: ttnpb.DATA_RATE_4,
							Rx1Frequency:     431000000,
						}
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 5) {
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-2")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-app-down-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-app-down-2")
				}

				setDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						PendingRequests: []*ttnpb.MACCommand{
							{
								CID: ttnpb.CID_DEV_STATUS,
							},
						},
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					RecentDownlinks: []*ttnpb.DownlinkMessage{
						lastDown,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.Resemble, []string{
						"mac_state",
						"queued_application_downlinks",
						"recent_downlinks",
						"session",
					})
					if a.So(resp.Device, should.NotBeNil) &&
						a.So(resp.Device.MACState, should.NotBeNil) &&
						a.So(resp.Device.MACState.LastConfirmedDownlinkAt, should.NotBeNil) {
						a.So([]time.Time{start, *resp.Device.MACState.LastConfirmedDownlinkAt, time.Now().Add(time.Second)}, should.BeChronological)
						setDevice.MACState.LastConfirmedDownlinkAt = resp.Device.MACState.LastConfirmedDownlinkAt
					}
					a.So(resp.Device, should.Resemble, setDevice)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device: setDevice,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows open/Rx1/application downlink/FOpts present/EU868/1.1",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, time.Now())
					}()
				}

				lastUp := &ttnpb.UplinkMessage{
					CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
					DeviceChannelIndex: 3,
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
					ReceivedAt: time.Now().Add(-time.Second),
					RxMetadata: deepcopy.Copy(rxMetadata).([]*ttnpb.RxMetadata),
					Settings: ttnpb.TxSettings{
						DataRateIndex: ttnpb.DATA_RATE_0,
						Frequency:     430000000,
					},
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RxWindowsAvailable: true,
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
							FCnt:           0x42,
							FPort:          0x1,
							FRMPayload:     []byte("testPayload"),
							Priority:       ttnpb.TxSchedulePriority_HIGHEST,
							SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
						},
					},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.Resemble, getPaths)

					go func() {
						dev, sets, err := req.Func(CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetRxMetadataGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0x60,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0x86,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* ResetConf */
								0x01, 0x01,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx1Delay:         ttnpb.RX_DELAY_3,
							Rx1DataRateIndex: ttnpb.DATA_RATE_0,
							Rx1Frequency:     431000000,
						}
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 5) {
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-2")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-app-down-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-app-down-2")
				}

				setDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						PendingRequests: []*ttnpb.MACCommand{
							{
								CID: ttnpb.CID_DEV_STATUS,
							},
						},
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					RecentDownlinks: []*ttnpb.DownlinkMessage{
						lastDown,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.Resemble, []string{
						"mac_state",
						"queued_application_downlinks",
						"recent_downlinks",
						"session",
					})
					if a.So(resp.Device, should.NotBeNil) &&
						a.So(resp.Device.MACState, should.NotBeNil) &&
						a.So(resp.Device.MACState.LastConfirmedDownlinkAt, should.NotBeNil) {
						a.So([]time.Time{start, *resp.Device.MACState.LastConfirmedDownlinkAt, time.Now().Add(time.Second)}, should.BeChronological)
						setDevice.MACState.LastConfirmedDownlinkAt = resp.Device.MACState.LastConfirmedDownlinkAt
					}
					a.So(resp.Device, should.Resemble, setDevice)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device: setDevice,
				}:
				}

				if !AssertDownlinkTaskAddRequest(ctx, env.DownlinkTasks.Add, func(reqCtx context.Context, ids ttnpb.EndDeviceIdentifiers, t time.Time, replace bool) bool {
					return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
						a.So(ids, should.Resemble, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}) &&
						a.So(replace, should.BeTrue) &&
						a.So(t, should.Resemble, setDevice.MACState.LastConfirmedDownlinkAt.Add(42*time.Second))
				},
					nil,
				) {
					t.Error("Downlink task add assertion failed")
					return false
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows open/Rx2/application downlink/no absolute time/no forced gateways/windows open/FOpts present/EU868/1.1",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, time.Now())
					}()
				}

				lastUp := &ttnpb.UplinkMessage{
					CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
					DeviceChannelIndex: 3,
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
					ReceivedAt: time.Now().Add(-time.Second),
					RxMetadata: deepcopy.Copy(rxMetadata).([]*ttnpb.RxMetadata),
					Settings: ttnpb.TxSettings{
						DataRateIndex: ttnpb.DATA_RATE_0,
						Frequency:     430000000,
					},
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RxWindowsAvailable: true,
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
							FCnt:           0x42,
							FPort:          0x1,
							FRMPayload:     []byte("testPayload"),
							Priority:       ttnpb.TxSchedulePriority_HIGHEST,
							SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
						},
					},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.Resemble, getPaths)

					go func() {
						dev, sets, err := req.Func(CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetRxMetadataGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				_, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0x60,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0x86,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* ResetConf */
								0x01, 0x01,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx1Delay:         ttnpb.RX_DELAY_3,
							Rx1DataRateIndex: ttnpb.DATA_RATE_0,
							Rx1Frequency:     431000000,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if !a.So(assertGetRxMetadataGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0x60,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0x86,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* ResetConf */
								0x01, 0x01,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_C,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx2DataRateIndex: ttnpb.DATA_RATE_1,
							Rx2Frequency:     420000000,
						}
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 5) {
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-2")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-app-down-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-app-down-2")
				}

				setDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						PendingRequests: []*ttnpb.MACCommand{
							{
								CID: ttnpb.CID_DEV_STATUS,
							},
						},
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					RecentDownlinks: []*ttnpb.DownlinkMessage{
						lastDown,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.Resemble, []string{
						"mac_state",
						"queued_application_downlinks",
						"recent_downlinks",
						"session",
					})
					if a.So(resp.Device, should.NotBeNil) &&
						a.So(resp.Device.MACState, should.NotBeNil) &&
						a.So(resp.Device.MACState.LastConfirmedDownlinkAt, should.NotBeNil) {
						a.So([]time.Time{start, *resp.Device.MACState.LastConfirmedDownlinkAt, time.Now().Add(time.Second)}, should.BeChronological)
						setDevice.MACState.LastConfirmedDownlinkAt = resp.Device.MACState.LastConfirmedDownlinkAt
					}
					a.So(resp.Device, should.Resemble, setDevice)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device: setDevice,
				}:
				}

				if !AssertDownlinkTaskAddRequest(ctx, env.DownlinkTasks.Add, func(reqCtx context.Context, ids ttnpb.EndDeviceIdentifiers, t time.Time, replace bool) bool {
					return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
						a.So(ids, should.Resemble, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}) &&
						a.So(replace, should.BeTrue) &&
						a.So(t, should.Resemble, setDevice.MACState.LastConfirmedDownlinkAt.Add(42*time.Second))
				},
					nil,
				) {
					t.Error("Downlink task add assertion failed")
					return false
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows open/Rx2/application downlink/absolute time within window/no forced gateways/windows open/FOpts present/EU868/1.1",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, time.Now())
					}()
				}

				lastUp := &ttnpb.UplinkMessage{
					CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
					DeviceChannelIndex: 3,
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
					ReceivedAt: time.Now().Add(-time.Second),
					RxMetadata: deepcopy.Copy(rxMetadata).([]*ttnpb.RxMetadata),
					Settings: ttnpb.TxSettings{
						DataRateIndex: ttnpb.DATA_RATE_0,
						Frequency:     430000000,
					},
				}

				absTime := time.Now().Add(10 * time.Second).UTC()

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RxWindowsAvailable: true,
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
							FCnt:           0x42,
							FPort:          0x1,
							FRMPayload:     []byte("testPayload"),
							Priority:       ttnpb.TxSchedulePriority_HIGHEST,
							SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							ClassBC: &ttnpb.ApplicationDownlink_ClassBC{
								AbsoluteTime: deepcopy.Copy(&absTime).(*time.Time),
							},
						},
					},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.Resemble, getPaths)

					go func() {
						dev, sets, err := req.Func(CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetRxMetadataGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0x60,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0x86,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* ResetConf */
								0x01, 0x01,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_C,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx2DataRateIndex: ttnpb.DATA_RATE_1,
							Rx2Frequency:     420000000,
							AbsoluteTime:     &absTime,
						}
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 5) {
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-2")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-app-down-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-app-down-2")
				}

				setDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						PendingRequests: []*ttnpb.MACCommand{
							{
								CID: ttnpb.CID_DEV_STATUS,
							},
						},
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					RecentDownlinks: []*ttnpb.DownlinkMessage{
						lastDown,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.Resemble, []string{
						"mac_state",
						"queued_application_downlinks",
						"recent_downlinks",
						"session",
					})
					if a.So(resp.Device, should.NotBeNil) &&
						a.So(resp.Device.MACState, should.NotBeNil) &&
						a.So(resp.Device.MACState.LastConfirmedDownlinkAt, should.NotBeNil) {
						a.So([]time.Time{start, *resp.Device.MACState.LastConfirmedDownlinkAt, time.Now().Add(time.Second)}, should.BeChronological)
						setDevice.MACState.LastConfirmedDownlinkAt = resp.Device.MACState.LastConfirmedDownlinkAt
					}
					a.So(resp.Device, should.Resemble, setDevice)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device: setDevice,
				}:
				}

				if !AssertDownlinkTaskAddRequest(ctx, env.DownlinkTasks.Add, func(reqCtx context.Context, ids ttnpb.EndDeviceIdentifiers, t time.Time, replace bool) bool {
					return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
						a.So(ids, should.Resemble, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}) &&
						a.So(replace, should.BeTrue) &&
						a.So(t, should.Resemble, setDevice.MACState.LastConfirmedDownlinkAt.Add(42*time.Second))
				},
					nil,
				) {
					t.Error("Downlink task add assertion failed")
					return false
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows open/Rx2/application downlink/absolute time outside window",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, time.Now())
					}()
				}

				lastUp := &ttnpb.UplinkMessage{
					CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
					DeviceChannelIndex: 3,
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
					ReceivedAt: time.Now().Add(-time.Second),
					RxMetadata: deepcopy.Copy(rxMetadata).([]*ttnpb.RxMetadata),
					Settings: ttnpb.TxSettings{
						DataRateIndex: ttnpb.DATA_RATE_0,
						Frequency:     430000000,
					},
				}

				absTime := time.Now().Add(42 * time.Hour).UTC()

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RxWindowsAvailable: true,
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
							FCnt:           0x42,
							FPort:          0x1,
							FRMPayload:     []byte("testPayload"),
							Priority:       ttnpb.TxSchedulePriority_HIGHEST,
							SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							ClassBC: &ttnpb.ApplicationDownlink_ClassBC{
								AbsoluteTime: deepcopy.Copy(&absTime).(*time.Time),
							},
						},
					},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.Resemble, getPaths)

					go func() {
						dev, sets, err := req.Func(CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeEmpty)
					a.So(resp.Device, should.NotBeNil)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{}:
				}

				if !AssertDownlinkTaskAddRequest(ctx, env.DownlinkTasks.Add, func(reqCtx context.Context, ids ttnpb.EndDeviceIdentifiers, t time.Time, replace bool) bool {
					return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
						a.So(ids, should.Resemble, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}) &&
						a.So(replace, should.BeTrue) &&
						a.So(t, should.Resemble, absTime.Add(-gsScheduleWindow))
				},
					nil,
				) {
					t.Error("Downlink task add assertion failed")
					return false
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows open/Rx2/expired application downlinks",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, time.Now())
					}()
				}

				lastUp := &ttnpb.UplinkMessage{
					CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
					DeviceChannelIndex: 3,
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
					ReceivedAt: time.Now().Add(-time.Second),
					RxMetadata: deepcopy.Copy(rxMetadata).([]*ttnpb.RxMetadata),
					Settings: ttnpb.TxSettings{
						DataRateIndex: ttnpb.DATA_RATE_0,
						Frequency:     430000000,
					},
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout:          DurationPtr(42 * time.Second),
						StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 0},
						StatusTimePeriodicity:  DurationPtr(0),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters:  *CopyMACParameters(eu868macParameters),
						DesiredParameters:  *CopyMACParameters(eu868macParameters),
						DeviceClass:        ttnpb.CLASS_C,
						LoRaWANVersion:     ttnpb.MAC_V1_1,
						RxWindowsAvailable: true,
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
							FCnt:           0x42,
							FPort:          0x1,
							FRMPayload:     []byte("testPayload"),
							Priority:       ttnpb.TxSchedulePriority_HIGHEST,
							SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							ClassBC: &ttnpb.ApplicationDownlink_ClassBC{
								AbsoluteTime: TimePtr(time.Now().Add(-2).UTC()),
							},
						},
						{
							CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
							FCnt:           0x42,
							FPort:          0x1,
							FRMPayload:     []byte("testPayload"),
							Priority:       ttnpb.TxSchedulePriority_HIGHEST,
							SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							ClassBC: &ttnpb.ApplicationDownlink_ClassBC{
								AbsoluteTime: TimePtr(time.Now().Add(-1).UTC()),
							},
						},
					},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *CopySessionKeys(sessionKeys),
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.Resemble, getPaths)

					go func() {
						dev, sets, err := req.Func(CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.Resemble, []string{
						"queued_application_downlinks",
					})
					if resp.Device != nil {
						a.So(resp.Device.QueuedApplicationDownlinks, should.BeEmpty)
					} else {
						a.So(resp.Device, should.BeNil)
					}
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "join-accept/windows open/Rx2/no active MAC state/window open/EU868/1.1",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, time.Now())
					}()
				}

				lastUp := &ttnpb.UplinkMessage{
					CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
					DeviceChannelIndex: 3,
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_JOIN_REQUEST,
						},
						Payload: &ttnpb.Message_JoinRequestPayload{JoinRequestPayload: &ttnpb.JoinRequestPayload{
							JoinEUI:  types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
							DevEUI:   types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
							DevNonce: types.DevNonce{0x00, 0x42},
						}},
					},
					ReceivedAt: time.Now().Add(-time.Second),
					RxMetadata: deepcopy.Copy(rxMetadata).([]*ttnpb.RxMetadata),
					Settings: ttnpb.TxSettings{
						DataRateIndex: ttnpb.DATA_RATE_0,
						Frequency:     430000000,
					},
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						JoinEUI:                &types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					PendingMACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedJoinAccept: &ttnpb.MACState_JoinAccept{
							Keys:    *CopySessionKeys(sessionKeys),
							Payload: bytes.Repeat([]byte{0x42}, 33),
							Request: ttnpb.JoinRequest{
								DevAddr: devAddr,
							},
						},
						RxWindowsAvailable: true,
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
							FCnt:           0x42,
							FPort:          0x1,
							FRMPayload:     []byte("testPayload"),
							Priority:       ttnpb.TxSchedulePriority_HIGHEST,
							SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
						},
					},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					SupportsJoin: true,
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.Resemble, getPaths)

					go func() {
						dev, sets, err := req.Func(CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetRxMetadataGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					bytes.Repeat([]byte{0x42}, 33),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGHEST,
							Rx1Delay:         ttnpb.RX_DELAY_5,
							Rx1DataRateIndex: ttnpb.DATA_RATE_0,
							Rx1Frequency:     431000000,
							Rx2DataRateIndex: ttnpb.DATA_RATE_1,
							Rx2Frequency:     420000000,
						}
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 3) {
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-1")
					a.So(lastDown.CorrelationIDs, should.Contain, "correlation-up-2")
				}

				setDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						JoinEUI:                &types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					PendingMACState: &ttnpb.MACState{
						CurrentParameters: *CopyMACParameters(eu868macParameters),
						DesiredParameters: *CopyMACParameters(eu868macParameters),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						PendingJoinRequest: &ttnpb.JoinRequest{
							DevAddr: devAddr,
						},
					},
					PendingSession: &ttnpb.Session{
						DevAddr:     devAddr,
						SessionKeys: *CopySessionKeys(sessionKeys),
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
							FCnt:           0x42,
							FPort:          0x1,
							FRMPayload:     []byte("testPayload"),
							Priority:       ttnpb.TxSchedulePriority_HIGHEST,
							SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
						},
					},
					RecentUplinks: []*ttnpb.UplinkMessage{
						CopyUplinkMessage(lastUp),
					},
					RecentDownlinks: []*ttnpb.DownlinkMessage{
						lastDown,
					},
					SupportsJoin: true,
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.Resemble, []string{
						"pending_mac_state.pending_join_request",
						"pending_mac_state.queued_join_accept",
						"pending_mac_state.rx_windows_available",
						"pending_session.dev_addr",
						"pending_session.keys",
						"recent_downlinks",
					})
					a.So(resp.Device, should.Resemble, setDevice)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device: setDevice,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)

			ns, ctx, env, stop := StartTest(t, Config{}, (1<<10)*test.Delay)
			defer stop()

			ns.downlinkPriorities = tc.DownlinkPriorities

			go func() {
				for ev := range env.Events {
					t.Logf("Event %s published with data %v", ev.Event.Name(), ev.Event.Data())
					ev.Response <- struct{}{}
				}
			}()

			<-env.DownlinkTasks.Pop

			processDownlinkTaskErrCh := make(chan error)
			go func() {
				err := ns.processDownlinkTask(ctx)
				select {
				case <-ctx.Done():
					t.Log("NetworkServer.processDownlinkTask took too long to return")
					return

				default:
					processDownlinkTaskErrCh <- err
				}
			}()

			res := tc.Handler(ctx, env)
			if !a.So(res, should.BeTrue) {
				t.Error("Test handler failed")
				return
			}
			select {
			case <-ctx.Done():
				t.Error("Timed out while waiting for NetworkServer.processDownlinkTask to return")
				return

			case err := <-processDownlinkTaskErrCh:
				if tc.ErrorAssertion != nil {
					a.So(tc.ErrorAssertion(t, err), should.BeTrue)
				} else {
					a.So(err, should.BeNil)
				}
			}
			close(processDownlinkTaskErrCh)
		})
	}
}

func TestGenerateDownlink(t *testing.T) {
	phy := test.Must(test.Must(band.GetByID(band.EU_863_870)).(band.Band).Version(ttnpb.PHY_V1_1_REV_B)).(band.Band)

	const appIDString = "process-downlink-test-app-id"
	appID := ttnpb.ApplicationIdentifiers{ApplicationID: appIDString}
	const devID = "process-downlink-test-dev-id"

	devAddr := types.DevAddr{0x42, 0xff, 0xff, 0xff}

	fNwkSIntKey := types.AES128Key{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	nwkSEncKey := types.AES128Key{0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	sNwkSIntKey := types.AES128Key{0x42, 0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	encodeMessage := func(msg *ttnpb.Message, ver ttnpb.MACVersion, confFCnt uint32) []byte {
		msg = deepcopy.Copy(msg).(*ttnpb.Message)
		mac := msg.GetMACPayload()

		if len(mac.FRMPayload) > 0 && mac.FPort == 0 {
			var key types.AES128Key
			switch ver {
			case ttnpb.MAC_V1_0, ttnpb.MAC_V1_0_1, ttnpb.MAC_V1_0_2:
				key = fNwkSIntKey
			case ttnpb.MAC_V1_1:
				key = nwkSEncKey
			default:
				panic(fmt.Errorf("unknown version %s", ver))
			}

			var err error
			mac.FRMPayload, err = crypto.EncryptDownlink(key, mac.DevAddr, mac.FCnt, mac.FRMPayload)
			if err != nil {
				t.Fatal("Failed to encrypt downlink FRMPayload")
			}
		}

		b, err := lorawan.MarshalMessage(*msg)
		if err != nil {
			t.Fatal("Failed to marshal downlink")
		}

		var key types.AES128Key
		switch ver {
		case ttnpb.MAC_V1_0, ttnpb.MAC_V1_0_1, ttnpb.MAC_V1_0_2:
			key = fNwkSIntKey
		case ttnpb.MAC_V1_1:
			key = sNwkSIntKey
		default:
			panic(fmt.Errorf("unknown version %s", ver))
		}

		mic, err := crypto.ComputeDownlinkMIC(key, mac.DevAddr, confFCnt, mac.FCnt, b)
		if err != nil {
			t.Fatal("Failed to compute MIC")
		}
		return append(b, mic[:]...)
	}

	encodeMAC := func(phy band.Band, cmds ...*ttnpb.MACCommand) (b []byte) {
		for _, cmd := range cmds {
			b = test.Must(lorawan.DefaultMACCommands.AppendDownlink(phy, b, *cmd)).([]byte)
		}
		return
	}

	for _, tc := range []struct {
		Name                         string
		Device                       *ttnpb.EndDevice
		Bytes                        []byte
		ApplicationDownlinkAssertion func(t *testing.T, down *ttnpb.ApplicationDownlink) bool
		DeviceAssertion              func(*testing.T, *ttnpb.EndDevice) bool
		Error                        error
	}{
		{
			Name: "1.1/no app downlink/no MAC/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
				},
				Session:           ttnpb.NewPopulatedSession(test.Randy, false),
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
				RecentUplinks: []*ttnpb.UplinkMessage{{
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
				}},
			},
			Error: errNoDownlink,
		},
		{
			Name: "1.1/no app downlink/status after 1 downlink/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACSettings: &ttnpb.MACSettings{
					StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 3},
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion:      ttnpb.MAC_V1_1,
					LastDevStatusFCntUp: 2,
				},
				Session: &ttnpb.Session{
					LastFCntUp: 4,
				},
				LoRaWANPHYVersion:       ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:         band.EU_863_870,
				LastDevStatusReceivedAt: TimePtr(time.Unix(42, 0)),
				RecentUplinks: []*ttnpb.UplinkMessage{{
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
				}},
			},
			Error: errNoDownlink,
		},
		{
			Name: "1.1/no app downlink/status after an hour/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACSettings: &ttnpb.MACSettings{
					StatusTimePeriodicity: DurationPtr(24 * time.Hour),
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
				},
				LoRaWANPHYVersion:       ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:         band.EU_863_870,
				LastDevStatusReceivedAt: TimePtr(time.Now()),
				Session:                 ttnpb.NewPopulatedSession(test.Randy, false),
				RecentUplinks: []*ttnpb.UplinkMessage{{
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
				}},
			},
			Error: errNoDownlink,
		},
		{
			Name: "1.1/no app downlink/no MAC/ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion:     ttnpb.MAC_V1_1,
					RxWindowsAvailable: true,
				},
				Session: &ttnpb.Session{
					DevAddr:       devAddr,
					LastNFCntDown: 41,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
				RecentUplinks: []*ttnpb.UplinkMessage{{
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_CONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{
							MACPayload: &ttnpb.MACPayload{
								FHDR: ttnpb.FHDR{
									FCnt: 24,
								},
							},
						},
					},
				}},
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: true,
								ADR: true,
							},
							FCnt: 42,
						},
					},
				},
			}, ttnpb.MAC_V1_1, 24),
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				return assertions.New(t).So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion:     ttnpb.MAC_V1_1,
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 42,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_CONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{
								MACPayload: &ttnpb.MACPayload{
									FHDR: ttnpb.FHDR{
										FCnt: 24,
									},
								},
							},
						},
					}},
				})
			},
		},
		{
			Name: "1.1/unconfirmed app downlink/no MAC/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion:     ttnpb.MAC_V1_1,
					RxWindowsAvailable: true,
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
				},
				QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
					{
						Confirmed:  false,
						FCnt:       42,
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
				RecentUplinks: []*ttnpb.UplinkMessage{{
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
				}},
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: false,
								ADR: true,
							},
							FCnt: 42,
						},
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
			}, ttnpb.MAC_V1_1, 0),
			ApplicationDownlinkAssertion: func(t *testing.T, down *ttnpb.ApplicationDownlink) bool {
				return assertions.New(t).So(down, should.Resemble, &ttnpb.ApplicationDownlink{
					Confirmed:  false,
					FCnt:       42,
					FPort:      1,
					FRMPayload: []byte("test"),
				})
			},
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				return assertions.New(t).So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion:     ttnpb.MAC_V1_1,
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr: devAddr,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
				})
			},
		},
		{
			Name: "1.1/unconfirmed app downlink/no MAC/ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion:     ttnpb.MAC_V1_1,
					RxWindowsAvailable: true,
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
				},
				QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
					{
						Confirmed:  false,
						FCnt:       42,
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
				RecentUplinks: []*ttnpb.UplinkMessage{{
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_CONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{
							MACPayload: &ttnpb.MACPayload{
								FHDR: ttnpb.FHDR{
									FCnt: 24,
								},
							},
						},
					},
				}},
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: true,
								ADR: true,
							},
							FCnt: 42,
						},
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
			}, ttnpb.MAC_V1_1, 24),
			ApplicationDownlinkAssertion: func(t *testing.T, down *ttnpb.ApplicationDownlink) bool {
				return assertions.New(t).So(down, should.Resemble, &ttnpb.ApplicationDownlink{
					Confirmed:  false,
					FCnt:       42,
					FPort:      1,
					FRMPayload: []byte("test"),
				})
			},
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				return assertions.New(t).So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion:     ttnpb.MAC_V1_1,
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr: devAddr,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_CONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{
								MACPayload: &ttnpb.MACPayload{
									FHDR: ttnpb.FHDR{
										FCnt: 24,
									},
								},
							},
						},
					}},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
				})
			},
		},
		{
			Name: "1.1/confirmed app downlink/no MAC/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
				},
				QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
					{
						Confirmed:  true,
						FCnt:       42,
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
				RecentUplinks: []*ttnpb.UplinkMessage{{
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
				}},
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_CONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: false,
								ADR: true,
							},
							FCnt: 42,
						},
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
			}, ttnpb.MAC_V1_1, 0),
			ApplicationDownlinkAssertion: func(t *testing.T, down *ttnpb.ApplicationDownlink) bool {
				return assertions.New(t).So(down, should.Resemble, &ttnpb.ApplicationDownlink{
					Confirmed:  true,
					FCnt:       42,
					FPort:      1,
					FRMPayload: []byte("test"),
				})
			},
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				a := assertions.New(t)
				if !a.So(dev.MACState, should.NotBeNil) {
					t.FailNow()
				}
				return a.So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion: ttnpb.MAC_V1_1,
						PendingApplicationDownlink: &ttnpb.ApplicationDownlink{
							Confirmed:  true,
							FCnt:       42,
							FPort:      1,
							FRMPayload: []byte("test"),
						},
					},
					Session: &ttnpb.Session{
						DevAddr:          devAddr,
						LastConfFCntDown: 42,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
				})
			},
		},
		{
			Name: "1.1/confirmed app downlink/no MAC/ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion:     ttnpb.MAC_V1_1,
					RxWindowsAvailable: true,
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
				QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
					{
						Confirmed:  true,
						FCnt:       42,
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
				RecentUplinks: []*ttnpb.UplinkMessage{{
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_CONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{
							MACPayload: &ttnpb.MACPayload{
								FHDR: ttnpb.FHDR{
									FCnt: 24,
								},
							},
						},
					},
				}},
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_CONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: true,
								ADR: true,
							},
							FCnt: 42,
						},
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
			}, ttnpb.MAC_V1_1, 24),
			ApplicationDownlinkAssertion: func(t *testing.T, down *ttnpb.ApplicationDownlink) bool {
				return assertions.New(t).So(down, should.Resemble, &ttnpb.ApplicationDownlink{
					Confirmed:  true,
					FCnt:       42,
					FPort:      1,
					FRMPayload: []byte("test"),
				})
			},
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				a := assertions.New(t)
				if !a.So(dev.MACState, should.NotBeNil) {
					t.FailNow()
				}
				return a.So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion:     ttnpb.MAC_V1_1,
						RxWindowsAvailable: true,
						PendingApplicationDownlink: &ttnpb.ApplicationDownlink{
							Confirmed:  true,
							FCnt:       42,
							FPort:      1,
							FRMPayload: []byte("test"),
						},
					},
					Session: &ttnpb.Session{
						DevAddr:          devAddr,
						LastConfFCntDown: 42,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
					},
					LoRaWANPHYVersion:          ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:            band.EU_863_870,
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_CONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{
								MACPayload: &ttnpb.MACPayload{
									FHDR: ttnpb.FHDR{
										FCnt: 24,
									},
								},
							},
						},
					}},
				})
			},
		},
		{
			Name: "1.1/no app downlink/status(count)/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACSettings: &ttnpb.MACSettings{
					StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 3},
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion:      ttnpb.MAC_V1_1,
					LastDevStatusFCntUp: 4,
				},
				Session: &ttnpb.Session{
					DevAddr:       devAddr,
					LastFCntUp:    99,
					LastNFCntDown: 41,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
				RecentUplinks: []*ttnpb.UplinkMessage{{
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
				}},
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: false,
								ADR: true,
							},
							FCnt: 42,
						},
						FPort: 0,
						FRMPayload: encodeMAC(
							phy,
							ttnpb.CID_DEV_STATUS.MACCommand(),
						),
					},
				},
			}, ttnpb.MAC_V1_1, 0),
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				a := assertions.New(t)
				if !a.So(dev.MACState, should.NotBeNil) {
					t.FailNow()
				}
				return a.So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACSettings: &ttnpb.MACSettings{
						StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 3},
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion:      ttnpb.MAC_V1_1,
						LastDevStatusFCntUp: 4,
						PendingRequests: []*ttnpb.MACCommand{
							ttnpb.CID_DEV_STATUS.MACCommand(),
						},
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastFCntUp:    99,
						LastNFCntDown: 42,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
				})
			},
		},
		{
			Name: "1.1/no app downlink/status(time/zero time)/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACSettings: &ttnpb.MACSettings{
					StatusTimePeriodicity: DurationPtr(time.Nanosecond),
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
				},
				Session: &ttnpb.Session{
					DevAddr:       devAddr,
					LastNFCntDown: 41,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
				RecentUplinks: []*ttnpb.UplinkMessage{{
					Payload: &ttnpb.Message{
						MHDR: ttnpb.MHDR{
							MType: ttnpb.MType_UNCONFIRMED_UP,
						},
						Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
					},
				}},
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: false,
								ADR: true,
							},
							FCnt: 42,
						},
						FPort: 0,
						FRMPayload: encodeMAC(
							phy,
							ttnpb.CID_DEV_STATUS.MACCommand(),
						),
					},
				},
			}, ttnpb.MAC_V1_1, 0),
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				a := assertions.New(t)
				if !a.So(dev.MACState, should.NotBeNil) {
					t.FailNow()
				}
				return a.So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACSettings: &ttnpb.MACSettings{
						StatusTimePeriodicity: DurationPtr(time.Nanosecond),
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion: ttnpb.MAC_V1_1,
						PendingRequests: []*ttnpb.MACCommand{
							ttnpb.CID_DEV_STATUS.MACCommand(),
						},
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 42,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
				})
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)

			logger := test.GetLogger(t)

			ctx := test.ContextWithT(test.Context(), t)
			ctx = log.NewContext(ctx, logger)
			ctx, cancel := context.WithTimeout(ctx, (1<<7)*test.Delay)
			defer cancel()

			c := component.MustNew(
				log.Noop,
				&component.Config{},
				component.WithClusterNew(func(context.Context, *config.Cluster, ...cluster.Option) (cluster.Cluster, error) {
					return &test.MockCluster{
						JoinFunc: test.ClusterJoinNilFunc,
					}, nil
				}),
			)
			c.FrequencyPlans = frequencyplans.NewStore(test.FrequencyPlansFetcher)
			err := c.Start()
			a.So(err, should.BeNil)

			ns := &NetworkServer{
				Component: c,
				ctx:       ctx,
				defaultMACSettings: ttnpb.MACSettings{
					StatusTimePeriodicity:  DurationPtr(0),
					StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 0},
				},
			}

			dev := CopyEndDevice(tc.Device)
			_, phy, err := getDeviceBandVersion(dev, ns.FrequencyPlans)
			if !a.So(err, should.BeNil) {
				t.Fail()
				return
			}

			genDown, genState, err := ns.generateDownlink(ctx, dev, phy, math.MaxUint16, math.MaxUint16)
			if tc.Error != nil {
				a.So(err, should.EqualErrorOrDefinition, tc.Error)
				a.So(genDown, should.BeNil)
				return
			}
			// TODO: Assert AS uplinks generated(https://github.com/TheThingsNetwork/lorawan-stack/issues/631).
			_ = genState

			if !a.So(err, should.BeNil) || !a.So(genDown, should.NotBeNil) {
				t.Fail()
				return
			}

			a.So(genDown.Payload, should.Resemble, tc.Bytes)
			if tc.ApplicationDownlinkAssertion != nil {
				a.So(tc.ApplicationDownlinkAssertion(t, genDown.ApplicationDownlink), should.BeTrue)
			} else {
				a.So(genDown.ApplicationDownlink, should.BeNil)
			}

			if tc.DeviceAssertion != nil {
				a.So(tc.DeviceAssertion(t, dev), should.BeTrue)
			} else {
				a.So(dev, should.Resemble, tc.Device)
			}
		})
	}
}
