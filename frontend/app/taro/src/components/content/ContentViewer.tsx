import React, {useMemo} from 'react';
import {View} from '@tarojs/components';
import {marked} from 'marked';
import type {ContentViewerProps} from './types';

// 语言别名
const LANG_ALIAS: Record<string, string> = {
    js: 'javascript',
    ts: 'typescript',
    py: 'python',
    sh: 'bash',
    shell: 'bash',
    yml: 'yaml',
    md: 'markdown',
};

function escapeHtml(text: string): string {
    const map: Record<string, string> = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;',
    };
    return text.replace(/[&<>"']/g, (m) => map[m]);
}

// 简化版 marked 渲染器
const renderer = new marked.Renderer();

renderer.codespan = (code) => {
    return `<code class="inline-code">${escapeHtml(code.text)}</code>`;
};

renderer.code = (code) => {
    const rawLang = (code.lang || 'plaintext').toLowerCase();
    const lang = LANG_ALIAS[rawLang] || rawLang;
    return `<pre class="code-block" data-lang="${rawLang || 'text'}"><code class="language-${lang}">${escapeHtml(code.text)}</code></pre>`;
};

renderer.heading = (heading) => {
    const inlineHtml = marked.parseInline(heading.text);
    return `<h${heading.depth} class="heading-anchor">${inlineHtml}</h${heading.depth}>`;
};

renderer.link = (link) => {
    const isExternal = link.href.startsWith('http') || link.href.startsWith('//');
    return `<View href="${link.href}" ${isExternal ? 'target="_blank" rel="noopener noreferrer"' : ''} class="markdown-link">${escapeHtml(link.text)}</View>`;
};

renderer.image = (image) => {
    return `<figure class="markdown-image"><Image src="${image.href}" alt="${image.text}" class="md-img" />${image.text ? `<figcaption>${image.text}</figcaption>` : ''}</figure>`;
};

renderer.paragraph = (token) => {
    return `<Text>${marked.parseInline(token.text)}</Text>\n`;
};

marked.setOptions({renderer});

const ContentViewer: React.FC<ContentViewerProps> = ({
    content = '',
    type = 'markdown',
    className = '',
}) => {
    const renderedContent = useMemo(() => {
        if (!content) return '';

        try {
            let html = '';
            switch (type) {
                case 'markdown':
                    html = marked.parse(content) as string;
                    return html;
                case 'html':
                    return content;
                case 'text':
                    return `<pre class="plain-text-block">${escapeHtml(content)}</pre>`;
                default:
                    return content;
            }
        } catch (error) {
            console.error('Error rendering content:', error);
            return `<Text class="content-error">Failed to render content</Text>`;
        }
    }, [content, type]);

    return (
        <View
          className={`content-viewer ${className}`}
        >
            <View dangerouslySetInnerHTML={{__html: renderedContent}} />
        </View>
    );
};

export default ContentViewer;
