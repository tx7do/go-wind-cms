import { useEffect, useMemo, useRef, useState } from 'react';
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
import { PLATFORM_CREDENTIAL_FIELDS, platformOptions } from './constants';

interface SocialAccountDrawerProps {
  open: boolean;
  data?: any;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 社交账号抽屉（新建 = 连接平台账号；编辑 = 改名称/启停/重输凭据）。
 *
 * 凭据模型：按所选平台动态渲染凭据字段，组装成 credentials JSON 提交；
 * 后端保存时同步调用平台验证，坏凭据直接报错不入库；编辑路径只提交非空键，
 * 未填键由服务端从已存凭据保留（增量合并，不会丢密钥）。凭据原文永远不可回读
 * （只有脱敏 hint），所以编辑时密码框留空 = 不修改凭据。
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
  const [platform, setPlatform] = useState<string | undefined>();

  const isEdit = !!data?.id;

  useEffect(() => {
    if (open) {
      setTimeout(() => {
        setPlatform(data?.platform);
        formRef.current?.setFieldsValue({
          platform: data?.platform,
          name: data?.name || '',
          handle: data?.handle || '',
          isEnabled: data ? data.isEnabled !== false : true,
          remark: data?.remark || '',
        });
      }, 0);
    }
  }, [open, data]);

  const credFields = useMemo(
    () => (platform ? PLATFORM_CREDENTIAL_FIELDS[platform] ?? [] : []),
    [platform],
  );

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

  // 凭据字段表单名前缀，避免与账号列字段撞名
  const credFormKey = (key: string) => `cred_${key}`;

  const handleSubmit = async (values: Record<string, any>) => {
    try {
      setConfirmLoading(true);

      const fields = PLATFORM_CREDENTIAL_FIELDS[values.platform] ?? [];
      // 只收集非空凭据键（编辑路径留空 = 保留已存值，服务端增量合并）
      const credentialValues: Record<string, string> = {};
      for (const f of fields) {
        const v = values[credFormKey(f.key)];
        if (v !== undefined && v !== null && String(v).trim() !== '') {
          credentialValues[f.key] = String(v).trim();
        }
      }

      // alsoHandle 平台：凭据键的值同时作为 handle 展示列
      const alsoHandleKey = fields.find((f) => f.alsoHandle)?.key;

      if (isEdit) {
        const maskKeys: string[] = ['name', 'isEnabled', 'remark'];
        const payload: Record<string, any> = {
          name: values.name,
          isEnabled: values.isEnabled ?? true,
          remark: values.remark,
        };
        const nextHandle =
          (alsoHandleKey && credentialValues[alsoHandleKey]) || values.handle;
        if (nextHandle && nextHandle !== data.handle) {
          payload.handle = nextHandle;
          maskKeys.push('handle');
        }
        if (Object.keys(credentialValues).length !== 0) {
          payload.credentials = JSON.stringify(credentialValues);
          maskKeys.push('credentials');
        }
        await updateMutation.mutateAsync({ id: data.id, values: payload, maskKeys });
      } else {
        // 新建：凭据必填字段由表单规则强校验后整包提交
        await createMutation.mutateAsync({
          data: {
            platform: values.platform,
            name: values.name,
            handle: (alsoHandleKey && credentialValues[alsoHandleKey]) || values.handle,
            isEnabled: values.isEnabled ?? true,
            remark: values.remark,
            credentials: JSON.stringify(credentialValues),
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
        fieldProps={{
          allowClear: false,
          showSearch: true,
          optionFilterProp: 'label',
          onChange: (value: string) => {
            setPlatform(value);
            // 切平台后清空所有凭据键的残留值
            Object.values(PLATFORM_CREDENTIAL_FIELDS)
              .flat()
              .forEach((f) => formRef.current?.setFieldValue(credFormKey(f.key), undefined));
          },
        }}
      />

      <ProFormText
        name="name"
        label={t('accountName')}
        placeholder={t('accountNamePlaceholder')}
        rules={[{ required: true, message: t('requiredAccountName') }]}
        fieldProps={{ allowClear: true }}
      />

      {/* handle 展示列：Bluesky 的登录标识即 handle，单独字段；其余平台由凭据键回填 */}
      {platform === 'PLATFORM_BLUESKY' && (
        <ProFormText
          name="handle"
          label={t('handle')}
          placeholder={t('handlePlaceholder')}
          rules={[{ required: true, message: t('requiredHandle') }]}
          extra={t('helpHandle')}
          fieldProps={{ allowClear: true }}
        />
      )}

      {/* 按平台动态渲染凭据字段；编辑模式下留空 = 不修改该键 */}
      {credFields.map((f) => {
        const name = credFormKey(f.key);
        const label = t(`credentialFields.${f.key}`);
        const helpKey = `helpCred.${f.key}`;
        const help = t(helpKey, { defaultValue: '' });
        const commonProps = {
          name,
          key: name,
          label,
          extra: help !== helpKey ? help : undefined,
        };
        if (f.secret) {
          return (
            <ProFormText.Password
              {...commonProps}
              placeholder={isEdit ? t('credentialEditPlaceholder') : t('credentialPlaceholder')}
              rules={isEdit ? [] : [{ required: f.required, message: t('requiredCredential', { label }) }]}
              fieldProps={{ allowClear: true, autoComplete: 'new-password' }}
            />
          );
        }
        return (
          <ProFormText
            {...commonProps}
            placeholder={t('credentialPlaceholder')}
            rules={isEdit ? [] : [{ required: f.required, message: t('requiredCredential', { label }) }]}
            fieldProps={{ allowClear: true }}
          />
        );
      })}

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
