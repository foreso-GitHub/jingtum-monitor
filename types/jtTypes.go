package types

//region define jt objects

type JtNode struct {
	Ip          string `json:"ip"`
	Port        string `json:"port"`
	Name        string `json:"name"`
	Online      bool
	Consensus   bool
	BlockNumber int
	//LatestBlock JtBlock
}

type JtNetwork struct {
	Name               string `json:"name"`
	NodeCount          int
	NodeList           []JtNode `json:"nodes"`
	OnlineNodeCount    int
	OnlineRate         float32
	ConsensusNodeCount int
	ConsensusRate      float32
	BlockNumber        int
	LatestBlock        JtBlock
}

//region jt response

type JtResponseJson struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Status  int         `json:"status"`
	Type    string      `json:"type"`
	Error   JtErrorJson `json:"error"`
}

type JtErrorJson struct {
	Count int    `json:"count"`
	Desc  string `json:"desc"`
	Info  string `json:"info"`
}

//endregion

//region blockNumber

type JtGetBlockNumberParam struct {
	Type string `json:"type"`
}

type JtGetBlockNumberResult struct {
	BlockNumber int `json:"result"`
}

type JtGetBlockNumberResponse struct {
	JtResponseJson
	Results []JtGetBlockNumberResult `json:"result"`
}

//endregion

//region block

type JtGetBlockParam struct {
	Number string `json:"number"`
	Full   bool   `json:"full"`
	//Ledger	string			`json:"ledger"`
}

type JtBlock struct {
	Accepted              bool     `json:"accepted"`
	Account_hash          string   `json:"account_hash"`
	Close_flags           int      `json:"close_flags"`
	Close_time            int      `json:"close_time"`
	Close_time_human      string   `json:"close_time_human"`
	Close_time_resolution int      `json:"close_time_resolution"`
	Closed                bool     `json:"closed"`
	Hash                  string   `json:"hash"`
	Ledger_hash           string   `json:"ledger_hash"`
	Ledger_index          string   `json:"ledger_index"`
	Parent_close_time     int      `json:"parent_close_time"`
	Parent_hash           string   `json:"parent_hash"`
	SeqNum                string   `json:"seqNum"`
	TotalCoins            string   `json:"totalCoins"`
	Total_coins           string   `json:"total_coins"`
	Transaction_hash      string   `json:"transaction_hash"`
	Transactions          []string `json:"transactions"`
}

type JtGetBlockResult struct {
	Block JtBlock `json:"result"`
}

type JtGetBlockResponse struct {
	JtResponseJson
	Results []JtGetBlockResult `json:"result"`
}

//endregion

//endregion

//region tps struct

type JtTps struct {
	Name       string
	Period     int
	BlockCount int
	Blocks     []JtBlock
	TxCount    int
	Tps        float64
}

type JtTpsStatus struct {
	CurrentBlockNumber int
	BlockMap           map[int]JtBlock
	Blocks             []JtBlock
	TotalBlockCount    int
	TotalPeriod        int
	TotalTxCount       int
	TotalTps           float64
	TpsMap             map[int]JtTps //key = 1, 12, 12*60, 12*60*24, 12*60*24*7, total
}

//endregion
