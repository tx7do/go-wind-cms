import React from 'react';
import {XIcon} from '@/plugins/xicon';
import {getCategoryName} from '@/api/hooks/category';
import {contentservicev1_Category} from '@/api/generated/app/service/v1';
import {useTranslations} from 'next-intl';
import styles from './HomeCategoryCard.module.css';

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
            className={`${styles.homeCategoryCard} scroll-reveal-item`}
            onClick={handleClick}
        >
            <div className={styles.homeCategoryCardBg}/>
            <div className={styles.homeCategoryCardContent}>
                <div className={styles.homeCategoryCardHeader}>
                    <div className={styles.homeCategoryIcon}>
                        <XIcon
                            name={getIconName(category.icon)}
                            size={48}
                        />
                    </div>
                    <div className={styles.homeCategoryInfo}>
                        <h3>{getCategoryName(category)}</h3>
                        <span className={styles.homePostCount}>
                            {t('article_count', {count: category.postCount || 0})}
                        </span>
                    </div>
                </div>
                <div className={styles.homeCategoryBadge}>
                    <XIcon name="carbon:time" size={14}/>
                    <span>{t('updated_days_ago', {days: 3})}</span>
                </div>
            </div>
        </div>
    );
};

export default HomeCategoryCard;
