import { useTranslation, type UseTranslationOptions } from 'react-i18next';
import { type SupportedLocale } from '@/locales';

// 类型安全 t 函数
export type TFunction = (key: string, options?: { [key: string]: any }) => string;

// 封装 useTranslation，自动推导命名空间
// 支持核心命名空间和扩展命名空间
export const useI18n = <N extends string = 'common'>(
  ns?: N | readonly N[],
  options?: UseTranslationOptions<any>,
) => {
  // 统一转成数组，默认使用 'common' 命名空间
  const namespace = ns ? (Array.isArray(ns) ? ns : [ns]) : ['common'];

  const { t, i18n, ...rest } = useTranslation(namespace as any, options);

  // 类型安全的 t 函数
  const typedT = ((key: string, options?: any) => {
    return t(key, options);
  }) as TFunction;

  // 切换语言
  const changeLocale = async (locale: SupportedLocale) => {
    await i18n.changeLanguage(locale);
  };

  return {
    t: typedT,
    i18n,
    changeLocale,
    ...rest,
  };
};
