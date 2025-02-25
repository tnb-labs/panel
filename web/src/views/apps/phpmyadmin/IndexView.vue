<script setup lang="ts">
defineOptions({
  name: 'apps-phpmyadmin-index'
})

import Editor from '@guolao/vue-monaco-editor'
import { NButton } from 'naive-ui'

import phpmyadmin from '@/api/apps/phpmyadmin'

const currentTab = ref('status')
const hostname = ref(window.location.hostname)
const port = ref(0)
const path = ref('')
const newPort = ref(0)
const url = computed(() => {
  return `http://${hostname.value}:${port.value}/${path.value}`
})

const { data: config } = useRequest(phpmyadmin.getConfig, {
  initialData: {
    config: ''
  }
})

const getInfo = async () => {
  const data = await phpmyadmin.info()
  path.value = data.path
  port.value = data.port
  newPort.value = data.port
}

const handleSave = () => {
  useRequest(phpmyadmin.port(newPort.value)).onSuccess(() => {
    window.$message.success('保存成功')
    getInfo()
  })
}

const handleSaveConfig = () => {
  useRequest(phpmyadmin.saveConfig(config.value)).onSuccess(() => {
    window.$message.success('保存成功')
  })
}

onMounted(() => {
  getInfo()
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button v-if="currentTab == 'status'" class="ml-16" type="primary" @click="handleSave">
        <TheIcon :size="18" icon="material-symbols:save-outline" />
        保存
      </n-button>
      <n-button
        v-if="currentTab == 'config'"
        class="ml-16"
        type="primary"
        @click="handleSaveConfig"
      >
        <TheIcon :size="18" icon="material-symbols:save-outline" />
        保存
      </n-button>
    </template>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="status" tab="状态">
        <n-space vertical>
          <n-card title="访问信息">
            <n-alert type="info">
              访问地址: <a :href="url" target="_blank">{{ url }}</a>
            </n-alert>
          </n-card>
          <n-card title="修改端口">
            <n-input-number v-model:value="newPort" :min="1" :max="65535" />
            修改 phpMyAdmin 访问端口
          </n-card>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="config" tab="修改配置">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是 phpMyAdmin 的 OpenResty
            配置文件，如果您不了解各参数的含义，请不要随意修改！
          </n-alert>
          <Editor
            v-model:value="config"
            language="ini"
            theme="vs-dark"
            height="60vh"
            mt-8
            :options="{
              automaticLayout: true,
              formatOnType: true,
              formatOnPaste: true
            }"
          />
        </n-space>
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>
