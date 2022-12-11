package luckperms

import (
	"encoding/json"
	"fmt"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/uuid"
	"io"
	"net/http"
	"simple-proxy/minimessage"
)

type UserResponse struct {
	UniqueId string   `json:"uniqueId"`
	Username string   `json:"username"`
	Nodes    []Node   `json:"nodes"`
	Metadata Metadata `json:"metadata"`
}

type PermissionCheck struct {
	Result string `json:"result"`
}

type Node struct {
	Key   string `json:"key"`
	Value bool   `json:"value"`
}

type Metadata struct {
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
}

var restIp = "http://172.17.0.1:10420"
var CachedData = make(map[uuid.UUID]UserResponse)

func DisplayName(player proxy.Player) *Text {
	return minimessage.Parse(fmt.Sprintf("%s%s", Prefix(player), player.Username()))
}

func CollectData(player proxy.Player) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/%s", restIp, player.ID().String()), nil)
	if err != nil {
		fmt.Printf("Req errored: %s", err)
		return
	}

	req.Close = true
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Printf("Do errored: %s", err)
		return
	}

	byteValue, _ := io.ReadAll(resp.Body)

	var userResponse UserResponse
	err = json.Unmarshal(byteValue, &userResponse)
	if err != nil {
		fmt.Printf("Unmarshal errored: %s", err)
		return
	}

	CachedData[player.ID()] = userResponse
}

func Prefix(player proxy.Player) string {
	return CachedData[player.ID()].Metadata.Prefix
}

func HasPermission(source command.Source, permission string) bool {
	//player, ok := source.(proxy.Player)
	//if !ok { // always return true if the source is the console
	//	return true
	//}
	//
	//req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/%s/permissionCheck?permission=%s", restIp, player.ID().String(), permission), nil)
	//if err != nil {
	//	fmt.Printf("Req errored: %s", err)
	//	return false
	//}
	//
	//req.Close = true
	//req.Header.Add("Content-Type", "application/json")
	//
	//resp, err := http.DefaultClient.Do(req)
	//if resp != nil {
	//	defer resp.Body.Close()
	//}
	//if err != nil {
	//	fmt.Printf("Do errored: %s", err)
	//	return false
	//}
	//
	//byteValue, _ := io.ReadAll(resp.Body)
	//
	//var permissionCheck PermissionCheck
	//err = json.Unmarshal(byteValue, &permissionCheck)
	//if err != nil {
	//	fmt.Printf("Unmarshal errored: %s", err)
	//	return false
	//}
	//
	//return permissionCheck.Result == "true"
	return false

	//
	//fmt.Printf("\nChecking permission %s\n", permission)
	//
	//for _, node := range CachedData[player.ID()].Nodes {
	//	if node.Key == permission {
	//		return node.Value
	//	}
	//}
	//
	//// check for wildcards
	//if permission == "*" {
	//	return false
	//}
	//
	//permission = strings.ReplaceAll(permission, ".*", "")
	//lastIndex := strings.LastIndex(permission, ".")
	//if lastIndex == -1 {
	//	return HasPermission(source, "*")
	//}
	//
	//return HasPermission(source, permission[:lastIndex]+".*")
}
