package services

import (
	"fmt"
	"net/http"

	"github.com/vanspaul/SmartMeterSystem/utils"
)

func MeterHandler(w http.ResponseWriter, r *http.Request) {
	d := r.FormValue("data")
	utils.Logger.Sugar().Debugf("Data: %v", d)
	fmt.Fprintln(w, "Ok")
}
