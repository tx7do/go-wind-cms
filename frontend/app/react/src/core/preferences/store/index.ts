import {create, type StoreApi, useStore} from 'zustand';
import {persist} from 'zustand/middleware';
import {createContext, useContext, type Context} from 'react';

import {defaultPreferences} from '../config/default';
import type {DeepPartial, Preferences} from '../types';
import {mergeDeep} from "../utils/merge";


export interface PreferencesState {
    preferences: Preferences;
    setPreferences: (overrides: DeepPartial<Preferences>) => void;
    resetPreferences: () => void;
    getPreference: <K extends keyof Preferences>(key: K) => Preferences[K];
}

// ==============================
// Store 工厂函数（每个请求/客户端独立实例）
// ==============================

export function createPreferencesStore(): StoreApi<PreferencesState> {
    return create<PreferencesState>()(
        persist(
            (set, get) => ({
                preferences: defaultPreferences,

                setPreferences: (overrides) => {
                    set((state) => ({
                        preferences: mergeDeep(state.preferences, overrides),
                    }));
                },

                resetPreferences: () => {
                    set({preferences: defaultPreferences});
                },

                getPreference: (key) => {
                    return get().preferences[key];
                },
            }),
            {
                name: 'app-preferences',
                partialize: (state) => ({preferences: state.preferences}),
            }
        )
    );
}

// ==============================
// React Context
// ==============================

export const PreferencesStoreContext: Context<StoreApi<PreferencesState> | null> = createContext<StoreApi<PreferencesState> | null>(null);

// ==============================
// 消费 Hook
// ==============================

/**
 * 在组件内访问 preferences store
 * 必须在 <PreferencesStoreProvider> 内部使用
 */
export function usePreferencesStore<T = PreferencesState>(selector?: (state: PreferencesState) => T): T {
    const storeApi = useContext(PreferencesStoreContext);
    if (!storeApi) {
        throw new Error('usePreferencesStore must be used within <PreferencesStoreProvider>');
    }
    return useStore(
        storeApi,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (selector ?? ((state: PreferencesState) => state)) as (state: PreferencesState) => T,
    );
}
