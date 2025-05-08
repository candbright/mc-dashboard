<template>
  <div class="save-list-container">
    <div class="page-header">
      <h1>存档列表</h1>
      <p class="subtitle">管理你的游戏存档和备份</p>
    </div>

    <div class="save-list-content">
      <div class="save-list-header">
        <el-input
          v-model="searchQuery"
          placeholder="搜索存档..."
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-upload
          class="upload-button"
          :action="uploadUrl"
          :show-file-list="false"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :before-upload="beforeUpload"
          :limit="1"
          :http-request="customUpload"
        >
          <el-button type="primary">
            <el-icon><Upload /></el-icon>
            上传存档
          </el-button>
        </el-upload>
      </div>

      <div class="save-list">
        <el-table
          :data="filteredSaves"
          v-loading="loading"
          border
          style="width: 100%"
          :highlight-current-row="false"
        >
          <el-table-column prop="name" label="存档名称" min-width="200">
            <template #default="{ row }">
              <div>
                <h3>{{ row.name }}</h3>
                <p class="save-meta">
                  <span>最后修改: {{ formatDate(row.lastModified) }}</span>
                  <span>大小: {{ formatSize(row.size) }}</span>
                </p>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" align="center">
            <template #default="{ row }">
              <div class="save-actions">
                <el-button type="primary" size="small" @click="applySave(row)">应用</el-button>
                <el-button type="danger" size="small" @click="deleteSave(row)">删除</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :background="true"
          prev-text="上一页"
          next-text="下一页"
          total-text="总计"
          :page-size-text="'条/页'"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 添加服务器选择对话框 -->
    <el-dialog
      v-model="serverDialogVisible"
      title="选择服务器"
      width="600px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="true"
      class="server-select-dialog"
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
import { Search, Upload, Monitor } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { uploadServerFile, requestSaveList, deleteSaveFile, applySave, requestServerInfos } from '@/api/mc/server'

