# \PullRequestsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PullRequestCreatePost**](PullRequestsAPI.md#PullRequestCreatePost) | **Post** /pullRequest/create | Создать PR и автоматически назначить до 2 ревьюверов из команды автора
[**PullRequestMergePost**](PullRequestsAPI.md#PullRequestMergePost) | **Post** /pullRequest/merge | Пометить PR как MERGED (идемпотентная операция)
[**PullRequestReassignPost**](PullRequestsAPI.md#PullRequestReassignPost) | **Post** /pullRequest/reassign | Переназначить конкретного ревьювера на другого из его команды



## PullRequestCreatePost

> PullRequestCreatePost201Response PullRequestCreatePost(ctx).PullRequestCreatePostRequest(pullRequestCreatePostRequest).Execute()

Создать PR и автоматически назначить до 2 ревьюверов из команды автора

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	pullRequestCreatePostRequest := *openapiclient.NewPullRequestCreatePostRequest("PullRequestId_example", "PullRequestName_example", "AuthorId_example") // PullRequestCreatePostRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PullRequestsAPI.PullRequestCreatePost(context.Background()).PullRequestCreatePostRequest(pullRequestCreatePostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PullRequestsAPI.PullRequestCreatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PullRequestCreatePost`: PullRequestCreatePost201Response
	fmt.Fprintf(os.Stdout, "Response from `PullRequestsAPI.PullRequestCreatePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPullRequestCreatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pullRequestCreatePostRequest** | [**PullRequestCreatePostRequest**](PullRequestCreatePostRequest.md) |  | 

### Return type

[**PullRequestCreatePost201Response**](PullRequestCreatePost201Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PullRequestMergePost

> PullRequestCreatePost201Response PullRequestMergePost(ctx).PullRequestMergePostRequest(pullRequestMergePostRequest).Execute()

Пометить PR как MERGED (идемпотентная операция)

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	pullRequestMergePostRequest := *openapiclient.NewPullRequestMergePostRequest("PullRequestId_example") // PullRequestMergePostRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PullRequestsAPI.PullRequestMergePost(context.Background()).PullRequestMergePostRequest(pullRequestMergePostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PullRequestsAPI.PullRequestMergePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PullRequestMergePost`: PullRequestCreatePost201Response
	fmt.Fprintf(os.Stdout, "Response from `PullRequestsAPI.PullRequestMergePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPullRequestMergePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pullRequestMergePostRequest** | [**PullRequestMergePostRequest**](PullRequestMergePostRequest.md) |  | 

### Return type

[**PullRequestCreatePost201Response**](PullRequestCreatePost201Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PullRequestReassignPost

> PullRequestReassignPost200Response PullRequestReassignPost(ctx).PullRequestReassignPostRequest(pullRequestReassignPostRequest).Execute()

Переназначить конкретного ревьювера на другого из его команды

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	pullRequestReassignPostRequest := *openapiclient.NewPullRequestReassignPostRequest("PullRequestId_example", "OldUserId_example") // PullRequestReassignPostRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PullRequestsAPI.PullRequestReassignPost(context.Background()).PullRequestReassignPostRequest(pullRequestReassignPostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PullRequestsAPI.PullRequestReassignPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PullRequestReassignPost`: PullRequestReassignPost200Response
	fmt.Fprintf(os.Stdout, "Response from `PullRequestsAPI.PullRequestReassignPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPullRequestReassignPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pullRequestReassignPostRequest** | [**PullRequestReassignPostRequest**](PullRequestReassignPostRequest.md) |  | 

### Return type

[**PullRequestReassignPost200Response**](PullRequestReassignPost200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

