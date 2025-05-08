<template>
  <div class="apple-container" :class="{ 'dark-theme': isDarkTheme, 'light-theme': !isDarkTheme }">
    <nav class="apple-nav">
      <div class="nav-content">
        <div class="nav-items">
          <div class="nav-item dropdown" v-for="(item, index) in getMenuItems()" :key="index">
            <router-link
              v-if="item.path"
              :to="{ name: item.path }"
              class="nav-link">
              {{ item.name }}
            </router-link>
            <span
              v-else
              class="nav-link"
              @click="handleMainClick(item)">
              {{ item.name }}
            </span>
            <div class="mega-menu" v-if="item.children">
              <div class="mega-menu-content">
                <div class="mega-menu-section">
                  <div class="mega-menu-items">
                    <router-link
                      v-for="(child, childIndex) in item.children"
                      :key="childIndex"
                      :to="{ name: child.path }"
                      class="mega-menu-item">
                      <span class="item-title">{{ child.name }}</span>
                      <span class="item-description">{{ child.description || '点击查看详情' }}</span>
                    </router-link>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="theme-toggle" @click="toggleTheme">
            <el-icon v-if="isDarkTheme" color="#f5f5f7"><Moon /></el-icon>
            <el-icon v-else color="#1d1d1f"><Sunny /></el-icon>
          </div>
        </div>
      </div>
    </nav>
    <main class="apple-main">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script>
import { Moon, Sunny } from '@element-plus/icons-vue'

export default {
  name: 'MainLayout',
  components: { Sunny, Moon },

  data() {
    return {
      isDarkTheme: false
    }
  },

  created() {
    // 从localStorage读取主题设置
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme) {
      this.isDarkTheme = savedTheme === 'dark'
    } else {
      // 检查系统主题偏好
      this.isDarkTheme = window.matchMedia('(prefers-color-scheme: dark)').matches
    }
  },

  methods: {
    getMenuItems() {
      return [
        {
          'name': '我的世界',
          'path': 'home',
          'children': [
            {
              'name': '服务器列表',
              'path': 'server_list',
              'description': '查看和管理所有可用的服务器'
            },
            {
              'name': '存档列表',
              'path': 'save_list',
              'description': '管理你的游戏存档和备份'
            }
          ]
        },
        {
          'name': '未开发',
          'path': 'impl'
        }
      ]
    },

    handleMainClick(item) {
      if (item.path) {
        this.$router.push({ name: item.path })
      }
    },

    toggleTheme() {
      this.isDarkTheme = !this.isDarkTheme
      localStorage.setItem('theme', this.isDarkTheme ? 'dark' : 'light')
    }
  }
}
</script>

<style scoped>
.apple-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-color);
  transition: all 0.3s;
}

.apple-nav {
  position: fixed;
  top: 0;
  width: 100%;
  z-index: 1000;
  background: var(--nav-bg);
  padding: 0;
}

.nav-items {
  display: flex;
  justify-content: center;
  height: 44px;
  padding: 0 20px;
}

.mega-menu {
  position: fixed;
  top: 44px;
  width: 100%;
  background: var(--menu-bg);
  backdrop-filter: saturate(180%) blur(20px);
  opacity: 0;
  visibility: hidden;
  transform: translateY(10px);
  transition: all 0.3s;
  box-shadow: var(--menu-shadow);
}

.apple-main {
  height: calc(100vh - 44px);
  width: 100vw;
  margin-top: 44px;
  transition: all 0.3s;
}

.light-theme {
  --bg-color: #ffffff;
  --nav-bg: #f5f5f7;
  --menu-bg: rgba(255, 255, 255, 0.98);
  --menu-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.dark-theme {
  --bg-color: #2c2c2c;
  --nav-bg: #1d1d1f;
  --menu-bg: rgba(29, 29, 31, 0.98);
  --menu-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.nav-item {
  position: relative;
  padding: 0 20px;
  height: 100%;
  display: flex;
  align-items: center;
}

.nav-link {
  color: #1d1d1f;
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
  transition: color 0.3s ease;
  display: inline-block;
}

.dark-theme .nav-link {
  color: #f5f5f7;
}

.nav-link:hover {
  color: #000000;
}

.dark-theme .nav-link:hover {
  color: #ffffff;
}

.mega-menu {
  position: fixed;
  top: 44px;
  left: 0;
  right: 0;
  width: 100vw;
  padding: 0;
  margin: 0;
  background-color: rgba(255, 255, 255, 0.98);
  backdrop-filter: saturate(180%) blur(20px);
  opacity: 0;
  visibility: hidden;
  transform: translateY(10px);
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  z-index: 999;
}

.dark-theme .mega-menu {
  background-color: rgba(29, 29, 31, 0.98);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.dropdown:hover .mega-menu {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
}

.mega-menu-content {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.mega-menu-section {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px 0;
}

.mega-menu-section h3 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 20px;
  color: #1d1d1f;
}

.dark-theme .mega-menu-section h3 {
  color: #f5f5f7;
}

.mega-menu-items {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.mega-menu-item {
  display: block;
  padding: 20px;
  text-decoration: none;
  border-radius: 12px;
  transition: background-color 0.3s ease;
}

.mega-menu-item:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.dark-theme .mega-menu-item:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.item-title {
  display: block;
  font-size: 18px;
  font-weight: 500;
  color: #1d1d1f;
  margin-bottom: 8px;
}

.dark-theme .item-title {
  color: #f5f5f7;
}

.item-description {
  display: block;
  font-size: 14px;
  color: #86868b;
}

.dark-theme .item-description {
  color: #a1a1a6;
}

.theme-toggle {
  display: grid;
  justify-content: center;
  align-items: center;
  background: none;
  border: none;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.theme-toggle:hover {
  transform: scale(1.1);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
