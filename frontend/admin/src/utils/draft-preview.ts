import { marked } from 'marked';

import { EditorType } from '@/components/common/Editor/types';

export interface DraftPreviewOptions {
  title: string;
  /** 正文：富文本为 HTML，Markdown 为源文本，其余类型按纯文本展示 */
  content: string;
  editorType?: EditorType | string;
  /** 顶部横幅文案（说明这是本地草稿、未保存到服务端） */
  bannerText: string;
  /** 浏览器标签页标题 */
  docTitle: string;
}

function escapeHtml(text: string) {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');
}

/**
 * 本地草稿预览：把当前编辑内容渲染成独立 HTML，经 Blob URL 打开新窗口。
 * 纯前端行为——不写服务端、不与 C 端 preview token 链路混淆；
 * 展示的是内容与排版效果，非真实 C 端站点皮肤（横幅文案中说明）。
 */
export function openDraftPreviewWindow(options: DraftPreviewOptions) {
  let bodyHtml: string;
  if (options.editorType === EditorType.MARKDOWN) {
    bodyHtml = marked.parse(options.content, { async: false });
  } else if (options.editorType === EditorType.RICH_TEXT) {
    bodyHtml = options.content;
  } else {
    bodyHtml = `<pre class="draft-plain">${escapeHtml(options.content)}</pre>`;
  }

  const html = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>${escapeHtml(options.docTitle)}</title>
<style>
  body {
    margin: 0;
    background: #f5f6f8;
    color: #1f2328;
    font-family: system-ui, -apple-system, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif;
  }
  .draft-banner {
    position: sticky;
    top: 0;
    padding: 10px 16px;
    text-align: center;
    font-size: 13px;
    color: #874d00;
    background: #fff7e6;
    border-bottom: 1px solid #ffe58f;
  }
  article {
    max-width: 860px;
    margin: 24px auto 48px;
    padding: 40px 48px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  }
  h1.draft-title {
    margin: 0 0 24px;
    font-size: 28px;
    line-height: 1.4;
  }
  article img { max-width: 100%; }
  article pre { overflow: auto; }
  .draft-plain { white-space: pre-wrap; word-break: break-word; font-family: inherit; }
</style>
</head>
<body>
<div class="draft-banner">${escapeHtml(options.bannerText)}</div>
<article>
<h1 class="draft-title">${escapeHtml(options.title)}</h1>
${bodyHtml}
</article>
</body>
</html>`;

  const blob = new Blob([html], { type: 'text/html;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  window.open(url, '_blank', 'noopener');
  // noopener 下拿不到窗口句柄，延迟回收即可：页面加载完成后 revoke 不影响已渲染文档
  window.setTimeout(() => URL.revokeObjectURL(url), 60_000);
}
