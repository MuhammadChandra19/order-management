import { createApp } from 'vue'
import App from './App.vue'
import '@/assets/index.css'
import { router } from './lib/routes'

createApp(App).use(router).mount('#app')