<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import {
  siteSettingTypeList,
  useLanguageStore,
  useSiteSettingStore,
} from '#/stores';

const siteSettingStore = useSiteSettingStore();
const languageStore = useLanguageStore();

const data = ref();
const languageOptions = ref<{ label: string; value: string }[]>([]);

const getTitle = computed(() =>
  data.value?.create
    ? $t('ui.modal.create', { moduleName: $t('page.siteSetting.moduleName') })
    : $t('ui.modal.update', { moduleName: $t('page.siteSetting.moduleName') }),
);

onMounted(async () => {
  try {
    const resp = await languageStore.listLanguage(undefined, {}, undefined, [
      'sortOrder',
    ]);
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

    try {
      await (data.value?.create
        ? siteSettingStore.createSiteSetting(values)
        : siteSettingStore.updateSiteSetting(data.value.row.id, values));

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
      baseFormApi.setValues(data.value?.row);

      setLoading(false);
    }
  },
});

function setLoading(loading: boolean) {
  drawerApi.setState({ loading });
}
</script>

<template>
  <Drawer :title="getTitle">
    <BaseForm />
  </Drawer>
</template>
