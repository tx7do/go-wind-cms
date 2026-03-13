import {createSlice, createAsyncThunk} from '@reduxjs/toolkit';
import {createCommentServiceClient} from '@/api/generated/app/service/v1';
import {requestClientRequestHandler} from '@/transport/rest';
import {
    commentservicev1_Comment,
} from '@/api/generated/app/service/v1';
import {makeOrderBy, makeQueryString} from "@/transport/rest/utils";

const commentService = createCommentServiceClient(requestClientRequestHandler);

export interface CommentState {
    list: commentservicev1_Comment[];
    detail: commentservicev1_Comment | null;
    loading: boolean;
    total: number;
}

const initialState: CommentState = {
    list: [],
    detail: null,
    loading: false,
    total: 0,
};

export const listComment = createAsyncThunk(
    'comment/listComment',
    async (
        params: {
            paging?: { page?: number; pageSize?: number };
            formValues?: object | undefined;
            fieldMask?: string | undefined;
            orderBy?: string[] | undefined;
        },
        {rejectWithValue}
    ) => {
        try {
            const formValues = {...(params.formValues || {})};
            const noPaging =
                params.paging?.page === undefined && params.paging?.pageSize === undefined;
            return await commentService.List({
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

export const getComment = createAsyncThunk(
    'comment/getComment',
    async (
        params: { id: number },
        {rejectWithValue}
    ) => {
        try {
            return await commentService.Get({id: params.id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const createComment = createAsyncThunk(
    'comment/createComment',
    async (values: Partial<commentservicev1_Comment>, {rejectWithValue}) => {
        try {
            return await commentService.Create({data: values as commentservicev1_Comment});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const updateComment = createAsyncThunk(
    'comment/updateComment',
    async (
        params: { id: number; values: Partial<commentservicev1_Comment> },
        {rejectWithValue}
    ) => {
        try {
            return await commentService.Update({
                id: params.id,
                data: params.values as commentservicev1_Comment,
                updateMask: Object.keys(params.values ?? []).join(','),
            });
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const deleteComment = createAsyncThunk(
    'comment/deleteComment',
    async (id: number, {rejectWithValue}) => {
        try {
            return await commentService.Delete({id});
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

const commentSlice = createSlice({
    name: 'comment',
    initialState,
    reducers: {
        clearCommentDetail: (state) => {
            state.detail = null;
        },
        resetComment: (state) => {
            state.list = [];
            state.detail = null;
            state.loading = false;
            state.total = 0;
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(listComment.pending, (state) => {
                state.loading = true;
            })
            .addCase(listComment.fulfilled, (state, action) => {
                state.loading = false;
                state.list = action.payload.items || [];
                state.total = action.payload.total || 0;
            })
            .addCase(listComment.rejected, (state) => {
                state.loading = false;
            })
            .addCase(getComment.pending, (state) => {
                state.loading = true;
            })
            .addCase(getComment.fulfilled, (state, action) => {
                state.loading = false;
                state.detail = action.payload || null;
            })
            .addCase(getComment.rejected, (state) => {
                state.loading = false;
            })
            .addCase(createComment.pending, (state) => {
                state.loading = true;
            })
            .addCase(createComment.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(createComment.rejected, (state) => {
                state.loading = false;
            })
            .addCase(updateComment.pending, (state) => {
                state.loading = true;
            })
            .addCase(updateComment.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(updateComment.rejected, (state) => {
                state.loading = false;
            })
            .addCase(deleteComment.pending, (state) => {
                state.loading = true;
            })
            .addCase(deleteComment.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(deleteComment.rejected, (state) => {
                state.loading = false;
            });
    },
});

export const {clearCommentDetail, resetComment} = commentSlice.actions;
export default commentSlice.reducer;
