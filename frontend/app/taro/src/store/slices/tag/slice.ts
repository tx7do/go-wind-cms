import {createSlice, createAsyncThunk} from '@reduxjs/toolkit';
import {createTagServiceClient} from '../../../api/generated/app/service/v1';
import {requestApi} from '@/transport/rest';
import {
    contentservicev1_Tag,
    contentservicev1_TagTranslation,
} from '@/api/generated/app/service/v1';
import {currentLocaleLanguageCode} from '@/i18n';
import {makeOrderBy, makeQueryString} from "@/transport/rest/utils";
import {ITagState} from "@/store/types";


const initialState: ITagState = {
    list: [],
    detail: null,
    loading: false,
    total: 0,
};

export const listTag = createAsyncThunk(
    'tag/listTag',
    async (
        params: {
            paging?: { page?: number; pageSize?: number };
            formValues?: object | undefined;
            fieldMask?: string | undefined | null;
            orderBy?: string[] | undefined | null;
        },
        {rejectWithValue}
    ) => {
        const tagService = createTagServiceClient(requestApi);

        try {
            const locale = currentLocaleLanguageCode();
            const formValues = {...(params.formValues || {}), locale};
            const noPaging =
                params.paging?.page === undefined && params.paging?.pageSize === undefined;
            return await tagService.List({
                fieldMask: params.fieldMask ?? undefined,
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

export const getTag = createAsyncThunk(
    'tag/getTag',
    async (
        params: { id: number },
        {rejectWithValue}
    ) => {
        const tagService = createTagServiceClient(requestApi);

        try {
            return await tagService.Get({id: params.id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const createTag = createAsyncThunk(
    'tag/createTag',
    async (values: Partial<contentservicev1_Tag>, {rejectWithValue}) => {
        const tagService = createTagServiceClient(requestApi);
        try {
            return await tagService.Create({data: values as contentservicev1_Tag});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const updateTag = createAsyncThunk(
    'tag/updateTag',
    async (
        params: { id: number; values: Partial<contentservicev1_Tag> },
        {rejectWithValue}
    ) => {
        const tagService = createTagServiceClient(requestApi);
        try {
            return await tagService.Update({
                id: params.id,
                data: params.values as contentservicev1_Tag,
                updateMask: Object.keys(params.values ?? []).join(','),
            });
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const deleteTag = createAsyncThunk(
    'tag/deleteTag',
    async (id: number, {rejectWithValue}) => {
        const tagService = createTagServiceClient(requestApi);
        try {
            return await tagService.Delete({id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

/**
 * 获取标签的翻译
 * @param tag 标签对象
 * @returns 当前语言的翻译，如果找不到则返回第一个翻译或 null
 */
export function getTranslation(tag: contentservicev1_Tag | null) {
    if (!tag || !tag?.translations || tag.translations.length === 0) return null;

    const locale = currentLocaleLanguageCode();
    // 优先查找当前语言的翻译
    const translation = tag.translations?.find(
        (t: contentservicev1_TagTranslation) => t.languageCode === locale
    );
    // 如果找不到，回退到第一个翻译
    return translation || tag.translations?.[0];
}

const tagSlice = createSlice({
    name: 'tag',
    initialState,
    reducers: {
        clearTagDetail: (state) => {
            state.detail = null;
        },
        resetTag: (state) => {
            state.list = [];
            state.detail = null;
            state.loading = false;
            state.total = 0;
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(listTag.pending, (state) => {
                state.loading = true;
            })
            .addCase(listTag.fulfilled, (state, action) => {
                state.loading = false;
                state.list = action.payload.items || [];
                state.total = action.payload.total || 0;
            })
            .addCase(listTag.rejected, (state) => {
                state.loading = false;
            })
            .addCase(getTag.pending, (state) => {
                state.loading = true;
            })
            .addCase(getTag.fulfilled, (state, action) => {
                state.loading = false;
                state.detail = action.payload || null;
            })
            .addCase(getTag.rejected, (state) => {
                state.loading = false;
            })
            .addCase(createTag.pending, (state) => {
                state.loading = true;
            })
            .addCase(createTag.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(createTag.rejected, (state) => {
                state.loading = false;
            })
            .addCase(updateTag.pending, (state) => {
                state.loading = true;
            })
            .addCase(updateTag.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(updateTag.rejected, (state) => {
                state.loading = false;
            })
            .addCase(deleteTag.pending, (state) => {
                state.loading = true;
            })
            .addCase(deleteTag.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(deleteTag.rejected, (state) => {
                state.loading = false;
            });
    },
});

export const {clearTagDetail, resetTag} = tagSlice.actions;
export default tagSlice.reducer;
