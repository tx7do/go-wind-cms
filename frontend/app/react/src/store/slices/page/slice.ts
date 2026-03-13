import {createSlice, createAsyncThunk} from '@reduxjs/toolkit';
import {createPageServiceClient} from '@/api/generated/app/service/v1';
import {requestApi} from '@/transport/rest';
import {
    contentservicev1_Page,
    contentservicev1_PageTranslation,
} from '@/api/generated/app/service/v1';
import {currentLocaleLanguageCode} from '@/i18n';
import {makeOrderBy, makeQueryString} from "@/transport/rest/utils";

export interface PageState {
    list: contentservicev1_Page[];
    detail: contentservicev1_Page | null;
    loading: boolean;
    total: number;
}

const initialState: PageState = {
    list: [],
    detail: null,
    loading: false,
    total: 0,
};

export const listPage = createAsyncThunk(
    'page/listPage',
    async (
        params: {
            paging?: { page?: number; pageSize?: number };
            formValues?: object | undefined;
            fieldMask?: string | undefined;
            orderBy?: string[] | undefined;
        },
        {rejectWithValue}
    ) => {
        const pageService = createPageServiceClient(requestApi);

        try {
            const locale = currentLocaleLanguageCode();
            const formValues = {...(params.formValues || {}), locale};
            const noPaging =
                params.paging?.page === undefined && params.paging?.pageSize === undefined;
            return await pageService.List({
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
            return rejectWithValue(error);
        }
    }
);

export const getPage = createAsyncThunk(
    'page/getPage',
    async (
        params: { id: number },
        {rejectWithValue}
    ) => {
        const pageService = createPageServiceClient(requestApi);

        try {
            return await pageService.Get({id: params.id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const createPage = createAsyncThunk(
    'page/createPage',
    async (values: Partial<contentservicev1_Page>, {rejectWithValue}) => {
        const pageService = createPageServiceClient(requestApi);

        try {
            return await pageService.Create({data: values as contentservicev1_Page});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const updatePage = createAsyncThunk(
    'page/updatePage',
    async (
        params: { id: number; values: Partial<contentservicev1_Page> },
        {rejectWithValue}
    ) => {
        const pageService = createPageServiceClient(requestApi);

        try {
            return await pageService.Update({
                id: params.id,
                data: params.values as contentservicev1_Page,
                updateMask: Object.keys(params.values ?? []).join(','),
            });
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const deletePage = createAsyncThunk(
    'page/deletePage',
    async (id: number, {rejectWithValue}) => {
        const pageService = createPageServiceClient(requestApi);

        try {
            return await pageService.Delete({id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

/**
 * 获取页面的翻译
 * @param page 页面对象
 * @returns 当前语言的翻译，如果找不到则返回第一个翻译或 null
 */
export function getTranslation(page: contentservicev1_Page) {
    if (!page?.translations || page.translations.length === 0) return null;

    const locale = currentLocaleLanguageCode();
    // 优先查找当前语言的翻译
    const translation = page.translations?.find(
        (t: contentservicev1_PageTranslation) => t.languageCode === locale
    );
    // 如果找不到，回退到第一个翻译
    return translation || page.translations?.[0];
}

const pageSlice = createSlice({
    name: 'page',
    initialState,
    reducers: {
        clearPageDetail: (state) => {
            state.detail = null;
        },
        resetPage: (state) => {
            state.list = [];
            state.detail = null;
            state.loading = false;
            state.total = 0;
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(listPage.pending, (state) => {
                state.loading = true;
            })
            .addCase(listPage.fulfilled, (state, action) => {
                state.loading = false;
                state.list = action.payload.items || [];
                state.total = action.payload.total || 0;
            })
            .addCase(listPage.rejected, (state) => {
                state.loading = false;
            })
            .addCase(getPage.pending, (state) => {
                state.loading = true;
            })
            .addCase(getPage.fulfilled, (state, action) => {
                state.loading = false;
                state.detail = action.payload || null;
            })
            .addCase(getPage.rejected, (state) => {
                state.loading = false;
            })
            .addCase(createPage.pending, (state) => {
                state.loading = true;
            })
            .addCase(createPage.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(createPage.rejected, (state) => {
                state.loading = false;
            })
            .addCase(updatePage.pending, (state) => {
                state.loading = true;
            })
            .addCase(updatePage.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(updatePage.rejected, (state) => {
                state.loading = false;
            })
            .addCase(deletePage.pending, (state) => {
                state.loading = true;
            })
            .addCase(deletePage.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(deletePage.rejected, (state) => {
                state.loading = false;
            });
    },
});

export const {clearPageDetail, resetPage} = pageSlice.actions;
export default pageSlice.reducer;
