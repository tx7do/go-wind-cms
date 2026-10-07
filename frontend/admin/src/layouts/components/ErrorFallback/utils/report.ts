interface ErrorReportPayload {
  message: string;
  stack?: string;
  componentStack?: string;
  url: string;
  userAgent: string;
  timestamp: number;
}

/**
 * 上报错误到监控系统（示例：Sentry/自研平台）
 */
export const reportError = (error: Error, componentStack?: string) => {
  if (import.meta.env.PROD) {
    const payload: ErrorReportPayload = {
      message: error.message,
      stack: error.stack,
      componentStack,
      url: window.location.href,
      userAgent: navigator.userAgent,
      timestamp: Date.now(),
    };

    // 示例：使用 fetch 上报（实际建议用 @sentry/react 或 logrocket）
    // fetch('/api/monitor/error', {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify(payload),
    // }).catch(() => {});

    console.debug('[ErrorReport] Payload:', payload);
  }
};
