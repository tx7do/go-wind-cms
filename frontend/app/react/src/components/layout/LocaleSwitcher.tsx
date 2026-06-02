'use client';

import React from 'react';
import {Globe} from 'lucide-react';
import {Button} from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

import {useI18n} from '@/i18n';

export const LocaleSwitcher: React.FC = () => {
    const {locale, changeLocale} = useI18n('menu');

    const localeOptions: { key: string; label: string }[] = [
        {key: 'zh-CN', label: '简体中文'},
        {key: 'en-US', label: 'English'},
    ];

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button
                    variant="ghost"
                    size="icon"
                    className="h-10 min-w-10 rounded-lg"
                    aria-label="Language switcher"
                >
                    <Globe className="h-4 w-4"/>
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
                {localeOptions.map(opt => (
                    <DropdownMenuItem
                        key={opt.key}
                        onClick={() => changeLocale(opt.key)}
                        className={locale === opt.key ? 'font-bold' : ''}
                    >
                        {opt.label}
                    </DropdownMenuItem>
                ))}
            </DropdownMenuContent>
        </DropdownMenu>
    );
};
