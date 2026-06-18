/**
 * 工具函数库
 * 包含项目中常用的工具方法
 */

/**
 * 异或(XOR)校验和计算
 * 用于485锁控板的异或校验
 * @param {string} hexStr - 十六进制字符串
 * @returns {string} 两位大写十六进制校验和
 */
export const xorChecksum = (hexStr) => {
    // 去除空格和非十六进制字符
    const cleanHex = hexStr.replace(/[^0-9A-Fa-f]/g, '');
    
    if (cleanHex.length % 2 !== 0) {
        throw new Error('十六进制字符串长度必须为偶数');
    }

    let checksum = 0;
    for (let i = 0; i < cleanHex.length; i += 2) {
        const byte = parseInt(cleanHex.substr(i, 2), 16);
        checksum ^= byte;
    }
    
    // 返回两位十六进制字符串，大写
    return checksum.toString(16).toUpperCase().padStart(2, '0');
};

/**
 * 格式化时间戳为日期字符串
 * @param {number} timestamp - 时间戳（毫秒）
 * @param {string} format - 格式化字符串，默认 'YYYY-MM-DD HH:mm:ss'
 * @returns {string} 格式化后的日期字符串
 */
export const formatDate = (timestamp, format = 'YYYY-MM-DD HH:mm:ss') => {
    const date = new Date(timestamp);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    
    return format
        .replace('YYYY', year)
        .replace('MM', month)
        .replace('DD', day)
        .replace('HH', hours)
        .replace('mm', minutes)
        .replace('ss', seconds);
};

/**
 * 计算十六进制数据长度（字节）并加5
 * @param {string} hexStr - 十六进制字符串
 * @returns {string} 长度（字节）+ 5 的十六进制格式
 */
export const getHexLengthPlus5 = (hexStr) => {
    // 去除空格和非十六进制字符
    const cleanHex = hexStr.replace(/[^0-9A-Fa-f]/g, '');
    // 每两个字符代表一个字节
    const byteLength = cleanHex.length / 2;
    // 返回十六进制格式（两位大写）
    return (byteLength + 5).toString(16).toUpperCase().padStart(2, '0');
};

/**
 * 生成完整的485锁控板命令字符串
 * @param {string} hexData - 十六进制数据部分
 * @returns {string} 完整命令字符串（5A5A + 长度 + 00A0 + 数据 + 校验和）
 */
export const generateLockCommand = (hexData) => {
    // 去除空格和非十六进制字符
    const cleanHex = hexData.replace(/[^0-9A-Fa-f]/g, '');
    
    // 计算长度+5
    // 每两个字符代表一个字节
    const byteLength = cleanHex.length / 2;
    // 返回十六进制格式（两位大写）
    const lengthHex = (byteLength + 5).toString(16).toUpperCase().padStart(2, '0');
    // console.log(byteLength,lengthHex);
    
    // 计算异或校验和
    const checksum = xorChecksum(cleanHex);
    
    // 拼接完整命令：5A5A + 长度 + 00A0 + 数据 + 校验和
    return `5A5A${lengthHex}00A0${cleanHex}${checksum}`;
};

/**
 * 获取本地存储的版本号
 * @returns {string} 当前版本号
 */
export const getLocalVersion = () => {
    try {
        return uni.getStorageSync('app_version') || '0.0.0';
    } catch (e) {
        return '0.0.0';
    }
};

/**
 * 保存版本号到本地存储
 * @param {string} version - 版本号
 */
export const setLocalVersion = (version) => {
    try {
        uni.setStorageSync('app_version', version);
    } catch (e) {
        console.error('保存版本号失败', e);
    }
};

/**
 * 比较两个版本号
 * @param {string} v1 - 版本号1
 * @param {string} v2 - 版本号2
 * @returns {number} 0: 相等, 1: v1 > v2, -1: v1 < v2
 */
