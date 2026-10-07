import type {PluginOption} from 'vite';

import archiver from 'archiver';
import fs from 'fs';
import path from 'path';

/**
 * 构建完成后将 dist 目录打包为 dist.zip，便于交付部署
 * （对齐 vue-vben 的构建产物行为，实现见 internal/vite-config/src/plugins/archiver.ts）
 * apply: 'build' 确保 dev server 不加载本插件，仅生产构建生效
 */
export const archiverPlugin = (): PluginOption => {
    return {
        apply: 'build',
        closeBundle: {
            handler() {
                setTimeout(async () => {
                    const zipOutputPath = path.join(process.cwd(), 'dist.zip');
                    try {
                        await zipFolder('dist', zipOutputPath);
                        console.log(`Folder has been zipped to: ${zipOutputPath}`);
                    } catch (error) {
                        console.error('Error zipping folder:', error);
                    }
                }, 0);
            },
            order: 'post',
        },
        enforce: 'post',
        name: 'vite:archiver',
    };
};

/**
 * 流式压缩指定目录为 zip 文件
 */
async function zipFolder(
    folderPath: string,
    outputPath: string,
): Promise<void> {
    return new Promise((resolve, reject) => {
        const output = fs.createWriteStream(outputPath);
        const archive = archiver('zip', {
            zlib: {level: 9}, // 最高压缩率
        });

        output.on('close', () => {
            console.log(
                `ZIP file created: ${outputPath} (${archive.pointer()} total bytes)`,
            );
            resolve();
        });

        archive.on('error', (err) => {
            reject(err);
        });

        archive.pipe(output);

        // 以流方式压缩目录，减少内存占用
        archive.directory(folderPath, false);

        archive.finalize();
    });
}
