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
        /** HSL raw 格式: "H S% L%"，直接映射到 CSS --primary 变量 */
        colorPrimary: '142.1 76.2% 36.3%',
        colorSuccess: '142.1 76.2% 36.3%',
        colorWarning: '38 92% 50%',
        colorDestructive: '0 84.2% 60.2%',
        radius: '0.6rem',
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
