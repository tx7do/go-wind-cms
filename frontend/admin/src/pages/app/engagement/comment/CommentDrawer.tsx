import { useEffect, useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  DrawerForm,
  ProFormSelect,
  ProFormSwitch,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { App } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';

import { useUpdateComment } from '@/api/hooks/comment';
import { statusOptions } from './constants';

interface CommentDrawerProps {
  open: boolean;
  data?: any;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 评论审核抽屉（仅用于审核：状态/垃圾/置顶；
 * 评论原文不在审核路径中修改）
 */
const CommentDrawer: React.FC<CommentDrawerProps> = ({
  open,
  data,
  onClose,
  onSuccess,
}) => {
  const { t } = useTranslation('comment');
  const formRef = useRef<ProFormInstance>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [confirmLoading, setConfirmLoading] = useState(false);

  // 编辑模式下设置表单值（destroyOnHidden 时需延迟赋值）
  useEffect(() => {
    if (open && data) {
      setTimeout(() => {
        formRef.current?.setFieldsValue({
          content: data.content || '',
          status: data.status,
          isSpam: data.isSpam ?? false,
          isSticky: data.isSticky ?? false,
        });
      }, 0);
    }
  }, [open, data]);

  // 更新 mutation
  const updateMutation = useUpdateComment({
    onSuccess: () => {
      message.success(t('updateSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listComments'] });
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
      if (!data?.id) return false;

      await updateMutation.mutateAsync({
        id: data.id,
        values: {
          status: values.status,
          isSpam: values.isSpam ?? false,
          isSticky: values.isSticky ?? false,
        },
        // content（评论原文）不在审核路径中修改，排除以避免误改用户原文
        maskKeys: ['status', 'isSpam', 'isSticky'],
      });
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
      title={t('edit')}
      open={open}
      onOpenChange={(visible) => {
        if (!visible) {
          formRef.current?.resetFields();
          onClose();
        }
      }}
      onFinish={handleSubmit}
      submitter={{
        searchConfig: {
          submitText: t('common:button.submit'),
          resetText: t('common:button.cancel'),
        },
        submitButtonProps: {
          loading: confirmLoading || updateMutation.isPending,
        },
        resetButtonProps: { onClick: onClose },
      }}
      drawerProps={{ destroyOnHidden: true, onClose, placement: 'left', size: 800 }}
    >
      <ProFormTextArea
        name="content"
        label={t('content')}
        rules={[{ required: true, message: t('requiredContent') }]}
        fieldProps={{ rows: 4, allowClear: true, disabled: true }}
      />

      <ProFormSelect
        name="status"
        label={t('status')}
        placeholder={t('statusPlaceholder')}
        rules={[{ required: true, message: t('requiredStatus') }]}
        options={statusOptions(t)}
        fieldProps={{ allowClear: true, showSearch: true, optionFilterProp: 'label' }}
      />

      <ProFormSwitch
        name="isSpam"
        label={t('isSpam')}
        help={t('helpIsSpam')}
      />

      <ProFormSwitch
        name="isSticky"
        label={t('isSticky')}
        help={t('helpIsSticky')}
      />
    </DrawerForm>
  );
};

export default CommentDrawer;
