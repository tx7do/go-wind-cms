export interface ApplicationConfig {
  apiURL: string;
}

export interface ApplicationConfigRaw {
  VITE_GLOB_API_URL: string;
}

/**
 * 由 vite-inject-app-config 注入的全局配置
 */
export function useAppConfig(
  env: Record<string, any>,
  _isProduction: boolean,
): ApplicationConfig {
  const config = (env as ApplicationConfigRaw);

  const { VITE_GLOB_API_URL } = config;

  return {
    apiURL: VITE_GLOB_API_URL,
  };
}
