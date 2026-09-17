import { useRouteError } from 'react-router-dom';
import { ErrorFallback } from './index';

/**
 * React Router v6 errorElement 专用适配器
 * 自动规范化错误对象并透传给 ErrorFallback
 */
export const RouteErrorFallback = () => {
  const routeError = useRouteError();

  // 规范化错误对象（React Router 可能返回 string/Response/Error）
  const normalizedError = routeError instanceof Error ? routeError : new Error(String(routeError));

  return <ErrorFallback error={normalizedError} />;
};

export default RouteErrorFallback;
