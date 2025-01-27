// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"
	time "time"

	types "github.com/gogo/protobuf/types"
)

func (dst *ApplicationPubSubIdentifiers) SetFields(src *ApplicationPubSubIdentifiers, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "application_ids":
			if len(subs) > 0 {
				newDst := &dst.ApplicationIdentifiers
				var newSrc *ApplicationIdentifiers
				if src != nil {
					newSrc = &src.ApplicationIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationIdentifiers = src.ApplicationIdentifiers
				} else {
					var zero ApplicationIdentifiers
					dst.ApplicationIdentifiers = zero
				}
			}
		case "pub_sub_id":
			if len(subs) > 0 {
				return fmt.Errorf("'pub_sub_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.PubSubID = src.PubSubID
			} else {
				var zero string
				dst.PubSubID = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationPubSub) SetFields(src *ApplicationPubSub, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "ids":
			if len(subs) > 0 {
				newDst := &dst.ApplicationPubSubIdentifiers
				var newSrc *ApplicationPubSubIdentifiers
				if src != nil {
					newSrc = &src.ApplicationPubSubIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationPubSubIdentifiers = src.ApplicationPubSubIdentifiers
				} else {
					var zero ApplicationPubSubIdentifiers
					dst.ApplicationPubSubIdentifiers = zero
				}
			}
		case "created_at":
			if len(subs) > 0 {
				return fmt.Errorf("'created_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CreatedAt = src.CreatedAt
			} else {
				var zero time.Time
				dst.CreatedAt = zero
			}
		case "updated_at":
			if len(subs) > 0 {
				return fmt.Errorf("'updated_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UpdatedAt = src.UpdatedAt
			} else {
				var zero time.Time
				dst.UpdatedAt = zero
			}
		case "format":
			if len(subs) > 0 {
				return fmt.Errorf("'format' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Format = src.Format
			} else {
				var zero string
				dst.Format = zero
			}
		case "base_topic":
			if len(subs) > 0 {
				return fmt.Errorf("'base_topic' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.BaseTopic = src.BaseTopic
			} else {
				var zero string
				dst.BaseTopic = zero
			}
		case "downlink_push":
			if len(subs) > 0 {
				newDst := dst.DownlinkPush
				if newDst == nil {
					newDst = &ApplicationPubSub_Message{}
					dst.DownlinkPush = newDst
				}
				var newSrc *ApplicationPubSub_Message
				if src != nil {
					newSrc = src.DownlinkPush
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkPush = src.DownlinkPush
				} else {
					dst.DownlinkPush = nil
				}
			}
		case "downlink_replace":
			if len(subs) > 0 {
				newDst := dst.DownlinkReplace
				if newDst == nil {
					newDst = &ApplicationPubSub_Message{}
					dst.DownlinkReplace = newDst
				}
				var newSrc *ApplicationPubSub_Message
				if src != nil {
					newSrc = src.DownlinkReplace
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkReplace = src.DownlinkReplace
				} else {
					dst.DownlinkReplace = nil
				}
			}
		case "uplink_message":
			if len(subs) > 0 {
				newDst := dst.UplinkMessage
				if newDst == nil {
					newDst = &ApplicationPubSub_Message{}
					dst.UplinkMessage = newDst
				}
				var newSrc *ApplicationPubSub_Message
				if src != nil {
					newSrc = src.UplinkMessage
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.UplinkMessage = src.UplinkMessage
				} else {
					dst.UplinkMessage = nil
				}
			}
		case "join_accept":
			if len(subs) > 0 {
				newDst := dst.JoinAccept
				if newDst == nil {
					newDst = &ApplicationPubSub_Message{}
					dst.JoinAccept = newDst
				}
				var newSrc *ApplicationPubSub_Message
				if src != nil {
					newSrc = src.JoinAccept
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.JoinAccept = src.JoinAccept
				} else {
					dst.JoinAccept = nil
				}
			}
		case "downlink_ack":
			if len(subs) > 0 {
				newDst := dst.DownlinkAck
				if newDst == nil {
					newDst = &ApplicationPubSub_Message{}
					dst.DownlinkAck = newDst
				}
				var newSrc *ApplicationPubSub_Message
				if src != nil {
					newSrc = src.DownlinkAck
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkAck = src.DownlinkAck
				} else {
					dst.DownlinkAck = nil
				}
			}
		case "downlink_nack":
			if len(subs) > 0 {
				newDst := dst.DownlinkNack
				if newDst == nil {
					newDst = &ApplicationPubSub_Message{}
					dst.DownlinkNack = newDst
				}
				var newSrc *ApplicationPubSub_Message
				if src != nil {
					newSrc = src.DownlinkNack
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkNack = src.DownlinkNack
				} else {
					dst.DownlinkNack = nil
				}
			}
		case "downlink_sent":
			if len(subs) > 0 {
				newDst := dst.DownlinkSent
				if newDst == nil {
					newDst = &ApplicationPubSub_Message{}
					dst.DownlinkSent = newDst
				}
				var newSrc *ApplicationPubSub_Message
				if src != nil {
					newSrc = src.DownlinkSent
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkSent = src.DownlinkSent
				} else {
					dst.DownlinkSent = nil
				}
			}
		case "downlink_failed":
			if len(subs) > 0 {
				newDst := dst.DownlinkFailed
				if newDst == nil {
					newDst = &ApplicationPubSub_Message{}
					dst.DownlinkFailed = newDst
				}
				var newSrc *ApplicationPubSub_Message
				if src != nil {
					newSrc = src.DownlinkFailed
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkFailed = src.DownlinkFailed
				} else {
					dst.DownlinkFailed = nil
				}
			}
		case "downlink_queued":
			if len(subs) > 0 {
				newDst := dst.DownlinkQueued
				if newDst == nil {
					newDst = &ApplicationPubSub_Message{}
					dst.DownlinkQueued = newDst
				}
				var newSrc *ApplicationPubSub_Message
				if src != nil {
					newSrc = src.DownlinkQueued
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkQueued = src.DownlinkQueued
				} else {
					dst.DownlinkQueued = nil
				}
			}
		case "location_solved":
			if len(subs) > 0 {
				newDst := dst.LocationSolved
				if newDst == nil {
					newDst = &ApplicationPubSub_Message{}
					dst.LocationSolved = newDst
				}
				var newSrc *ApplicationPubSub_Message
				if src != nil {
					newSrc = src.LocationSolved
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.LocationSolved = src.LocationSolved
				} else {
					dst.LocationSolved = nil
				}
			}

		case "provider":
			if len(subs) == 0 && src == nil {
				dst.Provider = nil
				continue
			} else if len(subs) == 0 {
				dst.Provider = src.Provider
				continue
			}

			subPathMap := _processPaths(subs)
			if len(subPathMap) > 1 {
				return fmt.Errorf("more than one field specified for oneof field '%s'", name)
			}
			for oneofName, oneofSubs := range subPathMap {
				switch oneofName {
				case "nats":
					if _, ok := dst.Provider.(*ApplicationPubSub_NATS); !ok {
						dst.Provider = &ApplicationPubSub_NATS{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Provider.(*ApplicationPubSub_NATS).NATS
						if newDst == nil {
							newDst = &ApplicationPubSub_NATSProvider{}
							dst.Provider.(*ApplicationPubSub_NATS).NATS = newDst
						}
						var newSrc *ApplicationPubSub_NATSProvider
						if src != nil {
							newSrc = src.GetNATS()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Provider.(*ApplicationPubSub_NATS).NATS = src.GetNATS()
						} else {
							dst.Provider.(*ApplicationPubSub_NATS).NATS = nil
						}
					}

				default:
					return fmt.Errorf("invalid oneof field: '%s.%s'", name, oneofName)
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationPubSubs) SetFields(src *ApplicationPubSubs, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "pubsubs":
			if len(subs) > 0 {
				return fmt.Errorf("'pubsubs' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Pubsubs = src.Pubsubs
			} else {
				dst.Pubsubs = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationPubSubFormats) SetFields(src *ApplicationPubSubFormats, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "formats":
			if len(subs) > 0 {
				return fmt.Errorf("'formats' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Formats = src.Formats
			} else {
				dst.Formats = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *GetApplicationPubSubRequest) SetFields(src *GetApplicationPubSubRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "ids":
			if len(subs) > 0 {
				newDst := &dst.ApplicationPubSubIdentifiers
				var newSrc *ApplicationPubSubIdentifiers
				if src != nil {
					newSrc = &src.ApplicationPubSubIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationPubSubIdentifiers = src.ApplicationPubSubIdentifiers
				} else {
					var zero ApplicationPubSubIdentifiers
					dst.ApplicationPubSubIdentifiers = zero
				}
			}
		case "field_mask":
			if len(subs) > 0 {
				return fmt.Errorf("'field_mask' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FieldMask = src.FieldMask
			} else {
				var zero types.FieldMask
				dst.FieldMask = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ListApplicationPubSubsRequest) SetFields(src *ListApplicationPubSubsRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "application_ids":
			if len(subs) > 0 {
				newDst := &dst.ApplicationIdentifiers
				var newSrc *ApplicationIdentifiers
				if src != nil {
					newSrc = &src.ApplicationIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationIdentifiers = src.ApplicationIdentifiers
				} else {
					var zero ApplicationIdentifiers
					dst.ApplicationIdentifiers = zero
				}
			}
		case "field_mask":
			if len(subs) > 0 {
				return fmt.Errorf("'field_mask' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FieldMask = src.FieldMask
			} else {
				var zero types.FieldMask
				dst.FieldMask = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *SetApplicationPubSubRequest) SetFields(src *SetApplicationPubSubRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "pubsub":
			if len(subs) > 0 {
				newDst := &dst.ApplicationPubSub
				var newSrc *ApplicationPubSub
				if src != nil {
					newSrc = &src.ApplicationPubSub
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationPubSub = src.ApplicationPubSub
				} else {
					var zero ApplicationPubSub
					dst.ApplicationPubSub = zero
				}
			}
		case "field_mask":
			if len(subs) > 0 {
				return fmt.Errorf("'field_mask' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FieldMask = src.FieldMask
			} else {
				var zero types.FieldMask
				dst.FieldMask = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationPubSub_NATSProvider) SetFields(src *ApplicationPubSub_NATSProvider, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "server_url":
			if len(subs) > 0 {
				return fmt.Errorf("'server_url' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ServerURL = src.ServerURL
			} else {
				var zero string
				dst.ServerURL = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationPubSub_Message) SetFields(src *ApplicationPubSub_Message, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "topic":
			if len(subs) > 0 {
				return fmt.Errorf("'topic' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Topic = src.Topic
			} else {
				var zero string
				dst.Topic = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
