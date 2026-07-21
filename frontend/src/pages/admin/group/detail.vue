<!-- 分组详情页面 -->
<template>
	<view class="container">
		<uv-navbar title="分组详情" :placeholder="true" @leftClick="goBack" />
		
		<scroll-view scroll-y class="content">
			<!-- 基本信息区域 -->
			<view class="section">
				<view class="section-title">
					<text class="title-text">基本信息</text>
				</view>
				
				<view class="form-item">
					<text class="label">分组名字:</text>
					<view class="input-wrap">
						<uv-input 
							v-model="groupDetail.name" 
							class="form-input" 
							placeholder="请输入分组名称"
						/>
					</view>
				</view>
				
				<view class="form-item">
					<text class="label">联系电话:</text>
					<view class="input-wrap">
						<uv-input 
							v-model="groupDetail.phone" 
							class="form-input" 
							type="number"
							placeholder="请输入联系电话"
						/>
					</view>
				</view>
				
				<view class="form-item">
					<text class="label">场地位置:</text>
					<view class="input-wrap">
						<uv-input 
							v-model="groupDetail.location" 
							class="form-input" 
							placeholder="请输入场地位置"
						/>
					</view>
				</view>

				<view class="form-item">
					<text class="label">分组类型:</text>
					<view class="radio-group">
						<view 
							v-for="option in typeOptions" 
							:key="option.value"
							class="radio-option"
							:class="{ active: groupDetail.type === option.value }"
							@click="groupDetail.type = option.value"
						>
							<view class="radio-circle">
								<view v-if="groupDetail.type === option.value" class="radio-dot"></view>
							</view>
							<text class="radio-label">{{ option.value }}</text>
						</view>
					</view>
				</view>
				
				<view class="form-item readonly">
					<text class="label">房间数量:</text>
					<view class="value-wrap">
						<text class="readonly-value">{{ groupDetail.count || 0 }}</text>
					</view>
				</view>
			</view>
			
			<!-- 设置区域 -->
			<view class="section">
				<view class="section-title">
					<text class="title-text">设置</text>
				</view>
				
				<view class="form-item">
					<text class="label" style="width: 140rpx;">绑定号码</text>
					<view class="radio-group-wrap">
						<view class="radio-group">
							<view 
								v-for="option in bindOptions" 
								:key="option.label"
								class="radio-option"
								:class="{ active: groupDetail.bindNumber === option.label }"
								@click="groupDetail.bindNumber = option.label"
							>
								<view class="radio-circle">
									<view v-if="groupDetail.bindNumber === option.label" class="radio-dot"></view>
								</view>
								<text class="radio-label">{{ option.label }}</text>
							</view>
						</view>
						<view class="help-btn" @click="showBindNumberHelp">
							<uv-icon name="question-circle" color="#3c9cff" size="24" />
						</view>
					</view>
				</view>
				
				<view class="form-item">
					<text class="label">消费推送</text>
					<view class="radio-group-wrap">
						<view class="radio-group">
							<view 
								v-for="option in toggleOptions" 
								:key="option.label"
								class="radio-option"
								:class="{ active: groupDetail.consumePush === option.label }"
								@click="groupDetail.consumePush = option.label"
							>
								<view class="radio-circle">
									<view v-if="groupDetail.consumePush === option.label" class="radio-dot"></view>
								</view>
								<text class="radio-label">{{ option.label }}</text>
							</view>
						</view>
						<view class="help-btn" @click="showConsumePushHelp">
							<uv-icon name="question-circle" color="#3c9cff" size="24" />
						</view>
					</view>
				</view>

				<view class="form-item clickable" v-if="groupDetail.type === '存柜'" @click="goToBillingRules">
					<text class="label">计费规则:</text>
					<view class="arrow-wrap">
						<text class="value-text">{{ groupDetail.ruleName || '请选择' }}</text>
						<uv-icon name="arrow-right" color="#ccc" size="20" />
					</view>
				</view>
				
				<view class="form-item clickable" @click="showQRCodeModal">
					<text class="label">生成二维码:</text>
					<view class="arrow-wrap">
						<text class="value-text">生成分组二维码</text>
						<uv-icon name="arrow-right" color="#ccc" size="20" />
					</view>
				</view>
			</view>

			<!-- 设备管理区域 -->
			<view class="section">
				<view class="section-title">
					<text class="title-text">设备管理</text>
					<view class="add-device-btn" @click="showAddDeviceModal">
						<text class="add-btn-text">添加设备</text>
					</view>
				</view>

				<!-- 设备列表 -->
				<view v-if="deviceList.length === 0" class="empty-device">
					<uv-empty mode="list" text="暂无绑定设备" />
				</view>

				<view v-else class="device-list">
					<view v-for="device in deviceList" :key="device.id" class="device-card">
						<view class="device-info">
							<view class="device-row">
								<text class="info-label">名称：</text>
								<text class="info-value">{{ device.name }}</text>
							</view>
							<view class="device-row">
								<text class="info-label">类型：</text>
								<text class="info-value">{{ device.type }}</text>
							</view>
							<view class="device-row">
								<text class="info-label">状态：</text>
								<text :class="['info-value', 'status-text', getDeviceStatusClass(device.status)]">{{ device.status }}</text>
							</view>
							<view v-if="device.type === '集控'">
								<view class="device-row">
									<text class="info-label">板号：</text>
									<text class="info-value">{{ device.boardNo }}</text>
								</view>
								<view class="device-row">
									<text class="info-label">锁定数量：</text>
									<text class="info-value">{{ device.lockCount }}</text>
								</view>
							</view>
							<view class="device-row">
								<text class="info-label">创建时间：</text>
								<text class="info-value">{{ formatDeviceTime(device.createdAt) }}</text>
							</view>
						</view>
						<view class="device-actions">
							<view class="action-btn edit-btn" @click="editDevice(device)">编辑</view>
							<view class="action-btn delete-btn" @click="deleteDevice(device)">删除</view>
						</view>
					</view>
				</view>
			</view>

			<!-- 修改设备弹出层 -->
			<view v-if="showEditDeviceSheet" class="action-sheet-mask" @click="closeEditDeviceModal">
				<view class="add-device-sheet" @click.stop>
					<view class="sheet-title">修改设备</view>
					
					<view class="edit-form">
						<text class="section-label">设备名称 <text class="required">*</text></text>
						<view class="form-item">
							<input
								v-model="editDeviceForm.name"
								class="name-input"
								placeholder="请输入设备名称"
								placeholder-class="input-placeholder"
								maxlength="10"
							/>
						</view>
					</view>

					<view class="edit-form">
						<text class="section-label">设备状态</text>
						<view class="status-options">
							<view 
								v-for="option in statusOptions" 
								:key="option.value"
								class="status-option"
								:class="{ active: editDeviceForm.status === option.value }"
								@click="editDeviceForm.status = option.value"
							>
								{{ option.value }}
							</view>
						</view>
					</view>

					<view class="edit-form">
						<text class="section-label">板号</text>
						<view class="form-item">
							<input
								v-model="editDeviceForm.boardNo"
								class="name-input"
								type="number"
								placeholder="请输入板号（01-99）"
								placeholder-class="input-placeholder"
								@blur="formatBoardNo"
							/>
						</view>
					</view>
					
					<view class="edit-form">
						<text class="section-label">锁定数量</text>
						<view class="form-item">
							<input
								v-model="editDeviceForm.lockCount"
								class="name-input"
								type="number"
								placeholder="请输入锁定数量"
								placeholder-class="input-placeholder"
							/>
						</view>
					</view>
					
					<view class="add-actions">
						<view class="add-btn cancel" @click="closeEditDeviceModal">取消</view>
						<view class="add-btn confirm" @click="updateDevice">确认修改</view>
					</view>
				</view>
			</view>

			<!-- 二维码弹窗 -->
			<view v-if="showQRCodeSheet" class="action-sheet-mask" @click="closeQRCodeModal">
				<view class="qr-code-sheet" @click.stop>
					<view class="sheet-title">分组二维码</view>
					
					<view class="qr-code-content">
						<view class="qr-code-wrapper">
							<uv-qrcode ref="qrcodeRef" :value="qrCodeContent" :size="200" :isQueueLoadImage="true" :h5SaveIsDownload="true" :h5DownloadName="`${groupDetail.name}.png`" />
						</view>
						<view class="qr-code-desc">
							<text class="desc-text">扫码访问分组</text>
							<text class="desc-hint">{{ qrCodeContent }}</text>
						</view>
					</view>
					
					<view class="qr-actions">
						<view class="qr-btn cancel" @click="closeQRCodeModal">关闭</view>
						<view class="qr-btn confirm" @click="saveQRCode">保存图片</view>
						<view class="qr-btn copy" @click="copyQRCode">复制地址</view>
					</view>
				</view>
			</view>

			<!-- 添加设备弹出层 -->
			<view v-if="showAddDeviceSheet" class="action-sheet-mask" @click="closeAddDeviceModal">
				<view class="add-device-sheet" @click.stop>
					<view class="sheet-title">添加设备</view>
					
					<view class="code-section">
						<text class="section-label">设备编码 <text class="required">*</text></text>
						<view class="code-input-wrapper">
							<input
								v-model="addDeviceForm.code"
								class="code-input"
								placeholder="请输入设备编码"
								placeholder-class="input-placeholder"
							/>
							<view class="scan-btn" @click="scanDeviceCode">
								<uv-icon name="scan" size="28" color="#3c9cff" />
							</view>
						</view>
					</view>
					
					<view class="add-actions">
						<view class="add-btn cancel" @click="closeAddDeviceModal">取消</view>
						<view class="add-btn confirm" @click="addDevice">确认添加</view>
					</view>
				</view>
			</view>

			<!-- 计费规则选择器 -->
			<uv-picker 
				ref="picker"
				title="选择计费规则"
				round="10"
				keyName="name"
				:columns="ruleDetail" 
				@confirm="onRuleConfirm"
				@cancel="picker.close()"
			/>
			
			<!-- 绑定号码说明弹框 -->
			<uv-modal 
				ref="bindNumberModal"
				title="绑定号码说明"
				content="关闭：不绑定手机号码，用户可直接使用；手动：用户需手动输入手机号码进行绑定；自动：系统自动获取用户手机号码进行绑定。"
				:showCancelButton="false"
				@confirm="bindNumberModal.close()"
			/>
			
			<!-- 消费推送说明弹框 -->
			<uv-modal 
				ref="consumePushModal"
				title="消费推送说明"
				content="关闭：不推送消费消息通知；开启：用户消费后推送消息通知到用户手机。"
				:showCancelButton="false"
				@confirm="consumePushModal.close()"
			/>
		</scroll-view>
		
		<!-- 底部按钮 -->
		<view class="bottom-bar">
			<uv-button 
				type="primary" 
				shape="circle" 
				@click="submitForm"
			>
				确认提交
			</uv-button>
		</view>
	</view>
