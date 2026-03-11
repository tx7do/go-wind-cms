import {createSlice, createAsyncThunk} from '@reduxjs/toolkit';

import {createNavigationServiceClient} from '@/api/generated/app/service/v1';
import {requestClientRequestHandler} from '@/transport/rest';
import {
    siteservicev1_Navigation,
} from '@/api/generated/app/service/v1';

const navigationService = createNavigationServiceClient(requestClientRequestHandler);

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
        params: { page?: number; pageSize?: number; query?: string },
        {rejectWithValue}
    ) => {
        try {
            return await navigationService.List({
                query: params.query,
                page: params.page,
                pageSize: params.pageSize,
                sorting: undefined,
            });
        } catch (error) {
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
        try {
            return await navigationService.Delete({id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

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
