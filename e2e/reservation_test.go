package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"la-warehouse/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	buffer := bytes.NewBufferString(string(bodyBytes))
	request, err := http.NewRequest("POST", "http://localhost:8012/rpc", buffer)
	require.NoError(t, err)

	cl := http.Client{}
	resp, err := cl.Do(request)
	require.NoError(t, err)

	require.Equal(t, resp.StatusCode, http.StatusOK, "status code must be 200")

	bytesArr, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var respRes RespReservation
	err = json.Unmarshal(bytesArr, &respRes)
	require.NoError(t, err)

	params := body.Params[0].Products
	for _, reserve := range respRes.Result.RespReservation {
		for i, product := range reserve.Products {
			assert.Equal(t, product.ProductID, params[i].ProductID, "product_id should be equal")
		}
	}
}
