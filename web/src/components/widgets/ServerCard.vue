<template>
  <el-card class="server-card" shadow="hover">
    <template #header>
      <div class="server-card__header">
        <h3 class="server-card__title">{{ server.name }}</h3>
        <p v-if="server.description" class="server-card__description">{{ server.description }}</p>
      </div>
    </template>

    <div class="server-card__content">
      <!-- 创建中的服务器显示进度条 -->
      <template v-if="isCreating">
        <div class="server-card__progress">
          <span class="server-card__progress-text">
            {{ isCompleting ? '创建完成' : (server.create_status?.current_task || '准备中...') }}
          </span>
          <el-progress
            :percentage="progressPercentage"
            :status="isCompleting ? 'success' : ''"
            class="server-card__progress-bar"
          />
        </div>
      </template>
      <!-- 正常服务器显示玩家信息和按钮 -->
      <template v-else>
        <div class="server-card__info animate-fade-in">
          <div class="server-card__stats">
            <div class="server-card__stat-item">
              <span>版本：</span>
              <span>{{ server.version || '未知' }}</span>
            </div>
            <div class="server-card__stat-item">
              <span>运行状态：</span>
              <el-tag :type="serverStatusType" size="small">
                {{ serverStatusText }}
              </el-tag>
            </div>
            <div class="server-card__stat-item">
              <span>在线玩家：</span>
              <span>{{ server.online_players || 0 }}/{{ server.max_players || 20 }}</span>
            </div>
          </div>

          <div class="server-card__actions">
            <el-button
              v-if="server.active"
              type="danger"
              size="small"
              @click="handleStopServer"
              class="server-card__action-btn"
            >
              <el-icon><VideoPause /></el-icon>
              <span>停止</span>
            </el-button>
            <el-button
              v-else
              type="success"
              size="small"
              @click="handleStartServer"
              class="server-card__action-btn"
            >
              <el-icon><VideoPlay /></el-icon>
              <span>启动</span>
            </el-button>
            <el-button
              type="primary"
              size="small"
              @click="handleManageServer"
              class="server-card__action-btn"
            >
              <el-icon><Setting /></el-icon>
              <span>管理</span>
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="handleUpdateServer"
              class="server-card__action-btn"
            >
              <el-icon><Edit /></el-icon>
              <span>编辑</span>
            </el-button>
            <el-button
              :disabled="server.active"
              type="danger"
              size="small"
              @click="handleDeleteServer"
              class="server-card__action-btn"
            >
              <el-icon><Delete /></el-icon>
              <span>删除</span>
            </el-button>
          </div>
        </div>
      </template>
    </div>
  </el-card>
</template>

<script>
import { ref, onMounted, onBeforeUnmount, computed, watch } from 'vue'
import { requestServerInfo } from '@/composables/server.js'
import {
  Edit,
  Delete,
  Setting,
  VideoPlay,
  VideoPause
} from '@element-plus/icons-vue'

