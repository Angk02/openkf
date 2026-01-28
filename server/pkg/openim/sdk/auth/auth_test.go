// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/openimsdk/openkf/server/pkg/openim/param/request"
	"github.com/openimsdk/openkf/server/pkg/openim/sdk/auth"
)

func TestGetAdminToken(t *testing.T) {
	api := os.Getenv("OPENIM_API_ADDRESS")
	if api == "" {
		t.Skip("set OPENIM_API_ADDRESS (e.g. http://127.0.0.1:10002) to run this integration test")
	}
	secret := os.Getenv("OPENIM_SECRET")
	if secret == "" {
		secret = "openIM123"
	}
	adminID := os.Getenv("OPENIM_ADMIN_ID")
	if adminID == "" {
		adminID = "imAdmin"
	}

	op := strconv.FormatInt(time.Now().UnixMilli(), 10)
	resp, err := auth.GetAdminToken(&request.AdminTokenParams{
		Secret: secret,
		UserID: adminID,
	}, op, api)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ErrCode != 0 {
		t.Fatalf("OpenIM error: errCode=%d errMsg=%s errDlt=%s", resp.ErrCode, resp.ErrMsg, resp.ErrDlt)
	}
	if resp.Data.Token == "" {
		t.Fatalf("empty token")
	}
}

func TestGetUserToken(t *testing.T) {
	api := os.Getenv("OPENIM_API_ADDRESS")
	if api == "" {
		t.Skip("set OPENIM_API_ADDRESS (e.g. http://127.0.0.1:10002) to run this integration test")
	}
	userID := os.Getenv("OPENIM_TEST_USER_ID")
	if userID == "" {
		t.Skip("set OPENIM_TEST_USER_ID to run this integration test")
	}

	secret := os.Getenv("OPENIM_SECRET")
	if secret == "" {
		secret = "openIM123"
	}
	adminID := os.Getenv("OPENIM_ADMIN_ID")
	if adminID == "" {
		adminID = "imAdmin"
	}

	op := strconv.FormatInt(time.Now().UnixMilli(), 10)
	adminResp, err := auth.GetAdminToken(&request.AdminTokenParams{
		Secret: secret,
		UserID: adminID,
	}, op+"-admin", api)
	if err != nil {
		t.Fatal(err)
	}
	if adminResp.ErrCode != 0 {
		t.Fatalf("OpenIM error: errCode=%d errMsg=%s errDlt=%s", adminResp.ErrCode, adminResp.ErrMsg, adminResp.ErrDlt)
	}

	userResp, err := auth.GetUserToken(&request.UserTokenParams{
		PlatformID: 5,
		UserID:     userID,
	}, op+"-user", api, adminResp.Data.Token)
	if err != nil {
		t.Fatal(err)
	}
	if userResp.ErrCode != 0 {
		t.Fatalf("OpenIM error: errCode=%d errMsg=%s errDlt=%s", userResp.ErrCode, userResp.ErrMsg, userResp.ErrDlt)
	}
	if userResp.Data.Token == "" {
		t.Fatalf("empty token")
	}
}
