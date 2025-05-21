<template>
    <el-container class="h-screen bg-gray-900 text-white">
        <!-- 顶部导航 -->
        <el-header class="z-10 flex-shrink-0 p-0">
            <nav class="bg-gray-800 bg-opacity-90 backdrop-blur-sm border-b border-gray-700">
                <div class="px-4 h-16 flex items-center justify-between">
                    <div class="flex items-center">
                        <i class="fas fa-cube text-green-500 text-2xl"></i>
                        <span class="ml-2 text-xl font-bold cursor-pointer transition-colors duration-300 hover:text-green-400" @click="$router.push('/home')">MC Dashboard</span>
                    </div>
                    <div class="flex items-center space-x-4">
                        <!-- 未登录显示登录按钮 -->
                        <template v-if="!isLoggedIn">
                            <router-link to="/login" class="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-md text-sm font-medium">
                                登录
                            </router-link>
                        </template>
                        <!-- 已登录显示用户信息 -->
                        <template v-else>
                            <div class="relative">
                                <div class="flex items-center space-x-3 cursor-pointer hover:bg-gray-700 px-3 py-2 rounded-lg transition-colors duration-200" @click="toggleDropdown">
                                    <img :src="userAvatar" alt="用户头像" class="w-8 h-8 rounded-full">
                                    <div class="flex flex-col">
                                        <span class="text-sm font-medium">{{ username }}</span>
                                        <span class="text-xs text-gray-400">{{ userEmail }}</span>
                                    </div>
                                    <i class="fas fa-chevron-down text-gray-400 text-xs"></i>
                                </div>
                                <!-- 下拉菜单 -->
                                <div v-show="showDropdown" class="absolute right-0 mt-2 w-48 bg-gray-800 rounded-lg shadow-lg py-1 z-50">
                                    <a href="#" class="block px-4 py-2 text-sm text-gray-300 hover:bg-gray-700">
                                        <i class="fas fa-user mr-2"></i>个人信息
                                    </a>
                                    <a href="#" class="block px-4 py-2 text-sm text-gray-300 hover:bg-gray-700">
                                        <i class="fas fa-cog mr-2"></i>系统设置
                                    </a>
                                    <div class="border-t border-gray-700 my-1"></div>
                                    <a href="#" class="block px-4 py-2 text-sm text-gray-300 hover:bg-gray-700" @click="handleLogout">
                                        <i class="fas fa-sign-out-alt mr-2"></i>退出登录
                                    </a>
                                </div>
                            </div>
                        </template>
                    </div>
                </div>
            </nav>
        </el-header>

        <!-- 主内容区 -->
        <el-main class="flex-1 p-0">
            <router-view v-slot="{ Component }">
                <transition name="fade" mode="out-in">
                    <component :is="Component" />
                </transition>
            </router-view>
        </el-main>
    </el-container>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { getAuthInfo, logout } from '@/composables/auth'
import { ElMessage } from 'element-plus'

export default {
  name: 'ContainerPage',
  setup() {
    const router = useRouter()
    const showDropdown = ref(false)

    // 获取认证信息
    const authInfo = computed(() => getAuthInfo())
    const isLoggedIn = computed(() => !!authInfo.value.token)
    const username = computed(() => authInfo.value.user?.nickname || '用户')
    const userEmail = computed(() => authInfo.value.user?.email || '')
    const userAvatar = computed(() => {
      return `https://api.dicebear.com/7.x/avataaars/svg?seed=${username.value}`
    })

    // 切换下拉菜单显示状态
    const toggleDropdown = () => {
      showDropdown.value = !showDropdown.value
    }

    // 处理退出登录
    const handleLogout = async() => {
      try {
        await logout()
        ElMessage.success('登出成功')
        router.push('/login')
      } catch (error) {
        console.error('登出失败:', error)
        // 即使后端登出失败，也清除本地存储并跳转到登录页
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        ElMessage.error('登出失败，请重新登录')
        router.push('/login')
      }
    }

    // 点击其他地方关闭下拉菜单
    const handleClickOutside = (event) => {
      const dropdown = document.querySelector('.relative')
      if (dropdown && !dropdown.contains(event.target)) {
        showDropdown.value = false
      }
    }

    // 添加全局点击事件监听
    document.addEventListener('click', handleClickOutside)

    return {
      isLoggedIn,
      username,
      userEmail,
      userAvatar,
      showDropdown,
      toggleDropdown,
      handleLogout
    }
  }
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}
</style>
