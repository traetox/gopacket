// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"code.google.com/p/gopacket"
	"encoding/binary"
	"fmt"
)

// LLDPTLVType is the type of each TLV value in a LinkLayerDiscovery packet.
type LLDPTLVType byte

const (
	LLDPTLVEnd             LLDPTLVType = 0
	LLDPTLVChassisID       LLDPTLVType = 1
	LLDPTLVPortID          LLDPTLVType = 2
	LLDPTLVTTL             LLDPTLVType = 3
	LLDPTLVPortDescription LLDPTLVType = 4
	LLDPTLVSysName         LLDPTLVType = 5
	LLDPTLVSysDescription  LLDPTLVType = 6
	LLDPTLVSysCapabilities LLDPTLVType = 7
	LLDPTLVMgmtAddress     LLDPTLVType = 8
	LLDPTLVOrgSpecific     LLDPTLVType = 127
)

// LinkLayerDiscoveryValue is a TLV value inside a LinkLayerDiscovery packet layer.
type LinkLayerDiscoveryValue struct {
	Type   LLDPTLVType
	Length uint16
	Value  []byte
}

// LLDPChassisIDSubType specifies the value type for a single LLDPChassisID.ID
type LLDPChassisIDSubType byte

const (
	LLDPChassisIDSubTypeReserved    LLDPChassisIDSubType = 0
	LLDPChassisIDSubTypeChassisComp LLDPChassisIDSubType = 1
	LLDPChassisIDSubtypeIfaceAlias  LLDPChassisIDSubType = 2
	LLDPChassisIDSubTypePortComp    LLDPChassisIDSubType = 3
	LLDPChassisIDSubTypeMACAddr     LLDPChassisIDSubType = 4
	LLDPChassisIDSubTypeNetworkAddr LLDPChassisIDSubType = 5
	LLDPChassisIDSubtypeIfaceName   LLDPChassisIDSubType = 6
	LLDPChassisIDSubTypeLocal       LLDPChassisIDSubType = 7
)

type LLDPChassisID struct {
	Subtype LLDPChassisIDSubType
	ID      []byte
}

// LLDPPortIDSubType specifies the value type for a single LLDPPortID.ID
type LLDPPortIDSubType byte

const (
	LLDPPortIDSubtypeReserved    LLDPPortIDSubType = 0
	LLDPPortIDSubtypeIfaceAlias  LLDPPortIDSubType = 1
	LLDPPortIDSubTypePortComp    LLDPPortIDSubType = 2
	LLDPPortIDSubTypeMACAddr     LLDPPortIDSubType = 3
	LLDPPortIDSubTypeNetworkAddr LLDPPortIDSubType = 4
	LLDPPortIDSubtypeIfaceName   LLDPPortIDSubType = 5
	LLDPPortIDSubTypeAgentCircuitID LLDPPortIDSubType = 6
	LLDPPortIDSubTypeLocal       LLDPPortIDSubType = 7
)

type LLDPPortID struct {
	Subtype LLDPPortIDSubType
	ID      []byte
}

// LinkLayerDiscovery is a packet layer containing the LinkLayer Discovery Protocol.
// See http:http://standards.ieee.org/getieee802/download/802.1AB-2009.pdf
// ChassisID, PortID and TTL are mandatory TLV's. Other values can be decoded
// with DecodeValues()
type LinkLayerDiscovery struct {
	baseLayer
	ChassisID LLDPChassisID
	PortID    LLDPPortID
	TTL       uint16
	Values    []LinkLayerDiscoveryValue
}

// LLDPOrgSpecificTLV is an Organisation-specific TLV
type LLDPOrgSpecificTLV struct {
	OUI     IEEEOUI
	SubType uint8
	Info    []byte
}

// LLDPCapabilities Types
const (
	LLDPCapsOther       uint16 = 1 << 0
	LLDPCapsRepeater    uint16 = 1 << 1
	LLDPCapsBridge      uint16 = 1 << 2
	LLDPCapsWLANAP      uint16 = 1 << 3
	LLDPCapsRouter      uint16 = 1 << 4
	LLDPCapsPhone       uint16 = 1 << 5
	LLDPCapsDocSis      uint16 = 1 << 6
	LLDPCapsStationOnly uint16 = 1 << 7
	LLDPCapsCVLAN       uint16 = 1 << 8
	LLDPCapsSVLAN       uint16 = 1 << 9
	LLDPCapsTmpr        uint16 = 1 << 10
)

