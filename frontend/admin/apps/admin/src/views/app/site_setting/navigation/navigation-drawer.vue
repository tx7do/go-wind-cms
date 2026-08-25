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
  navigationLocationList,
  PaginationQuery,
} from '#/api';

const data = ref<Record<string, any>>();
const languageOptions = ref<{ label: string; value: string }[]>([]);

const getTitle = computed(() =>
  data.value?.create
    ? $t('ui.modal.create', { moduleName: $t('page.navigation.moduleName') })
    : $t('ui.modal.update', { moduleName: $t('page.navigation.moduleName') }),
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
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
  },
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.navigation.name'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Select',
      fieldName: 'location',
      label: $t('page.navigation.location'),
      componentProps: {
        options: navigationLocationList,
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
      },
      rules: 'selectRequired',
    },
    {
      component: 'Select',
      fieldName: 'locale',
      label: $t('page.navigation.locale'),
      defaultValue: 'zh-CN',
      componentProps: {
        options: languageOptions,
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
      },
      rules: 'selectRequired',
    },
    {
      component: 'Switch',
      fieldName: 'isActive',
      label: $t('page.navigation.isActive'),
      defaultValue: true,
      componentProps: {
        class: 'w-auto',
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

    // 导航项（items）由独立的 navigation-item-list 子组件管理 CRUD，
    // 此 drawer 仅处理导航自身字段；payload 不含 items，避免整体替换语义
    // 误清空子项（itemsJson 往返解析失败会 items=[] 触发 CleanItems）。
    const payload = {
      isActive: values.isActive,
      locale: values.locale,
      location: values.location,
      name: values.name,
    };

    try {
      await (data.value?.create
        ? apiClient.navigationService.Create({ data: { ...payload } as any })
        : apiClient.navigationService.Update({
            id: data.value?.row?.id,
            data: { ...payload } as any,
            updateMask: makeUpdateMask(Object.keys(payload)),
          }));

      notification.success({
        message: data.value?.create
          ? $t('ui.notification.create_success')
          : $t('ui.notification.update_success'),
      });
      drawerApi.close();
    } catch {
      notification.error({
        message: data.value?.create
          ? $t('ui.notification.create_failed')
          : $t('ui.notification.update_failed'),
      });
    } finally {
      setLoading(false);
    }
  },

  onOpenChange(isOpen) {
    if (!isOpen) {
      return;
    }

    data.value = drawerApi.getData<Record<string, any>>();
    const row = data.value?.row;

    if (row) {
      baseFormApi.setValues({
        isActive: row.isActive,
        locale: row.locale,
        location: row.location,
        name: row.name,
      });
    } else {
      baseFormApi.setValues({
        isActive: true,
      });
    }

    setLoading(false);
  },
});

function setLoading(loading: boolean) {
  drawerApi.setState({ confirmLoading: loading });
}
</script>
<template>
  <Drawer :title="getTitle" class="w-[720px]">
    <BaseForm />
  </Drawer>
</template>
