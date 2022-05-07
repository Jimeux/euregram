<script setup lang="ts">
import {useI18n} from "vue-i18n"
import LanguageSelector from "./components/LanguageSelector.vue"
import Home from "./pages/Home.vue"

const {t} = useI18n()
</script>

<template>
  <q-layout view="hHh lpR fFf" class="bg-dark">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-toolbar-title class="text-secondary">
          {{ t("app.title") }}
        </q-toolbar-title>
        <LanguageSelector class="mr-5"/>
      </q-toolbar>
    </q-header>
    <q-page-container>
      <Home/>
    </q-page-container>
  </q-layout>
</template>
