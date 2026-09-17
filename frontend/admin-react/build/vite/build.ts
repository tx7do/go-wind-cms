import type {BuildOptions} from 'vite';
import {splitJSModules} from '../utils/helper';

/**
 * @description 分包配置
 */
export function buildOptions(): BuildOptions {
    return {
        chunkSizeWarningLimit: 1000, // 大于1000k才警告
        sourcemap: process.env.NODE_ENV !== 'production', // 非生产环境开启
        minify: 'esbuild',
        esbuild: {
            drop: ['console', 'debugger'],
        },
        rollupOptions: {
            output: {
                chunkFileNames: 'assets/js/[name].[hash].js',
                entryFileNames: 'assets/js/[name].[hash].js',
                assetFileNames: 'assets/[ext]/[name].[hash].[ext]',
                manualChunks(id) {
                    // JS模块
                    if (id.includes('node_modules')) {
                        return splitJSModules(id);
                    }
                },
            },
        },
    };
}
