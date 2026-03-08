<script setup lang="ts">
import {useRouter} from 'vue-router'

import {usePostStore} from '@/stores';
import type {contentservicev1_Post} from "@/api/generated/app/service/v1";
import {scrollToTop} from "@/utils";
import {formatDate} from "@/utils/date";

interface Props {
  post: contentservicev1_Post
  from?: string
  categoryId?: number
}

const props = withDefaults(defineProps<Props>(), {
  from: 'post-list',
  categoryId: undefined
})

const router = useRouter()
const postStore = usePostStore()

function handleViewPost() {
  const query: any = {from: props.from}
  if (props.categoryId) {
    query.categoryId = props.categoryId
  }

  router.push({
    path: `/post/${props.post.id}`,
    query
  })

  scrollToTop()
}
</script>

<template>
  <article
    class="post-card"
    @click="handleViewPost"
  >
    <div class="post-image">
      <img :src="postStore.getPostThumbnail(post)" :alt="postStore.getPostTitle(post)"/>
      <div class="image-overlay"/>
    </div>
    <div class="post-content">
      <h3 class="post-title">{{ postStore.getPostTitle(post) }}</h3>
      <p class="post-summary">{{ postStore.getPostSummary(post) }}</p>
      <div class="post-meta">
        <div class="meta-item">
          <span class="i-carbon:user"/>
          <span>{{ post.authorName }}</span>
        </div>
        <div class="meta-item">
          <span class="i-carbon:calendar"/>
          <span>{{ formatDate(post.createdAt) }}</span>
        </div>
        <div class="meta-item">
          <span class="i-carbon:view"/>
          <span>{{ post.visits || 0 }}</span>
        </div>
        <div class="meta-item">
          <span class="i-carbon:thumbs-up"/>
          <span>{{ post.likes || 0 }}</span>
        </div>
      </div>
    </div>
  </article>
</template>

<style scoped lang="less">
.post-card {
  background: var(--color-surface);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--color-border);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  height: 100%;

  &:hover {
    transform: translateY(-6px);
    box-shadow: 0 16px 32px rgba(0, 0, 0, 0.12);
    border-color: var(--color-brand);

    .post-image {
      .image-overlay {
        opacity: 0.5;
      }

      img {
        transform: scale(1.08);
      }
    }

    .post-title {
      color: var(--color-brand);
    }
  }

  .post-image {
    position: relative;
    width: 100%;
    height: 240px;
    flex-shrink: 0;
    overflow: hidden;
    background: var(--color-bg);

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
    }

    .image-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.15);
      transition: opacity 0.3s;
      opacity: 0;
    }
  }

  .post-content {
    padding: 24px;
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 12px;

    .post-title {
      font-size: 18px;
      font-weight: 700;
      margin: 0;
      color: var(--color-text-primary);
      line-height: 1.4;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      transition: color 0.3s;
    }

    .post-summary {
      color: var(--color-text-secondary);
      font-size: 14px;
      line-height: 1.6;
      margin: 0;
      display: -webkit-box;
      -webkit-line-clamp: 3;
      -webkit-box-orient: vertical;
      overflow: hidden;
      flex: 1;
    }

    .post-meta {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 12px;
      font-size: 13px;
      color: var(--color-text-secondary);
      padding-top: 12px;
      border-top: 1px solid var(--color-border);

      .meta-item {
        display: flex;
        align-items: center;
        gap: 6px;
        white-space: nowrap;

        span[class^="i-"] {
          font-size: 16px;
          opacity: 0.7;
        }
      }
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .post-card {
    .post-image {
      height: 200px;
    }

    .post-content {
      padding: 16px;
      gap: 10px;

      .post-title {
        font-size: 17px;
      }

      .post-summary {
        font-size: 13px;
      }

      .post-meta {
        font-size: 12px;
        flex-wrap: wrap;

        .meta-item {
          span[class^="i-"] {
            font-size: 14px;
          }
        }
      }
    }
  }
}
</style>
