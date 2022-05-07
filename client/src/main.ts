import {createApp} from 'vue'
import {Quasar, Notify} from 'quasar'
import i18n from "./plugins/vue-i18n"
import {createPinia} from 'pinia'
import App from './App.vue'

import "./plugins/axios" // init axios
import '@quasar/extras/roboto-font/roboto-font.css'
import '@quasar/extras/material-icons/material-icons.css'
import 'quasar/dist/quasar.css'

createApp(App)
    .use(i18n)
    .use(createPinia())
    .use(Quasar, {
      plugins: {
        Notify
      },
      config: {
        brand: {
          primary: '#a1bf14',
          secondary: '#3c3f41',
          accent: '#ffffff',

          dark: '#1d1d1d',

          positive: '#21BA45',
          negative: '#C10015',
          info: '#31CCEC',
          warning: '#F2C037'
        }
      }
    })
    .mount('#app')
