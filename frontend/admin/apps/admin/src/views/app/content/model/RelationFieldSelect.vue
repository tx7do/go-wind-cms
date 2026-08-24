<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';

import { i18n } from '@vben/locales';

import { Select } from 'ant-design-vue';

import {
  fetchListCategories,
  fetchListPages,
  fetchListPosts,
  fetchListTags,
} from '#/api';
import { PaginationQuery } from '#/transport/rest';

/**
 * relation 类型字段的实体引用选择器。
 *
 * 根据 relationConfig.targetEntityType 加载对应实体的列表，
 * 渲染为 Select 供管理员按显示名选择。值以 "entitytype:id" 格式
 * 存入 customFields（与后端 ValidateValues 的 parseEntityRef 对齐）。
 */
const props = defineProps<{
  targetEntityType?: string;
  modelValue?: string;
}>();
const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>();

interface Option {
  value: string;
  label: string;
}

const options = ref<Option[]>([]);
const loading = ref(false);
const selected = ref<string | undefined>(props.modelValue);

// targetEntityType → { fetch, displayField } 分派表
const dispatch: Record<string, { fetch: (p: PaginationQuery) => Promise<any>; displayField: string }> = {
  post: { fetch: fetchListPosts, displayField: 'title' },
  page: { fetch: fetchListPages, displayField: 'title' },
  category: { fetch: fetchListCategories, displayField: 'name' },
  tag: { fetch: fetchListTags, displayField: 'name' },
};

// 从实体的 translations[] 中按当前 locale 提取显示名
function pickDisplayName(entity: any, displayField: string): string {
  const currentLang = i18n.global.locale.value;
  const translation = entity.translations?.find(
    (t: any) => t.languageCode === currentLang,
  );
  return (
    translation?.[displayField] ||
    entity.translations?.[0]?.[displayField] ||
    ''
  );
}

async function loadOptions() {
  const entityType = props.targetEntityType?.toLowerCase();
  if (!entityType || !dispatch[entityType]) {
    options.value = [];
    return;
  }

  loading.value = true;
  try {
    const resp = await dispatch[entityType].fetch(
      new PaginationQuery({
        paging: { page: 1, pageSize: 200 },
      }),
    );
    const items = (resp?.items ?? []) as any[];
    options.value = items
      .filter((e) => e.id != null)
      .map((e) => ({
        value: `${entityType}:${e.id}`,
        label: `${pickDisplayName(e, dispatch[entityType].displayField)} (#${e.id})`,
      }));
  } catch {
    options.value = [];
  } finally {
    loading.value = false;
  }
}

onMounted(loadOptions);
watch(() => props.targetEntityType, loadOptions);

function handleChange(val: any) {
  selected.value = val as string | undefined;
  emit('update:modelValue', (val as string) ?? '');
}
</script>

<template>
  <Select
    :value="selected"
    :options="options"
    :loading="loading"
    :placeholder="$t('ui.placeholder.select')"
    allow-clear
    show-search
    :filter-option="(input: string, option: any) => option.label.toLowerCase().includes(input.toLowerCase())"
    @change="handleChange"
  />
</template>
