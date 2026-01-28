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

package msg

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/openimsdk/openkf/server/pkg/openim/param/request"
	"github.com/openimsdk/openkf/server/pkg/openim/sdk/auth"
	"github.com/openimsdk/openkf/server/pkg/openim/sdk/constant"
)

// TestAdminSendMsg test admin send msg function
func TestAdminSendMsg(t *testing.T) {
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
	recvID := os.Getenv("OPENIM_TEST_RECV_ID")
	if recvID == "" {
		t.Skip("set OPENIM_TEST_RECV_ID to run this integration test")
	}

	// test case
	testData := []struct {
		sendID           string
		recvID           string
		SenderPlatformID int
		content          string
		contentType      int
		sessionType      int
	}{
		{
			sendID:           adminID,
			recvID:           recvID,
			SenderPlatformID: constant.PLATFORMID_WEB,
			content:          "hello world!",
			contentType:      constant.CONTENT_TYPE_TEXT,
			sessionType:      constant.SESSION_TYPE_SINGLE_CHAT,
		},
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

	// range test case
	for _, data := range testData {
		res, err := AdminSendMsg(&request.MsgInfo{
			SendID:           data.sendID,
			RecvID:           data.recvID,
			SenderPlatformID: data.SenderPlatformID,
			Content:          &request.TextContent{Content: data.content},
			ContentType:      data.contentType,
			SessionType:      data.sessionType,
		},
			op+"-send",
			api,
			adminResp.Data.Token,
		)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%v", res)
	}
}
