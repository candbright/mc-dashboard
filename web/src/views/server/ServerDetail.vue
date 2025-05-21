<template>
  <div class="p-6 min-h-screen text-white">
    <!-- 返回按钮 -->
    <div class="mb-6">
      <el-button @click="back" text class="flex items-center gap-2 text-white hover:text-emerald-400">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
    </div>

    <!-- 基本信息 -->
    <el-card class="mb-6 border border-white/10 rounded-xl">
      <el-row justify="space-between" align="middle">
        <el-col>
          <h2 class="text-2xl font-semibold mb-4 text-white">{{ data.name }}</h2>
        </el-col>
        <el-col :span="24">
          <div class="mt-4">
            <el-row :gutter="20">
              <el-col :xs="24" :sm="12">
                <div class="flex items-center gap-2 mb-3">
                  <span class="text-gray-400 min-w-[100px]">服务器ID：</span>
                  <span class="text-white">{{ data.id }}</span>
                </div>
                <div class="flex items-center gap-2 mb-3">
                  <span class="text-gray-400 min-w-[100px]">服务器名称：</span>
                  <span class="text-white">{{ data.name }}</span>
                </div>
                <div class="flex items-center gap-2 mb-3">
                  <span class="text-gray-400 min-w-[100px]">版本信息：</span>
                  <span class="text-white">{{ data.version || '未知' }}</span>
                </div>
                <div class="flex items-center gap-2 mb-3">
                  <span class="text-gray-400 min-w-[100px]">运行状态：</span>
                  <el-tag :type="data.active ? 'success' : 'danger'" size="large">
                    {{ data.active ? '运行中' : '已停止' }}
                  </el-tag>
                </div>
              </el-col>
              <el-col :xs="24" :sm="12">
                <div class="flex items-center gap-2 mb-3">
                  <span class="text-gray-400 min-w-[100px]">CPU 使用率：</span>
                  <span class="text-white">45%</span>
                </div>
                <div class="flex items-center gap-2 mb-3">
                  <span class="text-gray-400 min-w-[100px]">内存使用：</span>
                  <span class="text-white">4GB/8GB</span>
                </div>
                <div class="flex items-center gap-2 mb-3">
                  <span class="text-gray-400 min-w-[100px]">在线玩家：</span>
                  <span class="text-white">{{ data.online_players || 0 }}/{{ data.max_players || 20 }}</span>
                </div>
                <div class="flex items-center gap-2 mb-3">
                  <span class="text-gray-400 min-w-[100px]">运行时间：</span>
                  <span class="text-white">3天2小时</span>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 服务器配置 -->
    <el-card class="mb-6 border border-white/10 rounded-xl">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-white">服务器配置</span>
          <el-button type="primary" @click="saveProperties" class="flex items-center gap-2">
            <el-icon><Document /></el-icon>
            <span>保存配置</span>
          </el-button>
        </div>
      </template>
      <el-form label-position="top">
        <el-tabs type="border-card" class="bg-gray-800">
          <!-- 基础设置 -->
          <el-tab-pane label="基础设置">
            <el-form-item label="服务器名称">
              <el-input v-model="data.server_properties['server-name']" class="w-full" />
            </el-form-item>
            <el-form-item label="世界名称">
              <el-input v-model="data.server_properties['level-name']" class="w-full" />
            </el-form-item>
            <el-form-item label="世界种子">
              <el-input v-model="data.server_properties['level-seed']" class="w-full" />
            </el-form-item>
            <el-form-item label="游戏模式">
              <el-select v-model="data.server_properties.gamemode" class="w-full">
                <el-option label="生存模式" value="survival" />
                <el-option label="创造模式" value="creative" />
                <el-option label="冒险模式" value="adventure" />
              </el-select>
            </el-form-item>
            <el-form-item label="难度">
              <el-select v-model="data.server_properties.difficulty" class="w-full">
                <el-option label="和平" value="peaceful" />
                <el-option label="简单" value="easy" />
                <el-option label="普通" value="normal" />
                <el-option label="困难" value="hard" />
              </el-select>
            </el-form-item>
            <el-form-item label="强制游戏模式">
              <el-switch
                v-model="data.server_properties['force-gamemode']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
          </el-tab-pane>

          <!-- 网络设置 -->
          <el-tab-pane label="网络设置">
            <el-form-item label="服务器端口">
              <el-input-number v-model="data.server_properties['server-port']" :min="1" :max="65535" class="w-full" />
            </el-form-item>
            <el-form-item label="IPv6端口">
              <el-input-number v-model="data.server_properties['server-portv6']" :min="1" :max="65535" class="w-full" />
            </el-form-item>
            <el-form-item label="最大玩家数">
              <el-input-number v-model="data.server_properties['max-players']" :min="1" :max="100" class="w-full" />
            </el-form-item>
            <el-form-item label="在线模式">
              <el-switch
                v-model="data.server_properties['online-mode']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
            <el-form-item label="启用局域网可见">
              <el-switch
                v-model="data.server_properties['enable-lan-visibility']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
            <el-form-item label="玩家空闲超时(分钟)">
              <el-input-number v-model="data.server_properties['player-idle-timeout']" :min="0" :max="1440" class="w-full" />
            </el-form-item>
          </el-tab-pane>

          <!-- 性能设置 -->
          <el-tab-pane label="性能设置">
            <el-form-item label="视距">
              <el-input-number v-model="data.server_properties['view-distance']" :min="1" :max="32" class="w-full" />
            </el-form-item>
            <el-form-item label="区块加载距离">
              <el-input-number v-model="data.server_properties['tick-distance']" :min="1" :max="12" class="w-full" />
            </el-form-item>
            <el-form-item label="最大线程数">
              <el-input-number v-model="data.server_properties['max-threads']" :min="1" :max="32" class="w-full" />
            </el-form-item>
            <el-form-item label="压缩算法">
              <el-select v-model="data.server_properties['compression-algorithm']" class="w-full">
                <el-option label="zlib" value="zlib" />
                <el-option label="snappy" value="snappy" />
              </el-select>
            </el-form-item>
            <el-form-item label="压缩阈值">
              <el-input-number v-model="data.server_properties['compression-threshold']" :min="0" :max="65535" class="w-full" />
            </el-form-item>
          </el-tab-pane>

          <!-- 安全设置 -->
          <el-tab-pane label="安全设置">
            <el-form-item label="启用白名单">
              <el-switch
                v-model="data.server_properties['allow-list']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
            <el-form-item label="默认玩家权限">
              <el-select v-model="data.server_properties['default-player-permission-level']" class="w-full">
                <el-option label="访客" value="visitor" />
                <el-option label="成员" value="member" />
                <el-option label="操作员" value="operator" />
              </el-select>
            </el-form-item>
            <el-form-item label="允许作弊">
              <el-switch
                v-model="data.server_properties['allow-cheats']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
            <el-form-item label="禁用自定义皮肤">
              <el-switch
                v-model="data.server_properties['disable-custom-skins']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
            <el-form-item label="禁用角色">
              <el-switch
                v-model="data.server_properties['disable-persona']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
          </el-tab-pane>

          <!-- 高级设置 -->
          <el-tab-pane label="高级设置">
            <el-form-item label="允许入站脚本调试">
              <el-switch
                v-model="data.server_properties['allow-inbound-script-debugging']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
            <el-form-item label="允许出站脚本调试">
              <el-switch
                v-model="data.server_properties['allow-outbound-script-debugging']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
            <el-form-item label="脚本调试器自动附加">
              <el-select v-model="data.server_properties['script-debugger-auto-attach']" class="w-full">
                <el-option label="禁用" value="disabled" />
                <el-option label="启用" value="enabled" />
              </el-select>
            </el-form-item>
            <el-form-item label="内容日志文件">
              <el-switch
                v-model="data.server_properties['content-log-file-enabled']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
            <el-form-item label="网络ID使用哈希">
              <el-switch
                v-model="data.server_properties['block-network-ids-are-hashes']"
                :active-value="'true'"
                :inactive-value="'false'"
              />
            </el-form-item>
            <el-form-item label="聊天限制">
              <el-select v-model="data.server_properties['chat-restriction']" class="w-full">
                <el-option label="无" value="None" />
                <el-option label="仅限白名单" value="Whitelist" />
                <el-option label="仅限操作员" value="Operator" />
              </el-select>
            </el-form-item>
          </el-tab-pane>
        </el-tabs>
      </el-form>
    </el-card>

    <!-- 控制台输出 -->
    <el-card class="border border-white/10 rounded-xl">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-white">控制台输出</span>
          <el-space>
            <el-button
              v-if="data.active"
              type="danger"
              @click="handleStopServer"
              class="flex items-center gap-2"
            >
              <el-icon><VideoPause /></el-icon>
              <span>停止服务器</span>
            </el-button>
            <el-button
              v-else
              type="success"
              @click="handleStartServer"
              class="flex items-center gap-2"
            >
              <el-icon><VideoPlay /></el-icon>
              <span>启动服务器</span>
            </el-button>
          </el-space>
        </div>
      </template>
      <div class="h-[400px] overflow-y-auto p-4 font-mono text-sm leading-relaxed text-gray-400 bg-black/20 rounded-lg">
        <pre class="whitespace-pre-wrap">{{ consoleLog }}</pre>
      </div>
    </el-card>
  </div>
