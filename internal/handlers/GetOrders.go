package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Repinoid/ku/internal/models"
	"github.com/Repinoid/ku/internal/securitate"
)

type OrdStruct struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

func GetOrders(rwr http.ResponseWriter, req *http.Request) {

	rwr.Header().Set("Content-Type", "application/json")

	UserID, err := securitate.DataBase.LoginByToken(rwr, req)
	if err != nil {
		return
	}

	db := securitate.DataBase.DB
	order := "select ordernumber as number, orderstatus as status, accrual, uploaded_at from orders where usercode=$1 order by uploaded_at ;"

	rows, err := db.Query(context.Background(), order, UserID) //
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}

	ord := OrdStruct{}
	orda := []OrdStruct{}
	var errScan error
	for rows.Next() {
		var tm time.Time
		errScan = rows.Scan(&ord.Number, &ord.Status, &ord.Accrual, &tm)
		ord.UploadedAt = tm.Format(time.RFC3339)
		if errScan != nil {
			break
		}
		orda = append(orda, ord)
	}
	rows.Close()
	if err := rows.Err(); err != nil || errScan != nil { // Err returns any error that occurred while reading. Err must only be called after the Rows is closed
		rwr.WriteHeader(http.StatusInternalServerError) // //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}
	if len(orda) == 0 {
		rwr.WriteHeader(http.StatusNoContent) // 204 No Content — сервер успешно обработал запрос, но в ответе были переданы только заголовки без тела сообщения
		fmt.Fprintf(rwr, `{"status":"StatusNoContent"}`)
		models.Sugar.Debug("No ORDERS\n")
		return
	}
	rwr.WriteHeader(http.StatusOK)
	models.Sugar.Debugf("orda[0].Status  \"%+v\"\n", orda[0].Status)
	//	fmt.Fprintf(rwr, `{"status":"StatusOK"}`)
	json.NewEncoder(rwr).Encode(orda)
}
