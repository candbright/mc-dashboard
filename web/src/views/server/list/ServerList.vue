<template>
  <div class="server-list-container">
    <div class="page-header">
      <h1>服务器列表</h1>
      <p class="subtitle">管理你的 Minecraft 服务器</p>
    </div>

    <div class="server-list-content">
      <el-table :data="tableData" v-loading="loading" border style="width: 100%"
                :highlight-current-row="false" @row-dblclick="handleRowDblClick">
        <el-table-column prop="id" label="ID" width="%10" align="center"></el-table-column>
        <el-table-column prop="name" label="服务器名称" width="%20" align="center"></el-table-column>
        <el-table-column prop="version" label="版本" width="%20" align="center">
          <template #default="{ row }">
            <div>
              <div v-if="row.version === ''">暂无</div>
              <div v-else>{{ row.version }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="%30" align="center">
          <template #default="{ row }">
            <div @dblclick.stop>
              <el-row justify="space-around">
                <el-col :span="2" v-if="row.downloading">
                  <el-button type="primary" loading circle/>
                </el-col>
                <el-col :span="2" v-if="!row.downloading && row.version !== ''">
                  <el-button type="primary" :icon="Upload" circle @click="upGrade"/>
                </el-col>
                <el-col :span="2" v-if="!row.downloading && row.version === ''">
                  <el-button type="primary" :loading="row.downloading" :icon="Download" circle
                             @click="downloadLatestVersion(row.id)"/>
                </el-col>
                <el-col :span="2" v-if="!row.active">
                  <el-button type="success" :disabled="!row.exist" :icon="CaretRight"
                             circle
                             @click="startServer(row)"/>
                </el-col>
                <el-col :span="2" v-if="row.active">
                  <el-button type="danger" :disabled="!row.exist" :icon="Close"
                             circle
                             @click="stopServer(row)"/>
                </el-col>
                <el-col :span="2">
                  <el-button type="primary" :disabled="!row.exist" :icon="Edit"
                             circle
                             @click="editServerInfo(row.id)"/>
                </el-col>
                <el-col :span="2">
                  <el-button type="danger" :disabled="!row.exist" :icon="Delete"
                             circle
                             @click="deleteServer(row.id)"/>
                </el-col>
              </el-row>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :background="true"
          prev-text="上一页"
          next-text="下一页"
          total-text="总计"
          :page-size-text="'条/页'"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { CaretRight, Close, Delete, Download, Edit, Upload } from '@element-plus/icons-vue'
import { downloadLatestServer, requestServerInfo, requestServerInfos, startServer, stopServer } from '@/api/mc/server.js'
import { ElNotification } from 'element-plus'

export default {
  name: 'ServerList',
  computed: {
    Close() {
      return Close
    },
    CaretRight() {
      return CaretRight
    },
    Edit() {
      return Edit
    },
    Delete() {
      return Delete
    },
    Download() {
      return Download
    },
    Upload() {
      return Upload
    }
  },
  data() {
    return {
      tableData: [],
      loading: false,
      currentPage: 1,
      pageSize: 10,
      total: 0
    }
  },
  async created() {
    await this.loadData()
  },
  methods: {
    async loadData() {
      this.loading = true
      try {
        const response = await requestServerInfos({
          page: this.currentPage,
          size: this.pageSize
        })
        this.tableData = response.data.items
        this.total = response.data.total
      } catch (error) {
        console.error('数据加载失败:', error)
        this.$message.error('数据加载失败')
      } finally {
        this.loading = false
      }
    },
    handleSizeChange(val) {
      this.pageSize = val
      this.loadData()
    },
    handleCurrentChange(val) {
      this.currentPage = val
      this.loadData()
    },
    refreshData() {
      requestServerInfos({
        page: this.currentPage,
        size: this.pageSize
      }).then(response => {
        this.tableData = response.data.items
        this.total = response.data.total
      }).catch(error => {
        console.error('数据刷新失败:', error)
        this.$message.error('数据刷新失败')
      })
    },
    handleRowDblClick(row, column, event) {
      this.$router.push({ name: 'server_info', params: { id: row.id }}).catch(err => {
        // 处理重复导航错误
        if (err.name !== 'NavigationDuplicated') {
          console.error(err)
        }
      })
    },
    downloadLatestVersion(id) {
      requestServerInfo(id).then(response => {
        if (response.data.downloading) {
          ElNotification({
            title: '下载进度',
            message: '正在下载中，请勿重复下载',
            type: 'warning',
            duration: 0
          })
          return
        }
        ElNotification({
          title: '下载进度',
          message: '开始下载',
          type: 'info',
          duration: 0
        })

        downloadLatestServer(id)

        let time = 6
        let success = false
        const timer = setInterval(() => {
          requestServerInfo(id).then(response => {
            if (!response.data.downloading && response.data.version !== '') {
              success = true
              time = 0
            } else {
              time--
            }
            if (time === 0) {
              if (success) {
                ElNotification({
                  title: '下载进度',
                  message: '下载成功',
                  type: 'success',
                  duration: 0
                })
              } else {
                ElNotification({
                  title: '下载进度',
                  message: '更新超时',
                  type: 'error',
                  duration: 0
                })
              }
              clearInterval(timer)
              this.refreshData()
            }
          })
        }, 2000)
      }).catch(error => {
        console.error('下载失败:', error)
        ElNotification({
          title: '下载进度',
          message: '下载失败，获取服务器信息失败',
          type: 'error',
          duration: 0
        })
      })
    },
    upGrade() {

    },
    editServerInfo(row) {

    },
    startServer(row) {
      ElNotification({
        title: row.name,
        message: '开启服务器',
        type: 'info',
        duration: 0
      })
      startServer(row.id).then(response => {
        let time = 6
        let success = false
        const timer = setInterval(() => {
          requestServerInfo(row.id).then(response => {
            if (response.active) {
              success = true
              time = 0
            } else {
              time--
            }
            if (time === 0) {
              if (success) {
                ElNotification({
                  title: row.name,
                  message: '开启服务器成功',
                  type: 'success',
                  duration: 0
                })
              } else {
                ElNotification({
                  title: row.name,
                  message: '更新超时',
                  type: 'error',
                  duration: 0
                })
              }
              clearInterval(timer)
              this.refreshData()
            }
          })
        }, 2000)
      })
    },
    stopServer(row) {
      ElNotification({
        title: row.name,
        message: '关闭服务器',
        type: 'info',
        duration: 0
      })
      stopServer(row.id).then(response => {
        let time = 6
        let success = false
        const timer = setInterval(() => {
          requestServerInfo(row.id).then(response => {
            if (!response.active) {
              success = true
              time = 0
            } else {
              time--
            }
            if (time === 0) {
              if (success) {
                ElNotification({
                  title: row.name,
                  message: '关闭服务器成功',
                  type: 'success',
                  duration: 0
                })
              } else {
                ElNotification({
                  title: row.name,
                  message: '更新超时',
                  type: 'error',
                  duration: 0
                })
              }
              clearInterval(timer)
              this.refreshData()
            }
          })
        }, 2000)
      })
    },
    deleteServer(id) {

    }
  }
}
</script>

<style scoped>
.server-list-container {
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

.server-list-content {
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
