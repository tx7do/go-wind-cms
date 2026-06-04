import {useTranslation} from 'react-i18next';
import LegalPage from '@/components/layout/LegalPage';

export default function DisclaimerPage() {
    const {t} = useTranslation();
    return (
        <LegalPage
            icon='carbon:warning'
            title={t('page.legal.disclaimer.title')}
            description={t('page.legal.disclaimer.description')}
            items={[
                t('page.legal.disclaimer.item_1'),
                t('page.legal.disclaimer.item_2'),
                t('page.legal.disclaimer.item_3'),
            ]}
        />
    );
}
