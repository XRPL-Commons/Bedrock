package caller

// CallResult represents the result of a contract call
type CallResult struct {
	TxHash            string                 `json:"txHash"`
	ReturnCode        int                    `json:"returnCode"`
	ReturnValue       string                 `json:"returnValue"`
	GasUsed           int64                  `json:"gasUsed"`
	Validated         bool                   `json:"validated"`
	TransactionResult string                 `json:"transactionResult"`
	Meta              map[string]interface{} `json:"meta"`
}

// CallConfig holds configuration for calling a contract function
type CallConfig struct {
	ContractAccount      string
	FunctionName         string
	NetworkURL           string
	WalletSeed           string
	ABIPath              string
	Parameters           map[string]interface{} // JSON parameters
	ComputationAllowance string
	Fee                  string
}
