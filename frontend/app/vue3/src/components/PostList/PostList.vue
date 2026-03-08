<script setup lang="ts">
import {useRouter} from 'vue-router'

import {$t} from '@/locales'
import PostCard from '@/components/PostCard';
import type {contentservicev1_Post} from "@/api/generated/app/service/v1";

interface Props {
  posts: contentservicev1_Post[]
  loading?: boolean
  showSkeleton?: boolean
  from?: string
  categoryId?: number
}

withDefaults(defineProps<Props>(), {
  loading: false,
  showSkeleton: true,
  from: 'post-list',
  categoryId: undefined
})
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
