<template>
  <div class="flex flex-col min-h-[calc(100vh-4rem)] p-4 pb-8">
    <!-- 页面头部 -->
    <div class="mb-4">
      <div class="flex justify-between items-center px-6 py-4">
        <div class="flex flex-col">
          <h2 class="text-2xl font-semibold mb-1">服务器列表</h2>
          <p class="text-sm text-gray-400">管理您的 Minecraft 服务器</p>
        </div>
        <el-space>
          <el-button type="info" @click="refreshData">
            <el-icon><Refresh /></el-icon>
            <span>刷新</span>
          </el-button>
          <el-button type="success" @click="handleShowCreateDialog">
            <el-icon><Plus /></el-icon>
            <span>新建服务器</span>
          </el-button>
        </el-space>
      </div>
    </div>

    <!-- 服务器列表内容 -->
    <div class="flex-1 relative" v-loading="loading">
      <el-row :gutter="20">
        <!-- 服务器卡片 -->
        <el-col :span="8" v-for="server in tableData" :key="server.id" class="mb-4">
          <server-card
            :server="server"
            @start="handleStartServer"
            @stop="handleStopServer"
            @manage="handleManageServer"
            @update="handleShowUpdateDialog"
            @delete="handleDeleteServer"
          />
        </el-col>

        <!-- 添加服务器卡片 -->
        <el-col :span="8" class="mb-4">
          <el-card
            class="h-full cursor-pointer transition-all duration-300 hover:-translate-y-1 flex justify-center items-center"
            shadow="hover"
            @click="handleShowCreateDialog"
          >
            <div class="flex flex-col justify-center items-center gap-4 p-8">
              <el-icon class="text-4xl text-gray-400"><Plus /></el-icon>
              <p class="text-gray-400">添加新服务器</p>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 创建服务器对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="创建新服务器"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-width="100px"
        @input="handleCreateFormInput"
      >
        <el-form-item label="服务器名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入服务器名称"/>
        </el-form-item>
        <el-form-item label="服务器描述" prop="description">
          <el-input v-model="createForm.description" type="textarea" :rows="3" placeholder="请输入服务器描述"/>
        </el-form-item>
        <el-form-item label="世界名称" prop="world_name">
          <el-input v-model="createForm.world_name" placeholder="请输入世界名称">
            <template #append>
              <el-tooltip content="不填写将使用默认世界名称" placement="top">
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="服务器版本" prop="version">
          <el-input v-model="createForm.version" placeholder="可选，不填则使用最新版本">
            <template #append>
              <el-tooltip content="不填写将使用最新版本" placement="top">
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-space>
          <el-button @click="createDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleCreateServer" :loading="creating" :disabled="!isCreateFormValid">
            创建
          </el-button>
        </el-space>
      </template>
    </el-dialog>

    <!-- 编辑服务器对话框 -->
    <el-dialog
      v-model="updateDialogVisible"
      title="更新服务器"
      width="500px"
      :close-on-click-modal="false"
      class="p-2 rounded-xl"
    >
      <el-form
        ref="updateFormRef"
        :model="updateForm"
        :rules="updateRules"
        label-width="100px"
        @input="handleUpdateFormInput"
      >
        <el-form-item label="服务器名称" prop="name">
          <el-input v-model="updateForm.name" placeholder="请输入服务器名称" />
        </el-form-item>
        <el-form-item label="服务器描述" prop="description">
          <el-input v-model="updateForm.description" type="textarea" :rows="3" placeholder="请输入服务器描述"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-space>
          <el-button @click="updateDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleUpdateServer" :loading="updating" :disabled="!isUpdateFormValid">
            保存
          </el-button>
        </el-space>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {
  QuestionFilled,
  Refresh,
  Plus
} from '@element-plus/icons-vue'
import { createServer, deleteServer, requestServerInfos, updateServer, startServer, stopServer } from '@/composables/server.js'
import { ElNotification, ElMessageBox } from 'element-plus'
import ServerCard from '@/components/widgets/ServerCard.vue'

