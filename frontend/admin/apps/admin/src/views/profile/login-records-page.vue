<script lang="ts" setup>
import { onMounted, ref } from 'vue';

import { useUserStore } from '@vben/stores';

import { Table } from 'ant-design-vue';

import {
  fetchListLoginAuditLogs,
  getLoginAuditLogActionTypeColor,
  getLoginAuditLogStatusColor,
  type auditservicev1_LoginAuditLog as LoginAuditLog,
  loginAuditLogActionTypeToName,
  loginAuditLogStatusToName,
  PaginationQuery,
} from '#/api';
import { $t } from '#/locales';

/**
 * 当前用户的登录记录（个人中心 tab）。
 *
 * 复用 LoginAuditLogService.List，按当前用户名过滤；
 * username 由后端 List 通用 query 谓词支持（与日志管理页同一链路）。
 */
const userStore = useUserStore();

const loading = ref(false);
const records = ref<LoginAuditLog[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(10);

const columns = [
  {
    title: $t('page.loginAuditLog.createdAt'),
    dataIndex: 'createdAt',
    key: 'createdAt',
    width: 170,
  },
  {
    title: $t('page.loginAuditLog.actionType'),
    dataIndex: 'actionType',
    key: 'actionType',
    width: 110,
  },
  {
    title: $t('page.loginAuditLog.status'),
    dataIndex: 'status',
    key: 'status',
    width: 100,
  },
  {
    title: $t('page.loginAuditLog.ipAddress'),
    dataIndex: 'ipAddress',
    key: 'ipAddress',
    width: 140,
  },
  {
    title: $t('page.loginAuditLog.geoLocation'),
    key: 'geoLocation',
  },
  {
    title: $t('page.loginAuditLog.platform'),
    key: 'platform',
  },
];

async function loadRecords() {
  loading.value = true;
  try {
    const resp = await fetchListLoginAuditLogs(
      new PaginationQuery({
        paging: { page: page.value, pageSize: pageSize.value },
        formValues: { username: userStore.userInfo?.username },
        orderBy: ['-created_at'],
      }),
    );
    records.value = (resp?.items ?? []) as LoginAuditLog[];
    total.value = Number(resp?.total ?? 0);
  } catch {
    records.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

function handlePageChange(nextPage: number, nextSize: number) {
  page.value = nextPage;
  pageSize.value = nextSize;
  loadRecords();
}

onMounted(loadRecords);
</script>

<template>
  <div class="p-4">
    <Table
      :columns="columns"
      :data-source="records"
      :loading="loading"
      :pagination="{
        current: page,
        pageSize: pageSize,
        total: total,
        showSizeChanger: true,
        showTotal: (t: number) => `${t}`,
        onChange: handlePageChange,
      }"
      row-key="id"
      size="small"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actionType'">
          <a-tag :color="getLoginAuditLogActionTypeColor(record.actionType)">
            {{ loginAuditLogActionTypeToName(record.actionType) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'status'">
          <a-tag :color="getLoginAuditLogStatusColor(record.status)">
            {{ loginAuditLogStatusToName(record.status) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'geoLocation'">
          <span>
            {{ record.geoLocation?.province }}{{ record.geoLocation?.city }}
          </span>
        </template>
        <template v-else-if="column.key === 'platform'">
          <span>
            {{ record.deviceInfo?.osName }} {{ record.deviceInfo?.browserName }}
          </span>
        </template>
      </template>
    </Table>
  </div>
</template>
