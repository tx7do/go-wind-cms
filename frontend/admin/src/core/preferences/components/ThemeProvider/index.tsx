import React, { useEffect, useLayoutEffect, useMemo, useState } from 'react';

import { ConfigProvider, Watermark, theme as antdTheme, type ThemeConfig } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import enUS from 'antd/locale/en_US';
// import jaJP from 'antd/locale/ja_JP';

import { usePreferencesStore } from '../../store';
import { DARK_PALETTE, darkThemeComponents, darkThemeTokens } from '../../config';
import { SECONDARY_TEXT_DARK, SECONDARY_TEXT_LIGHT } from '../../config/constants';
import { mixColor } from '@/utils/color';
import type { SupportedLanguagesType } from '../../types';

export interface ThemeProviderProps {
  children: React.ReactNode;
}

const localeMap: Record<SupportedLanguagesType, any> = {
  'zh-CN': zhCN,
  'en-US': enUS,
};

/**
 * 全局主题提供者
 * 负责：AntD 主题动态切换、系统深色模式同步、CSS 滤镜、水印、国际化
 */
export const ThemeProvider = ({ children }: ThemeProviderProps) => {
  // 1. 精确订阅偏好设置（分别订阅 theme 和 app 模块）
  const themePrefs = usePreferencesStore((state) => state.preferences.theme);
  const appPrefs = usePreferencesStore((state) => state.preferences.app);

  // 2. 系统主题监听（仅当 mode === 'auto' 时生效）
  const [systemIsDark, setSystemIsDark] = useState(false);
  useEffect(() => {
    if (themePrefs.mode !== 'auto') return;

    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    setSystemIsDark(mediaQuery.matches);

    const handleChange = (e: MediaQueryListEvent) => setSystemIsDark(e.matches);
    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, [themePrefs.mode]);

  // 计算最终生效的主题模式
  const effectiveMode = useMemo(() => {
    return themePrefs.mode === 'auto' ? (systemIsDark ? 'dark' : 'light') : themePrefs.mode;
  }, [themePrefs.mode, systemIsDark]);

  // 3. 生成 Ant Design 主题配置
  const themeConfig = useMemo((): ThemeConfig => {
    const algorithms = [];
    if (effectiveMode === 'dark') algorithms.push(antdTheme.darkAlgorithm);
    if (appPrefs.compact) algorithms.push(antdTheme.compactAlgorithm);

    const tokens: ThemeConfig['token'] = {
      colorPrimary: themePrefs.colorPrimary,
      colorSuccess: themePrefs.colorSuccess,
      colorWarning: themePrefs.colorWarning,
      colorError: themePrefs.colorDestructive,
      // §2.1 声明"react 端 colorInfo 同值"，但代码里从未设过 colorInfo/colorLink，
      // 于是链接与 info 态走 antd 默认 #1677FF（浅色实测：链接文字 #1677FF on 白 4.10:1，
      // 且与主色 #006BE6 不同蓝）。显式钉住主色后链接 = #006BE6 on 白 4.76:1。
      colorInfo: themePrefs.colorPrimary,
      // 暗色下不能直接用主色当链接色：darkAlgorithm 把 #006BE6 再派生成 #035EC7，
      // 实测"收起/文档/公告"等链接文字 2.90:1（暗色大底 #111827 上）；按 §2.1 文字档
      // 反解（主色 70% + 白 30%）= #4D97EE，同底 5.89:1。
      colorLink:
        effectiveMode === 'dark'
          ? mixColor(themePrefs.colorPrimary, '#ffffff', 0.7)
          : themePrefs.colorPrimary,
      borderRadius: Number.parseInt(themePrefs.radius) || 6,
      // colorTextDescription 驱动 .ant-typography-secondary / 面包屑分隔符 /
      // Descriptions 标签 / Empty 描述——antd 默认 rgba(0,0,0,.45) 实测只有 3.35:1（浅色）
      // 且暗色派生到 #6E7681 只有 3.86:1；两侧都按 §2.2 钉成"次要文字"档。
      colorTextDescription: effectiveMode === 'dark' ? SECONDARY_TEXT_DARK : SECONDARY_TEXT_LIGHT,
    };

    // 暗黑模式叠加蓝调深色色阶（背景/文字/边框/填充分层），亮色保持 antd 默认
    if (effectiveMode === 'dark') Object.assign(tokens, darkThemeTokens);

    // 暗黑模式下输入类控件用 L1 内嵌底色，在卡片上形成「凹陷」层次
    const components: ThemeConfig['components'] =
      effectiveMode === 'dark' ? darkThemeComponents : undefined;

    return {
      algorithm: algorithms.length > 0 ? algorithms : antdTheme.defaultAlgorithm,
      token: tokens,
      components,
    };
  }, [effectiveMode, appPrefs.compact, themePrefs]);

  // 4. 同步根背景色和文本颜色（解决暗黑模式白色闪烁）
  // 用 useLayoutEffect：主题切换走 startThemeViewTransition 时，flushSync 提交
  // 阶段就会同步执行 layout effect，保证 data-theme/根背景色在过渡新快照
  // 截取前落地（View Transition 的 update 阶段不跑 rAF，不能靠等帧）。
  useLayoutEffect(() => {
    const root = document.documentElement;
    const body = document.body;

    const bgLayout = effectiveMode === 'dark' ? DARK_PALETTE.bgLayout : '#f5f5f5';
    root.style.backgroundColor = bgLayout;
    body.style.backgroundColor = bgLayout;

    const colorText = effectiveMode === 'dark'
      ? 'rgba(230, 237, 243, 0.92)'
      : 'rgba(0, 0, 0, 0.85)';
    root.style.color = colorText;

    // 标记主题模式，供 CSS 选择器使用
    root.setAttribute('data-theme', effectiveMode);

    return () => {
      root.style.backgroundColor = '';
      body.style.backgroundColor = '';
      root.style.color = '';
      root.removeAttribute('data-theme');
    };
  }, [effectiveMode]);

  // 4.1 主色/语义色锚点变量：错误页插画等自定义 CSS 的取色源，同时是
  // styles/semantic-text.css 里「语义色作文字」降档的输入（换主色后文字档跟着联动）。
  // 写在 <html> 上——antd v6 的 --ant-* cssVar 只注入在包裹层内，<html> 作用域引用会落空。
  useEffect(() => {
    const root = document.documentElement;
    root.style.setProperty('--app-color-primary', themePrefs.colorPrimary);
    root.style.setProperty('--app-color-success', themePrefs.colorSuccess);
    root.style.setProperty('--app-color-warning', themePrefs.colorWarning);
    root.style.setProperty('--app-color-destructive', themePrefs.colorDestructive);
    return () => {
      root.style.removeProperty('--app-color-primary');
      root.style.removeProperty('--app-color-success');
      root.style.removeProperty('--app-color-warning');
      root.style.removeProperty('--app-color-destructive');
    };
  }, [themePrefs.colorPrimary, themePrefs.colorSuccess, themePrefs.colorWarning, themePrefs.colorDestructive]);

  // 5. 应用 CSS 滤镜（色弱 / 灰色模式）
  useEffect(() => {
    const root = document.documentElement;
    if (appPrefs.colorWeakMode) {
      root.style.filter = 'invert(80%)';
    } else if (appPrefs.colorGrayMode) {
      root.style.filter = 'grayscale(100%)';
    } else {
      root.style.filter = '';
    }
  }, [appPrefs.colorWeakMode, appPrefs.colorGrayMode]);

  // 6. 国际化配置
  const locale = localeMap[appPrefs.locale] ?? zhCN;

  // 7. 水印内容
  const watermarkText = appPrefs.watermark ? `${appPrefs.name} ${appPrefs.version}` : '';

  // 8. 水印颜色（跟随主题）
  const watermarkColor = useMemo(() => {
    return effectiveMode === 'dark'
      ? 'rgba(255, 255, 255, 0.1)'
      : 'rgba(0, 0, 0, 0.15)';
  }, [effectiveMode]);

  return (
    <ConfigProvider
      theme={themeConfig}
      locale={locale}
      wave={{ disabled: true }} // 关闭点击波纹，提升 Admin 操作手感
    >
      {watermarkText ? (
        <Watermark
          content={watermarkText}
          font={{ fontSize: 14, color: watermarkColor }}
          gap={[120, 120]}
          offset={[60, 60]}
          rotate={-22}
          zIndex={1}
        >
          {children}
        </Watermark>
      ) : (
        children
      )}
    </ConfigProvider>
  );
};
