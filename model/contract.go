package model

type ContractModel struct {
	TemplateName                string   `json:"templateName" form:"templateName"` // 合同模板名称
	ID                          int      `json:"id" form:"id"`                     // 合同ID
	TaskID                      int      `json:"taskId" form:"taskId"`             // 稿单ID
	PartyAName                  string   `json:"partyAName" form:"partyAName"`
	PartyBName                  string   `json:"partyBName" form:"partyBName"`                                   // 甲方姓名
	Amount                      float64  `json:"amount" form:"amount"`                                           // 合同姓名
	UpperAmount                 string   `json:"upperAmount" form:"upperAmount"`                                 // 合同价格（大写）
	Deposit                     float64  `json:"deposit" form:"deposit"`                                         // 支付定金
	UpperDeposit                string   `json:"upperDeposit" form:"upperDeposit"`                               // 支付定金（大写）
	PayDepositWithinDays        int      `json:"payDepositWithinDays" form:"payDepositWithinDays"`               // 几天内支付定金
	Balance                     float64  `json:"balance" form:"balance"`                                         // 支付尾款
	UpperBalance                string   `json:"upperBalance" form:"upperBalance"`                               // 支付尾款（大写）
	PayBalanceWithinDays        int      `json:"payBalanceWithinDays" form:"payBalanceWithinDays"`               // 几天内支付尾款
	DelayedDeliveryPenalty      *float64 `json:"delayedDeliveryPenalty" form:"delayedDeliveryPenalty"`           // 延迟交付违约金
	UpperDelayedDeliveryPenalty *string  `json:"upperDelayedDeliveryPenalty" form:"upperDelayedDeliveryPenalty"` // 延迟交付违约金（大写）
	DelayedPaymentPenalty       *float64 `json:"delayedPaymentPenalty" form:"delayedPaymentPenalty"`             // 延迟支付违约金
	UpperDelayedPaymentPenalty  *string  `json:"upperDelayedPaymentPenalty" form:"upperDelayedPaymentPenalty"`   // 延迟支付违约金（大写）
	ExpiresDate                 string   `json:"expiresDate" form:"expiresDate"`                                 // 到期日期（年-月-日）
	PartyASignature             string   `json:"partyASignature" form:"partyASignature"`                         // 甲方签名
	PartyASignDate              string   `json:"partyASignDate" form:"partyASignDate"`                           // 甲方签署日期
	PartyBSignature             string   `json:"partyBSignature" form:"partyBSignature"`                         // 乙方签名
	PartyBSignDate              string   `json:"partyBSignDate" form:"partyBSignDate"`                           // 乙方签署日期
}
