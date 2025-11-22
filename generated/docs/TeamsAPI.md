# \TeamsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TeamAddPost**](TeamsAPI.md#TeamAddPost) | **Post** /team/add | Создать команду с участниками (создаёт/обновляет пользователей)
[**TeamGetGet**](TeamsAPI.md#TeamGetGet) | **Get** /team/get | Получить команду с участниками



## TeamAddPost

> TeamAddPost201Response TeamAddPost(ctx).Team(team).Execute()

Создать команду с участниками (создаёт/обновляет пользователей)

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
	team := *openapiclient.NewTeam("TeamName_example", []openapiclient.TeamMember{*openapiclient.NewTeamMember("UserId_example", "Username_example", false)}) // Team | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TeamsAPI.TeamAddPost(context.Background()).Team(team).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TeamsAPI.TeamAddPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TeamAddPost`: TeamAddPost201Response
	fmt.Fprintf(os.Stdout, "Response from `TeamsAPI.TeamAddPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTeamAddPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **team** | [**Team**](Team.md) |  | 

### Return type

[**TeamAddPost201Response**](TeamAddPost201Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TeamGetGet

> Team TeamGetGet(ctx).TeamName(teamName).Execute()

Получить команду с участниками

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
	teamName := "teamName_example" // string | Уникальное имя команды

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TeamsAPI.TeamGetGet(context.Background()).TeamName(teamName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TeamsAPI.TeamGetGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TeamGetGet`: Team
	fmt.Fprintf(os.Stdout, "Response from `TeamsAPI.TeamGetGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTeamGetGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **teamName** | **string** | Уникальное имя команды | 

### Return type

[**Team**](Team.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

