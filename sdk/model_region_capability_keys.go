/*
Membership API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package sdk

import (
	"encoding/json"
	"fmt"
)

// RegionCapabilityKeys the model 'RegionCapabilityKeys'
type RegionCapabilityKeys string

// List of RegionCapabilityKeys
const (
	MODULE_LIST RegionCapabilityKeys = "MODULE_LIST"
	EE RegionCapabilityKeys = "EE"
)

// All allowed values of RegionCapabilityKeys enum
var AllowedRegionCapabilityKeysEnumValues = []RegionCapabilityKeys{
	"MODULE_LIST",
	"EE",
}

func (v *RegionCapabilityKeys) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := RegionCapabilityKeys(value)
	for _, existing := range AllowedRegionCapabilityKeysEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid RegionCapabilityKeys", value)
}

// NewRegionCapabilityKeysFromValue returns a pointer to a valid RegionCapabilityKeys
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewRegionCapabilityKeysFromValue(v string) (*RegionCapabilityKeys, error) {
	ev := RegionCapabilityKeys(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for RegionCapabilityKeys: valid values are %v", v, AllowedRegionCapabilityKeysEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v RegionCapabilityKeys) IsValid() bool {
	for _, existing := range AllowedRegionCapabilityKeysEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to RegionCapabilityKeys value
func (v RegionCapabilityKeys) Ptr() *RegionCapabilityKeys {
	return &v
}

type NullableRegionCapabilityKeys struct {
	value *RegionCapabilityKeys
	isSet bool
}

func (v NullableRegionCapabilityKeys) Get() *RegionCapabilityKeys {
	return v.value
}

func (v *NullableRegionCapabilityKeys) Set(val *RegionCapabilityKeys) {
	v.value = val
	v.isSet = true
}

func (v NullableRegionCapabilityKeys) IsSet() bool {
	return v.isSet
}

func (v *NullableRegionCapabilityKeys) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRegionCapabilityKeys(val *RegionCapabilityKeys) *NullableRegionCapabilityKeys {
	return &NullableRegionCapabilityKeys{value: val, isSet: true}
}

func (v NullableRegionCapabilityKeys) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRegionCapabilityKeys) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

