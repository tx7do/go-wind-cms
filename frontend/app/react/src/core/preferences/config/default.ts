import type {Preferences} from '../types';

const defaultPreferences: Preferences = {
    app: {
        name: 'GoWind CMS',
        title: 'GoWind Content Hub',
        defaultAvatar: '/default-avatar.png',
        locale: 'zh-CN',
        isMobile: false,
        compact: false,
        defaultPageSize: 10,
    },
    theme: {
        mode: 'auto',
        colorPrimary: 'hsl(212 100% 45%)',
        colorSuccess: 'hsl(144 57% 58%)',
        colorWarning: 'hsl(42 84% 61%)',
        colorDestructive: 'hsl(348 100% 61%)',
        radius: '0.5',
    },
    content: {
        hideSensitiveContent: true,
        compactMode: false,
        showRecommendations: true,
    },
    copyright: {
        enable: true,
        companyName: 'GoWind',
        companySiteLink: 'https://www.gowind.cloud',
        date: '2026',
        icp: '',
        icpLink: '',
    },
    widget: {
        themeToggle: true,
        languageToggle: true,
        globalSearch: true,
        backToTop: true,
    },
    transition: {
        enable: true,
        loading: true,
        name: 'fade-slide',
        progress: true,
    },
};

export {defaultPreferences};
