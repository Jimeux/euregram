<script setup lang="ts">
import {onMounted} from "vue"
import {useImageStore} from "../store/images"

const {images, fetchImages, selectImage} = useImageStore()
onMounted(fetchImages)
</script>

<template>
  <div style="max-width: 1300px; margin: 0 auto">
    <div class="row items-start q-pa-lg q-col-gutter-lg">
      <q-card class="col-md-4 col-sm-6 col-xs-12 bg-dark item" v-for="img in images" @click="selectImage(img)">
        <q-img :src="img.url" :ratio="1">
          <div class="absolute-bottom text-subtitle2 text-center">
            {{ img.caption }}
          </div>
        </q-img>
      </q-card>
    </div>
  </div>
</template>

<style scoped>
.item {
  cursor: pointer;
}
</style>
