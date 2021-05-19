package apigw

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fastjson"
)

type SubscriberBalance struct {
	AcctResID     int
	AcctResName   string
	BalType       int
	Balance       int
	EffectiveDate time.Time
	ExpiryDate    time.Time
}

type SubscriberService struct {
	ServiceCode   string
	ServiceName   string
	EffectiveDate time.Time
	ExpiryDate    time.Time
}

type Subscriber struct {
	Gw   *Client
	Data *fastjson.Value
	MDN  string
}

func NewSubscriber(api *Client, mdn string) (s *Subscriber, err error) {
	s = &Subscriber{
		Gw:  api,
		MDN: NormalizeMDN(mdn),
	}

	err = s.Query()

	return
}

func (s *Subscriber) Query() (err error) {
	req := map[string]string{
		"mdn": s.MDN,
	}

	var res []byte
	if res, err = s.Gw.Get("/crm/subscriber/query", "v1.0", req); err != nil {
		return
	}

	var parser fastjson.Parser
	s.Data, err = parser.ParseBytes(res)

	return
}

func (s *Subscriber) AddService(serviceCode string) (err error) {
	data := map[string]string{
		"mdn":         s.MDN,
		"serviceCode": serviceCode,
	}

	var req []byte
	if req, err = json.Marshal(data); err != nil {
		return
	}

	var res []byte
	if res, err = s.Gw.Post("/crm/service/buy", "v1.0", req); err != nil {
		return
	}

	var response AddServiceMessage
	if err = json.Unmarshal(res, &response); err != nil {
		return
	}

	if response.TransactionID == "" || response.TransactionID == "map[-nil:true]" {
		err = fmt.Errorf("Failure adding service: %#v", response)
	}

	return
}

func (s *Subscriber) Status() string {
	return string(s.Data.GetStringBytes("state"))
}

func (s *Subscriber) IMSI() string {
	return string(s.Data.GetStringBytes("imsi"))
}

func (s *Subscriber) PUK1() string {
	return string(s.Data.GetStringBytes("puk1"))
}

func (s *Subscriber) PUK2() string {
	return string(s.Data.GetStringBytes("puk2"))
}

func (s *Subscriber) MarketingCategory() string {
	return string(s.Data.GetStringBytes("marketingCategory"))
}

func (s *Subscriber) FraudLocked() string {
	return string(s.Data.GetStringBytes("fraudLocked"))
}

func (s *Subscriber) AccountNumber() string {
	return string(s.Data.GetStringBytes("acctNbr"))
}

func (s *Subscriber) ActiveDate() time.Time {
	return ToTime(string(s.Data.GetStringBytes("activeDate")))
}

func (s *Subscriber) ActiveEndDate() time.Time {
	return ToTime(string(s.Data.GetStringBytes("activeEndDate")))
}

func (s *Subscriber) TerminationDate() time.Time {
	return ToTime(string(s.Data.GetStringBytes("terminationDate")))
}

func (s *Subscriber) NextStateDate() time.Time {
	return ToTime(string(s.Data.GetStringBytes("nextStateDate")))
}

func (s *Subscriber) DueDate() time.Time {
	return ToTime(string(s.Data.GetStringBytes("dueDate")))
}

func (s *Subscriber) LastPaymentDate() time.Time {
	return ToTime(string(s.Data.GetStringBytes("lastPaymentDate")))
}

func (s *Subscriber) SettlementMethod() string {
	return string(s.Data.GetStringBytes("settlementMethod"))
}

func (s *Subscriber) ICCID() string {
	return string(s.Data.GetStringBytes("iccid"))
}

func (s *Subscriber) BirthPlace() string {
	return string(s.Data.GetStringBytes("birthPlace"))
}

func (s *Subscriber) BirthDay() string {
	return string(s.Data.GetStringBytes("birthday"))
}

func (s *Subscriber) CustomerName() string {
	return string(s.Data.GetStringBytes("customerName"))
}

func (s *Subscriber) CustomerType() string {
	return string(s.Data.GetStringBytes("customerType"))
}

