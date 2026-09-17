import { useRef, useState, useEffect } from 'react';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  DrawerForm,
  ProFormText,
  ProFormDigit,
  ProFormSelect,
  ProFormSwitch,
  ProFormTextArea,
  ProFormDependency,
} from '@ant-design/pro-components';
import { App } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';

import type { siteservicev1_SiteSetting } from '@/api/generated/admin/service/v1';
import {
  useCreateSiteSetting,
  useUpdateSiteSetting,
} from '@/api/hooks/site-setting';
import { fetchListLanguages } from '@/api/hooks/language';
import { PaginationQuery } from '@/core';
import { typeOptions } from './constants';

interface SiteSettingDrawerProps {
  open: boolean;
  mode: 'create' | 'edit';
  data?: any;
  onClose: () => void;
  onSuccess: () => void;
}

// 更新时不可变的身份字段：进 data 但不进 updateMask（对齐 vben 侧）
const IDENTITY_FIELDS = ['key', 'locale', 'siteId'];

// options 是后端 map<string,string>；前端用 optionsJson（JSON 文本）桥接编辑。
// 非对象/解析失败一律回退到空对象，避免脏数据落库。
function parseOptionsJson(json: string): Record<string, string> {
  try {
    const parsed = json ? (JSON.parse(json) as unknown) : null;
    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
      const result: Record<string, string> = {};
      for (const [k, v] of Object.entries(parsed as Record<string, unknown>)) {
        if (typeof v === 'string') result[k] = v;
      }
      return result;
    }
  } catch {
    // 忽略非法 JSON
  }
  return {};
}

function parseSelectOptions(optionsJson: string): { label: string; value: string }[] {
  const parsed = parseOptionsJson(optionsJson);
  return Object.entries(parsed).map(([value, label]) => ({ label, value }));
}

/**
 * 站点配置编辑/创建抽屉组件
 */