export default {
  name: 'SaveList',

  components: {
    Search,
    Upload,
    Monitor
  },

  data() {
    return {
      searchQuery: '',
      uploadUrl: '/api/server/upload',
      saves: [],
      total: 0,
      currentPage: 1,
      pageSize: 10,
      loading: false,
      // 添加服务器选择相关数据
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
        const response = await requestSaveList({
          page: this.currentPage,
          size: this.pageSize
        })
        this.saves = response.data.items
        this.total = response.data.total
      } catch (error) {
        console.error('获取存档列表失败:', error)
        ElMessage.error('获取存档列表失败')
      } finally {
        this.loading = false
      }
    },

    handleSizeChange(val) {
      this.pageSize = val
      this.fetchSaves()
    },

    handleCurrentChange(val) {
      this.currentPage = val
      this.fetchSaves()
    },

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
    },

    beforeUpload(file) {
      const isMcworld = file.name.endsWith('.mcworld')
      if (!isMcworld) {
        ElMessage.error('只能上传 MCWORLD 格式的存档文件！')
        return false
      }

      // 检查文件大小（限制为 500MB）
      const maxSize = 500 * 1024 * 1024
      if (file.size > maxSize) {
        ElMessage.error('文件大小不能超过 500MB！')
        return false
      }

      return true
    },

    async customUpload({ file }) {
      try {
        await uploadServerFile(file)
        ElMessage.success('存档上传成功')
        this.fetchSaves() // 刷新存档列表
      } catch (error) {
        console.error('上传错误:', error)
        this.handleUploadError(error)
      }
    },

    handleUploadSuccess() {
      // 由于使用了 customUpload，这个方法不会被调用
      // 保留此方法是为了兼容性
    },

    handleUploadError(error) {
      let message = '存档上传失败'
      if (error.response) {
        message = error.response.data?.message || message
      }
      ElMessage.error(message)
    },

    async fetchServers() {
      try {
        this.loading = true
        const response = await requestServerInfos()
        console.log('获取到的服务器列表:', response)

        // 处理不同可能的数据格式
        let serverList = []
        if (response.data) {
          if (Array.isArray(response.data)) {
            serverList = response.data
          } else if (response.data.items && Array.isArray(response.data.items)) {
            serverList = response.data.items
          } else if (typeof response.data === 'object') {
            // 如果是对象，尝试转换为数组
            serverList = Object.values(response.data)
          }
        }

        // 确保每个服务器对象都有必要的属性
        this.servers = serverList.map(server => ({
          id: server.id || '',
          name: server.name || '未命名服务器',
          version: server.version || '未知版本',
          active: Boolean(server.active)
        }))

        console.log('处理后的服务器列表:', this.servers)
      } catch (error) {
        console.error('获取服务器列表失败:', error)
        ElMessage.error('获取服务器列表失败')
        this.servers = []
      } finally {
        this.loading = false
      }
    },

    async applySave(save) {
      try {
        console.log('开始应用存档:', save)
        // 先获取服务器列表
        await this.fetchServers()
        console.log('服务器列表:', this.servers)

        if (this.servers.length === 0) {
          ElMessage.warning('没有可用的服务器')
          return
        }

        // 保存当前选中的存档
        this.currentSave = save
        // 显示服务器选择对话框
        this.serverDialogVisible = true
      } catch (error) {
        console.error('获取服务器列表失败:', error)
        ElMessage.error('获取服务器列表失败')
      }
    },

    async handleServerSelect(server) {
      try {
        console.log('选择服务器:', server)
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

        console.log('开始应用存档到服务器:', server.id, this.currentSave.name)
        await applySave(server.id, this.currentSave.name)
        ElMessage.success('存档应用成功')
        this.fetchSaves() // 刷新存档列表
      } catch (error) {
        if (error !== 'cancel') {
          console.error('应用存档失败:', error)
          ElMessage.error(error.response?.data?.message || '应用存档失败')
        }
      } finally {
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

        await deleteSaveFile(save.name)
        ElMessage.success('存档删除成功')
        this.fetchSaves() // 刷新存档列表
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除存档失败:', error)
          ElMessage.error(error.response?.data?.message || '删除存档失败')
        }
      }
    },

    tableRowClassName({ row }) {
      return row.active ? 'active-server-row' : ''
    }
  }
}
</script>

<style scoped>
.save-list-container {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 40px;
}

.page-header h1 {
  font-size: 32px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 10px;
}

.subtitle {
  font-size: 18px;
  color: var(--el-text-color-secondary);
}

.save-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  gap: 20px;
}

.search-input {
  width: 400px;
}

.save-list {
  margin-bottom: 20px;
}

.save-meta {
  display: flex;
  gap: 20px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-top: 8px;
}

.save-actions {
  display: flex;
  gap: 10px;
  justify-content: center;
}

@media (max-width: 768px) {
  .save-list-header {
    flex-direction: column;
  }

  .search-input {
    width: 100%;
  }
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.server-select-dialog :deep(.el-dialog__header) {
  margin-right: 0;
  padding: 20px 24px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.server-select-dialog :deep(.el-dialog__title) {
  font-size: 18px;
  font-weight: 600;
}

.dialog-content {
  padding: 20px 0;
}

.dialog-header {
  margin-bottom: 20px;
  padding: 0 24px;
}

.dialog-description {
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin: 0;
}

.server-list {
  max-height: 400px;
  overflow-y: auto;
  padding: 0 24px;
}

.server-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.server-name .el-icon {
  font-size: 16px;
  color: var(--el-text-color-secondary);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  padding: 20px 24px;
  border-top: 1px solid var(--el-border-color-lighter);
}

:deep(.active-server-row) {
  cursor: pointer;
  background-color: var(--el-color-primary-light-9);
}

:deep(.active-server-row:hover) {
  background-color: var(--el-color-primary-light-8);
}

:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table__row:hover) {
  background-color: var(--el-color-primary-light-9);
}
</style>
