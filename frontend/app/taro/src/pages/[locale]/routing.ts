import {DEFAULT_LANGUAGE, locales} from "@/i18n";

export const routing = {
    // 支持的语言列表
    locales: locales,

    // 默认语言
    defaultLocale: DEFAULT_LANGUAGE,
};

export type SupportedLocale = (typeof routing.locales)[number];

export function getSupportedLocale(locale: string): SupportedLocale | undefined {
    return routing.locales.find(l => l === locale);
}

export function isSupportedLocale(locale: string): locale is SupportedLocale {
    return getSupportedLocale(locale) !== undefined;
}
