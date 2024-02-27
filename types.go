package krakenapi

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"time"
)

// actions constants
const (
	BUY    = "b"
	SELL   = "s"
	MARKET = "m"
	LIMIT  = "l"
)

// KrakenResponse wraps the Kraken API JSON response
type KrakenResponse struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

// TimeResponse represents the server's time
type TimeResponse struct {
	// Unix timestamp
	Unixtime int64
	// RFC 1123 time format
	Rfc1123 string
}

// AssetPairsResponse includes asset pair informations
type AssetPairsResponse map[string]AssetPairInfo

// AssetPairInfo represents asset pair information
type AssetPairInfo struct {
	Altname            string      `json:"altname"`              // Alternate pair name
	AssetClassBase     string      `json:"aclass_base"`          // Asset class of base component
	Base               string      `json:"base"`                 // Asset ID of base component
	AssetClassQuote    string      `json:"aclass_quote"`         // Asset class of quote component
	Quote              string      `json:"quote"`                // Asset ID of quote component
	PairDecimals       int         `json:"pair_decimals"`        // Scaling decimal places for pair
	CostDecimals       int         `json:"cost_decimals"`        // Scaling decimal places for cost
	LotDecimals        int         `json:"lot_decimals"`         // Scaling decimal places for volume
	LotMultiplier      int         `json:"lot_multiplier"`       // Amount to multiply lot volume by to get currency volume
	LeverageBuy        []float64   `json:"leverage_buy"`         // Array of leverage amounts available when buying
	LeverageSell       []float64   `json:"leverage_sell"`        // Array of leverage amounts available when selling
	Fees               [][]float64 `json:"fees"`                 // Fee schedule array in [<volume>, <percent fee>] tuples
	FeesMaker          [][]float64 `json:"fees_maker"`           // Maker fee schedule array in [<volume>, <percent fee>] tuples (if on maker/taker)
	FeeVolumeCurrency  string      `json:"fee_volume_currency"`  // Volume discount currency
	MarginCall         int         `json:"margin_call"`          // Margin call level
	MarginStop         int         `json:"margin_stop"`          // Stop-out/liquidation margin level
	OrderMin           string      `json:"ordermin"`             // Minimum order size (in terms of base currency)
	CostMin            string      `json:"costmin"`              // Minimum order cost (in terms of quote currency)
	TickSize           string      `json:"tick_size"`            // Minimum increment between valid price levels
	Status             string      `json:"status"`               // Status of asset. Possible values: online, cancel_only, post_only, limit_only, reduce_only.
	LongPositionLimit  int         `json:"long_position_limit"`  // Maximum long margin position size (in terms of base currency)
	ShortPositionLimit int         `json:"short_position_limit"` // Maximum short margin position size (in terms of base currency)
}

// AssetsResponse includes asset informations
type AssetsResponse map[string]AssetInfo

// AssetInfo represents an asset information
type AssetInfo struct {
	// Alternate name
	Altname string
	// Asset class
	AssetClass string `json:"aclass"`
	// Scaling decimal places for record keeping
	Decimals int
	// Scaling decimal places for output display
	DisplayDecimals int `json:"display_decimals"`
}

// BalanceResponse represents the account's balances (list of currencies)
type BalanceResponse map[string]string

// TradeBalanceResponse struct used as the response for the TradeBalance method
type TradeBalanceResponse struct {
	EquivalentBalance         float64 `json:"eb,string"`
	TradeBalance              float64 `json:"tb,string"`
	MarginOP                  float64 `json:"m,string"`
	UnrealizedNetProfitLossOP float64 `json:"n,string"`
	CostBasisOP               float64 `json:"c,string"`
	CurrentValuationOP        float64 `json:"v,string"`
	Equity                    float64 `json:"e,string"`
	FreeMargin                float64 `json:"mf,string"`
	MarginLevel               float64 `json:"ml,string"`
}

// Fees includes fees information for different currencies
type Fees map[string]FeeInfo

// FeeInfo represents a fee information
type FeeInfo struct {
	Fee        float64 `json:"fee,string"`
	MinFee     float64 `json:"minfee,string"`
	MaxFee     float64 `json:"maxfee,string"`
	NextFee    float64 `json:"nextfee,string"`
	NextVolume float64 `json:"nextvolume,string"`
	TierVolume float64 `json:"tiervolume,string"`
}

// TradeVolumeResponse represents the response for trade volume
type TradeVolumeResponse struct {
	Volume    float64 `json:"volume,string"`
	Currency  string  `json:"currency"`
	Fees      Fees    `json:"fees"`
	FeesMaker Fees    `json:"fees_maker"`
}

// TickerResponse includes the requested ticker pairs
type TickerResponse map[string]PairTickerInfo

// DepositAddressesResponse is the response type of a DepositAddresses query to the Kraken API.
type DepositAddressesResponse []struct {
	Address  string `json:"address"`
	Expiretm string `json:"expiretm"`
	New      bool   `json:"new,omitempty"`
}

// WithdrawResponse is the response type of a Withdraw query to the Kraken API.
type WithdrawResponse struct {
	RefID string `json:"refid"`
}

