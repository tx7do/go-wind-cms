import {createAsyncThunk, createSlice} from '@reduxjs/toolkit';

import {
    contentservicev1_Category,
    contentservicev1_CategoryTranslation,
    createCategoryServiceClient,
} from '@/api/generated/app/service/v1';
import {requestApi} from '@/transport/rest';
import {makeOrderBy, makeQueryString, makeUpdateMask} from "@/transport/rest/utils";
import {currentLocaleLanguageCode} from "@/i18n";

export interface CategoryState {
    list: contentservicev1_Category[];
    detail: contentservicev1_Category | null;
    loading: boolean;
    total: number;
}

const initialState: CategoryState = {
    list: [],
    detail: null,
    loading: false,
    total: 0,
};

export const listCategory = createAsyncThunk(
    'category/listCategory',
    async (
        params: {
            paging?: { page?: number; pageSize?: number };
            formValues?: object | undefined;
            fieldMask?: string | undefined;
            orderBy?: string[] | undefined;
        },
        {rejectWithValue}
    ) => {
        const categoryService = createCategoryServiceClient(requestApi);

        try {
            const locale = currentLocaleLanguageCode();
            const formValues = {...(params.formValues || {}), locale};
            const noPaging =
                params.paging?.page === undefined && params.paging?.pageSize === undefined;
            return await categoryService.List({
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

export const getCategory = createAsyncThunk(
    'category/getCategory',
    async (
        params: { id: number; fieldMask?: string | undefined },
        {rejectWithValue}
    ) => {
        const categoryService = createCategoryServiceClient(requestApi);

        try {
            const locale = currentLocaleLanguageCode();
            return await categoryService.Get({
                id: params.id,
                locale,
                viewMask: params.fieldMask,
            });
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const createCategory = createAsyncThunk(
    'category/createCategory',
    async (values: Partial<contentservicev1_Category>, {rejectWithValue}) => {
        const categoryService = createCategoryServiceClient(requestApi);

        try {
            return await categoryService.Create({
                data: {
                    ...values,
                    translations: values.translations ?? [],
                    availableLanguages: values.availableLanguages ?? [],
                    customFields: values.customFields ?? [],
                    children: values.children ?? [],
                } as contentservicev1_Category,
            });
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);
export const updateCategory = createAsyncThunk(
    'category/updateCategory',
    async (
        params: { id: number; values: Partial<contentservicev1_Category> },
        {rejectWithValue}
    ) => {
        const categoryService = createCategoryServiceClient(requestApi);

        try {
            return await categoryService.Update({
                id: params.id,
                data: {
                    ...params.values,
                    translations: params.values.translations ?? [],
                    availableLanguages: params.values.availableLanguages ?? [],
                    customFields: params.values.customFields ?? [],
                    children: params.values.children ?? [],
                } as contentservicev1_Category,
                updateMask: makeUpdateMask(Object.keys(params.values ?? [])),
            });
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const deleteCategory = createAsyncThunk(
    'category/deleteCategory',
    async (id: number, {rejectWithValue}) => {
        const categoryService = createCategoryServiceClient(requestApi);

        try {
            return await categoryService.Delete({id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

/**
 * 获取分类的翻译
 */
export function getTranslation(category: contentservicev1_Category | null) {
    if (!category || !category?.translations || category.translations.length === 0) return null;
    const locale = currentLocaleLanguageCode();
    const translation = category.translations?.find(
        (t: contentservicev1_CategoryTranslation) => t.languageCode === locale
    );
    return translation || category.translations?.[0];
}

/**
 * 获取分类名称（支持国际化）
 * @param category 分类对象
 * @param t 可选的翻译函数，如果在组件外部调用且需要国际化支持，请传入 t('page.categories')
 */
export function getCategoryName(category: contentservicev1_Category | null, t?: (key: string) => string) {
    const translation = getTranslation(category);
    return translation?.name || (t ? t('category_untitled') : '未命名分类');
}

/**
 * 获取分类描述
 */
export function getCategoryDescription(category: contentservicev1_Category | null) {
    const translation = getTranslation(category);
    return translation?.description || '';
}

/**
 * 获取分类缩略图
 */
export function getCategoryThumbnail(category: contentservicev1_Category | null) {
    const translation = getTranslation(category);
    return translation?.thumbnail || '/placeholder.jpg';
}

const categorySlice = createSlice({
    name: 'category',
    initialState,
    reducers: {
        clearCategoryDetail: (state) => {
            state.detail = null;
        },
        resetCategory: (state) => {
            state.list = [];
            state.detail = null;
            state.loading = false;
            state.total = 0;
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(listCategory.pending, (state) => {
                state.loading = true;
            })
            .addCase(listCategory.fulfilled, (state, action) => {
                state.loading = false;
                state.list = action.payload.items || [];
                state.total = action.payload.total || 0;
            })
            .addCase(listCategory.rejected, (state) => {
                state.loading = false;
            })
            .addCase(getCategory.pending, (state) => {
                state.loading = true;
            })
            .addCase(getCategory.fulfilled, (state, action) => {
                state.loading = false;
                state.detail = action.payload || null;
            })
            .addCase(getCategory.rejected, (state) => {
                state.loading = false;
            })
            .addCase(createCategory.pending, (state) => {
                state.loading = true;
            })
            .addCase(createCategory.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(createCategory.rejected, (state) => {
                state.loading = false;
            })
            .addCase(updateCategory.pending, (state) => {
                state.loading = true;
            })
            .addCase(updateCategory.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(updateCategory.rejected, (state) => {
                state.loading = false;
            })
            .addCase(deleteCategory.pending, (state) => {
                state.loading = true;
            })
            .addCase(deleteCategory.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(deleteCategory.rejected, (state) => {
                state.loading = false;
            });
    },
});

export const {clearCategoryDetail, resetCategory} = categorySlice.actions;
export default categorySlice.reducer;
