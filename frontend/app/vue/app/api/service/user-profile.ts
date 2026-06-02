import {
  createUserProfileServiceClient,
} from '@/api/generated/app/service/v1';
import { requestApi } from '@/core/transport/rest';
import { makeUpdateMask, omit } from '@/core/transport/rest';

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
 * 获取当前用户
 */
export async function getMe() {
  try {
    return await getUserProfileService().GetUser({});
  } catch (error) {
    console.error('getMe failed:', error);
    return null;
  }
}

/**
 * 更新当前用户
 */
export async function updateUser(id: number, values: Record<string, any> = {}) {
  const password = values.password ?? null;
  const cleaned = omit(values, 'password');
  const updateMask = makeUpdateMask(Object.keys(cleaned ?? []));
  return await getUserProfileService().UpdateUser({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...cleaned,
    },
    password,
    // @ts-ignore proto generated code is error.
    updateMask,
  });
}

/**
 * 修改用户密码
 */
export async function changePassword(oldPassword: string, newPassword: string) {
  return await getUserProfileService().ChangePassword({
    oldPassword,
    newPassword,
  });
}

/**
 * 上传用户头像（Base64）
 */
export async function uploadAvatarBase64(imageBase64: string) {
  return await getUserProfileService().UploadAvatar({
    imageBase64,
  });
}

/**
 * 上传用户头像（图片URL）
 */
export async function uploadAvatarUrl(imageUrl: string) {
  return await getUserProfileService().UploadAvatar({
    imageUrl,
  });
}

/**
 * 删除用户头像
 */
export async function deleteAvatar() {
  return await getUserProfileService().DeleteAvatar({});
}

/**
 * 绑定手机号
 */
export async function bindPhone(phone: string, code: string) {
  return await getUserProfileService().BindContact({
    phone: { phone, code },
  });
}

/**
 * 绑定邮箱
 */
export async function bindEmail(email: string, verificationCode: string) {
  return await getUserProfileService().BindContact({
    email: { email, verificationCode },
  });
}

/**
 * 验证手机号
 */
export async function verifyPhone(phone: string, code: string, verificationId?: string) {
  return await getUserProfileService().VerifyContact({
    phone: { phone, code },
    verificationId,
  });
}

/**
 * 验证邮箱
 */
export async function verifyEmail(email: string, code: string, verificationId?: string) {
  return await getUserProfileService().VerifyContact({
    email: { email, code },
    verificationId,
  });
}
