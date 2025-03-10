// Copyright 2020 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/casbin/casnode/casdoor"
	"github.com/casbin/casnode/util"
	"github.com/casdoor/casdoor-go-sdk/auth"
)

func GetRankingRich() ([]*auth.User, error) {
	return casdoor.GetSortedUsers("score", 25), nil
}

func GetUser(id string) *auth.User {
	user := casdoor.GetUser(id)
	return user
}

func GetUsers() []*auth.User {
	users := casdoor.GetUsers()
	return users
}

func GetMemberNum() int {
	return casdoor.GetUserCount()
}

func UpdateMemberEditorType(user *auth.User, editorType string) (bool, error) {
	if user == nil {
		return false, fmt.Errorf("user is nil")
	}

	SetUserField(user, "editorType", editorType)
	return auth.UpdateUserForColumns(user, []string{"properties"})
}

func GetMemberEditorType(user *auth.User) string {
	return GetUserField(user, "editorType")
}

func UpdateMemberLanguage(user *auth.User, language string) (bool, error) {
	SetUserField(user, "language", language)
	return auth.UpdateUserForColumns(user, []string{"properties"})
}

func GetMemberLanguage(user *auth.User) string {
	return GetUserField(user, "language")
}

// GetMemberEmailReminder return member's email reminder status, and his email address.
func GetMemberEmailReminder(id string) (bool, string) {
	user := GetUser(id)
	if user == nil {
		return false, ""
	}

	return true, user.Email
}

func GetUserByEmail(email string) *auth.User {
	return casdoor.GetUserByEmail(email)
}

func GetMemberCheckinDate(user *auth.User) string {
	return GetUserField(user, "checkinDate")
}

func UpdateMemberCheckinDate(user *auth.User, checkinDate string) (bool, error) {
	SetUserField(user, "checkinDate", checkinDate)
	return auth.UpdateUserForColumns(user, []string{"properties"})
}

func GetUserName(user *auth.User) string {
	if user == nil {
		return ""
	}

	return user.Name
}

func CheckIsAdmin(user *auth.User) bool {
	if user == nil {
		return false
	}

	return user.IsAdmin
}

func GetMemberFileQuota(user *auth.User) int {
	if user == nil {
		return 0
	}

	return GetUserFieldInt(user, "fileQuota")
}

// UpdateMemberOnlineStatus updates member's online information.
func UpdateMemberOnlineStatus(user *auth.User, isOnline bool, lastActionDate string) (bool, error) {
	if user == nil {
		return false, fmt.Errorf("user is nil")
	}

	user.IsOnline = isOnline
	SetUserField(user, "lastActionDate", lastActionDate)
	return auth.UpdateUserForColumns(user, []string{"isOnline", "properties"})
}

func GetOnlineUserCount() int {
	return casdoor.GetOnlineUserCount()
}

type UpdateListItem struct {
	Table     string
	Attribute string
}

func AddMemberByNameAndEmailIfNotExist(username, email string) (*auth.User, error) {
	username = strings.ReplaceAll(username, " ", "")
	if username == "" {
		return nil, fmt.Errorf("username is empty")
	}

	email = strings.ReplaceAll(email, " ", "")
	if email == "" {
		return nil, fmt.Errorf("email is empty")
	}

	user, err := auth.GetUser(username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	username = strings.Split(email, "@")[0]
	user, err = auth.GetUser(username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	newUser := GetUserByEmail(email)
	if newUser == nil {
		properties := map[string]string{}
		properties["emailVerifiedTime"] = util.GetCurrentTime()
		properties["fileQuota"] = strconv.Itoa(DefaultUploadFileQuota)
		properties["renameQuota"] = strconv.Itoa(DefaultRenameQuota)

		newUser = &auth.User{
			Name:              username,
			CreatedTime:       util.GetCurrentTime(),
			UpdatedTime:       util.GetCurrentTime(),
			Id:                "",
			Type:              "",
			Password:          "",
			DisplayName:       "",
			Avatar:            "",
			Email:             email,
			Phone:             "",
			Location:          "",
			Address:           nil,
			Affiliation:       "",
			Title:             "",
			Homepage:          "",
			Tag:               "",
			Score:             getInitScore(),
			Ranking:           GetMemberNum() + 1,
			IsOnline:          false,
			IsAdmin:           false,
			IsGlobalAdmin:     false,
			IsForbidden:       false,
			SignupApplication: CasdoorApplication,
			Properties:        properties,
		}

		_, err = auth.AddUser(newUser)
		if err != nil {
			return newUser, err
		}
	}

	return newUser, nil
}
