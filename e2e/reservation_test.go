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
	// ---- РЕЗЕРВИРОВАНИЕ ПРОДУКТОВ В СКЛАДЕ ----
	whID := int64(1)
	body := CreateReservationBody(whID, []model.ProductsAmount{
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
	})
	resp := makeRequest(t, body)

	require.Equal(t, resp.StatusCode, http.StatusOK, "status code must be 200")

	bytesArr, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var respRes RespReservation
	err = json.Unmarshal(bytesArr, &respRes)
	require.NoError(t, err)

	amountSum := int64(0)
	params := body.Params[0].Products

	correctProducts := make(map[int64]int64)
	for i, product := range respRes.Result.RespReservation {
		assert.Equal(t, product.ProductID, params[i].ProductID, "product_id should be equal")

		if product.Status == "ok" {
			amountSum += params[i].Amount
			correctProducts[product.ProductID] = params[i].Amount
		}
	}

	// ---- ПРОВЕРЯЕМ, ЧТО КОЛ-ВО ЗАРЕЗЕРВИРОВАННОГО СОВПАДАЕТ С КОЛ-ВОМ УСПЕШНО ЗАРЕЗЕРВИРОВАННЫХ ----
	amountFromReq := RequestGetAmountProducts(t, whID)
	assert.Equal(t, amountFromReq, amountSum, "amount of products should be equal")

	respReservFree := RequestReservationFree(t, []model.ReqReservationFree{
		{
			WhID:     whID,
			Products: []int64{1, 2},
		},
	})

	// ---- ОСВОБОЖДАЕМ НЕСКОЛЬКО ТОВАРОВ ИЗ РЕЗЕРВА ----
	for _, product := range respReservFree.Result.RespReservation {
		if product.Status == "ok" {
			amount, exist := correctProducts[product.ProductID]
			if exist {
				amountSum -= amount
			}
		}
	}

	// ---- ПРОВЕРЯЕМ ЧТО КОЛ-ВО ТОВАРОВ ПОМЕНЯЛОСЬ И ЭТО КОЛ-ВО КОРРЕКТНО ----
	amountFromReq = RequestGetAmountProducts(t, whID)
	assert.Equal(t, amountFromReq, amountSum, "amount of products should be equal")
}

func CreateReservationBody(whID int64, products []model.ProductsAmount) ReqReservation {
	return ReqReservation{
		Method: RESERVATION_METHOD,
		Params: []model.ReqReservation{
			{
				WhID:     whID,
				Products: products,
			},
		},
		Id: 1,
	}
}

func RequestReservationFree(t *testing.T, params []model.ReqReservationFree) RespReservation {
	body := ReqReservationFree{
		Method: RESERVATION_FREE_METHOD,
		Params: params,
		Id:     1,
	}
	resp := makeRequest(t, body)

	bytesArr, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var reservFree RespReservation
	err = json.Unmarshal(bytesArr, &reservFree)
	require.NoError(t, err)

	return reservFree
}

func RequestGetAmountProducts(t *testing.T, whID int64) int64 {
	body := ReqGetAmountProducts{
		Method: GET_AMOUNT_PRODUCTS_METHOD,
		Params: []model.ReqGetAmountProducts{
			{
				WhID: whID,
			},
		},
		Id: 1,
	}
	resp := makeRequest(t, body)

	bytesArr, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var amount RespGetAmountProducts
	err = json.Unmarshal(bytesArr, &amount)
	require.NoError(t, err)

	return amount.Result.Amount
}

func makeRequest(t *testing.T, body interface{}) *http.Response {
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	buffer := bytes.NewBufferString(string(bodyBytes))
	request, err := http.NewRequest("POST", REQUEST_URL, buffer)
	require.NoError(t, err)

	cl := http.Client{}
	resp, err := cl.Do(request)
	require.NoError(t, err)
	return resp
}
