'use client';

import React, {useState} from 'react';
import {Input} from '@/components/ui/input';
import {Search} from 'lucide-react';
import {useTranslations} from 'next-intl';

import styles from './SearchBar.module.css';

export default function SearchBar() {
    const t = useTranslations('navbar.top');
    const [searchQuery, setSearchQuery] = useState('');

    const handleSearch = () => {
        console.log('Searching for:', searchQuery);
    };

    return (
        <div className={styles.searchBarWrapper}>
            <div className="relative w-full">
                <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground"/>
                <Input
                    className={`${styles.searchBar} pl-8`}
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    onKeyUp={handleSearch}
                    placeholder={t('search_placeholder')}
                />
            </div>
        </div>
    );
}
