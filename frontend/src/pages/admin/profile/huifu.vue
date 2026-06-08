<!-- 汇付天下API测试页面 -->
<template>
	<view class="container">
		<uv-navbar :title="'汇付测试'" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />

		<!-- 导航标签 -->
		<view class="tab-bar">
			<view v-for="tab in tabs" :key="tab.key" :class="['tab-item', { active: activeTab === tab.key }]"
				@click="switchTab(tab.key)">
				<text class="tab-text">{{ tab.name }}</text>
				<view class="tab-indicator" v-if="activeTab === tab.key"></view>
			</view>
		</view>

		<!-- 商户进件测试 -->
		<view v-if="activeTab === 'merchant'" class="content">
			<view class="tab-bars">
				<uv-tabs :list="merchantTabs" @click="switchMerchantTab" />
			</view>
			<view class="card" v-if="activeMerchantTab == 'enterprise'">
				<view class="card-title">企业商户进件</view>
				<view class="form-item">
					<text class="form-label">商户名称 *</text>
					<input class="form-input" v-model="enterpriseForm.regName" placeholder="请输入商户名称" />
				</view>
				<view class="form-item">
					<text class="form-label">商户简称</text>
					<input class="form-input" v-model="enterpriseForm.shortName" placeholder="请输入商户简称" />
				</view>
				<view class="form-item">
					<text class="form-label">小票名称</text>
					<input class="form-input" v-model="enterpriseForm.receiptName" placeholder="请输入小票名称" />
				</view>
				<view class="form-item">
					<text class="form-label">公司类型 *</text>
					<input class="form-input" v-model="enterpriseForm.entType" placeholder="1-企业" />
				</view>
				<view class="form-item">
					<text class="form-label">所属行业MCC *</text>
					<input class="form-input" v-model="enterpriseForm.mcc" placeholder="请输入MCC编码，如5411" />
				</view>
				<view class="form-item">
					<text class="form-label">经营类型 *</text>
					<input class="form-input" v-model="enterpriseForm.busiType" placeholder="1-实体，2-虚拟" />
				</view>
				<view class="form-item">
					<text class="form-label">场景类型 *</text>
					<input class="form-input" v-model="enterpriseForm.sceneType"
						placeholder="如：test, ONLINE, OFFLINE, ALL" />
				</view>
				<view class="form-item">
					<text class="form-label">证照编号 *</text>
					<input class="form-input" v-model="enterpriseForm.licenseCode" placeholder="请输入营业执照编号" />
				</view>
				<view class="form-item">
					<text class="form-label">证照有效期类型 *</text>
					<input class="form-input" v-model="enterpriseForm.licenseValidityType" placeholder="0-固定期限，1-长期" />
				</view>
				<view class="form-item">
					<text class="form-label">证照有效期开始日期 *</text>
					<input class="form-input" v-model="enterpriseForm.licenseBeginDate" placeholder="格式：yyyyMMdd" />
				</view>
				<view class="form-item" v-if="enterpriseForm.licenseValidityType === '0'">
					<text class="form-label">证照有效期截止日期</text>
					<input class="form-input" v-model="enterpriseForm.licenseEndDate" placeholder="格式：yyyyMMdd" />
				</view>
				<view class="form-item">
					<text class="form-label">注册区 *</text>
					<input class="form-input" v-model="enterpriseForm.regDistrictId" placeholder="请输入注册区编码，如350203" />
				</view>
				<view class="form-item">
					<text class="form-label">注册详细地址 *</text>
					<input class="form-input" v-model="enterpriseForm.regDetail" placeholder="请输入注册详细地址" />
				</view>
				<view class="form-item">
					<text class="form-label">经营区 *</text>
					<input class="form-input" v-model="enterpriseForm.districtId" placeholder="请输入经营区编码，如310104" />
				</view>
				<view class="form-item">
					<text class="form-label">经营详细地址 *</text>
					<input class="form-input" v-model="enterpriseForm.detailAddr" placeholder="请输入经营详细地址" />
				</view>
				<view class="form-item">
					<text class="form-label">法人姓名 *</text>
					<input class="form-input" v-model="enterpriseForm.legalName" placeholder="请输入法人姓名" />
				</view>
				<view class="form-item">
					<text class="form-label">法人证件类型 *</text>
					<input class="form-input" v-model="enterpriseForm.legalCertType" placeholder="00-身份证" />
				</view>
				<view class="form-item">
					<text class="form-label">法人证件号码 *</text>
					<input class="form-input" v-model="enterpriseForm.legalCertNo" placeholder="请输入法人身份证号码" />
				</view>
				<view class="form-item">
					<text class="form-label">法人证件有效期类型 *</text>
					<input class="form-input" v-model="enterpriseForm.legalCertValidityType"
						placeholder="0-固定期限，1-长期" />
				</view>
				<view class="form-item">
					<text class="form-label">法人证件有效期开始日期 *</text>
					<input class="form-input" v-model="enterpriseForm.legalCertBeginDate" placeholder="格式：yyyyMMdd" />
				</view>
				<view class="form-item" v-if="enterpriseForm.legalCertValidityType === '0'">
					<text class="form-label">法人证件有效期截止日期</text>
					<input class="form-input" v-model="enterpriseForm.legalCertEndDate" placeholder="格式：yyyyMMdd" />
				</view>
				<view class="form-item">
					<text class="form-label">管理员手机号 *</text>
					<input class="form-input" v-model="enterpriseForm.contactMobileNo" placeholder="请输入管理员手机号" />
				</view>
				<view class="form-item">
					<text class="form-label">管理员邮箱 *</text>
					<input class="form-input" v-model="enterpriseForm.contactEmail" placeholder="请输入管理员邮箱" />
				</view>
				<view class="form-item">
					<text class="form-label">管理员账号 *</text>
					<input class="form-input" v-model="enterpriseForm.loginName" placeholder="请输入管理员账号" />
				</view>
				<view class="form-item">
					<text class="form-label">开户名称 *</text>
					<input class="form-input" v-model="enterpriseForm.cardName" placeholder="对公账户开户名称" />
				</view>
				<view class="form-item">
					<text class="form-label">开户账号 *</text>
					<input class="form-input" v-model="enterpriseForm.cardNo" placeholder="对公账户账号" />
				</view>
				<view class="form-item">
					<text class="form-label">开户行联行号 *</text>
					<input class="form-input" v-model="enterpriseForm.branchCode" placeholder="请输入开户行联行号" />
				</view>
				<uv-button type="primary" text="提交企业进件" @click="testEnterpriseRegister" />
			</view>

			<view class="card" v-if="activeMerchantTab == 'personal'">
				<view class="card-title">个人商户进件</view>
				<view class="form-item">
					<text class="form-label">商户名称 *</text>
					<input class="form-input" v-model="personalForm.regName" placeholder="请输入商户名称" />
				</view>
				<view class="form-item">
					<text class="form-label">商户简称</text>
					<input class="form-input" v-model="personalForm.shortName" placeholder="请输入商户简称" />
				</view>
				<view class="form-item">
					<text class="form-label">所属行业MCC</text>
					<input class="form-input" v-model="personalForm.mcc" placeholder="请输入MCC编码" />
				</view>
				<view class="form-item">
					<text class="form-label">场景类型 *</text>
					<input class="form-input" v-model="personalForm.sceneType"
						placeholder="如：test, ONLINE, OFFLINE, ALL" />
				</view>
				<view class="form-item">
					<text class="form-label">经营区 *</text>
					<input class="form-input" v-model="personalForm.districtId" placeholder="请输入经营区编码，如310105" />
				</view>
				<view class="form-item">
					<text class="form-label">经营详细地址 *</text>
					<input class="form-input" v-model="personalForm.detailAddr" placeholder="请输入经营详细地址" />
				</view>
				<view class="form-item">
					<text class="form-label">负责人证件号码 *</text>
					<input class="form-input" v-model="personalForm.legalCertNo" placeholder="请输入负责人身份证号码" />
				</view>
				<view class="form-item">
					<text class="form-label">证件有效期开始日期 *</text>
					<input class="form-input" v-model="personalForm.legalCertBeginDate" placeholder="格式：yyyyMMdd" />
				</view>
				<view class="form-item">
					<text class="form-label">证件有效期截止日期 *</text>
					<input class="form-input" v-model="personalForm.legalCertEndDate" placeholder="格式：yyyyMMdd" />
				</view>
				<view class="form-item">
					<text class="form-label">负责人手机号 *</text>
					<input class="form-input" v-model="personalForm.contactMobileNo" placeholder="请输入负责人手机号" />
				</view>
				<view class="form-item">
					<text class="form-label">负责人邮箱 *</text>
					<input class="form-input" v-model="personalForm.contactEmail" placeholder="请输入负责人邮箱" />
				</view>
				<view class="form-item">
					<text class="form-label">结算卡开户名称 *</text>
					<input class="form-input" v-model="personalForm.cardName" placeholder="请输入开户名称" />
				</view>
				<view class="form-item">
					<text class="form-label">结算卡账号 *</text>
					<input class="form-input" v-model="personalForm.cardNo" placeholder="请输入银行卡号" />
				</view>
				<view class="form-item">
					<text class="form-label">开户行编码 *</text>
					<input class="form-input" v-model="personalForm.bankCode" placeholder="请输入开户行编码，如01030000" />
				</view>
				<uv-button type="primary" text="提交个人进件" @click="testPersonalRegister" />
			</view>

			<view class="card" v-if="activeMerchantTab == 'modify'">
				<view class="card-title">商户基本信息修改</view>
				<view class="form-item">
					<text class="form-label">汇付客户ID *</text>
					<input class="form-input" v-model="modifyForm.huifuId" placeholder="请输入汇付客户ID" />
				</view>
				<view class="form-item">
					<text class="form-label">商户简称</text>
					<input class="form-input" v-model="modifyForm.shortName" placeholder="请输入商户简称" />
				</view>
				<view class="form-item">
					<text class="form-label">所属行业</text>
					<input class="form-input" v-model="modifyForm.mcc" placeholder="请输入MCC编码" />
				</view>
				<view class="form-item">
					<text class="form-label">管理员姓名</text>
					<input class="form-input" v-model="modifyForm.contactName" placeholder="请输入管理员姓名" />
				</view>
				<view class="form-item">
					<text class="form-label">管理员邮箱</text>
					<input class="form-input" v-model="modifyForm.contactEmail" placeholder="请输入管理员邮箱" />
				</view>
				<view class="form-item">
					<text class="form-label">客服电话</text>
					<input class="form-input" v-model="modifyForm.servicePhone" placeholder="请输入客服电话" />
				</view>
				<view class="form-item">
					<text class="form-label">小票名称</text>
					<input class="form-input" v-model="modifyForm.receiptName" placeholder="请输入小票名称" />
				</view>
				<uv-button type="primary" text="提交修改" @click="testModifyMerchant" />
			</view>
		</view>

		<!-- 扫码支付测试 -->
		<view v-if="activeTab === 'payment'" class="content">
			<view class="tab-bars">
				<uv-tabs :list="paymentTabs" @click="switchPaymentTab" />
			</view>
			<view class="card" v-if="activePaymentTab == 'wxh5'">
				<view class="card-title">H5微信支付</view>
				<view class="form-item">
					<text class="form-label">商户号</text>
					<input class="form-input" v-model="h5PayForm.huifuId" placeholder="请输入商户号" />
				</view>
				<view class="form-item">
					<text class="form-label">交易金额</text>
					<input class="form-input" type="number" v-model="h5PayForm.amount" placeholder="请输入金额" />
				</view>
				<view class="form-item">
					<text class="form-label">商品描述</text>
					<input class="form-input" v-model="h5PayForm.goodsDesc" placeholder="请输入商品描述" />
				</view>
				<view class="form-item">
					<text class="form-label">异步通知地址</text>
					<input class="form-input" v-model="h5PayForm.notifyUrl" placeholder="请输入通知地址" />
				</view>
				<view class="form-item">
					<text class="form-label">前端回调地址</text>
					<input class="form-input" v-model="h5PayForm.frontUrl" placeholder="请输入回调地址" />
				</view>
				<view class="form-item">
					<text class="form-label">客户端IP</text>
					<input class="form-input" v-model="h5PayForm.clientIp" placeholder="请输入客户端IP" />
				</view>
				<view class="form-item">
					<text class="form-label">是否延迟交易</text>
					<input class="form-input" v-model="h5PayForm.delayAcctFlag" placeholder="N/Y" />
				</view>
				<uv-button type="primary" text="创建H5微信支付" @click="testH5WechatPay" />
			</view>

			<view class="card" v-if="activePaymentTab == 'wxjsapi'">
				<view class="card-title">微信公众号支付</view>
				<view class="form-item">
					<text class="form-label">商户号</text>
					<input class="form-input" v-model="jsapiPayForm.huifuId" placeholder="请输入商户号" />
				</view>
				<view class="form-item">
					<text class="form-label">交易金额</text>
					<input class="form-input" type="number" v-model="jsapiPayForm.amount" placeholder="请输入金额" />
				</view>
				<view class="form-item">
					<text class="form-label">商品描述</text>
					<input class="form-input" v-model="jsapiPayForm.goodsDesc" placeholder="请输入商品描述" />
				</view>
				<view class="form-item">
					<text class="form-label">交易类型</text>
					<input class="form-input" v-model="jsapiPayForm.tradeType" placeholder="T_JSAPI" />
				</view>
				<view class="form-item">
					<text class="form-label">OpenId</text>
					<input class="form-input" v-model="jsapiPayForm.openId" placeholder="请输入OpenId" />
				</view>
				<view class="form-item">
					<text class="form-label">公众号AppId</text>
					<input class="form-input" v-model="jsapiPayForm.appId" placeholder="请输入公众号AppId" />
				</view>
				<view class="form-item">
					<text class="form-label">异步通知地址</text>
					<input class="form-input" v-model="jsapiPayForm.notifyUrl" placeholder="请输入通知地址" />
				</view>
				<view class="form-item">
					<text class="form-label">交易有效期</text>
					<input class="form-input" v-model="jsapiPayForm.timeExpire" placeholder="格式：yyyyMMddHHmmss" />
				</view>
				<view class="form-item">
					<text class="form-label">是否延迟交易</text>
					<input class="form-input" v-model="jsapiPayForm.delayAcctFlag" placeholder="N/Y" />
				</view>
				<uv-button type="primary" text="创建公众号支付" @click="testJsapiWechatPay" />
			</view>

			<view class="card" v-if="activePaymentTab == 'zfbh5'">
				<view class="card-title">H5支付宝支付</view>
				<view class="form-item">
					<text class="form-label">商户号</text>
					<input class="form-input" v-model="alipayH5Form.huifuId" placeholder="请输入商户号" />
				</view>
				<view class="form-item">
					<text class="form-label">交易金额</text>
					<input class="form-input" type="number" v-model="alipayH5Form.amount" placeholder="请输入金额" />
				</view>
				<view class="form-item">
					<text class="form-label">商品描述</text>
					<input class="form-input" v-model="alipayH5Form.goodsDesc" placeholder="请输入商品描述" />
				</view>
				<view class="form-item">
					<text class="form-label">异步通知地址</text>
					<input class="form-input" v-model="alipayH5Form.notifyUrl" placeholder="请输入通知地址" />
				</view>
				<view class="form-item">
					<text class="form-label">前端回调地址</text>
					<input class="form-input" v-model="alipayH5Form.frontUrl" placeholder="请输入回调地址" />
				</view>
				<view class="form-item">
					<text class="form-label">客户端IP</text>
					<input class="form-input" v-model="alipayH5Form.clientIp" placeholder="请输入客户端IP" />
				</view>
				<view class="form-item">
					<text class="form-label">是否延迟交易</text>
					<input class="form-input" v-model="alipayH5Form.delayAcctFlag" placeholder="N/Y" />
				</view>
				<uv-button type="primary" text="创建H5支付宝支付" @click="testH5AlipayPay" />
			</view>

			<view class="card" v-if="activePaymentTab == 'wxmini'">
				<view class="card-title">微信小程序支付</view>
				<view class="form-item">
					<text class="form-label">商户号 *</text>
					<input class="form-input" v-model="wxMiniForm.huifuId" placeholder="请输入商户号" />
				</view>
				<view class="form-item">
					<text class="form-label">交易金额 *</text>
					<input class="form-input" type="number" v-model="wxMiniForm.amount" placeholder="请输入金额" />
				</view>
				<view class="form-item">
					<text class="form-label">商品描述 *</text>
					<input class="form-input" v-model="wxMiniForm.goodsDesc" placeholder="请输入商品描述" />
				</view>
				<view class="form-item">
					<text class="form-label">交易类型</text>
					<input class="form-input" v-model="wxMiniForm.tradeType" placeholder="T_MINIAPP" />
				</view>
				<view class="form-item">
					<text class="form-label">OpenId *</text>
					<input class="form-input" v-model="wxMiniForm.openId" placeholder="请输入OpenId" />
				</view>
				<view class="form-item">
					<text class="form-label">小程序AppId *</text>
					<input class="form-input" v-model="wxMiniForm.miniAppId" placeholder="请输入小程序AppId" />
				</view>
				<view class="form-item">
					<text class="form-label">子商户AppId</text>
					<input class="form-input" v-model="wxMiniForm.subAppId" placeholder="请输入子商户AppId" />
				</view>
				<view class="form-item">
					<text class="form-label">子商户OpenId</text>
					<input class="form-input" v-model="wxMiniForm.subOpenId" placeholder="请输入子商户OpenId" />
				</view>
				<view class="form-item">
					<text class="form-label">交易有效期</text>
					<input class="form-input" v-model="wxMiniForm.timeExpire" placeholder="格式：yyyyMMddHHmmss" />
				</view>
				<view class="form-item">
					<text class="form-label">场景类型</text>
					<input class="form-input" v-model="wxMiniForm.payScene" placeholder="如：02" />
				</view>
				<view class="form-item">
					<text class="form-label">异步通知地址 *</text>
					<input class="form-input" v-model="wxMiniForm.notifyUrl" placeholder="请输入通知地址" />
				</view>
				<view class="form-item">
					<text class="form-label">是否延迟交易</text>
					<input class="form-input" v-model="wxMiniForm.delayAcctFlag" placeholder="N/Y" />
				</view>
				<view class="form-item">
					<text class="form-label">商品编号</text>
					<input class="form-input" v-model="wxMiniForm.goodsId" placeholder="请输入商品编号" />
				</view>
				<view class="form-item">
					<text class="form-label">商品详情</text>
					<input class="form-input" v-model="wxMiniForm.goodsDetail" placeholder="请输入商品详情" />
				</view>
				<view class="form-item">
					<text class="form-label">附加信息</text>
					<input class="form-input" v-model="wxMiniForm.attach" placeholder="请输入附加信息" />
				</view>
				<uv-button type="primary" text="微信小程序支付" @click="testMiniWechatPay" />
			</view>

			<view class="card" v-if="activePaymentTab == 'zfbmini'">
				<view class="card-title">支付宝小程序支付</view>
				<view class="form-item">
					<text class="form-label">商户号 *</text>
					<input class="form-input" v-model="zfbMiniForm.huifuId" placeholder="请输入商户号" />
				</view>
				<view class="form-item">
					<text class="form-label">交易金额 *</text>
					<input class="form-input" type="number" v-model="zfbMiniForm.amount" placeholder="请输入金额" />
				</view>
				<view class="form-item">
					<text class="form-label">商品描述 *</text>
					<input class="form-input" v-model="zfbMiniForm.goodsDesc" placeholder="请输入商品描述" />
				</view>
				<view class="form-item">
					<text class="form-label">买家ID</text>
					<input class="form-input" v-model="zfbMiniForm.buyerId" placeholder="请输入买家ID" />
				</view>
				<view class="form-item">
					<text class="form-label">异步通知地址</text>
					<input class="form-input" v-model="zfbMiniForm.notifyUrl" placeholder="请输入通知地址" />
				</view>
				<view class="form-item">
					<text class="form-label">是否延迟交易</text>
					<input class="form-input" v-model="zfbMiniForm.delayAcctFlag" placeholder="N/Y" />
				</view>
				<view class="form-item">
					<text class="form-label">商品编号</text>
					<input class="form-input" v-model="zfbMiniForm.goodsId" placeholder="请输入商品编号" />
				</view>
				<view class="form-item">
					<text class="form-label">附加信息</text>
					<input class="form-input" v-model="zfbMiniForm.attach" placeholder="请输入附加信息" />
				</view>
				<uv-button type="primary" text="创建支付宝小程序支付" @click="testMiniAlipayPay" />
			</view>

			<view class="card" v-if="activePaymentTab == 'query'">
				<view class="card-title">交易查询</view>
				<view class="form-item">
					<text class="form-label">交易ID</text>
					<input class="form-input" v-model="queryForm.transId" placeholder="请输入交易ID" />
				</view>
				<uv-button type="primary" text="查询交易" @click="testQueryPay" />
			</view>

			<view class="card" v-if="activePaymentTab == 'refund'">
				<view class="card-title">交易退款</view>
				<view class="form-item">
					<text class="form-label">原交易ID</text>
					<input class="form-input" v-model="refundForm.orgTransId" placeholder="请输入原交易ID" />
				</view>
				<view class="form-item">
					<text class="form-label">退款金额</text>
					<input class="form-input" type="number" v-model="refundForm.amount" placeholder="请输入退款金额" />
				</view>
				<uv-button type="primary" text="提交退款" @click="testRefund" />
			</view>
		</view>

		<!-- 延时交易测试 -->
		<view v-if="activeTab === 'delayed'" class="content">
			<view class="tab-bars">
				<uv-tabs :list="delayedTabs" @click="switchDelayedTab" />
			</view>
			<view class="card" v-if="activeDelayedTab == 'confirm'">
				<view class="card-title">延时交易确认</view>
				<view class="form-item">
					<text class="form-label">商户号</text>
					<input class="form-input" v-model="delayedForm.huifuId" placeholder="请输入商户号" />
				</view>
				<view class="form-item">
					<text class="form-label">交易类型</text>
					<input class="form-input" v-model="delayedForm.payType" placeholder="ACCT_PAYMENT/QUICK_PAY/REMITTANCE_PAY" />
				</view>
				<view class="form-item">
					<text class="form-label">原交易请求日期</text>
					<input class="form-input" v-model="delayedForm.orgReqDate" placeholder="格式：yyyyMMdd" />
				</view>
				<view class="form-item">
					<text class="form-label">原交易请求流水号</text>
					<input class="form-input" v-model="delayedForm.orgReqSeqId" placeholder="请输入原交易请求流水号" />
				</view>
				<uv-button type="primary" text="延时交易确认" @click="testDelayedConfirm" />
			</view>

			<view class="card" v-if="activeDelayedTab == 'confirmQuery'">
				<view class="card-title">延时交易确认查询</view>
				<view class="form-item">
					<text class="form-label">商户号</text>
					<input class="form-input" v-model="delayedConfirmQueryForm.huifuId" placeholder="请输入商户号" />
				</view>
				<view class="form-item">
					<text class="form-label">原交易请求日期</text>
					<input class="form-input" v-model="delayedConfirmQueryForm.orgReqDate" placeholder="格式：yyyyMMdd" />
				</view>
				<view class="form-item">
					<text class="form-label">原交易请求流水号</text>
					<input class="form-input" v-model="delayedConfirmQueryForm.orgReqSeqId" placeholder="请输入原交易请求流水号" />
				</view>
				<uv-button type="primary" text="查询延时交易确认" @click="testDelayedConfirmQuery" />
			</view>

			<view class="card" v-if="activeDelayedTab == 'refund'">
				<view class="card-title">延时交易确认退款</view>
				<view class="form-item">
					<text class="form-label">商户号</text>
					<input class="form-input" v-model="delayedRefundForm.huifuId" placeholder="请输入商户号" />
				</view>
				<view class="form-item">
					<text class="form-label">原交易请求日期</text>
					<input class="form-input" v-model="delayedRefundForm.orgReqDate" placeholder="格式：yyyyMMdd" />
				</view>
				<view class="form-item">
					<text class="form-label">原交易请求流水号</text>
					<input class="form-input" v-model="delayedRefundForm.orgReqSeqId" placeholder="请输入原交易请求流水号" />
				</view>
				<uv-button type="primary" text="延时交易退款" @click="testDelayedRefund" />
			</view>

			<view class="card" v-if="activeDelayedTab == 'refundQuery'">
				<view class="card-title">延时交易确认退款查询</view>
				<view class="form-item">
					<text class="form-label">商户号</text>
					<input class="form-input" v-model="delayedRefundQueryForm.huifuId" placeholder="请输入商户号" />
				</view>
				<view class="form-item">
					<text class="form-label">原交易请求日期</text>
					<input class="form-input" v-model="delayedRefundQueryForm.orgReqDate" placeholder="格式：yyyyMMdd" />
				</view>
				<view class="form-item">
					<text class="form-label">原交易请求流水号</text>
					<input class="form-input" v-model="delayedRefundQueryForm.orgReqSeqId" placeholder="请输入原交易请求流水号" />
				</view>
				<view class="form-item">
					<text class="form-label">原退款全局流水号</text>
					<input class="form-input" v-model="delayedRefundQueryForm.orgHfSeqId" placeholder="请输入原退款全局流水号" />
				</view>
				<uv-button type="primary" text="查询延时交易退款" @click="testDelayedRefundQuery" />
			</view>
		</view>

		<!-- 分账测试 -->
		<view v-if="activeTab === 'profit'" class="content">
			<view class="card">
				<view class="card-title">提交分账</view>
				<view class="form-item">
					<text class="form-label">原交易ID</text>
					<input class="form-input" v-model="profitForm.orgTransId" placeholder="请输入原交易ID" />
				</view>
				<view class="form-item">
					<text class="form-label">分账账户ID</text>
					<input class="form-input" v-model="profitForm.accountId" placeholder="请输入分账账户ID" />
				</view>
				<view class="form-item">
					<text class="form-label">分账金额</text>
					<input class="form-input" type="number" v-model="profitForm.amount" placeholder="请输入分账金额" />
				</view>
				<uv-button type="primary" text="提交分账" @click="testProfitShare" />
			</view>
		</view>

		<!-- 响应结果展示 -->
		<view v-if="responseData" class="result-card">
			<view class="result-title">响应结果</view>
			<scroll-view class="result-content" scroll-y>
				<text>{{ responseData }}</text>
			</scroll-view>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive } from 'vue';

