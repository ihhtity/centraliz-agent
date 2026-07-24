<template>
	<view class="container">
		<uv-navbar :title="'开发助手'" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack">
			<template #right>
				<uv-icon name="more-circle" size="32" color="#333" @click="openSessionActionSheet()" />
			</template>
		</uv-navbar>
		<!-- 聊天区域 -->
		<scroll-view class="chat-container" scroll-y :scroll-top="scrollTop">
			<view v-if="messages.length === 0" class="empty-chat">
				<view class="empty-icon">
					<uv-icon name="server-man" size="64" color="#3c9cff" />
				</view>
				<text class="empty-title">欢迎使用开发助手</text>
				<text class="empty-desc">我可以帮您编写代码、调试接口、生成图表</text>
			</view>
			<!-- 消息列表 -->
			<view v-for="(msg, index) in messages" :key="msg.id" :id="'msg-' + index" class="message-item"
				:class="msg.role">
				<view class="avatar">
					<uv-icon v-if="msg.role === 'user'" name="account" size="24" color="#fff" />
					<uv-icon v-else name="server-man" size="24" color="#3c9cff" />
				</view>
				<view class="message-content">
					<view v-if="msg.thought" class="thought-bubble">
						<uv-icon name="error-circle" size="16" color="#faad14" />
						<text class="thought-text">{{ msg.thought }}</text>
					</view>
					<view class="message-text">
						<template v-for="(item, itemIndex) in parseMessage(msg.content)" :key="itemIndex">
							<pre v-if="item.type === 'code'" class="code-block">
								<code><text class="code-lang" v-if="item.language">{{ item.language }}</text>{{ item.content }}</code>
							</pre>
							<strong v-else-if="item.type === 'bold'" class="bold-text">{{ item.content }}</strong>
							<view v-else-if="item.type === 'list'" class="list-item">
								<text class="list-bullet">•</text>
								<text>{{ item.content }}</text>
							</view>
							<view v-else-if="item.type === 'orderedList'" class="list-item ordered">
								<text>{{ item.content }}</text>
							</view>
							<view v-else-if="item.type === 'quote'" class="quote-block">
								<text>{{ item.content }}</text>
							</view>
							<view v-else-if="item.type === 'table'" class="table-container">
								<table class="markdown-table">
									<template v-for="(row, rowIdx) in item.content" :key="rowIdx">
										<tr v-if="rowIdx === 0" class="table-header">
											<td v-for="(cell, cellIdx) in row.split('|').filter(c => c.trim())" :key="cellIdx">
												{{ cell.trim() }}
											</td>
										</tr>
										<tr v-else-if="row.includes('---')" class="table-divider"></tr>
										<tr v-else class="table-row">
											<td v-for="(cell, cellIdx) in row.split('|').filter(c => c.trim())" :key="cellIdx">
												{{ cell.trim() }}
											</td>
										</tr>
									</template>
								</table>
							</view>
							<view v-else-if="item.type === 'break'" class="line-break"></view>
							<text v-else class="text">{{ item.content }}</text>
						</template>
					</view>
					<view v-if="msg.tool_calls && msg.tool_calls.length > 0" class="tool-calls">
						<view class="tool-calls-title">
							<uv-icon name="empty-news" size="16" color="#3c9cff" />
							<text>工具调用</text>
						</view>
						<view v-for="(tool, idx) in msg.tool_calls" :key="idx" class="tool-call-item">
							<view class="tool-name">{{ tool.tool_name }}</view>
							<view class="tool-params">{{ JSON.stringify(tool.parameters, null, 2) }}</view>
						</view>
					</view>
					<view v-if="msg.tool_results && msg.tool_results.length > 0" class="tool-results">
						<view class="tool-results-title">
							<uv-icon name="checkmark-circle" size="16" color="#67c23a" />
							<text>执行结果</text>
						</view>
						<view v-for="(result, idx) in msg.tool_results" :key="idx" class="tool-result-item"
							:class="result.success ? 'success' : 'error'">
							<view class="tool-name">{{ result.tool_name }}</view>
							<view class="tool-result-output">{{ result.success ? result.output : result.error }}</view>
						</view>
					</view>
					<text class="time">{{ msg.created_at }}</text>
				</view>
			</view>
			<!-- 加载区域 -->
			<view v-if="loading" class="loading-item">
				<view class="avatar">
					<uv-icon name="server-man" size="24" color="#3c9cff" />
				</view>
				<view class="loading-bubble">
					<view class="custom-loading">
						<view class="loading-dot"></view>
						<view class="loading-dot"></view>
						<view class="loading-dot"></view>
					</view>
					<text class="loading-text">思考中...</text>
				</view>
			</view>
		</scroll-view>
		<!-- 快捷命令区域 -->
		<view class="quick-commands">
			<uv-button v-for="cmd in quickCommands" :key="cmd.label" size="mini" type="primary" plain
				@click="handleQuickCommand(cmd.prompt)">
				{{ cmd.label }}
			</uv-button>
		</view>
		<!-- 输入区域 -->
		<view class="input-area">
			<uv-input v-model="inputValue" placeholder="输入您的问题..." @confirm="handleSend" />
			<uv-button type="primary" size="mini" :loading="loading" :disabled="!inputValue.trim() || loading"
				@click="handleSend">
				发送
			</uv-button>
		</view>
		<!-- 确认执行弹窗 -->
		<uv-modal ref="confirmModalRef" title="确认执行" content="此操作将修改以下文件，请确认是否执行？" @confirm="handleConfirm"
			@cancel="closeConfirmModal">
			<view class="affected-files">
				<view v-for="file in affectedFiles" :key="file" class="file-item">
					<uv-icon name="file-text" size="20" color="#999" />
					<text>{{ file }}</text>
				</view>
			</view>
		</uv-modal>
		<!-- 会话管理弹窗 -->
		<uv-action-sheet ref="sessionActionSheetRef" :actions="sessionActions" title="会话管理"
			@select="handleSessionAction" />
		<!-- 新建会话弹窗 -->
		<uv-modal ref="createSessionModalRef" title="新建会话" :showCancelButton="true" @confirm="handleConfirmCreateSession"
			@cancel="closeCreateSessionModal">
			<view class="modal-content">
				<text class="modal-label">请输入会话标题：</text>
				<uv-input v-model="newSessionTitle" placeholder="例如：订单模块开发" />
			</view>
		</uv-modal>
	</view>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue';

