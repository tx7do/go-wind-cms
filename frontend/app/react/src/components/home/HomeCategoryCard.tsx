import React from 'react';
import {XIcon} from '@/plugins/xicon';
import {getCategoryName} from '@/api/hooks/category';
import {contentservicev1_Category} from '@/api/generated/app/service/v1';
import {useTranslations} from 'next-intl';

import {cn} from '@/lib/utils';

interface HomeCategoryCardProps {
    category: contentservicev1_Category;
    onClick?: (id: number) => void;
}

const HomeCategoryCard: React.FC<HomeCategoryCardProps> = ({category, onClick}) => {
    const t = useTranslations('page.home');

    const handleClick = () => {
        if (!category?.id) return;
        onClick?.(category.id);
    };

    const getIconName = (icon?: string): string => {
        if (!icon) return 'carbon:folder';
        return icon.includes(':') ? icon : `carbon:${icon}`;
    };

    return (
        <div
            className={cn(
                'group scroll-reveal-item relative flex min-h-[200px] h-full cursor-pointer flex-col',
                'rounded-2xl border border-primary/10 bg-card p-6 shadow-sm',
                'transition-all duration-300',
                'hover:border-primary/20 hover:shadow-lg hover:shadow-primary/10',
            )}
            onClick={handleClick}
        >
            <div className="relative z-1 flex h-full flex-col">
                <div className="mb-5 flex gap-5">
                    <div className={cn(
                        'mb-4 flex h-[70px] w-[70px] items-center justify-center',
                        'rounded-xl bg-primary/10 text-primary text-3xl',
                        'transition-all duration-300',
                        'group-hover:bg-primary group-hover:text-white',
                    )}>
                        <XIcon
                            name={getIconName(category.icon)}
                            size={48}
                        />
                    </div>
                    <div className="flex-1 w-full">
                        <h3 className="mb-2 text-lg font-extrabold leading-tight tracking-wide text-foreground transition-colors">
                            {getCategoryName(category)}
                        </h3>
                        <span className={cn(
                            'mb-3 inline-block rounded-lg bg-primary/10 px-3 py-1.5',
                            'text-sm font-semibold text-primary transition-all duration-300',
                            'group-hover:bg-primary/15',
                        )}>
                            {t('article_count', {count: category.postCount || 0})}
                        </span>
                    </div>
                </div>
                <div className={cn(
                    'mt-auto w-fit flex items-center gap-2 rounded-lg bg-primary/10 px-3.5 py-2',
                    'text-xs font-semibold text-primary transition-all duration-300',
                    'group-hover:bg-primary/15',
                )}>
                    <XIcon name="carbon:time" size={14}/>
                    <span>{t('updated_days_ago', {days: 3})}</span>
                </div>
            </div>
        </div>
    );
};

export default HomeCategoryCard;
