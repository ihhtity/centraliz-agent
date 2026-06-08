package huifu

import (
	"centraliz-backend/pkg/response"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/huifurepo/bspay-go-sdk/BsPaySdk"
	"github.com/huifurepo/bspay-go-sdk/ut/tool"
)

func WxMiniPay(c *gin.Context) {
	// 1. 数据初始化
	dgSDK, _ := BsPaySdk.NewBsPay(true, "./config/huifu.json")

	// 2.组装请求参数
	dgReq := BsPaySdk.V3TradePaymentJspayRequest{
		// 请求日期
		ReqDate: tool.GetCurrentDate(),
		// 请求流水号
		ReqSeqId: tool.GetReqSeqId(),
		// 商户号
		HuifuId: "6666000153390803",
		// 商品描述
		GoodsDesc: "微信小程序支付测试",
		// 交易金额
		TransAmt: "0.13",
		// 交易类型
		TradeType: "T_MINIAPP",
	}
	// 设置非必填字段
	dgReq.ExtendInfos = GetminiExtendInfos()

	_, _ = json.Marshal(dgReq)

	// 3. 发起API调用
	resp, err := dgSDK.V3TradePaymentJspayRequest(dgReq)
	if err != nil {
		fmt.Println("预下单失败:", err)
		response.Fail(c, 400, "预下单失败", err)
		return
	}

	// 定义 map 接收数据
	var pay_info map[string]interface{}
	// 反序列化
	// 获取data字段（已经是map[string]interface{}）
	dataMap, ok := resp["data"].(map[string]interface{})
	if !ok {
		response.Fail(c, 400, "响应数据格式错误", nil)
		return
	}
	// 获取pay_info并解析
	payInfoStr, ok := dataMap["pay_info"].(string)
	if !ok {
		response.Fail(c, 400, "响应数据格式错误", nil)
		return
	}
	err = json.Unmarshal([]byte(payInfoStr), &pay_info)
	if err != nil {
		response.Fail(c, 400, "解析支付信息失败", err)
		return
	}

	response.Success(c, pay_info)
}

/**
 * 非必填字段
 */
func GetminiExtendInfos() map[string]interface{} {
	// 设置非必填字段
	extendInfoMap := make(map[string]interface{})
	// 是否延迟交易
	extendInfoMap["delay_acct_flag"] = "N"
	// 微信小程序扩展参数集合
	extendInfoMap["wx_data"] = "{\"sub_appid\":\"wxb3699bf4d56598b7\",\"sub_openid\":\"o61IK7RGQ5lxj1M0XmKxqFZH8Pho\",\"detail\":{\"cost_price\":\"43.00\",\"receipt_id\":\"20220628132043853798\",\"goods_detail\":[{\"goods_id\":\"6934572310301\",\"goods_name\":\"太龙双黄连口服液\",\"price\":\"43.00\",\"quantity\":\"1\",\"wxpay_goods_id\":\"12235413214070356458058\"}]}}"
	// 分账对象
	// extendInfoMap["acct_split_bunch"] = "{\"acct_infos\":[{\"div_amt\":\"0.01\",\"huifu_id\":\"6666000109133323\"}]}"
	// 交易异步通知地址
	extendInfoMap["notify_url"] = "https://callback.service.com/xx"
	return extendInfoMap
}
