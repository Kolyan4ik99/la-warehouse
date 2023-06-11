package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"la-warehouse/internal/model"
	"net/http"
	"testing"
)

func TestReservation(t *testing.T) {
	body := ReqReservation{
		Method: RESERVATION_METHOD,
		Params: []model.ReqReservation{
			{
				WhID: 1,
				Products: []model.ProductsAmount{
					{
						ProductID: 1,
						Amount:    102,
					},
					{
						ProductID: 2,
						Amount:    234,
					},
					{
						ProductID: 3,
						Amount:    12,
					},
					{
						ProductID: 4,
						Amount:    520,
					},
					{
						ProductID: 11,
						Amount:    12,
					},
				},
			},
		},
		Id: 1,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	buffer := bytes.NewBufferString(string(bodyBytes))
	request, err := http.NewRequest("POST", "http://localhost:8012/rpc", buffer)
	if err != nil {
		t.Fatal(err)
	}
	cl := http.Client{}
	resp, err := cl.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.StatusCode)
	bytesArr, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bytesArr))

	var respRes RespReservation
	err = json.Unmarshal(bytesArr, &respRes)
	if err != nil {
		t.Fatal(err)
	}

}
