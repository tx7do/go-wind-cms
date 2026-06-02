'use client';

import React, {useCallback, useEffect, useMemo, useRef, useState} from 'react';
import XIcon from '@/plugins/xicon';

export interface TocItem {
    id: string;
    level: number;
    text: string;
    index: number;
}

export interface TableOfContentsProps {
    /** 文章内容容器的 ref */
    contentRef: React.RefObject<HTMLDivElement | null>;
    /** 内容标识，变化时触发重新生成目录（如 displayContent 字符串） */
    contentKey: string;
    /** 是否展开 */
    isExpanded: boolean;
    /** 折叠回调 */
    onCollapse: () => void;
    /** 展开回调 */
    onExpand: () => void;
    /** TOC 标题文字 */
    title: string;
}

const HEADING_OFFSET = 150;
const THROTTLE_DELAY = 200;

/**
 * 文章目录侧边栏
 *
 * - 自动从 contentRef 内的 h2/h3 生成目录
 * - 滚动时高亮当前章节
 * - 支持折叠/展开
 * - 点击跳转（用 getBoundingClientRect 计算，避免 scrollIntoView 冲突）
 * - 自动处理 URL hash 跳转
 */
export default function TableOfContents({
    contentRef,
    contentKey,
    isExpanded,
    onCollapse,
    onExpand,
    title,
}: TableOfContentsProps) {
    const [tableOfContents, setTableOfContents] = useState<TocItem[]>([]);
    const [activeHeading, setActiveHeading] = useState('');

    // 从实时 DOM 查询 heading 元素（避免 Shiki 重新渲染后引用过期）
    const getLiveHeading = useCallback((index: number): Element | null => {
        const contentEl = contentRef.current;
        if (!contentEl) return null;
        const headings = contentEl.querySelectorAll('h2, h3');
        return headings[index] || null;
    }, [contentRef]);

    // 生成目录
    const generateToc = useCallback(() => {
        const contentEl = contentRef.current;
        if (!contentEl) return;

        const headings = contentEl.querySelectorAll('h2, h3');
        const toc: TocItem[] = [];

        headings.forEach((heading, index) => {
            const id = `heading-${index}`;
            heading.setAttribute('id', id);

            toc.push({
                id,
                level: heading.tagName === 'H2' ? 2 : 3,
                text: heading.textContent || '',
                index,
            });
        });

        setTableOfContents(toc);
    }, [contentRef]);

    // 内容变化时重新生成（延迟以等待 marked + Shiki 完成）
    useEffect(() => {
        if (!contentKey) return;
        const t1 = setTimeout(generateToc, 100);
        const t2 = setTimeout(generateToc, 1500);
        return () => {
            clearTimeout(t1);
            clearTimeout(t2);
        };
    }, [contentKey, generateToc]);

    // 点击目录时定位
    const scrollToHeading = useCallback((id: string) => {
        const tocItem = tableOfContents.find(item => item.id === id);
        if (!tocItem) return;

        const element = getLiveHeading(tocItem.index);
        if (!element) return;

        element.setAttribute('id', id);

        const rect = element.getBoundingClientRect();
        const HEADER_OFFSET = 80;
        const targetPosition = window.scrollY + rect.top - HEADER_OFFSET;

        window.scrollTo({
            top: Math.max(0, targetPosition),
            behavior: 'smooth',
        });

        setActiveHeading(id);

        if (window.history.pushState) {
            const currentState = window.history.state || {};
            window.history.pushState(currentState, '', `#${id}`);
        }
    }, [tableOfContents, getLiveHeading]);

    // 页面加载时检查 URL hash，自动滚动
    useEffect(() => {
        if (!contentKey || tableOfContents.length === 0 || typeof window === 'undefined' || !window.location.hash) return;
        const hashId = window.location.hash.slice(1);
        const t = setTimeout(() => scrollToHeading(hashId), 500);
        return () => clearTimeout(t);
    }, [contentKey, tableOfContents.length, scrollToHeading]);

    // 节流滚动处理
    const handleScroll = useCallback(() => {
        if (tableOfContents.length === 0) return;
        let currentActive = '';
        for (const tocItem of tableOfContents) {
            const el = getLiveHeading(tocItem.index);
            if (el && el.getBoundingClientRect().top < HEADING_OFFSET) {
                currentActive = tocItem.id;
            }
        }
        if (currentActive !== activeHeading) {
            setActiveHeading(currentActive);
        }
    }, [tableOfContents, activeHeading, getLiveHeading]);

    const throttleTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const throttledScroll = useMemo(() => {
        return () => {
            if (!throttleTimerRef.current) {
                throttleTimerRef.current = setTimeout(() => {
                    handleScroll();
                    throttleTimerRef.current = null;
                }, THROTTLE_DELAY);
            }
        };
    }, [handleScroll]);

    useEffect(() => {
        window.addEventListener('scroll', throttledScroll);
        return () => window.removeEventListener('scroll', throttledScroll);
    }, [throttledScroll]);

    // 不渲染任何内容如果没有目录项
    if (tableOfContents.length === 0) return null;

    return (
        <>
            {/* 侧边栏 */}
            {isExpanded && (
                <aside className="w-60 shrink-0 max-md:hidden">
                    <nav className="sticky top-24 rounded-lg border border-border bg-card p-4 pr-5 border-r border-r-border/50">
                        <div className="mb-3 flex items-center justify-between">
                            <h3 className="flex items-center gap-2 text-sm font-semibold text-foreground">
                                <XIcon name="carbon:list"/>
                                <span>{title}</span>
                            </h3>
                            <button
                                onClick={onCollapse}
                                className="cursor-pointer rounded p-1 text-muted-foreground transition-colors hover:bg-primary/5 hover:text-primary"
                            >
                                <XIcon name="carbon:chevron-left"/>
                            </button>
                        </div>
                        <div className="space-y-0.5">
                            {tableOfContents.map(item => (
                                <a
                                    key={item.id}
                                    href={`#${item.id}`}
                                    className={`block truncate rounded px-2 py-1.5 text-sm transition-colors ${
                                        activeHeading === item.id
                                            ? 'bg-primary/10 font-medium text-primary'
                                            : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'
                                    } ${item.level === 3 ? 'pl-6 text-[13px]' : ''}`}
                                    onClick={(e) => {
                                        e.preventDefault();
                                        setTimeout(() => scrollToHeading(item.id), 10);
                                    }}
                                >
                                    {item.text}
                                </a>
                            ))}
                        </div>
                    </nav>
                </aside>
            )}

            {/* 折叠时的展开按钮 */}
            {!isExpanded && (
                <div className="fixed top-1/2 left-4 z-10 -translate-y-1/2">
                    <button
                        onClick={onExpand}
                        className="flex cursor-pointer items-center gap-2 rounded-lg border border-border bg-card px-3 py-2 text-sm text-foreground shadow-md transition-all hover:border-primary hover:bg-primary/5"
                    >
                        <XIcon name="carbon:list"/>
                        <span className="max-md:hidden">{title}</span>
                        <XIcon name="carbon:chevron-right"/>
                    </button>
                </div>
            )}
        </>
    );
}
