/*
Membership API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package sdk

import (
	"encoding/json"
)

// checks if the CreateOrganizationResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateOrganizationResponse{}

// CreateOrganizationResponse struct for CreateOrganizationResponse
type CreateOrganizationResponse struct {
	Data *OrganizationExpanded `json:"data,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CreateOrganizationResponse CreateOrganizationResponse

// NewCreateOrganizationResponse instantiates a new CreateOrganizationResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateOrganizationResponse() *CreateOrganizationResponse {
	this := CreateOrganizationResponse{}
	return &this
}

// NewCreateOrganizationResponseWithDefaults instantiates a new CreateOrganizationResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateOrganizationResponseWithDefaults() *CreateOrganizationResponse {
	this := CreateOrganizationResponse{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *CreateOrganizationResponse) GetData() OrganizationExpanded {
	if o == nil || IsNil(o.Data) {
		var ret OrganizationExpanded
		return ret
	}
	return *o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationResponse) GetDataOk() (*OrganizationExpanded, bool) {
	if o == nil || IsNil(o.Data) {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *CreateOrganizationResponse) HasData() bool {
	if o != nil && !IsNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given OrganizationExpanded and assigns it to the Data field.
func (o *CreateOrganizationResponse) SetData(v OrganizationExpanded) {
	o.Data = &v
}

func (o CreateOrganizationResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateOrganizationResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Data) {
		toSerialize["data"] = o.Data
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CreateOrganizationResponse) UnmarshalJSON(data []byte) (err error) {
	varCreateOrganizationResponse := _CreateOrganizationResponse{}

	err = json.Unmarshal(data, &varCreateOrganizationResponse)

	if err != nil {
		return err
	}

	*o = CreateOrganizationResponse(varCreateOrganizationResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "data")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCreateOrganizationResponse struct {
	value *CreateOrganizationResponse
	isSet bool
}

func (v NullableCreateOrganizationResponse) Get() *CreateOrganizationResponse {
	return v.value
}

func (v *NullableCreateOrganizationResponse) Set(val *CreateOrganizationResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateOrganizationResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateOrganizationResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateOrganizationResponse(val *CreateOrganizationResponse) *NullableCreateOrganizationResponse {
	return &NullableCreateOrganizationResponse{value: val, isSet: true}
}

func (v NullableCreateOrganizationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateOrganizationResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


