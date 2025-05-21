<template>
  <div class="flex flex-col min-h-[calc(100vh-4rem)] p-4 pb-8">
    <!-- 页面头部 -->
    <div class="mb-4">
      <div class="flex justify-between items-center px-6 py-4">
        <div class="flex flex-col">
          <h2 class="text-2xl font-semibold mb-1">存档列表</h2>
          <p class="text-sm text-gray-400">管理您的游戏存档和备份</p>
        </div>
        <el-space>
          <el-button type="info" @click="refreshData">
            <el-icon><Refresh /></el-icon>
            <span>刷新</span>
          </el-button>
        </el-space>
      </div>
    </div>

    <!-- 存档列表内容 -->
    <div class="flex-1 relative" v-loading="loading">
      <div class="mb-4 px-6">
        <el-input
          v-model="searchQuery"
          placeholder="搜索存档..."
          class="w-[400px]"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>

      <el-row :gutter="20" class="px-6">
        <!-- 存档卡片 -->
        <el-col :span="8" v-for="save in filteredSaves" :key="save.name" class="mb-4">
          <save-card
            :save="save"
            @apply="applySave"
            @delete="deleteSave"
          />
        </el-col>

        <!-- 添加上传存档卡片 -->
        <el-col :span="8" class="mb-4">
          <el-card
            class="h-full cursor-pointer transition-all duration-300 hover:-translate-y-1 flex justify-center items-center"
            shadow="hover"
          >
            <el-upload
              :action="uploadUrl"
              :show-file-list="false"
              :before-upload="beforeUpload"
              :limit="1"
              :http-request="customUpload"
            >
              <div class="flex flex-col justify-center items-center gap-4 p-8">
                <el-icon class="text-4xl text-gray-400"><Upload /></el-icon>
                <p class="text-gray-400">上传新存档</p>
              </div>
            </el-upload>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 服务器选择对话框 -->
    <el-dialog
      v-model="serverDialogVisible"
      title="选择服务器"
      width="600px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="true"
    >
      <div class="dialog-content">
        <div class="dialog-header">
          <p class="dialog-description">请选择要应用存档的服务器</p>
        </div>
        <div class="server-list">
          <el-table
            :data="servers"
            style="width: 100%"
            @row-click="handleServerSelect"
            v-loading="loading"
            :row-class-name="tableRowClassName"
            highlight-current-row
          >
            <el-table-column prop="name" label="服务器名称" min-width="200">
              <template #default="{ row }">
                <div class="server-name">
                  <el-icon><Monitor /></el-icon>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="version" label="版本" width="120" align="center">
              <template #default="{ row }">
                <el-tag size="small" effect="plain">{{ row.version }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="row.active ? 'success' : 'info'" effect="light">
                  {{ row.active ? '运行中' : '已停止' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="serverDialogVisible = false">取消</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { Search, Upload, Monitor, Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { uploadServerFile, requestSaveList, deleteSaveFile, applySave, requestServerInfos } from '@/composables/server'
import SaveCard from '@/components/widgets/SaveCard.vue'

export default {
  name: 'SaveList',
  components: {
    Search,
    Upload,
    Monitor,
    Refresh,
    SaveCard
  },
  data() {
    return {
      searchQuery: '',
      uploadUrl: '/api/server/upload',
      saves: [],
      loading: false,
      serverDialogVisible: false,
      servers: [],
      selectedServer: null,
      currentSave: null
    }
  },
  computed: {
    filteredSaves() {
      return this.saves.filter(save =>
        save.name.toLowerCase().includes(this.searchQuery.toLowerCase())
      )
    }
  },
  created() {
    this.fetchSaves()
  },
  methods: {
    async fetchSaves() {
      try {
        this.loading = true
        const response = await requestSaveList()

        if (response?.data) {
          const { items } = response.data
          this.saves = items || []
        } else {
          this.saves = []
          ElMessage.warning('获取存档列表失败')
        }
      } catch (error) {
        console.error('获取存档列表失败:', error)
        this.saves = []
        ElMessage.error(error.response?.data?.message || '获取存档列表失败')
      } finally {
        this.loading = false
      }
    },
    beforeUpload(file) {
      const isMcworld = file.name.endsWith('.mcworld')
      if (!isMcworld) {
        ElMessage.error('只能上传 MCWORLD 格式的存档文件！')
        return false
      }

      const maxSize = 500 * 1024 * 1024
      if (file.size > maxSize) {
        ElMessage.error('文件大小不能超过 500MB！')
        return false
      }

      return true
    },
    async customUpload({ file }) {
      try {
        this.loading = true
        await uploadServerFile(file)
        ElMessage.success('存档上传成功')
        await this.fetchSaves()
      } catch (error) {
        console.error('上传错误:', error)
        ElMessage.error(error.response?.data?.message || '存档上传失败')
      } finally {
        this.loading = false
      }
    },
    async fetchServers() {
      try {
        this.loading = true
        const response = await requestServerInfos({
          page: 0,
          size: 100,
          order: 'desc',
          order_by: 'created_at'
        })

        if (response?.data) {
          const { items } = response.data
          this.servers = items || []
        } else {
          this.servers = []
          ElMessage.warning('获取服务器列表失败')
        }
      } catch (error) {
        console.error('获取服务器列表失败:', error)
        this.servers = []
        ElMessage.error(error.response?.data?.message || '获取服务器列表失败')
      } finally {
        this.loading = false
      }
    },
    async applySave(save) {
      try {
        await this.fetchServers()

        if (this.servers.length === 0) {
          ElMessage.warning('没有可用的服务器')
          return
        }

        this.currentSave = save
        this.serverDialogVisible = true
      } catch (error) {
        console.error('获取服务器列表失败:', error)
        ElMessage.error(error.response?.data?.message || '获取服务器列表失败')
      }
    },
    async handleServerSelect(server) {
      try {
        this.selectedServer = server
        this.serverDialogVisible = false

        await ElMessageBox.confirm(
          `确定要将存档 "${this.currentSave.name}" 应用到服务器 "${server.name}" 吗？`,
          '应用存档',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        this.loading = true
        await applySave(this.currentSave.id, server.id)
        ElMessage.success('存档应用成功')
        await this.fetchSaves()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('应用存档失败:', error)
          ElMessage.error(error.response?.data?.message || '应用存档失败')
        }
      } finally {
        this.loading = false
        this.selectedServer = null
        this.currentSave = null
      }
    },
    async deleteSave(save) {
      try {
        await ElMessageBox.confirm(
          `确定要删除存档 "${save.name}" 吗？`,
          '删除存档',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        this.loading = true
        await deleteSaveFile(save.id)
        ElMessage.success('存档删除成功')
        await this.fetchSaves()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除存档失败:', error)
          ElMessage.error(error.response?.data?.message || '删除存档失败')
        }
      } finally {
        this.loading = false
      }
    },
    tableRowClassName({ row }) {
      return row.active ? 'active-server-row' : ''
    },
    async refreshData() {
      await this.fetchSaves()
    }
  }
}
</script>

<style scoped>
.server-name {
  @apply flex items-center gap-2;
}

.server-name .el-icon {
  @apply text-gray-400;
}

:deep(.active-server-row) {
  @apply cursor-pointer bg-emerald-500/10;
}

:deep(.active-server-row:hover) {
  @apply bg-emerald-500/20;
}

:deep(.el-table__row) {
  @apply cursor-pointer;
}

:deep(.el-table__row:hover) {
  @apply bg-emerald-500/10;
}
</style>
