import Taro from '@tarojs/taro';

/**
 * 跨端 Storage 适配器
 *
 * 背景：原 storage.class.ts 和 zustand persist 直接使用 window.localStorage /
 * window.sessionStorage，在微信小程序等无 window 对象的运行时会抛
 * "window is not defined" 导致 token 持久化失效甚至崩溃。
 *
 * 本适配器实现了浏览器 Storage 接口（getItem/setItem/removeItem/key/length/clear），
 * 内部按运行环境委托：
 *   - 浏览器/H5：直接转发到 window.localStorage / window.sessionStorage
 *   - 小程序：转发到 Taro 的同步存储 API（getStorageInfoSync 提供 keys 列表以支持 length/key）
 *
 * 注意：小程序存储的单条 value 上限与总体上限由平台决定，调用方应自行控制体积。
 */

/** 判断当前是否运行在浏览器/H5 环境（存在 window.localStorage） */
function isBrowserEnv(): boolean {
  return typeof window !== 'undefined' && typeof window.localStorage !== 'undefined';
}

/**
 * 判断当前环境是否有可用的存储（浏览器 localStorage 或 Taro 小程序存储）。
 * 用于 SSR / node 环境（如单元测试）的兜底判断。
 */
export function isStorageAvailable(): boolean {
  if (isBrowserEnv()) return true;
  // 小程序环境：Taro 在编译期注入，运行期 Taro.getStorageSync 可用即视为有存储
  try {
    return typeof Taro !== 'undefined' && typeof Taro.getStorageSync === 'function';
  } catch {
    return false;
  }
}

/**
 * 小程序端的 Storage 适配实现。
 * 用一个内部索引 key 维护所有 key 的有序列表，以模拟 Web Storage 的 length/key 接口。
 */
class MiniProgramStorage implements Storage {
  /** 索引 key：保存所有业务 key 的有序数组 */
  private static readonly INDEX_KEY = '__taro_storage_index__';

  private readIndex(): string[] {
    try {
      const keys = Taro.getStorageSync(MiniProgramStorage.INDEX_KEY);
      return Array.isArray(keys) ? keys : [];
    } catch {
      return [];
    }
  }

  private writeIndex(keys: string[]): void {
    try {
      Taro.setStorageSync(MiniProgramStorage.INDEX_KEY, keys);
    } catch (error) {
      console.warn('[MiniProgramStorage] writeIndex failed:', error);
    }
  }

  get length(): number {
    // 优先用平台提供的 storage 信息（更准确，不依赖自身索引）
    try {
      const info = Taro.getStorageInfoSync();
      return info.keys.filter((k) => k !== MiniProgramStorage.INDEX_KEY).length;
    } catch {
      return this.readIndex().length;
    }
  }

  clear(): void {
    try {
      Taro.clearStorageSync();
    } catch (error) {
      console.warn('[MiniProgramStorage] clear failed:', error);
    }
  }

  getItem(key: string): string | null {
    try {
      const val = Taro.getStorageSync(key);
      if (val === '' || val === undefined || val === null) {
        // 无法区分"空字符串"和"不存在"，小程序存储里 getStorageSync 对不存在的 key 返回 ''
        // 这里通过 storageInfo 的 keys 列表二次确认
        const info = Taro.getStorageInfoSync();
        if (!info.keys.includes(key)) return null;
        return val === undefined || val === null ? '' : String(val);
      }
      return typeof val === 'string' ? val : JSON.stringify(val);
    } catch {
      return null;
    }
  }

  key(index: number): string | null {
    try {
      const info = Taro.getStorageInfoSync();
      const keys = info.keys.filter((k) => k !== MiniProgramStorage.INDEX_KEY);
      return index >= 0 && index < keys.length ? keys[index] : null;
    } catch {
      return null;
    }
  }

  removeItem(key: string): void {
    try {
      Taro.removeStorageSync(key);
    } catch (error) {
      console.warn('[MiniProgramStorage] removeItem failed:', error);
    }
  }

  setItem(key: string, value: string): void {
    try {
      Taro.setStorageSync(key, value);
      // 维护索引（保证 key()/length 可用，且不重复）
      const keys = this.readIndex();
      if (!keys.includes(key)) {
        keys.push(key);
        this.writeIndex(keys);
      }
    } catch (error) {
      console.warn('[MiniProgramStorage] setItem failed:', error);
      throw error;
    }
  }
}

/**
 * 根据运行环境返回合适的本地存储实现（localStorage 语义）。
 * - 浏览器：window.localStorage
 * - 小程序：基于 Taro 同步 API 的适配器
 */
export function getLocalStorage(): Storage {
  if (isBrowserEnv()) {
    return window.localStorage;
  }
  return new MiniProgramStorage();
}

/**
 * 根据运行环境返回合适的会话存储实现（sessionStorage 语义）。
 * 小程序没有 sessionStorage 概念，退化为本地存储（与 localStorage 一致）。
 */
export function getSessionStorage(): Storage {
  if (isBrowserEnv() && typeof window.sessionStorage !== 'undefined') {
    return window.sessionStorage;
  }
  return new MiniProgramStorage();
}