// WithdrawInfoResponse is the response type showing withdrawal information for a selected withdrawal method.
type WithdrawInfoResponse struct {
	Method string    `json:"method"`
	Limit  big.Float `json:"limit"`
	Amount big.Float `json:"amount"`
	Fee    big.Float `json:"fee"`
}

// GetPairTickerInfo is a helper method that returns given `pair`'s `PairTickerInfo`
func (v *TickerResponse) GetPairTickerInfo(pair string) PairTickerInfo {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(pair)

	return f.Interface().(PairTickerInfo)
}

// PairTickerInfo represents ticker information for a pair
type PairTickerInfo struct {
	// Ask array(<price>, <whole lot volume>, <lot volume>)
	Ask []string `json:"a"`
	// Bid array(<price>, <whole lot volume>, <lot volume>)
	Bid []string `json:"b"`
	// Last trade closed array(<price>, <lot volume>)
	Close []string `json:"c"`
	// Volume array(<today>, <last 24 hours>)
	Volume []string `json:"v"`
	// Volume weighted average price array(<today>, <last 24 hours>)
	VolumeAveragePrice []string `json:"p"`
	// Number of trades array(<today>, <last 24 hours>)
	Trades []int `json:"t"`
	// Low array(<today>, <last 24 hours>)
	Low []string `json:"l"`
	// High array(<today>, <last 24 hours>)
	High []string `json:"h"`
	// Today's opening price
	OpeningPrice float64 `json:"o,string"`
}

// TradesResponse represents a list of the last trades
type TradesResponse struct {
	Last   int64
	Trades []TradeInfo
}

// TradesHistoryResponse represents a list of executed trade
type TradesHistoryResponse struct {
	Trades map[string]TradeHistoryInfo `json:"trades"`
	Count  int                         `json:"count"`
}

// TradeHistoryInfo represents a transaction
type TradeHistoryInfo struct {
	TransactionID string  `json:"ordertxid"`
	PostxID       string  `json:"postxid"`
	AssetPair     string  `json:"pair"`
	Time          float64 `json:"time"`
	Type          string  `json:"type"`
	OrderType     string  `json:"ordertype"`
	Price         float64 `json:"price,string"`
	Cost          float64 `json:"cost,string"`
	Fee           float64 `json:"fee,string"`
	Volume        float64 `json:"vol,string"`
	Margin        float64 `json:"margin,string"`
	Misc          string  `json:"misc"`
}

// TradeInfo represents a trades information
type TradeInfo struct {
	Price         string
	PriceFloat    float64
	Volume        string
	VolumeFloat   float64
	Time          int64
	Buy           bool
	Sell          bool
	Market        bool
	Limit         bool
	Miscellaneous string
}

// LedgersResponse represents an associative array of ledgers infos
type LedgersResponse struct {
	Ledger map[string]LedgerInfo `json:"ledger"`
}

// LedgerInfo Represents the ledger informations
type LedgerInfo struct {
	RefID   string    `json:"refid"`
	Time    float64   `json:"time"`
	Type    string    `json:"type"`
	Aclass  string    `json:"aclass"`
	Asset   string    `json:"asset"`
	Amount  big.Float `json:"amount"`
	Fee     big.Float `json:"fee"`
	Balance big.Float `json:"balance"`
}

// OrderTypes for AddOrder
const (
	OTMarket              = "market"
	OTLimit               = "limit"                  // (price = limit price)
	OTStopLoss            = "stop-loss"              // (price = stop loss price)
	OTTakeProfi           = "take-profit"            // (price = take profit price)
	OTStopLossProfit      = "stop-loss-profit"       // (price = stop loss price, price2 = take profit price)
	OTStopLossProfitLimit = "stop-loss-profit-limit" // (price = stop loss price, price2 = take profit price)
	OTStopLossLimit       = "stop-loss-limit"        // (price = stop loss trigger price, price2 = triggered limit price)
	OTTakeProfitLimit     = "take-profit-limit"      // (price = take profit trigger price, price2 = triggered limit price)
	OTTrailingStop        = "trailing-stop"          // (price = trailing stop offset)
	OTTrailingStopLimit   = "trailing-stop-limit"    // (price = trailing stop offset, price2 = triggered limit offset)
	OTStopLossAndLimit    = "stop-loss-and-limit"    // (price = stop loss price, price2 = limit price)
	OTSettlePosition      = "settle-position"
)

// OrderDescription represents an order description
type OrderDescription struct {
	Pair      string  `json:"pair"`          // Asset pair
	Type      string  `json:"type"`          // "buy" or "sell"
	OrderType string  `json:"ordertype"`     // "market" or "limit" or "stop-loss" or "take-profit" or "stop-loss-limit" or "take-profit-limit" or "trailing-stop" or "trailing-stop-limit" or "settle-position"
	Price     float64 `json:"price,string"`  // Limit price for "limit" orders. Trigger price for "stop-loss", "stop-loss-limit", "take-profit", "take-profit-limit", "trailing-stop" and "trailing-stop-limit orders"
	Price2    float64 `json:"price2,string"` // Limit price for "stop-loss-limit", "take-profit-limit" and "trailing-stop-limit orders"
	Leverage  string  `json:"leverage"`      // Amount of leverage
	Order     string  `json:"order"`         // Order description
	Close     string  `json:"close"`         // Conditional close order description (if conditional close set)
}

