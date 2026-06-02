'use client';

export {default as NavigationProgress} from './NavigationProgress';

/**
 * @deprecated 已被 NavigationProgress 顶部进度条替代
 *
 * 全屏遮罩会导致页面无法交互，体验较差。
 * 请使用 <NavigationProgress /> 或 loading.tsx 骨架屏方案。
 */
export default function GlobalLoading() {
    return null;
}