</template>

<script setup>
// 导入Vue响应式API和uniapp生命周期
import { ref, toRaw } from 'vue';
import { onLoad } from '@dcloudio/uni-app';
// 导入二维码工具函数
import { generateQRCodeContent, scanQRCode } from '@/utils/utils';

// 页面加载时获取分组ID并请求详情
onLoad((options) => {
	groupDetail.value.id = options.id || '';
	groupDetail.value.rulesId = options.rulesId || '';
	fetchGroupDetail();
});

// 分组详情数据（包含表单数据）
const groupDetail = ref({});
// 规则数据列表
const ruleDetail = ref([]);
// 商家数据
const merch = uni.getStorageSync('merch') || {};
// 选择器显示状态
const picker = ref(null);
// 绑定号码说明弹框
const bindNumberModal = ref(null);
// 消费推送说明弹框
const consumePushModal = ref(null);
// 分组类型选项
const typeOptions = [ { value: '存柜' }, { value: '零售' } ];
// 绑定号码选项
const bindOptions = [ { label: '自动' }, { label: '手动' }, { label: '关闭' } ];
// 消费推送
const toggleOptions = [ { label: '开启' }, { label: '关闭' }, ];
// 设备列表
const deviceList = ref([]);
// 添加设备弹出层状态
const showAddDeviceSheet = ref(false);
// 添加设备表单
const addDeviceForm = ref({ code: '' });
// 二维码弹出层状态
const showQRCodeSheet = ref(false);
// 二维码内容
const qrCodeContent = ref('');
// 二维码组件ref
const qrcodeRef = ref(null);
// 修改设备弹出层状态
const showEditDeviceSheet = ref(false);
// 修改设备表单
const editDeviceForm = ref({
	id: '',
	name: '',
	boardNo: '',
	lockCount: '',
	status: '在线'
});
// 设备状态选项
const statusOptions = [ { value: '在线' }, { value: '离线' }, { value: '维修' } ];

