import type {DeepPartial} from "../types";
import {isValidElement} from "react";

/**
 * 深度合并对象
 */
export function mergeDeep<T extends object>(target: T, source: DeepPartial<T>): T {
    const result: any = {...target};

    for (const key in source) {
        if (Object.prototype.hasOwnProperty.call(source, key)) {
            const sourceValue = source[key as keyof DeepPartial<T>];
            const targetValue = target[key as keyof T];

            // 跳过 React 元素、DOM 节点、函数等非纯对象
            if (
                isValidElement(sourceValue) || 
                sourceValue instanceof Element || 
                typeof sourceValue === 'function' ||
                sourceValue === null
            ) {
                console.warn(`[mergeDeep] 跳过非纯对象字段: ${String(key)}, type:`, typeof sourceValue);
                continue;
            }

            // 如果两个值都是普通对象，则递归合并
            if (
                isPlainObject(targetValue) &&
                isPlainObject(sourceValue)
            ) {
                result[key] = mergeDeep(
                    targetValue as Record<string, any>,
                    sourceValue as DeepPartial<any>
                );
            }
            // 否则直接使用源值覆盖（包括 undefined 的情况）
            else if (sourceValue !== undefined) {
                result[key] = sourceValue;
            }
        }
    }

    return result;
}

/**
 * 判断是否为普通对象
 */
export function isPlainObject(value: any): value is Record<string, any> {
    return (
        value !== null &&
        typeof value === 'object' &&
        Object.prototype.toString.call(value) === '[object Object]' &&
        !Array.isArray(value)
    );
}
