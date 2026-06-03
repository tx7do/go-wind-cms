import type { Preferences } from "../types";

const defaultPreferences: Preferences = {
  app: {
    name: "GoWind CMS",
    title: "GoWind Content Hub",
    version: "0.0.0",
    locale: "zh-CN",
    isMobile: false,
    defaultAvatar: "/default-avatar.png",
    dynamicTitle: true,
    defaultPageSize: 10,
    compact: false,
  },
  copyright: {
    enable: true,
    companyName: "GoWind",
    companySiteLink: "https://www.gowind.cloud",
    date: "2026",
    icp: "",
    icpLink: "",
  },
  logo: {
    enable: true,
    source: "/logo.png",
  },
  theme: {
    mode: "auto",
    colorPrimary: "142.1 76.2% 36.3%",
    colorSuccess: "142.1 76.2% 36.3%",
    colorWarning: "38 92% 50%",
    colorDestructive: "0 84.2% 60.2%",
    radius: "0.6rem",
  },
  content: {
    hideSensitiveContent: true,
    compactMode: false,
    showRecommendations: true,
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
    name: "fade-slide",
    progress: true,
  },
};

export { defaultPreferences };
