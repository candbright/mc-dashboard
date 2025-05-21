<template>
  <div class="text-white">
    <!-- 背景图 -->
    <div class="fixed inset-0 bg-cover bg-center bg-no-repeat" style="background-image: url('https://images.unsplash.com/photo-1587573089734-599851b2c3b5?ixlib=rb-1.2.1&auto=format&fit=crop&w=1920&q=80'); filter: brightness(0.3);"></div>

    <!-- 内容区 -->
    <div class="relative min-h-screen flex flex-col">

      <!-- 主要内容 -->
      <main class="flex-grow">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <!-- 标题和描述 -->
          <div class="text-center mb-16">
            <h1 class="text-5xl font-bold mb-6">Minecraft 服务器管理系统</h1>
            <p class="text-xl text-gray-300 max-w-2xl mx-auto">
              一个强大的 Minecraft 服务器管理平台，让您轻松管理多个服务器、存档和玩家
            </p>
          </div>

          <!-- 功能展示 -->
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <!-- 服务器管理 -->
            <div class="bg-gray-800 bg-opacity-90 backdrop-blur-sm rounded-lg p-6 hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
              <div class="w-12 h-12 bg-green-500 bg-opacity-20 rounded-lg flex items-center justify-center mb-4">
                <i class="fas fa-server text-2xl text-green-500"></i>
              </div>
              <h3 class="text-xl font-semibold mb-2">服务器管理</h3>
              <p class="text-gray-400 mb-4">
                轻松管理多个 Minecraft 服务器，实时监控状态，一键启动/停止
              </p>
              <router-link to="/server/list" class="inline-flex items-center text-green-500 hover:text-green-400">
                了解更多 <i class="fas fa-arrow-right ml-2"></i>
              </router-link>
            </div>

            <!-- 存档管理 -->
            <div class="bg-gray-800 bg-opacity-90 backdrop-blur-sm rounded-lg p-6 hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
              <div class="w-12 h-12 bg-green-500 bg-opacity-20 rounded-lg flex items-center justify-center mb-4">
                <i class="fas fa-archive text-2xl text-green-500"></i>
              </div>
              <h3 class="text-xl font-semibold mb-2">存档管理</h3>
              <p class="text-gray-400 mb-4">
                便捷的存档备份和恢复功能，确保您的世界数据安全
              </p>
              <router-link to="/save/list" class="inline-flex items-center text-green-500 hover:text-green-400">
                了解更多 <i class="fas fa-arrow-right ml-2"></i>
              </router-link>
            </div>

            <!-- 玩家管理 -->
            <div class="bg-gray-800 bg-opacity-90 backdrop-blur-sm rounded-lg p-6 hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
              <div class="w-12 h-12 bg-green-500 bg-opacity-20 rounded-lg flex items-center justify-center mb-4">
                <i class="fas fa-users text-2xl text-green-500"></i>
              </div>
              <h3 class="text-xl font-semibold mb-2">玩家管理</h3>
              <p class="text-gray-400 mb-4">
                简单的白名单管理，轻松控制服务器访问权限
              </p>
              <router-link to="/whitelist" class="inline-flex items-center text-green-500 hover:text-green-400">
                了解更多 <i class="fas fa-arrow-right ml-2"></i>
              </router-link>
            </div>
          </div>

          <!-- 开始使用按钮 -->
          <div class="text-center mt-16">
            <router-link
              to="/login"
              class="inline-flex items-center px-8 py-4 bg-green-500 hover:bg-green-600 rounded-lg text-white font-medium text-lg transition-colors duration-300"
            >
              <i class="fas fa-rocket mr-2"></i>
              立即开始使用
            </router-link>
          </div>
        </div>
      </main>

      <!-- 页脚 -->
      <footer class="bg-gray-800 bg-opacity-90 backdrop-blur-sm border-t border-gray-700">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div class="text-center text-gray-400 text-sm">
            <p>© 2024 MC Dashboard. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomePage'
}
</script>

<style scoped>
/* 添加页面过渡动画 */
.page-enter-active,
.page-leave-active {
  transition: opacity 0.1s ease;
}

.page-enter-from,
.page-leave-to {
  opacity: 0;
}
</style>