// LLDPCapabilities represents the capabilities of a device
type LLDPCapabilities struct {
	Other       bool
	Repeater    bool
	Bridge      bool
	WLANAP      bool
	Router      bool
	Phone       bool
	DocSis      bool
	StationOnly bool
	CVLAN       bool
	SVLAN       bool
	TMPR        bool
}

type LLDPSysCapabilities struct {
	SystemCap  LLDPCapabilities
	EnabledCap LLDPCapabilities
}

type LLDPMgmtAddressSubtype byte

// LLDP Management Address Subtypes
const (
	LLDPMgmtAddressSubtypeIPV4 LLDPMgmtAddressSubtype = 1
	LLDPMgmtAddressSubtypeIPV6 LLDPMgmtAddressSubtype = 2
)

type LLDPInterfaceSubtype byte

// LLDP Interface Subtypes
const (
	LLDPInterfaceSubtypeUnknown LLDPInterfaceSubtype = 1
	LLDPInterfaceSubtypeifIndex LLDPInterfaceSubtype = 2
	LLDPInterfaceSubtypeSysPort LLDPInterfaceSubtype = 3
)

type LLDPMgmtAddress struct {
	Subtype          LLDPMgmtAddressSubtype
	Address          []byte
	InterfaceSubtype LLDPInterfaceSubtype
	InterfaceNumber  uint32
	OID              string
}

// LinkLayerDiscoveryInfo represents the decoded details for a set of LinkLayerDiscoveryValues
type LinkLayerDiscoveryInfo struct {
	baseLayer
	PortDescription string
	SysName         string
	SysDescription  string
	SysCapabilities LLDPSysCapabilities
	MgmtAddress     LLDPMgmtAddress
	OrgTLVs []LLDPOrgSpecificTLV      // Private TLVs
	Unknown []LinkLayerDiscoveryValue // undecoded TLVs
}

type IEEEOUI uint32

// http://standards.ieee.org/develop/regauth/oui/oui.txt
const (
	IEEEOUI8021     IEEEOUI = 0x0080c2
	IEEEOUI8023     IEEEOUI = 0x00120f
	IEEEOUI8021Qbg  IEEEOUI = 0x0013BF
	IEEEOUICisco2   IEEEOUI = 0x000142
	IEEEOUITR41     IEEEOUI = 0x0012bb
	IEEEOUIProfinet IEEEOUI = 0x000ecf
)

/// IEEE 802.1 TLV Subtypes
const (
	LLDP8021SubtypePortVLANID       uint8 = 1
	LLDP8021SubtypeProtocolVLANID   uint8 = 2
	LLDP8021SubtypeVLANName         uint8 = 3
	LLDP8021SubtypeProtocolIdentity uint8 = 4
	LLDP8021SubtypeVDIUsageDigest   uint8 = 5
	LLDP8021SubtypeManagementVID    uint8 = 6
	LLDP8021SubtypeLinkAggregation  uint8 = 7
)

// VLAN Port Protocol ID options
const (
	LLDPProtocolVLANIDCapability byte = 1 << 0
	LLDPProtocolVLANIDStatus     byte = 1 << 1
)

type PortProtocolVLANID struct {
	Supported bool
	Enabled   bool
	ID        uint16
}

type VLANName struct {
	ID   uint16
	Name string
}

type ProtocolIdentity []byte

// LACP options
const (
	LLDPAggregationCapability byte = 1 << 0
	LLDPAggregationStatus     byte = 1 << 1
)

// IEEE 802.1 Link Aggregation parameters
type LinkAggregation8021 struct {
	Supported bool
	Enabled   bool
	PortID    uint32
}

type LLDPInfo8021 struct {
	PVID               uint16
	PPVIDs             []PortProtocolVLANID
	VLANNames          []VLANName
	ProtocolIdentities []ProtocolIdentity
	VIDUsageDigest     uint32
	ManagementVID      uint16
	LinkAggregation LinkAggregation8021
}

// IEEE 802.3 TLV Subtypes
const (
	LLDP8023SubtypeMACPHY          uint8 = 1
	LLDP8023SubtypeMDIPower        uint8 = 2
	LLDP8023SubtypeLinkAggregation uint8 = 3
	LLDP8023SubtypeMTU             uint8 = 4
)

// MACPHY options
const (
	LLDPMACPHYCapability byte = 1 << 0
	LLDPMACPHYStatus     byte = 1 << 1
)

