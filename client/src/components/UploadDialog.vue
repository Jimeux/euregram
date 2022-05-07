<script setup lang="ts">
import {ref, watch} from "vue"
import {useI18n} from "vue-i18n"
import {useImageStore} from "../store/images"

const {t} = useI18n()
const {uploadImage} = useImageStore()

const prompt = ref(false)
const selectedFile = ref<File | null>()
const selectedSrc = ref<string | null>()
const caption = ref<string>("")
const loading = ref<boolean>(false)

watch(selectedFile, async (value) =>
    selectedSrc.value = value == null ? null : URL.createObjectURL(value))

function clear() {
  selectedFile.value = null
  caption.value = ""
}

function canSubmit(): boolean {
  return selectedFile.value != null && caption.value !== ""
}

async function onSubmit() {
  if (!canSubmit()) return
  loading.value = true
  const success = await uploadImage(selectedFile.value!, caption.value)
  loading.value = false
  if (success) {
    prompt.value = false
    clear()
  }
}
</script>

<template>
  <div>
    <q-dialog v-model="prompt">
      <q-card style="min-width: 370px;" class="bg-secondary text-accent">
        <q-card-section class="q-mb-sm">
          <div class="text-h6">{{ t("createDialog.title") }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-file
              name="poster_file"
              dark
              borderless
              v-model="selectedFile"
              filled
              :label="t(`createDialog.filePlaceholder`)"
          />
        </q-card-section>

        <q-img
            v-if="selectedSrc != null"
            class="q-mb-md rounded"
            width="200"
            :src="selectedSrc"
        />

        <q-card-section class="q-pt-none">
          <q-input dark borderless filled v-model="caption" :label="t(`createDialog.captionPlaceholder`)"/>
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat :label="t(`createDialog.close`)" v-close-popup @click="clear"/>
          <q-btn flat :label="t(`createDialog.share`)" @click="onSubmit" :disable="!canSubmit()" :loading="loading"/>
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-page-sticky position="bottom-right" :offset="[40, 40]">
      <q-btn fab icon="add" color="primary" @click="prompt = true"/>
    </q-page-sticky>
  </div>
</template>
