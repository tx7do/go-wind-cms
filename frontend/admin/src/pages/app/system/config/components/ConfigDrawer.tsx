import { useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  DrawerForm,
  ProFormText,
  ProFormSelect,
  ProFormSwitch,
} from '@ant-design/pro-components';
import { App } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { configservicev1_Config as Config } from '@/api/generated/admin/service/v1';
import { useCreateConfig, useUpdateConfig } from '@/api/hooks/config';
import { getValueTypeOptions } from '../constants';

interface ConfigDrawerProps {
  open: boolean;
  mode: 'create' | 'edit';
  data?: Config;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 参数编辑/创建抽屉组件
 */
const ConfigDrawer: React.FC<ConfigDrawerProps> = ({
  open,
  mode,
  data,
  onClose,
  onSuccess,
}) => {
  const { t } = useTranslation('config');
  const formRef = useRef<ProFormInstance>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [confirmLoading, setConfirmLoading] = useState(false);

  // 创建参数
  const createMutation = useCreateConfig({
    onSuccess: () => {
      message.success(t('createSuccess'));
      onSuccess();
      onClose();
      queryClient.invalidateQueries({ queryKey: ['listConfigs'] });
    },
    onError: (error: Error) => {
      console.error('create config failed', error);
      message.error(error.message || t('createFailed'));
    },
  });

  // 更新参数
  const updateMutation = useUpdateConfig({
    onSuccess: () => {
      message.success(t('updateSuccess'));
      onSuccess();
      onClose();
      queryClient.invalidateQueries({ queryKey: ['listConfigs'] });
    },
    onError: (error: Error) => {
      console.error('update config failed', error);
      message.error(error.message || t('updateFailed'));
    },
  });

  // 提交表单
  const handleSubmit = async (values: any) => {
    setConfirmLoading(true);

    try {
      const payload = {
        ...values,
        valueType: values.valueType ?? 'STRING',
        isBuiltIn: values.isBuiltIn ?? false,
      };

      if (mode === 'create') {
        await createMutation.mutateAsync({ data: payload });
      } else if (data?.id) {
        await updateMutation.mutateAsync({ id: data.id, values: payload });
      }
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
      initialValues={
        mode === 'edit'
          ? { ...data }
          : {
              valueType: 'STRING',
              isBuiltIn: false,
            }
      }
      onFinish={handleSubmit}
      submitter={{
        searchConfig: {
          submitText: t('common:button.submit'),
          resetText: t('common:button.cancel'),
        },
        submitButtonProps: {
          loading: confirmLoading || createMutation.isPending || updateMutation.isPending,
        },
        resetButtonProps: {
          onClick: onClose,
        },
      }}
      drawerProps={{
        destroyOnHidden: true,
        onClose,
        size: 480,
      }}
    >
      <ProFormText
        name="name"
        label={t('name')}
        placeholder={t('namePlaceholder')}
        rules={[{ required: true, message: t('requiredName') }]}
        fieldProps={{
          allowClear: true,
        }}
      />

      <ProFormText
        name="key"
        label={t('key')}
        placeholder={t('keyPlaceholder')}
        rules={[{ required: true, message: t('requiredKey') }]}
        fieldProps={{
          allowClear: true,
        }}
        extra="e.g. sys.login.captchaEnabled"
      />

      <ProFormText
        name="value"
        label={t('value')}
        placeholder={t('valuePlaceholder')}
        rules={[{ required: true, message: t('requiredValue') }]}
        fieldProps={{
          allowClear: true,
        }}
      />

      <ProFormSelect
        name="valueType"
        label={t('valueType')}
        options={getValueTypeOptions(t)}
        fieldProps={{
          allowClear: false,
        }}
      />

      <ProFormSwitch
        name="isBuiltIn"
        label={t('isBuiltIn')}
      />
    </DrawerForm>
  );
};

export default ConfigDrawer;
