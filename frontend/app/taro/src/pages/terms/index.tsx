import {useTranslation} from 'react-i18next';
import LegalPage from '@/components/layout/LegalPage';
import {usePageTitle} from '@/hooks/usePageTitle';

export default function TermsPage() {
    const {t} = useTranslation();
    usePageTitle('page.title.terms');
    return (
        <LegalPage
          icon='carbon:document'
          title={t('page.legal.terms.title')}
          description={t('page.legal.terms.description')}
          items={[
                t('page.legal.terms.item_1'),
                t('page.legal.terms.item_2'),
                t('page.legal.terms.item_3'),
            ]}
        />
    );
}
