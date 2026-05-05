<script setup lang="ts">
import { ref, computed } from 'vue'
import { ContentEditor, type EditorType } from '@/components/ContentEditor'
import { $t } from '@/locales'

const content = ref('')
const editorType = ref<EditorType>('markdown')

// Type options for selector
const typeOptions = computed(() => [
  { label: $t('component.editor.markdown'), value: 'markdown' },
  { label: $t('component.editor.richtext'), value: 'richtext' },
  { label: $t('component.editor.code'), value: 'code' },
  { label: $t('component.editor.json'), value: 'json' },
  { label: $t('component.editor.text'), value: 'text' },
])

const exampleMarkdown = '# Welcome to Content Editor\n\n## Features\n\n- **Bold** and *italic* text\n- [Links](https://example.com)\n- `Inline code` and code blocks\n\n```javascript\nfunction hello() {\n  console.log(\'Hello, World!\');\n}\n```\n\n> This is a quote\n\n## Lists\n\n1. First item\n2. Second item\n3. Third item\n\n- Unordered item\n- Another item\n\n---\n\n![Image](/logo.png)\n'

function handleContentChange(value: string) {
  console.log('Content changed:', value)
}


function loadExample() {
  content.value = exampleMarkdown
}
</script>

<template>
  <div class="editor-demo">
    <n-card :title="`${$t('component.editor.markdown')} ${$t('component.editor.edit')}`" class="editor-card">
      <template #header-extra>
        <n-space>
          <n-select
            v-model:value="editorType"
            :options="typeOptions"
            style="width: 150px"
            size="small"
          />
          <n-button size="small" @click="loadExample">
            {{ $t('common.loadExample') }}
          </n-button>
        </n-space>
      </template>

      <div class="editor-container">
        <ContentEditor
          :key="editorType"
          v-model="content"
          :editor-type="editorType"
          @change="handleContentChange"
        />
      </div>

      <template #footer>
        <div class="editor-info">
          <n-text depth="3">
            {{ $t('common.currentType') }}: <n-tag size="small">{{ editorType }}</n-tag>
          </n-text>
          <n-text depth="3">
            {{ $t('common.contentLength') }}: {{ content.length }} {{ $t('common.characters') }}
          </n-text>
        </div>
      </template>
    </n-card>
  </div>
</template>

<style scoped lang="less">
.editor-demo {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
  height: calc(100vh - 100px);
  display: flex;
  flex-direction: column;

  .editor-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0; // Important for flex children with overflow

    :deep(.n-card__header) {
      flex-shrink: 0;
    }

    :deep(.n-card__content) {
      flex: 1;
      display: flex;
      flex-direction: column;
      padding: 1rem;
      overflow: hidden;
      min-height: 0;
    }

    .editor-container {
      flex: 1;
      display: flex;
      overflow: auto; // Allow scrolling if content is too large
      min-height: 0;

      :deep(.content-editor) {
        flex: 1;
        min-height: 300px;
        min-width: 100%;
      }
    }

    :deep(.n-card__footer) {
      flex-shrink: 0;
      border-top: 1px solid var(--color-border);
      padding: 1rem;
    }
  }

  .editor-info {
    display: flex;
    justify-content: space-between;
    gap: 2rem;
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .editor-demo {
    padding: 1rem;
    height: auto;

    .editor-card {
      height: 600px;
      min-height: 500px;

      .editor-info {
        gap: 1rem;
      }
    }
  }
}
</style>

