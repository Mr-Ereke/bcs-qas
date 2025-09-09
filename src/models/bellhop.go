package models

const (
	title  = "Price change"
	pushID = 100
	screen = "INSTRUMENT_CARD"
	tab    = "chart"
)

type BellhopTarget struct {
	Screen string `json:"screen"`
	Symbol string `json:"symbol"`
	Tab    string `json:"tab"`
}

type BellhopPushData struct {
	PushId int           `json:"pushId"`
	Target BellhopTarget `json:"target"`
}

type BellhopPayload struct {
	Title      string          `json:"title"`
	Body       string          `json:"body"`
	CustomerID uint            `json:"customer_id"`
	Data       BellhopPushData `json:"data"`
}

type BellhopRequest struct {
	Name    string         `json:"name"`
	Payload BellhopPayload `json:"payload"`
}

func BuildRequest(pushName string, customerId uint, body string, symbol string) *BellhopRequest {
	return &BellhopRequest{
		Name: pushName,
		Payload: BellhopPayload{
			Title:      title,
			Body:       body,
			CustomerID: customerId,
			Data: BellhopPushData{
				PushId: pushID,
				Target: BellhopTarget{
					Screen: screen,
					Symbol: symbol,
					Tab:    tab,
				},
			},
		},
	}
}
