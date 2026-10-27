import { List, Button } from 'antd';
import { useTranslation } from 'react-i18next';
import MfaManagement from './MfaManagement';

interface SecureSettingProps {
  onSwitchTab: (key: string) => void;
}

const SecureSetting = ({ onSwitchTab }: SecureSettingProps) => {
  const { t } = useTranslation('profile');

  return (
    <>
      <List
        dataSource={[
          {
            key: 'password',
            title: t('secureSetting.password'),
            description: t('secureSetting.passwordDesc'),
          },
        ]}
        renderItem={(item) => (
          <List.Item
            actions={[
              <Button type="link" onClick={() => onSwitchTab('2')}>
                {t('secureSetting.modify')}
              </Button>,
            ]}
          >
            <List.Item.Meta title={item.title} description={item.description} />
          </List.Item>
        )}
      />
      {/* 两步验证（TOTP）绑定/解绑管理 */}
      <MfaManagement />
    </>
  );
};

export default SecureSetting;
