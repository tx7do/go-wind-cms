<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import { Bold, Italic, Code, List, Strikethrough } from 'lucide-vue-next'
import { XIcon } from '@/plugins/xicon'

const props = withDefaults(defineProps<{
  modelValue: string
  submitting?: boolean
  placeholder?: string
  maxLength?: number
  submitLabel?: string
}>(), {
  submitting: false,
  placeholder: '',
  maxLength: 1000,
  submitLabel: '',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'submit': []
}>()

const editor = useEditor({
  extensions: [
    StarterKit,
    Placeholder.configure({
      placeholder: props.placeholder,
    }),
  ],
  content: props.modelValue,
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  },
  editorProps: {
    attributes: {
      class: 'prose dark:prose-invert focus:outline-none min-h-[120px] max-h-[300px] overflow-y-auto p-4 text-sm text-foreground',
    },
  },
})

// 外部 modelValue 变化时同步到编辑器（如提交后清空）
watch(() => props.modelValue, (val) => {
  if (!editor.value) return
  const currentHTML = editor.value.getHTML()
  if (val !== currentHTML) {
    editor.value.commands.setContent(val || '')
  }
})

const textLength = computed(() => editor.value?.getText().length ?? 0)
const isEmpty = computed(() => textLength.value === 0)
const isOverLimit = computed(() => textLength.value > props.maxLength)

function handleKeyDown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
    e.preventDefault()
    if (!isEmpty.value && !props.submitting) {
      emit('submit')
    }
  }
}

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<template>
  <div
    v-if="editor"
    class="w-full overflow-hidden rounded-xl border border-border bg-card transition-all focus-within:border-primary focus-within:ring-1 focus-within:ring-primary/40"
    @keydown="handleKeyDown"
  >
    <!-- Toolbar -->
    <div class="flex items-center gap-1 border-b border-border bg-muted/50 p-1.5">
      <button
        type="button"
        :class="[
          'inline-flex h-8 w-8 items-center justify-center rounded-md transition-colors',
          editor.isActive('bold')
            ? 'bg-muted text-foreground'
            : 'text-muted-foreground hover:bg-muted hover:text-foreground',
        ]"
        @click="editor.commands.toggleBold()"
      >
        <Bold class="h-4 w-4" />
      </button>

      <button
        type="button"
        :class="[
          'inline-flex h-8 w-8 items-center justify-center rounded-md transition-colors',
          editor.isActive('italic')
            ? 'bg-muted text-foreground'
            : 'text-muted-foreground hover:bg-muted hover:text-foreground',
        ]"
        @click="editor.commands.toggleItalic()"
      >
        <Italic class="h-4 w-4" />
      </button>

      <button
        type="button"
        :class="[
          'inline-flex h-8 w-8 items-center justify-center rounded-md transition-colors',
          editor.isActive('strike')
            ? 'bg-muted text-foreground'
            : 'text-muted-foreground hover:bg-muted hover:text-foreground',
        ]"
        @click="editor.commands.toggleStrike()"
      >
        <Strikethrough class="h-4 w-4" />
      </button>

      <button
        type="button"
        :class="[
          'inline-flex h-8 w-8 items-center justify-center rounded-md transition-colors',
          editor.isActive('codeBlock')
            ? 'bg-muted text-foreground'
            : 'text-muted-foreground hover:bg-muted hover:text-foreground',
        ]"
        @click="editor.commands.toggleCodeBlock()"
      >
        <Code class="h-4 w-4" />
      </button>

      <button
        type="button"
        :class="[
          'inline-flex h-8 w-8 items-center justify-center rounded-md transition-colors',
          editor.isActive('bulletList')
            ? 'bg-muted text-foreground'
            : 'text-muted-foreground hover:bg-muted hover:text-foreground',
        ]"
        @click="editor.commands.toggleBulletList()"
      >
        <List class="h-4 w-4" />
      </button>

      <div class="ml-auto flex items-center gap-2 text-xs text-muted-foreground">
        <XIcon icon="carbon:keyboard" :size="14" />
        <span class="max-sm:hidden">Ctrl + Enter</span>
      </div>
    </div>

    <!-- Editor Content -->
    <EditorContent :editor="editor" />

    <!-- Bottom Bar -->
    <div class="flex items-center justify-between border-t border-border/50 bg-muted/20 p-2">
      <span :class="['text-xs', isOverLimit ? 'text-destructive' : 'text-muted-foreground']">
        {{ textLength }} / {{ maxLength }}
      </span>
      <UiButton
        size="sm"
        :disabled="submitting || isEmpty || isOverLimit"
        class="gap-2"
        @click="emit('submit')"
      >
        <XIcon v-if="submitting" icon="carbon:circle-dash" class="h-4 w-4 animate-spin" />
        <XIcon v-else icon="carbon:send-alt" :size="16" />
        {{ submitLabel }}
      </UiButton>
    </div>
  </div>
</template>

<style scoped>
:deep(.tiptap) {
  outline: none;
}

:deep(.tiptap p.is-editor-empty:first-child::before) {
  content: attr(data-placeholder);
  float: left;
  color: hsl(var(--muted-foreground));
  pointer-events: none;
  height: 0;
}

:deep(.tiptap p) {
  margin: 0.25em 0;
}

:deep(.tiptap strong) {
  font-weight: 700;
}

:deep(.tiptap em) {
  font-style: italic;
}

:deep(.tiptap s) {
  text-decoration: line-through;
}

:deep(.tiptap code) {
  background: hsl(var(--muted));
  border-radius: 0.375rem;
  padding: 0.15em 0.4em;
  font-size: 0.875em;
  font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
}

:deep(.tiptap pre) {
  background: hsl(var(--muted));
  border-radius: 0.5rem;
  padding: 0.75em 1em;
  overflow-x: auto;
}

:deep(.tiptap pre code) {
  background: transparent;
  padding: 0;
  font-size: 0.85em;
}

:deep(.tiptap ul) {
  list-style-type: disc;
  padding-left: 1.5em;
}

:deep(.tiptap li) {
  margin: 0.15em 0;
}
</style>
