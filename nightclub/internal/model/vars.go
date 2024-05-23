package model

// HomestayOrder 交易状态 :  -1: 已取消 0:待支付 1:未使用 2:已使用  3:已过期
var SendGiftOrderTradeStateCancel int64 = -1
var SendGiftOrderTradeStateWaitPay int64 = 0
var SendGiftOrderTradeStateWaitUse int64 = 1
var SendGiftOrderTradeStateUsed int64 = 2
var SendGiftOrderTradeStateRefund int64 = 3
var SendGiftOrderTradeStateExpire int64 = 4

// 支付业务类型
var ThirdPaymentServiceTypeSendGiftOrder string = "homestayOrder" //送礼支付

// 支付方式
var ThirdPaymentPayModelWechatPay = "WECHAT_PAY" //微信支付

// 平台内支付状态
var ThirdPaymentPayTradeStateFAIL int64 = -1   //支付失败
var ThirdPaymentPayTradeStateWait int64 = 0    //待支付
var ThirdPaymentPayTradeStateSuccess int64 = 1 //支付成功
var ThirdPaymentPayTradeStateRefund int64 = 2  //已退款
