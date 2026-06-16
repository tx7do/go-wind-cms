<script lang="ts" setup>
import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { apiClient, makeUpdateMask, navigationItemLinkTypeList } from '#/api';

const data = ref<Record<string, any>>();

const getTitle = computed(() =>
  data.value?.create
    ? $t('ui.modal.create', {
        moduleName: $t('page.navigationItem.moduleName'),
      })
    : $t('ui.modal.update', {
        moduleName: $t('page.navigationItem.moduleName'),
      }),
);

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
      fieldName: 'title',
      label: $t('page.navigationItem.title'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'description',
      label: $t('page.navigationItem.description'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Input',
      fieldName: 'icon',
      label: $t('page.navigationItem.icon'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'url',
      label: $t('page.navigationItem.url'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Select',
      fieldName: 'linkType',
      label: $t('page.navigationItem.linkType'),
      defaultValue: 'LINK_TYPE_CUSTOM',
      componentProps: {
        options: navigationItemLinkTypeList,
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
      },
      rules: 'selectRequired',
    },
    {
      component: 'Switch',
      fieldName: 'isOpenNewTab',
      label: $t('page.navigationItem.isOpenNewTab'),
      defaultValue: false,
      componentProps: {
        class: 'w-auto',
      },
    },
    {
      component: 'Switch',
      fieldName: 'isInvalid',
      label: $t('page.navigationItem.isInvalid'),
      defaultValue: false,
      componentProps: {
        class: 'w-auto',
      },
    },
    {
      component: 'Textarea',
      fieldName: 'cssClass',
      label: $t('page.navigationItem.cssClass'),
      componentProps: {
        rows: 12,
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
    let items: any[] = [];

    if (values.itemsJson && String(values.itemsJson).trim()) {
      try {
        const parsed = JSON.parse(values.itemsJson);
        if (Array.isArray(parsed)) {
          items = parsed;
        } else {
          notification.error({
            message: $t('page.navigationItem.validation.itemsJsonInvalid'),
          });
          setLoading(false);
          return;
        }
      } catch {
        notification.error({
          message: $t('page.navigationItem.validation.itemsJsonInvalid'),
        });
        setLoading(false);
        return;
      }
    }

    const payload = {
      isActive: values.isActive,
      items,
      locale: values.locale,
      location: values.location,
      name: values.name,
    };

    try {
      await (data.value?.create
        ? apiClient.navigationItemService.Create({
            data: { ...payload } as any,
          })
        : apiClient.navigationItemService.Update({
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
      baseFormApi.setValues(data.value?.row);
    } else {
      baseFormApi.setValues({});
    }

    setLoading(false);
  },
});

function setLoading(loading: boolean) {
  drawerApi.setState({ loading });
}
</script>
<template>
  <Drawer :title="getTitle" class="w-[720px]">
    <BaseForm />
  </Drawer>
</template>
