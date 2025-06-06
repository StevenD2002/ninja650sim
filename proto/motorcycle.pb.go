// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.29.3
// source: proto/motorcycle.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Basic engine data
type EngineData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rpm              float64 `protobuf:"fixed64,1,opt,name=rpm,proto3" json:"rpm,omitempty"`
	ThrottlePosition float64 `protobuf:"fixed64,2,opt,name=throttle_position,json=throttlePosition,proto3" json:"throttle_position,omitempty"` // 0-100%
	Timestamp        int64   `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Power            float64 `protobuf:"fixed64,4,opt,name=power,proto3" json:"power,omitempty"`                                              // Horsepower
	Torque           float64 `protobuf:"fixed64,5,opt,name=torque,proto3" json:"torque,omitempty"`                                            // Nm
	EngineTemp       float64 `protobuf:"fixed64,6,opt,name=engine_temp,json=engineTemp,proto3" json:"engine_temp,omitempty"`                  // Celsius
	AfrCurrent       float64 `protobuf:"fixed64,7,opt,name=afr_current,json=afrCurrent,proto3" json:"afr_current,omitempty"`                  // Current Air/Fuel Ratio
	AfrTarget        float64 `protobuf:"fixed64,8,opt,name=afr_target,json=afrTarget,proto3" json:"afr_target,omitempty"`                     // Target Air/Fuel Ratio
	FuelInjectionMs  float64 `protobuf:"fixed64,9,opt,name=fuel_injection_ms,json=fuelInjectionMs,proto3" json:"fuel_injection_ms,omitempty"` // Fuel injection duration in ms
	IgnitionAdvance  float64 `protobuf:"fixed64,10,opt,name=ignition_advance,json=ignitionAdvance,proto3" json:"ignition_advance,omitempty"`  // Ignition timing in degrees BTDC
	Gear             int32   `protobuf:"varint,11,opt,name=gear,proto3" json:"gear,omitempty"`
	Speed            float64 `protobuf:"fixed64,12,opt,name=speed,proto3" json:"speed,omitempty"`                                         // km/h
	ClutchPosition   float64 `protobuf:"fixed64,13,opt,name=clutch_position,json=clutchPosition,proto3" json:"clutch_position,omitempty"` // 0-1
}

func (x *EngineData) Reset() {
	*x = EngineData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_motorcycle_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EngineData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EngineData) ProtoMessage() {}

func (x *EngineData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_motorcycle_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EngineData.ProtoReflect.Descriptor instead.
func (*EngineData) Descriptor() ([]byte, []int) {
	return file_proto_motorcycle_proto_rawDescGZIP(), []int{0}
}

func (x *EngineData) GetRpm() float64 {
	if x != nil {
		return x.Rpm
	}
	return 0
}

func (x *EngineData) GetThrottlePosition() float64 {
	if x != nil {
		return x.ThrottlePosition
	}
	return 0
}

func (x *EngineData) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *EngineData) GetPower() float64 {
	if x != nil {
		return x.Power
	}
	return 0
}

func (x *EngineData) GetTorque() float64 {
	if x != nil {
		return x.Torque
	}
	return 0
}

func (x *EngineData) GetEngineTemp() float64 {
	if x != nil {
		return x.EngineTemp
	}
	return 0
}

func (x *EngineData) GetAfrCurrent() float64 {
	if x != nil {
		return x.AfrCurrent
	}
	return 0
}

func (x *EngineData) GetAfrTarget() float64 {
	if x != nil {
		return x.AfrTarget
	}
	return 0
}

func (x *EngineData) GetFuelInjectionMs() float64 {
	if x != nil {
		return x.FuelInjectionMs
	}
	return 0
}

func (x *EngineData) GetIgnitionAdvance() float64 {
	if x != nil {
		return x.IgnitionAdvance
	}
	return 0
}

func (x *EngineData) GetGear() int32 {
	if x != nil {
		return x.Gear
	}
	return 0
}

func (x *EngineData) GetSpeed() float64 {
	if x != nil {
		return x.Speed
	}
	return 0
}

func (x *EngineData) GetClutchPosition() float64 {
	if x != nil {
		return x.ClutchPosition
	}
	return 0
}

// User input
type UserInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ThrottlePosition float64 `protobuf:"fixed64,1,opt,name=throttle_position,json=throttlePosition,proto3" json:"throttle_position,omitempty"` // 0-100%
	ClutchPosition   float64 `protobuf:"fixed64,2,opt,name=clutch_position,json=clutchPosition,proto3" json:"clutch_position,omitempty"`       // 0-1 (0=engaged, 1=disengaged)
	Gear             int32   `protobuf:"varint,3,opt,name=gear,proto3" json:"gear,omitempty"`                                                  // 0=Neutral, 1-6=Gears
}

func (x *UserInput) Reset() {
	*x = UserInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_motorcycle_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInput) ProtoMessage() {}

func (x *UserInput) ProtoReflect() protoreflect.Message {
	mi := &file_proto_motorcycle_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInput.ProtoReflect.Descriptor instead.
func (*UserInput) Descriptor() ([]byte, []int) {
	return file_proto_motorcycle_proto_rawDescGZIP(), []int{1}
}

func (x *UserInput) GetThrottlePosition() float64 {
	if x != nil {
		return x.ThrottlePosition
	}
	return 0
}

func (x *UserInput) GetClutchPosition() float64 {
	if x != nil {
		return x.ClutchPosition
	}
	return 0
}

func (x *UserInput) GetGear() int32 {
	if x != nil {
		return x.Gear
	}
	return 0
}

// A single row in a 2D map
type MapRow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []float64 `protobuf:"fixed64,1,rep,packed,name=values,proto3" json:"values,omitempty"`
}

func (x *MapRow) Reset() {
	*x = MapRow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_motorcycle_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapRow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapRow) ProtoMessage() {}

func (x *MapRow) ProtoReflect() protoreflect.Message {
	mi := &file_proto_motorcycle_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapRow.ProtoReflect.Descriptor instead.
func (*MapRow) Descriptor() ([]byte, []int) {
	return file_proto_motorcycle_proto_rawDescGZIP(), []int{2}
}

func (x *MapRow) GetValues() []float64 {
	if x != nil {
		return x.Values
	}
	return nil
}

// A 2D map (e.g., fuel, ignition)
type Map2D struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type            string    `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"` // "fuel", "ignition", or "afr"
	RpmBreakpoints  []float64 `protobuf:"fixed64,2,rep,packed,name=rpm_breakpoints,json=rpmBreakpoints,proto3" json:"rpm_breakpoints,omitempty"`
	LoadBreakpoints []float64 `protobuf:"fixed64,3,rep,packed,name=load_breakpoints,json=loadBreakpoints,proto3" json:"load_breakpoints,omitempty"`
	Values          []*MapRow `protobuf:"bytes,4,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *Map2D) Reset() {
	*x = Map2D{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_motorcycle_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Map2D) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Map2D) ProtoMessage() {}

func (x *Map2D) ProtoReflect() protoreflect.Message {
	mi := &file_proto_motorcycle_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Map2D.ProtoReflect.Descriptor instead.
func (*Map2D) Descriptor() ([]byte, []int) {
	return file_proto_motorcycle_proto_rawDescGZIP(), []int{3}
}

func (x *Map2D) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Map2D) GetRpmBreakpoints() []float64 {
	if x != nil {
		return x.RpmBreakpoints
	}
	return nil
}

func (x *Map2D) GetLoadBreakpoints() []float64 {
	if x != nil {
		return x.LoadBreakpoints
	}
	return nil
}

func (x *Map2D) GetValues() []*MapRow {
	if x != nil {
		return x.Values
	}
	return nil
}

// All ECU maps
type ECUMaps struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FuelMap     *Map2D `protobuf:"bytes,1,opt,name=fuel_map,json=fuelMap,proto3" json:"fuel_map,omitempty"`
	IgnitionMap *Map2D `protobuf:"bytes,2,opt,name=ignition_map,json=ignitionMap,proto3" json:"ignition_map,omitempty"`
	AfrMap      *Map2D `protobuf:"bytes,3,opt,name=afr_map,json=afrMap,proto3" json:"afr_map,omitempty"`
}

func (x *ECUMaps) Reset() {
	*x = ECUMaps{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_motorcycle_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ECUMaps) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ECUMaps) ProtoMessage() {}

func (x *ECUMaps) ProtoReflect() protoreflect.Message {
	mi := &file_proto_motorcycle_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ECUMaps.ProtoReflect.Descriptor instead.
func (*ECUMaps) Descriptor() ([]byte, []int) {
	return file_proto_motorcycle_proto_rawDescGZIP(), []int{4}
}

func (x *ECUMaps) GetFuelMap() *Map2D {
	if x != nil {
		return x.FuelMap
	}
	return nil
}

func (x *ECUMaps) GetIgnitionMap() *Map2D {
	if x != nil {
		return x.IgnitionMap
	}
	return nil
}

func (x *ECUMaps) GetAfrMap() *Map2D {
	if x != nil {
		return x.AfrMap
	}
	return nil
}

// Request for ECU maps
type MapsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MapsRequest) Reset() {
	*x = MapsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_motorcycle_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapsRequest) ProtoMessage() {}

func (x *MapsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_motorcycle_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapsRequest.ProtoReflect.Descriptor instead.
func (*MapsRequest) Descriptor() ([]byte, []int) {
	return file_proto_motorcycle_proto_rawDescGZIP(), []int{5}
}

// Request to update a map cell
type MapUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MapType string  `protobuf:"bytes,1,opt,name=map_type,json=mapType,proto3" json:"map_type,omitempty"` // "fuel", "ignition", or "afr"
	Rpm     float64 `protobuf:"fixed64,2,opt,name=rpm,proto3" json:"rpm,omitempty"`
	Load    float64 `protobuf:"fixed64,3,opt,name=load,proto3" json:"load,omitempty"`
	Value   float64 `protobuf:"fixed64,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *MapUpdateRequest) Reset() {
	*x = MapUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_motorcycle_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapUpdateRequest) ProtoMessage() {}