// 获取设备状态样式类
const getDeviceStatusClass = (status) => {
	switch (status) {
		case '在线':
			return 'online';
		case '离线':
			return 'offline';
		case '维修':
			return 'maintenance';
		default:
			return 'offline';
	}
};

// 格式化设备时间
const formatDeviceTime = (time) => {
	if (!time) return '-';
	return time.replace('T', ' ').substring(0, 19);
};

// 获取分组详情接口
const fetchGroupDetail = () => {
	uni.showLoading({ title: '加载中' });
	
	uni.$uv.http.get(`/group/${groupDetail.value.id}?merchsId=${merch.id}&rulesId=${groupDetail.value.rulesId}`, {
		custom: { auth: true }
	}).then((res) => {
		groupDetail.value = res.data.group;
		ruleDetail.value = [toRaw(res.data.rules)] || [];
		deviceList.value = res.data.devices || [];
		
		uni.hideLoading();
	}).catch((err) => {
		uni.hideLoading();
		console.error('获取分组详情失败:', err);
	});
};

// 提交表单
const submitForm = () => {
	// 校验分组名称
	if (!groupDetail.value.name || !groupDetail.value.name.trim()) {
		uni.showToast({ title: '请输入分组名称', icon: 'none', duration: 3000 });
		return;
	}

	// 校验联系电话
	if (!groupDetail.value.phone || !groupDetail.value.phone.trim()) {
		uni.showToast({ title: '请输入联系电话', icon: 'none', duration: 3000 });
		return;
	}

	// 校验手机号格式
	const phoneRegex = /^1[3-9]\d{9}$/;
	if (!phoneRegex.test(groupDetail.value.phone.trim())) {
		uni.showToast({ title: '请输入有效的手机号码', icon: 'none', duration: 3000 });
		return;
	}

	// 校验计费规则
	if (!groupDetail.value.rulesId && groupDetail.value.type === '存柜') {
		uni.showToast({ title: '请选择计费规则', icon: 'none', duration: 3000 });
		return;
	}

	// 校验设备：分组最少绑定一个锁定数量大于0、设备类型为集控的设备
	// const validDevices = deviceList.value.filter(device => {
	// 	return device.type === '集控' && device.lockCount && parseInt(device.lockCount) > 0;
	// });
	// if (validDevices.length === 0) {
	// 	uni.showToast({ title: '请至少绑定一个锁定数量大于0的集控设备', icon: 'none', duration: 3000 });
	// 	return;
	// }

	// 组装提交数据
	const data = {
		rulesId: parseInt(groupDetail.value.rulesId) || null,
		name: groupDetail.value.name.trim(),
		phone: groupDetail.value.phone.trim(),
		type: groupDetail.value.type.trim(),
		location: groupDetail.value.location ? groupDetail.value.location.trim() : null,
		rulename: groupDetail.value.ruleName ? groupDetail.value.ruleName.trim() : null,
		bind_number: groupDetail.value.bindNumber,
		consume_push: groupDetail.value.consumePush,
	};
	// console.log(data);return
	
	uni.showLoading({ title: '保存中' });
	
	// 调用更新接口
	uni.$uv.http.put(`/group/${groupDetail.value.id}`, data, {
		custom: { auth: true }
	}).then((res) => {
		uni.hideLoading();
		if (res.code !== 200) {
			uni.showToast({ title: res.msg || '保存失败', icon: 'none', duration: 3000 });
			return;
		}

		uni.showToast({ title: '保存成功', icon: 'success' });
		setTimeout(() => {
			goBack();
		}, 1000);
	}).catch((err) => {
		uni.hideLoading();
		console.error('保存分组失败:', err);
		uni.showToast({ title: '保存失败', icon: 'none' });
	});
};