</template>

<script>
import { addAllowList, deleteAllowList } from '@/composables/allow_list.js'
import { requestServerInfo, getConsoleLog, startServer, stopServer } from '@/composables/server.js'
import {
  ArrowLeft,
  Document,
  VideoPause,
  VideoPlay
} from '@element-plus/icons-vue'
import { setServerProperties } from '@/composables/server_properties.js'

export default {
  name: 'ServerDetail',
  components: {
    ArrowLeft,
    Document,
    VideoPause,
    VideoPlay
  },
  data() {
    return {
      activeName: 'server_info',
      data: {
        id: '',
        name: '',
        active: false,
        server_properties: {
          'server-name': '',
          'gamemode': 'survival',
          'difficulty': 'normal',
          'max-players': 20,
          'view-distance': 10
        },
        allow_list: []
      },
      loading: false,
      username: '',
      consoleLog: '',
      logTimer: null
    }
  },
  async created() {
    await this.loadData()
    this.startLogPolling()
  },
  beforeUnmount() {
    this.stopLogPolling()
  },
  methods: {
    async loadData() {
      this.loading = true
      try {
        const response = await requestServerInfo(this.$route.params.id)
        if (response?.data) {
          this.data = {
            ...this.data,
            ...response.data,
            server_properties: {
              ...this.data.server_properties,
              ...response.data.server_properties
            }
          }

          // 转换数值类型字段
          const numberFields = [
            'max-players',
            'server-port',
            'server-portv6',
            'tick-distance',
            'view-distance'
          ]

          numberFields.forEach(field => {
            if (this.data.server_properties[field]) {
              this.data.server_properties[field] =
                parseInt(this.data.server_properties[field])
            }
          })
        }
      } catch (error) {
        console.error('数据加载失败:', error)
        this.$message.error('数据加载失败')
      } finally {
        this.loading = false
      }
    },

    async saveProperties() {
      try {
        await setServerProperties(
          this.data.id,
          this.data.server_properties
        )
        this.$message.success('配置保存成功')
      } catch (error) {
        console.error('保存失败:', error)
        this.$message.error('配置保存失败')
      }
    },

    async addAllowList() {
      if (!this.username.trim()) {
        this.$message.warning('请输入有效的玩家名称')
        return
      }

      if (!this.data.active) {
        this.$message.warning('请先启动服务器后再添加白名单')
        return
      }

      try {
        await addAllowList(this.data.id, {
          name: this.username,
          xuid: '', // 需要实际获取XUID的逻辑
          ignoresPlayerLimit: false
        })

        this.data.allow_list.push({
          name: this.username,
          xuid: '',
          ignoresPlayerLimit: false
        })

        this.username = ''
        this.$message.success('添加成功')
      } catch (error) {
        console.error('添加失败:', error)
        this.$message.error('添加失败')
      }
    },

    async deleteAllowList() {
      if (!this.username.trim()) {
        this.$message.warning('请输入要删除的玩家名称')
        return
      }

      if (!this.data.active) {
        this.$message.warning('请先启动服务器后再删除白名单')
        return
      }

      try {
        await deleteAllowList(this.data.id, this.username)
        this.data.allow_list = this.data.allow_list.filter(
          item => item.name !== this.username
        )
        this.username = ''
        this.$message.success('删除成功')
      } catch (error) {
        console.error('删除失败:', error)
        this.$message.error('删除失败')
      }
    },

    handleRemove(index) {
      if (!this.data.active) {
        this.$message.warning('请先启动服务器后再删除白名单')
        return
      }
      const item = this.data.allow_list[index]
      this.username = item.name
      this.deleteAllowList()
    },

    back() {
      this.$router.go(-1)
    },

    handleTabClick(tab) {
      console.log('切换到标签:', tab.props.name)
    },

    startLogPolling() {
      this.logTimer = setInterval(async () => {
        if (this.data.active) {
          await this.fetchConsoleLog()
        }
      }, 2000) // 每2秒更新一次
    },

    stopLogPolling() {
      if (this.logTimer) {
        clearInterval(this.logTimer)
        this.logTimer = null
      }
    },

    async fetchConsoleLog() {
      try {
        const response = await getConsoleLog(this.data.id)
        if (response?.data?.content) {
          this.consoleLog = response.data.content
        }
      } catch (error) {
        console.error('获取控制台日志失败:', error)
      }
    },

    async handleStartServer() {
      try {
        await startServer(this.data.id)
        this.data.active = true
        this.startLogPolling()
        this.$message.success('服务器启动成功')
      } catch (error) {
        console.error('启动失败:', error)
        this.$message.error('服务器启动失败')
      }
    },

    async handleStopServer() {
      try {
        await stopServer(this.data.id)
        this.data.active = false
        this.stopLogPolling()
        this.$message.success('服务器停止成功')
      } catch (error) {
        console.error('停止失败:', error)
        this.$message.error('服务器停止失败')
      }
    }
  }
}
</script>