export default {
  name: 'ServerCard',
  components: {
    VideoPause,
    VideoPlay,
    Setting,
    Edit,
    Delete
  },
  props: {
    server: {
      type: Object,
      required: true,
      validator: (value) => {
        return value && typeof value === 'object' && 'id' in value
      }
    }
  },
  emits: ['start', 'stop', 'manage', 'update', 'delete'],
  setup(props, { emit }) {
    // 状态管理
    const localServer = ref(props.server)
    const statusTimer = ref(null)
    const loadingProgress = ref(0)
    const loadingTimer = ref(null)
    const isCompleting = ref(false)

    // 计算属性
    const isCreating = computed(() => localServer.value.create_status?.is_running || isCompleting.value)
    const serverStatusType = computed(() => localServer.value.active ? 'success' : 'info')
    const serverStatusText = computed(() => localServer.value.active ? '运行中' : '已停止')
    const progressPercentage = computed(() => {
      if (!isCreating.value) return 100
      return Math.floor(Math.min(loadingProgress.value, localServer.value.create_status?.percentage || 0))
    })

    // 监听器
    watch(() => props.server, (newVal) => {
      localServer.value = newVal
    }, { deep: true })

    // 方法
    const startStatusUpdate = () => {
      if (statusTimer.value) return
      statusTimer.value = setInterval(async() => {
        try {
          const response = await requestServerInfo(localServer.value.id)
          if (response.data) {
            localServer.value = response.data
            if (!localServer.value.create_status?.is_running && !isCompleting.value) {
              startCompletingAnimation()
            }
          }
        } catch (error) {
          console.error('获取服务器状态失败:', error)
          stopStatusUpdate()
          stopLoadingAnimation()
        }
      }, 1000)
    }

    const stopStatusUpdate = () => {
      if (statusTimer.value) {
        clearInterval(statusTimer.value)
        statusTimer.value = null
      }
    }

    const startLoadingAnimation = () => {
      if (loadingTimer.value) return
      loadingProgress.value = 0
      loadingTimer.value = setInterval(() => {
        if (loadingProgress.value < 95) {
          const realProgress = localServer.value.create_status?.percentage || 0
          const diff = realProgress - loadingProgress.value
          // 根据与真实进度的差距计算增长速度
          const speed = Math.max(0.1, Math.min(2, diff / 10))
          loadingProgress.value = Math.floor(loadingProgress.value + speed)
        }
      }, 100)
    }

    const startCompletingAnimation = () => {
      isCompleting.value = true
      const targetProgress = localServer.value.create_status?.percentage || 0
      if (targetProgress < 100) {
        loadingTimer.value = setInterval(() => {
          if (loadingProgress.value < 100) {
            const diff = 100 - loadingProgress.value
            const speed = Math.max(0.1, Math.min(1, diff / 20))
            loadingProgress.value = Math.floor(loadingProgress.value + speed)
          } else {
            stopLoadingAnimation()
            setTimeout(() => {
              isCompleting.value = false
              stopStatusUpdate()
            }, 1000)
          }
        }, 100)
      } else {
        stopLoadingAnimation()
        setTimeout(() => {
          isCompleting.value = false
          stopStatusUpdate()
        }, 1000)
      }
    }

    const stopLoadingAnimation = () => {
      if (loadingTimer.value) {
        clearInterval(loadingTimer.value)
        loadingTimer.value = null
      }
      loadingProgress.value = 100
    }

    // 事件处理
    const handleStartServer = () => emit('start', localServer.value)
    const handleStopServer = () => emit('stop', localServer.value)
    const handleManageServer = () => emit('manage', localServer.value)
    const handleUpdateServer = () => emit('update', localServer.value)
    const handleDeleteServer = () => emit('delete', localServer.value)

    // 生命周期钩子
    onMounted(() => {
      if (isCreating.value) {
        startStatusUpdate()
        startLoadingAnimation()
      }
    })

    onBeforeUnmount(() => {
      stopStatusUpdate()
      stopLoadingAnimation()
    })

    return {
      localServer,
      isCreating,
      isCompleting,
      serverStatusType,
      serverStatusText,
      progressPercentage,
      handleStartServer,
      handleStopServer,
      handleManageServer,
      handleUpdateServer,
      handleDeleteServer
    }
  }
}
</script>

<style lang="scss" scoped>
.server-card {
  @apply h-full transition-all duration-300 hover:-translate-y-1 hover:shadow-lg bg-gray-800 border border-white/10 rounded-2xl overflow-hidden flex flex-col;

  &__header {
    @apply flex flex-col h-[60px] justify-center;
  }

  &__title {
    @apply text-lg font-semibold text-white;
  }

  &__description {
    @apply text-sm text-gray-500 mt-1 line-clamp-1;
  }

  &__content {
    @apply flex-1 flex flex-col p-3;
  }

  &__progress {
    @apply flex-1 flex flex-col justify-center;

    &-text {
      @apply block mb-1 text-sm text-gray-400;
    }

    &-bar {
      @apply transition-all duration-300;

      :deep(.el-progress-bar__outer) {
        @apply bg-gray-700/50;
      }

      :deep(.el-progress-bar__inner) {
        @apply transition-all duration-300 ease-out;
        background: linear-gradient(90deg, #409eff 0%, #67c23a 100%);
      }
    }
  }

  &__info {
    @apply flex-1 flex flex-col;
  }

  &__stats {
    @apply space-y-1 text-gray-400 mb-2;
  }

  &__stat-item {
    @apply flex justify-between;
  }

  &__actions {
    @apply flex flex-wrap gap-4 mt-auto justify-between;
  }

  &__action-btn {
    @apply flex-1 min-w-[80px];
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease;
}

@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

:deep(.el-button) {
  margin-left: 0 !important;
}
</style>