// 打开计费规则选择器
const goToBillingRules = () => {
	if (ruleDetail.value.length === 0) {
		uni.showToast({ title: '暂无可用规则', icon: 'none' });
		return;
	}
	picker.value.open();
};

// 显示绑定号码说明弹框
const showBindNumberHelp = () => {
	bindNumberModal.value.open();
};

// 显示消费推送说明弹框
const showConsumePushHelp = () => {
	consumePushModal.value.open();
};

// 选择规则确认
const onRuleConfirm = (e) => {
	const selectedValue = e.value[0];
	groupDetail.value.ruleName = selectedValue.name;
	groupDetail.value.rulesId = selectedValue.id;
	picker.value.close();
};

// 显示二维码弹窗
const showQRCodeModal = () => {
	if (!groupDetail.value.id) {
		uni.showToast({ title: '分组ID为空', icon: 'none' });
		return;
	}
	qrCodeContent.value = generateQRCodeContent('group', groupDetail.value.type, groupDetail.value.id);
	showQRCodeSheet.value = true;
};

// 关闭二维码弹窗
const closeQRCodeModal = () => {
	showQRCodeSheet.value = false;
};

// 保存二维码图片
const saveQRCode = async () => {
	uni.showLoading({ title: '保存中' });
	
	try {
		// #ifdef H5
		// H5端使用uv-qrcode组件的save方法
		if (qrcodeRef.value && qrcodeRef.value.save) {
			await qrcodeRef.value.save({
				content: qrCodeContent.value,
				success: (res) => {
					console.log('二维码保存成功:', res);
					uni.showToast({ title: '保存成功', icon: 'success' });
				},
				fail: (err) => {
					console.error('二维码保存失败:', err);
					uni.showToast({ title: '保存失败', icon: 'none' });
				}
			});

			closeQRCodeModal();
		} else {
			uni.hideLoading();
			uni.showToast({ title: '组件未就绪', icon: 'none' });
		}
		// #endif
		
		// #ifndef H5
		// 小程序端使用uv-qrcode组件的save方法
		if (qrcodeRef.value && qrcodeRef.value.save) {
			await qrcodeRef.value.save({
				content: qrCodeContent.value,
				success: (res) => {
					console.log('二维码保存成功:', res);
					uni.showToast({ title: '保存成功', icon: 'success' });
				},
				fail: (err) => {
					console.error('二维码保存失败:', err);
					uni.showToast({ title: '保存失败', icon: 'none' });
				}
			});

			closeQRCodeModal();
		} else {
			uni.showToast({ title: '组件未就绪', icon: 'none' });
		}
		// #endif
	} catch (err) {
		console.error('保存二维码失败:', err);
		uni.showToast({ title: '保存失败', icon: 'none' });
	}
};

