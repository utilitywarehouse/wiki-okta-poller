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

//var WikiInterval string = os.Getenv("WikiInterval")

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
		fmt.Println("env variable wiki.oktaprovider not set")
		os.Exit(1)
	}

	// get duration
	WikiIntervalInt, _ := strconv.Atoi(os.Getenv("WikiInterval"))
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
			usercheck, _ := wikiapi.DoesUserExist(userdetails.Profile.Email)
			if usercheck {
				fmt.Println("--", userdetails.Profile.Email, "from okta ALREADY EXISTS in wiki")
				// add user to wiki group based on okta group
				if strings.Contains(wikigroupuserjson, userdetails.Profile.Email) {
					fmt.Println("--", userdetails.Profile.Email, "ALREADY EXISTS in wiki group -", wikigroup)
				} else {
					if wikiapi.AddUserNameToGroupName(wikigroup, userdetails.Profile.Email) {
						fmt.Println("added user", userdetails.Profile.Email, "to wiki group -", wikigroup)
					} else {
						fmt.Println("ERROR adding user", userdetails.Profile.Email, "to", wikigroup)
					}
				}
			} else {
				fmt.Println("user from okta MISSING in wiki", userdetails.Profile.Email, usercheck)
				// if user is missing create user with okta link and default group
				createuserresult, slug := wikiapi.CreateNewUser(userdetails.Profile.Email, ""+userdetails.Profile.FirstName+" "+userdetails.Profile.LastName+"", "placeholderpas$wordb3causeOKTA3213", "3")
				if createuserresult {
					fmt.Println("Created user in wiki", userdetails.Profile.Email)
				} else {
					fmt.Println("ERROR creating user in wiki", userdetails.Profile.Email, slug)
				}
				// now add user to wiki group based on okta group
				if wikiapi.AddUserNameToGroupName(wikigroup, userdetails.Profile.Email) {
					fmt.Println("added user", userdetails.Profile.Email, "to wiki group -", wikigroup)
				} else {
					fmt.Println("ERROR adding user", userdetails.Profile.Email, "to wiki group -", wikigroup)
				}
			}
		}
		// get list of users and details from wiki
		var wikiuserlistout wikiapi.GetUsersFromGroup
		json.Unmarshal([]byte(wikigroupuserjson), &wikiuserlistout)
		// check wiki user list against okta group and remove users not in okta group
		for _, wikiuserdetails := range wikiuserlistout.Groups.Single.Users {
			fmt.Println("testsetsetset", wikiuserdetails.Email)
			//fmt.Println("outoutoutoutout - okta", oktauserlist)
			if !strings.Contains(oktauserlist, wikiuserdetails.Email) {
				fmt.Println("user exists in wiki but not in okta", wikiuserdetails.Email)
				// remove user from wiki if condition met
				if wikiapi.RemoveUserNameFromGroupName(wikigroup, wikiuserdetails.Email) {
					fmt.Println("user removed from wiki group", wikiuserdetails.Email, "-", wikigroup)
				} else {
					fmt.Println("ERROR removing user", wikiuserdetails.Email, "from wiki group -", wikigroup)
				}
			}
		}
	}
}