// 大标签页
const tabs = [
	{ key: 'merchant', name: '商户进件' },
	{ key: 'payment', name: '扫码支付' },
	{ key: 'delayed', name: '延时交易' },
	{ key: 'profit', name: '分账管理' }
];
// 商户标签页
const merchantTabs = [
	{ key: 'enterprise', name: '企业进件' },
	{ key: 'personal', name: '个人进件' },
	{ key: 'modify', name: '基本信息' }
];
// 支付标签页
const paymentTabs = [
	{ key: 'wxh5', name: '微信H5支付' },
	{ key: 'wxjsapi', name: '微信公众号支付' },
	{ key: 'wxmini', name: '微信小程序支付' },
	{ key: 'zfbh5', name: '支付宝H5支付' },
	{ key: 'zfbmini', name: '支付宝小程序支付' },
	{ key: 'query', name: '交易查询' },
	{ key: 'refund', name: '交易退款' },
];
// 延时交易标签页
const delayedTabs = [
	{ key: 'confirm', name: '延时交易确认' },
	{ key: 'confirmQuery', name: '交易确认查询' },
	{ key: 'refund', name: '延时交易退款' },
	{ key: 'refundQuery', name: '交易退款查询' },
];
// 当前大标签
const activeTab = ref('payment');
// 当前商户标签
const activeMerchantTab = ref('enterprise');
// 当前支付标签
const activePaymentTab = ref('wxh5');
// 当前延时交易标签
const activeDelayedTab = ref('confirm');

