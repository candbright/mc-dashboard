/**
 * 通用工具类 Utils
 * 采用 ES6 模块化规范
 */

// ==================== 类型判断模块 ====================
const TypeUtils = {
  isNumber: value => typeof value === 'number' && Number.isFinite(value),
  isString: value => typeof value === 'string',
  isObject: value => value !== null && typeof value === 'object',
  isPrimitive: value => Object(value) !== value
}

// ==================== 数据处理模块 ====================
const DataUtils = {
  /**
   * 深度转换对象中的数值类型为字符串（支持嵌套对象/数组）
   * @param {Object|Array} data - 输入数据
   * @param {Object} [options]
   * @param {boolean} [options.recursive=true] - 是否递归处理
   * @param {boolean} [options.skipNull=true] - 是否跳过 null 值
   */
  convertNumbersToString(data, options = { recursive: true, skipNull: true }) {
    const processor = value => {
      if (TypeUtils.isNumber(value)) return value.toString()
      if (options.recursive && TypeUtils.isObject(value)) {
        return Array.isArray(value)
          ? value.map(processor)
          : this.convertNumbersToString(value, options)
      }
      return value
    }

    return Array.isArray(data)
      ? data.map(processor)
      : Object.entries(data).reduce((acc, [key, value]) => {
        acc[key] = processor(value)
        return acc
      }, {})
  },

  /**
   * 安全获取深层对象属性
   * @param {Object} obj - 目标对象
   * @param {string} path - 属性路径 (e.g. 'a.b.c')
   * @param {any} defaultValue - 默认值
   */
  deepGet(obj, path, defaultValue = null) {
    return path.split('.').reduce((acc, key) =>
      (acc && acc[key] !== undefined) ? acc[key] : defaultValue, obj)
  }
}

export const Utils = {
  ...TypeUtils,
  ...DataUtils,
  /**
   * 扩展工具类
   * @param {Object} extensions - 扩展方法集合
   */
  extend(extensions) {
    Object.assign(this, extensions)
  }
}

