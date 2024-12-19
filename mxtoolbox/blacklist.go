package mxtoolbox

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type BlackListMonitor struct {
	ApiUrl     string `json:"apiUrl"`
	Token      string `json:"token"`
	DomainName string `json:"domainName"`
}

type MxtoolboxResponse struct {
	MonitorUID   string `json:"MonitorUID"`
	ActionString string `json:"ActionString"`
	Action       struct {
		RawString          string      `json:"RawString"`
		Command            string      `json:"Command"`
		Address            string      `json:"Address"`
		CommandParts       []string    `json:"CommandParts"`
		Threshold          int         `json:"Threshold"`
		AllowPrivateIP     bool        `json:"AllowPrivateIP"`
		SpecialServer      interface{} `json:"SpecialServer"`
		IsAll              bool        `json:"IsAll"`
		HTTPRegex          interface{} `json:"HTTP_Regex"`
		TCPPort            int         `json:"TCP_Port"`
		SMTPPort           int         `json:"SMTP_Port"`
		MAILFLOWIP         interface{} `json:"MAILFLOW_IP"`
		ActionString       string      `json:"ActionString"`
		LarArgument        string      `json:"LarArgument"`
		HasAdvancedOptions bool        `json:"HasAdvancedOptions"`
		IsMonitorable      bool        `json:"IsMonitorable"`
	} `json:"Action"`
	LastTransition     string        `json:"LastTransition"`
	LastChecked        string        `json:"LastChecked"`
	MxRep              string        `json:"MxRep"`
	HistoryURL         string        `json:"HistoryUrl"`
	Name               string        `json:"Name"`
	TimeElapsed        string        `json:"TimeElapsed"`
	RecordCount        int           `json:"RecordCount"`
	LarUID             string        `json:"LarUID"`
	LarID              string        `json:"LarId"`
	FrequencyInMinutes int           `json:"FrequencyInMinutes"`
	CurrentStatus      int           `json:"CurrentStatus"`
	ExpiresIn          interface{}   `json:"ExpiresIn"`
	ExpirationInHours  int           `json:"ExpirationInHours"`
	IsActive           bool          `json:"IsActive"`
	DomainSubscription interface{}   `json:"DomainSubscription"`
	Tags               []interface{} `json:"Tags"`
	Failing            []interface{} `json:"Failing"`
	Warnings           []interface{} `json:"Warnings"`
	Policies           []struct {
		ID                           int           `json:"Id"`
		AlertGroupID                 int           `json:"AlertGroupId"`
		PolicyUID                    string        `json:"PolicyUID"`
		TagMaps                      []interface{} `json:"TagMaps"`
		EndPoints                    []interface{} `json:"EndPoints"`
		IsEnabled                    bool          `json:"IsEnabled"`
		Title                        string        `json:"Title"`
		URL                          interface{}   `json:"Url"`
		Description                  interface{}   `json:"Description"`
		PolicyFilters                []interface{} `json:"PolicyFilters"`
		NumMinutesToRepeat           interface{}   `json:"NumMinutesToRepeat"`
		NumMinutesToDelay            interface{}   `json:"NumMinutesToDelay"`
		SortOrder                    int           `json:"SortOrder"`
		StackedNotificationFrequency interface{}   `json:"StackedNotificationFrequency"`
		IsSendToAllTagsEnabled       bool          `json:"IsSendToAllTagsEnabled"`
		TagNames                     string        `json:"TagNames"`
	} `json:"Policies"`
	StatusSummary string `json:"StatusSummary"`
}

// IsDomainBlacklisted checks if the domain is blacklisted
func (bl BlackListMonitor) IsDomainBlacklisted() (bool, error) {
	slog.Info("Blacklist: Getting Domain Status for", slog.String("DomainName", bl.DomainName))
	client := &http.Client{}

	req, err := http.NewRequest("GET", bl.ApiUrl, nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Set("Authorization", bl.Token)
	slog.Info("Sending Request to mxtoolbox")
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error sending request: %s", err)
	}
	defer resp.Body.Close()
	var data MxtoolboxResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return false, fmt.Errorf("error decoding JSON: %s", err)
	}

	if data.StatusSummary != "Not Blacklisted" {
		return true, fmt.Errorf("domain %s is BlackListed", bl.DomainName)
	}
	slog.Info("Domain Check completed", slog.String("DomainName", bl.DomainName))
	return false, nil

}
