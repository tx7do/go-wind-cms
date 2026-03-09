<script setup lang="ts">
import {ref, watch, defineProps, defineExpose} from 'vue';
import {usePostStore} from '@/stores';
import {$t} from '@/locales';
import PostCard from '@/components/PostCard';
import type {contentservicev1_Post} from "@/api/generated/app/service/v1";
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';

interface Props {
  queryParams?: object;
  fieldMask?: string;
  orderBy?: string[];
  page?: number;
  pageSize?: number;
  showSkeleton?: boolean;
  from?: string;
  categoryId?: number;
}

const props = withDefaults(defineProps<Props>(), {
  queryParams: () => ({}),
  fieldMask: undefined,
  orderBy: undefined,
  page: 1,
  pageSize: 6,
  showSkeleton: true,
  from: 'post-list',
  categoryId: undefined
});

const postStore = usePostStore();
const posts = ref<contentservicev1_Post[]>([]);
const loading = ref(false);

async function fetchPosts() {
  loading.value = true;
  try {
    const res = await postStore.listPost(
      {page: props.page, pageSize: props.pageSize},
      props.queryParams,
      props.fieldMask,
      props.orderBy
    );
    posts.value = res.items || [];
  } catch (error) {
    posts.value = [];
    console.error('PostList fetch failed:', error);
  } finally {
    loading.value = false;
  }
}

function reload() {
  fetchPosts();
}

defineExpose({ reload });

watch(
  () => [props.queryParams, props.fieldMask, props.orderBy, props.page, props.pageSize],
  () => {
    fetchPosts();
  },
  { deep: true, immediate: true }
);

useLanguageChangeEffect(() => {
  fetchPosts();
});
</script>

<template>
  <div class="posts-list-container">
    <!-- Loading Skeleton -->
    <div v-if="loading && showSkeleton" class="posts-grid">
      <div v-for="i in 6" :key="i" class="post-card">
        <n-skeleton height="220px"/>
        <div style="padding: 24px;">
          <n-skeleton :width="'80%'" size="medium" style="margin-bottom: 12px;"/>
          <n-skeleton :width="'90%'" size="small" style="margin-bottom: 16px;"/>
          <div style="display: flex; gap: 12px;">
            <n-skeleton :width="60" size="small"/>
            <n-skeleton :width="60" size="small"/>
            <n-skeleton :width="60" size="small"/>
          </div>
        </div>
      </div>
    </div>

    <!-- Loaded Content -->
    <div v-else-if="!loading && posts.length > 0" class="posts-grid">
      <PostCard
        v-for="post in posts"
        :key="post.id"
        :post="post"
        :from="from"
        :category-id="categoryId"
      />
    </div>

    <!-- Empty State -->
    <div v-else-if="!loading && posts.length === 0" class="empty-state">
      <n-empty :description="$t('page.posts.no_results')"/>
    </div>
  </div>
</template>

<style scoped lang="less">
.posts-list-container {
  width: 100%;
}

.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.empty-state {
  padding: 60px 20px;
  text-align: center;
}

// 响应式
@media (max-width: 768px) {
  .posts-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>
