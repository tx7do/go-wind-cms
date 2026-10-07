import { message } from 'antd';
import i18n from 'i18next';

/**
 * 编辑器/图片上传的全局反馈（对齐 vben api/composables/file-transfer.ts 同名实现）：
 * - 上传进度 toast（antd message 按 key 复用同一实例更新内容）
 * - 失败归一提示：413/Request Entity Too Large → 文件过大；其余 → 上传失败
 * percent < 0 表示进度未知（拿不到 total），只显示上传中。
 */
const UPLOAD_PROGRESS_TOAST_KEY = 'upload-progress-toast';

export function showUploadProgress(percent: number) {
  const text =
    i18n.t('uploadFeedback.uploading', { ns: 'common' }) +
    (percent >= 0 ? ` ${Math.min(99, percent)}%` : '');
  message.loading({
    key: UPLOAD_PROGRESS_TOAST_KEY,
    content: text,
    duration: 0,
  });
}

/** 上传结束（无论成败）都必须调用，避免残留常驻 toast */
export function hideUploadProgress() {
  message.destroy(UPLOAD_PROGRESS_TOAST_KEY);
}

export function notifyUploadError(error: unknown) {
  let text = '';
  if (typeof error === 'string') {
    text = error;
  } else if (error instanceof Error) {
    text = error.message;
  } else {
    try {
      text = JSON.stringify(error ?? {});
    } catch {
      text = '';
    }
  }
  if (/413|Request Entity Too Large/i.test(text)) {
    message.error(i18n.t('uploadFeedback.fileTooLarge', { ns: 'common' }));
  } else {
    message.error(i18n.t('uploadFeedback.uploadFailed', { ns: 'common' }));
  }
}
