<script setup lang="ts">
import {useMessage} from 'naive-ui'
import {ref} from 'vue'

import {$t} from '@/locales'
import {ContentViewer} from '@/components/ContentViewer'
import type {commentservicev1_Comment} from "@/api/generated/app/service/v1"
import {formatDate} from "@/utils/date";

// --- Props ---
defineProps<{
  comments: commentservicev1_Comment[]
}>()

// --- Emits ---
const emit = defineEmits<{
  (e: 'reply', comment: commentservicev1_Comment, content: string): void
}>()

// --- 状态 ---
const message = useMessage()
const replyingCommentId = ref<number | null>(null)
const replyContent = ref('')
const submitting = ref(false)

function hasChildren(comment: commentservicev1_Comment): boolean {
  return !!comment.children && comment.children.length > 0
}

function isOwnerReply(comment: commentservicev1_Comment): boolean {
  // 博主回复：authorType 为 ADMIN 且有 replyToId (回复他人)
  return comment.authorType === 'AUTHOR_TYPE_ADMIN' && !!comment.replyToId
}

// --- 方法 ---
function handleReply(comment: commentservicev1_Comment) {
  replyingCommentId.value = comment.id
  replyContent.value = ''
}

function cancelReply() {
  replyingCommentId.value = null
  replyContent.value = ''
}

async function submitReply(comment: commentservicev1_Comment) {
  if (!replyContent.value.trim()) {
    message.warning($t('comment.empty_content'))
    return
  }

  if (submitting.value) return

  submitting.value = true
  try {
    emit('reply', comment, replyContent.value.trim())
    cancelReply()
    message.success($t('comment.comment_posted'))
  } catch (error) {
    console.error('Submit reply failed:', error)
    message.error($t('comment.submit_comment_failed'))
  } finally {
    submitting.value = false
  }
}

function handleShare(comment: commentservicev1_Comment) {
  // 复制链接到剪贴板
  const url = window.location.href.split('#')[0]
  const commentUrl = `${url}#comment-${comment.id}`

  navigator.clipboard.writeText(commentUrl).then(() => {
    message.success($t('comment.link_copied'))
  }).catch(() => {
    message.error($t('comment.copy_failed'))
  })
}
</script>

