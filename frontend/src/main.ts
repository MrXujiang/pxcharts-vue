import 'normalize.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import '@/assets/style/global.less'
import 'tdesign-vue-next/es/style/index.css'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
