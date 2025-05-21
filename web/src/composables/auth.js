import request from '@/pkg/request/request.js'

/**
 * 用户登录
 * @param {Object} params - 登录参数
 * @param {string} params.login_type - 登录类型：phone_password/email_password/phone_code
 * @param {string} params.account - 账号（手机号或邮箱）
 * @param {string} params.credential - 凭证（密码或验证码）
 * @returns {Promise} 返回登录结果，包含：
 *   - token: JWT令牌
 */
export function login(params) {
  return request.post('/login', params)
}

/**
 * 用户注册
 * @param {Object} params - 注册参数
 * @param {string} params.phone - 手机号
 * @param {string} params.email - 邮箱（可选）
 * @param {string} params.password - 密码
 * @returns {Promise} 返回注册结果，包含用户信息
 */
export function register(params) {
  return request.post('/register', params)
}

/**
 * 获取用户信息
 * @returns {Promise} 返回用户信息
 */
export function getUserInfo() {
  return request.get('/user')
}

/**
 * 退出登录
 * @returns {Promise} 返回登出结果
 */
export function logout() {
  return request.post('/user/logout').then(() => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  })
}

/**
 * 保存认证信息
 * @param {string} token - JWT令牌
 * @param {Object} user - 用户信息
 * @returns {void}
 */
export function saveAuthInfo(token, user) {
  localStorage.setItem('token', token)
  localStorage.setItem('user', JSON.stringify(user))
}

/**
 * 获取认证信息
 * @returns {Object} 返回认证信息
 *   - token: JWT令牌
 *   - user: 用户信息
 */
export function getAuthInfo() {
  const token = localStorage.getItem('token')
  const userStr = localStorage.getItem('user')
  let user = null
  if (userStr) {
    try {
      user = JSON.parse(userStr)
    } catch (error) {
      console.error('解析用户信息失败:', error)
      localStorage.removeItem('user')
    }
  }
  return {
    token,
    user
  }
}

export function sendVerificationCode(params) {
  return request.post('/send-code', params)
}