export default {
  name: 'ServerList',
  components: {
    QuestionFilled,
    Refresh,
    Plus,
    ServerCard
  },
  data() {
    return {
      tableData: [],
      loading: false,
      // 创建服务器相关数据
      createDialogVisible: false,
      creating: false,
      createForm: {
        name: '',
        description: '',
        world_name: '',
        version: ''
      },
      createRules: {
        name: [
          { required: true, message: '请输入服务器名称', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        description: [
          { max: 200, message: '描述不能超过200个字符', trigger: 'blur' }
        ],
        world_name: [
          { max: 50, message: '世界名称不能超过50个字符', trigger: 'blur' }
        ],
        version: [
          { pattern: /^\d+\.\d+\.\d+\.\d+$/, message: '版本号格式不正确，例如：1.19.50.02', trigger: 'blur' }
        ]
      },
      // 编辑服务器相关数据
      updateDialogVisible: false,
      updating: false,
      updateForm: {
        id: '',
        name: '',
        description: ''
      },
      updateRules: {
        name: [
          { required: true, message: '请输入服务器名称', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        description: [
          { max: 200, message: '描述不能超过200个字符', trigger: 'blur' }
        ]
      },
      // 状态更新定时器
      dataTimer: null,
      isCreateFormValid: false,
      isUpdateFormValid: false
    }
  },
  async created() {
    await this.loadData()
    // 启动状态更新
    this.startStatusUpdate()
  },
  beforeUnmount() {
    // 组件销毁前清除定时器
    this.stopStatusUpdate()
  },
  methods: {
    async loadData() {
      // 只在第一次加载时显示加载状态
      if (this.tableData.length === 0) {
        this.loading = true
      }

      try {
        const response = await requestServerInfos({
          page: this.currentPage,
          size: this.pageSize,
          order: 'desc',
          order_by: 'created_at'
        })

        if (response?.data) {
          const { total, items } = response.data
          this.tableData = items || []
          this.total = total || 0
        } else {
          this.tableData = []
          this.total = 0
          ElNotification({
            title: '警告',
            message: '获取服务器列表失败',
            type: 'warning'
          })
        }
      } catch (error) {
        console.error('数据加载失败:', error)
        this.tableData = []
        this.total = 0
        ElNotification({
          title: '错误',
          message: error.response?.data?.message || '数据加载失败',
          type: 'error'
        })
      } finally {
        this.loading = false
      }
    },
    handleManageServer(server) {
      this.$router.push({
        name: 'server_detail',
        params: { id: server.id }
      }).catch(err => {
        // 处理重复导航错误
        if (err.name !== 'NavigationDuplicated') {
          console.error(err)
        }
      })
    },
    handleDeleteServer(server) {
      ElMessageBox.confirm(
        '确定要删除这个服务器吗？此操作不可恢复。',
        '删除确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(() => {
        ElNotification({
          title: '删除中',
          message: '正在删除服务器...',
          type: 'info'
        })

        deleteServer(server.id).then(() => {
          ElNotification({
            title: '删除成功',
            message: '服务器已成功删除',
            type: 'success'
          })
          this.refreshData()
        }).catch(error => {
          console.error('删除服务器失败:', error)
          ElNotification({
            title: '删除失败',
            message: error.message || '删除服务器失败',
            type: 'error'
          })
        })
      }).catch(() => {
        // 用户取消删除操作
      })
    },
    handleCreateFormInput() {
      this.$refs.createFormRef?.validate((valid) => {
        this.isCreateFormValid = valid
      })
    },
    handleUpdateFormInput() {
      this.$refs.updateFormRef?.validate((valid) => {
        this.isUpdateFormValid = valid
      })
    },
    handleShowCreateDialog() {
      this.createForm = {
        name: '',
        description: '',
        world_name: '',
        version: ''
      }
      this.isCreateFormValid = false
      this.createDialogVisible = true
      this.$nextTick(() => {
        this.$refs.createFormRef?.clearValidate()
      })
    },
    async handleCreateServer() {
      if (!this.createForm.name) {
        ElNotification({
          title: '创建失败',
          message: '请输入服务器名称',
          type: 'error'
        })
        return
      }

      this.creating = true
      try {
        await createServer({
          name: this.createForm.name,
          description: this.createForm.description,
          world_name: this.createForm.world_name,
          version: this.createForm.version
        })

        ElNotification({
          title: '创建成功',
          message: '服务器创建成功',
          type: 'success'
        })

        this.createDialogVisible = false
        this.refreshData()
      } catch (error) {
        console.error('创建服务器失败:', error)
        ElNotification({
          title: '创建失败',
          message: error.response?.data?.message || '创建服务器失败',
          type: 'error'
        })
      } finally {
        this.creating = false
      }
    },
    // 启动状态更新
    startStatusUpdate() {
      this.stopStatusUpdate() // 确保不会重复启动

      // 服务器列表更新定时器
      this.dataTimer = setInterval(() => {
        this.loadData()
      }, 10000)
    },
    // 停止状态更新
    stopStatusUpdate() {
      if (this.dataTimer) {
        clearInterval(this.dataTimer)
        this.dataTimer = null
      }
    },
    // 刷新数据
    refreshData() {
      this.loadData()
    },
    handleShowUpdateDialog(server) {
      this.updateForm = {
        id: server.id,
        name: server.name || '',
        description: server.description || ''
      }
      this.isUpdateFormValid = false
      this.updateDialogVisible = true
      this.$nextTick(() => {
        this.$refs.updateFormRef?.clearValidate()
      })
    },
    async handleUpdateServer() {
      if (!this.updateForm.name) {
        ElNotification({
          title: '更新失败',
          message: '请输入服务器名称',
          type: 'error'
        })
        return
      }

      this.updating = true
      try {
        await updateServer(this.updateForm.id, {
          name: this.updateForm.name,
          description: this.updateForm.description
        })

        ElNotification({
          title: '更新成功',
          message: '服务器信息已更新',
          type: 'success'
        })

        this.updateDialogVisible = false
        this.refreshData()
      } catch (error) {
        console.error('更新服务器失败:', error)
        ElNotification({
          title: '更新失败',
          message: error.response?.data?.message || '更新服务器失败',
          type: 'error'
        })
      } finally {
        this.updating = false
      }
    },
    async handleStartServer(server) {
      try {
        ElNotification({
          title: server.name,
          message: '正在启动服务器...',
          type: 'info'
        })

        await startServer(server.id)

        ElNotification({
          title: server.name,
          message: '服务器启动成功',
          type: 'success'
        })

        this.refreshData()
      } catch (error) {
        console.error('启动服务器失败:', error)
        ElNotification({
          title: server.name,
          message: error.response?.data?.message || '启动服务器失败',
          type: 'error'
        })
      }
    },
    async handleStopServer(server) {
      try {
        ElNotification({
          title: server.name,
          message: '正在停止服务器...',
          type: 'info'
        })

        await stopServer(server.id)

        ElNotification({
          title: server.name,
          message: '服务器停止成功',
          type: 'success'
        })

        this.refreshData()
      } catch (error) {
        console.error('停止服务器失败:', error)
        ElNotification({
          title: server.name,
          message: error.response?.data?.message || '停止服务器失败',
          type: 'error'
        })
      }
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

@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease;
}
</style>

