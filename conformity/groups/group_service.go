package groups

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"terraform-provider-conformity/conformity/models"
	"terraform-provider-conformity/conformity/provider"
)

type GroupService struct {
	SuffixURL string
	Client    *provider.Client
}

type DeleteRespBody struct {
	Meta struct {
		Status string `json:"status"`
	} `json:"meta"`
}

type GroupDataList struct {
	Data []models.Group `json:"data"`
}

type GroupData struct {
	Data models.Group `json:"data"`
}

var singleton *GroupService
var once sync.Once

func GetGroupService(client *provider.Client) *GroupService {
	once.Do(func() {
		singleton = &GroupService{}
		singleton.SuffixURL = "/groups"
		singleton.Client = client
	})
	return singleton
}

func (gs *GroupService) DoCreateGroup(baseURL string, group *models.Group) (ng models.Group, err error) {
	path := baseURL + gs.SuffixURL
	log.Println("[DEBUG] DoCreateGroup of GS started")

	groupData := GroupData{}
	groupData.Data = *group
	log.Printf("[DEBUG] Before Group Create Marshal %v\n", groupData)

	rb, err := json.Marshal(groupData)
	if err != nil {
		return ng, err
	}
	log.Printf("[DEBUG] Request body is %s\n", string(rb))
	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(string(rb)))
	if err != nil {
		return ng, err
	}
	body, err := gs.Client.DoRequest(req)
	if err != nil {
		return ng, fmt.Errorf("unable to issue create group API call: %s", err.Error())
	}
	err = json.Unmarshal(body, &groupData)
	log.Printf("[DEBUG] After Group Create Unmarshal %v\n", groupData)
	if err != nil {
		log.Printf("[DEBUG] Unmarshal error %v\n", err)
		return ng, err
	}
	return groupData.Data, err
}

func (gs *GroupService) DoUpdateGroup(baseURL string, group *models.Group) (ng models.Group, err error) {
	path := baseURL + gs.SuffixURL + "/" + group.ID
	log.Println("[DEBUG] DoUpdateGroup of GS started")

	groupData := GroupData{}
	groupData.Data = *group

	rb, err := json.Marshal(groupData)
	if err != nil {
		return ng, err
	}
	log.Printf("[DEBUG] Request body is %s\n", string(rb))
	req, err := http.NewRequest(http.MethodPatch, path, strings.NewReader(string(rb)))
	if err != nil {
		return ng, err
	}
	body, err := gs.Client.DoRequest(req)
	if err != nil {
		return ng, fmt.Errorf("unable to issue create group API call: %s", err.Error())
	}
	err = json.Unmarshal(body, &groupData)
	if err != nil {
		return ng, err
	}
	return groupData.Data, err
}

func (gs *GroupService) DoGetGroup(baseURL string, groupID string) (group models.Group, err error) {
	path := baseURL + gs.SuffixURL + "/" + groupID
	log.Println("[DEBUG] DoGetGroup of GS has started!")

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return group, err
	}
	body, err := gs.Client.DoRequest(req)
	if err != nil {
		return group, fmt.Errorf("Unable to get Group details: %s", err.Error())
	}
	groups := GroupDataList{}
	err = json.Unmarshal(body, &groups)
	if err != nil && len(groups.Data) == 0 {
		log.Printf("[DEBUG] No Groups found or error %v\n\n", err)
		return group, err
	}
	return groups.Data[0], err
}

func (gs *GroupService) DoDeleteGroup(baseURL string, groupID string) (err error) {
	path := baseURL + gs.SuffixURL + "/" + groupID
	log.Printf("[DEBUG] DoDeleteGroup for group id %s\n", groupID)

	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	body, err := gs.Client.DoRequest(req)
	if err != nil {
		return fmt.Errorf("Unable to delete Group: %s", err.Error())
	}
	respBody := DeleteRespBody{}
	err = json.Unmarshal(body, &respBody)
	out, err := json.Marshal(respBody)
	log.Printf("[DEBUG] Delete group response struct = %s\n", string(out))
	if err != nil || respBody.Meta.Status != "deleted" {
		log.Printf("[DEBUG] Delete group failure with err = %s, body = %s\n", err, string(body))
		return fmt.Errorf("Unable to delete group : %s", err.Error())
	}
	return err
}

func (gs *GroupService) DoGetGroups(baseURL string) (g []map[string]interface{}, err error) {
	path := baseURL + gs.SuffixURL
	log.Println("[DEBUG] DoGetGroups of GS started")
	log.Printf("[DEBUG] Client is %v\n", gs.Client)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		log.Printf("[ERROR] Error while creating request %v\n", err)
		return nil, err
	}

	body, err := gs.Client.DoRequest(req)

	groupsData := GroupDataList{}
	err = json.Unmarshal(body, &groupsData)

	log.Printf("Groups: %v\n", groupsData)
	fgroups := flattenGroupsData(groupsData.Data)
	return fgroups, err
}

func flattenGroupsData(groups []models.Group) []map[string]interface{} {
	if groups != nil {
		fgroups := make([]map[string]interface{}, len(groups))

		for i, group := range groups {
			fgroup := make(map[string]interface{})

			fgroup["type"] = group.Type
			fgroup["id"] = group.ID
			fgroup["attributes_name"] = group.Attributes.Name
			fgroup["attributes_tags"] = group.Attributes.Tags
			fgroup["attributes_created_date"] = group.Attributes.CreatedDate
			fgroup["attributes_last_modified_date"] = group.Attributes.LastModifiedDate

			fgroup["relationships_organisation_data_type"] = group.Relationships.Organisation.Data.Type
			fgroup["relationships_organisation_data_id"] = group.Relationships.Organisation.Data.ID
			faccounts := make([]map[string]interface{}, len(group.Relationships.Accounts.Data))
			for j, account := range group.Relationships.Accounts.Data {
				faccount := make(map[string]interface{})
				faccount["type"] = account.Type
				faccount["id"] = account.ID
				faccounts[j] = faccount
			}
			fgroup["relationships_accounts_data"] = faccounts

			fgroups[i] = fgroup
		}
		return fgroups
	}

	return make([]map[string]interface{}, 0)
}
