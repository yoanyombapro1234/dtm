package dtmsvr

import (
	"fmt"
	"strings"
	"time"

	"github.com/yedf/dtm/common"
)

func CronPreparedOnce(expire time.Duration) {
	db := DbGet()
	ss := []SagaModel{}
	dbr := db.Model(&SagaModel{}).Where("update_time < date_sub(now(), interval ? second)", int(expire/time.Second)).Where("status = ?", "prepared").Find(&ss)
	common.PanicIfError(dbr.Error)
	writeTransLog("", "saga fetch prepared", fmt.Sprint(len(ss)), -1, "")
	if len(ss) == 0 {
		return
	}
	for _, sm := range ss {
		writeTransLog(sm.Gid, "saga touch prepared", "", -1, "")
		dbr = db.Model(&sm).Update("id", sm.ID)
		common.PanicIfError(dbr.Error)
		resp, err := common.RestyClient.R().SetQueryParam("gid", sm.Gid).Get(sm.TransQuery)
		common.PanicIfError(err)
		body := resp.String()
		if strings.Contains(body, "FAIL") {
			writeTransLog(sm.Gid, "saga canceled", "canceled", -1, "")
			dbr = db.Model(&sm).Where("status = ?", "prepared").Update("status", "canceled")
			common.PanicIfError(dbr.Error)
		} else if strings.Contains(body, "SUCCESS") {
			m := M{}
			steps := []M{}
			common.MustRemarshal(sm, &m)
			common.PanicIfError(err)
			common.MustUnmarshalString(m["steps"].(string), &steps)
			m["steps"] = steps
			err = rabbit.SendAndConfirm(RabbitmqConstCommited, m)
			common.PanicIfError(err)
		}
	}
}

func CronPrepared() {
	for {
		CronPreparedOnce(10 * time.Second)
	}
}

func CronCommitedOnce(expire time.Duration) {
	db := DbGet()
	ss := []SagaModel{}
	dbr := db.Model(&SagaModel{}).Where("update_time < date_sub(now(), interval ? second)", int(expire/time.Second)).Where("status = ?", "commited").Find(&ss)
	common.PanicIfError(dbr.Error)
	writeTransLog("", "saga fetch commited", fmt.Sprint(len(ss)), -1, "")
	if len(ss) == 0 {
		return
	}
	for _, sm := range ss {
		writeTransLog(sm.Gid, "saga touch commited", "", -1, "")
		dbr = db.Model(&sm).Update("id", sm.ID)
		common.PanicIfError(dbr.Error)
		ProcessCommitedSaga(sm.Gid)
	}
}

func CronCommited() {
	for {
		CronCommitedOnce(10 * time.Second)
	}
}