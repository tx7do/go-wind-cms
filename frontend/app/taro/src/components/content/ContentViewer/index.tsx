import React from 'react';
import {View, RichText} from '@tarojs/components';

import type {ContentViewerProps} from '../types';

import './index.scss';

const ContentViewer: React.FC<ContentViewerProps> = ({
                                                       content = '',
                                                       type = 'markdown',
                                                       className = ''
                                                     }) => {
  // 简化的 Markdown 转 HTML（仅支持基础语法）
  const parseMarkdown = (md: string): string => {
    if (!md) return '';

    let html = md;

    // 标题
    html = html.replace(/^### (.*$)/gim, '<h3>$1</h3>');
    html = html.replace(/^## (.*$)/gim, '<h2>$1</h2>');
    html = html.replace(/^# (.*$)/gim, '<h1>$1</h1>');

    // 粗体
    html = html.replace(/\*\*(.*?)\*\*/gim, '<strong>$1</strong>');

    // 斜体
    html = html.replace(/\*(.*?)\*/gim, '<em>$1</em>');

    // 删除线
    html = html.replace(/~~(.*?)~~/gim, '<del>$1</del>');

    // 代码块
    html = html.replace(/```([\s\S]*?)```/gim, '<pre class="code-block"><code>$1</code></pre>');

    // 行内代码
    html = html.replace(/`(.*?)`/gim, '<code class="inline-code">$1</code>');

    // 引用
    html = html.replace(/^> (.*$)/gim, '<blockquote>$1</blockquote>');

    // 无序列表
    html = html.replace(/^\* (.*$)/gim, '<li>$1</li>');
    html = html.replace(/^- (.*$)/gim, '<li>$1</li>');

    // 链接
    html = html.replace(/\[(.*?)\]\((.*?)\)/gim, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>');

    // 图片
    html = html.replace(/!\[(.*?)\]\((.*?)\)/gim, '<img src="$2" alt="$1" class="markdown-image" />');

    // 换行
    html = html.replace(/\n/gim, '<br/>');

    return html;
  };

  const getRenderedContent = () => {
    if (!content) return '';

    try {
      switch (type) {
        case 'markdown':
          return parseMarkdown(content);
        case 'html':
          return content;
        case 'text':
          return `<pre class="plain-text-block">${content}</pre>`;
        default:
          return content;
      }
    } catch (error) {
      console.error('Error rendering content:', error);
      return '<p class="content-error">Failed to render content</p>';
    }
  };

  return (
    <View className={`content-viewer ${className}`}>
      <RichText nodes={getRenderedContent()}/>
    </View>
  );
};

export default ContentViewer;