func (x *MapUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_motorcycle_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapUpdateRequest.ProtoReflect.Descriptor instead.
func (*MapUpdateRequest) Descriptor() ([]byte, []int) {
	return file_proto_motorcycle_proto_rawDescGZIP(), []int{6}
}

func (x *MapUpdateRequest) GetMapType() string {
	if x != nil {
		return x.MapType
	}
	return ""
}

func (x *MapUpdateRequest) GetRpm() float64 {
	if x != nil {
		return x.Rpm
	}
	return 0
}

func (x *MapUpdateRequest) GetLoad() float64 {
	if x != nil {
		return x.Load
	}
	return 0
}

func (x *MapUpdateRequest) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// ECU settings
type ECUSettings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FuelTrim         float64 `protobuf:"fixed64,1,opt,name=fuel_trim,json=fuelTrim,proto3" json:"fuel_trim,omitempty"`
	IgnitionTrim     float64 `protobuf:"fixed64,2,opt,name=ignition_trim,json=ignitionTrim,proto3" json:"ignition_trim,omitempty"`
	IdleRpm          float64 `protobuf:"fixed64,3,opt,name=idle_rpm,json=idleRpm,proto3" json:"idle_rpm,omitempty"`
	RevLimit         float64 `protobuf:"fixed64,4,opt,name=rev_limit,json=revLimit,proto3" json:"rev_limit,omitempty"`
	TempCompensation bool    `protobuf:"varint,5,opt,name=temp_compensation,json=tempCompensation,proto3" json:"temp_compensation,omitempty"`
}

