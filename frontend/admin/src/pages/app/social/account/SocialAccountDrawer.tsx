import { useEffect, useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  DrawerForm,
  ProFormSelect,
  ProFormSwitch,
  ProFormText,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { App } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';

import { useCreateSocialAccount, useUpdateSocialAccount } from '@/api/hooks/social';
import { platformOptions } from './constants';

interface SocialAccountDrawerProps {
  open: boolean;
  data?: any;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 社交账号抽屉（新建 = 连接平台账号；编辑 = 改名称/启停/重输凭据）。
 *
 * 凭据模型：表单收集 handle + App Password，组装成 credentials JSON 提交；
 * 后端保存时同步调用平台验证，坏凭据直接报错不入库。读取路径永远拿不到
 * 凭据原文（只有脱敏 hint），所以编辑时密码框留空 = 不修改凭据。
 */
const SocialAccountDrawer: React.FC<SocialAccountDrawerProps> = ({
  open,
  data,
  onClose,
  onSuccess,
}) => {
  const { t } = useTranslation('social');
  const formRef = useRef<ProFormInstance>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [confirmLoading, setConfirmLoading] = useState(false);

  const isEdit = !!data?.id;

  useEffect(() => {
    if (open) {
      setTimeout(() => {
        formRef.current?.setFieldsValue({
          platform: data?.platform,
          name: data?.name || '',
          handle: data?.handle || '',
          appPassword: '',
          isEnabled: data ? data.isEnabled !== false : true,
          remark: data?.remark || '',
        });
      }, 0);
    }
  }, [open, data]);

  const createMutation = useCreateSocialAccount({
    onSuccess: () => {
      message.success(t('createSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listSocialAccounts'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('createFailed'));
    },
  });

  const updateMutation = useUpdateSocialAccount({
    onSuccess: () => {
      message.success(t('updateSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listSocialAccounts'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('updateFailed'));
    },
  });

  const handleSubmit = async (values: Record<string, any>) => {
    try {
      setConfirmLoading(true);

      if (isEdit) {
        // 编辑：凭据只在重输密码时提交（空 = 不修改）
        const maskKeys: string[] = ['name', 'isEnabled', 'remark'];
        const payload: Record<string, any> = {
          name: values.name,
          isEnabled: values.isEnabled ?? true,
          remark: values.remark,
        };
        if (values.handle && values.handle !== data.handle) {
          payload.handle = values.handle;
          maskKeys.push('handle');
        }
        if (values.appPassword) {
          payload.credentials = JSON.stringify({
            handle: values.handle,
            app_password: values.appPassword,
          });
          maskKeys.push('credentials');
        }
        await updateMutation.mutateAsync({
          id: data.id,
          values: payload,
          maskKeys,
        });
      } else {
        // 新建：凭据必填，保存即验证
        await createMutation.mutateAsync({
          data: {
            platform: values.platform,
            name: values.name,
            handle: values.handle,
            isEnabled: values.isEnabled ?? true,
            remark: values.remark,
            credentials: JSON.stringify({
              handle: values.handle,
              app_password: values.appPassword,
            }),
          } as any,
        });
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
      title={isEdit ? t('editAccount') : t('createAccount')}
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
          loading: confirmLoading || createMutation.isPending || updateMutation.isPending,
        },
        resetButtonProps: { onClick: onClose },
      }}
      drawerProps={{ destroyOnHidden: true, onClose, size: 600 }}
    >
      <ProFormSelect
        name="platform"
        label={t('platform')}
        placeholder={t('platformPlaceholder')}
        rules={[{ required: true, message: t('requiredPlatform') }]}
        options={platformOptions(t)}
        disabled={isEdit}
        extra={t('helpPlatform')}
        fieldProps={{ allowClear: false, showSearch: true, optionFilterProp: 'label' }}
      />

      <ProFormText
        name="name"
        label={t('accountName')}
        placeholder={t('accountNamePlaceholder')}
        rules={[{ required: true, message: t('requiredAccountName') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="handle"
        label={t('handle')}
        placeholder={t('handlePlaceholder')}
        rules={[{ required: true, message: t('requiredHandle') }]}
        extra={t('helpHandle')}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText.Password
        name="appPassword"
        label={t('appPassword')}
        placeholder={isEdit ? t('appPasswordEditPlaceholder') : t('appPasswordPlaceholder')}
        rules={isEdit ? [] : [{ required: true, message: t('requiredAppPassword') }]}
        extra={
          isEdit
            ? t('helpAppPasswordEdit')
            : t('helpAppPasswordCreate')
        }
        fieldProps={{ allowClear: true, autoComplete: 'new-password' }}
      />

      <ProFormSwitch
        name="isEnabled"
        label={t('isEnabled')}
        help={t('helpIsEnabled')}
      />

      <ProFormTextArea
        name="remark"
        label={t('remark')}
        fieldProps={{ rows: 3, allowClear: true }}
      />
    </DrawerForm>
  );
};

export default SocialAccountDrawer;
