<script setup lang="ts">
import {computed} from "vue"
import {useI18n} from "vue-i18n"
import {storeToRefs} from "pinia"
import {useImageStore} from "../store/images"

const {t} = useI18n()

const store = useImageStore()
const {selected} = storeToRefs(store)
const {deselectImage} = store

const show = computed({
  get: () => selected != null,
  set: deselectImage
})
</script>

<template>
  <q-dialog v-model="show" v-if="selected">
    <q-card style="min-width: 600px;" class="bg-secondary text-accent">
      <q-card-section>
        <div class="text-h6">{{ selected.caption }}</div>
        <div class="text-subtitle2">by &nbsp;<span class="text-primary">{{ selected.username }}</span></div>
      </q-card-section>

      <q-img class="mx-4 rounded" width="200" :src="selected.url"/>

      <q-card-actions align="right" class="text-primary">
        <q-btn flat :label="t(`detailDialog.close`)" v-close-popup/>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
