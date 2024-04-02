package lknpd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

type (
	CustomerType string
	PaymentType  string
	CancelType   string
)

const (
	Individual    CustomerType = "FROM_INDIVIDUAL"
	LegalEntity   CustomerType = "FROM_LEGAL_ENTITY"
	ForeignAgency CustomerType = "FROM_FOREIGN_AGENCY"

	Cash    PaymentType = "CASH"
	Account PaymentType = "ACCOUNT"

	Cancel CancelType = "Чек сформирован ошибочно"
	Refund CancelType = "Возврат средств"
)

type (
	Service struct {
		Label    string          `json:"name"`
		Amount   decimal.Decimal `json:"amount"`
		Quantity int64           `json:"quantity"`
	}

	Customer struct {
		ContactPhone string       `json:"contactPhone"`
		DisplayName  string       `json:"displayName"`
		IncomeType   CustomerType `json:"incomeType"`
		Inn          string       `json:"inn"`
	}

	CreateSaleRequest struct {
		PaymentType                     PaymentType `json:"paymentType"`
		Client                          *Customer   `json:"client"`
		RequestTime                     time.Time   `json:"requestTime"`
		OperationTime                   time.Time   `json:"operationTime"`
		Services                        []*Service  `json:"services"`
		TotalAmount                     string      `json:"totalAmount"`
		IgnoreMaxTotalIncomeRestriction bool        `json:"ignoreMaxTotalIncomeRestriction"`
	}

	CreateSaleResponse struct {
		ApprovedReceiptUUID string `json:"approvedReceiptUuid"`
	}
)

func (o *Client) CreateSale(sale CreateSaleRequest) (orderId string, err error) {
	o.CheckTokenExpireIn()

	if sale.PaymentType == "" {
		sale.PaymentType = Cash
	}

	totalAmount := decimal.NewFromInt(0)
	for _, service := range sale.Services {
		totalAmount = totalAmount.Add(service.Amount.Mul(decimal.NewFromInt(service.Quantity)))
	}
	sale.TotalAmount = totalAmount.String()

	location, err := time.LoadLocation(o.timezone)
	if err != nil {
		return
	}
	sale.RequestTime = time.Now().In(location)
	if sale.OperationTime.IsZero() {
		sale.OperationTime = sale.RequestTime
	}

	client := resty.New()
	resp, err := client.R().SetBody(sale).
		SetHeader("Content-Type", "application/json").
		SetHeader("Referrer", "https://lknpd.nalog.ru/").
		SetHeader("Referrer-Policy", "strict-origin-when-cross-origin").
		SetHeader("Authorization", "Bearer "+o.accessToken).
		Post("https://lknpd.nalog.ru/api/v1/income")
	if err != nil {
		return
	}

	if resp.StatusCode() == 200 {
		var result CreateSaleResponse
		if err = json.Unmarshal(resp.Body(), &result); err != nil {
			return
		}
		orderId = result.ApprovedReceiptUUID
	} else {
		err = fmt.Errorf("status code: %d. msg: %s", resp.StatusCode(), resp.RawBody())
	}

	return
}

type (
	CancelSaleRequest struct {
		RequestTime   time.Time  `json:"requestTime"`
		OperationTime time.Time  `json:"operationTime"`
		CancelType    CancelType `json:"comment"`
		ReceiptUUID   string     `json:"receiptUuid"`
		PartnerCode   string     `json:"partnerCode"`
	}
)

func (o *Client) CancelSale(sale CancelSaleRequest) (err error) {
	o.CheckTokenExpireIn()

	location, err := time.LoadLocation(o.timezone)
	if err != nil {
		return
	}
	sale.RequestTime = time.Now().In(location)
	if sale.OperationTime.IsZero() {
		sale.OperationTime = sale.RequestTime
	}

	client := resty.New()
	resp, err := client.R().SetBody(sale).
		SetHeader("Content-Type", "application/json").
		SetHeader("Referrer", "https://lknpd.nalog.ru/").
		SetHeader("Referrer-Policy", "strict-origin-when-cross-origin").
		SetHeader("Authorization", "Bearer "+o.accessToken).
		Post("https://lknpd.nalog.ru/api/v1/cancel")
	if err != nil {
		return
	}

	if resp.StatusCode() != 200 {
		err = fmt.Errorf("status code: %d. msg: %+v", resp.StatusCode(), resp.RawBody())
	}

	return
}
