<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-900">
    <div class="max-w-md w-full space-y-8 p-8 bg-gray-800 rounded-lg shadow-lg page-transition">
      <!-- Logo -->
      <div class="text-center">
        <i class="fas fa-cube text-green-500 text-5xl"></i>
        <h2 class="mt-6 text-3xl font-bold text-white">登录</h2>
        <p class="mt-2 text-sm text-gray-400">欢迎使用 Minecraft 服务器管理面板</p>
      </div>

      <!-- 登录表单 -->
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="rounded-md shadow-sm space-y-4">
          <!-- 登录方式选择 -->
          <div>
            <label class="block text-sm font-medium text-gray-300">登录方式</label>
            <div class="mt-1 grid grid-cols-3 gap-2">
              <button
                type="button"
                :class="[
                  'px-3 py-2 rounded-md text-sm font-medium',
                  form.loginType === 'phone_password'
                    ? 'bg-green-500 text-white'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                ]"
                @click="handleLoginTypeChange('phone_password')"
              >
                手机密码
              </button>
              <button
                type="button"
                :class="[
                  'px-3 py-2 rounded-md text-sm font-medium',
                  form.loginType === 'email_password'
                    ? 'bg-green-500 text-white'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                ]"
                @click="handleLoginTypeChange('email_password')"
              >
                邮箱密码
              </button>
              <button
                type="button"
                :class="[
                  'px-3 py-2 rounded-md text-sm font-medium',
                  form.loginType === 'phone_code'
                    ? 'bg-green-500 text-white'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                ]"
                @click="handleLoginTypeChange('phone_code')"
              >
                手机验证码
              </button>
            </div>
          </div>

          <!-- 账号输入 -->
          <div>
            <label for="account" class="block text-sm font-medium text-gray-300">
              {{ form.loginType === 'email_password' ? '邮箱' : '手机号' }}
            </label>
            <div class="mt-1">
              <input
                id="account"
                v-model="form.account"
                :type="form.loginType === 'email_password' ? 'email' : 'tel'"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-600 rounded-md shadow-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-green-500 focus:border-green-500"
                :placeholder="form.loginType === 'email_password' ? '请输入邮箱' : '请输入手机号'"
                @input="validateAccount"
              >
              <p v-if="errors.account" class="mt-1 text-sm text-red-500">{{ errors.account }}</p>
            </div>
          </div>

          <!-- 密码/验证码输入 -->
          <div>
            <label for="credential" class="block text-sm font-medium text-gray-300">
              {{ form.loginType === 'phone_code' ? '验证码' : '密码' }}
            </label>
            <div class="mt-1 flex space-x-2">
              <input
                id="credential"
                v-model="form.credential"
                :type="form.loginType === 'phone_code' ? 'text' : 'password'"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-600 rounded-md shadow-sm bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-green-500 focus:border-green-500"
                :placeholder="form.loginType === 'phone_code' ? '请输入验证码' : '请输入密码'"
                @input="validateCredential"
              >
              <button
                v-if="form.loginType === 'phone_code'"
                type="button"
                :disabled="countdown > 0 || !isAccountValid"
                class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-500 hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
                @click="handleSendCode"
              >
                {{ countdown > 0 ? `${countdown}秒后重试` : '获取验证码' }}
              </button>
            </div>
            <p v-if="errors.credential" class="mt-1 text-sm text-red-500">{{ errors.credential }}</p>
          </div>
        </div>

        <!-- 记住我和忘记密码 -->
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember-me"
              v-model="form.rememberMe"
              type="checkbox"
              class="h-4 w-4 text-green-500 focus:ring-green-500 border-gray-600 rounded bg-gray-700"
            >
            <label for="remember-me" class="ml-2 block text-sm text-gray-300">记住我</label>
          </div>

          <div class="text-sm">
            <a href="#" class="font-medium text-green-500 hover:text-green-400">忘记密码？</a>
          </div>
        </div>

        <!-- 登录按钮 -->
        <div>
          <button
            type="submit"
            :disabled="loading || !isFormValid"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-500 hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? '登录中...' : '登录' }}
          </button>
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
import { login, getUserInfo } from '@/composables/auth.js'
import { ElMessage } from 'element-plus'