// 切换大标签
const switchTab = (key) => {
	if (activeTab.value === key) return;
	activeTab.value = key;
	// 切换时滚动到顶部
	uni.pageScrollTo({
		scrollTop: 0,
		duration: 300
	});
	// 清空之前的响应结果
	responseData.value = '';
};
// 切换商户标签
const switchMerchantTab = (e) => {
	if (activeMerchantTab.value === e.key) return;
	activeMerchantTab.value = e.key;
	// 切换时滚动到顶部
	uni.pageScrollTo({
		scrollTop: 0,
		duration: 300
	});
	// 清空之前的响应结果
	responseData.value = '';
};
// 切换支付标签
const switchPaymentTab = (e) => {
	if (activePaymentTab.value === e.key) return;
	activePaymentTab.value = e.key;
	// 切换时滚动到顶部
	uni.pageScrollTo({
		scrollTop: 0,
		duration: 300
	});
	// 清空之前的响应结果
	responseData.value = '';
};
// 切换延时交易标签
const switchDelayedTab = (e) => {
	if (activeDelayedTab.value === e.key) return;
	activeDelayedTab.value = e.key;
	// 切换时滚动到顶部
	uni.pageScrollTo({
		scrollTop: 0,
		duration: 300
	});
	// 清空之前的响应结果
	responseData.value = '';
};

