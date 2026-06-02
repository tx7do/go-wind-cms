'use client';

import {useTranslations} from 'next-intl';
import PageHero from '@/components/layout/PageHero';

export default function TermsPage() {
    const t = useTranslations('page.legal.terms');

    return (
        <div className="w-full">
            <PageHero
                title={t('title')}
                description={t('description')}
                icon="carbon:document"
                size="sm"
            />

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
