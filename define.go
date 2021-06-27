package apigw

import (
	"strings"
	"time"
)

type AddServiceMessage struct {
	MDN           string `json:"mdn,omitempty"`
	ServiceCode   string `json:"serviceCode,omitempty"`
	ServiceName   string `json:"serviceName,omitempty"`
	EffectiveDate string `json:"effectiveDate,omitempty"`
	ExpiryDate    string `json:"expireDate,omitempty"`
	TransactionID string `json:"transactionId,omitempty"`
	ReturnCode    string `json:"returnCode,omitempty"`
	ResultMsg     string `json:"resultMsg,omitempty"`
}

func NormalizeMDN(mdn string) string {
	if strings.HasPrefix(mdn, "6288") {
		return mdn
	}

	if strings.HasPrefix(mdn, "+62") {
		//remove leading +62
		mdn = strings.Replace(mdn, "+62", "62", 1)
	} else if strings.HasPrefix(mdn, "088") {
		//remove leading 0
		mdn = strings.Replace(mdn, "088", "6288", 1)
	} else if strings.HasPrefix(mdn, "88") {
		//remove leading 0
		mdn = strings.Replace(mdn, "88", "6288", 1)
	}
	return mdn
}

func ToTime(timestamp string) (t time.Time) {
	var err error

	if timestamp == "" {
		timestamp = "01/01/1900 00:00:00"
	}

	if t, err = time.Parse("02/01/2006 15:04:05", timestamp); err == nil {
		return
	}

	if t, err = time.Parse("2006-01-02 15:04:05", timestamp); err == nil {
		return
	}

	t, _ = time.Parse("2006-01-02 15:04:05", "1900-01-01 00:00:00")

	return
}

func ToDate(timestamp string) (t time.Time) {
	var err error

	if timestamp == "" {
		timestamp = "01/01/1900 00:00:00"
	}

	if t, err = time.Parse("02/01/2006 00:00:00", timestamp); err == nil {
		return
	}

	if t, err = time.Parse("2006-01-02 00:00:00", timestamp); err == nil {
		return
	}

	t, _ = time.Parse("2006-01-02 15:04:05", "1900-01-01 00:00:00")

	return
}
