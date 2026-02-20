/**
 * 深度克隆一个值，支持对象、数组、Map、Set、Date、RegExp 等类型，且能处理循环引用。
 * @param value
 */
export function deepClone<T>(value: T): T {
  // 尝试使用原生 structuredClone，但发生异常时回退到手写实现
  if (typeof (globalThis as any).structuredClone === 'function') {
    try {
      return (globalThis as any).structuredClone(value);
    } catch {
      // structuredClone 不能克隆某些类型（如 Vue 响应式代理、DOM 节点、window 等）
      // 回退到手写深拷贝
    }
  }

  const seen = new WeakMap<any, any>();

  const _clone = (v: any): any => {
    if (v === null || typeof v !== 'object') return v;
    if (v instanceof Date) return new Date(v.getTime());
    if (v instanceof RegExp) return new RegExp(v.source, v.flags);
    if (v instanceof Map) {
      if (seen.has(v)) return seen.get(v);
      const m = new Map();
      seen.set(v, m);
      for (const [k, val] of v) m.set(_clone(k), _clone(val));
      return m;
    }
    if (v instanceof Set) {
      if (seen.has(v)) return seen.get(v);
      const s = new Set();
      seen.set(v, s);
      for (const item of v) s.add(_clone(item));
      return s;
    }
    if (seen.has(v)) return seen.get(v);

    if (Array.isArray(v)) {
      const arr: any[] = [];
      seen.set(v, arr);
      for (let i = 0; i < v.length; i++) arr[i] = _clone(v[i]);
      return arr;
    }

    const obj: any = Object.create(Object.getPrototypeOf(v));
    seen.set(v, obj);
    for (const key of Reflect.ownKeys(v)) {
      obj[key as any] = _clone((v as any)[key as any]);
    }
    return obj;
  };

  return _clone(value);
}