// 复制二维码地址
const copyQRCode = () => {
	uni.setClipboardData({
		data: qrCodeContent.value,
		success: () => {
			uni.showToast({ title: '复制成功', icon: 'success' })
		},
		fail: () => {
			uni.showToast({ title: '复制失败', icon: 'none' })
		}
	})
}

// 显示添加设备弹出层
const showAddDeviceModal = () => {
	addDeviceForm.value = {
		code: ''
	};
	showAddDeviceSheet.value = true;
};

// 关闭添加设备弹出层
const closeAddDeviceModal = () => {
	showAddDeviceSheet.value = false;
};

// 扫描设备编码
const scanDeviceCode = async () => {
	try {
		const result = await scanQRCode({
			successText: '扫码成功',
			failText: '扫码失败'
		});
		
		addDeviceForm.value.code = result;
		
	} catch (err) {
		// 用户取消扫码不提示错误
		if (err.message !== '用户取消扫码') {
			console.error('扫码失败:', err);
		}
	}
};

// 添加设备
const addDevice = () => {
	if (!addDeviceForm.value.code.trim()) {
		uni.showToast({ title: '请输入设备编码', icon: 'none' });
		return;
	}

	const deviceCode = addDeviceForm.value.code.trim();
	
	// 检查设备是否已添加到当前分组
	const exists = deviceList.value.some(device => {
		return device.code && device.code.toLowerCase() === deviceCode.toLowerCase();
	});
	
	if (exists) {
		uni.showToast({ title: '该设备已添加', icon: 'none' });
		return;
	}

	uni.showLoading({ title: '绑定中' });
	
	uni.$uv.http.post('/device/bind-group', {
		code: deviceCode,
		groupId: parseInt(groupDetail.value.id),
		merchsId: merch.id
	}, {
		custom: { auth: true }
	}).then((res) => {
		uni.hideLoading();
		if (res.code !== 200) {
			uni.showToast({ title: res.msg || '绑定失败', icon: 'none' });
			return;
		}
		
		closeAddDeviceModal();
		uni.showToast({ title: '绑定成功', icon: 'success' });
		// 刷新设备列表
		fetchGroupDetail();
	}).catch((err) => {
		uni.hideLoading();
		console.error('绑定设备失败:', err);
		uni.showToast({ title: '绑定失败', icon: 'none' });
	});
};

// 编辑设备（打开弹出层）
const editDevice = (device) => {
	editDeviceForm.value = {
		id: device.id,
		name: device.name || '',
		boardNo: device.boardNo || '',
		lockCount: device.lockCount || '',
		status: device.status || '在线'
	};
	showEditDeviceSheet.value = true;
};

// 格式化板号（补零，确保2位数字）
const formatBoardNo = () => {
	let boardNo = parseInt(editDeviceForm.value.boardNo);
	// 补零
	editDeviceForm.value.boardNo = boardNo.toString().padStart(2, '0');
};

// 关闭修改设备弹出层
const closeEditDeviceModal = () => {
	showEditDeviceSheet.value = false;
};

