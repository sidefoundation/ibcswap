package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v4/modules/core/exported"
)

// ChannelKeeper defines the expected IBC channel keeper.
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, portID, channelID string) (channeltypes.Channel, bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	ChanCloseInit(ctx sdk.Context, portID, channelID string, chanCap *capabilitytypes.Capability) error
}

// PortKeeper defines the expected IBC port keeper.
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
}

// ScopedKeeper defines the expected IBC scoped keeper.
type ScopedKeeper interface {
	GetCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, bool)
	AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool
	ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error
}

// ICS4Wrapper defines the expected ICS4Wrapper for middleware
type ICS4Wrapper interface {
	SendPacket(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet ibcexported.PacketI) error
}
