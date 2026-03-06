<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import {useMessage} from 'naive-ui'

import {$t} from '@/locales'
import {useCommentStore} from '@/stores/modules/app'
import type {
  commentservicev1_Comment,
  commentservicev1_Comment_ContentType
} from "@/api/generated/app/service/v1"

// --- 类型定义 ---
interface CommentForm {
  content: string
  authorName: string
  email: string
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

const newComment = ref<CommentForm>({
  content: '',
  authorName: '',
  email: '',
})

// --- 计算属性 ---
const displayComments = computed(() => comments.value || [])

const hasComments = computed(() => displayComments.value.length > 0)

async function loadComments() {
  if (!props.objectId || !props.contentType) return
  if (loading.value) return

  loading.value = true
  try {
    const res = await commentStore.listComment(
      {
        page: 1,
        pageSize: 50
      },
      {
        objectId: props.objectId,
        contentType: props.contentType,
        status: 'STATUS_APPROVED'
      }
    )
    comments.value = res.items || []
    emit('update:comments', comments.value)
  } catch (error) {
    console.error('Load comments failed:', error)
  } finally {
    loading.value = false
  }
}

async function handleSubmitComment() {
  if (!newComment.value.content || !newComment.value.authorName || !newComment.value.email) {
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
      email: newComment.value.email,
      status: 'COMMENT_STATUS_PENDING',
    })
    message.success($t('comment.comment_submitted'))
    newComment.value = {content: '', authorName: '', email: ''}
    await loadComments()
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

// --- 生命周期 ---
onMounted(() => {
  loadComments()
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
      <h3>{{ $t('comment.write_comment') }}</h3>
      <n-form>
        <n-grid :cols="2" :x-gap="16">
          <n-form-item-gi :label="$t('comment.nickname')">
            <n-input
              v-model:value="newComment.authorName"
              :placeholder="$t('comment.enter_nickname')"
              size="large"
              :disabled="submitting"
            />
          </n-form-item-gi>
          <n-form-item-gi :label="$t('comment.email')">
            <n-input
              v-model:value="newComment.email"
              :placeholder="$t('comment.enter_email')"
              type="text"
              size="large"
              :disabled="submitting"
            />
          </n-form-item-gi>
        </n-grid>
        <n-form-item :label="$t('comment.comment_content')">
          <n-input
            v-model:value="newComment.content"
            type="textarea"
            :rows="5"
            :placeholder="$t('comment.write_comment')"
            size="large"
            :disabled="submitting"
          />
        </n-form-item>
        <n-form-item>
          <n-button
            type="primary"
            size="large"
            @click="handleSubmitComment"
            :loading="submitting"
          >
            <template #icon>
              <span class="i-carbon:send-alt"/>
            </template>
            {{ $t('comment.submit_comment') }}
          </n-button>
        </n-form-item>
      </n-form>
    </div>

    <!-- Comments List -->
    <div v-if="hasComments" class="comments-list">
      <CommentTree :comments="displayComments"/>
    </div>
    <n-empty v-else :description="$t('comment.no_comments')"
             style="margin-top: 40px;"/>
  </section>
</template>

<style scoped lang="less">
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
  box-shadow:
    0 4px 24px rgba(0, 0, 0, 0.06),
    0 8px 48px rgba(0, 0, 0, 0.04),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

// Comment Form
.comment-form {
  background: linear-gradient(135deg,
    var(--color-bg) 0%,
    rgba(168, 85, 247, 0.03) 100%);
  border-radius: 16px;
  padding: 40px;
  margin-bottom: 48px;
  border: 1px solid rgba(168, 85, 247, 0.15);
  box-shadow:
    0 4px 16px rgba(0, 0, 0, 0.04),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg,
      var(--color-brand) 0%,
      #764ba2 50%,
      var(--color-brand) 100%);
    opacity: 0.8;
  }

  h3 {
    font-size: 21px;
    font-weight: 600;
    margin: 0 0 28px 0;
    color: var(--color-text-primary);
    letter-spacing: -0.2px;
    display: flex;
    align-items: center;
    gap: 10px;

    &::after {
      content: '';
      flex: 1;
      height: 1px;
      background: linear-gradient(90deg,
        var(--color-border) 0%,
        transparent 100%);
    }
  }

  :deep(.n-form-item) {
    margin-bottom: 28px;
    transition: all 0.3s;

    &:hover {
      transform: translateX(2px);
    }
  }

  :deep(.n-button) {
    padding: 0 36px;
    height: 48px;
    font-size: 15px;
    font-weight: 600;
    border-radius: 12px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 4px 12px rgba(168, 85, 247, 0.2);

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 20px rgba(168, 85, 247, 0.3);
    }

    &:active {
      transform: translateY(0);
    }
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
    box-shadow:
      0 8px 24px rgba(168, 85, 247, 0.12),
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
