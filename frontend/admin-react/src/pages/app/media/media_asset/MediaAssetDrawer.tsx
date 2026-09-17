import { useEffect, useRef, useState } from 'react';
import type { mediaservicev1_MediaAsset } from '@/api/generated/admin/service/v1';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  DrawerForm,
  ProFormText,
  ProFormSelect,
  ProFormSwitch,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { App } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';

import {
  useCreateMediaAsset,
  useUpdateMediaAsset,
} from '@/api/hooks/media-asset';
import { assetTypeOptions } from './constants';

interface MediaAssetDrawerProps {
  open: boolean;
  mode: 'create' | 'edit';
  data?: any;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 媒体资源编辑抽屉组件（媒体文件由上传入口创建，此抽屉仅编辑元数据）
 */
const MediaAssetDrawer: React.FC<MediaAssetDrawerProps> = ({
  open,
  mode,
  data,
  onClose,
  onSuccess,
}) => {
  const { t } = useTranslation('media-asset');
  const formRef = useRef<ProFormInstance>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [confirmLoading, setConfirmLoading] = useState(false);

  // 编辑模式下设置表单值（destroyOnHidden 时需延迟赋值）
  useEffect(() => {
    if (open && mode === 'edit' && data) {
      setTimeout(() => {
        formRef.current?.setFieldsValue({
          filename: data.filename || '',
          type: data.type,
          title: data.title || '',
          caption: data.caption || '',
          altText: data.altText || '',
          isPrivate: data.isPrivate ?? false,
        });
      }, 0);
    }
  }, [open, mode, data]);

  // 创建 mutation
  const createMutation = useCreateMediaAsset({
    onSuccess: () => {
      message.success(t('createSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listMediaAssets'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('createFailed'));
    },
  });

  // 更新 mutation
  const updateMutation = useUpdateMediaAsset({
    onSuccess: () => {
      message.success(t('updateSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listMediaAssets'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('updateFailed'));
    },
  });

  // 提交表单
  const handleSubmit = async (values: Record<string, any>) => {
    try {
      setConfirmLoading(true);
      const payload: Record<string, any> = {
        filename: values.filename,
        type: values.type,
        title: values.title,
        caption: values.caption,
        altText: values.altText,
        isPrivate: values.isPrivate ?? false,
      };

      if (mode === 'edit' && data?.id) {
        // filename 是身份字段：进 data 但不进 updateMask（对齐 vben 侧）
        await updateMutation.mutateAsync({
          id: data.id,
          values: payload as unknown as mediaservicev1_MediaAsset,
          maskKeys: Object.keys(payload).filter((k) => k !== 'filename'),
        });
      } else {
        await createMutation.mutateAsync({ data: payload as unknown as mediaservicev1_MediaAsset });
      }
      return true;
    } catch {
      return false;
    } finally {
      setConfirmLoading(false);
    }
  };

  return (
    <DrawerForm
      formRef={formRef}
      title={mode === 'create' ? t('create') : t('edit')}
      open={open}
      onOpenChange={(visible) => {
        if (!visible) {
          formRef.current?.resetFields();
          onClose();
        }
      }}
      initialValues={{
        isPrivate: false,
      }}
      onFinish={handleSubmit}
      submitter={{
        searchConfig: {
          submitText: t('common:button.submit'),
          resetText: t('common:button.cancel'),
        },
        submitButtonProps: {
          loading: confirmLoading || createMutation.isPending || updateMutation.isPending,
        },
        resetButtonProps: { onClick: onClose },
      }}
      drawerProps={{ destroyOnHidden: true, onClose, placement: 'left', size: 800 }}
    >
      <ProFormText
        name="filename"
        label={t('filename')}
        placeholder={t('filenamePlaceholder')}
        rules={[{ required: true, message: t('requiredFilename') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormSelect
        name="type"
        label={t('type')}
        placeholder={t('typePlaceholder')}
        options={assetTypeOptions(t)}
        fieldProps={{ allowClear: true, showSearch: true, optionFilterProp: 'label' }}
      />

      <ProFormText
        name="title"
        label={t('title')}
        placeholder={t('titlePlaceholder')}
        fieldProps={{ allowClear: true }}
      />

      <ProFormTextArea
        name="caption"
        label={t('caption')}
        placeholder={t('captionPlaceholder')}
        fieldProps={{ rows: 2, allowClear: true }}
      />

      <ProFormText
        name="altText"
        label={t('altText')}
        placeholder={t('altTextPlaceholder')}
        fieldProps={{ allowClear: true }}
      />

      <ProFormSwitch name="isPrivate" label={t('isPrivate')} />
    </DrawerForm>
  );
};

export default MediaAssetDrawer;
