<template>
    <el-container class="h-full">
        <!-- 侧边栏 -->
        <el-aside class="w-64 flex-shrink-0 bg-gray-800 p-4 flex flex-col shadow-lg relative">
            <el-menu class="border-none space-y-2 flex-1" :default-active="activeMenu" router :collapse="false" background-color="rgb(31 41 55)" text-color="rgb(226 232 240)" active-text-color="rgb(255 255 255)">
                <el-menu-item v-for="(item, index) in filteredMenuItems" :key="index" :index="item.path" class="h-12 leading-[48px] my-1 rounded-lg transition-all duration-300 hover:bg-white/10">
                    <i :class="[getIconForItem(item), 'mr-2 text-lg w-6 text-center']"></i>
                    <span>{{ item.name }}</span>
                </el-menu-item>
            </el-menu>

            <!-- 系统信息 -->
            <div class="absolute bottom-4 left-4 right-4" v-show="showSystemInfo">
                <div class="bg-gray-700 rounded-lg p-4">
                    <div class="flex items-center justify-between mb-2">
                        <span class="text-sm text-gray-400">系统状态</span>
                        <span class="bg-green-500 text-xs px-2 py-1 rounded-full">正常</span>
                    </div>
                    <div class="space-y-1 text-sm text-gray-400">
                        <div class="flex justify-between">
                            <span>CPU 使用率</span>
                            <span>{{ formatCpuUsage(systemStatus.cpu_usage) }}</span>
                        </div>
                        <div class="flex justify-between">
                            <span>内存使用</span>
                            <span>{{ formatMemory(systemStatus.memory_used) }} / {{ formatMemory(systemStatus.memory_total) }}</span>
                        </div>
                        <div class="flex justify-between">
                            <span>磁盘空间</span>
                            <span>{{ formatDisk(systemStatus.disk_used) }} / {{ formatDisk(systemStatus.disk_total) }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </el-aside>

        <!-- 主内容区 -->
        <el-main class="overflow-y-auto p-0">
            <router-view v-slot="{ Component }">
                <transition name="fade" mode="out-in">
                    <component :is="Component" />
                </transition>
            </router-view>
        </el-main>
    </el-container>
</template>

<script>
import { requestSystemStatus } from '@/composables/system.js'

export default {
  name: 'DashboardPage',

  data() {
    return {
      // 系统状态数据
      systemStatus: {
        cpu_usage: 0,
        memory_total: 0,
        memory_used: 0,
        disk_total: 0,
        disk_used: 0,
        uptime: 0,
        os: '',
        arch: ''
      },
      systemTimer: null,
      showSystemInfo: true,
      loading: false
    }
  },

  computed: {
    filteredMenuItems() {
      return this.getMenuItems().filter(item => item.path)
    },
    activeMenu() {
      return this.$route.path
    }
  },

  created() {
    this.getSystemStatus()
    this.startSystemStatusUpdate()
    this.checkScreenHeight()
    window.addEventListener('resize', this.checkScreenHeight)
  },

  beforeUnmount() {
    this.stopSystemStatusUpdate()
    window.removeEventListener('resize', this.checkScreenHeight)
  },

  methods: {
    // 获取系统状态
    async getSystemStatus() {
      if (this.loading) return

      this.loading = true
      try {
        const response = await requestSystemStatus()
        if (response.data) {
          this.systemStatus = response.data
        } else {
          console.error('获取系统状态失败:', response?.message || '未知错误')
        }
      } catch (error) {
        console.error('获取系统状态失败:', error)
      } finally {
        this.loading = false
      }
    },

    // 格式化内存大小
    formatMemory(bytes) {
      if (!bytes) return '0 GB'
      const gb = bytes / (1024 * 1024 * 1024)
      return gb.toFixed(1) + ' GB'
    },

    // 格式化磁盘大小
    formatDisk(bytes) {
      if (!bytes) return '0 GB'
      const gb = bytes / (1024 * 1024 * 1024)
      return Math.round(gb) + ' GB'
    },

    // 格式化 CPU 使用率
    formatCpuUsage(usage) {
      if (!usage) return '0%'
      return usage.toFixed(1) + '%'
    },

    // 启动系统状态更新
    startSystemStatusUpdate() {
      this.stopSystemStatusUpdate()
      this.getSystemStatus() // 立即获取一次
      this.systemTimer = setInterval(() => {
        this.getSystemStatus()
      }, 10000) // 每10秒更新一次
    },

    // 停止系统状态更新
    stopSystemStatusUpdate() {
      if (this.systemTimer) {
        clearInterval(this.systemTimer)
        this.systemTimer = null
      }
    },

    getMenuItems() {
      return [{
        'name': '服务器列表',
        'path': '/server/list',
        'icon': 'fas fa-server'
      },
      {
        'name': '存档列表',
        'path': '/save/list',
        'icon': 'fas fa-archive'
      }
      ]
    },

    getIconForItem(item) {
      return item.icon || 'fas fa-circle'
    },

    checkScreenHeight() {
      this.showSystemInfo = window.innerHeight >= 700
    }
  }
}
</script>

<style scoped>
/* 过渡动画 */
.fade-enter-from,
.fade-leave-to {
    @apply opacity-0;
}

.fade-enter-active,
.fade-leave-active {
    @apply transition-opacity duration-200;
}

</style>