// 响应结果
const responseData = ref('');

// 企业商户进件表单
const enterpriseForm = reactive({
	regName: '',
	shortName: '',
	receiptName: '',
	entType: '1',
	mcc: '5411',
	busiType: '1',
	sceneType: 'test',
	licenseCode: '',
	licenseValidityType: '1',
	licenseBeginDate: '',
	licenseEndDate: '',
	regDistrictId: '350203',
	regDetail: '',
	districtId: '310104',
	detailAddr: '',
	legalName: '',
	legalCertType: '00',
	legalCertNo: '',
	legalCertValidityType: '1',
	legalCertBeginDate: '',
	legalCertEndDate: '',
	contactMobileNo: '',
	contactEmail: '',
	loginName: '',
	cardName: '',
	cardNo: '',
	branchCode: ''
});

// 个人商户进件表单
const personalForm = reactive({
	regName: '',
	shortName: '',
	mcc: '',
	sceneType: 'test',
	districtId: '310105',
	detailAddr: '',
	legalCertNo: '',
	legalCertBeginDate: '',
	legalCertEndDate: '',
	contactMobileNo: '',
	contactEmail: '',
	cardName: '',
	cardNo: '',
	bankCode: '01030000'
});

// 商户修改表单
const modifyForm = reactive({
	huifuId: '',
	shortName: '',
	mcc: '',
	contactName: '',
	contactEmail: '',
	servicePhone: '',
	receiptName: ''
});