const messages = ref([]); // 当前会话的消息列表
const inputValue = ref(''); // 输入框值
const loading = ref(false); // 加载状态
const confirmModalRef = ref(null); // 确认执行弹窗引用
const affectedFiles = ref([]); // 受影响的文件列表
const currentSessionId = ref(''); // 当前会话ID
const currentUserName = ref(''); // 当前用户名称
const sessionActionSheetRef = ref(null); // 会话管理弹窗引用
const sessions = ref([]); // 所有会话列表
const createSessionModalRef = ref(null); // 新建会话弹窗引用
const newSessionTitle = ref(''); // 新建会话标题
const scrollTop = ref(0);
// 会话管理操作列表
const sessionActions = ref([
	{ name: '新建会话', subname: '开始新的对话', color: '#3c9cff' },
	{ name: '切换会话', subname: '选择已有会话', color: '#666' },
	{ name: '删除会话', subname: '删除当前会话', color: '#f56c6c' },
	{ name: '取消', subname: '', color: '#999' },
]);
// 快捷命令列表
const quickCommands = [
	{ label: '修复 Bug', prompt: '请描述你遇到的 Bug，我来帮你分析和修复。' },
	{ label: '生成接口', prompt: '请告诉我需要创建的实体名称和字段，我来生成完整的 CRUD 接口。' },
	{ label: '优化代码', prompt: '请提供需要优化的代码，我来分析并给出优化建议。' },
	{ label: '添加功能', prompt: '请描述你想要添加的新功能，我来帮你实现。' },
	{ label: '生成图表', prompt: '请提供数据和图表类型，我来生成可视化图表。' },
	{ label: '代码审查', prompt: '请提供需要审查的代码，我来进行代码审查。' },
];
// 滚动到底部
const scrollToBottom = async () => {
	nextTick(() => {
    scrollTop.value = 9999;
  });
};

