<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import {useMessage} from 'naive-ui'

import {$t} from '@/locales'
import {useCommentStore} from '@/stores'
import type {
  commentservicev1_Comment,
  commentservicev1_Comment_ContentType
} from "@/api/generated/app/service/v1"

// --- 类型定义 ---
interface CommentForm {
  content: string
  authorName: string
  authorEmail: string
}

// --- Props ---
const props = defineProps<{
  objectId: number | null
  contentType: commentservicev1_Comment_ContentType | null
}>()

// --- Emits ---
const emit = defineEmits<{
  (e: 'update:comments', comments: commentservicev1_Comment[]): void
}>()

// --- Store ---
const commentStore = useCommentStore()

// --- 状态 ---
const message = useMessage()
const submitting = ref(false)
const loading = ref(false)
const comments = ref<commentservicev1_Comment[]>([])

// 分页相关状态
const currentPage = ref(1)
const pageSize = ref(20) // 每页加载 20 条评论
const total = ref(0)
const hasMore = ref(true)
const loadingMore = ref(false)

const newComment = ref<CommentForm>({
  content: '',
  authorName: '',
  authorEmail: '',
})

// --- 计算属性 ---
const displayComments = computed(() => comments.value || [])

const hasComments = computed(() => displayComments.value.length > 0)

const showLoadMore = computed(() => hasMore.value && !loadingMore.value)

