package wikiapi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/machinebox/graphql"
)

var WikiToken string = os.Getenv("WikiToken")
var WikiUrl string = os.Getenv("WikiUrl")
var Oktaprovider string = os.Getenv("Oktaprovider")

type DoesUserExistType struct {
	Users struct {
		Search []struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		} `json:"search"`
	} `json:"users"`
}

type AddUserToGroupType struct {
	Groups struct {
		AssignUser struct {
			ResponseResult struct {
				Succeeded bool `json:"succeeded"`
			} `json:"responseResult"`
		} `json:"assignUser"`
	} `json:"groups"`
}

type UnassignUserFromGroupType struct {
	Groups struct {
		UnassignUser struct {
			ResponseResult struct {
				Succeeded bool `json:"succeeded"`
			} `json:"responseResult"`
		} `json:"unassignUser"`
	} `json:"groups"`
}

type CreateNewUserType struct {
	Users struct {
		Create struct {
			ResponseResult struct {
				Succeeded bool   `json:"succeeded"`
				ErrorCode int    `json:"errorCode"`
				Slug      string `json:"slug"`
				Message   string `json:"message"`
			} `json:"responseResult"`
		} `json:"create"`
	} `json:"users"`
}

type DeleteUserType struct {
	Users struct {
		Delete struct {
			ResponseResult struct {
				Succeeded bool   `json:"succeeded"`
				ErrorCode int    `json:"errorCode"`
				Slug      string `json:"slug"`
				Message   string `json:"message"`
			} `json:"responseResult"`
		} `json:"delete"`
	} `json:"users"`
}

type GetGroupIDType struct {
	Groups struct {
		List []struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"list"`
	} `json:"groups"`
}

type GetUsersFromGroup struct {
	Groups struct {
		Single struct {
			Users []struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"users"`
		} `json:"single"`
	} `json:"groups"`
}

// add useremail to groupname
func AddUserNameToGroupName(groupname string, useremail string) bool {

	_, userId := DoesUserExist(useremail)
	groupId := GetGroupId(groupname)
	response := AddUserToGroupId(groupId, userId)
	return response
}

// remove useremail from groupname
func RemoveUserNameFromGroupName(groupname string, useremail string) bool {
	_, userId := DoesUserExist(useremail)
	groupId := GetGroupId(groupname)
	response := RemoveUserFromGroupId(groupId, userId)
	return response
}

// check if user exists and return id
func DoesUserExist(useremail string) (bool, int) {
	graphqlClient := graphql.NewClient(WikiUrl)
	graphqlRequest := graphql.NewRequest(`
  query{
      users {
        search(query: "` + useremail + `"){
       id, email
        }
      }
    }
`)
	graphqlRequest.Header.Set("Authorization", WikiToken)
	var graphqlResponse DoesUserExistType
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}

	if len(graphqlResponse.Users.Search) == 0 {
		return false, 999
	}
	if strings.Contains(graphqlResponse.Users.Search[0].Email, useremail) {
		return true, graphqlResponse.Users.Search[0].ID
	} else {
		return false, graphqlResponse.Users.Search[0].ID
	}
}

// create new user
func CreateNewUser(useremail string, username string, userpasswd string, usergroups string) (bool, string) {
	graphqlClient := graphql.NewClient(WikiUrl)
	graphqlRequest := graphql.NewRequest(`
  mutation{users{create(
    email: "` + useremail + `"
    name: "` + username + `"
    passwordRaw: "` + userpasswd + `"
    providerKey: "` + Oktaprovider + `"
    groups: [` + usergroups + `]
    mustChangePassword: false
    sendWelcomeEmail: false
    ){
       responseResult{succeeded,errorCode,slug,message}
      }
       }
  }
`)
	graphqlRequest.Header.Set("Authorization", WikiToken)
	var graphqlResponse CreateNewUserType
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		fmt.Println("ERROR:", err)
	}

	if graphqlResponse.Users.Create.ResponseResult.Succeeded {
		return graphqlResponse.Users.Create.ResponseResult.Succeeded, ""
	} else {
		return graphqlResponse.Users.Create.ResponseResult.Succeeded, graphqlResponse.Users.Create.ResponseResult.Slug
	}
}

