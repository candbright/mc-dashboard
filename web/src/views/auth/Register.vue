<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-900">
    <div class="max-w-md w-full space-y-8 p-8 bg-gray-800 rounded-lg shadow-lg page-transition">
      <!-- Logo -->
      <div class="text-center">
        <i class="fas fa-cube text-green-500 text-5xl"></i>
        <h2 class="mt-6 text-3xl font-bold text-white">注册</h2>
        <p class="mt-2 text-sm text-gray-400">创建您的 Minecraft 服务器管理账号</p>
      </div>

      <!-- 注册表单 -->
      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
        <div class="rounded-md shadow-sm space-y-4">
          <!-- 手机号输入 -->
          <div>
            <label for="phone" class="block text-sm font-medium text-gray-300">手机号</label>
            <div class="mt-1">
              <input
                id="phone"
                v-model="form.phone"
                type="tel"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-600 rounded-md shadow-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-green-500 focus:border-green-500"
                placeholder="请输入手机号"
              >
            </div>
          </div>

          <!-- 验证码输入 -->
          <div>
            <label for="verificationCode" class="block text-sm font-medium text-gray-300">验证码</label>
            <div class="mt-1 flex space-x-2">
              <input
                id="verificationCode"
                v-model="form.verificationCode"
                type="text"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-600 rounded-md shadow-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-green-500 focus:border-green-500"
                placeholder="请输入验证码"
              >
              <button
                type="button"
                @click="sendVerificationCode"
                :disabled="sendingCode || countdown > 0"
                class="px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-green-500 hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ sendingCode ? '发送中...' : countdown > 0 ? `${countdown}秒后重试` : '获取验证码' }}
              </button>
            </div>
          </div>

          <!-- 密码输入 -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-300">密码</label>
            <div class="mt-1">
              <input
                id="password"
                v-model="form.password"
                type="password"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-600 rounded-md shadow-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-green-500 focus:border-green-500"
                placeholder="请输入密码"
              >
            </div>
          </div>

          <!-- 确认密码 -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-300">确认密码</label>
            <div class="mt-1">
              <input
                id="confirmPassword"
                v-model="form.confirmPassword"
                type="password"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-600 rounded-md shadow-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-green-500 focus:border-green-500"
                placeholder="请再次输入密码"
              >
            </div>
          </div>
        </div>

        <!-- 注册按钮 -->
        <div>
          <button
            type="submit"
            :disabled="loading"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-500 hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? '注册中...' : '注册' }}
          </button>
        </div>

        <!-- 登录链接 -->
        <div class="text-center">
          <span class="text-sm text-gray-400">已有账号？</span>
          <router-link to="/login" class="ml-1 text-sm font-medium text-green-500 hover:text-green-400">
            立即登录
          </router-link>
        </div>

        <!-- 错误提示 -->
        <div v-if="error" class="mt-2 text-sm text-red-500 text-center">
          {{ error }}
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { register, sendVerificationCode } from '@/composables/auth.js'

export default {
  name: 'RegisterView',
  data() {
    return {
      form: {
        phone: '',
        verificationCode: '',
        password: '',
        confirmPassword: ''
      },
      loading: false,
      sendingCode: false,
      countdown: 0,
      error: ''
    }
  },
  methods: {
    async handleRegister() {
      try {
        // 验证两次密码是否一致
        if (this.form.password !== this.form.confirmPassword) {
          this.error = '两次输入的密码不一致'
          return
        }

        this.loading = true
        this.error = ''

        await register({
          phone: this.form.phone,
          password: this.form.password,
          code: this.form.verificationCode
        })

        // 注册成功后跳转到登录页
        this.$router.push('/login')
      } catch (error) {
        this.error = error.response?.data?.message || '注册失败，请稍后重试'
      } finally {
        this.loading = false
      }
    },

    async sendVerificationCode() {
      if (!this.form.phone) {
        this.error = '请输入手机号'
        return
      }

      try {
        this.sendingCode = true
        this.error = ''

        await sendVerificationCode({
          phone: this.form.phone
        })

        // 开始倒计时
        this.countdown = 60
        const timer = setInterval(() => {
          this.countdown--
          if (this.countdown <= 0) {
            clearInterval(timer)
          }
        }, 1000)
      } catch (error) {
        this.error = error.response?.data?.message || '发送验证码失败，请稍后重试'
      } finally {
        this.sendingCode = false
      }
    }
  }
}
</script>

<style>
/* 页面过渡动画 */
.page-transition {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 路由过渡动画 */
.page-enter-active,
.page-leave-active {
  transition: all 0.3s ease;
}

.page-enter-from,
.page-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
