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
                /* 空气动力学悬浮：与 PostCard / FeaturesSection 统一 */
                'transition-all duration-500 ease-out',
                'hover:-translate-y-2 hover:border-primary/30',
                'hover:shadow-[0_20px_40px_-12px_rgba(0,0,0,0.12),0_0_24px_-6px_hsl(var(--primary)/0.2)]',
            )}
            onClick={handleClick}
            style={{ willChange: 'transform, box-shadow' }}
        >
            {/* hover 风迹渐变蒙层 */}
            <div className="pointer-events-none absolute inset-0 bg-gradient-to-tr from-primary/5 via-transparent to-sky-400/5 opacity-0 transition-opacity duration-500 group-hover:opacity-100 rounded-2xl"/>
            <div className="relative z-1 flex h-full flex-col">
                <div className="mb-5 flex gap-5">
                    <div className={cn(
                        'mb-4 flex h-[70px] w-[70px] items-center justify-center',
                        'rounded-xl bg-primary/10 text-primary text-3xl',
                        'transition-all duration-500 ease-out',
                        'group-hover:scale-110 group-hover:bg-primary group-hover:text-white',
                        'group-hover:shadow-[0_8px_24px_-4px_hsl(var(--primary)/0.4)]',
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
