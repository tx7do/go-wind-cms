import {createSlice, createAsyncThunk} from '@reduxjs/toolkit';

import {createNavigationServiceClient} from '@/api/generated/app/service/v1';
import {requestApi} from '@/transport/rest';
import {
    siteservicev1_Navigation,
    siteservicev1_NavigationItem,
} from '@/api/generated/app/service/v1';
import {currentLocaleLanguageCode} from "@/i18n";
import {makeOrderBy, makeQueryString} from "@/transport/rest/utils";

export interface NavigationState {
    list: siteservicev1_Navigation[];
    detail: siteservicev1_Navigation | null;
    loading: boolean;
    total: number;
}

const initialState: NavigationState = {
    list: [],
    detail: null,
    loading: false,
    total: 0,
};

export const listNavigation = createAsyncThunk(
    'navigation/listNavigation',
    async (
        params: {
            paging?: { page?: number; pageSize?: number };
            formValues?: object | undefined;
            fieldMask?: string | undefined;
            orderBy?: string[] | undefined;
        },
        {rejectWithValue}
    ) => {
        const navigationService = createNavigationServiceClient(requestApi);

        try {
            const locale = currentLocaleLanguageCode();
            const formValues = {...(params.formValues || {}), locale};
            const noPaging =
                params.paging?.page === undefined && params.paging?.pageSize === undefined;
            return await navigationService.List({
                fieldMask: params.fieldMask,
                orderBy: makeOrderBy(params.orderBy),
                sorting: Array.isArray(params.orderBy) ? params.orderBy.map(o => ({
                    field: o,
                    direction: 'ASC'
                })) : undefined,
                query: makeQueryString(formValues, false),
                page: params.paging?.page,
                pageSize: params.paging?.pageSize,
                noPaging,
            });
        } catch (error) {
            console.error('xxxxxxxxxxxxxxxxxxxxxxxx', error);
            return rejectWithValue(error);
        }
    }
);

export const getNavigation = createAsyncThunk(
    'navigation/getNavigation',
    async (
        params: { id: number },
        {rejectWithValue}
    ) => {
        const navigationService = createNavigationServiceClient(requestApi);

        try {
            return await navigationService.Get({id: params.id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const createNavigation = createAsyncThunk(
    'navigation/createNavigation',
    async (values: Partial<siteservicev1_Navigation>, {rejectWithValue}) => {
        const navigationService = createNavigationServiceClient(requestApi);

        try {
            return await navigationService.Create({data: values as siteservicev1_Navigation});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const updateNavigation = createAsyncThunk(
    'navigation/updateNavigation',
    async (
        params: { id: number; values: Partial<siteservicev1_Navigation> },
        {rejectWithValue}
    ) => {
        const navigationService = createNavigationServiceClient(requestApi);

        try {
            return await navigationService.Update({
                id: params.id,
                data: params.values as siteservicev1_Navigation,
                updateMask: Object.keys(params.values ?? []).join(','),
            });
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const deleteNavigation = createAsyncThunk(
    'navigation/deleteNavigation',
    async (id: number, {rejectWithValue}) => {
        const navigationService = createNavigationServiceClient(requestApi);

        try {
            return await navigationService.Delete({id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

/**
 * 递归查找导航项
 * @param items 导航项列表
 * @param key 要查找的 ID
 * @returns 找到的导航项，未找到返回 null
 */
export function findNavItem(items: siteservicev1_NavigationItem[], key: number): siteservicev1_NavigationItem | null {
    for (const item of items) {
        if (item.id === key) return item;
        if (item.children) {
            const found = findNavItem(item.children, key);
            if (found) return found;
        }
    }
    return null;
}

const navigationSlice = createSlice({
    name: 'navigation',
    initialState,
    reducers: {
        clearNavigationDetail: (state) => {
            state.detail = null;
        },
        resetNavigation: (state) => {
            state.list = [];
            state.detail = null;
            state.loading = false;
            state.total = 0;
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(listNavigation.pending, (state) => {
                state.loading = true;
            })
            .addCase(listNavigation.fulfilled, (state, action) => {
                state.loading = false;
                state.list = action.payload.items || [];
                state.total = action.payload.total || 0;
            })
            .addCase(listNavigation.rejected, (state) => {
                state.loading = false;
            })
            .addCase(getNavigation.pending, (state) => {
                state.loading = true;
            })
            .addCase(getNavigation.fulfilled, (state, action) => {
                state.loading = false;
                state.detail = action.payload || null;
            })
            .addCase(getNavigation.rejected, (state) => {
                state.loading = false;
            })
            .addCase(createNavigation.pending, (state) => {
                state.loading = true;
            })
            .addCase(createNavigation.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(createNavigation.rejected, (state) => {
                state.loading = false;
            })
            .addCase(updateNavigation.pending, (state) => {
                state.loading = true;
            })
            .addCase(updateNavigation.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(updateNavigation.rejected, (state) => {
                state.loading = false;
            })
            .addCase(deleteNavigation.pending, (state) => {
                state.loading = true;
            })
            .addCase(deleteNavigation.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(deleteNavigation.rejected, (state) => {
                state.loading = false;
            });
    },
});

export const {clearNavigationDetail, resetNavigation} = navigationSlice.actions;
export default navigationSlice.reducer;