// From IANA-MAU-MIB (introduced by RFC 4836) - dot3MauType
const (
	LLDPMAUTypeUnknown         uint16 = 0
	LLDPMAUTypeAUI             uint16 = 1
	LLDPMAUType10Base5         uint16 = 2
	LLDPMAUTypeFOIRL           uint16 = 3
	LLDPMAUType10Base2         uint16 = 4
	LLDPMAUType10BaseT         uint16 = 5
	LLDPMAUType10BaseFP        uint16 = 6
	LLDPMAUType10BaseFB        uint16 = 7
	LLDPMAUType10BaseFL        uint16 = 8
	LLDPMAUType10BROAD36       uint16 = 9
	LLDPMAUType10BaseT_HD      uint16 = 10
	LLDPMAUType10BaseT_FD      uint16 = 11
	LLDPMAUType10BaseFL_HD     uint16 = 12
	LLDPMAUType10BaseFL_FD     uint16 = 13
	LLDPMAUType100BaseT4       uint16 = 14
	LLDPMAUType100BaseTX_HD    uint16 = 15
	LLDPMAUType100BaseTX_FD    uint16 = 16
	LLDPMAUType100BaseFX_HD    uint16 = 17
	LLDPMAUType100BaseFX_FD    uint16 = 18
	LLDPMAUType100BaseT2_HD    uint16 = 19
	LLDPMAUType100BaseT2_FD    uint16 = 20
	LLDPMAUType1000BaseX_HD    uint16 = 21
	LLDPMAUType1000BaseX_FD    uint16 = 22
	LLDPMAUType1000BaseLX_HD   uint16 = 23
	LLDPMAUType1000BaseLX_FD   uint16 = 24
	LLDPMAUType1000BaseSX_HD   uint16 = 25
	LLDPMAUType1000BaseSX_FD   uint16 = 26
	LLDPMAUType1000BaseCX_HD   uint16 = 27
	LLDPMAUType1000BaseCX_FD   uint16 = 28
	LLDPMAUType1000BaseT_HD    uint16 = 29
	LLDPMAUType1000BaseT_FD    uint16 = 30
	LLDPMAUType10GBaseX        uint16 = 31
	LLDPMAUType10GBaseLX4      uint16 = 32
	LLDPMAUType10GBaseR        uint16 = 33
	LLDPMAUType10GBaseER       uint16 = 34
	LLDPMAUType10GBaseLR       uint16 = 35
	LLDPMAUType10GBaseSR       uint16 = 36
	LLDPMAUType10GBaseW        uint16 = 37
	LLDPMAUType10GBaseEW       uint16 = 38
	LLDPMAUType10GBaseLW       uint16 = 39
	LLDPMAUType10GBaseSW       uint16 = 40
	LLDPMAUType10GBaseCX4      uint16 = 41
	LLDPMAUType2BaseTL         uint16 = 42
	LLDPMAUType10PASS_TS       uint16 = 43
	LLDPMAUType100BaseBX10D    uint16 = 44
	LLDPMAUType100BaseBX10U    uint16 = 45
	LLDPMAUType100BaseLX10     uint16 = 46
	LLDPMAUType1000BaseBX10D   uint16 = 47
	LLDPMAUType1000BaseBX10U   uint16 = 48
	LLDPMAUType1000BaseLX10    uint16 = 49
	LLDPMAUType1000BasePX10D   uint16 = 50
	LLDPMAUType1000BasePX10U   uint16 = 51
	LLDPMAUType1000BasePX20D   uint16 = 52
	LLDPMAUType1000BasePX20U   uint16 = 53
	LLDPMAUType10GBaseT        uint16 = 54
	LLDPMAUType10GBaseLRM      uint16 = 55
	LLDPMAUType1000BaseKX      uint16 = 56
	LLDPMAUType10GBaseKX4      uint16 = 57
	LLDPMAUType10GBaseKR       uint16 = 58
	LLDPMAUType10_1GBasePRX_D1 uint16 = 59
	LLDPMAUType10_1GBasePRX_D2 uint16 = 60
	LLDPMAUType10_1GBasePRX_D3 uint16 = 61
	LLDPMAUType10_1GBasePRX_U1 uint16 = 62
	LLDPMAUType10_1GBasePRX_U2 uint16 = 63
	LLDPMAUType10_1GBasePRX_U3 uint16 = 64
	LLDPMAUType10GBasePR_D1    uint16 = 65
	LLDPMAUType10GBasePR_D2    uint16 = 66
	LLDPMAUType10GBasePR_D3    uint16 = 67
	LLDPMAUType10GBasePR_U1    uint16 = 68
	LLDPMAUType10GBasePR_U3    uint16 = 69
)

