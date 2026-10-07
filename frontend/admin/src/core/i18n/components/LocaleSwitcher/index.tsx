import { Select } from 'antd';
import { useLocaleSync } from '@/core/i18n';
import { supportedLocales } from '@/locales';
import { usePreferencesStore } from '@/core/preferences';

const localeLabels: Record<string, string> = {
  'zh-CN': '🇨🇳 简体中文',
  'en-US': '🇺🇸 English',
};

export const LocaleSwitcher = () => {
  const { changeLocale } = useLocaleSync();

  return (
    <Select
      size="small"
      value={usePreferencesStore((s) => s.preferences.app.locale)}
      options={supportedLocales.map((lng) => ({
        label: localeLabels[lng],
        value: lng,
      }))}
      onChange={(value) => changeLocale(value)}
    />
  );
};
