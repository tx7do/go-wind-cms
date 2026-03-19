import React, {useMemo, useRef, useEffect} from 'react';
import {View, RichText} from '@tarojs/components';
import Taro from '@tarojs/taro';

import type {ContentViewerProps} from '../types';

import './index.scss';

// 简单的 HTML 转义（保留 emoji）
const escapeHtml = (text: string): string => {
  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;',
  };
  return text.replace(/[&<>"']/g, (m) => map[m]);
};

// 简化的 Markdown 解析器
const parseMarkdown = (md: string): string => {
  if (!md) return '';

  let html = md;

  // 处理 URL 后跟中文描述的情况
  html = html.replace(/(https?:\/\/[^\s，]+)(，[^ \n]+)/g, (_match, url, desc) => {
    return `[${url}](${url})${desc}`;
  });

  // 标题
  html = html.replace(/^### (.*$)/gim, '<h3 class="heading-anchor">$1</h3>');
  html = html.replace(/^## (.*$)/gim, '<h2 class="heading-anchor">$1</h2>');
  html = html.replace(/^# (.*$)/gim, '<h1 class="heading-anchor">$1</h1>');

  // 粗体
  html = html.replace(/\*\*(.*?)\*\*/gim, '<strong>$1</strong>');

  // 斜体
  html = html.replace(/\*(.*?)\*/gim, '<em>$1</em>');

  // 删除线
  html = html.replace(/~~(.*?)~~/gim, '<del>$1</del>');

  // 代码块
  html = html.replace(/```([\s\S]*?)```/gim, (_match, code) => {
    return `<pre class="code-block" data-lang="plaintext"><code>${escapeHtml(code)}</code></pre>`;
  });

  // 行内代码
  html = html.replace(/`(.*?)`/gim, '<code class="inline-code">$1</code>');

  // 引用
  html = html.replace(/^> (.*$)/gim, '<blockquote>$1</blockquote>');

  // 列表
  html = html.replace(/^[*\-] (.*$)/gim, '<li>$1</li>');
  html = html.replace(/^\d+\. (.*$)/gim, '<li>$1</li>');

  // 链接
  html = html.replace(/\[(.*?)]\((.*?)\)/gim, (_match, text, href) => {
    const isExternal = href.startsWith('http') || href.startsWith('//');
    return `<a href="${href}" ${isExternal ? 'target="_blank" rel="noopener noreferrer"' : ''} class="markdown-link">${text}</a>`;
  });

  // 图片
  html = html.replace(/!\[(.*?)]\((.*?)\)/gim, (_match, alt, src) => {
    const figcaption = alt ? `<figcaption>${alt}</figcaption>` : '';
    return `<figure class="markdown-image"><img src="${src}" alt="${alt}" class="md-img" />${figcaption}</figure>`;
  });

  // 表格 - 简化处理
  html = html.replace(/\|(.+)\|/gim, (_match, content) => {
    const cells = content.split('|').map((cell: string) => cell.trim());
    return `<tr>${cells.map((cell: any) => `<td>${cell}</td>`).join('')}</tr>`;
  });

  // 段落
  html = html.replace(/^(?!<[hpliuot]|<tr|<td)(.+)$/gim, '<p>$1</p>');

  // 换行
  html = html.replace(/\n/gim, '<br/>');

  return html;
};

const ContentViewer: React.FC<ContentViewerProps> = ({
                                                       content = '',
                                                       type = 'markdown',
                                                       className = ''
                                                     }) => {
  const containerRef = useRef<HTMLDivElement>(null);

  const getRenderedContent = useMemo(() => {
    if (!content) return '';

    try {
      switch (type) {
        case 'markdown':
          return parseMarkdown(content);
        case 'html':
          return content;
        case 'text':
          return `<pre class="plain-text-block">${escapeHtml(content)}</pre>`;
        default:
          return content;
      }
    } catch (error) {
      console.error('Error rendering content:', error);
      return '<p class="content-error">Failed to render content</p>';
    }
  }, [content, type]);

  // 在小程序中需要动态添加 ID 到标题
  useEffect(() => {
    if (containerRef.current && Taro.getEnv() === Taro.ENV_TYPE.WEB) {
      const headings = containerRef.current.querySelectorAll('h2, h3');
      headings.forEach((heading, index) => {
        if (!heading.id) {
          heading.setAttribute('id', `${heading.tagName.toLowerCase()}-${index}`);
        }
      });
    }
  }, [getRenderedContent]);

  return (
    <View ref={containerRef} className={`content-viewer ${className}`}>
      <RichText nodes={getRenderedContent}/>
    </View>
  );
};

export default ContentViewer;
