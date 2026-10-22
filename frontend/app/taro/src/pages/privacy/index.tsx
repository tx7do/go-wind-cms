import {useTranslation} from 'react-i18next';
import {usePageTitle} from '@/hooks/usePageTitle';
import LegalPage from '@/components/layout/LegalPage';

export default function PrivacyPage() {
    const {t} = useTranslation();
    usePageTitle('page.title.privacy');

    return (
        <LegalPage
          icon='carbon:shield-alert'
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
