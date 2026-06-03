import { defineNuxtPlugin } from '#app';
import { initPreferences, preferencesManager } from '@/core/preferences';
import type { SupportedLanguagesType } from '@/core/preferences/types';

function normalizeLocale(code: string): SupportedLanguagesType {
  if (code.includes('-')) return code as SupportedLanguagesType;
  const map: Record<string, SupportedLanguagesType> = { zh: 'zh-CN', en: 'en-US' };
  return (map[code] || code + '-CN') as SupportedLanguagesType;
}

export default defineNuxtPlugin(async (nuxtApp) => {
  if (import.meta.server) return;

  const i18n = nuxtApp.$i18n as any;

  // 在 initPreferences 之前，先用 i18n 的当前 locale 覆盖 preferences 缓存
  // 确保直接访问 /en-US/xxx 时 locale 不会是缓存的旧值
  if (i18n?.locale?.value) {
    const normalized = normalizeLocale(i18n.locale.value as string);
    try {
      const storageKey = 'gowind-cms:preferences';
      const raw = localStorage.getItem(storageKey);
      if (raw) {
        const cached = JSON.parse(raw);
        cached.app = cached.app || {};
        cached.app.locale = normalized;
        localStorage.setItem(storageKey, JSON.stringify(cached));
      }
      localStorage.setItem('gowind-cms:app-locale', JSON.stringify(normalized));
    } catch {
      // ignore
    }
  }

  await initPreferences({
    namespace: 'gowind-cms',
  });

  // 初始化后立即同步一次（防止 restoreLocaleFromI18n 未生效）
  if (i18n?.locale?.value) {
    const normalized = normalizeLocale(i18n.locale.value as string);
    const current = preferencesManager.getPreferences().app.locale;
    if (current !== normalized) {
      preferencesManager.updatePreferences({
        app: { locale: normalized },
      });
    }
  }

  // 监听 i18n locale 变化，持续同步到 preferencesManager
  if (i18n?.locale) {
    watch(
      () => i18n.locale.value as string,
      (newLocale: string) => {
        const normalized = normalizeLocale(newLocale);
        const current = preferencesManager.getPreferences().app.locale;
        if (current !== normalized) {
          preferencesManager.updatePreferences({
            app: { locale: normalized },
          });
        }
      },
    );
  }
});
