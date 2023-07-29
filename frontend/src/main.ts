import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(router)

app.config.globalProperties.imgbbSecret = import.meta.env.VITE_IMGBB_SECRET
export const useGlobals = () => app.config.globalProperties

app.mount('#app')
