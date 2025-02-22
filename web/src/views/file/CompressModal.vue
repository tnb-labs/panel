<script setup lang="ts">
import { NButton, NInput } from 'naive-ui'

import api from '@/api/panel/file'
import { generateRandomString, getBase } from '@/utils'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const path = defineModel<string>('path', { type: String, required: true })
const selected = defineModel<string[]>('selected', { type: Array, default: () => [] })
const format = ref('.zip')
const loading = ref(false)

const generateName = () => {
  return selected.value.length > 0
    ? `${getBase(selected.value[0])}-${generateRandomString(6)}${format.value}`
    : `${path.value}/${generateRandomString(8)}${format.value}`
}

const file = ref(generateName())

const ensureExtension = (extension: string) => {
  if (!file.value.endsWith(extension)) {
    file.value = `${getBase(file.value)}${extension}`
  }
}

const handleArchive = () => {
  ensureExtension(format.value)
  loading.value = true
  const message = window.$message.loading('正在压缩中...', {
    duration: 0
  })
  const paths = selected.value.map((item) => item.replace(path.value, '').replace(/^\//, ''))
  useRequest(api.compress(path.value, paths, file.value))
    .onSuccess(() => {
      show.value = false
      selected.value = []
      window.$message.success('压缩成功')
    })
    .onComplete(() => {
      message?.destroy()
      loading.value = false
      window.$bus.emit('file:refresh')
    })
}

onMounted(() => {
  watch(
    selected,
    () => {
      file.value = generateName()
    },
    { immediate: true }
  )
})
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="压缩"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-flex vertical>
      <n-form>
        <n-form-item label="待压缩">
          <n-dynamic-input v-model:value="selected" :min="1" />
        </n-form-item>
        <n-form-item label="压缩为">
          <n-input v-model:value="file" />
        </n-form-item>
        <n-form-item label="格式">
          <n-select
            v-model:value="format"
            :options="[
              { label: '.zip', value: '.zip' },
              { label: '.bz2', value: '.bz2' },
              { label: '.tar', value: '.tar' },
              { label: '.gz', value: '.gz' },
              { label: '.tar.gz', value: '.tar.gz' },
              { label: '.tgz', value: '.tgz' },
              { label: '.xz', value: '.xz' },
              { label: '.7z', value: '.7z' }
            ]"
            @update:value="ensureExtension"
          />
        </n-form-item>
      </n-form>
      <n-button :loading="loading" :disabled="loading" type="primary" @click="handleArchive">
        压缩
      </n-button>
    </n-flex>
  </n-modal>
</template>

<style scoped lang="scss"></style>
