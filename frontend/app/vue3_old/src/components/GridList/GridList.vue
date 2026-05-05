<script setup lang="ts">
import {ref, onMounted, reactive} from 'vue';

import {LoadingState} from "@/utils";

const props = withDefaults(
  defineProps<{
    pageSize?: number
    cols: number
    xGap?: number
    yGap?: number
    api: Function
    params?: any
    component: any
  }>(),
  {
    pageSize: 16,
    cols: 1,
    xGap: 16,
    yGap: 16,
  }
)

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: props.pageSize
});

const loading = new LoadingState();

const dataRef = ref<any[]>([]);

/**
 * 处理翻页
 * @param newPage
 */
function handlePageChange(newPage: number) {
  pagination.page = newPage;
  reload();
}

/**
 * 重新加载数据
 */
function reload() {
  loading.setLoadingState(true);

  props.api(false, pagination.page, pagination.pageSize, props.params)
    .then((res: any) => {
      if (res === undefined) {
        loading.setLoadingState(false);
        return
      }

      dataRef.value = res.items || [];
      pagination.pageCount = Math.ceil(res.total! / pagination.pageSize) || 0;

      loading.setLoadingState(false);
    }).catch(() => {
    loading.setLoadingState(false, true);
  });
}

onMounted(async () => {
  reload();
});

// 暴露方法给父组件
defineExpose({
  reload
});
</script>

<template>
  <n-layout>
    <n-layout-content>
      <n-grid :cols="props.cols" :x-gap="props.xGap" :y-gap="props.yGap">
        <n-grid-item v-for="(item, index) in dataRef" :key="index">
          <component :is="props.component" :data="item"/>
        </n-grid-item>
      </n-grid>
    </n-layout-content>
    <n-layout-footer class="footer">
      <n-pagination
        v-model:page="pagination.page"
        :page-size="pagination.pageSize"
        :page-count="pagination.pageCount"
        @update:page="handlePageChange"
      />
    </n-layout-footer>
  </n-layout>
</template>

<style scoped lang="less">
@import "@/styles/app.less";

.footer {
  display: flex;
  justify-content: center;
  padding: 16px 0;
  background-color: @content-color;
}
</style>
