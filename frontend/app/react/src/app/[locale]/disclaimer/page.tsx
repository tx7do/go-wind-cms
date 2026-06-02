'use client';

import {useTranslations} from 'next-intl';
import LegalPage from '@/components/layout/LegalPage';

export default function DisclaimerPage() {
    const t = useTranslations('page.legal.disclaimer');

    return (
        <LegalPage
            icon="carbon:warning"
            title={t('title')}
            description={t('description')}
            items={[t('item_1'), t('item_2'), t('item_3')]}
        />
    );
}
