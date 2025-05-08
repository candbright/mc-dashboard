<template>
  <div>
    <el-row align="middle" justify="start">
      <el-col :span="1">
        <el-button
          :icon="ArrowLeft"
          text
          circle
          @click="back"
          style="margin-bottom: 0; vertical-align: middle; color: #000000"
          size="large"
        />
      </el-col>
      <el-col :span="4">
        <h3>{{ data.name }}</h3>
      </el-col>
    </el-row>
    <el-tabs
      v-model="activeName"
      type="card"
      class="demo-tabs"
      @tab-click="handleTabClick"
    >
      <el-tab-pane label="服务器详情" name="server_info"></el-tab-pane>
      <el-tab-pane label="参数配置" name="server_properties"></el-tab-pane>
      <el-tab-pane label="白名单" name="allow_list"></el-tab-pane>
    </el-tabs>

    <!-- 服务器详情 -->
    <div v-if="activeName === 'server_info'" class="content-box">
      <h3>服务器详情</h3>
      <div v-loading="loading">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="服务器ID">{{ data.id }}</el-descriptions-item>
          <el-descriptions-item label="服务器名称">{{ data.name }}</el-descriptions-item>
          <el-descriptions-item label="版本号">{{ data.version }}</el-descriptions-item>
          <el-descriptions-item label="运行状态">
            <el-tag :type="!!data.active">
              {{ data.active ? '运行中' : '已停止' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="是否存在实例">
            <el-tag :type="!!data.exist">
              {{ data.exist ? '已创建' : '未创建' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="下载状态">
            <el-tag :type="!!data.downloading">
              {{ data.downloading ? '下载中' : '未下载' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>

    <!-- 参数配置 -->
    <div v-if="activeName === 'server_properties'" class="content-box">
      <h3>服务器配置</h3>
      <div v-loading="loading">
        <el-form label-width="220px" v-if="data">
          <el-form-item label="允许作弊">
            <el-switch v-model="data.server_properties['allow-cheats']"
                       active-value="true"
                       inactive-value="false"/>
          </el-form-item>

          <el-form-item label="游戏模式">
            <el-select v-model="data.server_properties.gamemode">
              <el-option label="生存模式" value="survival"/>
              <el-option label="创造模式" value="creative"/>
              <el-option label="冒险模式" value="adventure"/>
            </el-select>
          </el-form-item>

          <el-form-item label="最大玩家数">
            <el-input-number
              v-model="data.server_properties['max-players']"
              :min="1"
              :max="100"
              controls-position="right"/>
          </el-form-item>

          <el-form-item label="难度等级">
            <el-select v-model="data.server_properties.difficulty">
              <el-option label="和平" value="peaceful"/>
              <el-option label="简单" value="easy"/>
              <el-option label="普通" value="normal"/>
              <el-option label="困难" value="hard"/>
            </el-select>
          </el-form-item>

          <el-form-item label="视距距离">
            <el-input-number
              v-model="data.server_properties['view-distance']"
              :min="4"
              :max="64"
              controls-position="right"/>
          </el-form-item>

          <el-form-item label="白名单功能">
            <el-switch v-model="data.server_properties['allow-list']"
                       active-value="true"
                       inactive-value="false"/>
          </el-form-item>

          <el-form-item label="服务器端口">
            <el-input
              v-model="data.server_properties['server-port']"
              style="width: 200px"/>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="saveProperties">保存配置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <!-- 白名单管理 -->
    <div v-if="activeName === 'allow_list'" class="content-box">
      <h3>白名单管理</h3>
      <div v-loading="loading">
        <div style="margin-bottom: 20px">
          <el-input
            v-model="username"
            placeholder="请输入玩家名称"
            style="width: 240px"
            @keyup.enter="addAllowList"/>
          <el-button type="primary" @click="addAllowList">添加</el-button>
          <el-button type="danger" @click="deleteAllowList">移除</el-button>
        </div>

        <el-table :data="data.allow_list" border>
          <el-table-column prop="name" label="玩家名称" width="200"/>
          <el-table-column prop="xuid" label="XUID"/>
          <el-table-column label="绕过限制">
            <template #default="{ row }">
              <el-tag :type="row.ignoresPlayerLimit ? 'success' : 'info'">
                {{ row.ignoresPlayerLimit ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="{ $index }">
              <el-button
                type="danger"
                size="small"
                @click="handleRemove($index)">删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script>
import { addAllowList, deleteAllowList } from '@/api/mc/allow_list.js'
import { requestServerInfo } from '@/api/mc/server.js'
import { ArrowLeft } from '@element-plus/icons-vue'
import { setServerProperties } from '@/api/mc/server_properties.js'

export default {
  name: 'ServerListItem',
  data() {
    return {
      activeName: 'server_info',
      data: {
        server_properties: {},
        allow_list: []
      },
      loading: false,
      username: ''
    }
  },
  computed: {
    ArrowLeft() {
      return ArrowLeft
    }
  },
  async created() {
    await this.loadData()
  },
  methods: {
    async loadData() {
      this.loading = true
      try {
        const response = await requestServerInfo(this.$route.params.id)
        this.data = response.data

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
              parseInt(this.data.server_properties[field]);
          }
        });
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
        );
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

      try {
        await addAllowList(this.data.id, {
          name: this.username,
          xuid: '', // 需要实际获取XUID的逻辑
          ignoresPlayerLimit: false
        });

        this.data.allow_list.push({
          name: this.username,
          xuid: '',
          ignoresPlayerLimit: false
        });

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

      try {
        await deleteAllowList(this.data.id, this.username)
        this.data.allow_list = this.data.allow_list.filter(
          item => item.name !== this.username
        );
        this.username = ''
        this.$message.success('删除成功')
      } catch (error) {
        console.error('删除失败:', error)
        this.$message.error('删除失败')
      }
    },

    handleRemove(index) {
      const item = this.data.allow_list[index]
      this.username = item.name
      this.deleteAllowList()
    },

    back() {
      this.$router.go(-1)
    },

    handleTabClick(tab) {
      console.log('切换到标签:', tab.props.name)
    }
  }
};
</script>

<style scoped>
.content-box {
  padding: 20px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  margin-top: 20px;
  background: #fff;
}

.el-form {
  max-width: 800px;
}

.el-input-number {
  width: 200px;
}
</style>
