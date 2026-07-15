<script>
import { checkVersionUpdate, showUpdateModal } from '@/utils/utils.js';

export default {
  onLaunch: function () {
	// 获取设备信息
    this.getDeviceType();
    // #ifdef H5
	// 检查版本更新
    this.checkVersion();
    // #endif
  },
  onHide: function () {
    // console.log('App Hide')
  },
  methods: {
	// 获取设备信息
	getDeviceType() {
		let SystemInfo = uni.getSystemInfoSync()
		uni.setStorageSync('SystemInfo', SystemInfo);
		// console.log(SystemInfo)
	},
    // 检查版本更新
    async checkVersion() {
      try {
        const result = await checkVersionUpdate();
        if (result && result.hasUpdate) {
          // 显示更新提示
          showUpdateModal(result);
        }
      } catch (error) {
        console.error('版本检测失败:', error);
      }
    }
  }
}
</script>

<style lang="scss">
	@import '@climblee/uv-ui/index.scss';

	/* 全局样式 - 响应式布局适配 */
	/* #ifdef H5 */
	page {
		min-height: 100vh;
		background-color: #f8f8f8;
		position: relative;
	}

	/* 移动端 - 全宽显示 */
	@media (max-width: 767px) {
		page {
			max-width: 100%;
			margin: 0;
		}
	}

	/* 平板端 - 适中宽度居中 */
	@media (min-width: 768px) and (max-width: 1024px) {
		page {
			max-width: 768px;
			margin: 0 auto;
			box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
		}
	}

	/* PC端 - 移动端宽度居中显示 */
	@media (min-width: 1025px) {
		page {
			max-width: 450px;
			margin: 0 auto;
			box-shadow: 0 0 30px rgba(0, 0, 0, 0.15);
			background-color: #fff;
		}

		/* PC端背景色 */
		body {
			background-color: #f0f0f0;
		}
	}

	/* 对于需要全屏的元素（如弹窗、遮罩等），移除宽度限制 */
	.uni-mask,
	.uni-popup,
	.uni-popup__wrapper,
	.uni-popup__wrapper-box,
	.uni-modal,
	uni-popup,
	uni-popup-wrapper,
	uni-popup-wrapper-box,
	uni-modal,
	.uv-popup,
	.uv-popup__wrapper,
	.uv-popup__wrapper-box {
		max-width: none !important;
	}
	/* #endif */
</style>
