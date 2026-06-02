<script setup lang="ts">
import { XIcon } from '@/plugins/xicon'
import { cn } from '@/lib/utils'
import { formatDate } from '@/utils/date'

const props = defineProps<{
  comments: any[]
  onReply: (comment: any, content: string) => void
  onLoadChildren: (comment: any) => Promise<void>
}>()

const { t } = useI18n()

const replyingCommentId = ref<number | null>(null)
const replyContent = ref('')
const submitting = ref(false)
const expandedComments = ref<Set<number>>(new Set())
const loadingChildren = ref<Set<number>>(new Set())
const likedComments = ref<Set<number>>(new Set())

function hasChildren(comment: any): boolean {
  return !!comment.children || (comment.replyCount !== undefined && comment.replyCount > 0)
}

function isOwnerReply(comment: any): boolean {
  return comment.authorType === 'AUTHOR_TYPE_ADMIN' && !!comment.replyToId
}

function isExpanded(comment: any): boolean {
  return expandedComments.value.has(comment.id || 0)
}

function isLoading(comment: any): boolean {
  return loadingChildren.value.has(comment.id || 0)
}

function handleReply(comment: any) {
  replyingCommentId.value = comment.id || 0
  replyContent.value = ''
}

function cancelReply() {
  replyingCommentId.value = null
  replyContent.value = ''
}

async function submitReply(comment: any) {
  if (!replyContent.value.trim()) {
    alert(t('comment.empty_content'))
    return
  }
  if (submitting.value) return
  submitting.value = true
  try {
    props.onReply(comment, replyContent.value.trim())
    cancelReply()
  } catch (error) {
    console.error('Submit reply failed:', error)
    alert(t('comment.submit_comment_failed'))
  } finally {
    submitting.value = false
  }
}

function handleLike(comment: any) {
  const commentId = comment.id || 0
  if (likedComments.value.has(commentId)) {
    likedComments.value.delete(commentId)
  } else {
    likedComments.value.add(commentId)
  }
}

function handleShare(comment: any) {
  const baseUrl = window.location.href.split('#')[0]
  const commentUrl = `${baseUrl}#comment-${comment.id}`
  navigator.clipboard.writeText(commentUrl).catch(() => {})
}

async function toggleExpand(comment: any) {
  if (isExpanded(comment)) {
    expandedComments.value.delete(comment.id || 0)
  } else {
    if (!comment.children && comment.replyCount && comment.replyCount > 0) {
      loadingChildren.value.add(comment.id || 0)
      try {
        await props.onLoadChildren(comment)
        expandedComments.value.add(comment.id || 0)
      } catch (error) {
        console.error('Load children failed:', error)
      } finally {
        loadingChildren.value.delete(comment.id || 0)
      }
    } else {
      expandedComments.value.add(comment.id || 0)
    }
  }
}
</script>

