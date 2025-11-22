# PullRequestReassignPost200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Pr** | [**PullRequest**](PullRequest.md) |  | 
**ReplacedBy** | **string** | user_id нового ревьювера | 

## Methods

### NewPullRequestReassignPost200Response

`func NewPullRequestReassignPost200Response(pr PullRequest, replacedBy string, ) *PullRequestReassignPost200Response`

NewPullRequestReassignPost200Response instantiates a new PullRequestReassignPost200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPullRequestReassignPost200ResponseWithDefaults

`func NewPullRequestReassignPost200ResponseWithDefaults() *PullRequestReassignPost200Response`

NewPullRequestReassignPost200ResponseWithDefaults instantiates a new PullRequestReassignPost200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPr

`func (o *PullRequestReassignPost200Response) GetPr() PullRequest`

GetPr returns the Pr field if non-nil, zero value otherwise.

### GetPrOk

`func (o *PullRequestReassignPost200Response) GetPrOk() (*PullRequest, bool)`

GetPrOk returns a tuple with the Pr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPr

`func (o *PullRequestReassignPost200Response) SetPr(v PullRequest)`

SetPr sets Pr field to given value.


### GetReplacedBy

`func (o *PullRequestReassignPost200Response) GetReplacedBy() string`

GetReplacedBy returns the ReplacedBy field if non-nil, zero value otherwise.

### GetReplacedByOk

`func (o *PullRequestReassignPost200Response) GetReplacedByOk() (*string, bool)`

GetReplacedByOk returns a tuple with the ReplacedBy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplacedBy

`func (o *PullRequestReassignPost200Response) SetReplacedBy(v string)`

SetReplacedBy sets ReplacedBy field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


