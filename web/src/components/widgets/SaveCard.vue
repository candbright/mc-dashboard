<template>
  <el-card class="h-full transition-all duration-300 hover:-translate-y-1" shadow="hover">
    <template #header>
      <el-row justify="space-between" align="middle">
        <el-col>
          <h3 class="text-lg text-white font-semibold">{{ save.name }}</h3>
        </el-col>
      </el-row>
    </template>

    <el-space direction="vertical" :size="16" style="width: 100%">
      <div class="save-info">
        <div class="info-row">
          <span class="info-label">最后修改：</span>
          <span class="info-value">{{ formatDate(save.lastModified) }}</span>
        </div>
        <div class="info-row">
          <span class="info-label">文件大小：</span>
          <span class="info-value">{{ formatSize(save.size) }}</span>
        </div>

        <div class="button-group">
          <el-button type="primary" size="small" @click="$emit('apply', save)">
            <el-icon><Check /></el-icon>
            <span>应用</span>
          </el-button>
          <el-button type="danger" size="small" @click="$emit('delete', save)">
            <el-icon><Delete /></el-icon>
            <span>删除</span>
          </el-button>
        </div>
      </div>
    </el-space>
  </el-card>
</template>

<script>
import { Check, Delete } from '@element-plus/icons-vue'

export default {
  name: 'SaveCard',
  components: {
    Check,
    Delete
  },
  props: {
    save: {
      type: Object,
      required: true
    }
  },
  emits: ['apply', 'delete'],
  methods: {
    formatDate(date) {
      if (!date) return '未知'
      return new Date(date).toLocaleString()
    },
    formatSize(size) {
      if (!size) return '未知'
      const units = ['B', 'KB', 'MB', 'GB']
      let value = size
      let unitIndex = 0
      while (value >= 1024 && unitIndex < units.length - 1) {
        value /= 1024
        unitIndex++
      }
      return `${value.toFixed(2)} ${units[unitIndex]}`
    }
  }
}
</script>

<style scoped>
.save-info {
  @apply flex flex-col flex-1;
}

.info-row {
  @apply flex items-center mb-3;
}

.info-label {
  @apply text-gray-400 w-[100px] flex-shrink-0;
}

.info-value {
  @apply text-gray-200;
}

.button-group {
  @apply flex gap-2 flex-wrap mt-auto pt-4;
}

.button-group .el-button {
  @apply flex-1 min-w-[80px];
}

.button-group .el-button :deep(.el-icon) {
  @apply mr-1;
}
</style>
