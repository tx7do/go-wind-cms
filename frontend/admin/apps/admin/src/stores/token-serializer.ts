/**
 * AUD9-M5: 构造并注入 access store 的 token 加密序列化器。
 *
 * 落盘前对 accessToken/refreshToken 字段做 AES-CBC 加密（key 复用 VITE_AES_KEY，
 * 与密码加密同源），读盘后解密。其余字段（accessCodes/expireTime）明文保留。
 * 解密失败置 null，迫使重新登录。
 *
 * 本模块运行在 admin app（已声明 crypto-js 依赖），通过 @vben/stores 暴露的
 * setTokenSerializer 注入，使 @vben/stores 共享包无需直接依赖任何加密库。
 */
import CryptoJS from 'crypto-js';

import { setTokenSerializer, type TokenSerializer } from '@vben/stores';

const AES_KEY = import.meta.env.VITE_AES_KEY ?? '';

function encField(plain: string): string {
  return CryptoJS.AES.encrypt(plain, CryptoJS.enc.Utf8.parse(AES_KEY), {
    iv: CryptoJS.enc.Utf8.parse(AES_KEY),
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  }).toString();
}

function decField(cipher: string): string {
  return CryptoJS.AES.decrypt(cipher, CryptoJS.enc.Utf8.parse(AES_KEY), {
    iv: CryptoJS.enc.Utf8.parse(AES_KEY),
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  }).toString(CryptoJS.enc.Utf8);
}

/**
 * 注入加密序列化器。必须在 access store 首次实例化前调用。
 * 无 AES key 时跳过注入（保留明文默认序列化器）。
 */
export function setupTokenSerializer(): void {
  if (!AES_KEY) return;
  const serializer: TokenSerializer = {
    serialize(data: any): string {
      const clone = JSON.parse(JSON.stringify(data));
      if (typeof clone?.accessToken === 'string' && clone.accessToken) {
        clone.accessToken = encField(clone.accessToken);
      }
      if (typeof clone?.refreshToken === 'string' && clone.refreshToken) {
        clone.refreshToken = encField(clone.refreshToken);
      }
      return JSON.stringify(clone);
    },
    deserialize(data: string): any {
      const parsed = JSON.parse(data);
      if (typeof parsed?.accessToken === 'string' && parsed.accessToken) {
        try {
          parsed.accessToken = decField(parsed.accessToken);
        } catch {
          parsed.accessToken = null;
        }
      }
      if (typeof parsed?.refreshToken === 'string' && parsed.refreshToken) {
        try {
          parsed.refreshToken = decField(parsed.refreshToken);
        } catch {
          parsed.refreshToken = null;
        }
      }
      return parsed;
    },
  };
  setTokenSerializer(serializer);
}
