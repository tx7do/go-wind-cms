'use client';

import {useTranslations} from 'next-intl';
import LegalPage from '@/components/layout/LegalPage';

export default function PrivacyPage() {
    const t = useTranslations('page.legal.privacy');

    return (
        <LegalPage
            icon="carbon:locked"
            title={t('title')}
            description={t('description')}
            items={[t('item_1'), t('item_2'), t('item_3')]}
        />
    );
}