// From RFC 3636 - ifMauAutoNegCapAdvertisedBits
const (
	LLDPMAUPMDOther        uint16 = 1 << 15
	LLDPMAUPMD10BaseT      uint16 = 1 << 14
	LLDPMAUPMD10BaseT_FD   uint16 = 1 << 13
	LLDPMAUPMD100BaseT4    uint16 = 1 << 12
	LLDPMAUPMD100BaseTX    uint16 = 1 << 11
	LLDPMAUPMD100BaseTX_FD uint16 = 1 << 10
	LLDPMAUPMD100BaseT2    uint16 = 1 << 9
	LLDPMAUPMD100BaseT2_FD uint16 = 1 << 8
	LLDPMAUPMDFDXPAUSE     uint16 = 1 << 7
	LLDPMAUPMDFDXAPAUSE    uint16 = 1 << 6
	LLDPMAUPMDFDXSPAUSE    uint16 = 1 << 5
	LLDPMAUPMDFDXBPAUSE    uint16 = 1 << 4
	LLDPMAUPMD1000BaseX    uint16 = 1 << 3
	LLDPMAUPMD1000BaseX_FD uint16 = 1 << 2
	LLDPMAUPMD1000BaseT    uint16 = 1 << 1
	LLDPMAUPMD1000BaseT_FD uint16 = 1 << 0
)

// Inverted ifMauAutoNegCapAdvertisedBits if required
// (Some manufacturers misinterpreted the spec - 
// see https://bugs.wireshark.org/bugzilla/show_bug.cgi?id=1455)
const (
	LLDPMAUPMDOtherInv        uint16 = 1 << 0
	LLDPMAUPMD10BaseTInv      uint16 = 1 << 1
	LLDPMAUPMD10BaseT_FDInv   uint16 = 1 << 2
	LLDPMAUPMD100BaseT4Inv    uint16 = 1 << 3
	LLDPMAUPMD100BaseTXInv    uint16 = 1 << 4
	LLDPMAUPMD100BaseTX_FDInv uint16 = 1 << 5
	LLDPMAUPMD100BaseT2Inv    uint16 = 1 << 6
	LLDPMAUPMD100BaseT2_FDInv uint16 = 1 << 7
	LLDPMAUPMDFDXPAUSEInv     uint16 = 1 << 8
	LLDPMAUPMDFDXAPAUSEInv    uint16 = 1 << 9
	LLDPMAUPMDFDXSPAUSEInv    uint16 = 1 << 10
	LLDPMAUPMDFDXBPAUSEInv    uint16 = 1 << 11
	LLDPMAUPMD1000BaseXInv    uint16 = 1 << 12
	LLDPMAUPMD1000BaseX_FDInv uint16 = 1 << 13
	LLDPMAUPMD1000BaseTInv    uint16 = 1 << 14
	LLDPMAUPMD1000BaseT_FDInv uint16 = 1 << 15
)

type MACPHYConfigStatus struct {
	AutoNegSupported  bool
	AutoNegEnabled    bool
	AutoNegCapability uint16
	MAUType           uint16
}

// MDI Power options
const (
	LLDPMDIPowerPortClass    byte = 1 << 0
	LLDPMDIPowerCapability   byte = 1 << 1
	LLDPMDIPowerStatus       byte = 1 << 2
	LLDPMDIPowerPairsAbility byte = 1 << 3
)

type LLDPPowerType byte

type LLDPPowerSource byte

type LLDPPowerPriority byte

type PowerViaMDI struct {
	PortClassPSE    bool // false = PD
	PSESupported    bool
	PSEEnabled      bool
	PSEPairsAbility bool
	PSEPowerPair    uint8
	PSEClass        uint8
	PowerType       LLDPPowerType
	PowerSource     LLDPPowerSource
	PowerPriority   LLDPPowerPriority
	RequestedPower  uint16 // 1-510 Watts
	AllocatedPower  uint16 // 1-510 Watts
}

// IEEE 802.3 Link Aggregation parameters
type LinkAggregation8023 struct {
	Status byte
	PortID uint32
}

type LLDPInfo8023 struct {
	MACPHYConfigStatus
	PowerViaMDI
	LinkAggregation LinkAggregation8023
	MTU uint16
}

// IEEE 802.1Qbg TLV Subtypes
const (
	LLDP8021QbgEVB  uint8 = 0
	LLDP8021QbgCDCP uint8 = 1
	LLDP8021QbgVDP  uint8 = 2
)


// LLDPEVBCapabilities Types
const (
	LLDPEVBCapsSTD uint16 = 1 << 0
	LLDPEVBCapsRR  uint16 = 1 << 1
	LLDPEVBCapsRTE uint16 = 1 << 2
	LLDPEVBCapsECP uint16 = 1 << 3
	LLDPEVBCapsVDP uint16 = 1 << 4
)

