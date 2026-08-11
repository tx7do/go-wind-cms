'use client';

import React, {useState} from 'react';
import {Input} from '@/components/ui/input';
import {Search} from 'lucide-react';
import {useTranslations} from 'next-intl';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export default function SearchBar() {
    const t = useTranslations('navbar.top');
    const router = useI18nRouter();
    const [searchQuery, setSearchQuery] = useState('');

    const handleSearch = () => {
        const q = searchQuery.trim();
        if (!q) return;
        router.push(`/search?q=${encodeURIComponent(q)}`);
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
                    }}
                    placeholder={t('search_placeholder')}
                />
            </div>
        </div>
    );
}
