import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import store from './store'

const selector = '#app';

const app = createApp(App)
const mountEl = document.querySelector(selector)
app.config.globalProperties.$appSettings = {...mountEl.dataset}

app.use(store)
app.mount(selector)