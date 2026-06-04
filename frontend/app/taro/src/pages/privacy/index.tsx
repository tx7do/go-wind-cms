import {useTranslation} from 'react-i18next';
import LegalPage from '@/components/layout/LegalPage';
import {usePageTitle} from '@/hooks/usePageTitle';

export default function PrivacyPage() {
    const {t} = useTranslation();
    usePageTitle('page.title.privacy');
    return (
        <LegalPage
          icon='carbon:locked'
          title={t('page.legal.privacy.title')}
          description={t('page.legal.privacy.description')}
          items={[
                t('page.legal.privacy.item_1'),
                t('page.legal.privacy.item_2'),
                t('page.legal.privacy.item_3'),
            ]}
        />
    );
}
