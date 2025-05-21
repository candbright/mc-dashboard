# MC Dashboard Frontend

基于 Vue 3 + Element Plus 开发的 Minecraft 服务器管理面板前端项目。

## 技术栈

- Vue 3 - 渐进式 JavaScript 框架
- Element Plus - 基于 Vue 3 的组件库
- Vite - 下一代前端构建工具
- Pinia - Vue 的状态管理库
- Vue Router - Vue.js 的官方路由
- Axios - 基于 Promise 的 HTTP 客户端

## 开发环境要求

- Node.js 16.0+
- npm 7.0+

## 推荐开发工具

- [VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (禁用 Vetur)
- [Vue DevTools](https://devtools.vuejs.org/) - Vue 开发者工具

## 项目结构

```
web/
├── src/                # 源代码
│   ├── api/           # API 接口
│   ├── assets/        # 静态资源
│   ├── components/    # 公共组件
│   ├── composables/   # 组合式函数
│   ├── router/        # 路由配置
│   ├── store/         # 状态管理
│   ├── styles/        # 全局样式
│   ├── utils/         # 工具函数
│   └── views/         # 页面视图
├── public/            # 公共资源
├── index.html         # HTML 模板
├── vite.config.js     # Vite 配置
└── package.json       # 项目依赖
```

## 快速开始

1. 安装依赖
```bash
npm install
```

2. 启动开发服务器
```bash
npm run dev
```

3. 构建生产版本
```bash
npm run build
```

4. 代码检查
```bash
npm run lint
```

## 开发规范

### 代码风格
- 使用 ESLint 进行代码检查
- 使用 Prettier 进行代码格式化
- 遵循 Vue 3 组合式 API 风格指南

### 组件开发
- 组件文件使用 PascalCase 命名
- 组件名使用 PascalCase 命名
- 组件属性使用 kebab-case 命名

### 状态管理
- 使用 Pinia 进行状态管理
- 按功能模块拆分 store
- 使用组合式函数封装业务逻辑

### API 调用
- 使用 Axios 进行 HTTP 请求
- API 接口按模块分类
- 统一处理请求错误

## 环境变量

在项目根目录创建 `.env` 文件：

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_TITLE=MC Dashboard
```

## 部署说明

1. 构建生产版本
```bash
npm run build
```

2. 将 `dist` 目录下的文件部署到 Web 服务器

3. 配置 Web 服务器（以 Nginx 为例）：
```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /path/to/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://backend-server;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 常见问题

1. 开发环境跨域问题
   - 在 `vite.config.js` 中配置代理
   - 或使用浏览器插件临时禁用跨域限制

2. 组件热更新不生效
   - 检查组件命名是否符合规范
   - 确保组件正确导出

3. 构建失败
   - 检查 Node.js 版本是否符合要求
   - 清除 node_modules 后重新安装依赖

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License