// 加载评论 (分页)
async function loadComments(reset = false) {
  if (!props.objectId || !props.contentType) return
  if (loading.value || (loadingMore.value && !reset)) return

  if (reset) {
    loading.value = true
  } else {
    loadingMore.value = true
  }

  try {
    const res = await commentStore.listComment(
      {
        page: reset ? 1 : currentPage.value,
        pageSize: pageSize.value
      },
      {
        objectId: props.objectId,
        contentType: props.contentType,
        status: 'STATUS_APPROVED'
      }
    )

    const newComments = res.items || []
    total.value = res.total || 0

    if (reset) {
      comments.value = newComments
    } else {
      comments.value = [...comments.value, ...newComments]
    }

    // 判断是否还有更多数据
    hasMore.value = comments.value.length < total.value
    currentPage.value++

    emit('update:comments', comments.value)
  } catch (error) {
    console.error('Load comments failed:', error)
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

// 加载更多
function loadMoreComments() {
  if (hasMore.value && !loadingMore.value) {
    loadComments()
  }
}

async function handleSubmitComment() {
  if (!newComment.value.content || !newComment.value.authorName || !newComment.value.authorEmail) {
    message.warning($t('comment.fill_form_info'))
    return
  }
  if (!props.objectId || submitting.value) return

  submitting.value = true
  try {
    await commentStore.createComment({
      postId: props.objectId,
      content: newComment.value.content,
      authorName: newComment.value.authorName,
      authorEmail: newComment.value.authorEmail,
      status: 'COMMENT_STATUS_PENDING',
    })
    message.success($t('comment.comment_submitted'))
    newComment.value = {content: '', authorName: '', authorEmail: ''}
    currentPage.value = 1
    await loadComments(true)
    // 用户体验优化：提交后滚动到评论列表
    setTimeout(() => {
      document.querySelector('.comments-list')?.scrollIntoView({behavior: 'smooth'})
    }, 100)
  } catch (error) {
    console.error('Submit comment failed:', error)
    message.error($t('comment.submit_comment_failed'))
  } finally {
    submitting.value = false
  }
}

// 处理回复
async function handleReply(comment: commentservicev1_Comment, content: string) {
  if (!content.trim()) {
    message.warning($t('comment.empty_content'))
    return
  }

  if (!props.objectId || submitting.value) return

  submitting.value = true
  try {
    await commentStore.createComment({
      postId: props.objectId,
      content: content.trim(),
      authorName: comment.authorName, // 使用被回复评论的作者名
      authorEmail: comment.authorEmail, // 使用被回复评论的邮箱
      parentId: comment.id, // 父评论 ID
      replyToId: comment.id, // 回复的目标评论 ID
      status: 'COMMENT_STATUS_PENDING',
    })
    message.success($t('comment.comment_posted'))
    currentPage.value = 1
    await loadComments(true)
  } catch (error) {
    console.error('Submit reply failed:', error)
    message.error($t('comment.submit_comment_failed'))
  } finally {
    submitting.value = false
  }
}

// 加载子评论
async function loadChildren(parentComment: commentservicev1_Comment) {
  if (!props.objectId || !props.contentType) return

  try {
    const res = await commentStore.listComment(
      {
        page: 1,
        pageSize: 50 // 一次性加载所有子评论
      },
      {
        objectId: props.objectId,
        contentType: props.contentType,
        parentId: parentComment.id, // 根据父评论 ID 过滤
        status: 'STATUS_APPROVED'
      }
    )

    // 将子评论添加到父评论的 children 中
    parentComment.children = res.items || []
  } catch (error) {
    console.error('Load children failed:', error)
    throw error
  }
}

// --- 生命周期 ---
onMounted(() => {
  loadComments(true)
})
</script>

<template>
  <section class="comments-section">
    <div class="section-header">
      <h2>
        <span class="i-carbon:chat"/>
        {{ $t('comment.comments_count', {count: displayComments.length}) }}
      </h2>
    </div>

    <!-- Comment Form -->
    <div class="comment-form">
      <div class="form-header">
        <div class="form-icon">
          <span class="i-carbon:edit"/>
        </div>
        <h3>{{ $t('comment.write_comment') }}</h3>
      </div>
      <n-form>
        <n-grid :cols="2" :x-gap="24">
          <n-form-item-gi>
            <n-input
              v-model:value="newComment.authorName"
              :placeholder="$t('comment.nickname') + ' *'"
              size="large"
              :disabled="submitting"
              clearable
            >
              <template #prefix>
                <span class="i-carbon:user"/>
              </template>
            </n-input>
          </n-form-item-gi>
          <n-form-item-gi>
            <n-input
              v-model:value="newComment.authorEmail"
              :placeholder="$t('comment.email') + ' *'"
              type="email"
              size="large"
              :disabled="submitting"
              clearable
            >
              <template #prefix>
                <span class="i-carbon:email"/>
              </template>
            </n-input>
          </n-form-item-gi>
        </n-grid>
        <n-form-item>
          <n-input
            v-model:value="newComment.content"
            type="textarea"
            :rows="6"
            :placeholder="$t('comment.write_comment')"
            size="large"
            :disabled="submitting"
            show-count
            maxlength="1000"
          />
        </n-form-item>
        <n-form-item>
          <div class="form-actions">
            <n-button
              size="large"
              @click="handleSubmitComment"
              :loading="submitting"
              round
            >
              <template #icon>
                <span class="i-carbon:send-alt"/>
              </template>
              {{ $t('comment.submit_comment') }}
            </n-button>
            <span class="form-tip">
              <span class="i-carbon:information"/>
              {{ $t('comment.fill_form_info') }}
            </span>
          </div>
        </n-form-item>
      </n-form>
    </div>

    <!-- Comments List -->
    <div v-if="hasComments" class="comments-list">
      <CommentTree
        :comments="displayComments"
        @reply="handleReply"
        @load-children="loadChildren"
      />

      <!-- 加载更多按钮 -->
      <div v-if="showLoadMore" class="load-more-container">
        <n-button
          size="large"
          @click="loadMoreComments"
          :loading="loadingMore"
        >
          {{ $t('comment.load_more') }}
        </n-button>
      </div>

      <!-- 没有更多提示 -->
      <div v-else-if="!hasMore && hasComments" class="no-more-text">
        {{ $t('comment.no_more') }}
      </div>
    </div>
    <n-empty v-else :description="$t('comment.no_comments')"
             style="margin-top: 40px;"/>
  </section>
</template>

<style scoped lang="less">
@keyframes gradient-shift {
  0% {
    background-position: 0 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0 50%;
  }
}

// Section Header
.section-header {
  margin-bottom: 40px;
  position: relative;

  h2 {
    display: flex;
    align-items: center;
    gap: 14px;
    font-size: 30px;
    font-weight: 700;
    margin: 0;
    color: var(--color-text-primary);
    letter-spacing: -0.3px;

    span[class^="i-"] {
      font-size: 36px;
      color: var(--color-brand);
      filter: drop-shadow(0 2px 8px rgba(168, 85, 247, 0.3));
    }
  }
}

// Comments Section
.comments-section {
  max-width: 1200px;
  margin: 0 auto 40px;
  background: var(--color-surface);
  border-radius: 20px;
  padding: 56px 72px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06),
  0 8px 48px rgba(0, 0, 0, 0.04),
  inset 0 1px 0 rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

// Comment Form
.comment-form {
  background: linear-gradient(135deg,
  var(--color-surface) 0%,
  rgba(168, 85, 247, 0.02) 100%);
  border-radius: 20px;
  padding: 48px;
  margin-bottom: 48px;
  border: 1px solid rgba(168, 85, 247, 0.12);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08),
  0 2px 16px rgba(168, 85, 247, 0.06),
  inset 0 1px 0 rgba(255, 255, 255, 0.1);
  position: relative;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: linear-gradient(90deg,
    var(--color-brand) 0%,
    #a855f7 25%,
    #764ba2 50%,
    #a855f7 75%,
    var(--color-brand) 100%);
    background-size: 200% 100%;
    opacity: 0.9;
    animation: gradient-shift 3s ease infinite;
  }

  &:hover {
    border-color: var(--color-brand);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.1),
    0 4px 20px rgba(168, 85, 247, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
    transform: translateY(-2px);
  }

  .form-header {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 32px;
    padding-bottom: 24px;
    border-bottom: 2px solid rgba(168, 85, 247, 0.08);

    .form-icon {
      width: 48px;
      height: 48px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg,
      var(--color-brand) 0%,
      #a855f7 100%);
      border-radius: 14px;
      box-shadow: 0 4px 16px rgba(168, 85, 247, 0.3);
      color: #fff;

      span[class^="i-"] {
        font-size: 24px;
      }
    }

    h3 {
      font-size: 24px;
      font-weight: 700;
      margin: 0;
      color: var(--color-text-primary);
      letter-spacing: -0.3px;
      flex: 1;
    }
  }

  :deep(.n-input) {
    border-radius: 12px;
    background: var(--color-bg);
    border: 2px solid rgba(168, 85, 247, 0.08);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    .n-input__border {
      border: none;
    }

    &:hover {
      border-color: rgba(168, 85, 247, 0.2);
      transform: translateY(-1px);
    }

    &:focus-within {
      border-color: var(--color-brand);
      box-shadow: 0 0 0 3px rgba(168, 85, 247, 0.08);
      transform: translateY(0);
    }

    .n-input__prefix {
      color: var(--color-brand);
      font-size: 18px;
      margin-right: 8px;
    }

    .n-input__count-suffix {
      color: var(--color-text-secondary);
      font-size: 12px;
      background: rgba(168, 85, 247, 0.05);
      padding: 2px 8px;
      border-radius: 6px;
    }
  }

  :deep(.n-form-item) {
    margin-bottom: 24px;
    transition: all 0.3s;

    &:hover {
      transform: translateX(3px);
    }
  }

  .form-actions {
    display: flex;
    align-items: center;
    gap: 16px;
    flex-wrap: wrap;

    :deep(.n-button) {
      min-width: 180px;
      height: 50px;
      font-size: 16px;
      font-weight: 600;
      border-radius: 25px;
      background: linear-gradient(135deg,
      var(--color-brand) 0%,
      #a855f7 100%);
      border: none;
      box-shadow: 0 4px 16px rgba(168, 85, 247, 0.3);
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      color: #fff;

      &:hover {
        transform: translateY(-3px);
        box-shadow: 0 8px 24px rgba(168, 85, 247, 0.4);
      }

      &:active {
        transform: translateY(-1px);
      }

      &[disabled] {
        opacity: 0.6;
        cursor: not-allowed;
      }
    }

    .form-tip {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 13px;
      color: var(--color-text-secondary);
      padding: 8px 14px;
      background: rgba(168, 85, 247, 0.05);
      border-radius: 10px;
      border: 1px solid rgba(168, 85, 247, 0.08);

      span[class^="i-"] {
        font-size: 16px;
        color: var(--color-brand);
      }
    }
  }
}

html.dark .comment-form .form-actions {
  :deep(.n-button) {
    background: linear-gradient(135deg,
    rgba(99, 102, 241, 0.85) 0%,
    rgba(168, 85, 247, 0.7) 100%);
    box-shadow: 0 4px 16px rgba(99, 102, 241, 0.35);
    border: 1px solid rgba(168, 85, 247, 0.35);

    &:hover {
      transform: translateY(-3px);
      box-shadow: 0 8px 24px rgba(99, 102, 241, 0.45);
      border-color: rgba(168, 85, 247, 0.5);
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .comments-section {
    padding: 40px 24px !important;
    border-radius: 16px !important;
  }

  .section-header {
    margin-bottom: 32px;

    h2 {
      font-size: 24px;

      span[class^="i-"] {
        font-size: 28px;
      }
    }
  }

  .comment-form {
    padding: 32px 24px !important;
    margin-bottom: 36px !important;
    border-radius: 16px !important;

    .form-header {
      margin-bottom: 24px;
      padding-bottom: 20px;

      .form-icon {
        width: 40px;
        height: 40px;
        border-radius: 12px;

        span[class^="i-"] {
          font-size: 20px;
        }
      }

      h3 {
        font-size: 20px;
      }
    }
  }

  :deep(.n-input) {
    font-size: 15px;
  }

  .form-actions {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;

    :deep(.n-button) {
      min-width: 100%;
      height: 48px;
      font-size: 15px;
    }

    .form-tip {
      justify-content: center;
      font-size: 12px;
      padding: 6px 12px;
    }
  }
}

@media (max-width: 640px) {
  .comments-section {
    padding: 32px 20px !important;
    border-radius: 14px !important;
  }

  .section-header {
    margin-bottom: 28px;

    h2 {
      font-size: 22px;
      gap: 10px;

      span[class^="i-"] {
        font-size: 26px;
      }
    }
  }

  .comment-form {
    padding: 28px 20px !important;
    margin-bottom: 32px !important;
    border-radius: 14px !important;

    .form-header {
      margin-bottom: 20px;
      padding-bottom: 16px;

      .form-icon {
        width: 36px;
        height: 36px;
        border-radius: 10px;

        span[class^="i-"] {
          font-size: 18px;
        }
      }

      h3 {
        font-size: 18px;
      }
    }
  }

  :deep(.n-grid) {
    gap: 16px !important;
  }

  :deep(.n-input) {
    font-size: 14px;

    &.n-input--textarea {
      .n-input__textarea-el {
        min-height: 100px !important;
        font-size: 14px;
      }
    }
  }

  .form-actions {
    :deep(.n-button) {
      height: 46px;
      font-size: 15px;
    }

    .form-tip {
      font-size: 11px;
      padding: 6px 10px;
    }
  }

  .load-more-container {
    margin-top: 28px;
    padding-top: 20px;

    :deep(.n-button) {
      min-width: 140px;
      height: 42px;
      font-size: 14px;
    }
  }

  .no-more-text {
    padding: 16px;
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .comments-section {
    padding: 28px 16px !important;
    border-radius: 12px !important;
  }

  .section-header {
    margin-bottom: 24px;

    h2 {
      font-size: 20px;

      span[class^="i-"] {
        font-size: 24px;
      }
    }
  }

  .comment-form {
    padding: 24px 16px !important;
    margin-bottom: 28px !important;

    .form-header {
      margin-bottom: 16px;
      padding-bottom: 12px;

      .form-icon {
        width: 32px;
        height: 32px;
        border-radius: 8px;

        span[class^="i-"] {
          font-size: 16px;
        }
      }

      h3 {
        font-size: 17px;
      }
    }
  }

  :deep(.n-form-item) {
    margin-bottom: 20px;
  }

  :deep(.n-input) {
    font-size: 13px;

    &.n-input--textarea {
      .n-input__textarea-el {
        min-height: 90px !important;
        font-size: 13px;
      }
    }
  }

  .form-actions {
    :deep(.n-button) {
      height: 44px;
      font-size: 14px;
    }

    .form-tip {
      font-size: 11px;
      padding: 5px 8px;
    }
  }

  .load-more-container {
    margin-top: 24px;
    padding-top: 16px;

    :deep(.n-button) {
      min-width: 120px;
      height: 40px;
      font-size: 13px;
    }
  }

  .no-more-text {
    padding: 14px;
    font-size: 12px;
  }
}

// Comments List
.comments-list {
  display: flex;
  flex-direction: column;
  gap: 28px;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 64px;
    width: 2px;
    height: 100%;
    background: linear-gradient(180deg,
    rgba(168, 85, 247, 0.2) 0%,
    rgba(168, 85, 247, 0.05) 100%);
    z-index: 0;
  }
}

// 加载更多按钮
.load-more-container {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid rgba(168, 85, 247, 0.08);

  :deep(.n-button) {
    min-width: 160px;
    height: 44px;
    font-size: 15px;
    font-weight: 600;
    border-radius: 12px;
    background: linear-gradient(135deg,
    var(--color-brand) 0%,
    #764ba2 100%);
    border: none;
    box-shadow: 0 4px 12px rgba(168, 85, 247, 0.2);
    color: #fff;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 16px rgba(168, 85, 247, 0.3);
    }

    &:active {
      transform: translateY(0);
    }
  }
}

html.dark .load-more-container {
  border-top: 1px solid rgba(168, 85, 247, 0.15);

  :deep(.n-button) {
    background: linear-gradient(135deg,
    rgba(99, 102, 241, 0.8) 0%,
    rgba(139, 92, 246, 0.6) 100%);
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.25);
    border: 1px solid rgba(168, 85, 247, 0.3);

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 16px rgba(99, 102, 241, 0.35);
      border-color: rgba(168, 85, 247, 0.5);
    }

    &:active {
      transform: translateY(0);
    }
  }
}

// 没有更多提示
.no-more-text {
  text-align: center;
  padding: 20px;
  color: var(--color-text-secondary);
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.5px;
}

.comment-item {
  display: flex;
  gap: 20px;
  padding: 28px;
  background: linear-gradient(135deg,
  var(--color-bg) 0%,
  rgba(168, 85, 247, 0.02) 100%);
  border-radius: 16px;
  border: 1px solid rgba(168, 85, 247, 0.08);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 4px;
    height: 100%;
    background: linear-gradient(180deg,
    var(--color-brand) 0%,
    #764ba2 100%);
    opacity: 0;
    transition: opacity 0.3s;
  }

  &:hover {
    border-color: var(--color-brand);
    box-shadow: 0 8px 24px rgba(168, 85, 247, 0.12),
    0 4px 12px rgba(0, 0, 0, 0.04);
    transform: translateX(4px) translateY(-2px);

    &::before {
      opacity: 1;
    }
  }

  .comment-avatar {
    flex-shrink: 0;
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;

    &::after {
      content: '';
      position: absolute;
      top: 50%;
      left: 50%;
      width: calc(100% + 6px);
      height: calc(100% + 6px);
      transform: translate(-50%, -50%);
      border-radius: 50%;
      background: linear-gradient(135deg,
      var(--color-brand) 0%,
      #764ba2 100%);
      opacity: 0.15;
      filter: blur(6px);
      z-index: 0;
    }

    :deep(.n-avatar) {
      position: relative;
      z-index: 1;
      background: linear-gradient(135deg,
      var(--color-brand) 0%,
      #764ba2 100%);
      color: #fff;
      font-weight: 600;
      box-shadow: 0 2px 8px rgba(168, 85, 247, 0.2);
      border: 1px solid rgba(255, 255, 255, 0.1);
    }
  }

  .comment-body {
    flex: 1;
    min-width: 0;
    position: relative;

    .comment-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      padding-bottom: 12px;
      border-bottom: 1px solid rgba(168, 85, 247, 0.08);

      .comment-author {
        font-size: 17px;
        font-weight: 600;
        color: var(--color-text-primary);
        letter-spacing: -0.2px;
        display: flex;
        align-items: center;
        gap: 8px;

        &::before {
          content: '@';
          color: var(--color-brand);
          font-size: 14px;
          font-weight: 700;
        }
      }

      .comment-date {
        font-size: 13px;
        color: var(--color-text-secondary);
        display: flex;
        align-items: center;
        gap: 6px;

        &::before {
          content: '📅';
          font-size: 12px;
        }
      }
    }

    .comment-content {
      font-size: 15px;
      line-height: 1.8;
      color: var(--color-text-primary);
      padding-left: 4px;
    }
  }
}

// Responsive Design
@media (max-width: 1024px) {
  .comments-section {
    padding: 48px 56px;
  }

  .comment-form {
    padding: 36px;
  }
}

@media (max-width: 768px) {
  .comments-section {
    padding: 40px 36px;
    border-radius: 16px;
  }

  .section-header {
    margin-bottom: 32px;

    h2 {
      font-size: 26px;

      span[class^="i-"] {
        font-size: 30px;
      }
    }
  }

  .comment-form {
    padding: 28px;
    margin-bottom: 36px;
    border-radius: 14px;

    h3 {
      font-size: 19px;
      margin-bottom: 24px;
    }

    :deep(.n-form-item) {
      margin-bottom: 24px;
    }

    :deep(.n-grid) {
      gap: 14px !important;
    }

    :deep(.n-button) {
      padding: 0 32px;
      height: 46px;
      font-size: 14px;
    }
  }

  .comments-list {
    gap: 24px;

    &::before {
      left: 56px;
    }
  }

  .comment-item {
    padding: 24px;
    gap: 16px;
    border-radius: 14px;

    .comment-avatar {
      :deep(.n-avatar) {
        --n-size: 44px !important;
      }
    }

    .comment-body {
      .comment-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
        margin-bottom: 12px;
        padding-bottom: 10px;

        .comment-author {
          font-size: 16px;
        }

        .comment-date {
          font-size: 12px;
        }
      }

      .comment-content {
        font-size: 14px;
        line-height: 1.7;
      }
    }
  }
}

@media (max-width: 640px) {
  .comments-section {
    padding: 32px 28px;
    border-radius: 14px;
  }

  .section-header {
    margin-bottom: 24px;

    h2 {
      font-size: 22px;
      gap: 10px;

      span[class^="i-"] {
        font-size: 26px;
      }
    }
  }

  .comment-form {
    padding: 20px;
    margin-bottom: 28px;
    border-radius: 12px;

    h3 {
      font-size: 17px;
      margin-bottom: 20px;
    }

    :deep(.n-form-item) {
      margin-bottom: 20px;
    }

    :deep(.n-form-item-label) {
      font-size: 13px;
    }

    :deep(.n-input) {
      font-size: 13px;

      &.n-input--textarea {
        .n-input__textarea-el {
          min-height: 100px !important;
        }
      }
    }

    :deep(.n-button) {
      height: 44px;
      font-size: 14px;
      padding: 0 24px;
    }
  }

  .comments-list {
    gap: 20px;

    &::before {
      left: 48px;
    }
  }

  .comment-item {
    padding: 18px;
    gap: 12px;
    border-radius: 12px;

    .comment-avatar {
      :deep(.n-avatar) {
        --n-size: 40px !important;
        font-size: 14px;
      }
    }

    .comment-body {
      .comment-header {
        gap: 6px;
        margin-bottom: 10px;
        padding-bottom: 8px;

        .comment-author {
          font-size: 15px;
        }

        .comment-date {
          font-size: 11px;
        }
      }

      .comment-content {
        font-size: 13px;
        line-height: 1.65;
      }
    }
  }
}

@media (max-width: 480px) {
  .comments-section {
    padding: 28px 20px;
    border-radius: 12px;
  }

  .comment-form {
    padding: 16px;
    margin-bottom: 24px;
    border-radius: 10px;

    h3 {
      font-size: 16px;
    }
  }

  .comment-item {
    padding: 16px;
    gap: 10px;
    border-radius: 10px;

    .comment-avatar {
      :deep(.n-avatar) {
        --n-size: 36px !important;
      }
    }
  }
}
</style>
