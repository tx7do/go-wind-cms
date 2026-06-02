<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'
import { fetchListComment } from '@/api/composables/comment'
import { createComment } from '@/api/service/comment'

const props = defineProps<{
  objectId: number | null
  contentType: string | null
}>()

const { t } = useI18n()

const submitting = ref(false)
const loading = ref(false)
const comments = ref<any[]>([])

const currentPage = ref(1)
const pageSize = 20
const total = ref(0)
const hasMore = ref(true)
const loadingMore = ref(false)

const newComment = reactive({
  content: '',
  authorName: '',
  authorEmail: '',
})

const commentsListRef = ref<HTMLElement | null>(null)

const hasComments = computed(() => comments.value.length > 0)
const showLoadMore = computed(() => hasMore.value && !loadingMore.value)

async function loadComments(reset = false) {
  if (!props.objectId || !props.contentType) return
  if (loading.value || (loadingMore.value && !reset)) return

  if (reset) {
    loading.value = true
  } else {
    loadingMore.value = true
  }

  try {
    const res = await fetchListComment({
      paging: { page: reset ? 1 : currentPage.value, pageSize },
      formValues: {
        objectId: props.objectId,
        contentType: props.contentType,
        status: 'STATUS_APPROVED',
      },
    })
    const newComments = res?.items || []
    const newTotal = res?.total || 0

    if (reset) {
      comments.value = newComments
    } else {
      comments.value = [...comments.value, ...newComments]
    }

    total.value = newTotal
    const allComments = reset ? newComments : [...comments.value, ...newComments]
    hasMore.value = allComments.length < newTotal
    currentPage.value++
  } catch (error) {
    console.error('Load comments failed:', error)
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

function loadMoreComments() {
  if (hasMore.value && !loadingMore.value) {
    loadComments()
  }
}

async function handleSubmitComment() {
  const tmpDiv = document.createElement('div')
  tmpDiv.innerHTML = newComment.content
  const textContent = tmpDiv.textContent?.trim() || ''
  if (!textContent) {
    alert(t('comment.empty_content'))
    return
  }
  if (!newComment.authorName.trim()) {
    alert(t('comment.invalid_nickname'))
    return
  }
  if (!newComment.authorEmail.trim()) {
    alert(t('comment.invalid_email'))
    return
  }
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(newComment.authorEmail)) {
    alert(t('comment.invalid_email'))
    return
  }
  if (!props.objectId || submitting.value) return

  submitting.value = true
  try {
    await createComment({
      objectId: props.objectId,
      contentType: props.contentType ?? undefined,
      content: newComment.content,
      authorName: newComment.authorName,
      authorEmail: newComment.authorEmail,
      status: 'STATUS_PENDING',
    })
    alert(t('comment.comment_submitted'))
    newComment.content = ''
    newComment.authorName = ''
    newComment.authorEmail = ''
    currentPage.value = 1
    await loadComments(true)
    setTimeout(() => {
      commentsListRef.value?.scrollIntoView({ behavior: 'smooth' })
    }, 100)
  } catch (error) {
    console.error('Submit comment failed:', error)
    alert(t('comment.submit_comment_failed'))
  } finally {
    submitting.value = false
  }
}

async function handleReply(comment: any, content: string) {
  if (!content.trim()) {
    alert(t('comment.empty_content'))
    return
  }
  if (!props.objectId || submitting.value) return

  submitting.value = true
  try {
    await createComment({
      objectId: props.objectId,
      contentType: props.contentType ?? undefined,
      content: content.trim(),
      authorName: comment.authorName,
      authorEmail: comment.authorEmail,
      parentId: comment.id,
      replyToId: comment.id,
      status: 'STATUS_PENDING',
    })
    alert(t('comment.comment_posted'))
    currentPage.value = 1
    await loadComments(true)
  } catch (error) {
    console.error('Submit reply failed:', error)
    alert(t('comment.submit_comment_failed'))
  } finally {
    submitting.value = false
  }
}

async function loadChildren(parentComment: any) {
  if (!props.objectId || !props.contentType) return
  try {
    const res = await fetchListComment({
      paging: { page: 1, pageSize: 50 },
      formValues: {
        objectId: props.objectId,
        contentType: props.contentType,
        parentId: parentComment.id,
        status: 'STATUS_APPROVED',
      },
    })
    parentComment.children = res?.items || []
  } catch (error) {
    console.error('Load children failed:', error)
    throw error
  }
}

onMounted(() => {
  if (props.objectId && props.contentType) {
    loadComments(true)
  }
})
</script>

<template>
  <section :class="cn(
    'mx-auto mb-10 max-w-300 rounded-2xl border border-border bg-card p-14 shadow-sm backdrop-blur-sm',
    'max-md:rounded-xl max-md:p-10',
    'max-sm:rounded-xl max-sm:p-8',
  )">
    <!-- Section Header -->
    <div class="mb-10 max-md:mb-8">
      <h2 class="flex items-center gap-3.5 text-3xl font-bold tracking-tight text-foreground max-md:text-2xl max-sm:text-xl">
        <XIcon icon="carbon:chat" :size="36" />
        {{ t('comment.comments_count', { count: comments.length }) }}
      </h2>
    </div>

    <!-- Comment Form -->
    <div :class="cn(
      'relative mb-12 overflow-hidden rounded-2xl border border-primary/10 p-12',
      'bg-linear-to-br from-card to-primary/2 shadow-sm',
      'transition-all duration-400 hover:border-primary hover:shadow-md',
      'max-md:rounded-xl max-md:p-8 max-md:mb-9',
      'max-sm:p-7 max-sm:mb-8',
    )">
      <div class="absolute top-0 left-0 right-0 h-1 bg-primary opacity-90" />

      <div class="mb-8 flex items-center gap-2.5 max-md:mb-6">
        <XIcon icon="carbon:edit" :size="22" class="text-primary" />
        <h3 class="text-lg font-bold tracking-tight text-foreground max-md:text-base">
          {{ t('comment.write_comment') }}
        </h3>
      </div>

      <div class="flex flex-col gap-6">
        <div class="grid grid-cols-2 gap-6 max-md:grid-cols-1 max-md:gap-4">
          <input
            v-model="newComment.authorName"
            :placeholder="t('comment.nickname') + ' *'"
            :class="cn(
              'w-full rounded-lg border border-border bg-background px-4 py-3 text-foreground',
              'transition-all duration-300 hover:border-primary',
              'focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15',
              'disabled:cursor-not-allowed disabled:opacity-60',
            )"
            :disabled="submitting"
          />
          <input
            v-model="newComment.authorEmail"
            type="email"
            :placeholder="t('comment.email') + ' *'"
            :class="cn(
              'w-full rounded-lg border border-border bg-background px-4 py-3 text-foreground',
              'transition-all duration-300 hover:border-primary',
              'focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15',
              'disabled:cursor-not-allowed disabled:opacity-60',
            )"
            :disabled="submitting"
          />
        </div>
        <div class="flex flex-col gap-2">
          <textarea
            v-model="newComment.content"
            rows="5"
            :placeholder="t('comment.write_comment')"
            class="w-full min-h-[120px] resize-y rounded-lg border border-border bg-background px-4 py-3 text-sm text-foreground transition-all duration-300 hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="submitting"
            @keydown.ctrl.enter="handleSubmitComment"
            @keydown.meta.enter="handleSubmitComment"
          />
        </div>
        <div class="flex flex-wrap items-center gap-4 max-md:flex-col max-md:items-stretch">
          <UiButton :disabled="submitting" @click="handleSubmitComment">
            {{ submitting ? '...' : t('comment.submit_comment') }}
          </UiButton>
          <span class="flex items-center gap-2 rounded-lg border border-primary/10 bg-primary/5 px-3.5 py-2 text-[13px] text-muted-foreground max-md:justify-center max-md:text-xs max-md:px-3 max-md:py-1.5">
            <XIcon icon="carbon:information" :size="16" class="text-primary" />
            {{ t('comment.fill_form_info') }}
          </span>
        </div>
      </div>
    </div>

    <!-- Comments List -->
    <div v-if="hasComments" ref="commentsListRef" class="flex flex-col gap-7">
      <CommentTree
        :comments="comments"
        :on-reply="handleReply"
        :on-load-children="loadChildren"
      />

      <div v-if="showLoadMore" class="mt-8 flex justify-center border-t border-primary/8 pt-6">
        <UiButton size="lg" :disabled="loadingMore" @click="loadMoreComments">
          {{ loadingMore ? 'Loading...' : t('comment.load_more') }}
        </UiButton>
      </div>

      <div v-if="!hasMore && hasComments" class="py-5 text-center text-sm font-medium tracking-wide text-muted-foreground">
        {{ t('comment.no_more') }}
      </div>
    </div>

    <UiAppEmpty v-else-if="!loading" :description="t('comment.no_comments')" in-container />
  </section>
</template>
