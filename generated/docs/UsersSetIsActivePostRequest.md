# UsersSetIsActivePostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UserId** | **string** |  | 
**IsActive** | **bool** |  | 

## Methods

### NewUsersSetIsActivePostRequest

`func NewUsersSetIsActivePostRequest(userId string, isActive bool, ) *UsersSetIsActivePostRequest`

NewUsersSetIsActivePostRequest instantiates a new UsersSetIsActivePostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUsersSetIsActivePostRequestWithDefaults

`func NewUsersSetIsActivePostRequestWithDefaults() *UsersSetIsActivePostRequest`

NewUsersSetIsActivePostRequestWithDefaults instantiates a new UsersSetIsActivePostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUserId

`func (o *UsersSetIsActivePostRequest) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *UsersSetIsActivePostRequest) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *UsersSetIsActivePostRequest) SetUserId(v string)`

SetUserId sets UserId field to given value.


### GetIsActive

`func (o *UsersSetIsActivePostRequest) GetIsActive() bool`

GetIsActive returns the IsActive field if non-nil, zero value otherwise.

### GetIsActiveOk

`func (o *UsersSetIsActivePostRequest) GetIsActiveOk() (*bool, bool)`

GetIsActiveOk returns a tuple with the IsActive field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsActive

`func (o *UsersSetIsActivePostRequest) SetIsActive(v bool)`

SetIsActive sets IsActive field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