func (x *ECUSettings) Reset() {
	*x = ECUSettings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_motorcycle_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ECUSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ECUSettings) ProtoMessage() {}

func (x *ECUSettings) ProtoReflect() protoreflect.Message {
	mi := &file_proto_motorcycle_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ECUSettings.ProtoReflect.Descriptor instead.
func (*ECUSettings) Descriptor() ([]byte, []int) {
	return file_proto_motorcycle_proto_rawDescGZIP(), []int{7}
}

func (x *ECUSettings) GetFuelTrim() float64 {
	if x != nil {
		return x.FuelTrim
	}
	return 0
}

func (x *ECUSettings) GetIgnitionTrim() float64 {
	if x != nil {
		return x.IgnitionTrim
	}
	return 0
}

func (x *ECUSettings) GetIdleRpm() float64 {
	if x != nil {
		return x.IdleRpm
	}
	return 0
}

func (x *ECUSettings) GetRevLimit() float64 {
	if x != nil {
		return x.RevLimit
	}
	return 0
}

func (x *ECUSettings) GetTempCompensation() bool {
	if x != nil {
		return x.TempCompensation
	}
	return false
}

// Status response for updates
type UpdateStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *UpdateStatus) Reset() {
	*x = UpdateStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_motorcycle_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateStatus) ProtoMessage() {}

