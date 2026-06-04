import React from 'react';
import {Button} from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {XIcon} from '@/plugins/xicon';
import {useI18n} from '@/i18n';

export const LocaleSwitcher: React.FC = () => {
    const {locale, changeLocale} = useI18n('menu');

    const localeOptions: {key: string; label: string}[] = [
        {key: 'zh-CN', label: '简体中文'},
        {key: 'en-US', label: 'English'},
    ];

    return (
        <DropdownMenu>
            <DropdownMenuTrigger>
                <Button
                  variant='ghost'
                  size='icon'
                  className='h-10 min-w-10 rounded-lg'
                >
                    <XIcon name='carbon:earth' size={16} />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align='end'>
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
