import {create} from 'zustand';
import {persist} from 'zustand/middleware';

import {defaultPreferences} from '../config/default';
import type {DeepPartial, Preferences} from '../types';
import {mergeDeep} from "../utils/merge";


export interface PreferencesState {
    preferences: Preferences;
    setPreferences: (overrides: DeepPartial<Preferences>) => void;
    resetPreferences: () => void;
    getPreference: <K extends keyof Preferences>(key: K) => Preferences[K];
}

export const usePreferencesStore = create<PreferencesState>()(
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
            // v1 → v2：设计语言规范（docs/design-language.md）回归 vben 深蓝
            // hsl(212 100% 45%) 并统一默认圆角 8。历史上 v0→v1 曾把默认主色迁到
            // #3B82F6；本版把默认主题的两种旧蓝一并收敛到规范值。仅当主题仍为
            // 内置 default 时迁移，用户自选的主题色/圆角不受影响。
            version: 2,
            migrate: (persistedState, version) => {
                const state = (persistedState ?? {}) as { preferences?: Preferences };
                if (version < 2) {
                    const themePref = state.preferences?.theme;
                    if (themePref && themePref.builtinType === 'default') {
                        if (
                            themePref.colorPrimary === 'hsl(212 100% 45%)' ||
                            themePref.colorPrimary === '#3B82F6'
                        ) {
                            themePref.colorPrimary = 'hsl(212 100% 45%)';
                        }
                        if (!themePref.radius || themePref.radius === '6') {
                            themePref.radius = '8';
                        }
                        console.info(
                            '[Preferences] 迁移 v<2：默认主题对齐设计语言规范（hsl(212 100% 45%) + radius 8）'
                        );
                    }
                }
                return state as { preferences: Preferences };
            },
            partialize: (state) => ({preferences: state.preferences}),
        }
    )
);
