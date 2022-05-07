import axios from "axios"
import {reactive, ref} from "vue"
import {defineStore} from "pinia"
import {Notify} from "quasar"
import {Image} from "../data/image"
import {useI18n} from "vue-i18n"

export const useImageStore = defineStore("image", () => {
  const {t} = useI18n()
  const presignClient = axios.create() // use client without any default headers or interceptors

  // state
  const selected = ref<Image | null>(null)
  const images = reactive<Array<Image>>([])

  async function fetchImages() {
    const listResult = await axios.get("/list")
    if (listResult && listResult.data) images.push(...listResult.data.images)
  }

  function selectImage(img: Image) {
    selected.value = img
  }

  function deselectImage() {
    selected.value = null
  }

  async function uploadImage(file: File, caption: string): Promise<boolean> {
    try {
      const issueResult = await axios.post("/presign", {
        "contentType": file.type,
        "contentLength": file.size
      })
      const url = issueResult.data.url
      await presignClient.put(url, file, {headers: {"Content-Type": file.type}})
      const result = await axios.post("/persist", {url, caption})
      const img = result.data as Image
      images.unshift(img)
      Notify.create({
        type: 'positive',
        message: t("createDialog.uploadSuccess")
      })
      return true
    } catch (err) {
      Notify.create({
        type: 'negative',
        message: t("createDialog.uploadFailure")
      })
      return false
    }
  }

  return {
    images,
    selected,
    fetchImages,
    selectImage,
    deselectImage,
    uploadImage
  }
})
