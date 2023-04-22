import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import './assets/tailwind.css'

import { OhVueIcon, addIcons } from 'oh-vue-icons'
import { FaSort, FaSortUp, FaSortDown } from 'oh-vue-icons/icons'

addIcons(FaSort, FaSortUp, FaSortDown)

const app = createApp(App)

app.component('v-icon', OhVueIcon)

app.use(router)

app.mount('#app')
