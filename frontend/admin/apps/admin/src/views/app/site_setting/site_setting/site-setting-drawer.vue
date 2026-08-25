<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import {
  apiClient,
  fetchListLanguages,
  makeUpdateMask,
  PaginationQuery,
  siteSettingTypeList,
} from '#/api';

const data = ref();
const languageOptions = ref<{ label: string; value: string }[]>([]);

const getTitle = computed(() =>
  data.value?.create
    ? $t('ui.modal.create', { moduleName: $t('page.siteSetting.moduleName') })
    : $t('ui.modal.update', { moduleName: $t('page.siteSetting.moduleName') }),
);

onMounted(async () => {
  try {
    const resp = await fetchListLanguages(
      new PaginationQuery({ orderBy: ['sortOrder'] }),
    );
    languageOptions.value =
      resp.items?.map((lang) => ({
        label: lang.nativeName || lang.languageCode || '',
        value: lang.languageCode || '',
      })) || [];
  } catch (error) {
    console.error('Failed to load language list:', error);
  }
});

const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,
  // 所有表单项共用，可单独在表单内覆盖
  commonConfig: {
    // 所有表单项
    componentProps: {
      class: 'w-full',
    },
  },
  schema: [
    {
      component: 'InputNumber',
      fieldName: 'siteId',
      label: $t('page.siteSetting.siteId'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
        defaultValue: 1,
      },
      rules: 'required',
    },
    {
      component: 'Select',
      fieldName: 'type',
      label: $t('page.siteSetting.type'),
      rules: 'selectRequired',
      componentProps: {
        options: siteSettingTypeList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
    {
      component: 'Input',
      fieldName: 'key',
      label: $t('page.siteSetting.key'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'label',
      label: $t('page.siteSetting.label'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'placeholder',
      label: $t('page.siteSetting.placeholder'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Input',
      fieldName: 'group',
      label: $t('page.siteSetting.group'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'locale',
      label: $t('page.siteSetting.locale'),
      defaultValue: 'zh-CN',
      componentProps: {
        options: languageOptions,
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
      },
      rules: 'selectRequired',
    },
    {
      component: 'Input',
      fieldName: 'validationRegex',
      label: $t('page.siteSetting.validationRegex'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Textarea',
      fieldName: 'description',
      label: $t('page.siteSetting.description'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Switch',
      fieldName: 'isRequired',
      label: $t('page.siteSetting.isRequired'),
      componentProps: {
        class: 'w-auto',
      },
    },
    {
      component: 'Textarea',
      fieldName: 'optionsJson',
      label: $t('page.siteSetting.options'),
      componentProps: {
        placeholder: $t('page.siteSetting.options'),
        allowClear: true,
        class: 'w-full',
      },
      dependencies: {
        triggerFields: ['type'],
        if: (values) => values.type === 'SETTING_TYPE_SELECT',
      },
    },
    {
      component: 'Input',
      fieldName: 'value',
      label: $t('page.siteSetting.value'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      dependencies: {
        triggerFields: ['type'],
        if: (values) => values.type === 'SETTING_TYPE_TEXT',
      },
    },
    {
      component: 'Input',
      fieldName: 'value',
      label: $t('page.siteSetting.value'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      dependencies: {
        triggerFields: ['type'],
        if: (values) => values.type === 'SETTING_TYPE_URL',
      },
    },
    {
      component: 'Input',
      fieldName: 'value',
      label: $t('page.siteSetting.value'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      dependencies: {
        triggerFields: ['type'],
        if: (values) => values.type === 'SETTING_TYPE_EMAIL',
      },
    },
    {
      component: 'Textarea',
      fieldName: 'value',
      label: $t('page.siteSetting.value'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      dependencies: {
        triggerFields: ['type'],
        if: (values) => values.type === 'SETTING_TYPE_TEXTAREA',
      },
    },
    {
      component: 'Textarea',
      fieldName: 'value',
      label: $t('page.siteSetting.value'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      dependencies: {
        triggerFields: ['type'],
        if: (values) => values.type === 'SETTING_TYPE_JSON',
      },
    },
    {
      component: 'Input',
      fieldName: 'value',
      label: $t('page.siteSetting.value'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      dependencies: {
        triggerFields: ['type'],
        if: (values) => values.type === 'SETTING_TYPE_IMAGE',
      },
    },
    {
      component: 'InputNumber',
      fieldName: 'value',
      label: $t('page.siteSetting.value'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
        class: 'w-full',
      },
      dependencies: {
        triggerFields: ['type'],
        if: (values) => values.type === 'SETTING_TYPE_NUMBER',
      },
    },
    {
      component: 'Switch',
      fieldName: 'value',
      label: $t('page.siteSetting.value'),
      componentProps: {
        class: 'w-auto',
      },
      dependencies: {
        triggerFields: ['type'],
        if: (values) => values.type === 'SETTING_TYPE_BOOLEAN',
      },
    },
    {
      component: 'Select',
      fieldName: 'value',
      label: $t('page.siteSetting.value'),
      dependencies: {
        triggerFields: ['type', 'optionsJson'],
        if: (values) => values.type === 'SETTING_TYPE_SELECT',
        componentProps: (values) => {
          let options: { label: string; value: string }[] = [];
          try {
            const parsed = values.optionsJson
              ? (JSON.parse(values.optionsJson as string) as unknown)
              : null;
            if (parsed && typeof parsed === 'object') {
              options = Object.entries(parsed as Record<string, string>).map(
                ([v, label]) => ({ label: String(label), value: String(v) }),
              );
            }
          } catch {
            options = [];
          }
          return {
            options,
            placeholder: $t('ui.placeholder.select'),
            allowClear: true,
            class: 'w-full',
          };
        },
      },
    },
  ],
});

const [Drawer, drawerApi] = useVbenDrawer({
  onCancel() {
    drawerApi.close();
  },

  async onConfirm() {
    const validate = await baseFormApi.validate();
    if (!validate.valid) {
      return;
    }

    setLoading(true);
    const values = await baseFormApi.getValues();

    // optionsJson 是 options 的字符串桥接（前端编辑用），提交前还原回 options 对象。
    if ('optionsJson' in values) {
      const options = parseOptionsJson(values.optionsJson as string);
      (values as Record<string, any>).options = options;
      delete (values as Record<string, any>).optionsJson;
    }

    try {
      await (data.value?.create
        ? apiClient.siteSettingService.Create({ data: { ...values } as any })
        : apiClient.siteSettingService.Update({
            id: data.value.row.id,
            data: { ...values } as any,
            updateMask: makeUpdateMask(
              Object.keys(values).filter(
                (k) => !['key', 'locale', 'siteId'].includes(k),
              ),
            ),
          }));

      notification.success({
        message: data.value?.create
          ? $t('ui.notification.create_success')
          : $t('ui.notification.update_success'),
      });
    } catch {
      notification.error({
        message: data.value?.create
          ? $t('ui.notification.create_failed')
          : $t('ui.notification.update_failed'),
      });
    } finally {
      drawerApi.close();
      setLoading(false);
    }
  },

  onOpenChange(isOpen) {
    if (isOpen) {
      data.value = drawerApi.getData<Record<string, any>>();
      const row = data.value?.row;
      if (row) {
        // 记录的 options 对象序列化为 JSON 文本，供 optionsJson 字段编辑。
        row.optionsJson = serializeOptions(
          row.options as Record<string, string> | undefined,
        );
      }
      baseFormApi.setValues(row);

      setLoading(false);
    }
  },
});

function setLoading(loading: boolean) {
  drawerApi.setState({ confirmLoading: loading });
}

// options 是后端 map<string,string>；前端用 optionsJson（JSON 文本）桥接编辑。
// 非对象/解析失败一律回退到空对象，避免脏数据落库。
function serializeOptions(options: Record<string, string> | undefined): string {
  try {
    return options ? JSON.stringify(options) : '{}';
  } catch {
    return '{}';
  }
}

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
</script>

<template>
  <Drawer :title="getTitle">
    <BaseForm />
  </Drawer>
</template>
