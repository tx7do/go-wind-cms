import { flushSync } from 'react-dom';

export interface ThemeTransitionOrigin {
  x: number;
  y: number;
}

export interface ThemeTransitionOptions {
  /** 执行状态更新的函数，内部会以 flushSync 包裹以保证过渡快照时机 */
  update: () => void;
  /** 动效圆心（通常是点击点）；缺省时取视口中心 */
  origin?: ThemeTransitionOrigin;
}

/**
 * 主题切换的 View Transitions 圆形扩散动效（ele/ThemeSwitch 风格）：
 * 始终让新画面自圆心扩张覆盖旧画面。切暗时暗色从点击点立即绽开，
 * 不存在"旧画面收缩"方案里点击点附近久久残留亮色的问题。
 * 浏览器不支持或用户偏好减少动效（prefers-reduced-motion）时退化为直接切换。
 */
export function startThemeViewTransition({
  update,
  origin,
}: ThemeTransitionOptions): void {
  if (!document.startViewTransition || window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    update();
    return;
  }

  const x = origin?.x ?? window.innerWidth / 2;
  const y = origin?.y ?? window.innerHeight / 2;
  const endRadius = Math.hypot(
    Math.max(x, window.innerWidth - x),
    Math.max(y, window.innerHeight - y),
  );
  const clipPath = [
    `circle(0px at ${x}px ${y}px)`,
    `circle(${endRadius}px at ${x}px ${y}px)`,
  ];

  const transition = document.startViewTransition(() => {
    // React 状态更新走异步批处理，flushSync 保证回调返回前 DOM 与
    // useLayoutEffect 副作用（data-theme/根背景色）都已落地，否则过渡
    // 新快照会拍到旧画面。update 阶段渲染循环暂停（rAF 不触发），
    // 回调必须同步完成，不能等帧。
    flushSync(update);
  });

  transition.ready
    .then(() => {
      document.documentElement.animate(
        {
          clipPath,
        },
        {
          duration: 450,
          easing: 'ease-out',
          // 必须 forwards：不填充的话动画一结束 clip 立即回退为无裁剪，
          // 而过渡伪元素树要到 finished 后才拆除，中间多出的帧会裸绘
          // 底层旧快照造成闪屏。新画面本就在顶层，此处兜底语义不变。
          fill: 'forwards',
          pseudoElement: '::view-transition-new(root)',
        },
      );
    })
    .catch(() => {
      // best-effort：ready 失败（如过渡被浏览器跳过）只影响动效，主题切换本身已生效
    });
  transition.finished.catch(() => {
    // best-effort：finished 的拒绝仅代表过渡异常中止，无可恢复动作
  });
}
