# PullRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PullRequestId** | **string** |  | 
**PullRequestName** | **string** |  | 
**AuthorId** | **string** |  | 
**Status** | **string** |  | 
**AssignedReviewers** | **[]string** | user_id назначенных ревьюверов (0..2) | 
**CreatedAt** | Pointer to **NullableTime** |  | [optional] 
**MergedAt** | Pointer to **NullableTime** |  | [optional] 

## Methods

### NewPullRequest

`func NewPullRequest(pullRequestId string, pullRequestName string, authorId string, status string, assignedReviewers []string, ) *PullRequest`

NewPullRequest instantiates a new PullRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPullRequestWithDefaults

`func NewPullRequestWithDefaults() *PullRequest`

NewPullRequestWithDefaults instantiates a new PullRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPullRequestId

`func (o *PullRequest) GetPullRequestId() string`

GetPullRequestId returns the PullRequestId field if non-nil, zero value otherwise.

### GetPullRequestIdOk

`func (o *PullRequest) GetPullRequestIdOk() (*string, bool)`

GetPullRequestIdOk returns a tuple with the PullRequestId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPullRequestId

`func (o *PullRequest) SetPullRequestId(v string)`

SetPullRequestId sets PullRequestId field to given value.


### GetPullRequestName

`func (o *PullRequest) GetPullRequestName() string`

GetPullRequestName returns the PullRequestName field if non-nil, zero value otherwise.

### GetPullRequestNameOk

`func (o *PullRequest) GetPullRequestNameOk() (*string, bool)`

GetPullRequestNameOk returns a tuple with the PullRequestName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPullRequestName

`func (o *PullRequest) SetPullRequestName(v string)`

SetPullRequestName sets PullRequestName field to given value.


### GetAuthorId

`func (o *PullRequest) GetAuthorId() string`

GetAuthorId returns the AuthorId field if non-nil, zero value otherwise.

### GetAuthorIdOk

`func (o *PullRequest) GetAuthorIdOk() (*string, bool)`

GetAuthorIdOk returns a tuple with the AuthorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorId

`func (o *PullRequest) SetAuthorId(v string)`

SetAuthorId sets AuthorId field to given value.


### GetStatus

`func (o *PullRequest) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PullRequest) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PullRequest) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetAssignedReviewers

`func (o *PullRequest) GetAssignedReviewers() []string`

GetAssignedReviewers returns the AssignedReviewers field if non-nil, zero value otherwise.

### GetAssignedReviewersOk

`func (o *PullRequest) GetAssignedReviewersOk() (*[]string, bool)`

GetAssignedReviewersOk returns a tuple with the AssignedReviewers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAssignedReviewers

`func (o *PullRequest) SetAssignedReviewers(v []string)`

SetAssignedReviewers sets AssignedReviewers field to given value.


### GetCreatedAt

`func (o *PullRequest) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *PullRequest) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *PullRequest) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *PullRequest) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### SetCreatedAtNil

`func (o *PullRequest) SetCreatedAtNil(b bool)`

 SetCreatedAtNil sets the value for CreatedAt to be an explicit nil

### UnsetCreatedAt
`func (o *PullRequest) UnsetCreatedAt()`

UnsetCreatedAt ensures that no value is present for CreatedAt, not even an explicit nil
### GetMergedAt

`func (o *PullRequest) GetMergedAt() time.Time`

GetMergedAt returns the MergedAt field if non-nil, zero value otherwise.

### GetMergedAtOk

`func (o *PullRequest) GetMergedAtOk() (*time.Time, bool)`

GetMergedAtOk returns a tuple with the MergedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMergedAt

`func (o *PullRequest) SetMergedAt(v time.Time)`

SetMergedAt sets MergedAt field to given value.

### HasMergedAt

`func (o *PullRequest) HasMergedAt() bool`

HasMergedAt returns a boolean if a field has been set.

### SetMergedAtNil

`func (o *PullRequest) SetMergedAtNil(b bool)`

 SetMergedAtNil sets the value for MergedAt to be an explicit nil

### UnsetMergedAt
`func (o *PullRequest) UnsetMergedAt()`

UnsetMergedAt ensures that no value is present for MergedAt, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


