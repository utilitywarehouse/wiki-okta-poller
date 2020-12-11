package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"wiki-poller/oktaapi"
	"wiki-poller/wikiapi"
)

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func main() {

	// 0. check vars from env

	if len(oktaapi.Oktadomain) < 1 {
		fmt.Println("env variable oktadomain not set")
		os.Exit(1)
	}
	if len(oktaapi.Oktatoken) < 1 {
		fmt.Println("env variable oktatoken not set")
		os.Exit(1)
	}
	if len(wikiapi.WikiToken) < 1 {
		fmt.Println("env variable wikitoken not set")
		os.Exit(1)
	}
	if len(wikiapi.WikiUrl) < 1 {
		fmt.Println("env variable wikiurl not set")
		os.Exit(1)
	}
	if len(wikiapi.Oktaprovider) < 1 {
		fmt.Println("env variable oktaprovider not set")
		os.Exit(1)
	}

	// get duration
	WikiIntervalInt, _ := strconv.Atoi(os.Getenv("WikiInterval"))
	if WikiIntervalInt == 0 {
		fmt.Println("env variable WikiInterval not set")
		os.Exit(1)
	}
	WikiIntervalDuration := time.Duration(WikiIntervalInt)
	//run ticker
	doEvery(WikiIntervalDuration*time.Second, mainprog)
}

func mainprog(t time.Time) {

	// 1. get list of groups
	oktagroupstartswith := "uw_ag_wiki-editors-"
	var oktagrouplistoutput oktaapi.OktaListGroups
	oktagetbody := oktaapi.GetListOfGroups(oktagroupstartswith)
	json.Unmarshal([]byte(oktagetbody), &oktagrouplistoutput)
	if len(oktagetbody) < 3 {
		fmt.Println("No groups match prefix -", oktagroupstartswith, "- nothing to process")
	}

	// 2. for each group get email list and add to wiki
	for _, grouplistItem := range oktagrouplistoutput {
		oktauserlist := oktaapi.GetUserEmailsFromGroup(grouplistItem.ID)
		var oktauserlistout oktaapi.OktaGroup
		json.Unmarshal([]byte(oktauserlist), &oktauserlistout)
		// extrapolate wiki group from okta group name
		wikigroup := strings.TrimPrefix(grouplistItem.Profile.Name, oktagroupstartswith)
		// get group json from wiki
		wikigroupuserjson := wikiapi.GetUsersFromGroupName(wikigroup)
		// parse list of users and get details
		for _, userdetails := range oktauserlistout {
			// check if user exists

			useremail := strings.ToLower(userdetails.Profile.Email)

			usercheck, _ := wikiapi.DoesUserExist(useremail)
			if usercheck {
				// add user to wiki group based on okta group
				if !strings.Contains(wikigroupuserjson, useremail) {
					if wikiapi.AddUserNameToGroupName(wikigroup, useremail) {
						fmt.Println("1 added user", useremail, "to wiki group -", wikigroup)
					} else {
						fmt.Println("ERROR adding user", useremail, "to", wikigroup)
					}
				}
			} else {
				// if user is missing create user with okta link and default group
				createuserresult, slug := wikiapi.CreateNewUser(useremail, ""+userdetails.Profile.FirstName+" "+userdetails.Profile.LastName+"", "placeholderpas$wordb3causeOKTA3213", "3")
				if createuserresult {
					fmt.Println("Created user in wiki", useremail)
				} else {
					fmt.Println("ERROR creating user in wiki", useremail, slug)
				}
				// now add user to wiki group based on okta group
				if wikiapi.AddUserNameToGroupName(wikigroup, useremail) {
					fmt.Println("2 added user", useremail, "to wiki group -", wikigroup)
				} else {
					fmt.Println("ERROR adding user", useremail, "to wiki group -", wikigroup)
				}
			}
		}
		// get list of users and details from wiki
		var wikiuserlistout wikiapi.GetUsersFromGroup
		json.Unmarshal([]byte(wikigroupuserjson), &wikiuserlistout)
		// check wiki user list against okta group and remove users not in okta group
		for _, wikiuserdetails := range wikiuserlistout.Groups.Single.Users {
			if !strings.Contains(strings.ToLower(oktauserlist), strings.ToLower(wikiuserdetails.Email)) {
				// remove user from wiki if condition met
				if wikiapi.RemoveUserNameFromGroupName(wikigroup, wikiuserdetails.Email) {
					fmt.Println("user removed from wiki group", wikiuserdetails.Email, "-", wikigroup)
				} else {
					fmt.Println("ERROR removing user", wikiuserdetails.Email, "from wiki group -", wikigroup)
				}
			}
		}
	}
	fmt.Println("Completed poll:", time.Now())
}
