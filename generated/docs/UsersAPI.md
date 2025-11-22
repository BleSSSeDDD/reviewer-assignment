# \UsersAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**UsersGetReviewGet**](UsersAPI.md#UsersGetReviewGet) | **Get** /users/getReview | Получить PR&#39;ы, где пользователь назначен ревьювером
[**UsersSetIsActivePost**](UsersAPI.md#UsersSetIsActivePost) | **Post** /users/setIsActive | Установить флаг активности пользователя



## UsersGetReviewGet

> UsersGetReviewGet200Response UsersGetReviewGet(ctx).UserId(userId).Execute()

Получить PR'ы, где пользователь назначен ревьювером

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
	userId := "userId_example" // string | Идентификатор пользователя

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UsersAPI.UsersGetReviewGet(context.Background()).UserId(userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UsersAPI.UsersGetReviewGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UsersGetReviewGet`: UsersGetReviewGet200Response
	fmt.Fprintf(os.Stdout, "Response from `UsersAPI.UsersGetReviewGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUsersGetReviewGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **string** | Идентификатор пользователя | 

### Return type

[**UsersGetReviewGet200Response**](UsersGetReviewGet200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UsersSetIsActivePost

> UsersSetIsActivePost200Response UsersSetIsActivePost(ctx).UsersSetIsActivePostRequest(usersSetIsActivePostRequest).Execute()

Установить флаг активности пользователя

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
	usersSetIsActivePostRequest := *openapiclient.NewUsersSetIsActivePostRequest("UserId_example", false) // UsersSetIsActivePostRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UsersAPI.UsersSetIsActivePost(context.Background()).UsersSetIsActivePostRequest(usersSetIsActivePostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UsersAPI.UsersSetIsActivePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UsersSetIsActivePost`: UsersSetIsActivePost200Response
	fmt.Fprintf(os.Stdout, "Response from `UsersAPI.UsersSetIsActivePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUsersSetIsActivePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **usersSetIsActivePostRequest** | [**UsersSetIsActivePostRequest**](UsersSetIsActivePostRequest.md) |  | 

### Return type

[**UsersSetIsActivePost200Response**](UsersSetIsActivePost200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