// Order represents a single order
type Order struct {
	ReferenceID    string           `json:"refid"`             // Referral order transaction ID that created this order
	UserRef        int              `json:"userref"`           // User reference id
	Status         string           `json:"status"`            // "pending" or "open" or "closed" or "canceled" or "expired"
	OpenTime       float64          `json:"opentm"`            // Unix timestamp of when order was placed
	StartTime      float64          `json:"starttm"`           // Unix timestamp of order start time (or 0 if not set)
	ExpireTime     float64          `json:"expiretm"`          // Unix timestamp of order end time (or 0 if not set)
	Description    OrderDescription `json:"descr"`             // Order description info
	Volume         float64          `json:"vol,string"`        // Volume of order (base currency)
	VolumeExecuted float64          `json:"vol_exec,string"`   // Volume executed (base currency)
	Cost           float64          `json:"cost,string"`       // Total cost (quote currency unless)
	Fee            float64          `json:"fee,string"`        // Total fee (quote currency)
	Price          float64          `json:"price,string"`      // Average price (quote currency)
	StopPrice      float64          `json:"stopprice.string"`  // Stop price (quote currency)
	LimitPrice     float64          `json:"limitprice,string"` // Triggered limit price (quote currency, when limit based order type triggered)
	Misc           string           `json:"misc"`              // Comma delimited list of miscellaneous info
	OrderFlags     string           `json:"oflags"`            // Comma delimited list of order flags
}

// ClosedOrdersResponse represents a list of closed orders, indexed by id
type ClosedOrdersResponse struct {
	Closed map[string]Order `json:"closed"`
	Count  int              `json:"count"`
}

// OrderBookItem is a piece of information about an order.
type OrderBookItem struct {
	Price  float64
	Amount float64
	Ts     int64
}

// UnmarshalJSON takes a json array from kraken and converts it into an OrderBookItem.
func (o *OrderBookItem) UnmarshalJSON(data []byte) error {
	tmpStruct := struct {
		price  string
		amount string
		ts     int64
	}{}
	tmpArray := []interface{}{&tmpStruct.price, &tmpStruct.amount, &tmpStruct.ts}
	err := json.Unmarshal(data, &tmpArray)
	if err != nil {
		return err
	}

	o.Price, err = strconv.ParseFloat(tmpStruct.price, 64)
	if err != nil {
		return err
	}
	o.Amount, err = strconv.ParseFloat(tmpStruct.amount, 64)
	if err != nil {
		return err
	}
	o.Ts = tmpStruct.ts
	return nil
}

// DepthResponse is a response from kraken to Depth request.
type DepthResponse map[string]OrderBook

// OrderBook contains top asks and bids.
type OrderBook struct {
	Asks []OrderBookItem
	Bids []OrderBookItem
}

// OpenOrdersResponse response when opening an order
type OpenOrdersResponse struct {
	Open map[string]Order `json:"open"`
}

// AddOrderResponse response when adding an order
type AddOrderResponse struct {
	Description struct {
		Order string `json:"order"`
	} `json:"descr"`
	TxId []string `json:"txid"`
}

// CancelOrderResponse response when cancelling and order
type CancelOrderResponse struct {
	Count   int  `json:"count"`
	Pending bool `json:"pending"`
}

// QueryOrdersResponse response when checking all orders
type QueryOrdersResponse map[string]Order

// NewOHLC constructor for OHLC
func NewOHLC(input []interface{}) (*OHLC, error) {
	if len(input) != 8 {
		return nil, fmt.Errorf("the length is not 8 but %d", len(input))
	}

	tmp := new(OHLC)
	tmp.Time = time.Unix(int64(input[0].(float64)), 0)
	tmp.Open, _ = strconv.ParseFloat(input[1].(string), 64)
	tmp.High, _ = strconv.ParseFloat(input[2].(string), 64)
	tmp.Low, _ = strconv.ParseFloat(input[3].(string), 64)
	tmp.Close, _ = strconv.ParseFloat(input[4].(string), 64)
	tmp.Vwap, _ = strconv.ParseFloat(input[5].(string), 64)
	tmp.Volume, _ = strconv.ParseFloat(input[6].(string), 64)
	tmp.Count = int(input[7].(float64))

	return tmp, nil
}

// OHLC represents the "Open-high-low-close chart"
type OHLC struct {
	Time   time.Time `json:"time"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Vwap   float64   `json:"vwap"`
	Volume float64   `json:"volume"`
	Count  int       `json:"count"`
}

// OHLCResponse represents the OHLC's response
type OHLCResponse struct {
	Pair string  `json:"pair"`
	OHLC []*OHLC `json:"OHLC"`
	Last float64 `json:"last"`
}