// H5微信支付表单
const h5PayForm = reactive({
	huifuId: '6666000153390803',
	amount: '0.01',
	goodsDesc: '测试商品',
	notifyUrl: 'https://example.com/notify',
	frontUrl: 'https://example.com/return',
	clientIp: '192.168.1.100',
	delayAcctFlag: 'N'
});

// 微信公众号支付表单
const jsapiPayForm = reactive({
	huifuId: '6666000153390803',
	amount: '0.01',
	goodsDesc: '测试商品',
	tradeType: 'T_JSAPI',
	openId: 'wxfc42adf0c2f58bb6',
	appId: 'o5X3A6b0UTYT60Y5pZkmqMkOAsNQ',
	notifyUrl: 'https://example.com/notify',
	timeExpire: '',
	delayAcctFlag: 'N'
});

// H5支付宝支付表单
const alipayH5Form = reactive({
	huifuId: '6666000153390803',
	amount: '0.01',
	goodsDesc: '测试商品',
	notifyUrl: 'https://example.com/notify',
	frontUrl: 'https://example.com/return',
	clientIp: '192.168.1.100',
	delayAcctFlag: 'N'
});

// 微信小程序支付表单
const wxMiniForm = reactive({
	huifuId: '6666000153390803',
	amount: '0.01',
	goodsDesc: 'hibs自动化-通用版验证',
	tradeType: 'T_MINIAPP',
	openId: 'o61IK7RGQ5lxj1M0XmKxqFZH8Pho',
	miniAppId: 'wxb3699bf4d56598b7',
	subAppId: '',
	subOpenId: '',
	timeExpire: '',
	payScene: '02',
	notifyUrl: 'https://example.com/notify',
	delayAcctFlag: 'N',
	goodsId: '',
	goodsDetail: '',
	attach: ''
});

// 支付宝小程序支付表单
const zfbMiniForm = reactive({
	huifuId: '6666000153390803',
	amount: '0.01',
	goodsDesc: '测试商品',
	buyerId: '',
	notifyUrl: 'https://example.com/notify',
	delayAcctFlag: 'N',
	goodsId: '',
	attach: ''
});

// 查询表单
const queryForm = reactive({
	transId: ''
});

// 退款表单
const refundForm = reactive({
	orgTransId: '',
	amount: '0.01'
});

// 分账表单
const profitForm = reactive({
	orgTransId: '',
	accountId: '',
	amount: '0.01'
});

// 延时交易表单
const delayedForm = reactive({
	huifuId: '6666000153390803',
	payType: 'ACCT_PAYMENT',
	orgReqDate: '',
	orgReqSeqId: ''
});

// 延时交易确认查询表单
const delayedConfirmQueryForm = reactive({
	huifuId: '6666000153390803',
	orgReqDate: '',
	orgReqSeqId: ''
});

// 延时交易退款表单
const delayedRefundForm = reactive({
	huifuId: '6666000153390803',
	orgReqDate: '',
	orgReqSeqId: ''
});

// 延时交易退款查询表单
const delayedRefundQueryForm = reactive({
	huifuId: '6666000153390803',
	orgReqDate: '',
	orgReqSeqId: '',
	orgHfSeqId: ''
});

// 显示响应结果
const showResponse = (data) => {
	responseData.value = JSON.stringify(data, null, 2);
};

// 显示错误提示
const showError = (message) => {
	uni.showToast({
		title: message,
		icon: 'error',
		duration: 2000
	});
};

