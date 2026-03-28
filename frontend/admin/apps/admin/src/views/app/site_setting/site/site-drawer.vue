<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { siteStatusList, useLanguageStore, useSiteStore } from '#/stores';

const siteStore = useSiteStore();
const languageStore = useLanguageStore();

const data = ref();
const languageOptions = ref<{ label: string; value: string }[]>([]);

const getTitle = computed(() =>
  data.value?.create
    ? $t('ui.modal.create', { moduleName: $t('page.site.moduleName') })
    : $t('ui.modal.update', { moduleName: $t('page.site.moduleName') }),
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
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
  },
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.site.name'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'slug',
      label: $t('page.site.slug'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'domain',
      label: $t('page.site.domain'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('page.site.status'),
      rules: 'selectRequired',
      defaultValue: 'SITE_STATUS_ACTIVE',
      componentProps: {
        options: siteStatusList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
    {
      component: 'Switch',
      fieldName: 'isDefault',
      label: $t('page.site.isDefault'),
      help: $t('page.site.help.isDefault'),
      componentProps: {
        class: 'w-auto',
      },
    },
    {
      component: 'Select',
      fieldName: 'defaultLocale',
      label: $t('page.site.defaultLocale'),
      defaultValue: 'zh-CN',
      componentProps: {
        options: languageOptions,
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
      },
      help: $t('page.site.help.defaultLocale'),
    },
    {
      component: 'Input',
      fieldName: 'template',
      label: $t('page.site.template'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'theme',
      label: $t('page.site.theme'),
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
    const isCreate = data.value?.create;

    (async () => {
      try {
        await (isCreate
          ? siteStore.createSite(values)
          : siteStore.updateSite(data.value.row.id, values));

        notification.success({
          message: isCreate
            ? $t('ui.notification.create_success')
            : $t('ui.notification.update_success'),
        });
        drawerApi.close();
      } catch {
        notification.error({
          message: isCreate
            ? $t('ui.notification.create_failed')
            : $t('ui.notification.update_failed'),
        });
      } finally {
        setLoading(false);
      }
    })();
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
