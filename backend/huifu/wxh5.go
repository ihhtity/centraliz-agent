package huifu

import (
	"centraliz-backend/pkg/response"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/huifurepo/bspay-go-sdk/BsPaySdk"
	"github.com/huifurepo/bspay-go-sdk/ut/tool"
)

func WxH5Pay(c *gin.Context) {
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
		GoodsDesc: "微信H5支付测试",
		// 交易金额
		TransAmt: "0.13",
		// 交易类型
		TradeType: "T_H5",
	}
	// 设置非必填字段
	dgReq.ExtendInfos = Geth5ExtendInfos()

	// fmt.Println("请求时间:", tool.GetCurrentTime())
	_, _ = json.Marshal(dgReq)
	// fmt.Println("请求数据:", string(respStr))

	// 3. 发起API调用
	resp, err := dgSDK.V3TradePaymentJspayRequest(dgReq)
	if err != nil {
		response.Fail(c, 400, "预下单失败", err)
		return
	}

	response.Success(c, resp)
}

/**
 * 非必填字段
 */
func Geth5ExtendInfos() map[string]interface{} {
	// 设置非必填字段
	extendInfoMap := make(map[string]interface{})
	// 交易有效期
	// extendInfoMap["time_expire"] = tool.GetCurrentTime()
	// 微信参数集合
	extendInfoMap["wx_data"] = "{\"detail\":{\"cost_price\":\"43.00\",\"receipt_id\":\"20220628132043853798\",\"goods_detail\":[{\"goods_id\":\"6934572310301\",\"goods_name\":\"太龙双黄连口服液\",\"price\":\"43.00\",\"quantity\":\"1\",\"wxpay_goods_id\":\"12235413214070356458058\",\"device_info\":\"4\"}]}}"
	// 是否延迟交易
	extendInfoMap["delay_acct_flag"] = "N"
	// 分账对象
	// extendInfoMap["acct_split_bunch"] = "{\"acct_infos\":[{\"div_amt\":\"0.10\",\"huifu_id\":\"6666000109133323\"}]}"
	// 传入分账遇到优惠的处理规则
	extendInfoMap["term_div_coupon_type"] = "0"
	// 禁用信用卡标记
	extendInfoMap["limit_pay_type"] = "NO_CREDIT"
	// 场景类型
	extendInfoMap["pay_scene"] = "02"
	// 备注
	extendInfoMap["remark"] = "string"
	// 安全信息
	extendInfoMap["risk_check_data"] = "{\"ip_addr\":\"180.167.105.130\",\"base_station\":\"192.168.1.1\",\"latitude\":\"33.3\",\"longitude\":\"33.3\"}"
	// 设备信息
	extendInfoMap["terminal_device_data"] = "{\"device_type\":\"1\",\"device_ip\":\"10.10.0.1\",\"device_gps\":\"192.168.0.0\",\"devs_id\":\"SPINTP357338300264411\"}"
	// 异步通知地址
	extendInfoMap["notify_url"] = "http://www.baidu.com"
	return extendInfoMap
}
