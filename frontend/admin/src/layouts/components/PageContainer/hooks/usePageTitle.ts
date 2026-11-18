import React, { useEffect, useMemo, useCallback } from 'react';
import { usePreferencesStore } from '@/core/preferences/store';
import { useI18n } from '@/core/i18n';
import { useTranslation } from 'react-i18next';

interface UsePageTitleOptions {
  /** 手动传入的标题（优先级最高） */
  manual?: React.ReactNode;

  /** 路由对象中的标题 */
  routeTitle?: string;

  /** 默认标题（兜底） */
  defaultTitle?: string;

  /** 是否更新浏览器标签页标题 */
  updateDocumentTitle?: boolean;
}

/** 路由标题翻译：menu:/routes: 前缀与裸 key 统一走 routes 命名空间
 *  （带 defaultValue 兜底）。导出共享给全局 document.title 更新器
 *  （src/router/useDocumentTitle.ts），避免翻译规则双份漂移。 */
export const translateRouteTitle = (
  title: string,
  tMenu: (key: string, options?: { defaultValue?: string }) => string,
): string => {
  // 如果 title 是 i18n key (以 'menu:' 或 'routes:' 开头)，进行翻译
  if (title.startsWith('menu:')) {
    const keyName = title.substring(5);
    return tMenu(keyName, { defaultValue: title });
  } else if (title.startsWith('routes:')) {
    const keyName = title.substring(7);
    return tMenu(keyName, { defaultValue: title });
  } else if (title.startsWith('menu.')) {
    const keyName = title.substring(5);
    return tMenu(keyName, { defaultValue: title });
  } else if (title.startsWith('routes.')) {
    const keyName = title.substring(7);
    return tMenu(keyName, { defaultValue: title });
  } else {
    // 否则直接翻译
    return tMenu(title, { defaultValue: title });
  }
};

export const usePageTitle = ({
  manual,
  routeTitle,
  defaultTitle,
  updateDocumentTitle = true,
}: UsePageTitleOptions): React.ReactNode => {
  const { app } = usePreferencesStore((state) => state.preferences);
  const { t } = useI18n('common');
  const { t: tMenu } = useTranslation('routes');
  const resolvedDefaultTitle = defaultTitle ?? t('pageContainer.defaultTitle');

  // 翻译标题的工具函数
  const translateTitle = useCallback(
    (title: string): string => translateRouteTitle(title, tMenu),
    [tMenu],
  );

  // 计算最终标题（包含翻译）
  const title = useMemo(() => {
    if (manual) return manual;
    if (routeTitle) {
      // 如果 routeTitle 是 i18n key，进行翻译
      return translateTitle(routeTitle);
    }
    return resolvedDefaultTitle;
  }, [manual, routeTitle, resolvedDefaultTitle, translateTitle]);

  // 更新 document.title（配合 preferences.app.dynamicTitle）
  useEffect(() => {
    if (!updateDocumentTitle || !app.dynamicTitle) return;

    const originalTitle = document.title;
    const appName = app.name || t('app.defaultName');

    if (title && typeof title === 'string') {
      document.title = `${title} - ${appName}`;
    }

    return () => {
      document.title = originalTitle;
    };
  }, [title, app.name, app.dynamicTitle, updateDocumentTitle, t]);

  return title;
};