// delete userid
func DeleteUserId(userid int, replaceid int) (bool, string) {
	replaceidstr := strconv.Itoa(replaceid)
	useridstr := strconv.Itoa(userid)
	graphqlClient := graphql.NewClient(WikiUrl)
	graphqlRequest := graphql.NewRequest(`
  mutation{users{delete(
    id: ` + useridstr + `
    replaceId: ` + replaceidstr + `
    ){
      responseResult{succeeded,errorCode,slug,message}
      }
       }
  }
`)
	graphqlRequest.Header.Set("Authorization", WikiToken)
	var graphqlResponse DeleteUserType
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		fmt.Println("ERROR:", err)
	}

	if graphqlResponse.Users.Delete.ResponseResult.Succeeded {
		return graphqlResponse.Users.Delete.ResponseResult.Succeeded, ""
	} else {
		return graphqlResponse.Users.Delete.ResponseResult.Succeeded, graphqlResponse.Users.Delete.ResponseResult.Slug
	}
}

// add userid to groupid
func AddUserToGroupId(groupid int, userid int) bool {
	groupidstr := strconv.Itoa(groupid)
	useridstr := strconv.Itoa(userid)
	graphqlClient := graphql.NewClient(WikiUrl)
	graphqlRequest := graphql.NewRequest(`
  mutation{groups{assignUser(
  	groupId:` + groupidstr + `, 
 	userId:` + useridstr + `){
 	responseResult{succeeded}
 	 }
	}
}
`)
	graphqlRequest.Header.Set("Authorization", WikiToken)
	var graphqlResponse AddUserToGroupType
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		fmt.Println("ERROR:", err)
	}
	return graphqlResponse.Groups.AssignUser.ResponseResult.Succeeded
}

// remove userid from groupid
func RemoveUserFromGroupId(groupid int, userid int) bool {
	groupidstr := strconv.Itoa(groupid)
	useridstr := strconv.Itoa(userid)
	graphqlClient := graphql.NewClient(WikiUrl)
	graphqlRequest := graphql.NewRequest(`
  mutation{groups{unassignUser(
  	groupId:` + groupidstr + `, 
 	userId:` + useridstr + `){
 	responseResult{succeeded}
 	 }
	}
}
`)
	graphqlRequest.Header.Set("Authorization", WikiToken)
	var graphqlResponse UnassignUserFromGroupType
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		fmt.Println("ERROR:", err)
	}

	return graphqlResponse.Groups.UnassignUser.ResponseResult.Succeeded
}

// GetGroupId
func GetGroupId(groupname string) int {
	graphqlClient := graphql.NewClient(WikiUrl)
	graphqlRequest := graphql.NewRequest(`
  query {
    groups {
      list(filter: "` + groupname + `") {
        name
        id
      }
    }
  }
`)
	graphqlRequest.Header.Set("Authorization", WikiToken)
	var graphqlResponse GetGroupIDType
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		fmt.Println("ERROR:", err)
	}
	returnvalue := 0
	for _, groupItem := range graphqlResponse.Groups.List {
		if strings.ToLower(groupItem.Name) == strings.ToLower(groupname) {
			returnvalue = groupItem.ID
			break
		}
	}
	return returnvalue
}

// GetUsersFromGroupName
func GetUsersFromGroupName(groupname string) string {
	groupId := GetGroupId(groupname)
	groupidstr := strconv.Itoa(groupId)
	graphqlClient := graphql.NewClient(WikiUrl)
	graphqlRequest := graphql.NewRequest(`
  query {
    groups {
        single(id: ` + groupidstr + `){
			users{
			  name, email
			}
		}
    }
  }
`)
	graphqlRequest.Header.Set("Authorization", WikiToken)
	var graphqlResponse GetUsersFromGroup
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		fmt.Println("ERROR:", err)
	}
	body, _ := json.Marshal(&graphqlResponse)
	return string(body)
}
