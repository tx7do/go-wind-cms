/**
 * ContentEditor Component - Usage Examples
 *
 * This file demonstrates various ways to use the ContentEditor component
 */

// ============================================================
// Example 1: Basic Markdown Editor
// ============================================================
/*
<template>
  <ContentEditor v-model="content" type="markdown" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ContentEditor } from '@/components/ContentEditor'

const content = ref('# Hello World\n\nThis is a **Markdown** editor.')
</script>
*/

// ============================================================
// Example 2: Rich Text Editor with Custom Height
// ============================================================
/*
<template>
  <ContentEditor
    v-model="richContent"
    type="richtext"
    height="600px"
    placeholder="Start writing your content..."
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ContentEditor } from '@/components/ContentEditor'

const richContent = ref('')
</script>
*/

// ============================================================
// Example 3: Editor with Type Switching and Event Handlers
// ============================================================
/*
<template>
  <div>
    <ContentEditor
      v-model="content"
      :type="editorType"
      :show-toolbar="true"
      :show-type-switch="true"
      @change="handleContentChange"
      @type-change="handleTypeChange"
    />

    <n-card style="margin-top: 16px">
      <p>Current Type: {{ editorType }}</p>
      <p>Content Length: {{ content.length }}</p>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ContentEditor, type EditorType } from '@/components/ContentEditor'
import { useMessage } from 'naive-ui'

const message = useMessage()
const content = ref('')
const editorType = ref<EditorType>('markdown')

function handleContentChange(value: string) {
  console.log('Content updated:', value)
}

function handleTypeChange(type: EditorType) {
  editorType.value = type
  message.info(`Editor type changed to: ${type}`)
}
</script>
*/

// ============================================================
// Example 4: Disabled Editor with Initial Content
// ============================================================
/*
<template>
  <ContentEditor
    v-model="previewContent"
    type="markdown"
    :disabled="true"
    :show-toolbar="false"
    height="300px"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ContentEditor } from '@/components/ContentEditor'

const previewContent = ref(`
# Read-Only Content

This editor is **disabled** and shows preview-only content.

- Item 1
- Item 2
- Item 3
`)
</script>
*/

// ============================================================
// Example 5: Integration with Form
// ============================================================
/*
<template>
  <n-form :model="formData" @submit.prevent="handleSubmit">
    <n-form-item label="Title" path="title">
      <n-input v-model:value="formData.title" />
    </n-form-item>

    <n-form-item label="Content" path="content">
      <ContentEditor
        v-model="formData.content"
        type="markdown"
        height="400px"
        placeholder="Enter article content..."
      />
    </n-form-item>

    <n-form-item>
      <n-button type="primary" attr-type="submit">
        Submit
      </n-button>
    </n-form-item>
  </n-form>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { ContentEditor } from '@/components/ContentEditor'
import { useMessage } from 'naive-ui'

const message = useMessage()

const formData = reactive({
  title: '',
  content: ''
})

function handleSubmit() {
  if (!formData.title || !formData.content) {
    message.warning('Please fill in all fields')
    return
  }

  console.log('Form submitted:', formData)
  message.success('Content saved successfully!')
}
</script>
*/

// ============================================================
// Example 6: Multiple Editors for Different Languages
// ============================================================
/*
<template>
  <n-tabs type="card">
    <n-tab-pane name="en" tab="English">
      <ContentEditor
        v-model="content.en"
        type="markdown"
        placeholder="Enter English content..."
      />
    </n-tab-pane>

    <n-tab-pane name="zh" tab="中文">
      <ContentEditor
        v-model="content.zh"
        type="markdown"
        placeholder="输入中文内容..."
      />
    </n-tab-pane>
  </n-tabs>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { ContentEditor } from '@/components/ContentEditor'

const content = reactive({
  en: '',
  zh: ''
})
</script>
*/

// ============================================================
// Example 7: Auto-save with Debounce
// ============================================================
/*
<template>
  <div>
    <ContentEditor
      v-model="content"
      type="markdown"
      @change="handleContentChange"
    />

    <n-text v-if="saving" depth="3" style="margin-top: 8px">
      Saving...
    </n-text>
    <n-text v-else-if="lastSaved" depth="3" style="margin-top: 8px">
      Last saved: {{ lastSaved }}
    </n-text>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ContentEditor } from '@/components/ContentEditor'
import { useDebounceFn } from '@vueuse/core'

const content = ref('')
const saving = ref(false)
const lastSaved = ref('')

const saveContent = useDebounceFn(async (value: string) => {
  saving.value = true
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    console.log('Content saved:', value)
    lastSaved.value = new Date().toLocaleTimeString()
  } finally {
    saving.value = false
  }
}, 2000)

function handleContentChange(value: string) {
  saveContent(value)
}
</script>
*/

export {}