<template>
  <div class="comment-tree">
    <div v-for="comment in comments" :key="comment.id" class="comment-node">
      <!-- 评论主体 -->
      <div class="comment-item">
        <div class="comment-avatar">
          <n-avatar :size="48" round>
            {{ comment.authorName?.charAt(0) || 'U' }}
          </n-avatar>
        </div>
        <div class="comment-body">
          <div class="comment-header">
            <div class="author-info">
              <strong class="comment-author">
                <span v-if="isOwnerReply(comment)" class="owner-badge">
                  <span class="i-carbon:badge"/>
                  {{ $t('comment.owner_reply') }}
                </span>
                {{ comment.authorName }}
              </strong>
              <span v-if="comment.location && !isOwnerReply(comment)"
                    class="comment-location">{{ comment.location }}</span>
            </div>
            <span class="comment-date">{{ formatDate(comment.createdAt) }}</span>
          </div>
          <div class="comment-content">
            <ContentViewer :content="comment.content" type="text"/>
          </div>
          <div class="comment-actions">
            <span class="action-item">
              <span class="i-carbon:thumbs-up"/>
              {{ comment.likeCount || 0 }}
            </span>
            <span class="action-item" @click="handleReply(comment)">
              <span class="i-carbon:chat"/>
              {{ $t('comment.reply') }}
            </span>
            <span class="action-item" @click="handleShare(comment)">
              <span class="i-carbon:share"/>
              {{ $t('comment.share') }}
            </span>
          </div>

          <!-- 回复表单 -->
          <div v-if="replyingCommentId === comment.id" class="reply-form">
            <n-input
              v-model:value="replyContent"
              type="textarea"
              :rows="3"
              :placeholder="$t('comment.write_comment')"
              size="small"
              :disabled="submitting"
              @keydown.ctrl.enter="submitReply(comment)"
            />
            <div class="reply-form-actions">
              <n-button
                size="small"
                type="primary"
                @click="submitReply(comment)"
                :loading="submitting"
              >
                {{ $t('comment.submit_comment') }}
              </n-button>
              <n-button
                size="small"
                @click="cancelReply"
                :disabled="submitting"
              >
                {{ $t('comment.cancel') }}
              </n-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 递归渲染子评论 -->
      <div v-if="hasChildren(comment)" class="comment-children">
        <CommentTree :comments="comment.children!"/>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.comment-tree {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.comment-node {
  display: flex;
  flex-direction: column;
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

      .author-info {
        display: flex;
        align-items: center;
        gap: 12px;
        flex-wrap: wrap;

        .comment-author {
          font-size: 17px;
          font-weight: 600;
          color: var(--color-text-primary);
          letter-spacing: -0.2px;
          display: flex;
          align-items: center;
          gap: 8px;

          .owner-badge {
            background: linear-gradient(135deg,
            var(--color-brand) 0%,
            #764ba2 100%);
            color: #fff;
            font-size: 11px;
            font-weight: 700;
            padding: 3px 10px;
            border-radius: 8px;
            letter-spacing: 0;
            display: inline-flex;
            align-items: center;
            gap: 4px;
            box-shadow: 0 2px 8px rgba(168, 85, 247, 0.3);

            span[class^="i-"] {
              font-size: 13px;
            }
          }
        }

        .comment-location {
          font-size: 12px;
          color: var(--color-text-secondary);
          display: flex;
          align-items: center;
          gap: 4px;

          &::before {
            content: '📍';
            font-size: 11px;
          }
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
      margin-bottom: 16px;
    }

    .comment-actions {
      display: flex;
      gap: 20px;
      padding-top: 12px;
      border-top: 1px solid rgba(168, 85, 247, 0.05);

      .action-item {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        color: var(--color-text-secondary);
        cursor: pointer;
        transition: all 0.2s;

        span[class^="i-"] {
          font-size: 16px;
        }

        &:hover {
          color: var(--color-brand);
          transform: translateY(-1px);
        }
      }
    }

    // 回复表单
    .reply-form {
      margin-top: 16px;
      padding-top: 16px;
      border-top: 1px solid rgba(168, 85, 247, 0.08);

      :deep(.n-input) {
        border-radius: 8px;
        background: var(--color-bg);
      }

      .reply-form-actions {
        display: flex;
        gap: 12px;
        margin-top: 12px;
        justify-content: flex-end;

        :deep(.n-button) {
          min-width: 80px;
        }
      }
    }
  }
}

// 子评论嵌套样式
.comment-children {
  margin-top: 20px;
  padding-left: 64px;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 32px;
    width: 2px;
    height: 100%;
    background: linear-gradient(180deg,
    rgba(168, 85, 247, 0.2) 0%,
    rgba(168, 85, 247, 0.05) 100%);
  }

  .comment-node {
    margin-bottom: 20px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .comment-item {
    padding: 24px;
    gap: 16px;

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

        .author-info {
          flex-direction: column;
          align-items: flex-start;
          gap: 6px;
        }

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

      .comment-actions {
        gap: 16px;

        .action-item {
          font-size: 12px;
        }
      }
    }
  }

  .comment-children {
    padding-left: 48px;
    margin-top: 16px;

    &::before {
      left: 24px;
    }
  }
}

@media (max-width: 640px) {
  .comment-item {
    padding: 18px;
    gap: 12px;

    .comment-avatar {
      :deep(.n-avatar) {
        --n-size: 40px !important;
      }
    }

    .comment-body {
      .comment-header {
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

  .comment-children {
    padding-left: 36px;
    margin-top: 12px;

    &::before {
      left: 18px;
    }
  }
}
</style>