export default {
  name: 'LoginView',
  data() {
    return {
      form: {
        loginType: 'phone_password',
        account: '',
        credential: '',
        rememberMe: false
      },
      errors: {
        account: '',
        credential: ''
      },
      loading: false,
      error: '',
      countdown: 0,
      isAccountValid: false
    }
  },
  computed: {
    isFormValid() {
      return this.isAccountValid &&
             this.form.credential &&
             !this.errors.account &&
             !this.errors.credential
    }
  },
  methods: {
    handleLoginTypeChange(type) {
      this.form.loginType = type
      this.form.account = ''
      this.form.credential = ''
      this.errors.account = ''
      this.errors.credential = ''
      this.isAccountValid = false
    },
    validateAccount() {
      const { account, loginType } = this.form
      this.errors.account = ''

      if (!account) {
        this.errors.account = '请输入账号'
        this.isAccountValid = false
        return
      }

      if (loginType === 'email_password') {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        if (!emailRegex.test(account)) {
          this.errors.account = '请输入有效的邮箱地址'
          this.isAccountValid = false
          return
        }
      } else {
        const phoneRegex = /^1[3-9]\d{9}$/
        if (!phoneRegex.test(account)) {
          this.errors.account = '请输入有效的手机号'
          this.isAccountValid = false
          return
        }
      }

      this.isAccountValid = true
    },
    validateCredential() {
      const { credential, loginType } = this.form
      this.errors.credential = ''

      if (!credential) {
        this.errors.credential = loginType === 'phone_code' ? '请输入验证码' : '请输入密码'
        return
      }

      if (loginType !== 'phone_code' && credential.length < 6) {
        this.errors.credential = '密码长度不能少于6位'
      }
    },
    async handleLogin() {
      if (!this.isFormValid) return

      try {
        this.loading = true
        this.error = ''

        // 1. 调用登录接口
        const response = await login({
          login_type: this.form.loginType,
          account: this.form.account,
          credential: this.form.credential
        })

        // 2. 保存 token
        const token = response.data.token
        if (!token) {
          throw new Error('登录失败：未获取到 token')
        }

        // 3. 保存 token 到本地存储
        localStorage.setItem('token', token)

        // 4. 获取用户信息
        try {
          const userInfo = await getUserInfo()
          if (!userInfo.data) {
            throw new Error('获取用户信息失败')
          }
          // 5. 保存用户信息
          localStorage.setItem('user', JSON.stringify(userInfo.data))
        } catch (error) {
          console.error('获取用户信息失败:', error)
          // 如果获取用户信息失败，清除 token
          localStorage.removeItem('token')
          throw new Error('获取用户信息失败，请重新登录')
        }

        // 6. 处理记住我功能
        if (this.form.rememberMe) {
          localStorage.setItem('rememberedAccount', this.form.account)
        } else {
          localStorage.removeItem('rememberedAccount')
        }

        ElMessage.success('登录成功')
        this.$router.push('/')
      } catch (error) {
        this.error = error.response?.data?.message || error.message || '登录失败，请稍后重试'
        ElMessage.error(this.error)
      } finally {
        this.loading = false
      }
    },
    async handleSendCode() {
      if (this.countdown > 0 || !this.isAccountValid) return

      try {
        // TODO: 调用发送验证码接口
        this.countdown = 60
        const timer = setInterval(() => {
          this.countdown--
          if (this.countdown <= 0) {
            clearInterval(timer)
          }
        }, 1000)
        ElMessage.success('验证码已发送')
      } catch (error) {
        this.error = error.response?.data?.message || '发送验证码失败，请稍后重试'
        ElMessage.error(this.error)
      }
    }
  },
  created() {
    // 检查是否有记住的账号
    const rememberedAccount = localStorage.getItem('rememberedAccount')
    if (rememberedAccount) {
      this.form.account = rememberedAccount
      this.form.rememberMe = true
      this.validateAccount()
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