// 加载会话
onMounted(() => {
	loadSessions();
});
// 加载会话
const loadSessions = async () => {
	try {
		const res = await uni.$uv.http.get('/admin/assistant/sessions', { custom: { auth: true } });
		if (res.code === 200 && res.data?.sessions) {
			sessions.value = res.data.sessions;
			if (sessions.value.length > 0 && !currentSessionId.value) {
				const latestSession = sessions.value[0];
				currentSessionId.value = latestSession.id;
				currentUserName.value = latestSession.user_name || '';
				if (latestSession.messages && Array.isArray(latestSession.messages)) {
					messages.value = latestSession.messages.map(msg => ({
						...msg,
						tool_calls: msg.tool_calls || [],
						tool_results: msg.tool_results || [],
						created_at: formatTime(msg.created_at),
					}));
				} else {
					messages.value = [];
				}
				scrollToBottom();
			}
		}
	} catch (error) {
		console.error('加载会话失败', error);
	}
};
// 创建新会话
const createNewSession = () => {
	newSessionTitle.value = '';
	createSessionModalRef.value.open();
};
// 关闭新建会话弹窗
const closeCreateSessionModal = () => {
	createSessionModalRef.value.close();
};
// 确认创建会话
const handleConfirmCreateSession = async () => {
	if (!newSessionTitle.value.trim()) {
		uni.showToast({ title: '请输入会话标题', icon: 'none' });
		return;
	}

	closeCreateSessionModal();
	try {
		const res = await uni.$uv.http.post('/admin/assistant/sessions', { title: newSessionTitle.value.trim() }, { custom: { auth: true } });
		if (res.code === 200 && res.data?.session) {
			currentSessionId.value = res.data.session.id;
			currentUserName.value = res.data.session.user_name || '';
			messages.value = [];
			await loadSessions();
			uni.showToast({ title: '会话已创建', icon: 'success' });
		}
	} catch (error) {
		console.error('创建会话失败', error);
		uni.showToast({ title: '创建失败', icon: 'error' });
	}
};
// 显示会话选择器
const showSessionSelector = async () => {
	if (sessions.value.length === 0) {
		uni.showToast({ title: '暂无会话', icon: 'none' });
		return;
	}

	const sessionNames = sessions.value.map(s => `${s.title || '未命名'}`);

	uni.showActionSheet({
		itemList: sessionNames,
		success: async (res) => {
			const selectedSession = sessions.value[res.tapIndex];
			currentSessionId.value = selectedSession.id;
			currentUserName.value = selectedSession.user_name || '';
			messages.value = [];
			if (selectedSession.messages && Array.isArray(selectedSession.messages)) {
				messages.value = selectedSession.messages.map(msg => ({
					...msg,
					tool_calls: msg.tool_calls || [],
					tool_results: msg.tool_results || [],
					created_at: formatTime(msg.created_at),
				}));
			}
			
			scrollToBottom();
			uni.showToast({ title: '已切换会话', icon: 'success' });
		},
	});
};
// 删除当前会话
const deleteCurrentSession = async () => {
	if (!currentSessionId.value) {
		uni.showToast({ title: '请先选择会话', icon: 'none' });
		return;
	}

	uni.showModal({
		title: '确认删除',
		content: '确定要删除当前会话吗？删除后无法恢复。',
		success: async (res) => {
			if (res.confirm) {
				try {
					const deleteRes = await uni.$uv.http.delete(`/admin/assistant/sessions/${currentSessionId.value}`, { custom: { auth: true } });
					if (deleteRes.code === 200) {
						currentSessionId.value = '';
						currentUserName.value = '';
						messages.value = [];
						await loadSessions();
						uni.showToast({ title: '删除成功', icon: 'success' });
					}
				} catch (error) {
					console.error('删除会话失败', error);
					uni.showToast({ title: '删除失败', icon: 'error' });
				}
			}
		},
	});
};
// 解析消息内容
const parseMessage = (content) => {
	if (!content || typeof content !== 'string') {
		return [{ type: 'text', content: '' }];
	}

	if (content.length > 50000) {
		return [{ type: 'text', content: content.substring(0, 50000) + '...' }];
	}

	const parts = [];
	const lines = content.split('\n');
	let inCodeBlock = false;
	let codeContent = '';
	let codeLanguage = '';
	let inTable = false;
	let tableRows = [];

	for (let i = 0; i < lines.length; i++) {
		const line = lines[i];

		if (line.startsWith('```')) {
			if (inCodeBlock) {
				parts.push({ type: 'code', content: codeContent, language: codeLanguage });
				inCodeBlock = false;
				codeContent = '';
				codeLanguage = '';
			} else {
				inCodeBlock = true;
				codeContent = '';
				codeLanguage = line.slice(3).trim() || 'text';
			}
			continue;
		}

		if (inCodeBlock) {
			codeContent += (codeContent ? '\n' : '') + line;
			continue;
		}

		if (line.includes('|') && line.match(/^[|].*[|]$/)) {
			if (!inTable) {
				inTable = true;
				tableRows = [];
			}
			tableRows.push(line);
			if (i + 1 < lines.length && !lines[i + 1].includes('|')) {
				inTable = false;
				parts.push({ type: 'table', content: tableRows });
				tableRows = [];
			}
			continue;
		}

		if (line.startsWith('**') && line.endsWith('**')) {
			parts.push({ type: 'bold', content: line.slice(2, -2) });
		} else if (line.match(/^\d+\./)) {
			parts.push({ type: 'orderedList', content: line });
		} else if (line.startsWith('- ') || line.startsWith('* ')) {
			parts.push({ type: 'list', content: line.slice(2) });
		} else if (line.startsWith('> ')) {
			parts.push({ type: 'quote', content: line.slice(2) });
		} else if (line === '') {
			parts.push({ type: 'break', content: '' });
		} else {
			let processedLine = line;
			const boldMatches = line.match(/\*\*(.*?)\*\*/g);
			if (boldMatches && boldMatches.length > 0) {
				let lastIndex = 0;
				const subParts = [];
				boldMatches.forEach(match => {
					const start = line.indexOf(match, lastIndex);
					if (start > lastIndex) {
						subParts.push({ type: 'text', content: line.substring(lastIndex, start) });
					}
					subParts.push({ type: 'bold', content: match.slice(2, -2) });
					lastIndex = start + match.length;
				});
				if (lastIndex < line.length) {
					subParts.push({ type: 'text', content: line.substring(lastIndex) });
				}
				parts.push(...subParts);
			} else {
				parts.push({ type: 'text', content: line + '\n' });
			}
		}
	}

	if (inTable && tableRows.length > 0) {
		parts.push({ type: 'table', content: tableRows });
	}

	return parts;
};
// 处理发送消息
const handleSend = async () => {
	if (!inputValue.value.trim() || loading.value) return;

	const userMsg = {
		id: Date.now().toString(),
		role: 'user',
		content: inputValue.value.trim(),
		created_at: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
	};

	messages.value.push(userMsg);
	inputValue.value = '';
	loading.value = true;
	scrollToBottom();

	try {
		const res = await uni.$uv.http.post('/admin/assistant/chat', {
			session_id: currentSessionId.value || undefined,
			message: userMsg.content,
		}, { custom: { auth: true } });

		if (res.code === 200 && res.data?.session) {
			currentSessionId.value = res.data.session.id;
			currentUserName.value = res.data.session.user_name || '';

			const newMessages = res.data.session.messages;
			if (newMessages && Array.isArray(newMessages) && newMessages.length > 0) {
				const lastMsg = newMessages[newMessages.length - 1];
				if (lastMsg.role === 'assistant') {
					messages.value.push({
						id: lastMsg.id,
						role: lastMsg.role,
						content: lastMsg.content,
						thought: lastMsg.thought,
						tool_calls: lastMsg.tool_calls || [],
						tool_results: lastMsg.tool_results || [],
						created_at: formatTime(lastMsg.created_at),
					});
					scrollToBottom();
				}
			}

			if (res.data.session.requires_confirm) {
				affectedFiles.value = res.data.session.affected_files || [];
				confirmModalRef.value.open();
			}
		}
	} catch (error) {
		console.error('发送失败', error);
		uni.showToast({ title: '发送失败', icon: 'error' });
	} finally {
		loading.value = false;
	}
};
// 处理确认操作
const handleConfirm = async () => {
	confirmModalRef.value.close();
	loading.value = true;

	try {
		const res = await uni.$uv.http.post('/admin/assistant/confirm', {
			session_id: currentSessionId.value,
			confirm: true,
		}, { custom: { auth: true } });

		if (res.code === 200 && res.data?.session) {
			currentSessionId.value = res.data.session.id;

			const newMessages = res.data.session.messages;
			if (newMessages.length > messages.value.length) {
				const assistantMsg = newMessages[newMessages.length - 1];
				messages.value.push({
					id: assistantMsg.id,
					role: assistantMsg.role,
					content: assistantMsg.content,
					thought: assistantMsg.thought,
					tool_calls: assistantMsg.tool_calls || [],
					tool_results: assistantMsg.tool_results || [],
					created_at: formatTime(assistantMsg.created_at),
				});
			}
		}
	} catch (error) {
		console.error('确认失败', error);
		uni.showToast({ title: '操作失败', icon: 'error' });
	} finally {
		loading.value = false;
	}
};
// 处理取消操作
const closeConfirmModal = () => {
	confirmModalRef.value.close();
};
// 打开会话操作弹窗
const openSessionActionSheet = () => {
	sessionActionSheetRef.value.open();
};
// 处理会话操作
const handleSessionAction = async (data) => {
	switch (data.name) {
		case '新建会话':
			createNewSession();
			break;
		case '切换会话':
			await showSessionSelector();
			break;
		case '删除会话':
			await deleteCurrentSession();
			break;
	}
};
// 处理快捷命令
const handleQuickCommand = (prompt) => {
	inputValue.value = prompt;
};
// 格式化时间
const formatTime = (dateStr) => {
	const date = new Date(dateStr);
	return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
};
// 返回上一页
const goBack = () => {
	uni.redirectTo({
		url: '/pages/admin/profile/index'
	});
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	display: flex;
	flex-direction: column;
	background-color: #f5f7fa;
}