// 更新设备
const updateDevice = () => {
	if (!editDeviceForm.value.name.trim()) {
		uni.showToast({ title: '请输入设备名称', icon: 'none' });
		return;
	}
	
	if (editDeviceForm.value.name.length > 10) {
		uni.showToast({ title: '设备名称不能超过10个字', icon: 'none' });
		return;
	}
	
	// 校验板号：必须填写且在01-99之间
	if (editDeviceForm.value.boardNo === '') {
		uni.showToast({ title: '请输入板号', icon: 'none' });
		return;
	}
	const boardNo = parseInt(editDeviceForm.value.boardNo);
	if (isNaN(boardNo) || boardNo < 1 || boardNo > 99) {
		uni.showToast({ title: '板号必须在01-99之间', icon: 'none' });
		return;
	}
	
	const lockCount = parseInt(editDeviceForm.value.lockCount);
	if (editDeviceForm.value.lockCount !== '' && (isNaN(lockCount) || lockCount < 0)) {
		uni.showToast({ title: '锁定数量不能小于0', icon: 'none' });
		return;
	}

	uni.showLoading({ title: '保存中' });
	
	uni.$uv.http.put(`/device/${editDeviceForm.value.id}`, {
		name: editDeviceForm.value.name.trim(),
		boardNo: editDeviceForm.value.boardNo,
		lockCount: lockCount >= 0 ? lockCount : null,
		status: editDeviceForm.value.status
	}, {
		custom: { auth: true }
	}).then((res) => {
		uni.hideLoading();
		if (res.code !== 200) {
			uni.showToast({ title: res.msg || '修改失败', icon: 'none' });
			return;
		}
		
		closeEditDeviceModal();
		uni.showToast({ title: '修改成功', icon: 'success' });
		// 刷新设备列表
		fetchGroupDetail();
	}).catch((err) => {
		uni.hideLoading();
		console.error('修改设备失败:', err);
		uni.showToast({ title: '修改失败', icon: 'none' });
	});
};

// 移除设备（解除分组绑定）
const deleteDevice = (device) => {
	uni.showModal({
		title: '确认移除',
		content: `确定将设备"${device.name || device.code}"从该分组移除吗？`,
		confirmColor: '#fa3534',
		success: (res) => {
			if (res.confirm) {
				// 校验设备：分组最少绑定一个锁定数量大于0、设备类型为集控的设备
				const validDevices = deviceList.value.filter(d => {
					return d.type === '集控' && d.lockCount && parseInt(d.lockCount) > 0;
				});
				
				// 检查当前要删除的设备是否是有效设备
				const isCurrentDeviceValid = device.type === '集控' && device.lockCount && parseInt(device.lockCount) > 0;
				
				// 只有当要删除的是唯一的有效设备时才阻止删除
				if (isCurrentDeviceValid && validDevices.length === 1) {
					uni.showToast({ title: '请至少保留一个锁定数量大于0的集控设备', icon: 'none', duration: 3000 });
					return;
				}

				uni.showLoading({ title: '移除中' });
				
				uni.$uv.http.post('/device/unbind-group', {
					deviceId: device.id
				}, {
					custom: { auth: true }
				}).then((res) => {
					uni.hideLoading();
					if (res.code !== 200) {
						uni.showToast({ title: res.msg || '移除失败', icon: 'none' });
						return;
					}
					
					uni.showToast({ title: '移除成功', icon: 'success' });
					// 从列表中移除该设备
					deviceList.value = deviceList.value.filter(item => item.id !== device.id);
				}).catch((err) => {
					uni.hideLoading();
					console.error('移除设备失败:', err);
					uni.showToast({ title: '移除失败', icon: 'none' });
				});
			}
		}
	});
};

// 返回上一页
const goBack = () => {
	uni.redirectTo({
		url: '/pages/admin/group/manage'
	});
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	display: flex;
	flex-direction: column;
}

.content {
	width: 93%;
	margin: 0 auto;
	flex: 1;
	padding: 20rpx;
	padding-bottom: 160rpx;
}

.section {
	background: #fff;
	border-radius: 16rpx;
	padding: 24rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.06);
}

.section-title {
	padding-bottom: 20rpx;
	margin-bottom: 8rpx;
	border-bottom: 1rpx solid #f0f0f0;
	
	.title-text {
		font-size: 30rpx;
		font-weight: 600;
		color: #333;
	}
}

.form-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 20rpx 0;
	border-bottom: 1rpx solid #f5f5f5;
	
	&:last-child {
		border-bottom: none;
	}
	
	&.readonly {
		.input-wrap {
			background: #f8f9fa;
		}
	}
	
	&.clickable {
		&:active {
			background: #f8f9fa;
		}
	}
	
	.label {
		font-size: 28rpx;
		color: #666;
		font-weight: 500;
		width: 200rpx;
		flex-shrink: 0;
	}
}

.input-wrap {
	flex: 1;
	text-align: right;
	
	.form-input {
		font-size: 28rpx;
		text-align: right;
	}
}

.value-wrap {
	flex: 1;
	text-align: right;
}

.readonly-value {
	font-size: 28rpx;
	color: #999;
	background: #f8f9fa;
	padding: 12rpx 20rpx;
	border-radius: 8rpx;
	display: inline-block;
}

.radio-group {
	display: flex;
	gap: 40rpx;
}

.radio-group-wrap {
	display: flex;
	align-items: center;
	gap: 20rpx;
}

