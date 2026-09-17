// ==============================
// 上传前图片压缩（仅照片类源，画图/截图类 PNG 原样直传）
// ==============================

/** 参与压缩的源类型：相机照片常见封装。PNG/SVG/GIF 是图表/截图主载体，保持原样。 */
const COMPRESSIBLE_IMAGE_TYPES = new Set(['image/jpeg', 'image/webp']);
/** 长边上限，超过才等比缩小；不足则只重编码不缩放 */
const IMAGE_MAX_LONG_EDGE = 1920;
/** 低于该体积分辨率再高也不值得处理，直传 */
const IMAGE_COMPRESS_MIN_BYTES = 300 * 1024;
/** WebP 重编码质量因子 */
const IMAGE_REENCODE_QUALITY = 0.82;

async function reencodeImage(file: File): Promise<File | null> {
  try {
    const bitmap = await createImageBitmap(file, {
      imageOrientation: 'from-image',
    });
    const longEdge = Math.max(bitmap.width, bitmap.height);
    const scale = longEdge > IMAGE_MAX_LONG_EDGE ? IMAGE_MAX_LONG_EDGE / longEdge : 1;
    const width = Math.round(bitmap.width * scale);
    const height = Math.round(bitmap.height * scale);
    if (width < 1 || height < 1) {
      bitmap.close();
      return null;
    }
    const canvas = document.createElement('canvas');
    canvas.width = width;
    canvas.height = height;
    const ctx = canvas.getContext('2d');
    if (!ctx) {
      bitmap.close();
      return null;
    }
    ctx.drawImage(bitmap, 0, 0, width, height);
    bitmap.close();
    const blob = await new Promise<Blob | null>((resolve) => {
      canvas.toBlob(resolve, 'image/webp', IMAGE_REENCODE_QUALITY);
    });
    if (!blob || blob.size === 0 || blob.size >= file.size) {
      return null;
    }
    const name = `${file.name.replace(/\.[^.]+$/, '')}.webp`;
    return new File([blob], name, { type: 'image/webp' });
  } catch {
    return null;
  }
}

/**
 * 上传前的图片预处理：照片类（jpeg/webp）等比缩小到长边上限内并重编码为 WebP
 * （同时剥离 EXIF 等元数据）；任何一步失败或无收益都原样返回，绝不阻断上传。
 */
export async function compressImageFile(file: File): Promise<File> {
  if (!COMPRESSIBLE_IMAGE_TYPES.has(file.type)) {
    return file;
  }
  if (file.size < IMAGE_COMPRESS_MIN_BYTES) {
    return file;
  }
  const out = await reencodeImage(file);
  return out ?? file;
}