func (s *Subscriber) DefaultPricePlan() string {
	return string(s.Data.GetStringBytes("defaultPricePlan"))
}

func (s *Subscriber) DefaultPricePlanCode() string {
	return string(s.Data.GetStringBytes("defaultPricePlanCode"))
}

func (s *Subscriber) DocNumber() string {
	return string(s.Data.GetStringBytes("docNumber"))
}

func (s *Subscriber) DocType() string {
	return string(s.Data.GetStringBytes("docType"))
}

func (s *Subscriber) DocAddress() string {
	return string(s.Data.GetStringBytes("docAddress"))
}

func (s *Subscriber) Email() string {
	return string(s.Data.GetStringBytes("email"))
}

func (s *Subscriber) Gender() string {
	return string(s.Data.GetStringBytes("gender"))
}

func (s *Subscriber) MotherMaidenName() string {
	return string(s.Data.GetStringBytes("motherMaidenName"))
}

func (s *Subscriber) NextState() string {
	return string(s.Data.GetStringBytes("nextState"))
}

func (s *Subscriber) OfferId() int {
	return s.Data.GetInt("offerId")
}

func (s *Subscriber) OfferName() string {
	return string(s.Data.GetStringBytes("offerName"))
}

func (s *Subscriber) ProductCode() string {
	return string(s.Data.GetStringBytes("productCode"))
}

func (s *Subscriber) ProductName() string {
	return string(s.Data.GetStringBytes("productName"))
}

func (s *Subscriber) CustomerGrade() string {
	return string(s.Data.GetStringBytes("customerGrade"))
}

func (s *Subscriber) TotalCreditLimit() int {
	if !s.Data.Exists("totalCreditLimit") {
		return 0
	}
	return s.Data.GetInt("totalCreditLimit") / 100
}

func (s *Subscriber) Balance() int {
	if !s.Data.Exists("balance") {
		for _, balance := range s.Balances() {
			if balance.AcctResID == 1 {
				s.Data.Set("balance", fastjson.MustParse(fmt.Sprintf("%d", balance.Balance)))
				break
			}
		}

	}
	return s.Data.GetInt("balance") / 100
}

func (s *Subscriber) BonusBalance() int {
	if !s.Data.Exists("bonusPulsa") {
		total := 0
		for _, balance := range s.Balances() {
			switch balance.AcctResID {
			case 48, 69, 110:
				total += balance.Balance
			}
		}
		s.Data.Set("bonusPulsa", fastjson.MustParse(fmt.Sprintf("%d", total)))
	}
	return s.Data.GetInt("bonusPulsa") / 100
}

func (s *Subscriber) RemainingCreditLimit() int {
	if !s.Data.Exists("remainingCreditLimit") {
		return 0
	}
	return s.Data.GetInt("remainingCreditLimit") / 100
}

func (s *Subscriber) CurrentUsage() int {
	if !s.Data.Exists("currentUsage") {
		return 0
	}
	return s.Data.GetInt("currentUsage") / 100
}

func (s *Subscriber) Balances() []*SubscriberBalance {
	var balances []*SubscriberBalance
	for _, b := range s.Data.GetArray("balances") {
		balance := &SubscriberBalance{
			AcctResID:     b.GetInt("acctResID"),
			AcctResName:   string(b.GetStringBytes("acctResName")),
			BalType:       b.GetInt("balType"),
			Balance:       b.GetInt("balance"),
			EffectiveDate: ToTime(string(b.GetStringBytes("effDate"))),
			ExpiryDate:    ToTime(string(b.GetStringBytes("expDate"))),
		}
		balances = append(balances, balance)
	}
	return balances
}

func (s *Subscriber) Services() []*SubscriberService {
	var services []*SubscriberService
	for _, b := range s.Data.GetArray("services") {
		service := &SubscriberService{
			ServiceCode:   string(b.GetStringBytes("serviceCode")),
			ServiceName:   string(b.GetStringBytes("serviceName")),
			EffectiveDate: ToTime(string(b.GetStringBytes("effDate"))),
			ExpiryDate:    ToTime(string(b.GetStringBytes("expDate"))),
		}
		services = append(services, service)
	}
	return services
}
