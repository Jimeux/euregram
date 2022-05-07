import {createI18n} from "vue-i18n"
import {en} from "../i18n/en"
import {fr} from "../i18n/fr"
import {ja} from "../i18n/ja"

export default createI18n({
  locale: 'ja',
  fallbackLocale: 'en',
  messages: {
    en,
    fr,
    ja,
  },
})