// 企业商户进件
const testEnterpriseRegister = async () => {
	if (!enterpriseForm.regName) {
		showError('请填写商户名称');
		return;
	}
	if (!enterpriseForm.licenseCode) {
		showError('请填写证照编号');
		return;
	}
	if (!enterpriseForm.regDetail) {
		showError('请填写注册详细地址');
		return;
	}
	if (!enterpriseForm.detailAddr) {
		showError('请填写经营详细地址');
		return;
	}
	if (!enterpriseForm.legalName) {
		showError('请填写法人姓名');
		return;
	}
	if (!enterpriseForm.legalCertNo) {
		showError('请填写法人证件号码');
		return;
	}
	if (!enterpriseForm.contactMobileNo) {
		showError('请填写管理员手机号');
		return;
	}
	if (!enterpriseForm.contactEmail) {
		showError('请填写管理员邮箱');
		return;
	}
	if (!enterpriseForm.loginName) {
		showError('请填写管理员账号');
		return;
	}
	if (!enterpriseForm.cardName) {
		showError('请填写开户名称');
		return;
	}
	if (!enterpriseForm.cardNo) {
		showError('请填写开户账号');
		return;
	}
	if (!enterpriseForm.branchCode) {
		showError('请填写开户行联行号');
		return;
	}

	try {
		const res = await uni.$uv.http.post('/huifu/merchant/register/enterprise', {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			reg_name: enterpriseForm.regName,
			short_name: enterpriseForm.shortName,
			receipt_name: enterpriseForm.receiptName,
			ent_type: enterpriseForm.entType,
			mcc: enterpriseForm.mcc,
			busi_type: enterpriseForm.busiType,
			scene_type: enterpriseForm.sceneType,
			license_code: enterpriseForm.licenseCode,
			license_validity_type: enterpriseForm.licenseValidityType,
			license_begin_date: enterpriseForm.licenseBeginDate,
			license_end_date: enterpriseForm.licenseEndDate,
			reg_district_id: enterpriseForm.regDistrictId,
			reg_detail: enterpriseForm.regDetail,
			district_id: enterpriseForm.districtId,
			detail_addr: enterpriseForm.detailAddr,
			legal_name: enterpriseForm.legalName,
			legal_cert_type: enterpriseForm.legalCertType,
			legal_cert_no: enterpriseForm.legalCertNo,
			legal_cert_validity_type: enterpriseForm.legalCertValidityType,
			legal_cert_begin_date: enterpriseForm.legalCertBeginDate,
			legal_cert_end_date: enterpriseForm.legalCertEndDate,
			contact_mobile_no: enterpriseForm.contactMobileNo,
			contact_email: enterpriseForm.contactEmail,
			login_name: enterpriseForm.loginName,
			card_info: JSON.stringify({
				card_type: '0',
				card_name: enterpriseForm.cardName,
				card_no: enterpriseForm.cardNo,
				branch_code: enterpriseForm.branchCode
			})
		});
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 个人商户进件
const testPersonalRegister = async () => {
	if (!personalForm.regName) {
		showError('请填写商户名称');
		return;
	}
	if (!personalForm.detailAddr) {
		showError('请填写经营详细地址');
		return;
	}
	if (!personalForm.legalCertNo) {
		showError('请填写负责人证件号码');
		return;
	}
	if (!personalForm.legalCertBeginDate) {
		showError('请填写证件有效期开始日期');
		return;
	}
	if (!personalForm.legalCertEndDate) {
		showError('请填写证件有效期截止日期');
		return;
	}
	if (!personalForm.contactMobileNo) {
		showError('请填写负责人手机号');
		return;
	}
	if (!personalForm.contactEmail) {
		showError('请填写负责人邮箱');
		return;
	}
	if (!personalForm.cardName) {
		showError('请填写结算卡开户名称');
		return;
	}
	if (!personalForm.cardNo) {
		showError('请填写结算卡账号');
		return;
	}
	if (!personalForm.bankCode) {
		showError('请填写开户行编码');
		return;
	}

	try {
		const res = await uni.$uv.http.post('/huifu/merchant/register/personal', {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			reg_name: personalForm.regName,
			short_name: personalForm.shortName,
			mcc: personalForm.mcc,
			scene_type: personalForm.sceneType,
			district_id: personalForm.districtId,
			detail_addr: personalForm.detailAddr,
			legal_cert_no: personalForm.legalCertNo,
			legal_cert_begin_date: personalForm.legalCertBeginDate,
			legal_cert_end_date: personalForm.legalCertEndDate,
			contact_mobile_no: personalForm.contactMobileNo,
			contact_email: personalForm.contactEmail,
			card_info: JSON.stringify({
				card_name: personalForm.cardName,
				card_no: personalForm.cardNo,
				prov_id: '310000',
				area_id: '310100',
				bank_code: personalForm.bankCode,
				cert_type: '00',
				cert_no: personalForm.legalCertNo,
				cert_validity_type: personalForm.legalCertBeginDate && personalForm.legalCertEndDate ? '0' : '1',
				cert_begin_date: personalForm.legalCertBeginDate,
				cert_end_date: personalForm.legalCertEndDate,
				mp: personalForm.contactMobileNo
			})
		});
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 商户基本信息修改
const testModifyMerchant = async () => {
	if (!modifyForm.huifuId) {
		showError('请填写汇付客户ID');
		return;
	}

	try {
		const params = {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			huifu_id: modifyForm.huifuId
		};

		if (modifyForm.shortName) params.short_name = modifyForm.shortName;
		if (modifyForm.mcc) params.mcc = modifyForm.mcc;
		if (modifyForm.contactName) params.contact_name = modifyForm.contactName;
		if (modifyForm.contactEmail) params.contact_email = modifyForm.contactEmail;
		if (modifyForm.servicePhone) params.service_phone = modifyForm.servicePhone;
		if (modifyForm.receiptName) params.receipt_name = modifyForm.receiptName;

		const res = await uni.$uv.http.put('/huifu/merchant/modify', params);
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// H5微信支付
const testH5WechatPay = async () => {
	if (!h5PayForm.huifuId) {
		showError('请填写商户号');
		return;
	}
	if (!h5PayForm.amount) {
		showError('请填写交易金额');
		return;
	}
	try {
		const params = {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			huifu_id: h5PayForm.huifuId,
			product_id: 'KAZX',
			trans_amt: parseFloat(h5PayForm.amount),
			goods_desc: h5PayForm.goodsDesc || '测试商品',
			notify_url: h5PayForm.notifyUrl || 'https://example.com/notify',
			front_url: h5PayForm.frontUrl || 'https://example.com/return',
			client_ip: h5PayForm.clientIp || '192.168.1.100'
		};
		if (h5PayForm.delayAcctFlag && h5PayForm.delayAcctFlag !== '') {
			params.delay_acct_flag = h5PayForm.delayAcctFlag;
		}
		const res = await uni.$uv.http.post('/huifu/qrpay/h5/wechat', params);
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 微信公众号支付
const testJsapiWechatPay = async () => {
	if (!jsapiPayForm.huifuId) {
		showError('请填写商户号');
		return;
	}
	if (!jsapiPayForm.amount) {
		showError('请填写交易金额');
		return;
	}
	if (!jsapiPayForm.openId) {
		showError('请填写OpenId');
		return;
	}
	try {
		const params = {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			huifu_id: jsapiPayForm.huifuId,
			goods_desc: jsapiPayForm.goodsDesc || '测试商品',
			trade_type: jsapiPayForm.tradeType || 'T_JSAPI',
			trans_amt: parseFloat(jsapiPayForm.amount),
			notify_url: jsapiPayForm.notifyUrl || 'https://example.com/notify',
			open_id: jsapiPayForm.openId
		};
		if (jsapiPayForm.appId && jsapiPayForm.appId !== '') {
			params.app_id = jsapiPayForm.appId;
		}
		if (jsapiPayForm.timeExpire && jsapiPayForm.timeExpire !== '') {
			params.time_expire = jsapiPayForm.timeExpire;
		}
		if (jsapiPayForm.delayAcctFlag && jsapiPayForm.delayAcctFlag !== '') {
			params.delay_acct_flag = jsapiPayForm.delayAcctFlag;
		}
		const res = await uni.$uv.http.post('/huifu/qrpay/jsapi/wechat', params);
		showResponse(res);

		if (res.code === 200) {
			console.log("支付信息:", res.data);
			WeixinJSBridge.invoke('getBrandWCPayRequest', {
				"appId": res.data.appId,
				"timeStamp": res.data.timeStamp,
				"nonceStr": res.data.nonceStr,
				"package": res.data.package,
				"signType": res.data.signType,
				"paySign": res.data.paySign
			}, function (res) {
				console.log(res);
			});
		}
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// H5支付宝支付
const testH5AlipayPay = async () => {
	if (!alipayH5Form.huifuId) {
		showError('请填写商户号');
		return;
	}
	if (!alipayH5Form.amount) {
		showError('请填写交易金额');
		return;
	}
	try {
		const params = {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			huifu_id: alipayH5Form.huifuId,
			product_id: 'KAZX',
			trans_amt: parseFloat(alipayH5Form.amount),
			goods_desc: alipayH5Form.goodsDesc || '测试商品',
			notify_url: alipayH5Form.notifyUrl || 'https://example.com/notify',
			front_url: alipayH5Form.frontUrl || 'https://example.com/return',
			client_ip: alipayH5Form.clientIp || '192.168.1.100'
		};
		if (alipayH5Form.delayAcctFlag && alipayH5Form.delayAcctFlag !== '') {
			params.delay_acct_flag = alipayH5Form.delayAcctFlag;
		}
		const res = await uni.$uv.http.post('/huifu/qrpay/h5/alipay', params);
		showResponse(res);
		if (res.code === 0 && res.data && res.data.pay_url) {
			uni.setClipboardData({
				data: res.data.pay_url,
				success: () => {
					uni.showToast({
						title: '支付链接已复制',
						icon: 'success'
					});
				}
			});
		}
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 微信小程序支付
const testMiniWechatPay = async () => {
	if (!wxMiniForm.huifuId) {
		showError('请填写商户号');
		return;
	}
	if (!wxMiniForm.amount) {
		showError('请填写交易金额');
		return;
	}
	if (!wxMiniForm.openId) {
		showError('请填写OpenId');
		return;
	}
	if (!wxMiniForm.miniAppId) {
		showError('请填写小程序AppId');
		return;
	}
	try {
		const params = {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			huifu_id: wxMiniForm.huifuId,
			goods_desc: wxMiniForm.goodsDesc || '测试商品',
			trade_type: wxMiniForm.tradeType || 'T_MINIAPP',
			trans_amt: parseFloat(wxMiniForm.amount),
			notify_url: wxMiniForm.notifyUrl || 'https://example.com/notify',
			open_id: wxMiniForm.openId,
			mini_app_id: wxMiniForm.miniAppId
		};
		if (wxMiniForm.subAppId && wxMiniForm.subAppId !== '') {
			params.sub_app_id = wxMiniForm.subAppId;
		}
		if (wxMiniForm.subOpenId && wxMiniForm.subOpenId !== '') {
			params.sub_open_id = wxMiniForm.subOpenId;
		}
		if (wxMiniForm.timeExpire && wxMiniForm.timeExpire !== '') {
			params.time_expire = wxMiniForm.timeExpire;
		}
		if (wxMiniForm.payScene && wxMiniForm.payScene !== '') {
			params.pay_scene = wxMiniForm.payScene;
		}
		if (wxMiniForm.delayAcctFlag && wxMiniForm.delayAcctFlag !== '') {
			params.delay_acct_flag = wxMiniForm.delayAcctFlag;
		}
		if (wxMiniForm.goodsId && wxMiniForm.goodsId !== '') {
			params.goods_id = wxMiniForm.goodsId;
		}
		if (wxMiniForm.goodsDetail && wxMiniForm.goodsDetail !== '') {
			params.goods_detail = wxMiniForm.goodsDetail;
		}
		if (wxMiniForm.attach && wxMiniForm.attach !== '') {
			params.attach = wxMiniForm.attach;
		}
		const res = await uni.$uv.http.post('/huifu/qrpay/mini/wechat', params);
		showResponse(res);
		if (res.code === 200) {
			uni.requestPayment({
				"timeStamp": res.data.timeStamp,
				"nonceStr": res.data.nonceStr,
				"package": res.data.package,
				"signType": "RSA",
				"paySign": res.data.paySign,
			})
		}
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 支付宝小程序支付
const testMiniAlipayPay = async () => {
	if (!zfbMiniForm.huifuId) {
		showError('请填写商户号');
		return;
	}
	if (!zfbMiniForm.amount) {
		showError('请填写交易金额');
		return;
	}
	try {
		const params = {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			huifu_id: zfbMiniForm.huifuId,
			product_id: 'KAZX',
			trans_amt: parseFloat(zfbMiniForm.amount),
			goods_desc: zfbMiniForm.goodsDesc || '测试商品',
			notify_url: zfbMiniForm.notifyUrl || 'https://example.com/notify'
		};
		if (zfbMiniForm.buyerId && zfbMiniForm.buyerId !== '') {
			params.buyer_id = zfbMiniForm.buyerId;
		}
		if (zfbMiniForm.delayAcctFlag && zfbMiniForm.delayAcctFlag !== '') {
			params.delay_acct_flag = zfbMiniForm.delayAcctFlag;
		}
		if (zfbMiniForm.goodsId && zfbMiniForm.goodsId !== '') {
			params.goods_id = zfbMiniForm.goodsId;
		}
		if (zfbMiniForm.attach && zfbMiniForm.attach !== '') {
			params.attach = zfbMiniForm.attach;
		}
		const res = await uni.$uv.http.post('/huifu/qrpay/mini/alipay', params);
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 查询交易
const testQueryPay = async () => {
	if (!queryForm.transId) {
		showError('请输入交易ID');
		return;
	}
	try {
		const res = await uni.$uv.http.get(`/huifu/qrpay/query/${queryForm.transId}`);
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 交易退款
const testRefund = async () => {
	if (!refundForm.orgTransId) {
		showError('请输入原交易ID');
		return;
	}
	try {
		const res = await uni.$uv.http.post('/huifu/qrpay/refund', {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			huifu_id: '6666000108854952',
			org_trans_id: refundForm.orgTransId,
			refund_amt: parseFloat(refundForm.amount),
			refund_reason: '测试退款'
		});
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 延时交易确认
const testDelayedConfirm = async () => {
	if (!delayedForm.huifuId || !delayedForm.payType) {
		showError('请填写商户号和交易类型');
		return;
	}
	try {
		const params = {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			huifu_id: delayedForm.huifuId,
			pay_type: delayedForm.payType
		};
		if (delayedForm.orgReqDate) params.org_req_date = delayedForm.orgReqDate;
		if (delayedForm.orgReqSeqId) params.org_req_seq_id = delayedForm.orgReqSeqId;
		const res = await uni.$uv.http.post('/huifu/delayed/confirm', params);
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 延时交易确认查询
const testDelayedConfirmQuery = async () => {
	if (!delayedConfirmQueryForm.huifuId || !delayedConfirmQueryForm.orgReqDate || !delayedConfirmQueryForm.orgReqSeqId) {
		showError('请填写完整信息');
		return;
	}
	try {
		const res = await uni.$uv.http.post('/huifu/delayed/confirm/query', {
			huifu_id: delayedConfirmQueryForm.huifuId,
			org_req_date: delayedConfirmQueryForm.orgReqDate,
			org_req_seq_id: delayedConfirmQueryForm.orgReqSeqId
		});
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 延时交易退款
const testDelayedRefund = async () => {
	if (!delayedRefundForm.huifuId || !delayedRefundForm.orgReqDate || !delayedRefundForm.orgReqSeqId) {
		showError('请填写完整信息');
		return;
	}
	try {
		const res = await uni.$uv.http.post('/huifu/delayed/refund', {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			huifu_id: delayedRefundForm.huifuId,
			org_req_date: delayedRefundForm.orgReqDate,
			org_req_seq_id: delayedRefundForm.orgReqSeqId
		});
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 延时交易退款查询
const testDelayedRefundQuery = async () => {
	if (!delayedRefundQueryForm.huifuId || !delayedRefundQueryForm.orgReqDate) {
		showError('请填写商户号和原交易请求日期');
		return;
	}
	if (!delayedRefundQueryForm.orgReqSeqId && !delayedRefundQueryForm.orgHfSeqId) {
		showError('请填写原交易请求流水号或原退款全局流水号');
		return;
	}
	try {
		const params = {
			huifu_id: delayedRefundQueryForm.huifuId,
			org_req_date: delayedRefundQueryForm.orgReqDate
		};
		if (delayedRefundQueryForm.orgReqSeqId) params.org_req_seq_id = delayedRefundQueryForm.orgReqSeqId;
		if (delayedRefundQueryForm.orgHfSeqId) params.org_hf_seq_id = delayedRefundQueryForm.orgHfSeqId;
		const res = await uni.$uv.http.post('/huifu/delayed/refund/query', params);
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 提交分账
const testProfitShare = async () => {
	if (!profitForm.orgTransId || !profitForm.accountId) {
		showError('请填写完整信息');
		return;
	}
	try {
		const res = await uni.$uv.http.post('/huifu/profit/share', {
			req_seq_id: generateSeqId(),
			req_date: formatDate(),
			huifu_id: '6666000153390803',
			org_trans_id: profitForm.orgTransId,
			profit_detail: [{
				account_id: profitForm.accountId,
				profit_amt: parseFloat(profitForm.amount),
				profit_desc: '测试分账',
				profit_type: 'SHARE'
			}]
		});
		showResponse(res);
	} catch (e) {
		showError('请求失败: ' + (e.message || e));
	}
};

// 生成请求流水号
const generateSeqId = () => {
	const now = new Date();
	const timestamp = now.getTime().toString();
	const random = Math.random().toString(36).substr(2, 6).toUpperCase();
	return 'TEST' + timestamp + random;
};

// 格式化日期
const formatDate = () => {
	const now = new Date();
	const year = now.getFullYear();
	const month = String(now.getMonth() + 1).padStart(2, '0');
	const day = String(now.getDate()).padStart(2, '0');
	return `${year}${month}${day}`;
};

// 返回上一页
const goBack = () => {
	uni.navigateBack();
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

.tab-bar {
	display: flex;
	background: #fff;
	padding: 0 24rpx;
	border-bottom: 2rpx solid #f0f0f0;
	position: sticky;
	top: 0;
	z-index: 100;
	box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);

	.tab-item {
		flex: 1;
		text-align: center;
		padding: 32rpx 0;
		font-size: 28rpx;
		color: #999;
		position: relative;
		transition: all 0.3s ease;
		cursor: pointer;
		display: flex;
		flex-direction: column;
		align-items: center;

		.tab-text {
			transition: all 0.3s ease;
		}

		.tab-indicator {
			width: 48rpx;
			height: 6rpx;
			background: linear-gradient(90deg, #3c9cff, #6b8cff);
			border-radius: 3rpx;
			margin-top: 12rpx;
			animation: indicatorSlide 0.3s ease;
		}

		&.active {
			.tab-text {
				color: #3c9cff;
				font-weight: 600;
			}
		}

		&:active {
			opacity: 0.7;
		}
	}
}

.tab-bars {
	background: #e5e5e5;
	padding-bottom: 8px;
    margin-bottom: 10px;
}

@keyframes indicatorSlide {
	from {
		width: 0;
		opacity: 0;
	}

	to {
		width: 48rpx;
		opacity: 1;
	}
}

.content {
	padding: 24rpx;
}

.card {
	background: #fff;
	border-radius: 16rpx;
	padding: 32rpx;
	margin-bottom: 24rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);

	.card-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 32rpx;
		padding-left: 20rpx;
		border-left: 8rpx solid #3c9cff;
	}
}

.form-item {
	margin-bottom: 28rpx;

	.form-label {
		display: block;
		font-size: 28rpx;
		color: #666;
		margin-bottom: 16rpx;
	}

	.form-input {
		width: 100%;
		height: 88rpx;
		background: #f8f9fa;
		border-radius: 12rpx;
		padding: 0 24rpx;
		font-size: 28rpx;
		color: #333;
		box-sizing: border-box;
	}
}

.result-card {
	background: #fff;
	margin: 0 24rpx 24rpx;
	border-radius: 16rpx;
	padding: 24rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);

	.result-title {
		font-size: 28rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 16rpx;
	}

	.result-content {
		max-height: 500rpx;
		background: #f8f9fa;
		border-radius: 8rpx;
		padding: 16rpx;
		font-size: 24rpx;
		color: #666;
		word-break: break-all;
	}
}
</style>