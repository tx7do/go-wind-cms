import type { SkeletonBlockProps } from '../Skeleton/types.ts';

/**
 * 根据数据结构动态生成骨架配置
 * @param dataSample - 示例数据（用于推断结构）
 * @param options - 生成选项
 */
export const generateSkeletonBlocks = (
  dataSample: any,
  options: {
    maxDepth?: number;
    includeArrays?: boolean;
    arrayLength?: number;
  } = {},
): SkeletonBlockProps[] => {
  const { maxDepth = 2, includeArrays = true, arrayLength = 3 } = options;
  const blocks: SkeletonBlockProps[] = [];

  const traverse = (value: any, depth: number) => {
    if (depth > maxDepth) return;

    if (value === null || value === undefined) {
      blocks.push({ width: 60, height: '1em' });
    } else if (typeof value === 'string') {
      blocks.push({ width: `${Math.min(100, value.length * 2)}%`, height: '1em' });
    } else if (typeof value === 'number') {
      blocks.push({ width: 40, height: '1.2em' });
    } else if (typeof value === 'boolean') {
      blocks.push({ width: 30, height: 20, shape: 'round' });
    } else if (Array.isArray(value) && includeArrays) {
      for (let i = 0; i < Math.min(arrayLength, value.length); i++) {
        traverse(value[i], depth + 1);
      }
    } else if (typeof value === 'object') {
      Object.values(value).forEach((v) => traverse(v, depth + 1));
    }
  };

  traverse(dataSample, 0);
  return blocks;
};
