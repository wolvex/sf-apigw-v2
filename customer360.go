package apigw

import (
	"github.com/valyala/fastjson"
)

type Customer struct {
	Gw   *Client
	Data *fastjson.Value
	MDN  string
}

func NewCustomer(api *Client, mdn string) (s *Customer, err error) {
	s = &Customer{
		Gw:  api,
		MDN: NormalizeMDN(mdn),
	}

	return
}

func (s *Customer) GetICCID() (iccid string, err error) {
	req := map[string]string{
		"mdn": s.MDN,
	}

	var res []byte
	if res, err = s.Gw.Get("/customer360/v1/subscriber/iccid", "", req); err != nil {
		return
	}

	var parser fastjson.Parser
	var data *fastjson.Value
	if data, err = parser.ParseBytes(res); err != nil {
		return
	}

	iccid = string(data.GetStringBytes("data", "iccid"))

	return
}