export const compareVersions = (v1, v2) => {
    const parts1 = v1.split('.').map(Number);
    const parts2 = v2.split('.').map(Number);
    const maxLen = Math.max(parts1.length, parts2.length);
    
    for (let i = 0; i < maxLen; i++) {
        const p1 = parts1[i] || 0;
        const p2 = parts2[i] || 0;
        if (p1 > p2) return 1;
        if (p1 < p2) return -1;
    }
    return 0;
};

/**
 * 检查版本更新
 * @returns {Promise<Object|null>} 更新信息或null
 */
export const checkVersionUpdate = async () => {
    try {
        // 带时间戳防缓存，请求静态资源目录下的版本文件
        const response = await fetch('/static/version.json?t=' + Date.now());
        if (!response.ok) {
            console.error('版本检测失败:', response.status);
            return null;
        }
        
        // 获取响应内容类型
        const contentType = response.headers.get('content-type');
        if (!contentType || !contentType.includes('application/json')) {
            console.error('版本检测失败: 响应不是 JSON 格式');
            return null;
        }
        
        const remoteVersion = await response.json();
        
        // 验证版本数据
        if (!remoteVersion || !remoteVersion.version) {
            console.error('版本检测失败: 版本数据格式不正确');
            return null;
        }
        
        const localVersion = getLocalVersion();
        
        // 比较版本号
        const compareResult = compareVersions(remoteVersion.version, localVersion);
        
        if (compareResult > 0) {
            return {
                hasUpdate: true,
                currentVersion: localVersion,
                latestVersion: remoteVersion.version,
                updateTime: remoteVersion.updateTime,
                description: remoteVersion.description
            };
        }
        
        // 更新本地版本号（保持同步）
        setLocalVersion(remoteVersion.version);
        return {
            hasUpdate: false,
            currentVersion: remoteVersion.version,
            latestVersion: remoteVersion.version
        };
    } catch (error) {
        console.error('版本检测异常:', error);
        // 如果是 JSON 解析错误，可能是文件不存在，不提示用户
        if (error instanceof SyntaxError && error.message.includes('Unexpected token')) {
            console.warn('版本文件未找到或格式错误，跳过版本检测');
        }
        return null;
    }
};

/**
 * 显示更新提示弹窗
 * @param {Object} updateInfo - 更新信息
 */
export const showUpdateModal = (updateInfo) => {
    uni.showModal({
        title: '发现新版本',
        content: `当前版本: ${updateInfo.currentVersion}\n最新版本: ${updateInfo.latestVersion}\n\n更新说明:\n${updateInfo.description}`,
        showCancel: false,
        confirmText: '立即刷新',
        success: () => {
            // 保存新版本号并刷新页面
            setLocalVersion(updateInfo.latestVersion);
            // #ifdef H5
            window.location.reload();
            // #endif
            // #ifndef H5
            uni.reLaunch({ url: '/pages/login/login' });
            // #endif
        }
    });
};

/**
 * 防抖函数
 * 在事件被触发n秒后再执行回调，如果在这n秒内又被触发，则重新计时
 * @param {Function} func - 需要防抖的函数
 * @param {number} delay - 延迟时间（毫秒）
 * @param {boolean} immediate - 是否立即执行
 * @returns {Function} 防抖后的函数
 */
export const debounce = (func, delay = 300, immediate = false) => {
    let timeout = null;
    return function (...args) {
        const context = this;
        if (timeout) clearTimeout(timeout);
        if (immediate) {
            const callNow = !timeout;
            timeout = setTimeout(() => {
                timeout = null;
            }, delay);
            if (callNow) func.apply(context, args);
        } else {
            timeout = setTimeout(() => {
                func.apply(context, args);
            }, delay);
        }
    };
};

/**
 * 节流函数
 * 规定在一个单位时间内，只能触发一次函数，如果这个单位时间内触发多次函数，只有一次生效
 * @param {Function} func - 需要节流的函数
 * @param {number} limit - 时间间隔（毫秒）
 * @returns {Function} 节流后的函数
 */
