<script>
import { checkVersionUpdate, showUpdateModal } from '@/utils/utils.js';

export default {
  onLaunch: function () {
    // #ifdef H5
    // 只有H5环境才检查版本更新
    this.checkVersion();
    // #endif
  },
  onHide: function () {
    // console.log('App Hide')
  },
  methods: {
    /**
     * 检查版本更新
     */
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
</style>
