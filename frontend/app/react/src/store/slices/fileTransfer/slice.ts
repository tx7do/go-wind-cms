import {createSlice, createAsyncThunk} from '@reduxjs/toolkit';
import {createFileTransferServiceClient} from '@/api/generated/app/service/v1';
import {requestClientRequestHandler} from '@/transport/rest';
import {
    storageservicev1_UploadFileResponse,
    storageservicev1_DownloadFileResponse,
} from '@/api/generated/app/service/v1';
import {requestClient} from '@/transport/rest/rest-client';

const fileTransferService = createFileTransferServiceClient(requestClientRequestHandler);

export interface FileTransferState {
    list: storageservicev1_UploadFileResponse[];
    detail: storageservicev1_DownloadFileResponse | null;
    loading: boolean;
    total: number;
}

const initialState: FileTransferState = {
    list: [],
    detail: null,
    loading: false,
    total: 0,
};

export const downloadFileThunk = createAsyncThunk(
    'fileTransfer/downloadFile',
    async (
        params: { bucketName: string; objectName: string; preferPresignedUrl: boolean },
        {rejectWithValue}
    ) => {
        try {
            const {bucketName, objectName, preferPresignedUrl} = params;
            return await fileTransferService.DownloadFile({
                storageObject: {bucketName, objectName},
                preferPresignedUrl,
            });
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

export const uploadFileThunk = createAsyncThunk(
    'fileTransfer/uploadFile',
    async (
        params: {
            bucketName: string;
            fileDirectory: string;
            fileData: File;
            method?: 'post' | 'put';
            onUploadProgress?: (progressEvent: ProgressEvent) => void;
        },
        {rejectWithValue}
    ) => {
        try {
            const {bucketName, fileDirectory, fileData, method = 'post', onUploadProgress} = params;
            const storageObject = JSON.stringify({bucketName, fileDirectory});
            const formData = new FormData();
            formData.append('file', fileData);
            formData.append('storageObject', storageObject);
            formData.append('sourceFileName', fileData.name);
            formData.append('mime', fileData.type);
            formData.append('size', String(fileData.size));
            formData.append('method', method);
            await requestClient.upload(
                'app/v1/file/upload',
                formData,
                {onUploadProgress}
            );
            return {success: true};
        } catch (error) {
            return rejectWithValue(error);
        }
    }
);

function normalizeBase64(s: string): string {
    let str = s.replaceAll(/\s+/g, '');
    str = str.replaceAll('-', '+').replaceAll('_', '/');
    while (str.length % 4 !== 0) str += '=';
    return str;
}

function toBlob(data: unknown, type = 'application/octet-stream'): Blob {
    if (!data) return new Blob([], {type});
    if (data instanceof Blob) return data;
    if (data instanceof ArrayBuffer) return new Blob([data], {type});
    if (ArrayBuffer.isView(data)) {
        const view = new Uint8Array((data as ArrayBufferView).buffer, (data as ArrayBufferView).byteOffset, (data as ArrayBufferView).byteLength);
        const bytes = new Uint8Array(view.byteLength);
        bytes.set(view);
        return new Blob([bytes], {type});
    }
    if (typeof data === 'string') {
        const maybeBase64 = (data as string).includes('base64,') ? (data as string).split('base64,')[1] : data as string;
        const base64 = normalizeBase64(maybeBase64 ?? '');
        let binary = '';
        try {
            binary = atob(base64);
        } catch {
            return new Blob([], {type});
        }
        const len = binary.length;
        const arr = new Uint8Array(len);
        for (let i = 0; i < len; i++) {
            arr[i] = (binary.codePointAt(i) ?? 0) & 0xff;
        }
        return new Blob([arr], {type});
    }
    return new Blob([data as any], {type});
}

const fileTransferSlice = createSlice({
    name: 'fileTransfer',
    initialState,
    reducers: {
        clearFileDetail: (state) => {
            state.detail = null;
        },
        resetFileTransfer: (state) => {
            state.list = [];
            state.detail = null;
            state.loading = false;
            state.total = 0;
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(downloadFileThunk.pending, (state) => {
                state.loading = true;
            })
            .addCase(downloadFileThunk.fulfilled, (state, action) => {
                state.loading = false;
                state.detail = action.payload || null;
            })
            .addCase(downloadFileThunk.rejected, (state) => {
                state.loading = false;
            })
            .addCase(uploadFileThunk.pending, (state) => {
                state.loading = true;
            })
            .addCase(uploadFileThunk.fulfilled, (state) => {
                state.loading = false;
            })
            .addCase(uploadFileThunk.rejected, (state) => {
                state.loading = false;
            });
    },
});

export const {clearFileDetail, resetFileTransfer} = fileTransferSlice.actions;
export default fileTransferSlice.reducer;