func (x *UpdateStatus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_motorcycle_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateStatus.ProtoReflect.Descriptor instead.
func (*UpdateStatus) Descriptor() ([]byte, []int) {
	return file_proto_motorcycle_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateStatus) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UpdateStatus) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_motorcycle_proto protoreflect.FileDescriptor

var file_proto_motorcycle_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6d, 0x6f, 0x74, 0x6f, 0x72, 0x63,
	0x79, 0x63, 0x6c, 0x65, 0x22, 0xa2, 0x03, 0x0a, 0x0a, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x70, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x03, 0x72, 0x70, 0x6d, 0x12, 0x2b, 0x0a, 0x11, 0x74, 0x68, 0x72, 0x6f, 0x74, 0x74, 0x6c,
	0x65, 0x5f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x10, 0x74, 0x68, 0x72, 0x6f, 0x74, 0x74, 0x6c, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x14, 0x0a, 0x05, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x05, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f, 0x72, 0x71, 0x75, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x74, 0x6f, 0x72, 0x71, 0x75, 0x65, 0x12, 0x1f,
	0x0a, 0x0b, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x0a, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x54, 0x65, 0x6d, 0x70, 0x12,
	0x1f, 0x0a, 0x0b, 0x61, 0x66, 0x72, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x61, 0x66, 0x72, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x66, 0x72, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x61, 0x66, 0x72, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12,
	0x2a, 0x0a, 0x11, 0x66, 0x75, 0x65, 0x6c, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x6d, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x66, 0x75, 0x65, 0x6c,
	0x49, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x69,
	0x67, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x64, 0x76, 0x61, 0x6e, 0x63, 0x65, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x41,
	0x64, 0x76, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x65, 0x61, 0x72, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x67, 0x65, 0x61, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x70,
	0x65, 0x65, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64,
	0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6c, 0x75, 0x74, 0x63, 0x68, 0x5f, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x63, 0x6c, 0x75, 0x74, 0x63,
	0x68, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x75, 0x0a, 0x09, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x2b, 0x0a, 0x11, 0x74, 0x68, 0x72, 0x6f, 0x74, 0x74,
	0x6c, 0x65, 0x5f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x10, 0x74, 0x68, 0x72, 0x6f, 0x74, 0x74, 0x6c, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6c, 0x75, 0x74, 0x63, 0x68, 0x5f, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x63, 0x6c,
	0x75, 0x74, 0x63, 0x68, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x67, 0x65, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x67, 0x65, 0x61, 0x72,
	0x22, 0x20, 0x0a, 0x06, 0x4d, 0x61, 0x70, 0x52, 0x6f, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x01, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x22, 0x9b, 0x01, 0x0a, 0x05, 0x4d, 0x61, 0x70, 0x32, 0x44, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x27, 0x0a, 0x0f, 0x72, 0x70, 0x6d, 0x5f, 0x62, 0x72, 0x65, 0x61, 0x6b, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x01, 0x52, 0x0e, 0x72, 0x70, 0x6d, 0x42, 0x72,
	0x65, 0x61, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x6c, 0x6f, 0x61,
	0x64, 0x5f, 0x62, 0x72, 0x65, 0x61, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x01, 0x52, 0x0f, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x72, 0x65, 0x61, 0x6b, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x73, 0x12, 0x2a, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63, 0x6c,
	0x65, 0x2e, 0x4d, 0x61, 0x70, 0x52, 0x6f, 0x77, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x22, 0x99, 0x01, 0x0a, 0x07, 0x45, 0x43, 0x55, 0x4d, 0x61, 0x70, 0x73, 0x12, 0x2c, 0x0a, 0x08,
	0x66, 0x75, 0x65, 0x6c, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x6d, 0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x4d, 0x61, 0x70, 0x32,
	0x44, 0x52, 0x07, 0x66, 0x75, 0x65, 0x6c, 0x4d, 0x61, 0x70, 0x12, 0x34, 0x0a, 0x0c, 0x69, 0x67,
	0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x6d, 0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x4d, 0x61,
	0x70, 0x32, 0x44, 0x52, 0x0b, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x70,
	0x12, 0x2a, 0x0a, 0x07, 0x61, 0x66, 0x72, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x4d,
	0x61, 0x70, 0x32, 0x44, 0x52, 0x06, 0x61, 0x66, 0x72, 0x4d, 0x61, 0x70, 0x22, 0x0d, 0x0a, 0x0b,
	0x4d, 0x61, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x69, 0x0a, 0x10, 0x4d,
	0x61, 0x70, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x19, 0x0a, 0x08, 0x6d, 0x61, 0x70, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x61, 0x70, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x70,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x72, 0x70, 0x6d, 0x12, 0x12, 0x0a, 0x04,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x6c, 0x6f, 0x61, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xb4, 0x01, 0x0a, 0x0b, 0x45, 0x43, 0x55, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x75, 0x65, 0x6c, 0x5f, 0x74,
	0x72, 0x69, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x66, 0x75, 0x65, 0x6c, 0x54,
	0x72, 0x69, 0x6d, 0x12, 0x23, 0x0a, 0x0d, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x72, 0x69, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x69, 0x67, 0x6e, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x72, 0x69, 0x6d, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x64, 0x6c, 0x65,
	0x5f, 0x72, 0x70, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x69, 0x64, 0x6c, 0x65,
	0x52, 0x70, 0x6d, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x76, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x72, 0x65, 0x76, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x12, 0x2b, 0x0a, 0x11, 0x74, 0x65, 0x6d, 0x70, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x74, 0x65, 0x6d,
	0x70, 0x43, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x42, 0x0a,
	0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x32, 0xa9, 0x02, 0x0a, 0x13, 0x4d, 0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63, 0x6c, 0x65,
	0x53, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x43, 0x0a, 0x0c, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x12, 0x15, 0x2e, 0x6d, 0x6f, 0x74, 0x6f,
	0x72, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x70, 0x75, 0x74,
	0x1a, 0x16, 0x2e, 0x6d, 0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x45, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x3c,
	0x0a, 0x0a, 0x47, 0x65, 0x74, 0x45, 0x43, 0x55, 0x4d, 0x61, 0x70, 0x73, 0x12, 0x17, 0x2e, 0x6d,
	0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x4d, 0x61, 0x70, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x6d, 0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63,
	0x6c, 0x65, 0x2e, 0x45, 0x43, 0x55, 0x4d, 0x61, 0x70, 0x73, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x0c,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x43, 0x55, 0x4d, 0x61, 0x70, 0x12, 0x1c, 0x2e, 0x6d,
	0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x4d, 0x61, 0x70, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x6d, 0x6f, 0x74,
	0x6f, 0x72, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x0e, 0x53, 0x65, 0x74, 0x45, 0x43, 0x55,
	0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x17, 0x2e, 0x6d, 0x6f, 0x74, 0x6f, 0x72,
	0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x45, 0x43, 0x55, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x73, 0x1a, 0x18, 0x2e, 0x6d, 0x6f, 0x74, 0x6f, 0x72, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x42, 0x2b, 0x5a,
	0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x79, 0x6f, 0x75, 0x72,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x2f, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x36, 0x35,
	0x30, 0x73, 0x69, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_motorcycle_proto_rawDescOnce sync.Once
	file_proto_motorcycle_proto_rawDescData = file_proto_motorcycle_proto_rawDesc
)

