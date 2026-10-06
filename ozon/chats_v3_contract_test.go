package ozon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatListV3WireContract(t *testing.T) {
	body := `{"chats":[{"chat":{"created_at":"2026-10-06T00:00:00Z","chat_id":"chat-1","chat_status":"OPENED","chat_type":"FUTURE_CHAT"},"first_unread_message_id":3000000000118021931,"last_message_id":"30000000001280042740","unread_count":1}],"total_unread_count":5,"cursor":"next","has_next":"true","unknown":{"preserve":true}}`
	server := httptest.NewServer(requestContractHandler(t, http.MethodPost, "/v3/chat/list", `{"filter":{"chat_status":"OPENED","unread_only":true},"limit":100,"cursor":"previous"}`, body))
	defer server.Close()
	response, err := NewClient(WithURI(server.URL)).Chats().ListV3(context.Background(), &ListChatsV3Params{Filter: &ListChatsV3Filter{ChatStatus: "OPENED", UnreadOnly: true}, Limit: 100, Cursor: "previous"})
	if err != nil {
		t.Fatal(err)
	}
	if !response.HasNext || response.Cursor != "next" || len(response.Chats) != 1 || string(response.Chats[0].LastMessageId) != "30000000001280042740" || string(response.Chats[0].FirstUnreadMessageId) != "3000000000118021931" || response.Chats[0].Chat.ChatType != "FUTURE_CHAT" {
		t.Fatalf("response=%+v", response)
	}
	raw := response.RawResponseJSON()
	if string(raw) != body {
		t.Fatalf("raw=%s", raw)
	}
	raw[0] = 'x'
	if string(response.RawResponseJSON()) != body {
		t.Fatal("raw aliases snapshot")
	}
}

func TestCurrentPaginationFlagVariants(t *testing.T) {
	for _, flag := range []string{`true`, `"true"`, `"1"`, `false`, `"false"`, `"0"`} {
		var warehouses GetListOfWarehousesV2Response
		if err := json.Unmarshal([]byte(`{"warehouses":[],"has_next":`+flag+`}`), &warehouses); err != nil {
			t.Fatal(err)
		}
		var chats ListChatsV3Response
		if err := json.Unmarshal([]byte(`{"chats":[],"has_next":`+flag+`}`), &chats); err != nil {
			t.Fatal(err)
		}
		want := flag == `true` || flag == `"true"` || flag == `"1"`
		if warehouses.HasNext != want || chats.HasNext != want {
			t.Fatalf("flag %s: warehouses=%v chats=%v", flag, warehouses.HasNext, chats.HasNext)
		}
	}
	for _, flag := range []string{`"unexpected"`, `123`, `[]`} {
		var chats ListChatsV3Response
		if err := json.Unmarshal([]byte(`{"chats":[],"has_next":`+flag+`}`), &chats); err == nil {
			t.Fatalf("accepted %s", flag)
		}
	}
}