<style scoped>
:deep(.el-tabs__item) {
  @apply text-gray-400;
}

:deep(.el-tabs__item.is-active) {
  @apply text-emerald-400;
}

:deep(.el-tabs__active-bar) {
  @apply bg-emerald-400;
}

:deep(.el-tabs__nav-wrap::after) {
  @apply bg-white/10;
}

:deep(.el-form-item__label) {
  @apply text-gray-400;
}

:deep(.el-input__wrapper),
:deep(.el-textarea__wrapper) {
  @apply bg-gray-800 border-white/10;
}

:deep(.el-input__inner),
:deep(.el-textarea__inner) {
  @apply text-white;
}

:deep(.el-input-number__wrapper) {
  @apply bg-gray-800 border-white/10;
}

:deep(.el-input-number__decrease),
:deep(.el-input-number__increase) {
  @apply bg-gray-700 border-white/10 text-white;
}

:deep(.el-select .el-input__wrapper) {
  @apply bg-gray-800 border-white/10;
}

:deep(.el-select-dropdown__item) {
  @apply text-gray-400;
}

:deep(.el-select-dropdown__item.selected) {
  @apply text-emerald-400;
}

:deep(.el-switch__core) {
  @apply border-white/10;
}

:deep(.el-switch.is-checked .el-switch__core) {
  @apply bg-emerald-500 border-emerald-500;
}
</style>
