import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from '@/router/index.js'
import '@fortawesome/fontawesome-free/css/all.min.css'
import 'tailwindcss/tailwind.css'
import './styles/element-plus.scss'
const app = createApp(App)

app.use(ElementPlus)
app.use(router)
app.mount('#app')
