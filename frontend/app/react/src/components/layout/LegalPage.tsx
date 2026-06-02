'use client';

import React from 'react';
import PageHero from './PageHero';

interface LegalPageProps {
    icon: string;
    title: string;
    description: string;
    items: string[];
}

/**
 * 法定/静态信息页面通用模板
 * 用于 privacy / terms / disclaimer / contact 等结构一致的页面
 */
export default function LegalPage({icon, title, description, items}: LegalPageProps) {
    return (
        <div className="w-full">
            <PageHero
                title={title}
                description={description}
                icon={icon}
                size="sm"
            />
            <div className="w-full my-12 max-w-3xl rounded-xl border border-border bg-card p-8 shadow-sm max-md:mx-4 max-md:p-6">
                <ul className="space-y-4 text-foreground">
                    {items.map((item, index) => (
                        <li key={index} className="flex items-start gap-3">
                            <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary"/>
                            <span>{item}</span>
                        </li>
                    ))}
                </ul>
            </div>
        </div>
    );
}
