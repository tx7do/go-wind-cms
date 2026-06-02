'use client';

import React from 'react';
import {cn} from '@/lib/utils';

export interface StatItem {
    value: number | string;
    label: string;
}

export interface UserStatGridProps {
    stats: StatItem[];
    className?: string;
}

/**
 * 用户统计栏 — followers/following/posts/likes 的网格展示
 *
 * 用于用户中心页顶部
 */
const UserStatGrid: React.FC<UserStatGridProps> = ({stats, className}) => {
    return (
        <div className={cn(
            'mt-6 grid grid-cols-4 gap-4 border-t border-border pt-6 max-md:grid-cols-2',
            className,
        )}>
            {stats.map((stat, index) => (
                <div key={index} className="text-center">
                    <div className="text-2xl font-bold text-foreground">{stat.value}</div>
                    <div className="text-xs text-muted-foreground">{stat.label}</div>
                </div>
            ))}
        </div>
    );
};

export default UserStatGrid;
