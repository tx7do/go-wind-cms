'use client';

import {useTranslations} from 'next-intl';
import LegalPage from '@/components/layout/LegalPage';

export default function ContactPage() {
    const t = useTranslations('page.legal.contact');

    return (
        <LegalPage
            icon="carbon:email"
            title={t('title')}
            description={t('description')}
            items={[t('item_1'), t('item_2'), t('item_3')]}
        />
    );
}
