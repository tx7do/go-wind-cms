<template>
  <section class="popular-tags-container">
    <div class="tags-content">
      <div v-if="loading" class="tags-grid">
        <div v-for="i in skeletonCount" :key="i" class="tag-item">
          <n-skeleton width="60px" height="28px" :sharp="false"/>
        </div>
      </div>
      <div v-else class="tags-grid">
        <div
          v-for="tag in displayTags"
          :key="tag.id"
          class="tag-item"
          :style="{ '--tag-color': tag.color }"
          @click="handleViewTag(tag)"
        >
          <span class="tag-label">{{ tag.name }}</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import {ref, onMounted, defineExpose} from 'vue';
import {useTagStore} from '@/stores';
import type {contentservicev1_Tag} from '@/api/generated/app/service/v1';
import {$t} from '@/locales';
import {useLanguageChangeEffect} from '@/hooks/useLanguageChangeEffect';
import router from "@/router/router";

const props = defineProps<{
  title?: string;
  pageSize?: number;
  filter?: Record<string, any>;
  skeletonCount?: number;
}>();

const tagStore = useTagStore();
const tags = ref<contentservicev1_Tag[]>([]);
const loading = ref(false);

const displayTags = ref<{ id: number; name: string; color: string }[]>([]);

function handleViewTag(tag: any) {
  router.push(`/tag/${tag.id}`)
}

async function loadPopularTags() {
  loading.value = true;
  try {
    const res = await tagStore.listTag(
      {page: 1, pageSize: props.pageSize ?? 6},
      props.filter ?? {status: 'TAG_STATUS_ACTIVE', isFeatured: true}
    );
    tags.value = res.items || [];
    displayTags.value = tags.value.map((tag, index) => ({
      id: tag.id,
      name: tag.translations?.[0]?.name || $t('page.tags.tag_untitled'),
      color: tag.color || `hsl(${index * 60}, 100%, 50%)`,
    }));
  } catch (e) {
    tags.value = [];
    displayTags.value = [];
  }
  loading.value = false;
}

function reload() {
  loadPopularTags();
}

defineExpose({reload});

onMounted(() => {
  loadPopularTags();
});

useLanguageChangeEffect(() => {
  loadPopularTags();
}, {immediate: false, autoCleanup: true});
</script>

<style scoped lang="less">
.popular-tags-container {
  width: 100%;

  .tags-content {
    width: 100%;
  }

  .tags-title {
    font-size: 2.25rem;
    font-weight: 800;
    letter-spacing: -1px;
    margin: 0 0 2rem 0;
    line-height: 1.1;
    background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    text-shadow: 0 2px 8px rgba(99, 102, 241, 0.08), 0 1px 2px rgba(0, 0, 0, 0.08);
  }

  .tags-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 1.5rem;
    margin-bottom: 1rem;
  }

  .tag-item {
    display: flex;
    align-items: center;
    padding: 1rem 1.5rem;
    background: #fff;
    border: 1.5px solid #e0e7ef;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    color: var(--tag-color, #6366f1);
    border-radius: 16px;
    cursor: pointer;
    font-weight: 700;
    font-size: 1.05rem;
    letter-spacing: 0.1px;
    transition: all 0.3s;
    min-width: 90px;

    .tag-label {
      font-size: 1.05rem;
      color: var(--tag-color, #6366f1);
      font-weight: 700;
      letter-spacing: 0.1px;
      transition: color 0.3s;
      word-break: break-all;
    }

    &:hover {
      transform: translateY(-4px);
      background: var(--tag-color, #6366f1);
      border-color: var(--tag-color, #6366f1);
      box-shadow: 0 8px 20px rgba(99, 102, 241, 0.3);

      .tag-label {
        color: #fff;
      }
    }
  }
}

html.dark .popular-tags-container {
  .tags-title {
    font-size: 2.25rem;
    font-weight: 800;
    letter-spacing: -1px;
    margin: 0 0 2rem 0;
    line-height: 1.1;
    background: linear-gradient(135deg, #a5b4fc 0%, #6366f1 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    text-shadow: 0 2px 8px #23263a44;
    color: #e0e7ef;
  }

  .tag-item {
    background: #1e293b;
    border: 1.5px solid rgba(51, 65, 85, 0.6);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    color: var(--tag-color, #a5b4fc);
    border-radius: 16px;

    .tag-label {
      color: var(--tag-color, #a5b4fc);
      font-weight: 700;
    }

    &:hover {
      transform: translateY(-4px);
      background: #1e293b;
      border-color: var(--tag-color, #6366f1);
      box-shadow: 0 8px 20px rgba(99, 102, 241, 0.25);

      .tag-label {
        color: #fff;
      }
    }
  }
}

@media (max-width: 900px) {
  .popular-tags-container {
    padding: 0 0.5rem;

    .tags-grid {
      gap: 0.5rem;
    }

    .tag-item {
      padding: 0.6rem 1rem;
      font-size: 0.92rem;
      border-radius: 12px;
      min-width: 60px;
    }

    .tag-label {
      font-size: 0.92rem;
    }
  }
}
</style>