.chat-container {
	flex: 1;
	width: 96%;
	margin: 40rpx auto;
}

.empty-chat {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 100rpx 0;

	.empty-icon {
		margin-bottom: 24rpx;
	}

	.empty-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 16rpx;
	}

	.empty-desc {
		font-size: 26rpx;
		color: #999;
	}
}

.message-item {
	display: flex;
	gap: 20rpx;
	margin-bottom: 32rpx;

	&.user {
		flex-direction: row-reverse;

		.avatar {
			background: linear-gradient(135deg, #3c9cff, #2b85e4);
		}

		.message-content {
			background: linear-gradient(135deg, #3c9cff, #2b85e4);
			border-radius: 24rpx 24rpx 0 24rpx;

			.text {
				color: #fff;
			}

			.time {
				color: rgba(255, 255, 255, 0.6);
			}

			.thought-bubble {
				background: rgba(0, 0, 0, 0.1);
			}
		}
	}

	&.assistant {
		.avatar {
			background: #fff;
			border: 2rpx solid #e8e8e8;
		}

		.message-content {
			background: #fff;
			border-radius: 24rpx 24rpx 24rpx 0;
			box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.06);
		}
	}

	.avatar {
		width: 80rpx;
		height: 80rpx;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.message-content {
		max-width: 75%;
		padding: 24rpx;

		.thought-bubble {
			display: flex;
			align-items: flex-start;
			gap: 12rpx;
			padding: 16rpx;
			background: #fffbe6;
			border-radius: 16rpx;
			margin-bottom: 16rpx;

			.thought-text {
				font-size: 24rpx;
				color: #ad8b00;
				line-height: 1.5;
			}
		}

		.message-text {
			line-height: 1.6;
		}

		.text {
			font-size: 28rpx;
			color: #333;
			word-break: break-all;
		}

		.bold-text {
			font-size: 28rpx;
			color: #333;
			font-weight: 600;
		}

		.code-block {
			background: #f6f8fa;
			border-radius: 12rpx;
			padding: 20rpx;
			margin: 16rpx 0;
			overflow-x: auto;

			code {
				font-family: 'Monaco', 'Menlo', monospace;
				font-size: 24rpx;
				color: #333;
				line-height: 1.5;
				word-break: break-all;
			}
		}

		.list-item {
			display: flex;
			align-items: flex-start;
			gap: 12rpx;
			padding: 8rpx 0;

			&.ordered {
				padding-left: 16rpx;
			}

			.list-bullet {
				color: #3c9cff;
				font-weight: bold;
				font-size: 28rpx;
			}
		}

		.table-container {
			overflow-x: auto;
			margin: 20rpx 0;
			border-radius: 12rpx;
			border: 2rpx solid #e8e8e8;

			.markdown-table {
				width: 100%;
				border-collapse: collapse;

				.table-header td {
					background: #fafafa;
					padding: 16rpx 12rpx;
					font-size: 26rpx;
					font-weight: 600;
					color: #333;
					border-bottom: 2rpx solid #e8e8e8;
				}

				.table-divider {
					height: 0;
				}

				.table-row td {
					padding: 16rpx 12rpx;
					font-size: 26rpx;
					color: #666;
					border-bottom: 1rpx solid #e8e8e8;

					&:last-child {
						border-right: none;
					}
				}
			}
		}

		.line-break {
			height: 16rpx;
		}

		.code-lang {
			display: block;
			font-size: 22rpx;
			color: #999;
			margin-bottom: 8rpx;
			font-weight: 600;
		}

		.quote-block {
			border-left: 6rpx solid #3c9cff;
			padding: 16rpx 20rpx;
			margin: 16rpx 0;
			background: #f8fafc;
			border-radius: 0 12rpx 12rpx 0;

			text {
				font-size: 26rpx;
				color: #666;
				font-style: italic;
			}
		}

		.tool-calls {
			margin-top: 20rpx;
			padding-top: 20rpx;
			border-top: 2rpx solid #f0f0f0;

			.tool-calls-title {
				display: flex;
				align-items: center;
				gap: 8rpx;
				font-size: 26rpx;
				font-weight: 600;
				color: #333;
				margin-bottom: 16rpx;
			}

			.tool-call-item {
				background: #f8fafc;
				border-radius: 12rpx;
				padding: 16rpx;
				margin-bottom: 12rpx;

				.tool-name {
					font-size: 26rpx;
					font-weight: 600;
					color: #3c9cff;
					margin-bottom: 8rpx;
				}

				.tool-params {
					font-family: 'Monaco', 'Menlo', monospace;
					font-size: 22rpx;
					color: #666;
					word-break: break-all;
					line-height: 1.5;
				}
			}
		}

		.tool-results {
			margin-top: 20rpx;
			padding-top: 20rpx;
			border-top: 2rpx solid #f0f0f0;

			.tool-results-title {
				display: flex;
				align-items: center;
				gap: 8rpx;
				font-size: 26rpx;
				font-weight: 600;
				color: #333;
				margin-bottom: 16rpx;
			}

			.tool-result-item {
				border-radius: 12rpx;
				padding: 16rpx;
				margin-bottom: 12rpx;

				&.success {
					background: #f6ffed;
					border: 2rpx solid #b7eb8f;
				}

				&.error {
					background: #fff2f0;
					border: 2rpx solid #ffccc7;
				}

				.tool-name {
					font-size: 26rpx;
					font-weight: 600;
					margin-bottom: 8rpx;
				}

				&.success .tool-name {
					color: #67c23a;
				}

				&.error .tool-name {
					color: #f56c6c;
				}

				.tool-result-output {
					font-family: 'Monaco', 'Menlo', monospace;
					font-size: 22rpx;
					color: #333;
					word-break: break-all;
					line-height: 1.5;
				}
			}
		}

		.time {
			display: block;
			font-size: 22rpx;
			color: #999;
			margin-top: 12rpx;
			text-align: right;
		}
	}
}

.loading-item {
	display: flex;
	gap: 20rpx;
	margin-bottom: 32rpx;

	.avatar {
		width: 80rpx;
		height: 80rpx;
		border-radius: 50%;
		background: #fff;
		border: 2rpx solid #e8e8e8;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.loading-bubble {
		background: #fff;
		border-radius: 24rpx 24rpx 24rpx 0;
		padding: 24rpx 32rpx;
		box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.06);
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 12rpx;
	}

	.custom-loading {
		display: flex;
		gap: 8rpx;

		.loading-dot {
			width: 16rpx;
			height: 16rpx;
			border-radius: 50%;
			background: #3c9cff;
			animation: dotPulse 1.4s infinite ease-in-out;

			&:nth-child(1) {
				animation-delay: -0.32s;
			}

			&:nth-child(2) {
				animation-delay: -0.16s;
			}

			&:nth-child(3) {
				animation-delay: 0s;
			}
		}
	}

	@keyframes dotPulse {

		0%,
		80%,
		100% {
			opacity: 0.3;
			transform: scale(0.8);
		}

		40% {
			opacity: 1;
			transform: scale(1);
		}
	}

	.loading-text {
		font-size: 24rpx;
		color: #999;
	}
}

.quick-commands {
	display: flex;
	gap: 16rpx;
	padding: 16rpx 24rpx;
	background: #fff;
	border-top: 1rpx solid #f0f0f0;
	flex-wrap: wrap;
}

.input-area {
	display: flex;
	gap: 16rpx;
	padding: 20rpx 24rpx;
	background: #fff;
	border-top: 1rpx solid #f0f0f0;

	:deep(.uv-input) {
		flex: 1;
	}
}

.affected-files {
	padding: 24rpx 0;

	.file-item {
		display: flex;
		align-items: center;
		gap: 16rpx;
		padding: 16rpx 0;
		font-size: 26rpx;
		color: #666;
		border-bottom: 1rpx solid #f0f0f0;

		&:last-child {
			border-bottom: none;
		}
	}
}
</style>