.help-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 8rpx;
}

.radio-option {
	display: flex;
	align-items: center;
	gap: 12rpx;
	
	&.active {
		.radio-circle {
			border-color: #3c9cff;
			background: #fff;
			
			.radio-dot {
				background: #3c9cff;
				transform: scale(1);
			}
		}
		
		.radio-label {
			color: #3c9cff;
			font-weight: 500;
		}
	}
}

.radio-circle {
	width: 36rpx;
	height: 36rpx;
	border-radius: 50%;
	border: 3rpx solid #d9d9d9;
	display: flex;
	align-items: center;
	justify-content: center;
	transition: all 0.2s;
}

.radio-dot {
	width: 20rpx;
	height: 20rpx;
	border-radius: 50%;
	background: #d9d9d9;
	transform: scale(0);
	transition: all 0.2s;
}

.radio-label {
	font-size: 28rpx;
	color: #333;
}

.arrow-wrap {
	display: flex;
	align-items: center;
	gap: 12rpx;
}

.value-text {
	font-size: 28rpx;
	color: #999;
}

.bottom-bar {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	padding: 20rpx 30rpx;
	padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
	background: #fff;
	box-shadow: 0 -4rpx 12rpx rgba(0, 0, 0, 0.08);
}

/* 设备管理样式 */
.section-title {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.add-device-btn {
	display: flex;
	align-items: center;
	gap: 8rpx;
	padding: 8rpx 16rpx;
	background: rgba(60, 156, 255, 0.1);
	border-radius: 20rpx;
	
	.add-btn-text {
		font-size: 24rpx;
		color: #3c9cff;
	}
}

.empty-device {
	padding: 40rpx 0;
}

.device-list {
	display: flex;
	flex-direction: column;
	gap: 16rpx;
}

.device-card {
	display: flex;
	align-items: center;
	gap: 16rpx;
	padding: 20rpx;
	background: #f8f9fa;
	border-radius: 12rpx;
}

.device-icon {
	width: 64rpx;
	height: 64rpx;
	border-radius: 12rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
	
	&.device-server {
		background: linear-gradient(135deg, #3c9cff, #2b85e4);
	}
	
	&.device-control {
		background: linear-gradient(135deg, #a855f7, #9333ea);
	}
	
	&.device-face {
		background: linear-gradient(135deg, #10b981, #059669);
	}
	
	&.device-camera {
		background: linear-gradient(135deg, #f59e0b, #d97706);
	}
	
	&.device-default {
		background: linear-gradient(135deg, #6b7280, #4b5563);
	}
}

.device-card {
	display: flex;
	align-items: center;
	background: #fff;
	border-radius: 16rpx;
	padding: 24rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
}

.device-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	gap: 8rpx;
	
	.device-row {
		display: flex;
		align-items: center;
	}
	
	.info-label {
		font-size: 26rpx;
		color: #999;
		flex-shrink: 0;
	}
	
	.info-value {
		font-size: 26rpx;
		color: #333;
		font-weight: 500;
		
		&.status-text {
			&.online {
				color: #19be6b;
			}
			
			&.offline {
				color: #f56c6c;
			}
			
			&.maintenance {
				color: #f59e0b;
			}
		}
	}
}

.device-actions {
	display: flex;
	flex-direction: column;
	gap: 12rpx;
	flex-shrink: 0;
	margin-left: 30rpx;
	
	.action-btn {
		font-size: 26rpx;
		padding: 12rpx 24rpx;
		border-radius: 8rpx;
		text-align: center;
		white-space: nowrap;
		
		&.edit-btn {
			background: #3c9cff;
			color: #fff;
		}
		
		&.delete-btn {
			background: #f56c6c;
			color: #fff;
		}
	}
}

/* 添加设备弹出层样式 */
.action-sheet-mask {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	align-items: flex-end;
	z-index: 1000;
}

.add-device-sheet {
	width: 100%;
	background: #fff;
	border-radius: 24rpx 24rpx 0 0;
	padding-bottom: calc(40rpx + env(safe-area-inset-bottom));
	animation: slideUp 0.3s ease;
	max-height: 80vh;
	overflow-y: auto;
}

@keyframes slideUp {
	from {
		transform: translateY(100%);
	}
	to {
		transform: translateY(0);
	}
}

.sheet-title {
	text-align: center;
	padding: 40rpx;
	font-size: 32rpx;
	font-weight: 600;
	color: #333;
	border-bottom: 1rpx solid #f0f0f0;
}

.type-section {
	padding: 30rpx;

	.section-label {
		font-size: 28rpx;
		color: #666;
		margin-bottom: 20rpx;
		display: block;

		.required {
			color: #fa3534;
			margin-left: 4rpx;
		}
	}
}

.edit-form {
	padding: 0 30rpx;

	.section-label {
		font-size: 28rpx;
		color: #666;
		margin-bottom: 20rpx;
		display: block;

		.required {
			color: #fa3534;
			margin-left: 4rpx;
		}
	}
	
	.form-item {
		margin-bottom: 0;
	}
	
	.name-input {
		width: 100%;
		height: 88rpx;
		padding: 0 24rpx;
		background: #f8f9fa;
		border-radius: 12rpx;
		font-size: 28rpx;
		color: #333;
		box-sizing: border-box;
	}
	
	.input-placeholder {
		color: #999;
	}
}

.status-options {
	display: flex;
	gap: 24rpx;
}

.status-option {
	flex: 1;
	height: 80rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	border: 1rpx solid #e8e8e8;
	border-radius: 12rpx;
	font-size: 28rpx;
	color: #666;
	transition: all 0.3s;

	&.active {
		background: #2979ff;
		border-color: #2979ff;
		color: #fff;
	}
}

.code-section {
	padding: 30rpx;

	.section-label {
		font-size: 28rpx;
		color: #666;
		margin-bottom: 20rpx;
		display: block;

		.required {
			color: #fa3534;
			margin-left: 4rpx;
		}
	}

	.code-input-wrapper {
		display: flex;
		gap: 16rpx;
		align-items: center;
	}

	.code-input {
		flex: 1;
		height: 88rpx;
		padding: 0 24rpx;
		background: #f8f9fa;
		border-radius: 12rpx;
		font-size: 28rpx;
		color: #333;
		box-sizing: border-box;
	}

	.scan-btn {
		width: 88rpx;
		height: 88rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #f8f9fa;
		border-radius: 12rpx;
	}
}

.add-actions {
	display: flex;
	padding: 0 30rpx 30rpx;
	gap: 24rpx;
	
	.add-btn {
		flex: 1;
		height: 88rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 12rpx;
		font-size: 30rpx;
		font-weight: 500;
		
		&.cancel {
			background: #f5f5f5;
			color: #666;
		}
		
		&.confirm {
			background: #3c9cff;
			color: #fff;
		}
	}
}

/* 二维码弹窗样式 */
.qr-code-sheet {
	width: 100%;
	background: #fff;
	border-radius: 32rpx 32rpx 0 0;
	padding-bottom: env(safe-area-inset-bottom);
	max-height: 80vh;
	overflow-y: auto;
}

.qr-code-content {
	padding: 48rpx 32rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 32rpx;
}

.qr-code-wrapper {
	width: 360rpx;
	height: 360rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	background: #ffffff;
	border: 2rpx solid #e8e8e8;
	border-radius: 20rpx;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
	overflow: hidden;
	flex-shrink: 0;
}

.qr-code-wrapper canvas {
	width: 100% !important;
	height: 100% !important;
	display: block;
}

/* 微信小程序特殊处理：禁止二维码区域滚动 */
/* #ifdef MP-WEIXIN */
.qr-code-content {
	padding: 48rpx 32rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 32rpx;
	touch-action: none;
}

.qr-code-wrapper {
	pointer-events: none;
}
/* #endif */

.qr-code-desc {
	text-align: center;
	width: 100%;
	
	.desc-text {
		font-size: 32rpx;
		font-weight: 600;
		color: #333333;
		display: block;
		margin-bottom: 16rpx;
	}
	
	.desc-hint {
		font-size: 24rpx;
		color: #999999;
		display: block;
		word-break: break-all;
		max-width: 100%;
		line-height: 1.5;
		padding: 0 16rpx;
	}
}

.qr-actions {
	display: flex;
	padding: 0 32rpx 32rpx;
	gap: 24rpx;
	border-top: 1rpx solid #f5f5f5;
	padding-top: 28rpx;
	
	.qr-btn {
		flex: 1;
		height: 88rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 44rpx;
		font-size: 30rpx;
		font-weight: 500;
		
		&.cancel {
			background: #f5f7fa;
			color: #666666;
			border: 1rpx solid #e8e8e8;
		}
		
		&.confirm {
			background: linear-gradient(135deg, #3c9cff 0%, #2979ff 100%);
			color: #ffffff;
			box-shadow: 0 4rpx 16rpx rgba(60, 156, 255, 0.3);
		}

		&.copy {
			background: linear-gradient(135deg, #ffc83c 0%, #ffab29 100%);
			color: #ffffff;
			box-shadow: 0 0.125rem 0.5rem rgb(255 60 60 / 30%);
		}
	}
}
</style>