<template>
  <div v-if="comments && comments.length > 0" class="flex flex-col gap-7">
    <div v-for="comment in comments" :key="comment.id" class="flex flex-col">
      <!-- Comment body -->
      <div :class="cn(
        'group relative flex gap-5 overflow-hidden rounded-2xl border border-border bg-linear-to-br from-card to-primary/2 p-7',
        'shadow-sm transition-all duration-400',
        'hover:border-primary hover:shadow-md hover:translate-x-1 hover:-translate-y-0.5',
        'max-md:p-6 max-md:gap-4',
        'max-sm:p-4 max-sm:gap-3',
      )">
        <div class="absolute top-0 left-0 h-full w-1 bg-primary opacity-0 transition-opacity duration-300 group-hover:opacity-100" />

        <div class="shrink-0">
          <UiAvatar class="h-12 w-12 ring-2 ring-border/30 shadow-md">
            <UiAvatarFallback class="bg-primary text-primary-foreground font-semibold">
              {{ comment.authorName?.charAt(0) || 'U' }}
            </UiAvatarFallback>
          </UiAvatar>
        </div>
        <div class="min-w-0 flex-1">
          <!-- Header -->
          <div :class="cn(
            'mb-4 flex items-center justify-between border-b border-primary/8 pb-3',
            'max-md:flex-col max-md:items-start max-md:gap-2 max-md:mb-3 max-md:pb-2.5',
            'max-sm:gap-1.5 max-sm:mb-2 max-sm:pb-1.5',
          )">
            <div class="flex flex-wrap items-center gap-3 max-md:flex-col max-md:items-start max-md:gap-1.5">
              <strong class="flex items-center gap-2 text-[17px] font-semibold tracking-tight text-foreground max-md:text-base max-sm:text-[15px]">
                <span class="text-sm font-bold text-primary">@</span>
                <span v-if="isOwnerReply(comment)" class="inline-flex items-center gap-1 rounded-lg bg-primary px-2.5 py-0.5 text-[11px] font-bold text-primary-foreground shadow-sm max-sm:text-[10px] max-sm:px-1.5 max-sm:py-0.5">
                  <XIcon icon="carbon:badge" :size="13" />
                  {{ t('comment.owner_reply') }}
                </span>
                {{ comment.authorName }}
              </strong>
              <span v-if="comment.location && !isOwnerReply(comment)" class="flex items-center gap-1 text-xs text-muted-foreground">
                {{ comment.location }}
              </span>
            </div>
            <span class="flex items-center gap-1.5 text-[13px] text-muted-foreground max-md:text-xs max-sm:text-[11px]">
              {{ formatDate(comment.createdAt) }}
            </span>
          </div>

          <!-- Content -->
          <div class="mb-4 pl-1 text-[15px] leading-relaxed text-foreground max-md:text-sm max-sm:text-[13px]">
            {{ comment.content }}
          </div>

          <!-- Actions -->
          <div :class="cn(
            'flex gap-5 border-t border-primary/5 pt-3',
            'max-md:flex-wrap max-md:gap-4',
            'max-sm:gap-2.5',
          )">
            <span
              :class="cn(
                'flex cursor-pointer items-center gap-1.5 text-[13px] text-muted-foreground transition-all duration-200 select-none hover:text-primary hover:-translate-y-0.5',
                'max-md:text-xs max-sm:text-[11px]',
                likedComments.has(comment.id || 0) && 'text-primary font-semibold',
              )"
              @click="handleLike(comment)"
            >
              <XIcon icon="carbon:thumbs-up" :size="16" />
              {{ comment.likeCount || 0 }}
            </span>
            <span
              class="flex cursor-pointer items-center gap-1.5 text-[13px] text-muted-foreground transition-all duration-200 select-none hover:text-primary hover:-translate-y-0.5 max-md:text-xs max-sm:text-[11px]"
              @click="handleReply(comment)"
            >
              <XIcon icon="carbon:chat" :size="16" />
              {{ t('comment.reply') }}
            </span>
            <span
              class="flex cursor-pointer items-center gap-1.5 text-[13px] text-muted-foreground transition-all duration-200 select-none hover:text-primary hover:-translate-y-0.5 max-md:text-xs max-sm:text-[11px]"
              @click="handleShare(comment)"
            >
              <XIcon icon="carbon:share" :size="16" />
              {{ t('comment.share') }}
            </span>
            <span
              v-if="comment.replyCount && comment.replyCount > 0"
              class="flex cursor-pointer items-center gap-1.5 text-[13px] font-medium text-primary transition-all duration-200 select-none hover:text-primary max-md:text-xs max-sm:text-[11px]"
              @click="toggleExpand(comment)"
            >
              <XIcon :icon="isExpanded(comment) ? 'carbon:chevron-up' : 'carbon:chevron-down'" :size="18" />
              {{ isExpanded(comment) ? t('comment.hide_replies') : t('comment.view_replies', { count: comment.replyCount }) }}
            </span>
          </div>

          <!-- Reply form -->
          <div v-if="replyingCommentId === comment.id" class="mt-4 border-t border-primary/8 pt-4 max-md:mt-3 max-md:pt-3">
            <textarea
              v-model="replyContent"
              rows="3"
              :placeholder="t('comment.write_comment')"
              class="w-full min-h-20 resize-y rounded-lg border border-border bg-muted px-4 py-3 text-sm text-foreground transition-all duration-300 hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15 disabled:cursor-not-allowed disabled:opacity-60 disabled:bg-background"
              :disabled="submitting"
              @keydown.ctrl.enter="submitReply(comment)"
              @keydown.meta.enter="submitReply(comment)"
            />
            <div class="mt-2 mb-3 text-right text-xs text-muted-foreground">
              {{ replyContent.length }} / 1000
            </div>
            <div class="flex justify-end gap-3 max-md:flex-col max-md:gap-2">
              <UiButton size="sm" :disabled="submitting" @click="submitReply(comment)">
                {{ submitting ? '...' : t('comment.submit_comment') }}
              </UiButton>
              <UiButton variant="outline" size="sm" :disabled="submitting" @click="cancelReply">
                {{ t('comment.cancel') }}
              </UiButton>
            </div>
          </div>
        </div>
      </div>

      <!-- Recursive children -->
      <div v-if="hasChildren(comment) && isExpanded(comment)" :class="cn(
        'relative mt-5 pl-16',
        'max-md:pl-12 max-md:mt-4',
        'max-sm:pl-9 max-sm:mt-3',
      )">
        <div class="absolute top-0 left-8 h-full w-0.5 bg-linear-to-b from-primary/20 to-primary/5 max-md:left-6 max-sm:left-4.5 max-sm:w-[1.5px]" />
        <div v-if="isLoading(comment)" class="flex items-center justify-center gap-3 py-5 text-sm text-muted-foreground">
          <UiSpinner size="sm" class="text-primary" />
          <span>{{ t('comment.loading') }}</span>
        </div>
        <CommentTree
          v-if="!isLoading(comment) && comment.children"
          :comments="comment.children"
          :on-reply="onReply"
          :on-load-children="onLoadChildren"
        />
      </div>
    </div>
  </div>
</template>
