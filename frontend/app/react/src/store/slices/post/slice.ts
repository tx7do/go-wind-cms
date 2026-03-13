import {createSlice, createAsyncThunk} from '@reduxjs/toolkit';
import {createPostServiceClient} from "@/api/generated/app/service/v1";
import {requestApi} from "@/transport/rest";
import {
    contentservicev1_Post,
    contentservicev1_PostTranslation,
} from '@/api/generated/app/service/v1';
import {currentLocaleLanguageCode} from "@/i18n";
import {makeOrderBy, makeQueryString} from "@/transport/rest/utils";

export interface PostState {
    list: contentservicev1_Post[];
    detail: contentservicev1_Post | null;
    loading: boolean;
    total: number;
}

// 初始状态
const initialState: PostState = {
    list: [],
    detail: null,
    loading: false,
    total: 0,
};

// 获取文章列表
export const listPost = createAsyncThunk(
    'post/listPost',
    async (
        params: {
            paging?: { page?: number; pageSize?: number };
            formValues?: object | undefined;
            fieldMask?: string | undefined;
            orderBy?: string[] | undefined;
        },
        {rejectWithValue}
    ) => {
        const postService = createPostServiceClient(requestApi);

        try {
            const locale = currentLocaleLanguageCode();
            const formValues = {...(params.formValues || {}), locale};
            const noPaging =
                params.paging?.page === undefined && params.paging?.pageSize === undefined;
            return await postService.List({
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

// 获取单篇文章
export const getPost = createAsyncThunk(
    'post/getPost',
    async (
        params: { id: number },
        {rejectWithValue}
    ) => {
        const postService = createPostServiceClient(requestApi);

        try {
            return await postService.Get({id: params.id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

// 创建文章
export const createPost = createAsyncThunk(
    'post/createPost',
    async (values: Partial<contentservicev1_Post>, {rejectWithValue}) => {
        const postService = createPostServiceClient(requestApi);

        try {
            return await postService.Create({data: values as contentservicev1_Post});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

// 更新文章
export const updatePost = createAsyncThunk(
    'post/updatePost',
    async (
        params: { id: number; values: Partial<contentservicev1_Post> },
        {rejectWithValue}
    ) => {
        const postService = createPostServiceClient(requestApi);

        try {
            return await postService.Update({
                id: params.id,
                data: params.values as contentservicev1_Post,
                updateMask: Object.keys(params.values ?? []).join(','),
            });
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

// 删除文章
export const deletePost = createAsyncThunk(
    'post/deletePost',
    async (id: number, {rejectWithValue}) => {
        const postService = createPostServiceClient(requestApi);

        try {
            return await postService.Delete({id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

const postSlice = createSlice({
    name: 'post',
    initialState,
    reducers: {
        // 清空文章详情
        clearPostDetail: (state) => {
            state.detail = null;
        },
        // 重置状态
        resetPost: (state) => {
            state.list = [];
            state.detail = null;
            state.loading = false;
            state.total = 0;
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(listPost.pending, (state) => {
                state.loading = true;
            })
            .addCase(listPost.fulfilled, (state, action) => {
                state.loading = false;
                state.list = action.payload.items || [];
                state.total = action.payload.total || 0;
            })
            .addCase(listPost.rejected, (state) => {
                state.loading = false;
            })
            .addCase(getPost.pending, (state) => {
                state.loading = true;
            })
            .addCase(getPost.fulfilled, (state, action) => {
                state.loading = false;
                state.detail = action.payload || null;
            })
            .addCase(getPost.rejected, (state) => {
                state.loading = false;
            })
            .addCase(createPost.pending, (state) => {
                state.loading = true;
            })
            .addCase(createPost.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(createPost.rejected, (state) => {
                state.loading = false;
            })
            .addCase(updatePost.pending, (state) => {
                state.loading = true;
            })
            .addCase(updatePost.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(updatePost.rejected, (state) => {
                state.loading = false;
            })
            .addCase(deletePost.pending, (state) => {
                state.loading = true;
            })
            .addCase(deletePost.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(deletePost.rejected, (state) => {
                state.loading = false;
            });
    },
});

export const {clearPostDetail, resetPost} = postSlice.actions;

/**
 * 获取帖子的翻译
 */
export function getTranslation(post: contentservicev1_Post): contentservicev1_PostTranslation | null {
    if (!post?.translations || post.translations.length === 0) return null;

    const locale = currentLocaleLanguageCode();
    // 优先查找当前语言的翻译
    const translation = post.translations?.find((t: contentservicev1_PostTranslation) => t.languageCode === locale);
    // 如果找不到，回退到第一个翻译
    return translation || post.translations?.[0];
}

/**
 * 获取帖子标题
 */
export function getPostTitle(post: contentservicev1_Post): string {
    const translation = getTranslation(post);
    return translation?.title || '';
}

/**
 * 获取帖子摘要
 */
export function getPostSummary(post: contentservicev1_Post): string {
    const translation = getTranslation(post);
    return translation?.summary || '';
}

/**
 * 获取帖子缩略图
 */
export function getPostThumbnail(post: contentservicev1_Post): string {
    const translation = getTranslation(post);
    return translation?.thumbnail || '/placeholder.jpg';
}

/**
 * 获取帖子内容
 */
export function getPostContent(post: contentservicev1_Post): string {
    const translation = getTranslation(post);
    return translation?.content || '';
}

export default postSlice.reducer;