// LLDPEVBCapabilities represents the EVB capabilities of a device
type LLDPEVBCapabilities struct {
	StandardBridging            bool
	ReflectiveRelay             bool
	RetransmissionTimerExponent bool
	EdgeControlProtocol         bool
	VSIDiscoveryProtocol        bool
}

type LLDPEVBSettings struct {
	Supported      LLDPEVBCapabilities
	Enabled        LLDPEVBCapabilities
	SupportedVSIs  uint16
	ConfiguredVSIs uint16
	RTEExponent    uint8
}

type LLDPInfo8021Qbg struct {
	EVBSettings LLDPEVBSettings
}


// LayerType returns gopacket.LayerTypeLinkLayerDiscovery.
func (c *LinkLayerDiscovery) LayerType() gopacket.LayerType {
	return LayerTypeLinkLayerDiscovery
}

func decodeLinkLayerDiscovery(data []byte, p gopacket.PacketBuilder) error {
	var vals []LinkLayerDiscoveryValue
	vData := data[0:]
	for len(vData) > 0 {
		nbit := vData[0] & 0x01
		t := LLDPTLVType(vData[0] >> 1)
		val := LinkLayerDiscoveryValue{Type: t, Length: uint16(nbit<<8 + vData[1])}
		if val.Length > 0 {
			val.Value = vData[2 : val.Length+2]
		}
		vals = append(vals, val)
		if t == LLDPTLVEnd {
			break
		}
		if len(vData) < int(2+val.Length) {
			return fmt.Errorf("Malformed LinkLayerDiscovery Header")
		}
		vData = vData[2+val.Length:]
	}
	if len(vals) < 4 {
		return fmt.Errorf("Missing mandatory LinkLayerDiscovery TLV")
	}
	c := &LinkLayerDiscovery{}
	gotEnd := false
	for _, v := range vals {
		switch v.Type {
		case LLDPTLVEnd:
			gotEnd = true
		case LLDPTLVChassisID:
			if len(v.Value) < 2 {
				return fmt.Errorf("Malformed LinkLayerDiscovery ChassisID TLV")
			}
			c.ChassisID.Subtype = LLDPChassisIDSubType(v.Value[0])
			c.ChassisID.ID = v.Value[1:]
		case LLDPTLVPortID:
			if len(v.Value) < 2 {
				return fmt.Errorf("Malformed LinkLayerDiscovery PortID TLV")
			}
			c.PortID.Subtype = LLDPPortIDSubType(v.Value[0])
			c.PortID.ID = v.Value[1:]
		case LLDPTLVTTL:
			if len(v.Value) < 2 {
				return fmt.Errorf("Malformed LinkLayerDiscovery TTL TLV")
			}
			c.TTL = binary.BigEndian.Uint16(v.Value[0:2])
		default:
			c.Values = append(c.Values, v)
		}
	}
	if c.ChassisID.Subtype == 0 || c.PortID.Subtype == 0 || !gotEnd {
		return fmt.Errorf("Missing mandatory LinkLayerDiscovery TLV")
	}
	c.contents = data
	p.AddLayer(c)

	info := &LinkLayerDiscoveryInfo{}
	var errors []error
	var ok bool
	for _, v := range c.Values {
		switch v.Type {
		case LLDPTLVPortDescription:
			info.PortDescription = string(v.Value)
		case LLDPTLVSysName:
			info.SysName = string(v.Value)
		case LLDPTLVSysDescription:
			info.SysDescription = string(v.Value)
		case LLDPTLVSysCapabilities:
			if ok, errors = checkLLDPTLVLen(v, 4, errors); ok {
				info.SysCapabilities.SystemCap = getCapabilities(binary.BigEndian.Uint16(v.Value[0:2]))
				info.SysCapabilities.EnabledCap = getCapabilities(binary.BigEndian.Uint16(v.Value[2:4]))
			}
		case LLDPTLVMgmtAddress:
			if ok, errors = checkLLDPTLVLen(v, 9, errors); ok {
				mlen := v.Value[0]
				if ok, errors = checkLLDPTLVLen(v, int(mlen+7), errors); !ok {
					continue
				}
				info.MgmtAddress.Subtype = LLDPMgmtAddressSubtype(v.Value[1])
				info.MgmtAddress.Address = v.Value[2 : mlen+1]
				info.MgmtAddress.InterfaceSubtype = LLDPInterfaceSubtype(v.Value[mlen+1])
				info.MgmtAddress.InterfaceNumber = binary.BigEndian.Uint32(v.Value[mlen+2 : mlen+6])
				olen := v.Value[mlen+6]
				if ok, errors = checkLLDPTLVLen(v, int(mlen+6+olen), errors); ok {
					info.MgmtAddress.OID = string(v.Value[mlen+9 : mlen+9+olen])
				}
			}
		case LLDPTLVOrgSpecific:
			if ok, errors = checkLLDPTLVLen(v, 4, errors); !ok {
				continue
			}
			info.OrgTLVs = append(info.OrgTLVs, LLDPOrgSpecificTLV{IEEEOUI(binary.BigEndian.Uint32(append([]byte{byte(0)}, v.Value[0:3]...))), uint8(v.Value[3]), v.Value[4:]})
		}
	}
	p.AddLayer(info)
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

func (l *LinkLayerDiscoveryInfo) Decode8021() (info LLDPInfo8021, err error) {
	var errors []error
	var ok bool
	for _, o := range l.OrgTLVs {
		if o.OUI != IEEEOUI8021 {
			continue;
		}
		switch o.SubType {
		case LLDP8021SubtypePortVLANID:
			if ok, errors = checkLLDPOrgSpecificLen(o, 2, errors); ok {
				info.PVID = binary.BigEndian.Uint16(o.Info[0:2])
			}
		case LLDP8021SubtypeProtocolVLANID:
			if ok, errors = checkLLDPOrgSpecificLen(o, 3, errors); ok {
				sup := (o.Info[0]&LLDPProtocolVLANIDCapability > 0)
				en := (o.Info[0]&LLDPAggregationStatus > 0)
				id := binary.BigEndian.Uint16(o.Info[1:3])
				info.PPVIDs = append(info.PPVIDs, PortProtocolVLANID{sup, en, id})
			}
		case LLDP8021SubtypeVLANName:
			if ok, errors = checkLLDPOrgSpecificLen(o, 2, errors); ok {
				id := binary.BigEndian.Uint16(o.Info[0:2])
				info.VLANNames = append(info.VLANNames, VLANName{id, string(o.Info[3:])})
			}
		case LLDP8021SubtypeProtocolIdentity:
			if ok, errors = checkLLDPOrgSpecificLen(o, 1, errors); ok {
				l := int(o.Info[0])
				if l > 0 {
					info.ProtocolIdentities = append(info.ProtocolIdentities, o.Info[1:1+l])
				}
			}
		case LLDP8021SubtypeVDIUsageDigest:
			if ok, errors = checkLLDPOrgSpecificLen(o, 4, errors); ok {
				info.VIDUsageDigest = binary.BigEndian.Uint32(o.Info[0:4])
			}
		case LLDP8021SubtypeManagementVID:
			if ok, errors = checkLLDPOrgSpecificLen(o, 2, errors); ok {
				info.ManagementVID = binary.BigEndian.Uint16(o.Info[0:2])
			}
		case LLDP8021SubtypeLinkAggregation:
			if ok, errors = checkLLDPOrgSpecificLen(o, 5, errors); ok {
				sup := (o.Info[0]&LLDPAggregationCapability > 0)
				en := (o.Info[0]&LLDPAggregationStatus > 0)
				info.LinkAggregation = LinkAggregation8021{sup, en, binary.BigEndian.Uint32(o.Info[1:5])}
			}
		}
	}
	if len(errors) > 0 {
		err = errors[0]
	}
	return
}

func (l *LinkLayerDiscoveryInfo) Decode8023() (info LLDPInfo8023, err error) {
	var errors []error
	var ok bool
	for _, o := range l.OrgTLVs {
		if o.OUI != IEEEOUI8023 {
			continue;
		}
		switch o.SubType {
		case LLDP8023SubtypeMACPHY:
			if ok, errors = checkLLDPOrgSpecificLen(o, 5, errors); ok {
				sup := (o.Info[0]&LLDPMACPHYCapability > 0)
				en := (o.Info[0]&LLDPMACPHYStatus > 0)
				ca := binary.BigEndian.Uint16(o.Info[1:3])
				mau := binary.BigEndian.Uint16(o.Info[3:5])
				info.MACPHYConfigStatus = MACPHYConfigStatus{sup, en, ca, mau}
			}
		case LLDP8023SubtypeMDIPower:
			if ok, errors = checkLLDPOrgSpecificLen(o, 3, errors); ok {
				info.PowerViaMDI.PortClassPSE = (o.Info[0]&LLDPMDIPowerPortClass > 0)
				info.PowerViaMDI.PSESupported = (o.Info[0]&LLDPMDIPowerCapability > 0)
				info.PowerViaMDI.PSEEnabled = (o.Info[0]&LLDPMDIPowerStatus > 0)
				info.PowerViaMDI.PSEPairsAbility = (o.Info[0]&LLDPMDIPowerPairsAbility > 0)
				info.PowerViaMDI.PSEPowerPair = uint8(o.Info[1])
				info.PowerViaMDI.PSEClass = uint8(o.Info[2])
				if len(o.Info) >= 8 {
					info.PowerViaMDI.PowerType = LLDPPowerType((o.Info[3] & 0xc0) >> 6)
					info.PowerViaMDI.PowerSource = LLDPPowerSource((o.Info[3] & 0x30) >> 4)
					if info.PowerViaMDI.PowerType == 1 || info.PowerViaMDI.PowerType == 3 {
						info.PowerViaMDI.PowerSource += 128 // For Stringify purposes
					}
					info.PowerViaMDI.PowerPriority = LLDPPowerPriority(o.Info[4] & 0x0f)
					info.PowerViaMDI.RequestedPower = binary.BigEndian.Uint16(o.Info[5:7])
					info.PowerViaMDI.AllocatedPower = binary.BigEndian.Uint16(o.Info[7:8])
				}
			}
		case LLDP8023SubtypeLinkAggregation:
			if ok, errors = checkLLDPOrgSpecificLen(o, 5, errors); ok {
				info.LinkAggregation = LinkAggregation8023{o.Info[0], binary.BigEndian.Uint32(o.Info[1:5])}
			}
		case LLDP8023SubtypeMTU:
			if ok, errors = checkLLDPOrgSpecificLen(o, 2, errors); ok {
				info.MTU = binary.BigEndian.Uint16(o.Info[0:2])
			}
		}
	}
	if len(errors) > 0 {
		err = errors[0]
	}
	return
}

func (l *LinkLayerDiscoveryInfo) Decode8021Qbg() (info LLDPInfo8021Qbg, err error) {
	var errors []error
	var ok bool
	for _, o := range l.OrgTLVs {
		if o.OUI != IEEEOUI8021Qbg {
			continue;
		}
		switch o.SubType {
		case LLDP8021QbgEVB:
			if ok, errors = checkLLDPOrgSpecificLen(o, 9, errors); ok {
				info.EVBSettings.Supported = getEVBCapabilities(binary.BigEndian.Uint16(o.Info[0:2]))
				info.EVBSettings.Enabled = getEVBCapabilities(binary.BigEndian.Uint16(o.Info[2:4]))
				info.EVBSettings.SupportedVSIs = binary.BigEndian.Uint16(o.Info[4:6])
				info.EVBSettings.ConfiguredVSIs = binary.BigEndian.Uint16(o.Info[6:8])
				info.EVBSettings.RTEExponent = uint8(o.Info[8])
			}
		}
	}
	if len(errors) > 0 {
		err = errors[0]
	}
	return
}

// LayerType returns gopacket.LayerTypeLinkLayerDiscoveryInfo.
func (c *LinkLayerDiscoveryInfo) LayerType() gopacket.LayerType {
	return LayerTypeLinkLayerDiscoveryInfo
}


func getCapabilities(v uint16) (c LLDPCapabilities) {
	c.Other = (v&LLDPCapsOther > 0)
	c.Repeater = (v&LLDPCapsRepeater > 0)
	c.Bridge = (v&LLDPCapsBridge > 0)
	c.WLANAP = (v&LLDPCapsWLANAP > 0)
	c.Router = (v&LLDPCapsRouter > 0)
	c.Phone = (v&LLDPCapsPhone > 0)
	c.DocSis = (v&LLDPCapsDocSis > 0)
	c.StationOnly = (v&LLDPCapsStationOnly > 0)
	c.CVLAN = (v&LLDPCapsCVLAN > 0)
	c.SVLAN = (v&LLDPCapsSVLAN > 0)
	c.TMPR = (v&LLDPCapsTmpr > 0)
	return
}

func getEVBCapabilities(v uint16) (c LLDPEVBCapabilities) {
	c.StandardBridging = (v & LLDPEVBCapsSTD) > 0
	c.StandardBridging = (v & LLDPEVBCapsSTD) > 0
	c.ReflectiveRelay = (v & LLDPEVBCapsRR) > 0
	c.RetransmissionTimerExponent = (v & LLDPEVBCapsRTE) > 0
	c.EdgeControlProtocol = (v & LLDPEVBCapsECP) > 0
	c.VSIDiscoveryProtocol = (v & LLDPEVBCapsVDP) > 0
	return
}

func (t LLDPTLVType) String() (s string) {
	switch t {
	case LLDPTLVEnd:
		s = "TLV End"
	case LLDPTLVChassisID:
		s = "Chassis ID"
	case LLDPTLVPortID:
		s = "Port ID"
	case LLDPTLVTTL:
		s = "TTL"
	case LLDPTLVPortDescription:
		s = "Port Description"
	case LLDPTLVSysName:
		s = "System Name"
	case LLDPTLVSysDescription:
		s = "System Description"
	case LLDPTLVSysCapabilities:
		s = "System Capabilities"
	case LLDPTLVMgmtAddress:
		s = "Management Address"
	case LLDPTLVOrgSpecific:
		s = "Organisation Specific"
	default:
		s = "Unknown"
	}
	return
}

func (t LLDPChassisIDSubType) String() (s string) {
	switch t {
	case LLDPChassisIDSubTypeReserved:
		s = "Reserved"
	case LLDPChassisIDSubTypeChassisComp:
		s = "Chassis Component"
	case LLDPChassisIDSubtypeIfaceAlias:
		s = "Interface Alias"
	case LLDPChassisIDSubTypePortComp:
		s = "Port Component"
	case LLDPChassisIDSubTypeMACAddr:
		s = "MAC Address"
	case LLDPChassisIDSubTypeNetworkAddr:
		s = "Network Address"
	case LLDPChassisIDSubtypeIfaceName:
		s = "Interface Name"
	case LLDPChassisIDSubTypeLocal:
		s = "Local"
	default:
		s = "Unknown"
	}
	return
}

func (t LLDPPortIDSubType) String() (s string) {
	switch t {
	case LLDPPortIDSubtypeReserved:
		s = "Reserved"
	case LLDPPortIDSubtypeIfaceAlias:
		s = "Interface Alias"
	case LLDPPortIDSubTypePortComp:
		s = "Port Component"
	case LLDPPortIDSubTypeMACAddr:
		s = "MAC Address"
	case LLDPPortIDSubTypeNetworkAddr:
		s = "Network Address"
	case LLDPPortIDSubtypeIfaceName:
		s = "Interface Name"
	case LLDPPortIDSubTypeAgentCircuitID:
		s = "Agent Circuit ID"
	case LLDPPortIDSubTypeLocal:
		s = "Local"
	default:
		s = "Unknown"
	}
	return
}

func (t LLDPMgmtAddressSubtype) String() (s string) {
	switch t {
	case LLDPMgmtAddressSubtypeIPV4:
		s = "IPv4"
	case LLDPMgmtAddressSubtypeIPV6:
		s = "IPv6"
	default:
		s = "Unknown"
	}
	return
}

func (t LLDPInterfaceSubtype) String() (s string) {
	switch t {
	case LLDPInterfaceSubtypeUnknown:
		s = "Unknown"
	case LLDPInterfaceSubtypeifIndex:
		s = "IfIndex"
	case LLDPInterfaceSubtypeSysPort:
		s = "System Port Number"
	default:
		s = "Unknown"
	}
	return
}

func (t LLDPPowerType) String() (s string) {
	switch t {
	case 0:
		s = "Type 2 PSE Device"
	case 1:
		s = "Type 2 PD Device"
	case 2:
		s = "Type 1 PSE Device"
	case 3:
		s = "Type 1 PD Device"
	default:
		s = "Unknown"
	}
	return
}

func (t LLDPPowerSource) String() (s string) {
	switch t {
	// PD Device
	case 0:
		s = "Unknown"
	case 1:
		s = "PSE"
	case 2:
		s = "Local"
	case 3:
		s = "PSE and Local"
	// PSE Device  (Actual value  + 128)
	case 128:
		s = "Unknown"
	case 129:
		s = "Primary Power Source"
	case 130:
		s = "Backup Power Source"
	default:
		s = "Unknown"
	}
	return
}

func (t LLDPPowerPriority) String() (s string) {
	switch t {
	case 0:
		s = "Unknown"
	case 1:
		s = "Critical"
	case 2:
		s = "High"
	case 3:
		s = "Low"
	default:
		s = "Unknown"
	}
	return
}

func checkLLDPTLVLen(v LinkLayerDiscoveryValue, l int, e []error) (ok bool, errors []error) {
	errors = e
	if ok = (len(v.Value) >= l); !ok {
		errors = append(errors, fmt.Errorf("Invalid TLV %v length %d (wanted mimimum %v", v.Type, len(v.Value), l))
	}
	return
}

func checkLLDPOrgSpecificLen(o LLDPOrgSpecificTLV, l int, e []error) (ok bool, errors []error) {
	errors = e
	if ok = (len(o.Info) >= l); !ok {
		errors = append(errors, fmt.Errorf("Invalid Org Specific TLV %v length %d (wanted minimum %v)", o.SubType, len(o.Info), l))
	}
	return
}