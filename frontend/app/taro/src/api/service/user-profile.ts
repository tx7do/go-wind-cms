import {
  createUserProfileServiceClient,
  type identityservicev1_User,
} from '@/api/generated/app/service/v1';
import { requestApi } from '@/core';

let _instance: ReturnType<typeof createUserProfileServiceClient> | null = null;

/**
 * 获取用户资料服务单例（延迟初始化）
 */
export function getUserProfileService() {
  if (!_instance) {
    _instance = createUserProfileServiceClient(requestApi);
  }
  return _instance;
}

/**
 * 获取当前登录用户信息
 */
export async function getMe(): Promise<identityservicev1_User | null> {
  try {
    return await getUserProfileService().GetUser({});
  } catch (error) {
    console.error('getMe failed:', error);
    return null;
  }
}
