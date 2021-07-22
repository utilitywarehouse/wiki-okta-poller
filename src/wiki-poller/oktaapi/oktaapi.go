package oktaapi

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var Oktadomain string = os.Getenv("Oktadomain")
var Oktatoken string = os.Getenv("Oktatoken")

type OktaGroup []struct {
	ID              string      `json:"id"`
	Status          string      `json:"status"`
	Created         time.Time   `json:"created"`
	Activated       time.Time   `json:"activated"`
	StatusChanged   time.Time   `json:"statusChanged"`
	LastLogin       time.Time   `json:"lastLogin"`
	LastUpdated     time.Time   `json:"lastUpdated"`
	PasswordChanged interface{} `json:"passwordChanged"`
	Type            struct {
		ID string `json:"id"`
	} `json:"type"`
	Profile struct {
		LastName             string      `json:"lastName"`
		HireDate             string      `json:"hireDate"`
		Manager              string      `json:"manager"`
		DisplayName          string      `json:"displayName"`
		NickName             string      `json:"nickName"`
		SecondEmail          interface{} `json:"secondEmail"`
		ManagerID            string      `json:"managerId"`
		Title                string      `json:"title"`
		Login                string      `json:"login"`
		EmployeeNumber       string      `json:"employeeNumber"`
		Division             string      `json:"division"`
		GithubUsername       string      `json:"githubUsername"`
		FirstName            string      `json:"firstName"`
		LastBamboohrSyncTime time.Time   `json:"lastBamboohrSyncTime"`
		MobilePhone          interface{} `json:"mobilePhone"`
		Organization         string      `json:"organization"`
		MiddleName           string      `json:"middleName"`
		Location             string      `json:"location"`
		JobClassificationID  string      `json:"jobClassificationID"`
		Department           string      `json:"department"`
		StartDate            string      `json:"startDate"`
		Email                string      `json:"email"`
	} `json:"profile"`
	Credentials struct {
		Password struct {
		} `json:"password"`
		RecoveryQuestion struct {
			Question string `json:"question"`
		} `json:"recovery_question"`
		Provider struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"provider"`
	} `json:"credentials"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

type OktaListGroups []struct {
	ID                    string    `json:"id"`
	Created               time.Time `json:"created"`
	LastUpdated           time.Time `json:"lastUpdated"`
	LastMembershipUpdated time.Time `json:"lastMembershipUpdated"`
	ObjectClass           []string  `json:"objectClass"`
	Type                  string    `json:"type"`
	Profile               struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"profile"`
	Links struct {
		Logo []struct {
			Name string `json:"name"`
			Href string `json:"href"`
			Type string `json:"type"`
		} `json:"logo"`
		Users struct {
			Href string `json:"href"`
		} `json:"users"`
		Apps struct {
			Href string `json:"href"`
		} `json:"apps"`
	} `json:"_links"`
}

func GetJson(url string) string {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", Oktatoken)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyoutput, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(bodyoutput)

}

// get list of groups from okta
func GetListOfGroups(groupstartswith string) string {

	url := "https://" + Oktadomain + "/api/v1/groups?q=" + groupstartswith + "&limit=1000"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", Oktatoken)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyoutput, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(bodyoutput)

}

// get email addresses from test group
func GetUserEmailsFromGroup(oktagroupid string) string {

	userurl := "https://" + Oktadomain + "/api/v1/groups/" + oktagroupid + "/users?limit=1000"
	getbody := GetJson(userurl)
	return getbody
}
