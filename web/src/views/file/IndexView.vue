<script setup lang="ts">
defineOptions({
  name: 'file-index'
})

import { useFileStore } from '@/store'
import CompressModal from '@/views/file/CompressModal.vue'
import ListTable from '@/views/file/ListTable.vue'
import PathInput from '@/views/file/PathInput.vue'
import PermissionModal from '@/views/file/PermissionModal.vue'
import ToolBar from '@/views/file/ToolBar.vue'
import type { Marked } from '@/views/file/types'

const fileStore = useFileStore()

const selected = ref<string[]>([])
const marked = ref<Marked[]>([])
const markedType = ref<string>('copy')

const compress = ref(false)
const permission = ref(false)
</script>

<template>
  <common-page show-footer>
    <n-flex vertical :size="20">
      <path-input v-model:path="fileStore.path" />
      <tool-bar
        v-model:path="fileStore.path"
        v-model:selected="selected"
        v-model:marked="marked"
        v-model:markedType="markedType"
        v-model:compress="compress"
        v-model:permission="permission"
      />
      <list-table
        v-model:path="fileStore.path"
        v-model:selected="selected"
        v-model:marked="marked"
        v-model:markedType="markedType"
        v-model:compress="compress"
        v-model:permission="permission"
      />
      <compress-modal
        v-model:show="compress"
        v-model:path="fileStore.path"
        v-model:selected="selected"
      />
      <permission-modal v-model:show="permission" v-model:selected="selected" />
    </n-flex>
  </common-page>
</template>