export const throttle = (func, limit = 300) => {
    let inThrottle = false;
    return function (...args) {
        const context = this;
        if (!inThrottle) {
            func.apply(context, args);
            inThrottle = true;
            setTimeout(() => {
                inThrottle = false;
            }, limit);
        }
    };
};

/**
 * 生成二维码内容
 * @param {string} type - 二维码类型（总码:total, 分组:group, 房间:room）
 * @param {string} groupType - 分组类型（分组:group, 房间:room）
 * @param {number} id - ID
 * @returns {string} 二维码内容URL
 */
export const generateQRCodeContent = (type, groupType, id) => {
    return `https://centraliz.bsldtech.cn/#/pages/user/index/index?type=${type}&groupType=${groupType}&id=${id}`;
};

/**
 * 扫码功能（适配小程序和H5）
 * @param {Object} options - 配置参数
 * @param {string} [options.successText='扫码成功'] - 成功提示文本
 * @param {string} [options.failText='扫码失败'] - 失败提示文本
 * @returns {Promise<string>} 扫码结果（二维码内容）
 */
export const scanQRCode = (options = {}) => {
    const { successText = '扫码成功', failText = '扫码失败' } = options;
    
    return new Promise((resolve, reject) => {
        // #ifdef H5
        // H5端使用摄像头扫码
        uni.scanCode({
            success: (res) => {
                // console.log('H5扫码成功:', res);
                uni.showToast({ title: successText, icon: 'success' });
                resolve(res.result);
            },
            fail: (err) => {
                // console.error('H5扫码失败:', err);
                uni.showToast({ title: failText, icon: 'none' });
                reject(err);
            }
        });
        // #endif
        
        // #ifndef H5
        // 小程序端使用系统扫码
        uni.scanCode({
            onlyFromCamera: true,
            scanType: ['qrCode'],
            success: (res) => {
                // console.log('小程序扫码成功:', res);
                uni.showToast({ title: successText, icon: 'success' });
                resolve(res.result);
            },
            fail: (err) => {
                // console.error('小程序扫码失败:', err);
                // 用户取消扫码不算错误
                if (err.errMsg && err.errMsg.includes('cancel')) {
                    reject(new Error('用户取消扫码'));
                } else {
                    uni.showToast({ title: failText, icon: 'none' });
                    reject(err);
                }
            }
        });
        // #endif
    });
};

// 记录操作日志
export const recordOperationLog = (data = {}) => {
	try {
		const res = {
			merchsId: parseInt(data.merchsId),
			devicesId: parseInt(data.devicesId),
			roomId: parseInt(data.roomId),
			code: data.code,
			deviceName: data.deviceName,
			roomName: data.roomName,
			type: uni.getStorageSync('SystemInfo').model || '手机',
            phone: data.phone || '',
			control: data.control,
			status: data.status,
			occupant: data.occupant,
		}
		uni.$uv.http.post('/device/log', res, {
			custom: { auth: true }
		})
	} catch (e) {
		console.error('记录操作日志失败:', e);
	}
		
}

// 默认导出所有工具函数
export default {
    xorChecksum, // 计算校验和
    formatDate, // 格式化日期
    getHexLengthPlus5, // 获取十六进制长度并添加5
    generateLockCommand, // 生成锁命令
    getLocalVersion, // 获取本地版本号
    setLocalVersion, // 设置本地版本号
    compareVersions, // 对比版本号
    checkVersionUpdate, // 检查版本更新
    showUpdateModal, // 显示更新弹窗
    debounce, // 防抖函数
    throttle, // 节流函数
    generateQRCodeContent, // 生成二维码内容
    scanQRCode, // 扫码功能（适配小程序和H5）
    recordOperationLog, // 记录锁操作日志
};
