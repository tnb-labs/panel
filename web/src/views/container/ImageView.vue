<script setup lang="ts">
import { NButton, NDataTable, NFlex, NInput, NPopconfirm, NTag } from 'naive-ui'

import container from '@/api/panel/container'
import { formatDateTime } from '@/utils'

const pullModel = ref({
  name: '',
  auth: false,
  username: '',
  password: ''
})
const pullModal = ref(false)
const selectedRowKeys = ref<any>([])

const columns: any = [
  { type: 'selection', fixed: 'left' },
  {
    title: 'ID',
    key: 'id',
    minWidth: 400,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '容器数',
    key: 'containers',
    width: 100,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '镜像',
    key: 'repo_tags',
    minWidth: 200,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any): any {
      return h(NFlex, null, {
        default: () =>
          row.repo_tags.map((tag: any) =>
            h(NTag, null, {
              default: () => tag
            })
          )
      })
    }
  },
  {
    title: '大小',
    key: 'size',
    width: 150,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 200,
    resizable: true,
    render(row: any) {
      return formatDateTime(row.created_at)
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    align: 'center',
    hideInExcel: true,
    render(row: any) {
      return [
        h(
          NPopconfirm,
          {
            onPositiveClick: async () => {
              await handleDelete(row)
            }
          },
          {
            default: () => {
              return '确定删除吗？'
            },
            trigger: () => {
              return h(
                NButton,
                {
                  size: 'small',
                  type: 'error'
                },
                {
                  default: () => '删除'
                }
              )
            }
          }
        )
      ]
    }
  }
]

const { loading, data, page, total, pageSize, pageCount, refresh } = usePagination(
  (page, pageSize) => container.imageList(page, pageSize),
  {
    initialData: { total: 0, list: [] },
    initialPageSize: 20,
    total: (res: any) => res.total,
    data: (res: any) => res.items
  }
)

const handleDelete = async (row: any) => {
  useRequest(container.imageRemove(row.id)).onSuccess(() => {
    refresh()
    window.$message.success('删除成功')
  })
}

const handlePrune = () => {
  useRequest(container.imagePrune()).onSuccess(() => {
    refresh()
    window.$message.success('清理成功')
  })
}

const handlePull = () => {
  loading.value = true
  useRequest(container.imagePull(pullModel.value))
    .onSuccess(() => {
      refresh()
      window.$message.success('拉取成功')
    })
    .onComplete(() => {
      loading.value = false
      pullModal.value = false
    })
}

onMounted(() => {
  refresh()
})
</script>

<template>
  <n-flex vertical :size="20">
    <n-flex>
      <n-button type="primary" @click="pullModal = true">拉取镜像</n-button>
      <n-button type="primary" @click="handlePrune" ghost>清理镜像</n-button>
    </n-flex>
    <n-data-table
      striped
      remote
      :loading="loading"
      :scroll-x="1000"
      :data="data"
      :columns="columns"
      :row-key="(row: any) => row.id"
      v-model:checked-row-keys="selectedRowKeys"
      v-model:page="page"
      v-model:pageSize="pageSize"
      :pagination="{
        page: page,
        pageCount: pageCount,
        pageSize: pageSize,
        itemCount: total,
        showQuickJumper: true,
        showSizePicker: true,
        pageSizes: [20, 50, 100, 200]
      }"
    />
  </n-flex>
  <n-modal
    v-model:show="pullModal"
    preset="card"
    title="拉取镜像"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-form :model="pullModel">
      <n-form-item path="name" label="镜像名">
        <n-input
          v-model:value="pullModel.name"
          type="text"
          @keydown.enter.prevent
          placeholder="docker.io/php:8.3-fpm"
        />
      </n-form-item>
      <n-form-item path="auth" label="验证">
        <n-switch v-model:value="pullModel.auth" />
      </n-form-item>
      <n-form-item v-if="pullModel.auth" path="username" label="用户名">
        <n-input
          v-model:value="pullModel.username"
          type="text"
          @keydown.enter.prevent
          placeholder="输入用户名"
        />
      </n-form-item>
      <n-form-item v-if="pullModel.auth" path="password" label="密码">
        <n-input
          v-model:value="pullModel.password"
          type="password"
          show-password-on="click"
          @keydown.enter.prevent
          placeholder="输入密码"
        />
      </n-form-item>
    </n-form>
    <n-button type="info" block :loading="loading" :disabled="loading" @click="handlePull">
      提交
    </n-button>
  </n-modal>
</template>
