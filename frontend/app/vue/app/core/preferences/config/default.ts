import type { Preferences } from "../types";

const defaultPreferences: Preferences = {
  app: {
    name: "GoWind CMS",
    version: "0.0.0",
    locale: "zh-CN",
    isMobile: false,
    defaultAvatar: "/default-avatar.png",
    dynamicTitle: true,
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
    builtinType: "default",
    colorPrimary: "hsl(220 100% 55%)",
    mode: "auto",
    radius: "0.5",
  },
  transition: {
    enable: true,
    loading: true,
    name: "fade-slide",
    progress: true,
  },
};

export { defaultPreferences };