func file_proto_motorcycle_proto_rawDescGZIP() []byte {
	file_proto_motorcycle_proto_rawDescOnce.Do(func() {
		file_proto_motorcycle_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_motorcycle_proto_rawDescData)
	})
	return file_proto_motorcycle_proto_rawDescData
}

var file_proto_motorcycle_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_motorcycle_proto_goTypes = []interface{}{
	(*EngineData)(nil),       // 0: motorcycle.EngineData
	(*UserInput)(nil),        // 1: motorcycle.UserInput
	(*MapRow)(nil),           // 2: motorcycle.MapRow
	(*Map2D)(nil),            // 3: motorcycle.Map2D
	(*ECUMaps)(nil),          // 4: motorcycle.ECUMaps
	(*MapsRequest)(nil),      // 5: motorcycle.MapsRequest
	(*MapUpdateRequest)(nil), // 6: motorcycle.MapUpdateRequest
	(*ECUSettings)(nil),      // 7: motorcycle.ECUSettings
	(*UpdateStatus)(nil),     // 8: motorcycle.UpdateStatus
}
var file_proto_motorcycle_proto_depIdxs = []int32{
	2, // 0: motorcycle.Map2D.values:type_name -> motorcycle.MapRow
	3, // 1: motorcycle.ECUMaps.fuel_map:type_name -> motorcycle.Map2D
	3, // 2: motorcycle.ECUMaps.ignition_map:type_name -> motorcycle.Map2D
	3, // 3: motorcycle.ECUMaps.afr_map:type_name -> motorcycle.Map2D
	1, // 4: motorcycle.MotorcycleSimulator.StreamEngine:input_type -> motorcycle.UserInput
	5, // 5: motorcycle.MotorcycleSimulator.GetECUMaps:input_type -> motorcycle.MapsRequest
	6, // 6: motorcycle.MotorcycleSimulator.UpdateECUMap:input_type -> motorcycle.MapUpdateRequest
	7, // 7: motorcycle.MotorcycleSimulator.SetECUSettings:input_type -> motorcycle.ECUSettings
	0, // 8: motorcycle.MotorcycleSimulator.StreamEngine:output_type -> motorcycle.EngineData
	4, // 9: motorcycle.MotorcycleSimulator.GetECUMaps:output_type -> motorcycle.ECUMaps
	8, // 10: motorcycle.MotorcycleSimulator.UpdateECUMap:output_type -> motorcycle.UpdateStatus
	8, // 11: motorcycle.MotorcycleSimulator.SetECUSettings:output_type -> motorcycle.UpdateStatus
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_motorcycle_proto_init() }
func file_proto_motorcycle_proto_init() {
	if File_proto_motorcycle_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_motorcycle_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EngineData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_motorcycle_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_motorcycle_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapRow); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_motorcycle_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Map2D); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_motorcycle_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ECUMaps); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_motorcycle_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_motorcycle_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapUpdateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_motorcycle_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ECUSettings); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_motorcycle_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_motorcycle_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_motorcycle_proto_goTypes,
		DependencyIndexes: file_proto_motorcycle_proto_depIdxs,
		MessageInfos:      file_proto_motorcycle_proto_msgTypes,
	}.Build()
	File_proto_motorcycle_proto = out.File
	file_proto_motorcycle_proto_rawDesc = nil
	file_proto_motorcycle_proto_goTypes = nil
	file_proto_motorcycle_proto_depIdxs = nil
}
