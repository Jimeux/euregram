<script setup lang="ts">
import {ref, watch} from "vue"
import {useI18n} from "vue-i18n"

const {locale} = useI18n()

const languages = [
  {value: "ja", label: "🍣"},
  {value: "en", label: "🍔"},
  {value: "fr", label: "🥖"},
]

// todo extract to a store
let defaultLang = languages[0]
const str = localStorage.getItem("auth")
if (str) {
  const auth = JSON.parse(str)
  if (auth.locale) {
    const [lang, country] = auth.locale.split("-")
    const found = languages.find(l => l.value === lang)
    if (found) {
      defaultLang = found
      locale.value = defaultLang.value
    }
  }
}

const model = ref<Record<string, string>>(defaultLang)
watch(model, (value) => locale.value = value.value)
</script>

<template>
  <q-select class="locale-changer q-mr-md"
            dark
            borderless
            popup-content-class="text-h5"
            hide-dropdown-icon
            v-model="model"
            :options="languages"/>
</template>

<style scoped>
.locale-changer {
  font-size: 24px;
}
</style>
