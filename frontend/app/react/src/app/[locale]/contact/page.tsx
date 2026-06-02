'use client';

import {useTranslations} from 'next-intl';

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用

export default function ContactPage() {
    const t = useTranslations('page.legal.contact');

    return (
        <div className="w-full">
            {/* Hero Section */}
            <section className="flex min-h-[300px] items-center justify-center border-b border-border bg-gradient-to-br from-primary/10 via-background to-background py-20">
                <div className="w-full max-w-3xl px-8 text-center">
                    <h1 className="mb-4 text-4xl font-bold text-foreground max-md:text-3xl">
                        {t('title')}
                    </h1>
                    <p className="text-lg text-muted-foreground max-md:text-base">
                        {t('description')}
                    </p>
                </div>
            </section>

            {/* Content Section */}
            <div className="w-full my-12 max-w-3xl rounded-xl border border-border bg-card p-8 shadow-sm max-md:mx-4 max-md:p-6">
                <ul className="space-y-4 text-foreground">
                    <li className="flex items-start gap-3">
                        <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary"/>
                        <span>{t('item_1')}</span>
                    </li>
                    <li className="flex items-start gap-3">
                        <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary"/>
                        <span>{t('item_2')}</span>
                    </li>
                    <li className="flex items-start gap-3">
                        <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary"/>
                        <span>{t('item_3')}</span>
                    </li>
                </ul>
            </div>
        </div>
    );
}
