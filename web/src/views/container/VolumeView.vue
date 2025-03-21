<script setup lang="ts">
import { NButton, NDataTable, NInput, NPopconfirm } from 'naive-ui'

import container from '@/api/panel/container'
import { formatDateTime } from '@/utils'

const createModel = ref({
  name: '',
  driver: 'local',
  options: [],
  labels: []
})

const options = [{ label: 'local', value: 'local' }]

const createModal = ref(false)

const selectedRowKeys = ref<any>([])

const columns: any = [
  { type: 'selection', fixed: 'left' },
  {
    title: '名称',
    key: 'name',
    minWidth: 150,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '驱动',
    key: 'driver',
    width: 100,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '范围',
    key: 'scope',
    width: 100,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '挂载点',
    key: 'mount_point',
    resizable: true,
    minWidth: 150,
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
  (page, pageSize) => container.volumeList(page, pageSize),
  {
    initialData: { total: 0, list: [] },
    initialPageSize: 20,
    total: (res: any) => res.total,
    data: (res: any) => res.items
  }
)

const handleDelete = async (row: any) => {
  useRequest(container.volumeRemove(row.id)).onSuccess(() => {
    refresh()
    window.$message.success('删除成功')
  })
}

const handlePrune = () => {
  useRequest(container.volumePrune()).onSuccess(() => {
    refresh()
    window.$message.success('清理成功')
  })
}

const handleCreate = () => {
  loading.value = true
  useRequest(container.volumeCreate(createModel.value))
    .onSuccess(() => {
      refresh()
      window.$message.success('创建成功')
    })
    .onComplete(() => {
      loading.value = false
      createModal.value = false
    })
}

onMounted(() => {
  refresh()
})
</script>

<template>
  <n-flex vertical :size="20">
    <n-flex>
      <n-button type="primary" @click="createModal = true">创建卷</n-button>
      <n-button type="primary" @click="handlePrune" ghost>清理卷</n-button>
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
    v-model:show="createModal"
    preset="card"
    title="创建卷"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-form :model="createModel">
      <n-form-item path="name" label="卷名">
        <n-input v-model:value="createModel.name" type="text" @keydown.enter.prevent />
      </n-form-item>
      <n-form-item path="driver" label="驱动">
        <n-select
          :options="options"
          v-model:value="createModel.driver"
          type="text"
          @keydown.enter.prevent
        >
        </n-select>
      </n-form-item>
      <n-form-item path="env" label="标签">
        <n-dynamic-input
          v-model:value="createModel.labels"
          preset="pair"
          key-placeholder="标签名"
          value-placeholder="标签值"
        />
      </n-form-item>
      <n-form-item path="env" label="选项">
        <n-dynamic-input
          v-model:value="createModel.options"
          preset="pair"
          key-placeholder="选项名"
          value-placeholder="选项值"
        />
      </n-form-item>
    </n-form>
    <n-button type="info" block :loading="loading" :disabled="loading" @click="handleCreate">
      提交
    </n-button>
  </n-modal>
</template>
