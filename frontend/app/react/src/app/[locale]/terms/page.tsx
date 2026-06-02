'use client';

import {useTranslations} from 'next-intl';
import LegalPage from '@/components/layout/LegalPage';

export default function TermsPage() {
    const t = useTranslations('page.legal.terms');

    return (
        <LegalPage
            icon="carbon:document"
            title={t('title')}
            description={t('description')}
            items={[t('item_1'), t('item_2'), t('item_3')]}
        />
    );
}
