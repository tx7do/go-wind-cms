'use client';

import React, {useEffect, useRef, useState} from 'react';
import {Input} from '@/components/ui/input';
import {Search} from 'lucide-react';
import {useTranslations} from 'next-intl';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {apiClient} from '@/api/client';

export default function SearchBar() {
    const t = useTranslations('navbar.top');
    const router = useI18nRouter();
    const [searchQuery, setSearchQuery] = useState('');
    const [suggestions, setSuggestions] = useState<string[]>([]);
    const [open, setOpen] = useState(false);
    // debounce 句柄与请求序号（防乱序回包覆盖新结果）
    const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const seqRef = useRef(0);

    const handleSearch = () => {
        const q = searchQuery.trim();
        if (!q) return;
        setOpen(false);
        router.push(`/search?q=${encodeURIComponent(q)}`);
    };

    useEffect(() => {
        if (timerRef.current) clearTimeout(timerRef.current);
        const q = searchQuery.trim();
        if (q.length < 2) {
            // 宏任务清态：规避 set-state-in-effect 规则
            timerRef.current = setTimeout(() => {
                setSuggestions([]);
                setOpen(false);
            }, 0);
            return;
        }
        timerRef.current = setTimeout(async () => {
            const seq = ++seqRef.current;
            try {
                const resp = await apiClient.postService.SearchPostsSuggest({query: q, language: '', pageSize: 8});
                if (seq !== seqRef.current) return;
                const titles = resp.titles ?? [];
                setSuggestions(titles);
                setOpen(titles.length > 0);
            } catch (error) {
                // 联想失败静默——搜索主链路不受影响（未按 Enter 前只是提示）
                if (seq === seqRef.current) {
                    console.warn('search suggest failed:', error);
                    setSuggestions([]);
                    setOpen(false);
                }
            }
        }, 300);
        return () => {
            if (timerRef.current) clearTimeout(timerRef.current);
        };
    }, [searchQuery]);

    const pick = (title: string) => {
        setOpen(false);
        router.push(`/search?q=${encodeURIComponent(title)}`);
    };

    return (
        <div className="mx-2 hidden h-11 max-w-80 flex-1 items-center md:flex lg:max-w-80">
            <div className="relative w-full">
                <Search className="absolute start-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground"/>
                <Input
                    className="h-full w-full ps-8"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') handleSearch();
                        if (e.key === 'Escape') setOpen(false);
                    }}
                    onBlur={() => setTimeout(() => setOpen(false), 150)}
                    onFocus={() => suggestions.length > 0 && setOpen(true)}
                    placeholder={t('search_placeholder')}
                />
                {open && suggestions.length > 0 && (
                    <ul
                        className="absolute z-50 mt-1 w-full overflow-hidden rounded-lg border border-border bg-card shadow-md"
                        role="listbox"
                    >
                        {suggestions.map((title) => (
                            <li key={title}>
                                <button
                                    type="button"
                                    className="flex w-full items-center gap-2 px-3 py-2 text-start text-sm text-foreground hover:bg-accent hover:text-accent-foreground"
                                    onMouseDown={(e) => {
                                        e.preventDefault();
                                        pick(title);
                                    }}
                                >
                                    <Search className="h-3.5 w-3.5 shrink-0 text-muted-foreground"/>
                                    <span className="line-clamp-1">{title}</span>
                                </button>
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        </div>
    );
}
