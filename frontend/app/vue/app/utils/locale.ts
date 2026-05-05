import {computed} from 'vue';

type LocaleOption = {
    code: string;
    iso?: string;
    name?: string;
};

function useCurrentLocaleOption() {
    const {locale, locales} = useI18n();

    return computed(() =>
        (locales.value as LocaleOption[]).find((item) => item.code === locale.value),
    );
}

export function useLocaleIso() {
    const currentLocale = useCurrentLocaleOption();
    return computed(() => currentLocale.value?.iso);
}

export function useLocaleName() {
    const currentLocale = useCurrentLocaleOption();
    return computed(() => currentLocale.value?.name);
}
