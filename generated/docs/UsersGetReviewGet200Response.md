# UsersGetReviewGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UserId** | **string** |  | 
**PullRequests** | [**[]PullRequestShort**](PullRequestShort.md) |  | 

## Methods

### NewUsersGetReviewGet200Response

`func NewUsersGetReviewGet200Response(userId string, pullRequests []PullRequestShort, ) *UsersGetReviewGet200Response`

NewUsersGetReviewGet200Response instantiates a new UsersGetReviewGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUsersGetReviewGet200ResponseWithDefaults

`func NewUsersGetReviewGet200ResponseWithDefaults() *UsersGetReviewGet200Response`

NewUsersGetReviewGet200ResponseWithDefaults instantiates a new UsersGetReviewGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUserId

`func (o *UsersGetReviewGet200Response) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *UsersGetReviewGet200Response) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *UsersGetReviewGet200Response) SetUserId(v string)`

SetUserId sets UserId field to given value.


### GetPullRequests

`func (o *UsersGetReviewGet200Response) GetPullRequests() []PullRequestShort`

GetPullRequests returns the PullRequests field if non-nil, zero value otherwise.

### GetPullRequestsOk

`func (o *UsersGetReviewGet200Response) GetPullRequestsOk() (*[]PullRequestShort, bool)`

GetPullRequestsOk returns a tuple with the PullRequests field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPullRequests

`func (o *UsersGetReviewGet200Response) SetPullRequests(v []PullRequestShort)`

SetPullRequests sets PullRequests field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