const SiteSettingDrawer: React.FC<SiteSettingDrawerProps> = ({
  open,
  mode,
  data,
  onClose,
  onSuccess,
}) => {
  const { t } = useTranslation('site-setting');
  const formRef = useRef<ProFormInstance>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [confirmLoading, setConfirmLoading] = useState(false);
  const [languageOptions, setLanguageOptions] = useState<
    { label: string; value: string }[]
  >([]);

  useEffect(() => {
    if (!open) return;
    (async () => {
      try {
        const resp = await fetchListLanguages(
          new PaginationQuery({ orderBy: ['sortOrder'] }),
        );
        setLanguageOptions(
          resp.items?.map((lang: any) => ({
            label: lang.nativeName || lang.languageCode || '',
            value: lang.languageCode || '',
          })) || [],
        );
      } catch (error) {
        console.error('Failed to load language list:', error);
      }
    })();
  }, [open]);

  // 编辑模式下设置表单值（destroyOnHidden 时需延迟赋值）。
  // value 在后端是 string，按类型还原为 number/boolean 供对应控件回显。
  useEffect(() => {
    if (open && mode === 'edit' && data) {
      setTimeout(() => {
        let value: any = data.value;
        if (data.type === 'SETTING_TYPE_NUMBER') {
          value = data.value === '' || data.value == null ? undefined : Number(data.value);
        } else if (data.type === 'SETTING_TYPE_BOOLEAN') {
          value = String(data.value) === 'true';
        }
        formRef.current?.setFieldsValue({
          siteId: data.siteId ?? 1,
          type: data.type,
          key: data.key || '',
          label: data.label || '',
          placeholder: data.placeholder || '',
          group: data.group || '',
          locale: data.locale || 'zh-CN',
          validationRegex: data.validationRegex || '',
          description: data.description || '',
          isRequired: data.isRequired ?? false,
          value,
          // 桥接回填：options map<string,string> → JSON 文本
          optionsJson: data.options ? JSON.stringify(data.options) : '',
        });
      }, 0);
    }
  }, [open, mode, data]);

  // 创建 mutation
  const createMutation = useCreateSiteSetting({
    onSuccess: () => {
      message.success(t('createSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listSiteSettings'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('createFailed'));
    },
  });

  // 更新 mutation
  const updateMutation = useUpdateSiteSetting({
    onSuccess: () => {
      message.success(t('updateSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listSiteSettings'] });
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

      // optionsJson 是 options 的字符串桥接（前端编辑用），提交前还原回 options 对象
      const payload: Record<string, any> = {
        siteId: values.siteId,
        type: values.type,
        key: values.key,
        label: values.label,
        placeholder: values.placeholder,
        group: values.group,
        locale: values.locale,
        validationRegex: values.validationRegex,
        description: values.description,
        isRequired: values.isRequired ?? false,
        options: values.type === 'SETTING_TYPE_SELECT'
          ? parseOptionsJson(values.optionsJson ?? '')
          : {},
      };

      // value 按类型归一为 string（后端 proto 为 string；BOOLEAN 用 'true'/'false'）
      switch (values.type) {
        case 'SETTING_TYPE_BOOLEAN':
          payload.value = values.value ? 'true' : 'false';
          break;
        case 'SETTING_TYPE_NUMBER':
          payload.value = values.value == null ? '' : String(values.value);
          break;
        default:
          payload.value = values.value ?? '';
      }

      if (mode === 'edit' && data?.id) {
        await updateMutation.mutateAsync({
          id: data.id,
          values: payload as unknown as siteservicev1_SiteSetting,
          maskKeys: Object.keys(payload).filter(
            (k) => !IDENTITY_FIELDS.includes(k),
          ),
        });
      } else {
        await createMutation.mutateAsync({ data: payload as unknown as siteservicev1_SiteSetting });
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
        siteId: 1,
        locale: 'zh-CN',
        isRequired: false,
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
      <ProFormDigit
        name="siteId"
        label={t('siteId')}
        placeholder={t('siteIdPlaceholder')}
        rules={[{ required: true, message: t('requiredSiteId') }]}
        fieldProps={{ min: 1, precision: 0 }}
      />

      <ProFormSelect
        name="type"
        label={t('type')}
        placeholder={t('typePlaceholder')}
        rules={[{ required: true, message: t('requiredType') }]}
        options={typeOptions(t)}
        fieldProps={{ allowClear: true, showSearch: true, optionFilterProp: 'label' }}
      />

      <ProFormText
        name="key"
        label={t('key')}
        placeholder={t('keyPlaceholder')}
        rules={[{ required: true, message: t('requiredKey') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="label"
        label={t('label')}
        placeholder={t('labelPlaceholder')}
        rules={[{ required: true, message: t('requiredLabel') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="placeholder"
        label={t('placeholderLabel')}
        placeholder={t('placeholderFieldPlaceholder')}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="group"
        label={t('group')}
        placeholder={t('groupPlaceholder')}
        fieldProps={{ allowClear: true }}
      />

      <ProFormSelect
        name="locale"
        label={t('locale')}
        placeholder={t('localePlaceholder')}
        rules={[{ required: true, message: t('requiredLocale') }]}
        options={languageOptions}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="validationRegex"
        label={t('validationRegex')}
        placeholder={t('validationRegexPlaceholder')}
        fieldProps={{ allowClear: true }}
      />

      <ProFormTextArea
        name="description"
        label={t('description')}
        placeholder={t('descriptionPlaceholder')}
        fieldProps={{ rows: 2, allowClear: true }}
      />

      <ProFormSwitch
        name="isRequired"
        label={t('isRequired')}
        help={t('helpIsRequired')}
      />

      {/* value 按类型切换控件形态（对齐 vben 的 9 分支条件渲染） */}
      <ProFormDependency name={['type', 'optionsJson']}>
        {({ type, optionsJson }) => {
          if (!type) return null;
          switch (type) {
            case 'SETTING_TYPE_BOOLEAN':
              return <ProFormSwitch name="value" label={t('value')} />;
            case 'SETTING_TYPE_NUMBER':
              return (
                <ProFormDigit
                  name="value"
                  label={t('value')}
                  fieldProps={{ precision: 0 }}
                />
              );
            case 'SETTING_TYPE_TEXTAREA':
            case 'SETTING_TYPE_JSON':
              return (
                <ProFormTextArea
                  name="value"
                  label={t('value')}
                  placeholder={t('valuePlaceholder')}
                  fieldProps={{ rows: 4, allowClear: true }}
                />
              );
            case 'SETTING_TYPE_SELECT': {
              const options = parseSelectOptions(optionsJson ?? '');
              return (
                <ProFormSelect
                  name="value"
                  label={t('value')}
                  options={options}
                  help={t('helpSelectValue')}
                  fieldProps={{ allowClear: true }}
                />
              );
            }
            default:
              return (
                <ProFormText
                  name="value"
                  label={t('value')}
                  placeholder={t('valuePlaceholder')}
                  fieldProps={{ allowClear: true }}
                />
              );
          }
        }}
      </ProFormDependency>

      <ProFormDependency name={['type']}>
        {({ type }) =>
          type === 'SETTING_TYPE_SELECT' ? (
            <ProFormTextArea
              name="optionsJson"
              label={t('options')}
              placeholder={t('optionsPlaceholder')}
              help={t('helpOptions')}
              fieldProps={{ rows: 3, allowClear: true }}
            />
          ) : null
        }
      </ProFormDependency>
    </DrawerForm>
  );
};

export default SiteSettingDrawer;

