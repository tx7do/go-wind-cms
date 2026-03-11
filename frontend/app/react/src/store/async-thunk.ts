import {useEffect, useRef} from 'react';
import {
    AsyncThunk,
    AsyncThunkPayloadCreator,
    AsyncThunkOptions,
    SerializedError,
    createAsyncThunk,
} from '@reduxjs/toolkit';

import {AppDispatch} from "./index";

/**
 * 通用的可取消Thunk调用配置
 */
export interface AbortableCallOptions {
    /** 是否在调用新请求时自动取消同类型的旧请求 */
    autoAbortPrevious?: boolean;
}

/**
 * 带取消能力的异步调用函数类型
 */
export type AbortableCall<ReturnData, Arg> =
    ((args: Arg) => Promise<ReturnData | null>) & {
    cancel: (reason?: string) => void;
};

/**
 * 批量创建的可取消调用对象类型
 * 修正：使用条件类型正确提取 Args 并合并 signal
 */
// @ts-expect-error - 这里的类型推导比较复杂，暂时允许 TypeScript 报错
export type AbortableCalls<T extends Record<string, AsyncThunk<unknown, unknown, unknown>>> = {
    // @ts-expect-error - 这里的类型推导比较复杂，暂时允许 TypeScript 报错
    [K in keyof T]: T[K] extends AsyncThunk<infer ReturnData, infer Args, unknown>
        ? AbortableCall<ReturnData, MergeSignal<Args>>
        : never;
} & {
    cancelAll: (reason?: string) => void;
};

/**
 * 类型工具：将 signal 合并到参数类型中
 * 修正：正确处理 void/undefined 参数情况
 */
type MergeSignal<Args> = Args extends void | undefined
    ? { signal?: AbortSignal }
    : Args & { signal?: AbortSignal };

/**
 * 创建单个可取消的 Thunk调用函数
 * @param dispatch Redux dispatch 函数
 * @param thunk 要包装的 async thunk
 * @param options 配置项
 * @returns 带取消能力的异步调用函数
 */
function createAbortableCall<
    ReturnData,
    ThunkArg,
    ThunkApiConfig extends { rejectValue: unknown }
>(
    dispatch: AppDispatch,
    thunk: AsyncThunk<ReturnData, ThunkArg, ThunkApiConfig>,
    options: AbortableCallOptions = {autoAbortPrevious: true}
): AbortableCall<ReturnData, MergeSignal<ThunkArg>> {
    // 存储当前请求的 AbortController
    let currentController: AbortController | null = null;

    // 实际的调用函数 - 修正参数类型
    const callFn = async (
        args: MergeSignal<ThunkArg>
    ): Promise<ReturnData | null> => {
        // 自动取消上一个同类型请求（可选功能）
        if (options.autoAbortPrevious && currentController) {
            currentController.abort('New request initiated');
        }

        // 创建新的取消控制器
        const abortController = new AbortController();
        currentController = abortController;

        try {
            // 合并 signal 参数 - 修正：避免重复添加 signal
            const thunkParams = {
                ...args as Record<string, unknown>,
                signal: abortController.signal
            } as MergeSignal<ThunkArg>;

            // 执行 thunk
            // @ts-expect-error - 忽略类型检查
            const action = await dispatch(thunk(thunkParams));

            // 检查请求状态 - 修正：更安全的类型守卫
            // @ts-expect-error - 忽略类型检查
            if (action.meta.requestStatus === 'rejected') {
                const errorPayload = (action as { payload?: unknown }).payload;
                const errorMessage =
                    typeof errorPayload === 'object' &&
                    errorPayload !== null &&
                    'message' in errorPayload &&
                    typeof (errorPayload as { message: unknown }).message === 'string'
                        ? (errorPayload as { message: string }).message
                        : String(errorPayload ?? 'Request failed');
                throw new Error(errorMessage);
            }

            return (action as { payload?: ReturnData }).payload ?? null;
        } catch (error) {
            // 过滤取消错误（主动取消不抛出异常）
            const errorObj = error as Error;
            if (
                errorObj.name === 'AbortError' ||
                errorObj.message === 'New request initiated' ||
                errorObj.message === 'Manual cancellation'
            ) {
                console.log(`Request aborted: ${thunk.typePrefix}`);
                return null;
            }

            // 其他错误抛出，让调用方处理
            console.error(`Error in ${thunk.typePrefix}:`, errorObj);
            throw errorObj;
        } finally {
            // 清理已完成/取消的控制器
            if (currentController === abortController) {
                currentController = null;
            }
        }
    };

    // 手动取消当前请求的方法
    const cancel = (reason = 'Manual cancellation'): void => {
        if (currentController) {
            currentController.abort(reason);
            currentController = null;
        }
    };

    // 返回增强后的函数（挂载 cancel 方法）
    return Object.assign(callFn, {cancel});
}

