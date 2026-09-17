import { useEffect, useState } from 'react';
import { Button, Image, Upload, App, Spin } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

import { uploadMediaAsset } from '@/api/hooks/media-asset';
import { compressImageFile } from '@/utils/image';

interface PostImageFieldProps {
  value?: string;
  onChange?: (value: string) => void;
  width?: number;
  height?: number;
}

/**
 * 缩略图/OG 图字段：预览卡片 + 上传/更换/移除。
 * 上传期间的本地即时预览（blob URL），成功后切换为远端地址。
 */
const PostImageField = ({
  value,
  onChange,
  width = 128,
  height = 80,
}: PostImageFieldProps) => {
  const { t } = useTranslation('post');
  const { message } = App.useApp();
  const [uploading, setUploading] = useState(false);
  // 上传进度百分比（-1 表示未知/不可得）
  const [uploadPercent, setUploadPercent] = useState(-1);
  const [localPreview, setLocalPreview] = useState('');

  // blob URL 生命周期：替换/卸载时回收
  useEffect(() => {
    return () => {
      if (localPreview) {
        URL.revokeObjectURL(localPreview);
      }
    };
  }, [localPreview]);

  const displayUrl = localPreview || value || '';

  const handleUpload = async (file: File) => {
    if (!file) return;
    const preview = URL.createObjectURL(file);
    setLocalPreview(preview);
    setUploading(true);
    setUploadPercent(-1);
    try {
      const compressed = await compressImageFile(file);
      const resp = await uploadMediaAsset(
        {},
        compressed,
        (progressEvent: { loaded?: number; total?: number }) => {
          const total = progressEvent?.total ?? 0;
          setUploadPercent(
            total > 0 ? Math.floor(((progressEvent.loaded ?? 0) / total) * 100) : -1,
          );
        },
      );
      const url = (resp as { objectName?: string })?.objectName || '';
      if (!url) {
        message.error(t('publishFailed'));
        return;
      }
      onChange?.(url);
    } catch (error) {
      console.error('Upload image failed:', error);
      message.error((error as Error)?.message || t('publishFailed'));
    } finally {
      setUploading(false);
      setUploadPercent(-1);
      setLocalPreview('');
    }
  };

  return (
    <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
      {/* 预览卡片：有图可点击放大，无图显示上传占位 */}
      <div
        style={{
          position: 'relative',
          flexShrink: 0,
          overflow: 'hidden',
          borderRadius: 6,
          border: '1px solid var(--ant-color-border-secondary)',
          width,
          height,
        }}
      >
        {displayUrl ? (
          <Image src={displayUrl} width={width} height={height} style={{ objectFit: 'cover' }} />
        ) : (
          <Upload
            accept="image/*"
            showUploadList={false}
            beforeUpload={(file) => {
              handleUpload(file);
              return false;
            }}
          >
            <div
              style={{
                width: width - 2,
                height: height - 2,
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                gap: 4,
                cursor: 'pointer',
                color: 'var(--ant-color-text-tertiary)',
                background: 'var(--ant-color-fill-quaternary)',
              }}
            >
              <PlusOutlined style={{ fontSize: 18 }} />
              <span style={{ fontSize: 12 }}>{t('upload')}</span>
            </div>
          </Upload>
        )}

        {uploading && (
          <div
            style={{
              position: 'absolute',
              inset: 0,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              gap: 8,
              background: 'rgba(0,0,0,0.4)',
            }}
          >
            <Spin size="small" />
            <span style={{ fontSize: 12, color: '#fff' }}>
              {t('uploading')}
              {uploadPercent >= 0 ? ` ${uploadPercent}%` : ''}
            </span>
          </div>
        )}
      </div>

      {/* 操作按钮仅在已有图片时出现；无图时占位卡本身即上传入口 */}
      {value ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <Upload
            accept="image/*"
            showUploadList={false}
            beforeUpload={(file) => {
              handleUpload(file);
              return false;
            }}
          >
            <Button size="small" loading={uploading}>
              {t('replaceImage')}
            </Button>
          </Upload>
          <Button size="small" danger onClick={() => onChange?.('')}>
            {t('removeImage')}
          </Button>
        </div>
      ) : null}
    </div>
  );
};

export default PostImageField;
