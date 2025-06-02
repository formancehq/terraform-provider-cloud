# OrganizationFeature

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**OrganizationID** | **string** |  | 
**Name** | **string** |  | 
**CreatedAt** | **time.Time** |  | 

## Methods

### NewOrganizationFeature

`func NewOrganizationFeature(organizationID string, name string, createdAt time.Time, ) *OrganizationFeature`

NewOrganizationFeature instantiates a new OrganizationFeature object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrganizationFeatureWithDefaults

`func NewOrganizationFeatureWithDefaults() *OrganizationFeature`

NewOrganizationFeatureWithDefaults instantiates a new OrganizationFeature object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganizationID

`func (o *OrganizationFeature) GetOrganizationID() string`

GetOrganizationID returns the OrganizationID field if non-nil, zero value otherwise.

### GetOrganizationIDOk

`func (o *OrganizationFeature) GetOrganizationIDOk() (*string, bool)`

GetOrganizationIDOk returns a tuple with the OrganizationID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationID

`func (o *OrganizationFeature) SetOrganizationID(v string)`

SetOrganizationID sets OrganizationID field to given value.


### GetName

`func (o *OrganizationFeature) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *OrganizationFeature) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *OrganizationFeature) SetName(v string)`

SetName sets Name field to given value.


### GetCreatedAt

`func (o *OrganizationFeature) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *OrganizationFeature) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *OrganizationFeature) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