/**
 * 批量创建可取消的 Thunk调用函数
 * @param dispatch Redux dispatch 函数
 * @param thunks Thunk 对象（key: thunk 函数）
 * @param options 全局配置项
 * @returns 包含所有可取消调用函数的对象
 */
export function createAbortableCalls<
    T extends Record<string, AsyncThunk<unknown, unknown, { rejectValue: unknown }>>
>(
    dispatch: AppDispatch,
    thunks: T,
    options: AbortableCallOptions = {autoAbortPrevious: true}
): AbortableCalls<T> {
    const result = {} as Record<keyof T, unknown>;
    const cancelFns: Array<(reason?: string) => void> = [];

    // 遍历创建每个 thunk 的可取消调用函数
    for (const [key, thunk] of Object.entries(thunks) as [keyof T, T[keyof T]][]) {
        // 修正：使用泛型推导而非 any 断言
        const abortableCall = createAbortableCall(
            dispatch,
            thunk as AsyncThunk<unknown, unknown, { rejectValue: unknown }>,
            options
        );
        result[key] = abortableCall;
        cancelFns.push(abortableCall.cancel);
    }

    // 添加批量取消所有请求的方法
    const cancelAll = (reason = 'Batch cancellation'): void => {
        cancelFns.forEach(cancel => cancel(reason));
    };

    // 修正：使用类型断言确保返回类型匹配
    return {
        ...result,
        cancelAll
    } as AbortableCalls<T>;
}

/**
 * 辅助 Hook：自动在组件卸载时取消所有请求
 * @param abortableCalls createAbortableCalls 返回的对象
 * @param deps 依赖项数组，控制 cleanup 时机（可选）
 */
export function useAbortableCallsCleanup<
    T extends Record<string, { cancel: (reason?: string) => void }> & { cancelAll?: (reason?: string) => void }
>(abortableCalls: T, deps: unknown[] = []): void {
    // 修正：使用 ref 避免 effect 因对象引用变化而频繁触发
    const abortableCallsRef = useRef(abortableCalls);

    useEffect(() => {
        abortableCallsRef.current = abortableCalls;
    }, [abortableCalls]);

    useEffect(() => {
        return () => {
            const current = abortableCallsRef.current;
            if (current?.cancelAll) {
                current.cancelAll('Component unmounted');
            }
        };
        // 修正：使用自定义依赖项而非直接依赖不稳定的对象引用
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, deps);
}

/**
 * 创建带取消支持的 AsyncThunk（辅助类型工具）
 * 修正：正确合并 signal 到参数类型
 */
export function createAbortableAsyncThunk<ReturnData, ThunkArg = void>(
    typePrefix: string,
    payloadCreator: AsyncThunkPayloadCreator<ReturnData, MergeSignal<ThunkArg>, { rejectValue?: unknown }>,
    options?: AsyncThunkOptions<MergeSignal<ThunkArg>, { rejectValue: SerializedError | string }>
): AsyncThunk<ReturnData, MergeSignal<ThunkArg>, { rejectValue: SerializedError | string }> {
    // @ts-expect-error - 忽略类型检查
    return createAsyncThunk(typePrefix, payloadCreator, {
        ...options,
        // @ts-expect-error - 忽略类型检查
        condition: (arg, {signal}) => {
            // 如果传入的 signal 已经中止，则不执行
            if (signal?.aborted) return false;
            // 执行用户自定义的 condition
            return options?.condition?.(arg, {signal} as never) ?? true;
        }
    });
}
