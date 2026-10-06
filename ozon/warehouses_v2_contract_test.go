package ozon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWarehouseListV2Contract(t *testing.T) {
	server := httptest.NewServer(requestContractHandler(t, http.MethodPost, "/v2/warehouse/list", `{"limit":200,"cursor":"page-1","warehouse_ids":["20605650762000"]}`, `{"cursor":"page-2","has_next":true,"warehouses":[{"warehouse_id":20605650762000,"name":"Warehouse","pause_at":null,"created_at":"2026-10-01T00:00:00Z","updated_at":"2026-10-02T00:00:00Z","postings_limit":-1,"working_days":["MONDAY","FUTURE_DAY"],"address_info":{"address":"address","latitude":55.49,"longitude":38.17,"utc":"UTC+03:00"},"first_mile":{"type":"PICK_UP","timeslot_id":287231,"timeslot_from":"20:59","timeslot_to":"21:00","dropoff_point_id":"1020002075314000"},"timetable":{"timetable_from":"2026-10-01T00:00:00Z","timetable_to":"2026-10-02T00:00:00Z","working_hours":[{"time_from":"2026-10-01T00:00:00Z","time_to":"2026-10-02T00:00:00Z"}]}}]}`))
	defer server.Close()
	response, err := NewClient(WithURI(server.URL)).Warehouses().GetListOfWarehousesV2(context.Background(), &GetListOfWarehousesV2Params{Limit: 200, Cursor: "page-1", WarehouseIds: []string{"20605650762000"}})
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || !response.HasNext || response.Cursor != "page-2" || len(response.Warehouses) != 1 {
		t.Fatalf("response=%+v", response)
	}
	w := response.Warehouses[0]
	if w.WarehouseId != 20605650762000 || w.PostingsLimit != -1 || w.PauseAt != nil || w.CreatedAt.IsZero() || w.UpdatedAt.IsZero() || w.FirstMile.TimeslotId != 287231 || w.AddressInfo.UTC != "UTC+03:00" || w.WorkingDays[1] != "FUTURE_DAY" || len(w.Timetable.WorkingHours) != 1 || w.Timetable.WorkingHours[0].TimeFrom.IsZero() {
		t.Fatalf("warehouse=%+v", w)
	}
}
